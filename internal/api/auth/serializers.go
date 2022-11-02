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
