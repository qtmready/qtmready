package workflows

import (
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
	meta := &defs.HydratedRepoEvent{}

	{
		payload := &defs.HydrateRepoEventPayload{
			RepoID:         push.Repository.ID,
			InstallationID: push.Installation.ID,
			Email:          push.Pusher.Email,
			Branch:         repos.BranchNameFromRef(push.Ref),
		}
		if err := workflow.ExecuteActivity(ctx, acts.HydratePushEvent, payload).Get(ctx, meta); err != nil {
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
		SetParents(meta.ParentID).
		SetScope(events.ScopePush).
		SetAction(action).
		SetSource(meta.Repo.Url).
		SetOrg(meta.Repo.OrgID).
		SetSubjectName(events.SubjectNameRepos).
		SetSubjectID(meta.Repo.ID).
		SetPayload(&proto)

	if meta.Team != nil {
		event.SetTeam(meta.Team.ID)
	}

	if meta.User != nil {
		event.SetUser(meta.User.ID)
	}

	if err := pulse.Persist(ctx, event); err != nil {
		return err
	}

	hevent := &defs.HydratedQuantmEvent[eventsv1.Push]{Event: event, Meta: meta}

	return workflow.ExecuteActivity(ctx, acts.SignalRepoWithGithubPush, hevent).Get(ctx, nil)
}
