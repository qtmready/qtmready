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
	"time"

	"go.breu.io/ctrlplane/internal/shared"
	"go.temporal.io/sdk/workflow"
)

const (
	unLockTimeOutStackMutex             time.Duration = time.Minute * 30 //TODO: adjust this
	OnPullRequestWorkflowPRSignalsLimit               = 1000             // TODO: adjust this
)

type (
	Activities struct{}
	Workflows  struct {
	}
)

// ChangesetController controls the rollout lifecycle for one changeset.
func (w *Workflows) ChangesetController(id string) error {
	return nil
}

// this workflow will be started on stack creation

// activity to start mutex workflow
// wait on signal for pr with repo id

// acquire lock on stack
// get stack from repo id
// get repo from stack
// compute changeset idempotency key
// signal sentinal to start orchestration
func (w *Workflows) OnPullRequestWorkflow(ctx workflow.Context, stackID string, stackName string) error {

	logger := workflow.GetLogger(ctx)
	currentWorkflowID := workflow.GetInfo(ctx).WorkflowExecution.ID
	pullRequestSignalName := currentWorkflowID
	resourceID := currentWorkflowID

	// execute activity to start mutex workflow
	logger.Info("executing SignalWithStartMutexWorkflowActivity")

	m := NewMutex(currentWorkflowID, resourceID, unLockTimeOutStackMutex)
	m.Init(ctx)

	payload := &shared.PullRequestSignal{}
	var prSignalsCounter int = 0
	for {
		// return continue as new if this workflow has processes signals upto a limit
		if prSignalsCounter >= OnPullRequestWorkflowPRSignalsLimit {
			return workflow.NewContinueAsNewError(ctx, w.OnPullRequestWorkflow, stackID, stackName)
		}

		// Wait for PR event
		workflow.GetSignalChannel(ctx, pullRequestSignalName).Receive(ctx, payload)
		prSignalsCounter++

		logger.Info("Pull request signal received from Github Workflow:", payload.SenderWorkflowID)
		unlockFunc, err := m.Lock(ctx)
		if err != nil {
			logger.Info("Error in acquiring lock", err)
		}

		time.Sleep(time.Second * 5)
		unlockFunc()
	}
	return nil
}
