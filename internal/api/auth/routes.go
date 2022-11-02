// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
//
// This software is made available by Breu, Inc., under the terms of the 
// Breu Community License Agreement ("BCL Agreement"), version 1.0, found at  
// https://www.breu.io/license/community. By installating, downloading, 
// accessing, using or distrubting any of the software, you agree to the  
// terms of the license agreement. 
//
// The above copyright notice and the subsequent license agreement shall be 
// included in all copies or substantial portions of the software. 
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, 
// IMPLIED, STATUTORY, OR OTHERWISE, AND SPECIFICALLY DISCLAIMS ANY WARRANTY OF 
// MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE 
// SOFTWARE. 
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT 
// LIMITED TO, LOST PROFITS OR ANY CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, 
// OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, ARISING 
// OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY  
// APPLICABLE LAW. 

package auth

import (
	"net/http"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/entities"
	"go.breu.io/ctrlplane/internal/shared"
)

// CreateRoutes is for creating auth related routes.
func CreateRoutes(g *echo.Group, middlewares ...echo.MiddlewareFunc) {
	r := &Routes{}
	g.POST("/register", r.register)
	g.POST("/login", r.login)

	akg := g.Group("/api-keys", Middleware)
	akr := &APIKeyRoutes{}
	akg.POST("/team", akr.team)
	akg.POST("/user", akr.user)
	akg.GET("/validate", akr.validate)
}

type (
	Routes       struct{}
	APIKeyRoutes struct{}
)

// @Summary     Registers a new user.
// @Description Registers a new user.
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body     RegistrationRequest true "RegistrationRequest"
// @Success     201  {object} RegistrationResponse
// @Failure     400  {object} echo.HTTPError
// @Router      /auth/register [post]
//
// register is a handler for /auth/register endpoint.
func (routes *Routes) register(ctx echo.Context) error {
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

	user := &entities.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
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

	return ctx.JSON(http.StatusCreated, &RegistrationResponse{Team: team, User: user})
}

// @Summary     Get short lived JWT token along with a refresh token.
// @Description Get short lived JWT token along with a refresh token.
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body     LoginRequest true "LoginRequest"
// @Success     200  {object} TokenResponse
// @Failure     400  {object} echo.HTTPError
// @Failure     401  {object} echo.HTTPError
// @Router      /auth/login [post]
//
// login gets a short lived JWT token along with a refresh token.
func (routes *Routes) login(ctx echo.Context) error {
	request := &LoginRequest{}

	if err := ctx.Bind(request); err != nil {
		return err
	}

	if err := ctx.Validate(request); err != nil {
		return err
	}

	params := db.QueryParams{"email": "'" + request.Email + "'"}
	user := &entities.User{}

	if err := db.Get(user, params); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	if user.VerifyPassword(request.Password) {
		access, _ := GenerateAccessToken(user.ID.String(), user.TeamID.String())
		refresh, _ := GenerateRefreshToken(user.ID.String(), user.TeamID.String())

		return ctx.JSON(http.StatusOK, &TokenResponse{AccessToken: access, RefreshToken: refresh})
	}

	return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
}

// @Summary     Create a new API Key for team.
// @Description Create a new API Key for team.
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body     CreateAPIKeyRequest true "CreateAPIKeyRequest"
// @Success     201  {object} CreateAPIKeyResponse
// @Failure     400  {object} echo.HTTPError
// @Failure     401  {object} echo.HTTPError
// @Router      /auth/api-keys/team [post]
//
// team creates a new API Key for team.
func (routes *APIKeyRoutes) team(ctx echo.Context) error {
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

	return ctx.JSON(http.StatusCreated, &CreateAPIKeyResponse{Key: key})
}

// @Summary     Create a new API Key for user.
// @Description Create a new API Key for user.
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body     CreateAPIKeyRequest true "CreateAPIKeyRequest"
// @Success     201  {object} CreateAPIKeyResponse
// @Failure     400  {object} echo.HTTPError
// @Failure     401  {object} echo.HTTPError
// @Router      /auth/api-keys/user [post]
//
// user creates a new API Key for user.
func (routes *APIKeyRoutes) user(ctx echo.Context) error {
	request := &CreateAPIKeyRequest{}
	if err := ctx.Bind(request); err != nil {
		return err
	}

	if err := ctx.Validate(request); err != nil {
		return err
	}

	id, _ := gocql.ParseUUID(ctx.Get("user_id").(string))
	guard := &entities.Guard{}
	key := guard.NewForUser(request.Name, id)

	if err := guard.Save(); err != nil {
		shared.Logger.Error("error saving guard", "error", err)
		return err
	}

	return ctx.JSON(http.StatusCreated, &CreateAPIKeyResponse{Key: key})
}

// @Summary     Parses the header and validates the API Key.
// @Description Parses the header and validates the API Key.
// @Tags        auth
// @Accept      json
// @Produce     json
// @Failure     400 {object} echo.HTTPError
// @Failure     401 {object} echo.HTTPError
// @Router      /auth/api-keys/validate [get]
//
// validate an API Key.
func (routes *APIKeyRoutes) validate(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, &ValidateAPIKeyResponse{IsValid: true})
}
