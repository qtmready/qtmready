package core

import (
	"net/http"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"

	"go.breu.io/ctrlplane/internal/auth"
	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/entity"
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
	app := &entity.Stack{Name: request.Name, TeamID: teamID}

	if err := db.Save(app); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, app)
}

func (s *ServerHandler) ListStacks(ctx echo.Context) error {
	stacks := make([]entity.Stack, 0)
	params := db.QueryParams{"team_id": ctx.Get("team_id").(string)}

	if err := db.Filter(&entity.Stack{}, &stacks, params); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, stacks)
}

func (s *ServerHandler) GetStack(ctx echo.Context) error {
	stack := &entity.Stack{}
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

	stackID, _ := gocql.ParseUUID(request.StackId.String())
	app := &entity.Repo{
		StackID:       stackID,
		ProviderID:    request.ProviderId,
		DefaultBranch: request.DefaultBranch,
		IsMonorepo:    request.IsMonorepo,
		Provider:      request.Provider,
	}

	if err := db.Save(app); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, app)
}

func (s *ServerHandler) ListRepos(ctx echo.Context) error {
	repos := make([]entity.Repo, 0)
	params := db.QueryParams{"team_id": ctx.Get("team_id").(string)}

	if err := db.Filter(&entity.Repo{}, &repos, params); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, repos)
}

func (s *ServerHandler) GetRepo(ctx echo.Context) error {
	repo := &entity.Repo{}
	params := db.QueryParams{"id": "'" + ctx.Param("id") + "'", "team_id": ctx.Get("team_id").(string)}

	if err := db.Get(repo, params); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, repo)
}
