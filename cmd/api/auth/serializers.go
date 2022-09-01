package auth

import "go.breu.io/ctrlplane/internal/entities"

type (
	// RegistrationRequest is the http request for user registration
	RegistrationRequest struct {
		Email           string `json:"email" validate:"required,email"`
		Password        string `json:"password" validate:"required"`
		FirstName       string `json:"first_name" validate:"required"`
		LastName        string `json:"last_name" validate:"required"`
		ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
		TeamName        string `json:"team_name" validate:"required"`
	}

	// RegistrationResponse is the http response after user registration is done
	RegistrationResponse struct {
		User *entities.User `json:"user"`
		Team *entities.Team `json:"team"`
	}

	// LoginRequest is the http request for login
	LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	RefreshTokenRequest struct {
		Token string `json:"token"`
	}

	TokenResponse struct {
		Token string `json:"token"`
	}
)
