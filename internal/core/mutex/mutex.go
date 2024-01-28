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
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/shared"
)

const (
	DefaultTimeout = 30 * time.Minute
)

const (
	WorkflowSignalLocked  shared.WorkflowSignal = "locked"
	WorkflowSignalAcquire shared.WorkflowSignal = "acquire"
	WorkflowSignalRelease shared.WorkflowSignal = "release"
)

type (
	// Mutex defines the signature for the workflow mutex. This workflow is meant to control the access to a resource.
	Mutex interface {
		Start(ctx workflow.Context) error   // Start the mutex workflow.
		Acquire(ctx workflow.Context) error // Acquire aquires the lock.
		Release(ctx workflow.Context) error // Release releases the lock.
		// SetContext(workflow.Context)        // SetContext sets the workflow context for the current mutex workflow exececution.
	}

	Option func(*Lock)

	Contexts struct {
		caller workflow.Context
		mutex  workflow.Context
	}

	// Lock is the implementation of the Mutex interface.
	//
	// FIXME: Although it gets the job done for now, but it is not an ideal design. The mutex should hold the lock regardless of the caller.
	// We should be able to call the mutex from any workflow and it should be able to acquire the lock.
	Lock struct {
		contexts    *Contexts
		ID          string        // ID of the mutex. The format is `{resource type}.{resource ID}`.
		Timeout     time.Duration // Timeout for the mutex. After this timeout, the lock is automagically released.
		ExecutionID string
	}
)

func (m *Lock) Start(ctx workflow.Context) error {
	if err := m.validate(); err != nil {
		shared.Logger().Error("unable to validate mutex", "error", err)
		return err
	}

	logger := workflow.GetLogger(ctx)
	opts := shared.Temporal().
		Queue(shared.CoreQueue).
		ChildWorkflowOptions(
			shared.WithWorkflowBlock("mutex"),
			shared.WithWorkflowBlockID(m.ID),
		)
	// GetChildWorkflowOptions("mutex", m.Id)
	cctx := workflow.WithChildOptions(ctx, opts)

	logger.Info("mutex: starting workflow ...", "resource ID", m.ID, "with timeout", m.Timeout)

	var exe workflow.Execution
	if err := workflow.
		ExecuteChildWorkflow(cctx, Workflow, m.Timeout).
		GetChildWorkflowExecution().
		Get(cctx, &exe); err != nil {
		logger.Error("mutex: unable to start.", "error", err)
		return NewStartWorkflowError(m.ExecutionID)
	}

	m.ExecutionID = exe.ID

	logger.Info(
		"mutex: workflow started, waiting for lock to be acquired.",
		"resource ID", m.ID,
		"workflow ID", exe.ID,
		"run ID", exe.RunID,
	)

	return nil
}

func (m *Lock) Acquire(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)
	caller := workflow.GetInfo(ctx)
	logger.Info(
		"mutex: acquiring lock. sending signal to mutex workflow ...",
		"resource ID", m.ID,
		"caller", caller.WorkflowType.Name,
		"caller ID", caller.WorkflowExecution.ID,
	)

	if err := workflow.
		SignalExternalWorkflow(
			ctx,
			m.ExecutionID,
			"",
			WorkflowSignalAcquire.String(),
			caller.WorkflowExecution.ID,
		).
		Get(ctx, nil); err != nil {
		return NewAcquireLockError(m.ExecutionID)
	}

	logger.Info("mutex: acquiring lock. signal sent successfully, waiting for lock ... ", "resource ID", m.ID)
	workflow.GetSignalChannel(ctx, WorkflowSignalLocked.String()).Receive(ctx, nil)
	logger.Info("mutex: lock acquired.", "resource ID", m.ID)

	return nil
}

func (m *Lock) Release(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("mutex: releasing lock. sending signal to mutex workflow ...", "resource ID", m.ID)

	caller := workflow.GetInfo(ctx)

	if err := workflow.SignalExternalWorkflow(
		ctx,
		m.ExecutionID,
		"",
		WorkflowSignalRelease.String(),
		caller.WorkflowExecution.ID,
	).Get(ctx, nil); err != nil {
		return NewReleaseLockError(m.ExecutionID)
	}

	return nil
}

// validate validates if the mutex is properly configured.
func (m *Lock) validate() error {
	if m.contexts == nil {
		return ErrNilContext
	}

	if m.ID == "" {
		return ErrNoResourceID
	}

	return nil
}

// WithCallerContext sets the workflow context for the workflow that is invoking the mutex.
func WithCallerContext(ctx workflow.Context) Option {
	return func(m *Lock) {
		m.contexts.caller = ctx
	}
}

// WithID sets the resource ID for the mutex workflow.
func WithID(id string) Option {
	return func(m *Lock) {
		m.ID = id
	}
}

// WithTimeout sets the timeout for the mutex workflow.
func WithTimeout(timeout time.Duration) Option {
	return func(m *Lock) {
		m.Timeout = timeout
	}
}

// New returns a new Mutex.
// it should always be called with at least WithCallerContext and WithResource.
// If WithTimeout is not called, it defaults to DefaultTimeout.
//
// Example:
//
//	m := mutex.New(
//	  mutex.WithCallerContext(ctx),
//	  mutex.WithID("resource-type.resource-id"),
//	  mutex.WithTimeout(30*time.Minute),
//	)
//	if err := m.Start(); err != nil {/*handle error*/}
//	if err := m.Acquire(); err != nil {/*handle error*/}
//	if err := m.Release(); err != nil {/*handle error*/}
func New(opts ...Option) Mutex {
	m := &Lock{Timeout: DefaultTimeout, contexts: &Contexts{}}
	for _, opt := range opts {
		opt(m)
	}

	return m
}
