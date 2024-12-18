package activities

import (
	"context"

	"go.breu.io/quantm/internal/hooks/github/defs"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// PullRequestReview groups all the activities required for the Github Pull Request review.
	PullRequestReview struct{}
)

func (prr *PullRequestReview) HydrateGithubPREvent(
	ctx context.Context, params *defs.HydrateRepoEventPayload,
) (*defs.HydratedRepoEvent, error) {
	return HydrateRepoEvent(ctx, params)
}

func (prr *PullRequestReview) SignalRepoWithGithubPR(
	ctx context.Context, hydrated *defs.HydratedQuantmEvent[eventsv1.PullRequestReview],
) error {
	return SignalRepo(ctx, hydrated)
}
