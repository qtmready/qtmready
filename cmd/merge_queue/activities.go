package main

import (
	"context"

	"go.temporal.io/sdk/activity"
)

// Activity method to process a signal.
func (w *MergeQueueWorkflows) ProcessSignalActivity(ctx context.Context, signal Signal) error {
	activity.GetLogger(ctx).Info("Processing signal", "Branch", signal.merge_queue_signal.Branch)
	// Implement the processing logic for the signal
	return nil
}
