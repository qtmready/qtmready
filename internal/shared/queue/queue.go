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
	// queue defines the basic queue.
	queue struct {
		name                Name   // The name of the temporal queue.
		prefix              string // The prefix for the Workflow ID.
		workflowMaxAttempts int32  // The maximum number of attempts for a workflow.
	}
)

func (q Name) ToString() string {
	return string(q)
}

func (q *queue) Name() string {
	return q.name.ToString()
}

func (q *queue) Prefix() string {
	return q.prefix
}

func (q *queue) WorkflowID(options ...WorkflowIDOption) string {
	id := NewWorkflowIDCreator(options...)
	if id.IsChild() {
		return id.String(q)
	}

	return id.String(nil)
}

func (q *queue) CreateWorkflowOptions(options ...WorkflowIDOption) client.StartWorkflowOptions {
	return client.StartWorkflowOptions{
		ID:          q.WorkflowID(options...),
		TaskQueue:   q.Name(),
		RetryPolicy: &sdktemporal.RetryPolicy{MaximumAttempts: q.workflowMaxAttempts},
	}
}

func (q *queue) CreateChildWorkflowOptions(options ...WorkflowIDOption) workflow.ChildWorkflowOptions {
	return workflow.ChildWorkflowOptions{
		WorkflowID:  q.WorkflowID(options...),
		RetryPolicy: &sdktemporal.RetryPolicy{MaximumAttempts: q.workflowMaxAttempts},
	}
}

// WithName sets the queue name and the prefix for the workflow ID.
func WithName(name Name) QueueOption {
	return func(q Queue) {
		q.(*queue).name = name
		q.(*queue).prefix = DefaultPrefix + name.ToString()
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
