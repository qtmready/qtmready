package main

import (
	"context"

	"go.temporal.io/sdk/activity"
)

// Activity method to process a signal.
func (w *MergeQueueWorkflows) ProcessSignalActivity(ctx context.Context, q Queue) error {
	activity.GetLogger(ctx).Info("Processing signal", "pull_request_id", q.pull_request_id)
	// Implement the processing logic for the signal
	return nil
}
