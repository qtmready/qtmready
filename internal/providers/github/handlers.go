// Copyright © 2022, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLATING, DOWNLOADING, ACCESSING, USING OR DISTRUBTING ANY OF
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

package github

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/shared"
)

// handleInstallationEvent handles GitHub App installation event.
func handleInstallationEvent(ctx echo.Context) error {
	payload := &InstallationEventPayload{}
	if err := ctx.Bind(payload); err != nil {
		return err
	}

	shared.Logger.Info("installation event received ...")

	workflows := &Workflows{}
	opts := shared.Temporal.
		Queues[shared.ProvidersQueue].
		GetWorkflowOptions("github", strconv.FormatInt(payload.Installation.ID, 10), InstallationEvent.String())

	exe, err := shared.Temporal.Client.SignalWithStartWorkflow(
		ctx.Request().Context(),
		opts.ID,
		WebhookInstallationEventSignal.String(),
		payload,
		opts,
		workflows.OnInstall,
	)
	if err != nil {
		shared.Logger.Error("unable to signal ...", "options", opts, "error", err)
		return nil
	}

	shared.Logger.Debug("installation event handled ...", "options", opts, "execution", exe.GetRunID())

	return ctx.JSON(http.StatusCreated, &WorkflowResponse{RunID: exe.GetRunID(), Status: WorkflowQueued})
}

// handlePushEvent handles GitHub push event.
func handlePushEvent(ctx echo.Context) error {
	payload := &PushEventPayload{}
	if err := ctx.Bind(payload); err != nil {
		shared.Logger.Error("unable to bind payload ...", "error", err)
		return err
	}

	// the value will be `NoCommit` if we have a tag push, or squash merge.
	if payload.After == NoCommit {
		return ctx.JSON(http.StatusOK, &WorkflowResponse{RunID: db.NullUUID, Status: WorkflowSkipped})
	}

	w := &Workflows{}
	opts := shared.Temporal.
		Queues[shared.ProvidersQueue].
		GetWorkflowOptions(
			"github",
			strconv.FormatInt(payload.Installation.ID, 10),
			"repo",
			strconv.FormatInt(payload.Repository.ID, 10),
			PushEvent.String(),
			"ref",
			payload.After)

	exe, err := shared.Temporal.Client.ExecuteWorkflow(context.Background(), opts, w.OnPush, payload)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, &WorkflowResponse{RunID: exe.GetRunID(), Status: WorkflowQueued})
}

// handlePullRequestEvent handles GitHub pull request event.
func handlePullRequestEvent(ctx echo.Context) error {
	payload := PullRequestEventPayload{}
	if err := ctx.Bind(&payload); err != nil {
		shared.Logger.Error("unable to bind payload ...", "error", err)
		return err
	}

	w := &Workflows{}
	opts := shared.Temporal.
		Queues[shared.ProvidersQueue].
		GetWorkflowOptions(
			"github",
			strconv.FormatInt(payload.Installation.ID, 10),
			"repo",
			strconv.FormatInt(payload.Repository.ID, 10),
			PullRequestEvent.String(),
			strconv.FormatInt(payload.PullRequest.ID, 10),
		)

	switch payload.Action {
	case "opened":
		exe, err := shared.Temporal.Client.ExecuteWorkflow(context.Background(), opts, w.OnPullRequest, payload)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusCreated, &WorkflowResponse{RunID: exe.GetRunID(), Status: WorkflowQueued})
	default:
		err := shared.Temporal.Client.SignalWorkflow(context.Background(), opts.ID, "", PullRequestEvent.String(), payload)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, &WorkflowResponse{RunID: db.NullUUID, Status: WorkflowSkipped})
	}
}
