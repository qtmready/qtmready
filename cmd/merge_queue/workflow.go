package main

import (
	"errors"
	"time"

	"go.temporal.io/sdk/workflow"
)

// QueueCtrl is the main workflow.
func QueueCtrl(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Merge queue workflow started")

	queue := NewQueue()

	logger.Info("Queue init", "queue", queue)

	if queue == nil {
		logger.Error("Failed to create queue")
		return errors.New("failed to create queue")
	}

	logger.Info("Queue created", "queue", queue)

	// Simulate adding branches to the queue
	queue.push("branch1")
	queue.push("branch2")
	queue.push("branch3")

	logger.Info("Queue pushed branches", "branches", queue.Branches)

	opts := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
	ctx = workflow.WithActivityOptions(ctx, opts)

	// Execute the activity to process the queue
	err := workflow.ExecuteActivity(ctx, BranchPop, queue).Get(ctx, nil)
	if err != nil {
		return err
	}

	logger.Info("Merge queue workflow completed")

	return nil
}
