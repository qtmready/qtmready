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
// This involves hydrating the event with repository, installation, user, team metadata hydrated details and original payload, and
// finally signaling the repository.
func PullRequest(ctx workflow.Context, pr *defs.PR) error {
	acts := &activities.PullRequest{}
	hydrated := &defs.HydratedRepoEvent{}

	ctx = dispatch.WithDefaultActivityContext(ctx)

	email := ""
	if pr.GetSenderEmail() != nil {
		email = *pr.GetSenderEmail()
	}

	payload := &defs.HydratedRepoEventPayload{
		RepoID:         pr.GetRepositoryID(),
		InstallationID: pr.GetInstallationID(),
		Email:          email,
		Branch:         repos.BranchNameFromRef(pr.GetHeadBranch()),
	}

	if err := workflow.ExecuteActivity(ctx, acts.HydrateGithubPREvent, payload).Get(ctx, hydrated); err != nil {
		return err
	}

	if pr.Label != nil {
		return handle_label(ctx, pr, hydrated)
	}

	return handle_pr(ctx, pr, hydrated)
}

func handle_pr(ctx workflow.Context, pr *defs.PR, repo_evt *defs.HydratedRepoEvent) error {
	acts := &activities.PullRequest{}
	proto := cast.PullRequestToProto(pr)
	// handle actions
	event := events.
		New[eventsv1.RepoHook, eventsv1.PullRequest]().
		SetHook(eventsv1.RepoHook_REPO_HOOK_GITHUB).
		SetScope(events.ScopePr).
		SetAction(events.Action(pr.GetAction())). // TODO - handle the PR actions
		SetSource(repo_evt.GetRepoUrl()).
		SetOrg(repo_evt.GetOrgID()).
		SetSubjectName(events.SubjectNameRepos).
		SetSubjectID(repo_evt.GetRepoID()).
		SetPayload(&proto)

	if repo_evt.GetParentID() != uuid.Nil {
		event.SetParents(repo_evt.GetParentID())
	}

	if repo_evt.GetTeam() != nil {
		event.SetTeam(repo_evt.GetTeamID())
	}

	if repo_evt.GetUser() != nil {
		event.SetUser(repo_evt.GetUserID())
	}

	if err := pulse.Persist(ctx, event); err != nil {
		return err
	}

	hevent := &defs.HydratedQuantmEvent[eventsv1.PullRequest]{Event: event, Meta: repo_evt, Signal: repos.SignalPullRequest}

	return workflow.ExecuteActivity(ctx, acts.SignalRepoWithGithubPR, hevent).Get(ctx, nil)
}

func handle_label(ctx workflow.Context, pr *defs.PR, repo_evt *defs.HydratedRepoEvent) error {
	acts := &activities.PullRequest{}

	proto := cast.PullRequestLabelToProto(pr)
	if proto == nil {
		return nil
	}

	event := events.
		New[eventsv1.RepoHook, eventsv1.MergeQueue]().
		SetHook(eventsv1.RepoHook_REPO_HOOK_GITHUB).
		SetScope(events.ScopeMergeQueue).
		SetSource(repo_evt.GetRepoUrl()).
		SetOrg(repo_evt.GetOrgID()).
		SetSubjectName(events.SubjectNameRepos).
		SetSubjectID(repo_evt.GetRepoID()).
		SetPayload(proto)

	switch pr.GetAction() {
	case "labeled":
		event.SetAction(events.EventActionAdded)
	case "unlabeled":
		event.SetAction(events.EventActionRemoved)
	default:
		return nil
	}

	if repo_evt.GetParentID() != uuid.Nil {
		event.SetParents(repo_evt.GetParentID())
	}

	if repo_evt.GetTeam() != nil {
		event.SetTeam(repo_evt.GetTeamID())
	}

	if repo_evt.GetUser() != nil {
		event.SetUser(repo_evt.GetUserID())
	}

	if err := pulse.Persist(ctx, event); err != nil {
		return err
	}

	hevent := &defs.HydratedQuantmEvent[eventsv1.MergeQueue]{Event: event, Meta: repo_evt, Signal: repos.SignalMergeQueue}

	return workflow.ExecuteActivity(ctx, acts.SignalRepoWithGithubMergeQueue, hevent).Get(ctx, nil)
}
