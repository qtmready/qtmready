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

package auth

import (
	"net/http"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/shared"
)

type (
	// SecurityHandler is the base security handler for the API. It is meant to be embedded in other handlers.
	//
	// Usage:
	//  package {name}
	//
	//  import "go.breu.io/ctrlplane/internal/auth"
	//
	//  type ServerHandler struct {
	//    *auth.SecurityHandler
	//  }
	SecurityHandler struct{ Middleware echo.MiddlewareFunc }
	ServerHandler   struct{ *SecurityHandler } // ServerHandler for auth
)

// NewServerHandler creates a new ServerHandler.
func NewServerHandler(security echo.MiddlewareFunc) *ServerHandler {
	return &ServerHandler{
		SecurityHandler: &SecurityHandler{Middleware: security},
	}
}

// SecureHandler wraps the handler with the security middleware.
func (s *SecurityHandler) SecureHandler(handler echo.HandlerFunc, ctx echo.Context) error {
	err := s.Middleware(handler)(ctx)
	return err
}

// Register registers a new user.
func (s *ServerHandler) Register(ctx echo.Context) error {
	shared.Logger.Debug("investigating code path")
	request := &RegisterationRequest{}

	// Translating request to json
	if err := ctx.Bind(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	// Validating request
	if err := ctx.Validate(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	// Validating team
	team := &Team{Name: request.TeamName}
	if err := ctx.Validate(team); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	user := &User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	}
	if err := ctx.Validate(user); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	if err := db.Save(team); err != nil {
		// TODO: cleanup created team.
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	user.TeamID = team.ID
	if err := db.Save(user); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusCreated, &RegisterationResponse{Team: team, User: user})
}

// Login returns JWT tokens if the email & password are correct.
func (s *ServerHandler) Login(ctx echo.Context) error {
	request := &LoginRequest{}

	if err := ctx.Bind(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	if err := ctx.Validate(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	params := db.QueryParams{"email": "'" + request.Email + "'"}
	user := &User{}

	if err := db.Get(user, params); err != nil {
		return shared.NewAPIError(http.StatusNotFound, ErrInvalidCredentials)
	}

	if user.VerifyPassword(request.Password) {
		access, _ := GenerateAccessToken(user.ID.String(), user.TeamID.String())
		refresh, _ := GenerateRefreshToken(user.ID.String(), user.TeamID.String())

		return ctx.JSON(http.StatusOK, &TokenResponse{AccessToken: &access, RefreshToken: &refresh})
	}

	return shared.NewAPIError(http.StatusUnauthorized, ErrInvalidCredentials)
}

// CreateTeamAPIKey creates an API key for the Team.
func (s *ServerHandler) CreateTeamAPIKey(ctx echo.Context) error {
	request := &CreateAPIKeyRequest{}

	if err := ctx.Bind(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	if err := ctx.Validate(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	id, _ := gocql.ParseUUID(ctx.Get("team_id").(string))
	guard := &Guard{}
	key := guard.NewForTeam(id)

	if err := guard.Save(); err != nil {
		shared.Logger.Error("error saving guard", "error", err)
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusCreated, &CreateAPIKeyResponse{Key: &key})
}

// CreateUserAPIKey creates an API Key for the User.
func (s *ServerHandler) CreateUserAPIKey(ctx echo.Context) error {
	request := &CreateAPIKeyRequest{}

	if err := ctx.Bind(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	if err := ctx.Validate(request); err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	id, _ := gocql.ParseUUID(ctx.Get("user_id").(string))
	guard := &Guard{}
	key := guard.NewForUser(*request.Name, id)

	if err := guard.Save(); err != nil {
		shared.Logger.Error("error saving guard", "error", err)
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusCreated, &CreateAPIKeyResponse{Key: &key})
}

// ValidateAPIKey validates the X-API-KEY header and return a boolean. This is mainly required for otel collector.
func (s *ServerHandler) ValidateAPIKey(ctx echo.Context) error {
	valid := "valid"
	return ctx.JSON(http.StatusOK, &APIKeyValidationResponse{Message: &valid})
}

func (s *ServerHandler) ListTeams(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}

func (s *ServerHandler) GetTeam(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}

func (s *ServerHandler) CreateTeam(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}

func (s *ServerHandler) AddUserToTeam(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, nil)
}
