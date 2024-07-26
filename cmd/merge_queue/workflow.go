package main

import (
	"log"
	"sort"
	"time"

	"go.temporal.io/sdk/workflow"
)

// Workflow method for managing the merge queue.
func (w *MergeQueueWorkflows) MergeQueueWorkflow(ctx workflow.Context) error {
	workflow.GetLogger(ctx).Info("Merge Queue Workflow started.")

	// Listen for signals to add tasks to the queue
	for {
		var signal Signal

		workflow.GetSignalChannel(ctx, "mergeQueueSignal").Receive(ctx, &signal)

		// Add the signal to the queue
		signal.priority = w.calculatePriority(signal)
		w.MergeQueue = append(w.MergeQueue, &signal)

		// Process the queue
		w.processQueue(ctx)
	}
}

func (w *MergeQueueWorkflows) calculatePriority(signal Signal) float64 {
	age := time.Since(signal.created_at).Seconds()
	return 1.0 / (1.0 + age) // Example: simple inverse age
}

func (w *MergeQueueWorkflows) processQueue(ctx workflow.Context) {
	// Sort the queue by priority (higher priority first)
	sort.SliceStable(w.MergeQueue, func(i, j int) bool {
		return w.MergeQueue[i].priority < w.MergeQueue[j].priority
	})

	opts := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
	ctx = workflow.WithActivityOptions(ctx, opts)

	// Process each task in the queue
	for len(w.MergeQueue) > 0 {
		task := w.MergeQueue[0]
		w.MergeQueue = w.MergeQueue[1:]

		err := workflow.ExecuteActivity(ctx, w.ProcessSignalActivity, task).Get(ctx, nil)
		if err != nil {
			log.Println("ProcessSignalActivity/error", err.Error())
		}
	}
}
