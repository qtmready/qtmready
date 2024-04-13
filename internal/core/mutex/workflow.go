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
	MutexStatusReleasing MutexStatus = "releasing"
	MutexStatusReleased  MutexStatus = "released"
	MutexStatusTimeout   MutexStatus = "timeout"
)

// Workflow is the mutex workflow. It is responsible for controlling the access to a resource. It should never be called directly, instead
// use the New function to create a new mutex, and use the Prepare, Acquire,Release and Cleanup functions to interact with the mutex.
//
// It works by listening to the following signals:
//   - WorkflowSignalPrepare: Prepares the lock by adding the caller to the queue.
//   - WorkflowSignalAcquire: Acquires the lock.
//   - WorkflowSignalRelease: Releases the lock.
//   - WorkflowSignalCleanup: Clean up shutdowns the workflow, if there are no more locks in the queue.
func Workflow(ctx workflow.Context, info *Info) error {
	wfinfo(ctx, info, "mutex: workflow started")

	persist := true                                               // persist is used to keep the workflow running.
	acquirer := info.Caller                                       // acquirer is the workflow that is trying to acquire the lock
	status := MutexStatusAcquiring                                // status is the current status of the lock.
	queue := &SafeMap{internal: make(map[string]time.Duration)}   // queue is the pool of workflows waiting to acquire the lock.
	orphans := &SafeMap{internal: make(map[string]time.Duration)} // orphans is the pool of workflows that have timed out.

	// coroutine to listen for prepare signals
	workflow.Go(ctx, func(ctx workflow.Context) {
		for persist {
			rx := &Info{}
			workflow.GetSignalChannel(ctx, WorkflowSignalPrepare.String()).Receive(ctx, rx)
			wfinfo(ctx, rx, "mutex: preparing ...", slog.Int("pool_size", len(queue.internal)))

			persist = true

			queue.Add(ctx, rx.Caller.WorkflowExecution.ID, rx.Timeout)

			wfinfo(ctx, rx, "mutex: new lock request recieved, waiting in queue ...", slog.Int("pool_size", len(queue.internal)))
		}
	})

	// coroutine to listen for cleanup signals
	workflow.Go(ctx, func(ctx workflow.Context) {
		for persist {
			rx := &Info{}
			workflow.GetSignalChannel(ctx, WorkflowSignalCleanup.String()).Receive(ctx, rx)
			wfinfo(ctx, rx, "mutex: cleanup requested ...", slog.Int("pool_size", len(queue.internal)))

			if len(queue.internal) == 0 {
				persist = false
			}

			err := workflow.
				SignalExternalWorkflow(ctx, rx.Caller.WorkflowExecution.ID, "", WorkflowSignalCleanupDone.String(), persist).
				Get(ctx, nil)

			if err != nil {
				wfwarn(ctx, rx, "mutex: unable cleanup", err)
			}

			wfinfo(ctx, info, "mutex: clean up done!", slog.Int("pool_size", len(queue.internal)), slog.Bool("persist", persist))
		}
	})

	// main loop
	for persist {
		wfinfo(ctx, info, "mutex: waiting for lock ...")

		workflow.GetSignalChannel(ctx, WorkflowSignalAcquire.String()).Receive(ctx, acquirer)

		info.Caller = acquirer
		timeout, ok := queue.Get(ctx, acquirer.WorkflowExecution.ID)

		if !ok {
			wfwarn(ctx, info, "mutex: unable to find the timeout against the requested lock, aborting ...", nil)
			continue
		}

		wfinfo(
			ctx, info, "mutex: lock requested ...",
			slog.String("requested_by", acquirer.WorkflowExecution.ID), slog.String("timeout", timeout.String()),
		)

		if err := workflow.
			SignalExternalWorkflow(ctx, acquirer.WorkflowExecution.ID, "", WorkflowSignalLocked.String(), ok).
			Get(ctx, nil); err != nil {
			wfwarn(ctx, info, "mutex: unable to acquire lock, retrying ...", err)
			continue
		}

		status = MutexStatusLocked

		wfinfo(ctx, info, "mutex: lock acquired", slog.String("acquired_by", acquirer.WorkflowExecution.ID))

		for {
			wfinfo(ctx, info, "mutex: waiting for release or timeout ...")
			selector := workflow.NewSelector(ctx)
			selector.AddReceive(
				workflow.GetSignalChannel(ctx, WorkflowSignalRelease.String()),
				_release(ctx, &status, info, queue, orphans, acquirer),
			)
			selector.AddFuture(
				workflow.NewTimer(ctx, timeout),
				_timeout(ctx, &status, info, queue, orphans, timeout),
			)

			selector.Select(ctx)

			// if the lock is released or timed out, set the status to acquiring and break the loop.
			if status == MutexStatusReleased || status == MutexStatusTimeout {
				status = MutexStatusAcquiring
				break
			}
		}
	}

	wfinfo(ctx, info, "mutex: shutdown!")

	return nil
}

// _release is a channel handler that is called when lock is to be released.
// TODO - handle the case when the lock is found in the orphans pool.
func _release(
	ctx workflow.Context, status *MutexStatus, info *Info, queue, orphans *SafeMap, acquirer *workflow.Info,
) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		releaser := &Info{}
		rx.Receive(ctx, releaser)

		wfinfo(ctx, info, "mutex: releasing ...", slog.Int("pool_size", len(queue.internal)))

		if releaser.Caller.WorkflowExecution.ID == acquirer.WorkflowExecution.ID {
			*status = MutexStatusReleasing
			_, ok := orphans.Get(ctx, acquirer.WorkflowExecution.ID)

			err := workflow.
				SignalExternalWorkflow(ctx, acquirer.WorkflowExecution.ID, "", WorkflowSignalReleased.String(), ok).
				Get(ctx, nil)

			if err != nil {
				wfwarn(ctx, info, "mutex: unable to release lock, retrying ...", err)

				return
			}

			queue.Remove(ctx, acquirer.WorkflowExecution.ID)

			*status = MutexStatusReleased

			wfinfo(ctx, info, "mutex: release done!", slog.Int("pool_size", len(queue.internal)))
		}
	}
}

// _timeout is a future handler that is called when the lock has timed out.
// TODO - what happens at the acquirer when the lock times out?
func _timeout(ctx workflow.Context, status *MutexStatus, info *Info, pool, orphans *SafeMap, timeout time.Duration) shared.FutureHandler {
	return func(future workflow.Future) {
		if *status == MutexStatusReleasing {
			wfinfo(ctx, info, "mutex: lock timeout reached", slog.Duration("timeout", timeout))
			pool.Remove(ctx, info.Caller.WorkflowExecution.ID)
			orphans.Add(ctx, info.Caller.WorkflowExecution.ID, timeout)

			*status = MutexStatusTimeout
		}
	}
}
