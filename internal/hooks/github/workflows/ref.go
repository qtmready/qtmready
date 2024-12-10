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

// The Ref workflow processes GitHub webhook ref events, converting the defs.WebhookRef payload into a QuantmEvent.
// This involves hydrating the event with repository, installation, user, and team metadata, determining the
// event action (create or delete), constructing and persisting a QuantmEvent encompassing the hydrated details
// and original payload, and finally signaling the repository.
func Ref(ctx workflow.Context, payload *defs.WebhookRef, event defs.WebhookEvent) error {
	acts := &activities.Ref{}
	ctx = dispatch.WithDefaultActivityContext(ctx)
	logger := workflow.GetLogger(ctx)

	proto := cast.RefToProto(payload)
	meta := &defs.HydratedRepoEvent{}

	{
		payload := &defs.HydrateRepoEventPayload{
			RepoID:            payload.Repository.ID,
			InstallationID:    payload.Installation.ID,
			ShouldFetchParent: false,
		}
		if err := workflow.ExecuteActivity(ctx, acts.HydrateGithubRefEvent, payload).Get(ctx, meta); err != nil {
			return err
		}
	}

	signal := repos.SignalRef
	scope := events.ScopeBranch

	if payload.RefType != "branch" {
		logger.Warn("ref: unhandled ref event", "type", payload.RefType)
		return nil
	}

	action := events.ActionCreated
	if event == defs.WebhookEventDelete {
		action = events.ActionDeleted
	}

	evt := events.
		New[eventsv1.RepoHook, eventsv1.GitRef]().
		SetHook(eventsv1.RepoHook_REPO_HOOK_GITHUB).
		SetScope(scope).
		SetAction(action).
		SetSource(meta.GetRepoUrl()).
		SetOrg(meta.GetOrgID()).
		SetSubjectName(events.SubjectNameRepos).
		SetSubjectID(meta.GetRepoID()).
		SetPayload(&proto)

	if meta.GetParentID() != uuid.Nil {
		evt.SetParents(meta.GetParentID())
	}

	if meta.GetTeam() != nil {
		evt.SetTeam(meta.GetTeamID())
	}

	if meta.GetUser() != nil {
		evt.SetUser(meta.GetUserID())
	}

	if err := pulse.Persist(ctx, evt); err != nil {
		return err
	}

	hevent := &defs.HydratedQuantmEvent[eventsv1.GitRef]{Event: evt, Meta: meta, Signal: signal}

	return workflow.ExecuteActivity(ctx, acts.SignalRepoWithGithubRef, hevent).Get(ctx, nil)
}
