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

package core

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/ctrlplane/internal/shared"
)

type (
	UnlockFunc func() error

	// Mutex holds the informations for a specific mutex resource.
	Mutex struct {
		CurrentWorkflowID string
		MutexWorkflowID   string
		ResourceID        string
		UnlockTimeout     time.Duration
	}
)

// NewMutex creates a mutex.
func NewMutex(currentWorkflowID string, resourceID string, unlockTimeout time.Duration) *Mutex {
	return &Mutex{
		CurrentWorkflowID: currentWorkflowID,
		ResourceID:        resourceID,
		UnlockTimeout:     unlockTimeout,
	}
}

// Init starts a mutex workflow for the given resourceID.
// The input parameter resourceID is the ID of the resource to be locked and
// unlockTimeout is the timeout after which the resource will be released automatically.
func (m *Mutex) Init(ctx workflow.Context) error {
	w := &Workflows{}

	logger := workflow.GetLogger(ctx)
	opts := shared.Temporal.Queues[shared.CoreQueue].GetChildWorkflowOptions("mutex", m.ResourceID)
	ctx = workflow.WithChildOptions(ctx, opts)

	// execute child workflow and wait for it to spawn
	var execution workflow.Execution
	err := workflow.
		ExecuteChildWorkflow(ctx, w.MutexWorkflow, m.ResourceID, m.UnlockTimeout).
		GetChildWorkflowExecution().Get(ctx, execution)

	if err != nil {
		logger.Error("Parent execution received child execution failure.", "Error", err)
		return err
	}

	m.MutexWorkflowID = execution.ID
	logger.Info("Started Child Mutex Workflow", "ID", execution.ID, "RunID", execution.RunID, "Error:", err)

	return err
}

// Lock sends request lock signal to the mutex workflow and waits for acknowledgement.
// It returns an unlock function that should be called to release the resource.
func (m *Mutex) Lock(ctx workflow.Context) (UnlockFunc, error) {
	// TODO: resource - mutex workflow id map?
	logger := workflow.GetLogger(ctx)

	// request mutex workflow to acquire lock
	logger.Info("Lock: sending acquire lock signal")

	_ = workflow.
		SignalExternalWorkflow(ctx, m.MutexWorkflowID, "", WorkflowSignalRequestLock.String(), m.CurrentWorkflowID).
		Get(ctx, nil)

	logger.Info("Lock: waiting to acquire lock")

	// wait for the acknowledgement from mutex workflow that lock has been acquired
	workflow.GetSignalChannel(ctx, WorkflowSignalLockAcquired.String()).Receive(ctx, nil)
	logger.Info("Lock: lock acquired")

	unlockFunc := func() error {
		logger.Info("Lock: sending release lock signal")
		return workflow.SignalExternalWorkflow(ctx, m.MutexWorkflowID, "", WorkflowSignalReleaseLock.String(), m.CurrentWorkflowID).Get(ctx, nil)
	}

	return unlockFunc, nil
}

// MutexWorkflow will lock a resource specified by resourceId. The resource wil be automatically released after unlockTimeout.
// MutexWorkflow waits for a request lock signal, sends acknowledgement to the workflow requesting lock and then waits for release lock signal.
func (w *Workflows) MutexWorkflow(ctx workflow.Context, resourceID string, unlockTimeout time.Duration) error {
	logger := workflow.GetLogger(ctx)
	selector := workflow.NewSelector(ctx)

	logger.Info("MutexWorkflow: started", "currentWorkflowID:", workflow.GetInfo(ctx).WorkflowExecution.ID)

	// get lock request from workflows on request lock channel
	var requestLockWorkflowID string

	releaseLockWorkflowID := ""
	requestLockCh := workflow.GetSignalChannel(ctx, WorkflowSignalRequestLock.String())
	releaseLockCh := workflow.GetSignalChannel(ctx, WorkflowSignalReleaseLock.String())

	for {
		// wait for the acquire lock signal
		logger.Info("Waiting for acquire lock request")
		requestLockCh.Receive(ctx, &requestLockWorkflowID)

		logger.Info("Aquire lock request received from workflow ID: " + requestLockWorkflowID)

		// send lock acquired ack to sender workflow
		err := workflow.SignalExternalWorkflow(ctx, requestLockWorkflowID, "", WorkflowSignalLockAcquired.String(), nil).Get(ctx, nil)

		if err != nil {
			// .Get(ctx, nil) blocks until the signal is sent.
			// If the senderWorkflowID is closed (terminated/canceled/timeouted/completed/etc), this would return error. In this case we release the lock
			// immediately instead of failing the mutex workflow. Mutex workflow failing would lead to all workflows that have sent requestLock will be waiting.
			logger.Info("SignalExternalWorkflow error", "Error", err)
			continue
		}

		// start timer, add future which will execute when timer expires
		// TODO: how will the lock release on timeout?
		selector.AddFuture(workflow.NewTimer(ctx, unlockTimeout), func(f workflow.Future) {
			logger.Info("MutexWorkflow: unlockTimeout exceeded")
		})

		// wait for a release lock signal from the workflow that has acquired the lock
		logger.Info("MutexWorkflow: wait for release lock signal")

		for {
			releaseLockCh.Receive(ctx, &releaseLockWorkflowID)

			if releaseLockWorkflowID == requestLockWorkflowID {
				logger.Info("MutexWorkflow: release lock signal received from workflow: " + releaseLockWorkflowID)
				break
			}

			logger.Info("MutexWorkflow: release signal received from a workflow that is not holding the lock. sending workflow ID:", releaseLockWorkflowID)
		}
	}
}
