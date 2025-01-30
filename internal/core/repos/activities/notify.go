package activities

import (
	"context"
	"log/slog"

	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	Notify struct{}
)

func (n *Notify) LinesExceeded(ctx context.Context, evt *events.Event[eventsv1.ChatHook, eventsv1.Diff]) error {
	if err := kernel.Get().ChatHook(evt.Context.Hook).NotifyLinesExceed(ctx, evt); err != nil {
		slog.Warn("unable to notify on chat", "error", err.Error())
		return err
	}

	return nil
}

func (n *Notify) MergeConflict(ctx context.Context, evt *events.Event[eventsv1.ChatHook, eventsv1.Merge]) error {
	if err := kernel.Get().ChatHook(evt.Context.Hook).NotifyMergeConflict(ctx, evt); err != nil {
		slog.Warn("unable to notify on chat", "error", err.Error())
		return err
	}

	return nil
}
