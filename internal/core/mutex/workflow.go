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
)

// MutexWorkflow is the mutex workflow. It controls access to a resource.
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
func MutexWorkflow(ctx workflow.Context, state *MutexState) error {
	state.restore(ctx)

	shutdown, shutdownfn := workflow.NewFuture(ctx)

	workflow.Go(ctx, state.on_prepare(ctx))
	workflow.Go(ctx, state.on_cleanup(ctx, shutdownfn))

	for state.Persist {
		state.logger.info(state.Handler.Info.WorkflowExecution.ID, "main_loop", "waiting for lock request ...")

		found := true
		acquirer := workflow.NewSelector(ctx)

		acquirer.AddReceive(workflow.GetSignalChannel(ctx, WorkflowSignalAcquire.String()), state.on_acquire(ctx))
		acquirer.AddFuture(shutdown, state.on_terminate(ctx))

		acquirer.Select(ctx)

		if !state.Persist {
			break // Shutdown signal received
		}

		if !found {
			state.logger.info(state.Handler.Info.WorkflowExecution.ID, "main_loop", "lock not found in the pool, retrying ...")
			state.to_acquiring(ctx)

			continue
		}

		state.logger.info(state.Handler.Info.WorkflowExecution.ID, "main_loop", "lock acquired!")
		state.to_locked(ctx)

		for {
			state.logger.info(state.Handler.Info.WorkflowExecution.ID, "main_loop", "waiting for release or timeout ...")

			releaser := workflow.NewSelector(ctx)

			releaser.AddReceive(
				workflow.GetSignalChannel(ctx, WorkflowSignalRelease.String()),
				state.on_release(ctx),
			)
			releaser.AddFuture(workflow.NewTimer(ctx, state.Timeout), state.on_abort(ctx))

			releaser.Select(ctx)

			if state.Status == MutexStatusReleased || state.Status == MutexStatusTimeout {
				state.to_acquiring(ctx)
				break
			}
		}
	}

	_ = workflow.Sleep(ctx, 500*time.Millisecond)

	state.logger.info(state.Handler.Info.WorkflowExecution.ID, "shutdown", "shutdown!")

	return nil
}
