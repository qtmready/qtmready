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
	"context"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"

	"go.breu.io/ctrlplane/internal/auth"
	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/shared"
)

type (
	ServerHandler struct {
		*auth.SecurityHandler
	}
)

// NewServerHandler creates a new ServerHandler.
func NewServerHandler(security echo.MiddlewareFunc) *ServerHandler {
	return &ServerHandler{
		SecurityHandler: &auth.SecurityHandler{Middleware: security},
	}
}

func (s *ServerHandler) CreateStack(ctx echo.Context) error {
	request := &StackCreateRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	teamID, _ := gocql.ParseUUID(ctx.Get("team_id").(string))
	stack := &Stack{Name: request.Name, TeamID: teamID}

	if err := db.Save(stack); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	/*
		start infinite stack workflow
		reason for going with infinite workflow instead of starting with signal is to follow the
		temporal guideline which state that workflow ids should not be resued
	*/
	w := &Workflows{}
	opts := shared.Temporal.Queues[shared.CoreQueue].
		GetWorkflowOptions("core", "stack", request.Name, "stackId", stack.ID.String())

	exe, err := shared.Temporal.Client.ExecuteWorkflow(context.Background(), opts, w.OnPullRequestWorkflow, request.Name)
	if err != nil {
		// TODO: remove stack if workflow not started? or always start this workflow with signal so it can be started on pull request (if not already running)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	shared.Logger.Info("started workflow: ", opts.ID, " run ID: ", exe.GetRunID())

	return ctx.JSON(http.StatusCreated, stack)
}

func (s *ServerHandler) ListStacks(ctx echo.Context) error {
	stacks := make([]Stack, 0)
	params := db.QueryParams{"team_id": ctx.Get("team_id").(string)}

	if err := db.Filter(&Stack{}, &stacks, params); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, stacks)
}

func (s *ServerHandler) GetStack(ctx echo.Context) error {
	stack := &Stack{}
	params := db.QueryParams{"slug": "'" + ctx.Param("slug") + "'", "team_id": ctx.Get("team_id").(string)}

	if err := db.Get(stack, params); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, stack)
}

func (s *ServerHandler) CreateRepo(ctx echo.Context) error {
	request := &RepoCreateRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	stack := &Repo{
		StackID:       request.StackID,
		ProviderID:    request.ProviderID,
		DefaultBranch: request.DefaultBranch,
		IsMonorepo:    request.IsMonorepo,
		Provider:      request.Provider,
	}

	if err := db.Save(stack); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, stack)
}

func (s *ServerHandler) ListRepos(ctx echo.Context) error {
	repos := make([]Repo, 0)
	params := db.QueryParams{"team_id": ctx.Get("team_id").(string)}

	if err := db.Filter(&Repo{}, &repos, params); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, repos)
}

func (s *ServerHandler) GetRepo(ctx echo.Context) error {
	repo := &Repo{}
	params := db.QueryParams{"id": "'" + ctx.Param("id") + "'", "team_id": ctx.Get("team_id").(string)}

	if err := db.Get(repo, params); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, repo)
}
