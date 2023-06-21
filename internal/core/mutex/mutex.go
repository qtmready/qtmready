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

	"go.breu.io/ctrlplane/internal/shared"
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
	Release func() error

	// Mutex defines the signature for the workflow mutex. This workflow is meant to control the access to a resource.
	Mutex interface {
		Start() error                // Start the mutex workflow.
		Acquire() error              // Acquire aquires the lock.
		Release() error              // Release releases the lock.
		SetContext(workflow.Context) // SetContext sets the workflow context for the current mutex workflow exececution.
	}

	MutexOption func(*MUtex)

	Contexts struct {
		caller workflow.Context
		mutex  workflow.Context
	}

	// mutex is the implementation of the Mutex interface.
	//
	// Although it gets the job done for now, but it is not an ideal design. The mutex should hold the lock regardless of the caller.
	// We should be able to call the mutex from any workflow and it should be able to acquire the lock.
	MUtex struct {
		Contexts *Contexts
		Id       string        // ID of the mutex. The format is `{resource type}.{resource ID}`.
		Timeout  time.Duration // Timeout for the mutex. After this timeout, the lock is automagically released.
	}
)

func (m *MUtex) Start() error {
	if err := m.validate(); err != nil {
		shared.Logger().Error("unable to validate mutex", "error", err)
		return err
	}

	logger := workflow.GetLogger(m.Contexts.caller)
	opts := shared.Temporal().
		Queue(shared.CoreQueue).
		ChildWorkflowOptions(
			shared.WithWorkflowBlock("mutex"),
			shared.WithWorkflowBlockID(m.Id),
		)
		// GetChildWorkflowOptions("mutex", m.Id)
	ctx := workflow.WithChildOptions(m.Contexts.caller, opts)
	logger.Info("mutex: starting workflow ...", "resource ID", m.Id, "with timeout", m.Timeout)

	var exe workflow.Execution
	if err := workflow.
		ExecuteChildWorkflow(ctx, Workflow, m.Timeout).
		GetChildWorkflowExecution().
		Get(ctx, &exe); err != nil {
		logger.Error("mutex: unable to start.", "error", err)
		return NewStartWorkflowError(m.Contexts.caller)
	}

	m.SetContext(ctx)
	logger.Info(
		"mutex: workflow started, waiting for lock to be acquired.",
		"resource ID", m.Id,
		"workflow ID", exe.ID,
		"run ID", exe.RunID,
	)

	return nil
}

func (m *MUtex) Acquire() error {
	shared.Logger().Debug("Acquire lock", "m", m)
	logger := workflow.GetLogger(m.Contexts.caller)
	caller := workflow.GetInfo(m.Contexts.caller)
	mutex := workflow.GetInfo(m.Contexts.mutex)
	logger.Info(
		"mutex: acquiring lock. sending signal to mutex workflow ...",
		"resource ID", m.Id,
		"caller", caller.WorkflowType.Name,
		"caller ID", caller.WorkflowExecution.ID,
	)

	if err := workflow.
		SignalExternalWorkflow(
			m.Contexts.caller,
			mutex.WorkflowExecution.ID,
			"",
			WorkflowSignalAcquire.String(),
			caller.WorkflowExecution.ID,
		).
		Get(m.Contexts.caller, nil); err != nil {
		return NewAcquireLockError(m.Contexts.mutex)
	}

	logger.Info("mutex: acquiring lock. signal sent successfully, waiting for lock ... ", "resource ID", m.Id)
	workflow.GetSignalChannel(m.Contexts.caller, WorkflowSignalLocked.String()).Receive(m.Contexts.caller, nil)
	logger.Info("mutex: lock acquired.", "resource ID", m.Id)

	return nil
}

func (m *MUtex) Release() error {
	logger := workflow.GetLogger(m.Contexts.caller)
	logger.Info("mutex: releasing lock. sending signal to mutex workflow ...", "resource ID", m.Id)

	caller := workflow.GetInfo(m.Contexts.caller)
	mutex := workflow.GetInfo(m.Contexts.mutex)

	if err := workflow.SignalExternalWorkflow(
		m.Contexts.caller,
		mutex.WorkflowExecution.ID,
		"",
		WorkflowSignalRelease.String(),
		caller.WorkflowExecution.ID,
	).Get(m.Contexts.caller, nil); err != nil {
		return NewReleaseLockError(m.Contexts.mutex)
	}

	return nil
}

func (m *MUtex) SetContext(ctx workflow.Context) {
	m.Contexts.mutex = ctx
}

// validate validates if the mutex is properly configured.
func (m *MUtex) validate() error {
	if m.Contexts == nil {
		return ErrNilContext
	}

	if m.Id == "" {
		return ErrNoResourceID
	}

	return nil
}

// WithCallerContext sets the workflow context for the workflow that is invoking the mutex.
func WithCallerContext(ctx workflow.Context) MutexOption {
	return func(m *MUtex) {
		m.Contexts.caller = ctx
	}
}

// WithID sets the resource ID for the mutex workflow.
func WithID(id string) MutexOption {
	return func(m *MUtex) {
		m.Id = id
	}
}

// WithTimeout sets the timeout for the mutex workflow.
func WithTimeout(timeout time.Duration) MutexOption {
	return func(m *MUtex) {
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
func New(opts ...MutexOption) Mutex {
	m := &MUtex{Timeout: DefaultTimeout, Contexts: &Contexts{}}
	for _, opt := range opts {
		opt(m)
	}

	return m
}
