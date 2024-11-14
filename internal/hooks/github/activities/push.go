package githubacts

import (
	"context"

	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/events"
	githubcast "go.breu.io/quantm/internal/hooks/github/cast"
	githubdefs "go.breu.io/quantm/internal/hooks/github/defs"
	commonv1 "go.breu.io/quantm/internal/proto/ctrlplane/common/v1"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// Push groups all the activities required for the Github Push.
	Push struct{}
)

func (p *Push) ConvertToPushEvent(
	ctx context.Context, payload *githubdefs.Push,
) (*githubdefs.Eventory[commonv1.RepoHook, eventsv1.Push], error) {
	// Populate and set the quantum event
	params := &githubdefs.RepoEventPayload{
		InstallationID: payload.InstallationID(),
		RepoID:         payload.RepoID(),
		Action:         events.EventActionCreated,
		Scope:          events.EventScopePush,
	}

	resp, err := PopulateRepoEvent[commonv1.RepoHook, eventsv1.Push](ctx, params)
	if err != nil {
		return nil, err
	}

	resp.Event.Payload = *githubcast.PushToProto(payload)

	return resp, nil
}

func (p *Push) SignalCoreRepo(
	ctx context.Context, repo *entities.GetRepoRow, signal queues.Signal, payload any,
) error {
	return SignalCoreRepo(ctx, repo, signal, payload)
}
