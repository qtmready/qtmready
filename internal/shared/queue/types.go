package queue

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

const (
	DefaultPrefix              string = "ai.ctrlplane." // Prefix used to prepend the workflow id.
	DefaultWorkflowMaxAttempts int32  = 10              // Default workflow max Attempts.
)

type (
	// Name is the name of the queue.
	Name string

	// Queue defines the common interface for utilizing the Temporal queue.
	Queue interface {
		// Name gets the name of the queue as string.
		Name() string

		// Prefix gets the prefix of the queue as string.
		Prefix() string

		// WorkflowID gets the workflow id given the options. In most cases, the workflow id is called via GetWorkflowOptions
		// or GetChildWorkflowOptions. However, when we need to signal a workflow, this method comes in handy.
		WorkflowID(options ...WorkflowOptionProvider) string

		// WorkflowOptions creates the workflow options for the queue given WorkflowIDOptions.
		WorkflowOptions(options ...WorkflowOptionProvider) client.StartWorkflowOptions

		// ChildWorkflowOptions creates the child workflow options for the queue given WorkflowIDOptions.
		ChildWorkflowOptions(options ...WorkflowOptionProvider) workflow.ChildWorkflowOptions

		// Worker creates a worker for the queue.
		Worker(client client.Client) worker.Worker
	}

	// QueueOption is the option for a queue.
	QueueOption func(Queue)

	// Queues is a map of queues.
	Queues map[Name]Queue

	// Workers is a map of workers against queues.
	Workers map[Name]worker.Worker

	// WorkflowOptions defines the interface for creating workflow options.
	WorkflowOptions interface {
		IsChild() bool            // IsChild returns true if the workflow id is a child workflow id.
		ParentWorkflowID() string // ParentWorkflowID returns the parent workflow id.
		Suffix() string           // Suffix santizes the suffix of the workflow id and then formats it as a string.
	}

	// WorkflowOptionProvider provides the option to create a workflow id.
	WorkflowOptionProvider func(WorkflowOptions)
)
