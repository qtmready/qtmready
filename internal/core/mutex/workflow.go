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

// Workflow is the mutex workflow. It controls access to a resource.
//
// IMPORTANT: Do not use this function directly. Instead, use mutex.New to create and interact with mutex instances.
//
// The workflow consists of three main event loops:
//  1. Main loop: Handles acquiring, releasing, and timing out of locks.
//  2. Prepare loop: Listens for and handles preparation of lock requests.
//  3. Cleanup loop: Manages the cleanup process and potential workflow shutdown.
//
// It operates as a state machine, transitioning between MutexStatus states:
// Acquiring -> Locked -> Releasing -> Released (or Timeout)
//
// Uses two pools to manage lock requests:
//   - Main pool: Tracks active lock requests and currently held locks.
//   - Orphans pool: Tracks locks that have timed out.
//
// Responds to several signals:
//   - WorkflowSignalPrepare: Prepares a new lock request.
//   - WorkflowSignalAcquire: Attempts to acquire a lock.
//   - WorkflowSignalRelease: Releases a held lock.
//   - WorkflowSignalCleanup: Initiates the cleanup process.
func Workflow(ctx workflow.Context, starter *Handler) error {
	info := workflow.GetInfo(ctx)
	logger := NewMutexControllerLogger(ctx, info.WorkflowExecution.ID)
	logger.info(starter.Info.WorkflowExecution.ID, "start", "mutex: workflow started",
		"resource_id", starter.ResourceID,
		"mutex_id", info.WorkflowExecution.ID)

	handler := starter
	status := MutexStatusAcquiring
	pool := NewPool(ctx)
	orphans := NewPool(ctx)
	shutdown, shutdownfn := workflow.NewFuture(ctx)

	workflow.Go(ctx, prepare(ctx, handler, pool, logger))
	workflow.Go(ctx, cleanup(ctx, handler, pool, shutdownfn, logger))

	for {
		logger.info(handler.Info.WorkflowExecution.ID, "main_loop", "mutex: waiting for lock request ...")

		found := true
		timeout := time.Duration(0)
		acquirer := workflow.NewSelector(ctx)

		acquirer.AddReceive(workflow.GetSignalChannel(ctx, WorkflowSignalAcquire.String()), acquire(ctx, &handler, pool, &timeout, &found, logger))
		acquirer.AddFuture(shutdown, terminate(ctx, handler, logger))

		acquirer.Select(ctx)

		if handler == nil {
			break // Shutdown signal received
		}

		if !found {
			logger.info(handler.Info.WorkflowExecution.ID, "main_loop", "mutex: lock not found in the pool, retrying ...")

			status = MutexStatusAcquiring

			continue
		}

		logger.info(handler.Info.WorkflowExecution.ID, "main_loop", "mutex: lock acquired!")

		status = MutexStatusLocked

		for {
			logger.info(handler.Info.WorkflowExecution.ID, "main_loop", "mutex: waiting for release or timeout ...")

			releaser := workflow.NewSelector(ctx)

			releaser.AddReceive(
				workflow.GetSignalChannel(ctx, WorkflowSignalRelease.String()),
				release(ctx, handler, &status, pool, orphans, logger),
			)
			releaser.AddFuture(workflow.NewTimer(ctx, timeout), abort(ctx, handler, &status, pool, orphans, timeout, logger))

			releaser.Select(ctx)

			if status == MutexStatusReleased || status == MutexStatusTimeout {
				status = MutexStatusAcquiring
				break
			}
		}
	}

	_ = workflow.Sleep(ctx, 500*time.Millisecond)

	logger.info(info.WorkflowExecution.ID, "shutdown", "mutex: shutdown!")

	return nil
}

func prepare(ctx workflow.Context, handler *Handler, pool *Pool, logger *MutexLogger) func(workflow.Context) {
	return func(ctx workflow.Context) {
		for {
			rx := &Handler{}
			workflow.GetSignalChannel(ctx, WorkflowSignalPrepare.String()).Receive(ctx, rx)

			logger.info(rx.Info.WorkflowExecution.ID, "prepare", "mutex: prepare request received ...",
				"pool_size", pool.size(),
				"mutex_id", handler.Info.WorkflowExecution.ID)
			pool.add(ctx, rx.Info.WorkflowExecution.ID, rx.Timeout)
			logger.info(rx.Info.WorkflowExecution.ID, "prepare", "mutex: prepared!",
				"pool_size", pool.size(),
				"mutex_id", handler.Info.WorkflowExecution.ID)
		}
	}
}

func acquire(ctx workflow.Context, handler **Handler, pool *Pool, timeout *time.Duration, found *bool, logger *MutexLogger) shared.ChannelHandler {
	return func(channel workflow.ReceiveChannel, more bool) {
		rx := &Handler{}
		channel.Receive(ctx, rx)

		*handler = rx
		*timeout, *found = pool.get(rx.Info.WorkflowExecution.ID)

		logger.info(rx.Info.WorkflowExecution.ID, "acquire", "mutex: lock request received ...",
			"requested_by", rx.Info.WorkflowExecution.ID,
			"mutex_id", (*handler).Info.WorkflowExecution.ID)

		if err := workflow.
			SignalExternalWorkflow(ctx, rx.Info.WorkflowExecution.ID, "", WorkflowSignalLocked.String(), *found).
			Get(ctx, nil); err != nil {
			logger.warn(rx.Info.WorkflowExecution.ID, "acquire", "mutex: unable to acquire lock, retrying ...",
				"error", err,
				"mutex_id", (*handler).Info.WorkflowExecution.ID)
		}
	}
}

func release(ctx workflow.Context, handler *Handler, status *MutexStatus, pool, orphans *Pool, logger *MutexLogger) shared.ChannelHandler {
	return func(channel workflow.ReceiveChannel, more bool) {
		rx := &Handler{}
		channel.Receive(ctx, rx)

		logger.info(rx.Info.WorkflowExecution.ID, "release", "mutex: releasing ...",
			"pool_size", pool.size(),
			"mutex_id", handler.Info.WorkflowExecution.ID)

		if rx.Info.WorkflowExecution.ID == handler.Info.WorkflowExecution.ID {
			*status = MutexStatusReleasing
			_, ok := orphans.get(handler.Info.WorkflowExecution.ID)

			err := workflow.
				SignalExternalWorkflow(ctx, handler.Info.WorkflowExecution.ID, "", WorkflowSignalReleased.String(), ok).
				Get(ctx, nil)

			if err != nil {
				logger.warn(handler.Info.WorkflowExecution.ID, "release", "mutex: unable to release lock, retrying ...",
					"error", err,
					"mutex_id", handler.Info.WorkflowExecution.ID)
				return
			}

			pool.remove(ctx, handler.Info.WorkflowExecution.ID)

			*status = MutexStatusReleased

			logger.info(handler.Info.WorkflowExecution.ID, "release", "mutex: release done!",
				"pool_size", pool.size(),
				"mutex_id", handler.Info.WorkflowExecution.ID)
		}
	}
}

func abort(ctx workflow.Context, handler *Handler, status *MutexStatus, pool, orphans *Pool, timeout time.Duration, logger *MutexLogger) shared.FutureHandler {
	return func(future workflow.Future) {
		if *status == MutexStatusLocked && *status != MutexStatusReleasing && timeout > 0 {
			logger.info(handler.Info.WorkflowExecution.ID, "abort", "mutex: timeout reached, releasing ...",
				"timeout", timeout,
				"mutex_id", handler.Info.WorkflowExecution.ID)
			pool.remove(ctx, handler.Info.WorkflowExecution.ID)
			orphans.add(ctx, handler.Info.WorkflowExecution.ID, timeout)

			*status = MutexStatusTimeout
		}

		logger.warn(handler.Info.WorkflowExecution.ID, "abort", "mutex: ignoring timeout ...",
			"mutex_id", handler.Info.WorkflowExecution.ID)
	}
}

func cleanup(ctx workflow.Context, handler *Handler, pool *Pool, fn workflow.Settable, logger *MutexLogger) func(workflow.Context) {
	return func(ctx workflow.Context) {
		for {
			rx := &Handler{}
			workflow.GetSignalChannel(ctx, WorkflowSignalCleanup.String()).Receive(ctx, rx)

			logger.info(rx.Info.WorkflowExecution.ID, "cleanup", "mutex: cleanup requested ...",
				"pool_size", pool.size(),
				"mutex_id", handler.Info.WorkflowExecution.ID)

			if pool.size() == 0 {
				logger.info(rx.Info.WorkflowExecution.ID, "cleanup", "mutex: requesting shutdown ...",
					"pool_size", pool.size(),
					"mutex_id", handler.Info.WorkflowExecution.ID)
				fn.Set(rx, nil)
				return
			}

			_ = workflow.
				SignalExternalWorkflow(ctx, rx.Info.WorkflowExecution.ID, "", WorkflowSignalCleanupDone.String(), false).
				Get(ctx, nil)

			workflow.GetSignalChannel(ctx, WorkflowSignalCleanupDoneAck.String()).Receive(ctx, nil)

			logger.info(rx.Info.WorkflowExecution.ID, "cleanup", "mutex: cleanup request processed!",
				"pool_size", pool.size(),
				"mutex_id", handler.Info.WorkflowExecution.ID)
		}
	}
}

func terminate(ctx workflow.Context, handler *Handler, logger *MutexLogger) shared.FutureHandler {
	return func(future workflow.Future) {
		rx := &Handler{}
		_ = future.Get(ctx, rx)

		logger.info(rx.Info.WorkflowExecution.ID, "terminate", "mutex: shutdown request received ...",
			"mutex_id", handler.Info.WorkflowExecution.ID)
		logger.info(rx.Info.WorkflowExecution.ID, "terminate", "mutex: shutting down ...",
			"mutex_id", handler.Info.WorkflowExecution.ID)
	}
}
