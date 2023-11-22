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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

func handleLabelEvent(ctx echo.Context) error {
	shared.Logger().Info("handleLabelEvent")

	body, _ := io.ReadAll(ctx.Request().Body)
	ctx.Request().Body = io.NopCloser(bytes.NewBuffer(body))

	// Decode the string into a map[string]interface{}
	var requestBody map[string]interface{}
	if err := json.Unmarshal(body, &requestBody); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	pr_id := requestBody["number"].(float64)

	repo_info, _ := requestBody["repository"].(map[string]interface{})
	repo_id := repo_info["id"].(float64)

	curr_label_info, _ := requestBody["label"].(map[string]interface{})
	label := curr_label_info["name"].(string)

	shared.Logger().Info("PR info", "repo-id", repo_id)
	shared.Logger().Info("PR info", "PR-id", pr_id)
	shared.Logger().Info("PR info", "label", label)

	if label == "quantm ready" {
		// Construct the API endpoint for merging the pull request
		pr_id_str := strconv.FormatFloat(pr_id, 'f', -1, 64)
		repo_id_str := strconv.FormatFloat(repo_id, 'f', -1, 64)
		url := fmt.Sprintf("https://api.github.com/repositories/%s/pulls/%s/merge", repo_id_str, pr_id_str)

		// Prepare the request body (optional parameters can be included here)
		requestBody := map[string]interface{}{
			"commit_title": "Merge PR via Quantum",
		}

		// Convert the request body to JSON
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			// return err
		}

		// Create a new HTTP client
		client := &http.Client{}

		// Create a new POST request
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			// return err
		}

		//TODO:
		accessToken := "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" //TODO: need to add access token using env or some other way

		// Set the request headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken)

		// Perform the request
		resp, err := client.Do(req)
		if err != nil {
			// return err
		}
		defer resp.Body.Close()

		// Check the response status
		if resp.StatusCode != http.StatusOK {
			// return fmt.Errorf("Failed to merge pull request. Status: %s", resp.Status)
		}

		fmt.Println("Pull request merged successfully.")
	}

	return ctx.JSON(http.StatusOK, &WorkflowResponse{})
}

// handlePullRequestEvent handles GitHub pull request event.
func handlePullRequestEvent(ctx echo.Context) error {
	payload := &PullRequestEvent{}
	if err := ctx.Bind(payload); err != nil {
		shared.Logger().Error("unable to bind payload ...", "error", err)
		return err
	}

	w := &Workflows{}
	opts := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(strconv.FormatInt(payload.Installation.ID, 10)),
			shared.WithWorkflowElement("repo"),
			shared.WithWorkflowElementID(strconv.FormatInt(payload.Repository.ID, 10)),
			shared.WithWorkflowMod(WebhookEventPullRequest.String()),
			shared.WithWorkflowModID(strconv.FormatInt(payload.PullRequest.ID, 10)),
		)

	switch payload.Action {
	case "opened":
		exe, err := shared.Temporal().
			Client().
			ExecuteWorkflow(context.Background(), opts, w.OnPullRequestEvent, payload)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return ctx.JSON(http.StatusCreated, &WorkflowResponse{RunID: exe.GetRunID(), Status: WorkflowStatusQueued})
	default:
		err := shared.Temporal().
			Client().
			SignalWorkflow(context.Background(), opts.ID, "", WebhookEventPullRequest.String(), payload)
		if err != nil {
			shared.Logger().Error("unable to signal ...", "options", opts, "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

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
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, &WorkflowResponse{RunID: exe.GetID(), Status: WorkflowStatusQueued})
}
