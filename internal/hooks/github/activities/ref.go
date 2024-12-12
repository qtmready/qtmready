package activities

import (
	"context"

	"go.breu.io/quantm/internal/hooks/github/defs"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// Ref groups all the activities required for the Github Ref.
	Ref struct{}
)

// HydrateGithubRefEvent hydrates the branch event with the given parameters.
func (b *Ref) HydrateGithubRefEvent(ctx context.Context, params *defs.HydratedRepoEventPayload) (*defs.HydratedRepoEvent, error) {
	return HydrateRepoEvent(ctx, params)
}

// SignalRepoWithGithubRef signals the repository with the hydrated branch event.
func (b *Ref) SignalRepoWithGithubRef(ctx context.Context, hydrated *defs.HydratedQuantmEvent[eventsv1.GitRef]) error {
	return SignalRepo(ctx, hydrated)
}
