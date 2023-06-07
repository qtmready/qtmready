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
	}

	// QueueOption is the option for a queue.
	QueueOption func(Queue)

	// Queues is a map of queues.
	Queues map[Name]Queue

	// WorkflowOptions is the interface for creating a workflow id.
	WorkflowOptions interface {
		IsChild() bool            // IsChild returns true if the workflow id is a child workflow id.
		ParentWorkflowID() string // ParentWorkflowID returns the parent workflow id.
		Suffix() string           // Suffix santizes the suffix of the workflow id and then formats it as a string.
	}

	// WorkflowOptionProvider is the option for creating a workflow id.
	WorkflowOptionProvider func(WorkflowOptions)
)
