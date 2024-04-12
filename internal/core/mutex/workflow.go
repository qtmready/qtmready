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

type (
	MutexStatus = string
)

var (
	MutexStatusAcquiring MutexStatus = "acquiring"
	MutexStatusLocked    MutexStatus = "locked"
	MutexStatusReleased  MutexStatus = "released"
	MutexStatusTimeout   MutexStatus = "timeout"
)

// Workflow is the mutex workflow. It is responsible for controlling the access to a resource. It should always be started as
// a child workflow.
// It works by listening to two signals:
//   - WorkflowSignalAcquire: this signal is sent by the caller to acquire the lock.
//   - WorkflowSignalRelease: this signal is sent by the caller to release the lock.
func Workflow(ctx workflow.Context, info *Info) error {
	wfinfo(ctx, info, "mutex: workflow started")

	persist := true                                            // persist is used to keep the workflow running.
	acquirer := info.Caller                                    // acquirer is the workflow that is trying to acquire the lock
	releaser := &workflow.Info{}                               // releaser is the workflow that is trying to release the lock
	status := MutexStatusAcquiring                             // status is the current status of the lock.
	queue := &Pool{internal: make(map[string]time.Duration)}   // queue is the pool of workflows waiting to acquire the lock.
	orphans := &Pool{internal: make(map[string]time.Duration)} // orphans is the pool of workflows that have timed out.

	// coroutine to listen for prepare signals
	workflow.Go(ctx, func(ctx workflow.Context) {
		for persist {
			rx := &Info{}
			workflow.GetSignalChannel(ctx, WorkflowSignalPrepare.String()).Receive(ctx, rx)
			wfinfo(ctx, rx, "mutex: received prepare signal", slog.Int("pool_size", len(queue.internal)))
			queue.Add(rx.Caller.WorkflowExecution.ID, rx.Timeout)
			wfinfo(ctx, rx, "mutex: added to pool", slog.Int("pool_size", len(queue.internal)))
		}
	})

	// coroutine to listen for cleanup signals
	workflow.Go(ctx, func(ctx workflow.Context) {
		for persist {
			rx := &Info{}
			workflow.GetSignalChannel(ctx, WorkflowSignalCleanup.String()).Receive(ctx, rx)
			wfinfo(ctx, rx, "mutex: received cleanup signal", slog.Int("pool_size", len(queue.internal)))

			if len(queue.internal) == 0 {
				persist = false

				continue
			}

			err := workflow.
				SignalExternalWorkflow(ctx, rx.Caller.WorkflowExecution.ID, "", WorkflowSignalCleanupDone.String(), persist).
				Get(ctx, nil)

			if err != nil {
				wfwarn(ctx, rx, "mutex: unable cleanup", err)
			}

			wfinfo(ctx, info, "mutex: skipping cleanup", slog.Int("pool_size", len(queue.internal)))
		}
	})

	// main loop
	for persist {
		wfinfo(ctx, info, "mutex: waiting for lock to be acquired ...")

		workflow.GetSignalChannel(ctx, WorkflowSignalAcquire.String()).Receive(ctx, acquirer)
		info.Caller = acquirer
		status = MutexStatusLocked
		timeout, ok := queue.Read(acquirer.WorkflowExecution.ID)

		if err := workflow.
			SignalExternalWorkflow(ctx, acquirer.WorkflowExecution.ID, "", WorkflowSignalLocked.String(), ok).
			Get(ctx, nil); err != nil {
			wferr(ctx, info, "mutex: unable to acquire lock.", err)
			continue
		}

		if !ok {
			wfwarn(ctx, info, "mutex: unable to read the lock from pool", nil)
			continue
		}

		wfinfo(ctx, info, "mutex: lock acquired", slog.String("acquired_by", acquirer.WorkflowExecution.ID))

		for {
			selector := workflow.NewSelector(ctx)
			selector.AddReceive(
				workflow.GetSignalChannel(ctx, WorkflowSignalRelease.String()),
				_release(ctx, &status, info, queue, orphans, acquirer, releaser),
			)
			selector.AddFuture(
				workflow.NewTimer(ctx, timeout),
				_timeout(ctx, &status, info, queue, orphans, timeout),
			)

			selector.Select(ctx)

			if status == MutexStatusReleased || status == MutexStatusTimeout {
				break
			}
		}
	}

	wfinfo(ctx, info, "mutex: workflow completed")

	return nil
}

// _release is a channel handler that is called when the lock is released.
// TODO - handle the case when the lock is found in the orphans pool.
func _release(
	ctx workflow.Context, status *MutexStatus, info *Info, queue, orphans *Pool, acquirer, releaser *workflow.Info,
) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		rx.Receive(ctx, releaser)

		if releaser.WorkflowExecution.ID == acquirer.WorkflowExecution.ID {
			wfinfo(ctx, info, "mutex: releasing lock ...", slog.String("released_by", releaser.WorkflowExecution.ID))

			_, ok := orphans.Read(releaser.WorkflowExecution.ID)

			err := workflow.
				SignalExternalWorkflow(ctx, acquirer.WorkflowExecution.ID, "", WorkflowSignalReleased.String(), ok).
				Get(ctx, nil)

			if err != nil {
				wfwarn(ctx, info, "mutex: unable to release lock", err)

				return
			}

			queue.Remove(acquirer.WorkflowExecution.ID)

			*status = MutexStatusReleased

			wfinfo(ctx, info, "mutex: removed from pool", slog.Int("pool_size", len(queue.internal)))
		}
	}
}

// _timeout is a future handler that is called when the lock has timed out.
// TODO - what happens at the acquirer when the lock times out?
func _timeout(ctx workflow.Context, status *MutexStatus, info *Info, pool, orphans *Pool, timeout time.Duration) shared.FutureHandler {
	return func(future workflow.Future) {
		wfinfo(ctx, nil, "mutex: lock timeout reached", slog.Duration("timeout", timeout))
		pool.Remove(info.Caller.WorkflowExecution.ID)
		orphans.Add(info.Caller.WorkflowExecution.ID, timeout)

		*status = MutexStatusTimeout
	}
}
