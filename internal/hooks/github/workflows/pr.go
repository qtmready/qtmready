package workflows

import (
	"github.com/google/uuid"
	"go.breu.io/durex/dispatch"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/repos"
	"go.breu.io/quantm/internal/events"
	"go.breu.io/quantm/internal/hooks/github/activities"
	"go.breu.io/quantm/internal/hooks/github/cast"
	"go.breu.io/quantm/internal/hooks/github/defs"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
	"go.breu.io/quantm/internal/pulse"
)

// The Pr workflow processes GitHub webhook pull request events, converting the defs.Pr payload into a QuantmEvent.
// This involves hydrating the event with repository, installation, user, and team metadata.
// hydrated details and original payload, and finally signaling the repository.
func Pr(ctx workflow.Context, pr *defs.PR) error {
	acts := &activities.Pr{}
	ctx = dispatch.WithDefaultActivityContext(ctx)

	proto := cast.PrToProto(pr)
	hre := &defs.HydratedRepoEvent{} // hre -> hydrated repo event

	email := ""
	if pr.GetSenderEmail() != nil {
		email = *pr.GetSenderEmail()
	}

	{
		payload := &defs.HydrateRepoEventPayload{
			RepoID:         pr.GetRepositoryID(),
			InstallationID: pr.GetInstallationID(),
			Email:          email,
			Branch:         repos.BranchNameFromRef(pr.GetHeadBranch()),
		}
		if err := workflow.ExecuteActivity(ctx, acts.HydrateGithubPullRequestEvent, payload).Get(ctx, hre); err != nil {
			return err
		}
	}

	// handle actions
	event := events.
		New[eventsv1.RepoHook, eventsv1.PullRequest]().
		SetHook(eventsv1.RepoHook_REPO_HOOK_GITHUB).
		SetScope(events.ScopePr).
		SetAction(events.Action(pr.GetAction())).
		SetSource(hre.GetRepoUrl()).
		SetOrg(hre.GetOrgID()).
		SetSubjectName(events.SubjectNameRepos).
		SetSubjectID(hre.GetRepoID()).
		SetPayload(&proto)

	if hre.GetParentID() != uuid.Nil {
		event.SetParents(hre.GetParentID())
	}

	if hre.GetTeam() != nil {
		event.SetTeam(hre.GetTeamID())
	}

	if hre.GetUser() != nil {
		event.SetUser(hre.GetUserID())
	}

	if err := pulse.Persist(ctx, event); err != nil {
		return err
	}

	hevent := &defs.HydratedQuantmEvent[eventsv1.PullRequest]{Event: event, Meta: hre, Signal: repos.SignalPR}

	return workflow.ExecuteActivity(ctx, acts.SignalRepoWithGithubPR, hevent).Get(ctx, nil)
}
