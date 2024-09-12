// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2022, 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.


package github

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"go.breu.io/quantm/internal/shared"
)

// handleInstallationEvent handles GitHub App installation event.
func handleInstallationEvent(ctx echo.Context) error {
	payload := &InstallationEvent{}
	if err := ctx.Bind(payload); err != nil {
		return err
	}

	shared.Logger().Info("installation event received ...", "action", payload.Action)

	workflows := &Workflows{}
	opts := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(payload.Installation.ID.String()),
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

// handlePushEvent handles GitHub push/create event.
func handlePushEvent(ctx echo.Context) error {
	payload := &PushEvent{}
	if err := ctx.Bind(payload); err != nil {
		shared.Logger().Error("unable to bind payload ...", "error", err)
		return err
	}

	shared.Logger().Info("repo event received ...", "action", payload.Installation.ID)

	// the value will be `NoCommit` if we have a tag push.
	// TODO: handle tag.
	if payload.After == NoCommit {
		return ctx.JSON(http.StatusNoContent, nil)
	}

	w := &Workflows{}
	event := WebhookEvent(ctx.Request().Header.Get("X-GitHub-Event"))
	delievery := ctx.Request().Header.Get("X-GitHub-Delivery")
	opts := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(payload.Installation.ID.String()),
			shared.WithWorkflowElement("repo"),
			shared.WithWorkflowElementID(payload.Repository.ID.String()),
			shared.WithWorkflowMod(event.String()),
			shared.WithWorkflowModID(delievery),
		)

	_, err := shared.Temporal().Client().ExecuteWorkflow(context.Background(), opts, w.OnPushEvent, payload)
	if err != nil {
		shared.Logger().Error("unable to signal OnPushEvent ...", "options", opts, "error", err)
		return nil
	}

	return ctx.NoContent(http.StatusNoContent)
}

// handlePushEvent handles GitHub push/create event.
func handleCreateOrDeleteEvent(ctx echo.Context) error {
	payload := &CreateOrDeleteEvent{}
	if err := ctx.Bind(payload); err != nil {
		shared.Logger().Error("unable to bind payload ...", "error", err)
		return err
	}

	shared.Logger().Info("repo event received ...", "installation", payload.Installation.ID)

	w := &Workflows{}

	event := WebhookEvent(ctx.Request().Header.Get("X-GitHub-Event"))
	if event == WebhookEventCreate {
		payload.IsCreated = true
	}

	delievery := ctx.Request().Header.Get("X-GitHub-Delivery")
	opts := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(payload.Installation.ID.String()),
			shared.WithWorkflowElement("repo"),
			shared.WithWorkflowElementID(payload.Repository.ID.String()),
			shared.WithWorkflowMod(event.String()),
			shared.WithWorkflowModID(delievery),
		)

	_, err := shared.Temporal().Client().ExecuteWorkflow(context.Background(), opts, w.OnCreateOrDeleteEvent, payload)
	if err != nil {
		shared.Logger().Error("unable to signal OnPushEvent ...", "options", opts, "error", err)
		return nil
	}

	return ctx.NoContent(http.StatusNoContent)
}

// handleReleaseEvent handles GitHub release event.
func handlePullRequestEvent(ctx echo.Context) error {
	payload := &PullRequestEvent{}
	if err := ctx.Bind(payload); err != nil {
		shared.Logger().Error("unable to bind payload ...", "error", err)
		return err
	}

	shared.Logger().Info("pull request event received ...", "action", payload.Action)

	w := &Workflows{}
	delievery := ctx.Request().Header.Get("X-GitHub-Delivery")
	opts := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(payload.Installation.ID.String()),
			shared.WithWorkflowElement("repo"),
			shared.WithWorkflowElementID(payload.Repository.ID.String()),
			shared.WithWorkflowMod(WebhookEventPullRequest.String()),
			shared.WithWorkflowModID(payload.PullRequest.Number.String()),
			shared.WithWorkflowProp("action", payload.Action),
			shared.WithWorkflowProp("id", delievery),
		)

	_, err := shared.Temporal().Client().ExecuteWorkflow(context.Background(), opts, w.OnPullRequestEvent, payload)
	if err != nil {
		shared.Logger().Error("unable to signal OnPullRequestEvent ...", "options", opts, "error", err)
		return nil
	}

	return ctx.NoContent(http.StatusNoContent)
}

// handleInstallationRepositoriesEvent handles GitHub installation repositories event.
func handleInstallationRepositoriesEvent(ctx echo.Context) error {
	payload := &InstallationRepositoriesEvent{}
	if err := ctx.Bind(payload); err != nil {
		return err
	}

	shared.Logger().Info("installation repositories event received...", "action", payload.Action)

	w := &Workflows{}
	opts := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(payload.Installation.ID.String()),
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

// handlePullRequestReviewEvent handles GitHub pull request review event.
func handlePullRequestReviewEvent(ctx echo.Context) error {
	payload := &PullRequestReviewEvent{}
	if err := ctx.Bind(payload); err != nil {
		shared.Logger().Error("unable to bind payload ...", "error", err)
		return err
	}

	shared.Logger().Info("pull request review event received ...", "action", payload.Action)

	w := &Workflows{}
	delievery := ctx.Request().Header.Get("X-GitHub-Delivery")
	opts := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(payload.Installation.ID.String()),
			shared.WithWorkflowElement("repo"),
			shared.WithWorkflowElementID(payload.Repository.ID.String()),
			shared.WithWorkflowMod(WebhookEventPullRequestReview.String()),
			shared.WithWorkflowModID(payload.PullRequest.Number.String()),
			shared.WithWorkflowProp("action", payload.Action),
			shared.WithWorkflowProp("id", delievery),
		)

	_, err := shared.Temporal().Client().ExecuteWorkflow(context.Background(), opts, w.OnPullRequestReviewEvent, payload)
	if err != nil {
		shared.Logger().Error("unable to signal OnPullRequestReviewEvent ...", "options", opts, "error", err)
		return nil
	}

	return ctx.NoContent(http.StatusNoContent)
}

// handlePullRequestReviewCommentEvent handles GitHub pull request review comment event.
func handlePullRequestReviewCommentEvent(ctx echo.Context) error {
	payload := &PullRequestReviewCommentEvent{}
	if err := ctx.Bind(payload); err != nil {
		shared.Logger().Error("unable to bind payload ...", "error", err)
		return err
	}

	shared.Logger().Info("pull request review comment event received ...", "action", payload.Action)

	w := &Workflows{}
	delievery := ctx.Request().Header.Get("X-GitHub-Delivery")
	opts := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(payload.Installation.ID.String()),
			shared.WithWorkflowElement("repo"),
			shared.WithWorkflowElementID(payload.Repository.ID.String()),
			shared.WithWorkflowMod(WebhookEventPullRequestReview.String()),
			shared.WithWorkflowModID(payload.PullRequest.Number.String()),
			shared.WithWorkflowProp("action", payload.Action),
			shared.WithWorkflowProp("id", delievery),
		)

	_, err := shared.Temporal().Client().ExecuteWorkflow(context.Background(), opts, w.OnPullRequestReviewCommentEvent, payload)
	if err != nil {
		shared.Logger().Error("unable to signal OnPullRequestReviewCommentEvent ...", "options", opts, "error", err)
		return nil
	}

	return ctx.NoContent(http.StatusNoContent)
}

func handleWorkflowRunEvent(ctx echo.Context) error {
	shared.Logger().Debug("workflow-run event received.")

	payload := &GithubWorkflowRunEvent{}
	if err := ctx.Bind(payload); err != nil {
		shared.Logger().Error("unable to bind payload ...", "error", err)
		return err
	}

	shared.Logger().Info("workflow run event received ...", "action", payload.Action)

	workflows := &Workflows{}
	opts := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(payload.Repository.Name),
			shared.WithWorkflowElement("workflow_run"),
			shared.WithWorkflowElementID(payload.WR.ID.String()),
		)

	exe, err := shared.Temporal().Client().ExecuteWorkflow(
		ctx.Request().Context(),
		opts,
		workflows.OnWorkflowRunEvent,
		payload,
	)
	if err != nil {
		shared.Logger().Error("unable to signal OnWorkflowRunEvent ...", "options", opts, "error", err)
		return nil
	}

	return ctx.JSON(http.StatusCreated, &WorkflowResponse{RunID: exe.GetRunID(), Status: WorkflowStatusQueued})
}
