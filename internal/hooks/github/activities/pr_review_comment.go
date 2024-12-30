package activities

import (
	"context"

	"go.breu.io/quantm/internal/hooks/github/defs"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// PullRequestReviewComment groups all the activities required for the Github Pull Request review comment.
	PullRequestReviewComment struct{}
)

func (prr *PullRequestReviewComment) HydrateGithubPREvent(
	ctx context.Context, params *defs.HydratedRepoEventPayload,
) (*defs.HydratedRepoEvent, error) {
	return HydrateRepoEvent(ctx, params)
}

func (prr *PullRequestReviewComment) SignalRepoWithGithubPR(
	ctx context.Context, hydrated *defs.HydratedQuantmEvent[eventsv1.PullRequestReview],
) error {
	return SignalRepo(ctx, hydrated)
}
