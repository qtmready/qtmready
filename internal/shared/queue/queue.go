// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package queue

import (
	"go.temporal.io/sdk/client"
	sdktemporal "go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type (
	// Name is the name of the queue.
	Name string

	// Queue is the interface for a queue.
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

		// GetWorkflowOptions returns the workflow options for the queue.
		// GetWorkflowOptions takes the same arguments as CreateWorkflowID.
		GetWorkflowOptions(string, ...string) client.StartWorkflowOptions

		// GetName gets the name of the queue as string.
		GetName() string

		// GetChildWorkflowOptions gets the child workflow options.
		GetChildWorkflowOptions(string, ...string) workflow.ChildWorkflowOptions
	}

	// QueueOption is the option for a queue.
	QueueOption func(Queue)

	// Queues is a map of queues.
	Queues map[Name]Queue

	// queue defines the basic queue.
	queue struct {
		Name                Name   // The name of the queue.
		Prefix              string // The prefix for the Workflow ID.
		WorkflowMaxAttempts int32  // The maximum number of attempts for a workflow.
	}
)

const (
	DefaultPrefix              string = "ai.ctrlplane." // Prefix used to prepend the workflow id.
	DefaultWorkflowMaxAttempts int32  = 10              //
)

func (q Name) ToString() string {
	return string(q)
}

// GetName gets the name as string from Name.
func (q *queue) GetName() string {
	return q.Name.ToString()
}

// CreateWorkflowID creates the unique workflow ID from the workflow sender and appropriate arguments.
func (q *queue) CreateWorkflowID(sender string, args ...string) string {
	prepend := format(q.Prefix, sender)
	return format(append([]string{prepend}, args...)...)
}

// GetWorkflowOptions returns the workflow options for the queue.
func (q *queue) GetWorkflowOptions(sender string, args ...string) client.StartWorkflowOptions {
	opts := client.StartWorkflowOptions{
		ID:        q.CreateWorkflowID(sender, args...),
		TaskQueue: q.GetName(),
		// WorkflowIDReusePolicy: 3, // client.WorkflowIDReusePolicyRejectDuplicate
	}
	retryPolicy := &sdktemporal.RetryPolicy{MaximumAttempts: q.WorkflowMaxAttempts}
	opts.RetryPolicy = retryPolicy

	return opts
}

// GetChildWorkflowOptions gets the child workflow options.
func (q *queue) GetChildWorkflowOptions(sender string, args ...string) workflow.ChildWorkflowOptions {
	return workflow.ChildWorkflowOptions{
		WorkflowID:  q.CreateWorkflowID(sender, args...),
		RetryPolicy: &sdktemporal.RetryPolicy{MaximumAttempts: q.WorkflowMaxAttempts},
	}
}

// WithName sets the queue name and the prefix for the workflow ID.
func WithName(name Name) QueueOption {
	return func(q Queue) {
		q.(*queue).Name = name
		q.(*queue).Prefix = DefaultPrefix + name.ToString()
	}
}

// WithWorkflowMaxAttempts sets the maximum number of attempts for a workflow.
func WithWorkflowMaxAttempts(attempts int32) QueueOption {
	return func(q Queue) {
		q.(*queue).WorkflowMaxAttempts = attempts
	}
}

// NewQueue creates a new queue with the given options.
func NewQueue(opts ...QueueOption) Queue {
	q := &queue{WorkflowMaxAttempts: DefaultWorkflowMaxAttempts}
	for _, opt := range opts {
		opt(q)
	}

	return q
}
