package main

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
)

// BranchPop processes the branches in the queue.
func BranchPop(ctx context.Context, q *Queue) error {
	activity.GetLogger(ctx).Info("Processing signal", "branches", q.Branches)

	for !q.is_empty() {
		branch := q.pop()
		if branch != nil {
			activity.GetLogger(ctx).Info("Processing branch", "branch", *branch)

			// Simulate processing the branch
			time.Sleep(time.Second)
		}
	}

	return nil
}
