// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
//
// This software is made available by Breu, Inc., under the terms of the Breu  
// Community License Agreement, Version 1.0 located at  
// http://www.breu.io/breu-community-license/v1. BY INSTALLING, DOWNLOADING,  
// ACCESSING, USING OR DISTRIBUTING ANY OF THE SOFTWARE, YOU AGREE TO THE TERMS  
// OF SUCH LICENSE AGREEMENT. 

package auth

import (
	"go.breu.io/ctrlplane/internal/entities"
)

type (
	// RegistrationRequest is the http request for user registration.
	RegistrationRequest struct {
		Email           string `json:"email" validate:"required,email"`
		Password        string `json:"password" validate:"required"`
		FirstName       string `json:"first_name" validate:"required"`
		LastName        string `json:"last_name" validate:"required"`
		ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
		TeamName        string `json:"team_name" validate:"required"`
	}

	// RegistrationResponse is the http response after user registration is done.
	RegistrationResponse struct {
		User *entities.User `json:"user"`
		Team *entities.Team `json:"team"`
	}

	// LoginRequest is the http request for login.
	LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	RefreshTokenRequest struct {
		Token string `json:"token"`
	}

	TokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	// TODO: This is generic request for all API's to include pagination and sorting.
	RequestBuilder[T any] struct {
		Data  *T
		Error *error
	}

	CreateAPIKeyRequest struct {
		Name string `json:"name" validate:"required"`
	}

	CreateAPIKeyResponse struct {
		Key string `json:"key"`
	}

	ValidateAPIKeyResponse struct {
		IsValid bool `json:"is_valid"`
	}
)
