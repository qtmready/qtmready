package pulse

import (
	"context"

	"go.breu.io/durex/dispatch"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

func Persist[H events.Hook, P events.Payload](ctx workflow.Context, event *events.Event[H, P]) error {
	ctx = dispatch.WithDefaultActivityContext(ctx)
	flat := event.Flatten()

	var future workflow.Future

	switch any(flat.Hook).(type) {
	case eventsv1.RepoHook:
		future = workflow.ExecuteActivity(ctx, PersistRepoEvent, flat)
	case eventsv1.MessagingHook:
		future = workflow.ExecuteActivity(ctx, PersistMessagingEvent, flat)
	}

	return future.Get(ctx, nil)
}

func PersistRepoEvent(ctx context.Context, flat events.Flat[eventsv1.RepoHook]) error { return nil }

func PersistMessagingEvent(ctx context.Context, flat events.Flat[eventsv1.MessagingHook]) error {
	return nil
}
