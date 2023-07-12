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

package shared

import (
	"github.com/go-playground/validator/v10"
	"github.com/gocql/gocql"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/shared/queue"
)

// workflow related shared types and contants.
type (
	WorkflowSignal string // WorkflowSignal is the name of a workflow signal.

	// PullRequestSignal is the sent to PR workflows to trigger a deployment.
	PullRequestSignal struct {
		RepoID           gocql.UUID
		SenderWorkflowID string
		TriggerID        int64
	}

	FutureHandler  func(workflow.Future)               // FutureHandler is the signature of the future handler function.
	ChannelHandler func(workflow.ReceiveChannel, bool) // ChannelHandler is the signature of the channel handler function.

	WorkflowOption = queue.WorkflowOptions
)

var (
	WithWorkflowParent    = queue.WithWorkflowParent    // Sets the parent workflow ID.
	WithWorkflowBlock     = queue.WithWorkflowBlock     // Sets the block name for the workflow ID.
	WithWorkflowBlockID   = queue.WithWorkflowBlockID   // Sets the block value for the workflow ID.
	WithWorkflowElement   = queue.WithWorkflowElement   // Sets the element name for the workflow ID.
	WithWorkflowElementID = queue.WithWorkflowElementID // Sets the element value for the workflow ID.
	WithWorkflowMod       = queue.WithWorkflowMod       // Sets the modifier for the workflow ID.
	WithWorkflowModID     = queue.WithWorkflowModID     // Sets the modifier value for the workflow ID.
	WithWorkflowProp      = queue.WithWorkflowProp      // Sets the property for the workflow ID.
	NewWorkflowOptions    = queue.NewWorkflowOptions    // Creates a new workflow ID. (see queue.NewWorkflowOptions)
)

// queue definitions.
const (
	CoreQueue      queue.Name = "core"      // core queue
	ProvidersQueue queue.Name = "providers" // messaging related to providers
	MutexQueue     queue.Name = "mutex"     // mutex workflow queue
)

// workflow signal definitions.
const (
	WorkflowSignalDeploymentStarted WorkflowSignal = "deployment_trigger"
)

/*
 * Methods for WorkflowSignal.
 */
func (w WorkflowSignal) String() string { return string(w) }

type (
	// EchoValidator is a wrapper for the instantiated validator.
	EchoValidator struct {
		Validator *validator.Validate
	}
)

func (ev *EchoValidator) Validate(i any) error {
	return ev.Validator.Struct(i)
}
