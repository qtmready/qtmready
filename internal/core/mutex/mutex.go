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

	MutexOption func(*mutex)

	Contexts struct {
		caller workflow.Context
		mutex  workflow.Context
	}

	// mutex is the implementation of the Mutex interface.
	//
	// Although it gets the job done for now, but it is not an ideal design. The mutex should hold the lock regardless of the caller.
	// We should be able to call the mutex from any workflow and it should be able to acquire the lock.
	mutex struct {
		contexts *Contexts
		id       string        // ID of the mutex. The format is `{resource type}.{resource ID}`.
		timeout  time.Duration // Timeout for the mutex. After this timeout, the lock is automagically released.
	}
)

func (m *mutex) Start() error {
	if err := m.validate(); err != nil {
		shared.Logger().Error("unable to validate mutex", "error", err)
		return err
	}

	logger := workflow.GetLogger(m.contexts.caller)
	opts := shared.Temporal().
		Queue(shared.CoreQueue).
		GetChildWorkflowOptions("mutex", m.id)
	ctx := workflow.WithChildOptions(m.contexts.caller, opts)
	logger.Info("mutex: starting workflow ...", "resource ID", m.id, "with timeout", m.timeout)

	var exe workflow.Execution
	if err := workflow.
		ExecuteChildWorkflow(ctx, Workflow, m.timeout).
		GetChildWorkflowExecution().
		Get(ctx, &exe); err != nil {
		logger.Error("mutex: unable to start.", "error", err)
		return NewStartWorkflowError(m.contexts.caller)
	}

	m.SetContext(ctx)
	logger.Info(
		"mutex: workflow started, waiting for lock to be acquired.",
		"resource ID", m.id,
		"workflow ID", exe.ID,
		"run ID", exe.RunID,
	)

	return nil
}

func (m *mutex) Acquire() error {
	logger := workflow.GetLogger(m.contexts.caller)
	caller := workflow.GetInfo(m.contexts.caller)
	mutex := workflow.GetInfo(m.contexts.mutex)
	logger.Info(
		"mutex: acquiring lock. sending signal to mutex workflow ...",
		"resource ID", m.id,
		"caller", caller.WorkflowType.Name,
		"caller ID", caller.WorkflowExecution.ID,
	)

	if err := workflow.
		SignalExternalWorkflow(
			m.contexts.caller,
			mutex.WorkflowExecution.ID,
			"",
			WorkflowSignalAcquire.String(),
			caller.WorkflowExecution.ID,
		).
		Get(m.contexts.caller, nil); err != nil {
		return NewAcquireLockError(m.contexts.mutex)
	}

	logger.Info("mutex: acquiring lock. signal sent successfully, waiting for lock ... ", "resource ID", m.id)
	workflow.GetSignalChannel(m.contexts.caller, WorkflowSignalLocked.String()).Receive(m.contexts.caller, nil)
	logger.Info("mutex: lock acquired.", "resource ID", m.id)

	return nil
}

func (m *mutex) Release() error {
	logger := workflow.GetLogger(m.contexts.caller)
	logger.Info("mutex: releasing lock. sending signal to mutex workflow ...", "resource ID", m.id)

	caller := workflow.GetInfo(m.contexts.caller)
	mutex := workflow.GetInfo(m.contexts.mutex)

	if err := workflow.SignalExternalWorkflow(
		m.contexts.caller,
		mutex.WorkflowExecution.ID,
		"",
		WorkflowSignalRelease.String(),
		caller.WorkflowExecution.ID,
	).Get(m.contexts.caller, nil); err != nil {
		return NewReleaseLockError(m.contexts.mutex)
	}

	return nil
}

func (m *mutex) SetContext(ctx workflow.Context) {
	m.contexts.mutex = ctx
}

// validate validates if the mutex is properly configured.
func (m *mutex) validate() error {
	if m.contexts == nil {
		return ErrNilContext
	}

	if m.id == "" {
		return ErrNoResourceID
	}

	return nil
}

// WithCallerContext sets the workflow context for the workflow that is invoking the mutex.
func WithCallerContext(ctx workflow.Context) MutexOption {
	return func(m *mutex) {
		m.contexts.caller = ctx
	}
}

// WithID sets the resource ID for the mutex workflow.
func WithID(id string) MutexOption {
	return func(m *mutex) {
		m.id = id
	}
}

// WithTimeout sets the timeout for the mutex workflow.
func WithTimeout(timeout time.Duration) MutexOption {
	return func(m *mutex) {
		m.timeout = timeout
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
	m := &mutex{timeout: DefaultTimeout, contexts: &Contexts{}}
	for _, opt := range opts {
		opt(m)
	}

	return m
}
