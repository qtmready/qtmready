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
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"go.breu.io/quantm/internal/shared/queues"
)

// handleInstallationEvent handles GitHub App installation event.
func handleInstallationEvent(ctx echo.Context) error {
	payload := &InstallationEvent{}
	if err := ctx.Bind(payload); err != nil {
		return err
	}

	slog.Info("installation event received ...", "action", payload.Action)

	workflows := &Workflows{}
	exe, err := queues.Providers().SignalWithStartWorkflow(
		ctx.Request().Context(),
		// TODO: we should probably take the action from the payload.
		InstallationWebhookWorkflowOptions(payload.Installation.ID, "install"),
		WorkflowSignalInstallationEvent,
		payload,
		workflows.OnInstallationEvent,
	)

	if err != nil {
		slog.Error("unable to signal ...", "error", err)
		return err
	}

	slog.Debug("installation event handled ...", "execution", exe.GetRunID())

	return ctx.JSON(http.StatusCreated, &WorkflowResponse{RunID: exe.GetRunID(), Status: WorkflowStatusQueued})
}

// handleInstallationRepositoriesEvent handles GitHub installation repositories event.
func handleInstallationRepositoriesEvent(ctx echo.Context) error {
	payload := &InstallationRepositoriesEvent{}
	if err := ctx.Bind(payload); err != nil {
		return err
	}

	slog.Info("installation repositories event received...", "action", payload.Action)

	w := &Workflows{}

	exe, err := queues.Providers().ExecuteWorkflow(
		ctx.Request().Context(),
		InstallationWebhookWorkflowOptions(payload.Installation.ID, WebhookEventInstallationRepositories.String()),
		w.OnInstallationRepositoriesEvent,
		payload,
	)
	if err != nil {
		slog.Error("error dispatching workflow", "error", err)
		return err
	}

	return ctx.JSON(http.StatusOK, &WorkflowResponse{RunID: exe.GetID(), Status: WorkflowStatusQueued})
}

// handlePushEvent handles GitHub push/create event.
func handlePushEvent(ctx echo.Context) error {
	payload := &PushEvent{}
	if err := ctx.Bind(payload); err != nil {
		slog.Error("unable to bind payload ...", "error", err)
		return err
	}

	slog.Info("repo event received ...", "action", payload.Installation.ID)

	// the value will be `NoCommit` if we have a tag push.
	// TODO: handle tag.
	if payload.After == NoCommit {
		return ctx.NoContent(http.StatusNoContent)
	}

	w := &Workflows{}
	event := WebhookEvent(ctx.Request().Header.Get("X-GitHub-Event"))
	delievery := ctx.Request().Header.Get("X-GitHub-Delivery")

	_, err := queues.Providers().ExecuteWorkflow(
		ctx.Request().Context(),
		RepoWebhookWorkflowOptions(payload.Installation.ID, payload.Repository.Name, event.String(), delievery),
		w.OnPushEvent,
		payload,
	)
	if err != nil {
		slog.Error("error dispatching workflow", "error", err, "event", event)
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

// handlePushEvent handles GitHub push/create event.
func handleCreateOrDeleteEvent(ctx echo.Context) error {
	payload := &CreateOrDeleteEvent{}
	if err := ctx.Bind(payload); err != nil {
		slog.Error("unable to bind payload ...", "error", err)
		return err
	}

	slog.Info("repo event received ...", "installation", payload.Installation.ID)

	w := &Workflows{}

	event := WebhookEvent(ctx.Request().Header.Get("X-GitHub-Event"))
	if event == WebhookEventCreate {
		payload.IsCreated = true
	}

	delievery := ctx.Request().Header.Get("X-GitHub-Delivery")

	_, err := queues.Providers().ExecuteWorkflow(
		ctx.Request().Context(),
		RepoWebhookWorkflowOptions(payload.Installation.ID, payload.Repository.Name, event.String(), delievery),
		w.OnCreateOrDeleteEvent,
		payload,
	)
	if err != nil {
		slog.Error("error dispatching workflow", "error", err, "event", event)
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

// handleReleaseEvent handles GitHub release event.
func handlePullRequestEvent(ctx echo.Context) error {
	payload := &PullRequestEvent{}
	if err := ctx.Bind(payload); err != nil {
		slog.Error("unable to bind payload ...", "error", err)
		return err
	}

	slog.Info("pull request event received ...", "action", payload.Action)

	w := &Workflows{}
	delievery := ctx.Request().Header.Get("X-GitHub-Delivery")

	_, err := queues.Providers().ExecuteWorkflow(
		ctx.Request().Context(),
		RepoWebhookWorkflowOptions(payload.Installation.ID, payload.Repository.Name, WebhookEventPullRequest.String(), delievery),
		w.OnPullRequestEvent,
		payload,
	)
	if err != nil {
		slog.Error("error dispatching workflow", "error", err, "event", WebhookEventPullRequest)
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

// handlePullRequestReviewEvent handles GitHub pull request review event.
func handlePullRequestReviewEvent(ctx echo.Context) error {
	payload := &PullRequestReviewEvent{}
	if err := ctx.Bind(payload); err != nil {
		slog.Error("unable to bind payload ...", "error", err)
		return err
	}

	slog.Info("pull request review event received ...", "action", payload.Action)

	w := &Workflows{}
	delievery := ctx.Request().Header.Get("X-GitHub-Delivery")

	_, err := queues.Providers().ExecuteWorkflow(
		ctx.Request().Context(),
		RepoWebhookWorkflowOptions(payload.Installation.ID, payload.Repository.Name, WebhookEventPullRequestReview.String(), delievery),
		w.OnPullRequestReviewEvent,
		payload,
	)
	if err != nil {
		slog.Error("error dispatching workflow", "error", err, "event", WebhookEventPullRequestReview)
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

// handlePullRequestReviewCommentEvent handles GitHub pull request review comment event.
func handlePullRequestReviewCommentEvent(ctx echo.Context) error {
	payload := &PullRequestReviewCommentEvent{}
	if err := ctx.Bind(payload); err != nil {
		slog.Error("unable to bind payload ...", "error", err)
		return err
	}

	slog.Info("pull request review comment event received ...", "action", payload.Action)

	w := &Workflows{}
	delievery := ctx.Request().Header.Get("X-GitHub-Delivery")

	_, err := queues.Providers().ExecuteWorkflow(
		ctx.Request().Context(),
		RepoWebhookWorkflowOptions(
			payload.Installation.ID, payload.Repository.Name, WebhookEventPullRequestReviewComment.String(), delievery,
		),
		w.OnPullRequestReviewCommentEvent,
		payload,
	)
	if err != nil {
		slog.Error("error dispatching workflow", "error", err, "event", WebhookEventPullRequestReviewComment)
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func handleWorkflowRunEvent(ctx echo.Context) error {
	slog.Debug("workflow-run event received.")

	payload := &GithubWorkflowRunEvent{}
	if err := ctx.Bind(payload); err != nil {
		slog.Error("unable to bind payload ...", "error", err)
		return err
	}

	slog.Info("workflow run event received ...", "action", payload.Action)

	workflows := &Workflows{}
	delievery := ctx.Request().Header.Get("X-GitHub-Delivery")

	exe, err := queues.Providers().ExecuteWorkflow(
		ctx.Request().Context(),
		RepoWebhookWorkflowOptions(payload.Installation.ID, payload.Repository.Name, WebhookEventWorkflowRun.String(), delievery),
		workflows.OnWorkflowRunEvent,
		payload,
	)
	if err != nil {
		slog.Error("error dispatching workflow", "error", err, "event", WebhookEventWorkflowRun)
		return err
	}

	return ctx.JSON(http.StatusCreated, &WorkflowResponse{RunID: exe.GetRunID(), Status: WorkflowStatusQueued})
}
