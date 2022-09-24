// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 

package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/entities"
	"go.breu.io/ctrlplane/internal/shared"
)

func CreateRoutes(g *echo.Group, middlewares ...echo.MiddlewareFunc) {
	r := &AuthRoutes{}
	g.POST("/register", r.register)
	g.POST("/login", r.login)
}

type AuthRoutes struct{}

// register is a handler for /auth/register endpoint
func (routes *AuthRoutes) register(ctx echo.Context) error {
	request := &RegistrationRequest{}

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

	user := &entities.User{FirstName: request.FirstName, LastName: request.LastName, Email: request.Email, Password: request.Password}
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

	return ctx.JSON(http.StatusCreated, &RegistrationResponse{Team: team, User: user})
}

// login is a handler for /auth/login endpoint
func (routes *AuthRoutes) login(ctx echo.Context) error {
	request := &LoginRequest{}

	// Translating request to json
	if err := ctx.Bind(request); err != nil {
		return err
	}

	// Validating request
	if err := ctx.Validate(request); err != nil {
		return err
	}

	params := db.QueryParams{"email": "'" + request.Email + "'"}
	user := &entities.User{}

	if err := db.Get(user, params); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	if user.VerifyPassword(request.Password) {
		access, _ := shared.GenerateAccessToken(user.ID.String(), user.TeamID.String())
		refresh, _ := shared.GenerateRefreshToken(user.ID.String(), user.TeamID.String())

		return ctx.JSON(http.StatusOK, &TokenResponse{AccessToken: access, RefreshToken: refresh})
	}

	return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
}
