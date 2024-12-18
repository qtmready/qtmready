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

// The PullRequest workflow processes GitHub webhook pull request events, converting the defs.PullRequest payload into a QuantmEvent.
// This involves hydrating the event with repository, installation, user, and team metadata.
// hydrated details and original payload, and finally signaling the repository.
func PullRequest(ctx workflow.Context, pr *defs.PR) error {
	acts := &activities.PullRequest{}
	ctx = dispatch.WithDefaultActivityContext(ctx)

	proto := cast.PrToProto(pr)
	hydrated := &defs.HydratedRepoEvent{} // hre -> hydrated repo event

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
		if err := workflow.ExecuteActivity(ctx, acts.HydrateGithubPREvent, payload).Get(ctx, hydrated); err != nil {
			return err
		}
	}

	// handle actions
	event := events.
		New[eventsv1.RepoHook, eventsv1.PullRequest]().
		SetHook(eventsv1.RepoHook_REPO_HOOK_GITHUB).
		SetScope(events.ScopePr).
		SetAction(events.Action(pr.GetAction())). // TODO - handle the PR actions
		SetSource(hydrated.GetRepoUrl()).
		SetOrg(hydrated.GetOrgID()).
		SetSubjectName(events.SubjectNameRepos).
		SetSubjectID(hydrated.GetRepoID()).
		SetPayload(&proto)

	if hydrated.GetParentID() != uuid.Nil {
		event.SetParents(hydrated.GetParentID())
	}

	if hydrated.GetTeam() != nil {
		event.SetTeam(hydrated.GetTeamID())
	}

	if hydrated.GetUser() != nil {
		event.SetUser(hydrated.GetUserID())
	}

	if err := pulse.Persist(ctx, event); err != nil {
		return err
	}

	hevent := &defs.HydratedQuantmEvent[eventsv1.PullRequest]{Event: event, Meta: hydrated, Signal: repos.SignalPR}

	return workflow.ExecuteActivity(ctx, acts.SignalRepoWithGithubPR, hevent).Get(ctx, nil)
}

func PullRequestLabel(ctx workflow.Context, pr *defs.PR) error {
	acts := &activities.PullRequest{}
	ctx = dispatch.WithDefaultActivityContext(ctx)

	proto := cast.PrLabelToProto(pr)
	hydrated := &defs.HydratedRepoEvent{} // hre -> hydrated repo event

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
		if err := workflow.ExecuteActivity(ctx, acts.HydrateGithubPREvent, payload).Get(ctx, hydrated); err != nil {
			return err
		}
	}

	// handle actions
	event := events.
		New[eventsv1.RepoHook, eventsv1.PullRequestLabel]().
		SetHook(eventsv1.RepoHook_REPO_HOOK_GITHUB).
		SetScope(events.ScopePrLabel).
		SetAction(events.Action(pr.GetAction())). // TODO - handle the PR actions
		SetSource(hydrated.GetRepoUrl()).
		SetOrg(hydrated.GetOrgID()).
		SetSubjectName(events.SubjectNameRepos).
		SetSubjectID(hydrated.GetRepoID()).
		SetPayload(&proto)

	if hydrated.GetParentID() != uuid.Nil {
		event.SetParents(hydrated.GetParentID())
	}

	if hydrated.GetTeam() != nil {
		event.SetTeam(hydrated.GetTeamID())
	}

	if hydrated.GetUser() != nil {
		event.SetUser(hydrated.GetUserID())
	}

	if err := pulse.Persist(ctx, event); err != nil {
		return err
	}

	hevent := &defs.HydratedQuantmEvent[eventsv1.PullRequestLabel]{Event: event, Meta: hydrated, Signal: repos.SignalPRLabel}

	return workflow.ExecuteActivity(ctx, acts.SignalRepoWithGithubPR, hevent).Get(ctx, nil)
}
