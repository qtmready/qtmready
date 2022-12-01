package auth

import (
	"net/http"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/entities"
	"go.breu.io/ctrlplane/internal/shared"
)

type (
	ServerHandler struct{}
)

func (s *ServerHandler) Register(ctx echo.Context) error {
	request := &RegisterationRequest{}

	// Translating request to json
	if err := ctx.Bind(request); err != nil {
		return err
	}

	// Validating request
	if err := ctx.Validate(request); err != nil {
		return err
	}

	// Validating team
	team := &entities.Team{Name: request.TeamName}
	if err := ctx.Validate(team); err != nil {
		return err
	}

	user := &entities.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     string(request.Email),
		Password:  request.Password,
	}
	if err := ctx.Validate(user); err != nil {
		return err
	}

	if err := db.Save(team); err != nil {
		return err
	}

	user.TeamID = team.ID
	if err := db.Save(user); err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, &RegisterationResponse{Team: team, User: user})
}

func (s *ServerHandler) Login(ctx echo.Context) error {
	request := &LoginRequest{}

	if err := ctx.Bind(request); err != nil {
		return err
	}

	if err := ctx.Validate(request); err != nil {
		return err
	}

	params := db.QueryParams{"email": "'" + string(request.Email) + "'"}
	user := &entities.User{}

	if err := db.Get(user, params); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	if user.VerifyPassword(request.Password) {
		access, _ := GenerateAccessToken(user.ID.String(), user.TeamID.String())
		refresh, _ := GenerateRefreshToken(user.ID.String(), user.TeamID.String())

		return ctx.JSON(http.StatusOK, &TokenResponse{AccessToken: &access, RefreshToken: &refresh})
	}

	return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
}

func (s *ServerHandler) CreateTeamAPIKey(ctx echo.Context) error {
	request := &CreateAPIKeyRequest{}

	if err := ctx.Bind(request); err != nil {
		return err
	}

	if err := ctx.Validate(request); err != nil {
		return err
	}

	id, _ := gocql.ParseUUID(ctx.Get("team_id").(string))
	guard := &entities.Guard{}
	key := guard.NewForTeam(id)

	if err := guard.Save(); err != nil {
		shared.Logger.Error("error saving guard", "error", err)
		return err
	}

	return ctx.JSON(http.StatusCreated, &CreateAPIKeyResponse{Key: &key})
}

func (s *ServerHandler) CreateUserAPIKey(ctx echo.Context) error {
	request := &CreateAPIKeyRequest{}

	if err := ctx.Bind(request); err != nil {
		return err
	}

	if err := ctx.Validate(request); err != nil {
		return err
	}

	id, _ := gocql.ParseUUID(ctx.Get("user_id").(string))
	guard := &entities.Guard{}
	key := guard.NewForUser(*request.Name, id)

	if err := guard.Save(); err != nil {
		shared.Logger.Error("error saving guard", "error", err)
		return err
	}

	return ctx.JSON(http.StatusCreated, &CreateAPIKeyResponse{Key: &key})
}

func (s *ServerHandler) ValidateAPIKey(ctx echo.Context) error {
	valid := "valid"
	return ctx.JSON(http.StatusOK, &ValidateAPIKeyResponse{Message: &valid})
}
