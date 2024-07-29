package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

var (
	workflow_id         = "merge-queue-" + uuid.NewString()
	temporal_task_queue = "merge-queue-task-queue"
	queue_signal_id     = "merge_queue_signal"
)

func main() {
	// Create a Temporal client with custom host:port
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Create a Temporal worker
	w := worker.New(c, temporal_task_queue, worker.Options{})

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
		q := Queue{
			pull_request_id: uuid.NewString(),
			created_at:      time.Now(),
		}

		// Use SignalWithStartWorkflow to signal the workflow and start it if necessary
		options := client.StartWorkflowOptions{
			ID:                       workflow_id,
			TaskQueue:                temporal_task_queue,
			WorkflowExecutionTimeout: time.Hour, // Set the desired timeout
		}

		// Signal the workflow and start it if it doesn't exist
		_, err = c.SignalWithStartWorkflow(
			context.Background(),
			workflow_id,     // Workflow ID
			queue_signal_id, // Signal name
			q,               // Signal input
			options,         // Workflow start options
			workflows.MergeQueueWorkflow,
		)

		if err != nil {
			log.Println("Unable to signal workflow:", err.Error())
		}

		// Sleep to simulate interval between signals
		time.Sleep(10 * time.Second)
	}
}
