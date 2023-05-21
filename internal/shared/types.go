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
	"github.com/gocql/gocql"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/ctrlplane/internal/shared/queue"
)

// workflow types.
type (
	Queue          = queue.Name
	WorkflowSignal string // WorkflowSignal is the name of a workflow signal.

	PullRequestSignal struct {
		RepoID           gocql.UUID
		SenderWorkflowID string
		TriggerID        int64
	}

	FutureHandler  func(workflow.Future)               // FutureHandler is the signature of the future handler function.
	ChannelHandler func(workflow.ReceiveChannel, bool) // ChannelHandler is the signature of the channel handler function.
)

const (
	CoreQueue      Queue = "core"      // core queue
	ProvidersQueue Queue = "providers" // messaging related to providers
	MutexQueue     Queue = "mutex"     // mutex workflow queue
)

// workflow signals.
const (
	WorkflowSignalTriggerDeployment WorkflowSignal = "deployment_trigger"
)

/*
 * Methods for WorkflowSignal.
 */
func (w WorkflowSignal) String() string { return string(w) }
