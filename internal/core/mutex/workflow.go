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

	_ = state.set_query_state(ctx)

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

	_ = workflow.Sleep(ctx, 500*time.Millisecond) // NOTE: This is a hack to wait for the signal from the cleanup loop.

	state.logger.info(state.Handler.Info.WorkflowExecution.ID, "shutdown", "shutdown!")

	return nil
}
