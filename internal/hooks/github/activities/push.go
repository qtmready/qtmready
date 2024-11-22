package activities

import (
	"context"

	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/core/repos"
	"go.breu.io/quantm/internal/events"
	"go.breu.io/quantm/internal/hooks/github/cast"
	"go.breu.io/quantm/internal/hooks/github/defs"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// Push groups all the activities required for the Github Push.
	Push struct{}
)

func (p *Push) ConvertToPushEvent(
	ctx context.Context, payload *defs.Push,
) (*defs.RepoEvent[eventsv1.RepoHook, eventsv1.Push], error) {
	// Populate and set the quantum event
	params := &defs.RepoEventPayload{
		InstallationID: payload.InstallationID(),
		RepoID:         payload.RepoID(),
		Email:          payload.PusherEmail(),
		Action:         events.EventActionCreated,
		Scope:          events.EventScopePush,
	}

	resp, err := PopulateRepoEvent[eventsv1.RepoHook, eventsv1.Push](ctx, params)
	if err != nil {
		return nil, err
	}

	resp.Event.Payload = *cast.PushToProto(payload)

	return resp, nil
}

func (p *Push) SignalCoreRepo(
	ctx context.Context, repo *repos.HypdratedRepo, signal queues.Signal, payload any,
) error {
	return SignalCoreRepo(ctx, repo, signal, payload)
}
