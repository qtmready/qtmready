package core

import (
	"net/http"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/entities"
	"go.breu.io/ctrlplane/internal/shared"
)

// CreateRoutes creates the routes for the app
func CreateRoutes(g *echo.Group, middlewares ...echo.MiddlewareFunc) {
	apps := &AppRoutes{}
	g.POST("/apps", apps.Create)
	g.GET("/apps", apps.List)
	g.GET("/apps/:slug", apps.Get)

	g.POST("/apps/:slug/repos", apps.CreateAppRepo)
	g.GET("/apps/:slug/repos", apps.GetAppRepos)
}

type (
	AppRoutes struct{}
)

// Create creates a new app
func (a *AppRoutes) Create(ctx echo.Context) error {
	request := &AppCreateRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	teamID, _ := gocql.ParseUUID(shared.GetTeamIDFromContext(ctx))
	app := &entities.App{Name: request.Name, TeamID: teamID}
	if err := db.Save(app); err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, app)
}

// Get gets an app by slug
func (a *AppRoutes) Get(ctx echo.Context) error {
	app := &entities.App{}
	params := db.QueryParams{"slug": "'" + ctx.Param("slug") + "'"}

	if err := db.Get(app, params); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	return ctx.JSON(http.StatusOK, app)
}

// List lists all apps associated with the team
func (a *AppRoutes) List(ctx echo.Context) error {
	result := make([]entities.App, 0)
	params := db.QueryParams{"team_id": shared.GetTeamIDFromContext(ctx)}

	if err := db.Filter(&entities.App{}, &result, params); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}

// GetAppRepos gets an app repos by slug
func (a *AppRoutes) GetAppRepos(ctx echo.Context) error {
	result := make([]entities.AppRepo, 0)
	app := &entities.App{}

	params := db.QueryParams{"slug": "'" + ctx.Param("slug") + "'", "team_id": shared.GetTeamIDFromContext(ctx)}
	if err := db.Get(app, params); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "app not found")
	}

	params = db.QueryParams{"app_id": app.ID.String()}
	if err := db.Filter(&entities.AppRepo{}, &result, params); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, result)
}

// CreateAppRepo creates a new app repo
func (a *AppRoutes) CreateAppRepo(ctx echo.Context) error {
	request := &AppRepoCreateRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	app := &entities.App{}
	params := db.QueryParams{"slug": "'" + ctx.Param("slug") + "'", "team_id": shared.GetTeamIDFromContext(ctx)}
	if err := db.Get(app, params); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	switch request.Provider {
	case "github":
		return a.github(ctx, request, app)
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "unsupported git provider")
	}
}

func (a *AppRoutes) github(ctx echo.Context, request *AppRepoCreateRequest, app *entities.App) error {
	if err := db.Get(&entities.GithubRepo{}, db.QueryParams{"id": request.RepoID.String()}); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "repo not found")
	}

	repo := &entities.AppRepo{
		AppID:         app.ID,
		RepoID:        request.RepoID,
		DefaultBranch: request.DefaultBranch,
		IsMonorepo:    request.IsMonorepo,
		Provider:      request.Provider,
	}

	if err := db.Save(repo); err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, repo)
}
