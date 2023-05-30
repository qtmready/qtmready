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

// Workflow is the mutex workflow. It is responsible for controlling the access to a resource. It should always be started as
// a child workflow.
// It works by listening to two signals:
//   - acquire: this signal is sent by the caller to acquire the lock.
//   - release: this signal is sent by the caller to release the lock.
//
// The workflow will block until the lock is acquired. Once acquired, it will block until the lock is released.
func Workflow(ctx workflow.Context, timout time.Duration) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("mutex: workflow started with ...", "workflow ID", workflow.GetInfo(ctx).WorkflowExecution.ID)

	// selector := workflow.NewSelector(ctx)
	acquirer := ""
	releaser := ""
	acquireCh := workflow.GetSignalChannel(ctx, WorkflowSignalAcquire.String())
	releaseCh := workflow.GetSignalChannel(ctx, WorkflowSignalRelease.String())

	for {
		logger.Info("mutex: waiting for acquire lock signal ...")
		acquireCh.Receive(ctx, &acquirer)
		logger.Info("mutex: lock acquired. signaling acquirer ....", "acquired by", acquirer)

		if err := workflow.SignalExternalWorkflow(ctx, acquirer, "", WorkflowSignalLocked.String(), nil).Get(ctx, nil); err != nil {
			logger.Error("mutex: failed to signal acquirer.", "error", err)
			continue
		}

		for {
			logger.Info("mutex: waiting for release lock signal ...")
			releaseCh.Receive(ctx, &releaser)

			if releaser == acquirer {
				logger.Info("mutex: releasing lock ....", "released by", releaser)
				break
			}

			logger.Info("mutex: lock is acquired by another workflow")
		}
	}
}
