package activities

import (
	"context"

	"go.breu.io/quantm/internal/hooks/github/defs"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// PullRequest groups all the activities required for the Github Pull Request.
	PullRequest struct{}
)

func (pr *PullRequest) HydrateGithubPREvent(ctx context.Context, params *defs.HydratedRepoEventPayload) (*defs.HydratedRepoEvent, error) {
	return HydrateRepoEvent(ctx, params)
}

func (p *PullRequest) SignalRepoWithGithubPR(ctx context.Context, hydrated *defs.HydratedQuantmEvent[eventsv1.PullRequest]) error {
	return SignalRepo(ctx, hydrated)
}

func (p *PullRequest) SignalRepoWithGithubMergeQueue(ctx context.Context, hydrated *defs.HydratedQuantmEvent[eventsv1.MergeQueue]) error {
	return SignalRepo(ctx, hydrated)
}
