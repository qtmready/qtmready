package activities

import (
	"context"

	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// Activities groups all the activities for the slack provider.
	Activities struct{}
)

func (a *Activities) NotifyLinesExceed(
	ctx context.Context, event *events.Event[eventsv1.RepoHook, eventsv1.Diff],
) error {
	return nil
}
