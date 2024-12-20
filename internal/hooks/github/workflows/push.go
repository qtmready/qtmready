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

// The Push workflow processes GitHub webhook push events, converting the defs.Push payload into a QuantmEvent.
// This involves hydrating the event with repository, installation, user, and team metadata, determining the
// event action (create, delete, or force push), constructing and persisting a QuantmEvent encompassing the
// hydrated details and original payload, and finally signaling the repository.
func Push(ctx workflow.Context, push *defs.Push) error {
	acts := &activities.Push{}
	ctx = dispatch.WithDefaultActivityContext(ctx)

	proto := cast.PushToProto(push)
	hre := &defs.HydratedRepoEvent{} // hre -> hydrated repo event

	{
		payload := &defs.HydratedRepoEventPayload{
			RepoID:         push.GetRepositoryID(),
			InstallationID: push.GetInstallationID(),
			Email:          push.GetPusherEmail(),
			Branch:         repos.BranchNameFromRef(push.GetRef()),
		}
		if err := workflow.ExecuteActivity(ctx, acts.HydrateGithubPushEvent, payload).Get(ctx, hre); err != nil {
			return err
		}
	}

	action := events.ActionCreated

	if push.Deleted {
		action = events.ActionDeleted
	}

	if push.Forced {
		action = events.ActionForced
	}

	event := events.
		New[eventsv1.RepoHook, eventsv1.Push]().
		SetHook(eventsv1.RepoHook_REPO_HOOK_GITHUB).
		SetScope(events.ScopePush).
		SetAction(action).
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

	hevent := &defs.HydratedQuantmEvent[eventsv1.Push]{Event: event, Meta: hre, Signal: repos.SignalPush}

	return workflow.ExecuteActivity(ctx, acts.SignalRepoWithGithubPush, hevent).Get(ctx, nil)
}
