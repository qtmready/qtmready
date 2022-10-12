// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package core

import (
	"errors"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/entities"
)

var (
	ErrNotFound            = errors.New("not found")
	ErrUnsupportedProvider = errors.New("unsupported provider")
)

// CreateRoutes creates the routes for the app.
func CreateRoutes(g *echo.Group, middlewares ...echo.MiddlewareFunc) {
	apps := &AppRoutes{}
	g.POST("/apps", apps.create)
	g.GET("/apps", apps.list)
	g.GET("/apps/:slug", apps.get)

	repos := &AppRepoRoutes{}

	g.POST("/apps/:slug/repos", repos.create)
	g.GET("/apps/:slug/repos", repos.list)
}

type (
	AppRoutes     struct{}
	AppRepoRoutes struct{}
)

// @Summary     List all apps for a team.
// @Description List all apps for a team.
// @Tags        core
// @Accept      json
// @Produce     json
// @Success     200 {array}  entities.App
// @Failure     400 {object} echo.HTTPError
// @Router      /apps [get]
//
// list all apps associated with the team.
func (routes *AppRoutes) list(ctx echo.Context) error {
	result := make([]entities.App, 0)
	params := db.QueryParams{"team_id": ctx.Get("team_id").(string)}

	if err := db.Filter(&entities.App{}, &result, params); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}

// @Summary     Create a new app.
// @Description Create a new app.
// @Tags        core
// @Accept      json
// @Produce     json
// @Param       body body     AppCreateRequest true "AppCreateRequest"
// @Success     201  {object} entities.App
// @Failure     400  {object} echo.HTTPError
// @Router      /apps [post]
//
// create a new app.
func (routes *AppRoutes) create(ctx echo.Context) error {
	request := &AppCreateRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	teamID, _ := gocql.ParseUUID(ctx.Get("team_id").(string))
	app := &entities.App{Name: request.Name, TeamID: teamID}

	if err := db.Save(app); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, app)
}

// @Summary     Get an app by slug.
// @Description Get an app by slug.
// @Tags        core
// @Accept      json
// @Produce     json
// @Param       slug path     string true "App slug"
// @Success     200  {object} entities.App
// @Failure     400  {object} echo.HTTPError
// @Router      /apps/{slug} [get]
//
// get an app by slug.
func (routes *AppRoutes) get(ctx echo.Context) error {
	app := &entities.App{}
	params := db.QueryParams{"slug": "'" + ctx.Param("slug") + "'"}

	if err := db.Get(app, params); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, ErrNotFound)
	}

	return ctx.JSON(http.StatusOK, app)
}

// @Summary     List all repos given an app.
// @Description List all repos given an app.
// @Tags        core
// @Accept      json
// @Produce     json
// @Param       slug path     string true "App slug"
// @Success     200  {array}  entities.Repo
// @Failure     400  {object} echo.HTTPError
// @Router      /apps/{slug}/repos [get]
//
// list all repos associated with an app.
func (routes *AppRepoRoutes) list(ctx echo.Context) error {
	result := make([]entities.Repo, 0)
	app := &entities.App{}

	params := db.QueryParams{"slug": "'" + ctx.Param("slug") + "'", "team_id": ctx.Get("team_id").(string)}
	if err := db.Get(app, params); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, ErrNotFound)
	}

	params = db.QueryParams{"app_id": app.ID.String()}
	if err := db.Filter(&entities.Repo{}, &result, params); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}

// @Summary     Create a new repo for an app.
// @Description Create a new repo for an app.
// @Tags        core
// @Accept      json
// @Produce     json
// @Param       slug path     string               true "App slug"
// @Param       body body     AppRepoCreateRequest true "AppRepoCreateRequest"
// @Success     201  {object} entities.Repo
// @Failure     400  {object} echo.HTTPError
// @Router      /apps/{slug}/repos [post]
//
// create a new repo for an app.
func (routes *AppRepoRoutes) create(ctx echo.Context) error {
	request := &AppRepoCreateRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	app := &entities.App{}
	params := db.QueryParams{"slug": "'" + ctx.Param("slug") + "'", "team_id": ctx.Get("team_id").(string)}

	if err := db.Get(app, params); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, ErrNotFound)
	}

	switch request.Provider {
	case "github":
		return routes.github(ctx, request, app)
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, ErrUnsupportedProvider)
	}
}

// github creates associates an app with a github repo.
func (routes *AppRepoRoutes) github(ctx echo.Context, request *AppRepoCreateRequest, app *entities.App) error {
	if err := db.Get(&entities.GithubRepo{}, db.QueryParams{"id": request.RepoID.String()}); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, ErrNotFound)
	}

	repo := &entities.Repo{
		AppID:         app.ID,
		RepoID:        request.RepoID,
		DefaultBranch: request.DefaultBranch,
		IsMonorepo:    request.IsMonorepo,
		Provider:      request.Provider,
	}

	if err := db.Save(repo); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, repo)
}
