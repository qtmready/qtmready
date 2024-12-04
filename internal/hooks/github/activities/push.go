package activities

import (
	"context"
	"time"

	"go.breu.io/quantm/internal/hooks/github/defs"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// Push groups all the activities required for the Github Push.
	Push struct{}
)

func (p *Push) HydrateGithubPushEvent(ctx context.Context, params *defs.HydrateRepoEventPayload) (*defs.HydratedRepoEvent, error) {
	time.Sleep(2 * time.Second) // FIXME: this is a quick hack to get the parent id.

	return HydrateRepoEvent(ctx, params)
}

func (p *Push) SignalRepoWithGithubPush(ctx context.Context, hydrated *defs.HydratedQuantmEvent[eventsv1.Push]) error {
	return SignalRepo(ctx, hydrated)
}
