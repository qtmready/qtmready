package activities

import (
	"context"
	"time"

	"go.breu.io/quantm/internal/hooks/github/defs"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// Pr groups all the activities required for the Github Pull Request.
	Pr struct{}
)

func (pr *Pr) HydrateGithubPullRequestEvent(ctx context.Context, params *defs.HydrateRepoEventPayload) (*defs.HydratedRepoEvent, error) {
	time.Sleep(2 * time.Second) // FIXME: this is a quick hack to get the parent id.

	return HydrateRepoEvent(ctx, params)
}

func (p *Pr) SignalRepoWithGithubPR(ctx context.Context, hydrated *defs.HydratedQuantmEvent[eventsv1.PullRequest]) error {
	return SignalRepo(ctx, hydrated)
}
