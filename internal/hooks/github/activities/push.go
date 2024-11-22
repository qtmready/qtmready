package activities

import (
	"context"

	"go.breu.io/quantm/internal/hooks/github/defs"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// Push groups all the activities required for the Github Push.
	Push struct{}
)

func (p *Push) HydratePushEvent(ctx context.Context, params *defs.HydrateRepoEventPayload) (*defs.HydratedRepoEvent, error) {
	return HydrateRepoEvent(ctx, params)
}

func (p *Push) SignalGithubPushEventToRepo(ctx context.Context, hydrated *defs.HydratedQuantmEvent[eventsv1.Push]) error {
	return SignalToRepo(ctx, hydrated)
}
