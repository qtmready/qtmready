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
	"log/slog"
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/shared"
)

const (
	DefaultTimeout = 0 * time.Minute // DefaultTimeout is the default timeout for the mutex.
)

const (
	WorkflowSignalPrepare        shared.WorkflowSignal = "mutex__prepare"
	WorkflowSignalAcquire        shared.WorkflowSignal = "mutex__acquire"
	WorkflowSignalLocked         shared.WorkflowSignal = "mutex__locked"
	WorkflowSignalRelease        shared.WorkflowSignal = "mutex__release"
	WorkflowSignalReleased       shared.WorkflowSignal = "mutex__released"
	WorkflowSignalCleanup        shared.WorkflowSignal = "mutex__cleanup"
	WorkflowSignalCleanupDone    shared.WorkflowSignal = "mutex__cleanup_done"
	WorkflowSignalCleanupDoneAck shared.WorkflowSignal = "mutex__cleanup_done_ack"
	WorkflowSignalShutDown       shared.WorkflowSignal = "mutex__shutdown"
)

type (
	// Mutex defines the signature for the workflow mutex. This workflow is meant to control the access to a resource.
	Mutex interface {
		Prepare(ctx workflow.Context) error // Prepares the mutex for use.
		Acquire(ctx workflow.Context) error // Acquire aquires the lock.
		Release(ctx workflow.Context) error // Release releases the lock.
		Cleanup(ctx workflow.Context) error // Cleanup attempts to shutdown the mutex workflow, if it is no longer needed.
	}

	Option func(*Handler)

	// Handler is the Mutex handler.
	Handler struct {
		ResourceID string              `json:"resource_id"` // ResourceID identifies the resource being locked.
		Info       *workflow.Info      `json:"info"`        // Info holds the workflow info that requests the mutex.
		Execution  *workflow.Execution `json:"execution"`   // Info holds the workflow info that holds the mutex.
		Timeout    time.Duration       `json:"timeout"`     // Timeout sets the timeout, after which the lock is automatically released.
	}
)

func (info *Handler) Prepare(ctx workflow.Context) error {
	if err := info.validate(); err != nil {
		wferr(ctx, info, "mutex handler: unable to validate mutex", err)
		return err
	}

	wfinfo(ctx, info, "mutex handler: preparing ...")

	opts := workflow.ActivityOptions{StartToCloseTimeout: info.Timeout}
	ctx = workflow.WithActivityOptions(ctx, opts)

	exe := &workflow.Execution{}
	if err := workflow.ExecuteActivity(ctx, PrepareMutexActivity, info).Get(ctx, exe); err != nil {
		wfwarn(ctx, info, "mutex handler: unable to prepare mutex", err)
		return NewPrepareMutexError(info.ResourceID)
	}

	info.Execution = exe

	wfinfo(ctx, info, "mutex handler: prepared!", slog.String("mutex", info.Execution.ID))

	return nil
}

func (info *Handler) Acquire(ctx workflow.Context) error {
	wfinfo(ctx, info, "mutex handler: requesting lock ...")

	ok := true

	if err := workflow.
		SignalExternalWorkflow(ctx, info.Execution.ID, "", WorkflowSignalAcquire.String(), info).
		Get(ctx, nil); err != nil {
		wfwarn(ctx, info, "mutex handler: unable to request lock ...", err)
		return NewAcquireLockError(info.ResourceID)
	}

	wfinfo(ctx, info, "mutex handler: waiting for lock ...")
	workflow.GetSignalChannel(ctx, WorkflowSignalLocked.String()).Receive(ctx, &ok)
	wfinfo(ctx, info, "mutex handler: lock acquired!")

	if ok {
		return nil
	}

	return NewAcquireLockError(info.ResourceID)
}

func (info *Handler) Release(ctx workflow.Context) error {
	wfinfo(ctx, info, "mutex handler: requesting release ...")

	orphan := false

	if err := workflow.
		SignalExternalWorkflow(ctx, info.Execution.ID, "", WorkflowSignalRelease.String(), info).
		Get(ctx, nil); err != nil {
		wfwarn(ctx, info, "mutex handler: unable to request release ...", err)
		return NewReleaseLockError(info.ResourceID)
	}

	wfinfo(ctx, info, "mutex handler: waiting for release ...")
	workflow.GetSignalChannel(ctx, WorkflowSignalReleased.String()).Receive(ctx, &orphan)

	if orphan {
		wfwarn(ctx, info, "mutex handler: lock released, orphaned!", nil)
	} else {
		wfinfo(ctx, info, "mutex handler: lock released, done!")
	}

	return nil
}

func (info *Handler) Cleanup(ctx workflow.Context) error {
	wfinfo(ctx, info, "mutex handler: requesting cleanup ...")

	shutdown := false

	if err := workflow.
		SignalExternalWorkflow(ctx, info.Execution.ID, "", WorkflowSignalCleanup.String(), info).
		Get(ctx, nil); err != nil {
		wferr(ctx, info, "mutex handler: unable to clean up", err)
		return NewCleanupMutexError(info.ResourceID)
	}

	wfinfo(ctx, info, "mutex handler: waiting for cleanup ...")
	workflow.GetSignalChannel(ctx, WorkflowSignalCleanupDone.String()).Receive(ctx, &shutdown)

	if err := workflow.
		SignalExternalWorkflow(ctx, info.Execution.ID, "", WorkflowSignalCleanupDoneAck.String(), shutdown).
		Get(ctx, nil); err != nil {
		wfwarn(ctx, info, "mutex handler: unable to acknowledge cleanup", err)
		return NewCleanupMutexError(info.ResourceID)
	}

	if shutdown {
		wfinfo(ctx, info, "mutex handler: cleanup done!")
	} else {
		wfwarn(ctx, info, "mutex handler: unable to clean up, mutex in use", nil)
	}

	return nil
}

// validate validates if the mutex is properly configured.
func (lock *Handler) validate() error {
	if lock.ResourceID == "" {
		return ErrNoResourceID
	}

	if lock.Info == nil {
		return ErrNilContext
	}

	return nil
}

// WithResourceID sets the resource ID for the mutex workflow. We start with the assumption that a valid resource ID will be provided.
// The lock must always be held against the ids of core entities e.g. Stack, Repo or Resource. and the format may look like
// ${entity_type}.${entity_id}.mutex
//   - entity type e.g stack, repo, resource
//   - entity id e.g. the database id.
//
// for some cases, this may be made easy by getting the id of the parent workflow info e.g. if we are running stack controller, we can
// get the stack controller id, which would be in the format "io.quantm.stack.${stack_id}" and then adding the "mutex" suffix. Alernatively
// this can be set explicitly as well as "io.quantm.repo.${repo_id}.branch.${branch_name}.mutex". This is the format that should be used
// when holding locks against specific resources like repos or artifacts or cloud resources. This is a judgement call. The goal is, we
// should be able to arrive at the lock id regardless of the context.
func WithResourceID(id string) Option {
	return func(m *Handler) {
		m.ResourceID = id
	}
}

// WithHandler sets the handler for the mutex workflow.
func WithHandler(ctx workflow.Context) Option {
	return func(m *Handler) {
		m.Info = workflow.GetInfo(ctx)
	}
}

// WithTimeout sets the timeout for the mutex workflow.
func WithTimeout(timeout time.Duration) Option {
	return func(m *Handler) {
		m.Timeout = timeout
	}
}

// New returns a new Mutex.
// It should always be called with WithHandler and WithResourceID options.
// If WithTimeout is not called, it defaults to DefaultTimeout.
//
// NOTE - This workflow assumes that we always start a workflow from within a workflow. We might have to change this design if our use case
// require locking a resource from outside a workflow context. The chances of that happening are really really slim.
//
// # Usage:
//
//	m := mutex.New(
//	  mutex.WithHandler(ctx),
//	  mutex.WithResourceID("id"),
//	  mutex.WithTimeout(30*time.Minute), // Optional
//	)
//	if err := m.Prepare(ctx); err != nil {/*handle error*/}
//	if err := m.Acquire(ctx); err != nil {/*handle error*/}
//	// do work.
//	if err := m.Release(ctx); err != nil {/*handle error*/}
//	// attempt to shutdown the mutex.
//	if err := m.Cleanup(ctx); err != nil {/*handle error*/}
func New(opts ...Option) Mutex {
	m := &Handler{Timeout: DefaultTimeout}
	for _, opt := range opts {
		opt(m)
	}

	return m
}
