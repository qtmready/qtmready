package queue

import (
	"go.temporal.io/sdk/client"
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
		// CreateWorkflowID creates an idempotency key. Sometimes we need to signal the workflow from a completely disconnected part of the
		// application. For us, it is important to arrive at the same workflow ID regardless of the conditions.
		// We try to follow the block, element, modifier pattern popularized by advocates of mantainable CSS. For more info,
		// https://getbem.com.
		//
		// Example:
		// For the block github with installation id 123, the element being the repository with id 456, and the modifier being the
		// pull request with id 789, we would call
		//   GetWorkflowOptions("github", "123", "repository", "456", "pullrequest", "789")
		CreateWorkflowID(string, ...string) string

		// Name gets the name of the queue as string.
		Name() string

		// Prefix gets the prefix of the queue as string.
		Prefix() string

		// GetWorkflowOptions returns the workflow options for the queue.
		// GetWorkflowOptions takes the same arguments as CreateWorkflowID.
		GetWorkflowOptions(string, ...string) client.StartWorkflowOptions

		// GetChildWorkflowOptions gets the child workflow options.
		GetChildWorkflowOptions(string, ...string) workflow.ChildWorkflowOptions
	}

	// QueueOption is the option for a queue.
	QueueOption func(Queue)

	// Queues is a map of queues.
	Queues map[Name]Queue

	// WorkflowID is the interface for creating a workflow id.
	WorkflowID interface {
		IsChild() bool       // IsChild returns true if the workflow id is a child workflow id.
		String(Queue) string // String returns the workflow id as a string.
	}

	// WorkflowIDOption is the option for creating a workflow id.
	WorkflowIDOption func(WorkflowID)
)
