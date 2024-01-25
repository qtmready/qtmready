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

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
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

func (s *ServerHandler) CreateBlueprint(ctx echo.Context) error {
	request := &BlueprintCreateRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	blueprint := &Blueprint{
		Name:           request.Name,
		Regions:        request.Regions,
		StackID:        request.StackID,
		RolloutBudget:  request.RolloutBudget,
		ProviderConfig: request.ProviderConfig,
	}

	if err := db.Save(blueprint); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, blueprint)
}

func (s *ServerHandler) GetBlueprint(ctx echo.Context) error {
	blueprint := &Blueprint{}
	params := db.QueryParams{"stack_id": ctx.Param("stack_id")}

	if err := db.Get(blueprint, params); err != nil {
		return shared.NewAPIError(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, blueprint)
}

func (s *ServerHandler) CreateWorkload(ctx echo.Context) error {
	request := &WorkloadCreateRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	workload := &Workload{
		Name:       request.Name,
		Kind:       request.Kind,
		RepoID:     request.RepoID,
		RepoPath:   request.RepoPath,
		StackID:    request.StackID,
		ResourceID: request.ResourceID,
		Container:  request.Container,
	}

	if err := db.Save(workload); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, workload)
}

func (s *ServerHandler) GetWorkload(ctx echo.Context) error {
	stackid := ctx.QueryParam("stack_id")
	repoid := ctx.QueryParam("repo_id")

	if repoid != "" {
		workload := &Workload{}
		params := db.QueryParams{"repo_id": repoid}

		if err := db.Get(workload, params); err != nil {
			return shared.NewAPIError(http.StatusNotFound, err)
		}

		return ctx.JSON(http.StatusOK, workload)
	} else if stackid != "" {
		workloads := make([]Workload, 0)
		params := db.QueryParams{"stack_id": stackid}

		if err := db.Filter(&Workload{}, &workloads, params); err != nil {
			return shared.NewAPIError(http.StatusInternalServerError, err)
		}

		return ctx.JSON(http.StatusOK, workloads)
	}

	return ctx.JSON(http.StatusBadRequest, "No query paramters provided")
}

func (s *ServerHandler) CreateResource(ctx echo.Context) error {
	request := &ResourceCreateRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	resource := &Resource{
		Name:     request.Name,
		Provider: request.Provider,
		StackID:  request.StackID,
		// TODO: check why assinging directly wasn't working
		Driver:      Driver(request.Driver),
		Config:      request.Config,
		IsImmutable: request.Immutable,
	}

	if err := db.Save(resource); err != nil {
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, resource)
}

func (s *ServerHandler) GetResource(ctx echo.Context) error {
	resources := make([]Resource, 0)
	params := db.QueryParams{"stack_id": ctx.Param("stack_id")}

	if err := db.Filter(&Resource{}, &resources, params); err != nil {
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, resources)
}

func (s *ServerHandler) CreateStack(ctx echo.Context) error {
	request := &StackCreateRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	teamID, _ := gocql.ParseUUID(ctx.Get("team_id").(string))
	stack := &Stack{Name: request.Name, TeamID: teamID}

	if err := db.Save(stack); err != nil {
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	/*
		start infinite stack workflow
		reason for going with infinite workflow instead of starting with signal is to follow the
		temporal guideline which state that workflow ids should not be resued
	*/
	w := &Workflows{}
	opts := shared.Temporal().
		Queue(shared.CoreQueue).
		WorkflowOptions(
			shared.WithWorkflowBlock("stack"),
			shared.WithWorkflowBlockID(stack.ID.String()),
		)

	exe, err := shared.Temporal().Client().ExecuteWorkflow(context.Background(), opts, w.StackController, stack.ID.String())
	if err != nil {
		// TODO: remove stack if workflow not started? or always start this workflow with signal so it can be started on pull request
		// (if not already running)
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	shared.Logger().Info("started workflow: ", "ID", opts.ID, " Run ID: ", exe.GetRunID())

	return ctx.JSON(http.StatusCreated, stack)
}

func (s *ServerHandler) ListStacks(ctx echo.Context) error {
	stacks := make([]Stack, 0)
	params := db.QueryParams{"team_id": ctx.Get("team_id").(string)}

	if err := db.Filter(&Stack{}, &stacks, params); err != nil {
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, stacks)
}

func (s *ServerHandler) GetStack(ctx echo.Context) error {
	stack := &Stack{}
	params := db.QueryParams{"slug": "'" + ctx.Param("slug") + "'", "team_id": ctx.Get("team_id").(string)}

	if err := db.Get(stack, params); err != nil {
		return shared.NewAPIError(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, stack)
}

func (s *ServerHandler) CreateRepo(ctx echo.Context) error {
	request := &RepoCreateRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	repo := &Repo{
		Name:          request.Name,
		StackID:       request.StackID,
		ProviderID:    request.ProviderID,
		DefaultBranch: request.DefaultBranch,
		IsMonorepo:    request.IsMonorepo,
		Provider:      request.Provider,
	}

	if err := db.Save(repo); err != nil {
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, repo)
}

func (s *ServerHandler) ListRepos(ctx echo.Context) error {
	repos := make([]Repo, 0)
	params := db.QueryParams{"team_id": ctx.Get("team_id").(string)}

	if err := db.Filter(&Repo{}, &repos, params); err != nil {
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, repos)
}

func (s *ServerHandler) GetRepo(ctx echo.Context) error {
	repo := &Repo{}
	params := db.QueryParams{"id": ctx.Param("id")}

	if err := db.Get(repo, params); err != nil {
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, repo)
}
