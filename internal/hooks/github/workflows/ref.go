package workflows

import (
	"github.com/google/uuid"
	"go.breu.io/durex/dispatch"
	"go.temporal.io/sdk/workflow"

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

	proto := cast.RefToProto(payload)
	meta := &defs.HydratedRepoEvent{}

	{
		hydratePayload := &defs.HydrateRepoEventPayload{
			RepoID:            payload.Repository.ID,
			InstallationID:    payload.Installation.ID,
			ShouldFetchParent: false,
		}
		if err := workflow.ExecuteActivity(ctx, acts.HydrateGithubRefEvent, hydratePayload).Get(ctx, meta); err != nil {
			return err
		}
	}

	scope := events.ScopeBranch

	if payload.RefType != "branch" {
		return nil
	}

	action := events.ActionCreated
	if !payload.IsCreated {
		action = events.ActionDeleted
	}

	evt := events.
		New[eventsv1.RepoHook, eventsv1.GitRef]().
		SetHook(eventsv1.RepoHook_REPO_HOOK_GITHUB).
		SetScope(scope).
		SetAction(action).
		SetSource(meta.Repo.Url).
		SetOrg(meta.Repo.OrgID).
		SetSubjectName(events.SubjectNameRepos).
		SetSubjectID(meta.Repo.ID).
		SetPayload(&proto)

	if meta.ParentID != uuid.Nil {
		evt.SetParents(meta.ParentID)
	}

	if meta.Team != nil {
		evt.SetTeam(meta.Team.ID)
	}

	if meta.User != nil {
		evt.SetUser(meta.User.ID)
	}

	if err := pulse.Persist(ctx, evt); err != nil {
		return err
	}

	hevent := &defs.HydratedQuantmEvent[eventsv1.GitRef]{Event: evt, Meta: meta}

	return workflow.ExecuteActivity(ctx, acts.SignalRepoWithGithubRef, hevent).Get(ctx, nil)
}
