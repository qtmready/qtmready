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

func (s *ServerHandler) CreateRepo(ctx echo.Context) error {
	request := &RepoCreateRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	teamID, _ := gocql.ParseUUID(ctx.Get("team_id").(string))
	data, err := Instance().
		RepoIO(request.Provider).
		GetRepoData(ctx.Request().Context(), request.CtrlID.String())

	if err != nil {
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	repo := &Repo{
		Name:                data.Name,
		DefaultBranch:       data.DefaultBranch,
		IsMonorepo:          request.IsMonorepo,
		Provider:            request.Provider,
		ProviderID:          data.ProviderID,
		CtrlID:              request.CtrlID,
		Threshold:           request.Threshold,
		TeamID:              teamID,
		MessageProvider:     request.MessageProvider,
		MessageProviderData: request.MessageProviderData,
		StaleDuration:       request.StaleDuration,
	}

	if err := db.Save(repo); err != nil {
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	if err := Instance().
		RepoIO(request.Provider).
		SetEarlyWarning(ctx.Request().Context(), request.CtrlID.String(), true); err != nil {
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (s *ServerHandler) ListRepos(ctx echo.Context) error {
	repos := make([]Repo, 0)
	params := db.QueryParams{"team_id": ctx.Get("team_id").(string)}

	if err := db.Filter(&Repo{}, &repos, params); err != nil {
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, repos)
}

func (s *ServerHandler) GetRepo(ctx echo.Context, id string) error {
	repo := &Repo{}
	params := db.QueryParams{"id": id}

	if err := db.Get(repo, params); err != nil {
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, repo)
}
