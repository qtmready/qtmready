package queue

import (
	"fmt"

	"go.temporal.io/sdk/client"
	sdktemporal "go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

type (
	// queue defines the basic queue.
	queue struct {
		name                Name   // The name of the temporal queue.
		prefix              string // The prefix for the Workflow ID.
		workflowMaxAttempts int32  // The maximum number of attempts for a workflow.
	}
)

func (q Name) String() string {
	return string(q)
}

func (q *queue) Name() string {
	return q.name.String()
}

func (q *queue) Prefix() string {
	return q.prefix
}

func (q *queue) WorkflowID(options ...WorkflowOptionProvider) string {
	prefix := ""
	opts := NewWorkflowOptions(options...)

	if opts.IsChild() {
		prefix = opts.ParentWorkflowID()
	} else {
		prefix = q.prefix
	}

	return fmt.Sprintf("%s.%s", prefix, opts.Suffix())
}

func (q *queue) WorkflowOptions(options ...WorkflowOptionProvider) client.StartWorkflowOptions {
	return client.StartWorkflowOptions{
		ID:          q.WorkflowID(options...),
		TaskQueue:   q.Name(),
		RetryPolicy: &sdktemporal.RetryPolicy{MaximumAttempts: q.workflowMaxAttempts},
	}
}

func (q *queue) ChildWorkflowOptions(options ...WorkflowOptionProvider) workflow.ChildWorkflowOptions {
	return workflow.ChildWorkflowOptions{
		WorkflowID:  q.WorkflowID(options...),
		RetryPolicy: &sdktemporal.RetryPolicy{MaximumAttempts: q.workflowMaxAttempts},
	}
}

func (q *queue) Worker(c client.Client) worker.Worker {
	options := worker.Options{OnFatalError: func(err error) { panic(err) }, EnableSessionWorker: true}
	return worker.New(c, q.Name(), options)
}

// WithName sets the queue name and the prefix for the workflow ID.
func WithName(name Name) QueueOption {
	return func(q Queue) {
		q.(*queue).name = name
		q.(*queue).prefix = DefaultPrefix + name.String()
	}
}

// WithWorkflowMaxAttempts sets the maximum number of attempts for a workflow.
func WithWorkflowMaxAttempts(attempts int32) QueueOption {
	return func(q Queue) {
		q.(*queue).workflowMaxAttempts = attempts
	}
}

// NewQueue creates a new queue with the given options.
func NewQueue(opts ...QueueOption) Queue {
	q := &queue{workflowMaxAttempts: DefaultWorkflowMaxAttempts}
	for _, opt := range opts {
		opt(q)
	}

	return q
}
