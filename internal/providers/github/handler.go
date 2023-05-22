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

	"go.breu.io/ctrlplane/internal/auth"
	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/shared"
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

	payload := &CompleteInstallationSignal{request.InstallationId, request.SetupAction, teamID}

	workflows := &Workflows{}
	opts := shared.Temporal().
		Queue(shared.ProvidersQueue).
		GetWorkflowOptions("github", strconv.Itoa(int(payload.InstallationID)), WebhookEventInstallation.String())

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

func (s *ServerHandler) GithubWebhook(ctx echo.Context) error {
	signature := ctx.Request().Header.Get("X-Hub-Signature-256")

	if signature == "" {
		return ctx.JSON(http.StatusUnauthorized, ErrMissingHeaderGithubSignature)
	}

	// NOTE: We are reading the request body twice. This is not ideal.
	body, _ := io.ReadAll(ctx.Request().Body)
	ctx.Request().Body = io.NopCloser(bytes.NewBuffer(body))

	if err := Github.VerifyWebhookSignature(body, signature); err != nil {
		return shared.NewAPIError(http.StatusUnauthorized, err)
	}

	headerEvent := ctx.Request().Header.Get("X-GitHub-Event")
	if headerEvent == "" {
		return shared.NewAPIError(http.StatusBadRequest, ErrMissingHeaderGithubEvent)
	}

	shared.Logger().Debug("headerEvent", "headerEvent", headerEvent)

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
