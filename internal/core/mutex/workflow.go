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
	MutexState = string
)

var (
	LockAcquired MutexState = "acquired"
	LockReleased MutexState = "released"
	Timeout      MutexState = "timeout"
)

// Workflow is the mutex workflow. It is responsible for controlling the access to a resource. It should always be started as
// a child workflow.
// It works by listening to two signals:
//   - WorkflowSignalAcquire: this signal is sent by the caller to acquire the lock.
//   - WorkflowSignalRelease: this signal is sent by the caller to release the lock.
//
// FIXME: Re-evaluate Lock.Start() behavior.
//
// Currently, `Lock.Start()` starts the worker as a child workflow. While this ensures termination when the parent workflow finishes,
// it conflicts with the purpose of blocking on a resource. Although Temporal's idempotency key prevents duplicate workflows for the
// same resource, subsequent workflows attempting to start the child workflow will encounter errors.
//
// Should we instead use a signal-based approach for starting the lock workflow? This would allow other parts of the core system
// to acquire and wait on the lock. However, a graceful cleanup mechanism would be necessary.
func Workflow(ctx workflow.Context, timeout time.Duration) error {
	var (
		acquirer, releaser string
		state              MutexState
	)

	logger := workflow.GetLogger(ctx)
	logger.Info("mutex: workflow started with ...", "workflow ID", workflow.GetInfo(ctx).WorkflowExecution.ID)

	// selector := workflow.NewSelector(ctx)
	acquireCh := workflow.GetSignalChannel(ctx, WorkflowSignalAcquire.String())
	releaseCh := workflow.GetSignalChannel(ctx, WorkflowSignalRelease.String())

	for {
		state = LockReleased

		logger.Info("mutex: waiting for acquire lock signal ...")
		acquireCh.Receive(ctx, &acquirer)
		logger.Info("mutex: lock acquired. signaling acquirer ....", "acquired by", acquirer)

		if err := workflow.SignalExternalWorkflow(ctx, acquirer, "", WorkflowSignalLocked.String(), nil).Get(ctx, nil); err != nil {
			logger.Error("mutex: failed to signal acquirer.", "error", err)
			continue
		}

		for {
			logger.Info("mutex: waiting for release lock signal ...")

			selector := workflow.NewSelector(ctx)
			selector.AddReceive(releaseCh, onRelease(ctx, state, &acquirer, &releaser))
			selector.AddFuture(workflow.NewTimer(ctx, timeout), onTimeOut(ctx, state, timeout))

			selector.Select(ctx)

			if state == LockReleased || state == Timeout {
				break
			}
		}
	}
}

func onRelease(ctx workflow.Context, state MutexState, acquirer, releaser *string) shared.ChannelHandler {
	logger := workflow.GetLogger(ctx)

	return func(channel workflow.ReceiveChannel, more bool) {
		channel.Receive(ctx, releaser)

		if *releaser == *acquirer {
			logger.Info("mutex: releasing lock ....", "released by", releaser)

			state = LockReleased
		}
	}
}

func onTimeOut(ctx workflow.Context, state MutexState, timeout time.Duration) shared.FutureHandler {
	logger := workflow.GetLogger(ctx)

	return func(future workflow.Future) {
		logger.Info("mutex: lock timeout reached, releasing lock ...", "timeout", timeout)

		state = Timeout
	}
}
