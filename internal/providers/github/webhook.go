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

package github

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

// handleInstallationEvent handles GitHub App installation event.
func handleInstallationEvent(ctx echo.Context) error {
	payload := &InstallationEvent{}
	if err := ctx.Bind(payload); err != nil {
		return err
	}

	shared.Logger().Info("installation event received ...")

	workflows := &Workflows{}
	opts := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(strconv.FormatInt(payload.Installation.ID, 10)),
			shared.WithWorkflowElement(WebhookEventInstallation.String()),
		)

	exe, err := shared.Temporal().Client().SignalWithStartWorkflow(
		ctx.Request().Context(),
		opts.ID,
		WorkflowSignalInstallationEvent.String(),
		payload,
		opts,
		workflows.OnInstallationEvent,
	)
	if err != nil {
		shared.Logger().Error("unable to signal ...", "options", opts, "error", err)
		return nil
	}

	shared.Logger().Debug("installation event handled ...", "options", opts, "execution", exe.GetRunID())

	return ctx.JSON(http.StatusCreated, &WorkflowResponse{RunID: exe.GetRunID(), Status: WorkflowStatusQueued})
}

// handlePushEvent handles GitHub push event.
func handlePushEvent(ctx echo.Context) error {
	payload := &PushEvent{}
	if err := ctx.Bind(payload); err != nil {
		shared.Logger().Error("unable to bind payload ...", "error", err)
		return err
	}

	// the value will be `NoCommit` if we have a tag push, or squash merge.
	if payload.After == NoCommit {
		return ctx.JSON(http.StatusOK, &WorkflowResponse{RunID: db.NullUUID, Status: WorkflowStatusSkipped})
	}

	w := &Workflows{}
	opts := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(strconv.FormatInt(payload.Installation.ID, 10)),
			shared.WithWorkflowElement("repo"),
			shared.WithWorkflowElementID(strconv.FormatInt(payload.Repository.ID, 10)),
			shared.WithWorkflowMod(WebhookEventPush.String()),
			shared.WithWorkflowProp("ref", payload.Ref),
		)

	exe, err := shared.Temporal().Client().ExecuteWorkflow(context.Background(), opts, w.OnPushEvent, payload)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, &WorkflowResponse{RunID: exe.GetRunID(), Status: WorkflowStatusQueued})
}

// handlePullRequestEvent handles GitHub pull request event.
func handlePullRequestEvent(ctx echo.Context) error {
	payload := &PullRequestEvent{}
	if err := ctx.Bind(payload); err != nil {
		shared.Logger().Error("unable to bind payload ...", "error", err)
		return err
	}

	shared.Logger().Info("handlePullRequestEvent executing...")

	w := &Workflows{}
	opts := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(strconv.FormatInt(payload.Installation.ID, 10)),
			shared.WithWorkflowElement("repo"),
			shared.WithWorkflowElementID(strconv.FormatInt(payload.Repository.ID, 10)),
			shared.WithWorkflowMod(WebhookEventPullRequest.String()),
			shared.WithWorkflowModID(fmt.Sprintf(
				"%s-%s",
				strconv.FormatInt(payload.PullRequest.ID, 10),
				payload.Action,
			)),
		)

	switch payload.Action {
	case "opened":
		shared.Logger().Info("PR", "status", "open")
		exe, err := shared.Temporal().
			Client().
			ExecuteWorkflow(context.Background(), opts, w.OnPullRequestEvent, payload)

		if err != nil {
			return shared.NewAPIError(http.StatusInternalServerError, err)
		}

		return ctx.JSON(http.StatusCreated, &WorkflowResponse{RunID: exe.GetRunID(), Status: WorkflowStatusQueued})

	case "labeled":
		shared.Logger().Info("PR", "status", "label")
		exe, err := shared.Temporal().
			Client().
			ExecuteWorkflow(context.Background(), opts, w.OnLabelEvent, payload)

		if err != nil {
			return shared.NewAPIError(http.StatusInternalServerError, err)
		}

		return ctx.JSON(http.StatusCreated, &WorkflowResponse{RunID: exe.GetRunID(), Status: WorkflowStatusQueued})

	default:
		shared.Logger().Debug("handlePullRequestEvent default closing...")
		// err := shared.Temporal().
		// 	Client().
		// 	SignalWorkflow(context.Background(), opts.ID, "", WebhookEventPullRequest.String(), payload)
		// if err != nil {
		// 	shared.Logger().Error("unable to signal ...", "options", opts, "error", err)
		// 	return shared.NewAPIError(http.StatusInternalServerError, err)
		// }

		return ctx.JSON(http.StatusOK, &WorkflowResponse{RunID: db.NullUUID, Status: WorkflowStatusSkipped})
	}
}

// handleInstallationRepositoriesEvent handles GitHub installation repositories event.
func handleInstallationRepositoriesEvent(ctx echo.Context) error {
	payload := &InstallationRepositoriesEvent{}
	if err := ctx.Bind(payload); err != nil {
		return err
	}

	w := &Workflows{}
	opts := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(strconv.FormatInt(payload.Installation.ID, 10)),
			shared.WithWorkflowElement(WebhookEventInstallationRepositories.String()),
		)

	exe, err := shared.Temporal().
		Client().
		ExecuteWorkflow(context.Background(), opts, w.OnInstallationRepositoriesEvent, payload)
	if err != nil {
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, &WorkflowResponse{RunID: exe.GetID(), Status: WorkflowStatusQueued})
}
