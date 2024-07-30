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
	workflowID        = "merge-queue-" + uuid.NewString()
	temporalTaskQueue = "merge-queue-task-queue"
)

// Register the workflow and activities with the worker.
func registerWorkflowsAndActivities(w worker.Worker) {
	w.RegisterWorkflow(QueueCtrl)
	w.RegisterActivity(BranchPop)
}

// TODO - refine.
func main() {
	// Create the Temporal client.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Create a worker for the task queue.
	w := worker.New(c, temporalTaskQueue, worker.Options{})
	registerWorkflowsAndActivities(w)

	// Start the worker.
	err = w.Start()
	if err != nil {
		log.Println("Unable to start worker", err)
	}
	defer w.Stop()

	// Start the workflow.
	options := client.StartWorkflowOptions{
		ID:                       workflowID,
		TaskQueue:                temporalTaskQueue,
		WorkflowExecutionTimeout: time.Hour, // Set the desired timeout
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, QueueCtrl)
	if err != nil {
		log.Println("Unable to execute workflow", err)
	}

	// Wait for the workflow to complete.
	var result string

	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Println("Unable to get workflow result", err)
	}

	log.Println("Succuesss")
}
