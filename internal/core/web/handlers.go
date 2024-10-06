// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2023, 2024.
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

package web

import (
	"net/http"
	"strings"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/core/kernel"
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

// create core repo handler will create or update the repo with its message provider info (channel).
func (s *ServerHandler) CreateRepo(ctx echo.Context) error {
	request := &defs.RepoCreateRequest{}
	repo := &defs.Repo{}

	if err := ctx.Bind(request); err != nil {
		return err
	}

	teamID, _ := gocql.ParseUUID(ctx.Get("team_id").(string))
	data, err := kernel.Instance().
		RepoIO(request.Provider).
		GetProviderInfo(ctx.Request().Context(), request.CtrlID.String())

	if err != nil {
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	// Get the core repo by CtrlID if it exit update the record otherwise create the record
	err = db.Get(repo, db.QueryParams{"ctrl_id": request.CtrlID.String()})
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	if err != nil && strings.Contains(err.Error(), "not found") {
		repo = &defs.Repo{
			Name:                data.RepoName,
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
	} else {
		repo.IsMonorepo = request.IsMonorepo
		repo.Provider = request.Provider
		repo.Threshold = request.Threshold
		repo.MessageProvider = request.MessageProvider
		repo.MessageProviderData = request.MessageProviderData
		repo.StaleDuration = request.StaleDuration
	}

	if err := db.Save(repo); err != nil {
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	if err := kernel.Instance().
		RepoIO(request.Provider).
		SetEarlyWarning(ctx.Request().Context(), request.CtrlID.String(), true); err != nil {
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (s *ServerHandler) ListRepos(ctx echo.Context) error {
	repos := make([]defs.Repo, 0)
	params := db.QueryParams{"team_id": ctx.Get("team_id").(string)}

	if err := db.Filter(&defs.Repo{}, &repos, params); err != nil {
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, repos)
}

func (s *ServerHandler) GetRepo(ctx echo.Context, id string) error {
	repo := &defs.Repo{}
	params := db.QueryParams{"id": id}

	if err := db.Get(repo, params); err != nil {
		return shared.NewAPIError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, repo)
}
