package kernel

import (
	"context"

	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	Chat interface {
		// NotifyLinesExceed sends a message indicating a line exceed message.
		//
		// This method must not be called from the workflow.
		NotifyLinesExceed(ctx context.Context, event *events.Event[eventsv1.ChatHook, eventsv1.Diff]) error

		// NotifyMergeConflict sends a message indicating a merge conflict.
		//
		// This method must not be called from the workflow.
		NotifyMergeConflict(ctx context.Context, event *events.Event[eventsv1.ChatHook, eventsv1.Merge]) error
	}
)
