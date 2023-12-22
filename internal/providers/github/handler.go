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
	"io"
	"net/http"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

type (
	ServerHandler struct{ *auth.SecurityHandler }
)

// NewServerHandler creates a new ServerHandler.
func NewServerHandler(middleware echo.MiddlewareFunc) *ServerHandler {
	return &ServerHandler{
		SecurityHandler: &auth.SecurityHandler{Middleware: middleware},
	}
}

func (s *ServerHandler) GithubCompleteInstallation(ctx echo.Context) error {
	request := &CompleteInstallationRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	teamID, err := gocql.ParseUUID(ctx.Get("team_id").(string))

	if err != nil {
		return err
	}

	payload := &CompleteInstallationSignal{request.InstallationID, request.SetupAction, teamID}

	workflows := &Workflows{}
	opts := shared.Temporal().
		Queue(shared.ProvidersQueue).
		// GetWorkflowOptions("github", strconv.Itoa(int(payload.InstallationID)), WebhookEventInstallation.String())
		WorkflowOptions(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(strconv.Itoa(int(payload.InstallationID))),
			shared.WithWorkflowElement(WebhookEventInstallation.String()),
		)

	exe, err := shared.Temporal().
		Client().
		SignalWithStartWorkflow(
			ctx.Request().Context(),
			opts.ID,
			WorkflowSignalCompleteInstallation.String(),
			payload,
			opts,
			workflows.OnInstallationEvent,
		)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, &WorkflowResponse{RunID: exe.GetID(), Status: WorkflowStatusQueued})
}

func (s *ServerHandler) GithubGetRepos(ctx echo.Context) error {
	result := make([]Repo, 0)
	if err := db.Filter(
		&Repo{},
		&result,
		db.QueryParams{"team_id": ctx.Get("team_id").(string)},
	); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}

func (s *ServerHandler) GithubGetInstallations(ctx echo.Context) error {
	result := make([]Installation, 0)
	if err := db.Filter(
		&Installation{},
		&result,
		db.QueryParams{"team_id": ctx.Get("team_id").(string)},
	); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}

// GithubArtifactReady API is called by github action after building and pushing the build artifact
// After receiving pull request webhook, Quantum waits for the artifact ready event to start deployment.
func (s *ServerHandler) GithubArtifactReady(ctx echo.Context) error {
	request := &ArtifactReadyRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	workflowID := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowID(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(request.InstallationID),
			shared.WithWorkflowElement("repo"),
			shared.WithWorkflowElementID(request.RepoID),
			shared.WithWorkflowMod(WebhookEventPullRequest.String()),
			shared.WithWorkflowModID(request.PullRequestID),
		)

	payload := &ArtifactReadySignal{Image: request.Image, Digest: request.Digest, Registry: request.Registry.String()}

	err := shared.Temporal().Client().SignalWorkflow(ctx.Request().Context(), workflowID, "", WorkflowSignalArtifactReady.String(), payload)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, &WorkflowResponse{RunID: workflowID, Status: WorkflowStatusSignaled})
}

func (s *ServerHandler) GithubActionResult(ctx echo.Context) error {
	shared.Logger().Info("GithubActionResult method triggered.")

	request := &GithubActionResultRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	shared.Logger().Debug("GithubActionResult", "request", request)

	if request.Branch == "main" {
		shared.Logger().Info("GithubActionResult", "action", "No action needed")
		return nil
	}

	workflowID := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowID(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(request.RepoID),
			shared.WithWorkflowElement("branch"),
			shared.WithWorkflowElementID(request.Branch),
		)

	result := make([]Repo, 0)
	if err := db.Filter(
		&Repo{},
		&result,
		db.QueryParams{"github_id": request.RepoID},
	); err != nil {
		return err
	}

	installationID := result[0].InstallationID
	payload := &GithubActionResult{
		Branch:         request.Branch,
		RepoID:         request.RepoID,
		RepoName:       request.RepoName,
		RepoOwner:      request.RepoOwner,
		InstallationID: installationID,
	}

	workflows := &Workflows{}
	opts := shared.Temporal().
		Queue(shared.ProvidersQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("github"),
			shared.WithWorkflowBlockID(strconv.FormatInt(installationID, 10)),
			shared.WithWorkflowElement("repo"),
			shared.WithWorkflowElementID(request.RepoID),
		)

	_, err := shared.Temporal().Client().SignalWithStartWorkflow(
		ctx.Request().Context(),
		opts.ID,
		WorkflowSignalActionResult.String(),
		payload,
		opts,
		workflows.OnGithubActionResult,
		payload,
	)
	if err != nil {
		shared.Logger().Error("unable to signal ...", "options", opts, "error", err)
		return nil
	}

	return ctx.JSON(http.StatusOK, &WorkflowResponse{RunID: workflowID, Status: WorkflowStatusSignaled})
}

func (s *ServerHandler) GithubWebhook(ctx echo.Context) error {
	signature := ctx.Request().Header.Get("X-Hub-Signature-256")

	if signature == "" {
		return ctx.JSON(http.StatusUnauthorized, ErrMissingHeaderGithubSignature)
	}

	// NOTE: We are reading the request body twice. This is not ideal.
	body, _ := io.ReadAll(ctx.Request().Body)
	ctx.Request().Body = io.NopCloser(bytes.NewBuffer(body))

	if err := Instance().VerifyWebhookSignature(body, signature); err != nil {
		return shared.NewAPIError(http.StatusUnauthorized, err)
	}

	headerEvent := ctx.Request().Header.Get("X-GitHub-Event")
	if headerEvent == "" {
		return shared.NewAPIError(http.StatusBadRequest, ErrMissingHeaderGithubEvent)
	}

	shared.Logger().Debug("GithubWebhook", "headerEvent", headerEvent)
	// Uncomment for debugging!
	// var jsonMap map[string]interface{}
	// json.Unmarshal([]byte(string(body)), &jsonMap)
	// shared.Logger().Debug("GithubWebhook", "body", jsonMap)

	event := WebhookEvent(headerEvent)
	handlers := WebhookEventHandlers{
		WebhookEventInstallation:             handleInstallationEvent,
		WebhookEventInstallationRepositories: handleInstallationRepositoriesEvent,
		WebhookEventPush:                     handlePushEvent,
		WebhookEventPullRequest:              handlePullRequestEvent,
	}

	if handle, exists := handlers[event]; exists {
		return handle(ctx)
	}

	return shared.NewAPIError(http.StatusBadRequest, ErrInvalidEvent)
}
