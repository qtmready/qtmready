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

package mutex

import (
	"errors"
	"fmt"

	"go.temporal.io/sdk/workflow"
)

var (
	ErrNilContext   = errors.New("contexts not initialized")
	ErrNoResourceID = errors.New("no resource ID provided")
)

type (
	acquireLockError struct {
		context workflow.Context // the workflow context for the mutex itself.
	}

	releaseLockError struct {
		context workflow.Context // the workflow context for the mutex itself.
	}

	startWorkflowError struct {
		context workflow.Context // the workflow contex for the workflow that is requesting to start the distributed mutex.
	}
)

func (e *acquireLockError) Error() string {
	return fmt.Sprintf("%s: failed to acquire lock.", workflow.GetInfo(e.context).WorkflowExecution.ID)
}

func (e *releaseLockError) Error() string {
	return fmt.Sprintf("%s: failed to release lock.", workflow.GetInfo(e.context).WorkflowExecution.ID)
}

func (e *startWorkflowError) Error() string {
	return fmt.Sprintf("%s: failed to start workflow.", workflow.GetInfo(e.context).WorkflowExecution.ID)
}

// NewAcquireLockError creates a new acquire lock error.
func NewAcquireLockError(ctx workflow.Context) error {
	return &acquireLockError{context: ctx}
}

// NewReleaseLockError creates a new release lock error.
func NewReleaseLockError(ctx workflow.Context) error {
	return &releaseLockError{context: ctx}
}

// NewStartWorkflowError creates a new start workflow error.
func NewStartWorkflowError(ctx workflow.Context) error {
	return &startWorkflowError{context: ctx}
}
