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
func Workflow(ctx workflow.Context, starter *Handler) error {
	wfinfo(ctx, starter, "mutex: workflow started")

	persist := true                // persist is used to keep the workflow running.
	handler := &Handler{}          // handler is the active lock request.
	status := MutexStatusAcquiring // status is the current status of the lock.
	pool := NewSimpleMap()         // queue heolds the lock requests.
	orphans := NewSimpleMap()      // orphans holds the lock requests that have timed out.
	shutdown, shutdownfn := workflow.NewFuture(ctx)

	// coroutines responsible for scheduling the lock request or scheduling a graceful shutdown of the mutex workflow.
	workflow.Go(ctx, _prepare(ctx, starter, &pool, &persist))
	workflow.Go(ctx, _cleanup(ctx, starter, &pool, shutdownfn))

	// main loop to handle
	//  - acquiring the lock
	//  - garbage collection
	//  - releasing the lock
	//  - timeout
	for persist {
		wfinfo(ctx, starter, "mutex: waiting for lock request ...")

		found := true
		cleanup := false
		timeout := time.Duration(0)
		acquirer := workflow.NewSelector(ctx)

		acquirer.AddReceive(workflow.GetSignalChannel(ctx, WorkflowSignalAcquire.String()), _acquire(ctx, handler, &pool, &timeout, &found))
		acquirer.AddFuture(shutdown, _shutdown(ctx, handler, &persist))

		acquirer.Select(ctx)

		// cleanup signal received and processed. queue is empty, so shutdown the workflow.
		if !persist || cleanup {
			wfinfo(ctx, starter, "mutex: cleanup done, shutting down ...")

			continue
		}

		// lock found in queue, set the status to locked and continue to next step, else set the status to acquiring and restart the loop.
		if !found {
			wfinfo(ctx, handler, "mutex: lock not found in the pool, retrying ...")

			status = MutexStatusAcquiring

			continue
		} else {
			wfinfo(ctx, handler, "mutex: lock acquired!")

			status = MutexStatusLocked
		}

		for {
			wfinfo(ctx, handler, "mutex: waiting for release or timeout ...")
			releaser := workflow.NewSelector(ctx)

			releaser.AddReceive(
				workflow.GetSignalChannel(ctx, WorkflowSignalRelease.String()), _release(ctx, handler, &status, &pool, &orphans),
			)
			releaser.AddFuture(workflow.NewTimer(ctx, timeout), _abort(ctx, handler, &status, &pool, &orphans, timeout))

			releaser.Select(ctx)

			// if the lock is released or timed out, set the status to acquiring and break the loop.
			if status == MutexStatusReleased || status == MutexStatusTimeout {
				status = MutexStatusAcquiring
				break
			}
		}
	}

	wfinfo(ctx, starter, "mutex: shutdown!")

	return nil
}

// _prepare is a coroutine that listens to the prepare signal and adds the lock request to the queue.
func _prepare(ctx workflow.Context, handler *Handler, pool *Pool, persist *bool) shared.CoroutineHandler {
	wfinfo(ctx, handler, "mutex: setting up workflow to prepare signal ....")

	return func(ctx workflow.Context) {
		for *persist {
			rx := &Handler{}
			workflow.GetSignalChannel(ctx, WorkflowSignalPrepare.String()).Receive(ctx, rx)

			wfinfo(ctx, rx, "mutex: prepare request received ...", slog.Int("pool_size", pool.Size()))

			pool.Add(rx.Info.WorkflowExecution.ID, rx.Timeout)

			wfinfo(ctx, rx, "mutex: prepared!", slog.Int("pool_size", pool.Size()))
		}
	}
}

// _acquire is a channel handler that is called when the lock is to be acquired.
func _acquire(ctx workflow.Context, handler *Handler, pool *Pool, timeout *time.Duration, found *bool) shared.ChannelHandler {
	return func(channel workflow.ReceiveChannel, more bool) {
		rx := &Handler{}
		channel.Receive(ctx, rx)

		*handler = *rx
		*timeout, *found = pool.Get(rx.Info.WorkflowExecution.ID)

		wfinfo(ctx, handler, "mutex: lock request received ...", slog.String("requested_by", rx.Info.WorkflowExecution.ID))

		if err := workflow.
			SignalExternalWorkflow(ctx, rx.Info.WorkflowExecution.ID, "", WorkflowSignalLocked.String(), *found).
			Get(ctx, nil); err != nil {
			wfwarn(ctx, handler, "mutex: unable to acquire lock, retrying ...", err)
		}
	}
}

// _release is a channel handler that is called when lock is to be released.
// TODO - handle the case when the lock is found in the orphans pool.
func _release(ctx workflow.Context, handler *Handler, status *MutexStatus, pool, orphans *Pool) shared.ChannelHandler {
	return func(channel workflow.ReceiveChannel, more bool) {
		rx := &Handler{}
		channel.Receive(ctx, rx)

		wfinfo(ctx, handler, "mutex: releasing ...", slog.Int("pool_size", pool.Size()))

		if rx.Info.WorkflowExecution.ID == handler.Info.WorkflowExecution.ID {
			*status = MutexStatusReleasing
			_, ok := orphans.Get(handler.Info.WorkflowExecution.ID)

			err := workflow.
				SignalExternalWorkflow(ctx, handler.Info.WorkflowExecution.ID, "", WorkflowSignalReleased.String(), ok).
				Get(ctx, nil)

			if err != nil {
				wfwarn(ctx, handler, "mutex: unable to release lock, retrying ...", err)

				return
			}

			pool.Remove(handler.Info.WorkflowExecution.ID)

			*status = MutexStatusReleased

			wfinfo(ctx, handler, "mutex: release done!", slog.Int("pool_size", pool.Size()))
		}
	}
}

// _abort is a future handler that is called when the lock has timed out.
// TODO - what happens at the acquirer when the lock times out?
func _abort(ctx workflow.Context, handler *Handler, status *MutexStatus, pool, orphans *Pool, timeout time.Duration) shared.FutureHandler {
	return func(future workflow.Future) {
		if *status == MutexStatusLocked && *status != MutexStatusReleasing && timeout > 0 {
			wfinfo(ctx, handler, "mutex: timeout reached, releasing ...", slog.Duration("timeout", timeout))
			pool.Remove(handler.Info.WorkflowExecution.ID)
			orphans.Add(handler.Info.WorkflowExecution.ID, timeout)

			*status = MutexStatusTimeout
		}

		wfwarn(ctx, handler, "mutex: ignoring timeout ...", nil)
	}
}

// cleanup is a coroutine that listens to the cleanup signal and shuts down the workflow if the queue is empty.
func _cleanup(ctx workflow.Context, handler *Handler, pool *Pool, fn workflow.Settable) shared.CoroutineHandler {
	wfinfo(ctx, handler, "mutex: setting up workflow for cleanup signals ....")

	shutdown := false

	return func(ctx workflow.Context) {
		for !shutdown {
			rx := &Handler{}
			workflow.GetSignalChannel(ctx, WorkflowSignalCleanup.String()).Receive(ctx, rx)

			wfinfo(ctx, handler, "mutex: cleanup requested ...", slog.Int("pool_size", pool.Size()))

			if pool.Size() == 0 {
				wfinfo(ctx, handler, "mutex: requesting shutdown ...", slog.Int("pool_size", pool.Size()))
				fn.Set(rx, nil)

				shutdown = true
			}

			workflow.SignalExternalWorkflow(ctx, rx.Info.WorkflowExecution.ID, "", WorkflowSignalCleanupDone.String(), shutdown)
			wfinfo(ctx, handler, "mutex: cleanup done!", slog.Int("pool_size", pool.Size()))
		}
	}
}

func _shutdown(ctx workflow.Context, handler *Handler, persist *bool) shared.FutureHandler {
	return func(future workflow.Future) {
		rx := &Handler{}
		_ = future.Get(ctx, rx)

		wfinfo(ctx, handler, "mutex: shutdown requested ...")

		*persist = false

		wfinfo(ctx, handler, "mutex: shutting down ...")
	}
}
