package main

import (
	"context"
	"log"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Create a Temporal client with custom host:port
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Create a Temporal worker
	w := worker.New(c, "merge-queue-task-queue", worker.Options{})

	// Create an instance of MergeQueueWorkflows
	workflows := &MergeQueueWorkflows{}

	// Register the workflow and activity with the worker
	w.RegisterWorkflow(workflows.MergeQueueWorkflow)
	w.RegisterActivity(workflows.ProcessSignalActivity)

	// Start the worker
	err = w.Start()
	if err != nil {
		log.Println("Unable to start worker", err.Error())
	}

	// Signal handling to add tasks to the queue
	for {
		// Simulate receiving a merge queue signal
		signal := Signal{
			merge_queue_signal: &MergeQueueSignal{PullRequestID: 1, Branch: "feature1"},
			created_at:         time.Now(),
		}

		// Use SignalWithStartWorkflow to signal the workflow and start it if necessary
		options := client.StartWorkflowOptions{
			ID:                       "mergeQueueWorkflowID",
			TaskQueue:                "merge-queue-task-queue",
			WorkflowExecutionTimeout: time.Hour, // Set the desired timeout
		}

		// Signal the workflow and start it if it doesn't exist
		_, err = c.SignalWithStartWorkflow(
			context.Background(),
			"mergeQueueWorkflowID", // Workflow ID
			"mergeQueueSignal",     // Signal name
			signal,                 // Signal input
			options,                // Workflow start options
			workflows.MergeQueueWorkflow,
		)

		if err != nil {
			log.Println("Unable to signal workflow:", err.Error())
		}

		// Sleep to simulate interval between signals
		time.Sleep(10 * time.Second)
	}
}
