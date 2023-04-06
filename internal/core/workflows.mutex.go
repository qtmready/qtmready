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
	"context"
	"time"

	"go.breu.io/ctrlplane/internal/shared"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

const (
	//TODO: define these as shared.WorkflowSignal
	AcquireLockSignalName = "acquire-lock-event" // AcquireLockSignalName signal channel name for lock acquisition
	RequestLockSignalName = "request-lock-event" // RequestLockSignalName channel name for request lock
	ReleaseLockSignalName = "release-lock-event" // RequestLockSignalName channel name for request lock
)

type (
	UnlockFunc func() error

	Mutex struct {
		currentWorkflowID string
		mutexWorkflowID   string
		resourceID        string
		unlockTimeout     time.Duration
	}
)

// NewMutex initializes mutex
func NewMutex(currentWorkflowID string, resourceID string, unlockTimeout time.Duration) *Mutex {
	return &Mutex{
		currentWorkflowID: currentWorkflowID,
		resourceID:        resourceID,
		unlockTimeout:     unlockTimeout,
	}
}

func (m *Mutex) Init(ctx workflow.Context) error {

	activityCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Hour * 6, //TODO: should we set this? is 6 hours enough?
		StartToCloseTimeout:    time.Minute,
	})

	var execution workflow.Execution
	err := workflow.ExecuteActivity(activityCtx, activities.StartMutexWorkflowActivity, m.resourceID, unLockTimeOutStackMutex).Get(ctx, &execution)
	m.mutexWorkflowID = execution.ID

	return err
}

// Lock - locks mutex
func (m *Mutex) Lock(ctx workflow.Context) (UnlockFunc, error) {

	// TODO: resource - mutex workflow id map?
	logger := workflow.GetLogger(ctx)

	// request mutex workflow to acquire lock
	logger.Info("Lock: sending acquire lock signal")
	workflow.SignalExternalWorkflow(ctx, m.mutexWorkflowID, "", RequestLockSignalName, m.currentWorkflowID).Get(ctx, nil)
	logger.Info("Lock: waiting to acquire lock")

	// wait for the acknowledgement from mutex workflow that lock has been acquired
	workflow.GetSignalChannel(ctx, AcquireLockSignalName).Receive(ctx, nil)
	logger.Info("Lock: lock acquired")

	unlockFunc := func() error {
		logger.Info("Lock: sending release lock signal")
		return workflow.SignalExternalWorkflow(ctx, m.mutexWorkflowID, "", ReleaseLockSignalName, m.currentWorkflowID).Get(ctx, nil)
	}
	return unlockFunc, nil
}

/*
StartMutexWorkflowActivity will start a mutex workflow that will wait on "RequestLockSignalName" to lock the resource
resourceID: ID of the resource to be locked
unlockTimeout: timeout after which the resource will be released automatically
*/
func (a *Activities) StartMutexWorkflowActivity(ctx context.Context, resourceID string, unlockTimeout time.Duration) (*workflow.Execution, error) {

	w := &Workflows{}
	opts := shared.Temporal.Queues[shared.CoreQueue].
		GetWorkflowOptions(
			"mutex",
			resourceID,
		)

	println("Execute mutex workflow ")
	wr, err := shared.Temporal.Client.ExecuteWorkflow(ctx, opts, w.MutexWorkflow, resourceID, unlockTimeout)

	activity.GetLogger(ctx).Info("Started Workflow", "WorkflowID", wr.GetID(), "RunID", wr.GetRunID(), "Error:", err)
	return &workflow.Execution{ID: wr.GetID(), RunID: wr.GetRunID()}, err
}

/*
Mutex workflow will lock a resource specified by "resourceId"
Get request lock signal name
*/
func (w *Workflows) MutexWorkflow(ctx workflow.Context, resourceId string, unlockTimeout time.Duration) error {

	logger := workflow.GetLogger(ctx)
	selector := workflow.NewSelector(ctx)

	logger.Info("MutexWorkflow: started", "currentWorkflowID:", workflow.GetInfo(ctx).WorkflowExecution.ID)

	// get lock request from workflows on request lock channel
	var requestLockWorkflowID string
	var releaseLockWorkflowID string = ""
	requestLockCh := workflow.GetSignalChannel(ctx, RequestLockSignalName)
	releaseLockCh := workflow.GetSignalChannel(ctx, ReleaseLockSignalName)

	for {
		// wait for the acquire lock signal
		logger.Info("Waiting for acquire lock request")
		requestLockCh.Receive(ctx, &requestLockWorkflowID)

		logger.Info("Aquire lock request received from workflow ID: " + requestLockWorkflowID)

		// send lock acquired ack to sender workflow
		err := workflow.SignalExternalWorkflow(ctx, requestLockWorkflowID, "", AcquireLockSignalName, true).Get(ctx, nil)

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
