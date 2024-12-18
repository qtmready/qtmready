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

func PullRequestReview(ctx workflow.Context, prr *defs.PrReview) error {
	acts := &activities.PullRequestReview{}
	ctx = dispatch.WithDefaultActivityContext(ctx)

	proto := cast.PrReviewToProto(prr)
	hydrated := &defs.HydratedRepoEvent{}

	email := ""
	if prr.GetSenderEmail() != nil {
		email = *prr.GetSenderEmail()
	}

	{
		payload := &defs.HydrateRepoEventPayload{
			RepoID:         prr.GetRepositoryID(),
			InstallationID: prr.GetInstallationID(),
			Email:          email,
			Branch:         repos.BranchNameFromRef(prr.GetHeadBranch()),
		}
		if err := workflow.ExecuteActivity(ctx, acts.HydrateGithubPREvent, payload).Get(ctx, hydrated); err != nil {
			return err
		}
	}

	// handle actions
	event := events.
		New[eventsv1.RepoHook, eventsv1.PullRequestReview]().
		SetHook(eventsv1.RepoHook_REPO_HOOK_GITHUB).
		SetScope(events.ScopePr).
		SetSource(hydrated.GetRepoUrl()).
		SetOrg(hydrated.GetOrgID()).
		SetSubjectName(events.SubjectNameRepos).
		SetSubjectID(hydrated.GetRepoID()).
		SetPayload(&proto)

	switch prr.Action {
	case "submitted":
		event.SetActionCreated()
	case "edited":
		event.SetActionUpdated()
	case "dismissed":
		event.SetActionDismissed()
	default:
		return nil
	}

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

	hevent := &defs.HydratedQuantmEvent[eventsv1.PullRequestReview]{Event: event, Meta: hydrated, Signal: repos.SignalPRReview}

	return workflow.ExecuteActivity(ctx, acts.SignalRepoWithGithubPR, hevent).Get(ctx, nil)
}
