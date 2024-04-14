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
	"sync"
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
func Workflow(ctx workflow.Context, lock *Info) error {
	wfinfo(ctx, lock, "mutex: workflow started")

	persist := true                                                                     // persist is used to keep the workflow running.
	active := &Info{}                                                                   // active is the active lock request.
	status := MutexStatusAcquiring                                                      // status is the current status of the lock.
	queue := &SafeMap{Mutex: &sync.Mutex{}, Internal: make(map[string]time.Duration)}   // queue is the pool for schduled lock requests.
	orphans := &SafeMap{Mutex: &sync.Mutex{}, Internal: make(map[string]time.Duration)} // orphans is the pool for timed out locks.

	// coroutine that schedules requests for the lock against the queue.
	workflow.Go(ctx, func(ctx workflow.Context) {
		for persist {
			rx := &Info{}
			workflow.GetSignalChannel(ctx, WorkflowSignalPrepare.String()).Receive(ctx, rx)
			wfinfo(ctx, rx, "mutex: preparing ...", slog.Int("pool_size", len(queue.Internal)))

			queue.Add(ctx, rx.Caller.WorkflowExecution.ID, rx.Timeout)

			wfinfo(ctx, rx, "mutex: new lock request recieved, waiting in queue ...", slog.Int("pool_size", len(queue.Internal)))
		}
	})

	// main loop to handle
	//  - acquiring the lock
	//  - garbage collection
	//  - releasing the lock
	//  - timeout
	for persist {
		wfinfo(ctx, lock, "mutex: waiting for lock ...")

		found := true
		timeout := time.Duration(0)
		acquirer := workflow.NewSelector(ctx)

		acquirer.AddReceive(workflow.GetSignalChannel(ctx, WorkflowSignalAcquire.String()), _acquire(ctx, active, queue, &timeout, &found))
		acquirer.AddReceive(workflow.GetSignalChannel(ctx, WorkflowSignalCleanup.String()), _cleanup(ctx, active, queue, &persist))
		acquirer.Select(ctx)

		// cleanup signal received and processed. queue is empty, so shutdown the workflow.
		if !persist {
			break
		}

		// lock found in queue, set the status to locked and continue to next step, else set the status to acquiring and restart the loop.
		if found {
			status = MutexStatusLocked
		} else {
			status = MutexStatusAcquiring

			continue
		}

		wfinfo(ctx, active, "mutex: lock acquired!")

		for {
			wfinfo(ctx, active, "mutex: waiting for release or timeout ...")
			releaser := workflow.NewSelector(ctx)

			releaser.AddReceive(
				workflow.GetSignalChannel(ctx, WorkflowSignalRelease.String()), _release(ctx, active, &status, queue, orphans),
			)
			releaser.AddFuture(workflow.NewTimer(ctx, timeout), _timeout(ctx, active, &status, queue, orphans, timeout))

			releaser.Select(ctx)

			// if the lock is released or timed out, set the status to acquiring and break the loop.
			if status == MutexStatusReleased || status == MutexStatusTimeout {
				status = MutexStatusAcquiring
				break
			}
		}
	}

	wfinfo(ctx, lock, "mutex: shutdown!")

	return nil
}

func _acquire(ctx workflow.Context, lock *Info, queue *SafeMap, timeout *time.Duration, found *bool) shared.ChannelHandler {
	return func(channel workflow.ReceiveChannel, more bool) {
		rx := &Info{}
		channel.Receive(ctx, rx)

		lock.Caller = rx.Caller
		*timeout, *found = queue.Get(ctx, rx.Caller.WorkflowExecution.ID)

		wfinfo(ctx, lock, "mutex: lock received ...", slog.String("requested_by", rx.Caller.WorkflowExecution.ID))

		if err := workflow.
			SignalExternalWorkflow(ctx, rx.Caller.WorkflowExecution.ID, "", WorkflowSignalLocked.String(), *found).
			Get(ctx, nil); err != nil {
			wfwarn(ctx, lock, "mutex: unable to acquire lock, retrying ...", err)
		}
	}
}

func _cleanup(ctx workflow.Context, lock *Info, queue *SafeMap, persist *bool) shared.ChannelHandler {
	return func(channel workflow.ReceiveChannel, more bool) {
		rx := &Info{}
		channel.Receive(ctx, rx)

		wfinfo(ctx, lock, "mutex: cleanup requested ...", slog.Int("pool_size", len(queue.Internal)))

		if len(queue.Internal) == 0 {
			*persist = false
		}

		err := workflow.
			SignalExternalWorkflow(ctx, rx.Caller.WorkflowExecution.ID, "", WorkflowSignalCleanupDone.String(), *persist).
			Get(ctx, nil)

		if err != nil {
			wfwarn(ctx, lock, "mutex: unable cleanup", err)
		}

		wfinfo(ctx, lock, "mutex: scheduling cleanup!", slog.Int("pool_size", len(queue.Internal)), slog.Bool("persist", *persist))
	}
}

// _release is a channel handler that is called when lock is to be released.
// TODO - handle the case when the lock is found in the orphans pool.
func _release(ctx workflow.Context, active *Info, status *MutexStatus, queue, orphans *SafeMap) shared.ChannelHandler {
	return func(channel workflow.ReceiveChannel, more bool) {
		rx := &Info{}
		channel.Receive(ctx, rx)

		wfinfo(ctx, active, "mutex: releasing ...", slog.Int("pool_size", len(queue.Internal)))

		if rx.Caller.WorkflowExecution.ID == active.Caller.WorkflowExecution.ID {
			*status = MutexStatusReleasing
			_, ok := orphans.Get(ctx, active.Caller.WorkflowExecution.ID)

			err := workflow.
				SignalExternalWorkflow(ctx, active.Caller.WorkflowExecution.ID, "", WorkflowSignalReleased.String(), ok).
				Get(ctx, nil)

			if err != nil {
				wfwarn(ctx, active, "mutex: unable to release lock, retrying ...", err)

				return
			}

			queue.Remove(ctx, active.Caller.WorkflowExecution.ID)

			*status = MutexStatusReleased

			wfinfo(ctx, active, "mutex: release done!", slog.Int("pool_size", len(queue.Internal)))
		}
	}
}

// _timeout is a future handler that is called when the lock has timed out.
// TODO - what happens at the acquirer when the lock times out?
func _timeout(ctx workflow.Context, active *Info, status *MutexStatus, pool, orphans *SafeMap, timeout time.Duration) shared.FutureHandler {
	return func(future workflow.Future) {
		if *status == MutexStatusLocked && *status != MutexStatusReleasing && timeout > 0 {
			wfinfo(ctx, active, "mutex: timeout reached, releasing ...", slog.Duration("timeout", timeout))
			pool.Remove(ctx, active.Caller.WorkflowExecution.ID)
			orphans.Add(ctx, active.Caller.WorkflowExecution.ID, timeout)

			*status = MutexStatusTimeout
		}

		wfwarn(ctx, active, "mutex: ignoring timeout ...", nil)
	}
}
