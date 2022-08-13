package auth

import "go.breu.io/ctrlplane/internal/entities"

type (
	// Registration request serializer and validator
	RegistrationRequest struct {
		Email           string `json:"email" validate:"required,email"`
		Password        string `json:"password" validate:"required"`
		FirstName       string `json:"first_name" validate:"required"`
		LastName        string `json:"last_name" validate:"required"`
		ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
		TeamName        string `json:"team_name" validate:"required"`
	}

	// Registration Response Serializer
	RegistrationResponse struct {
		User *entities.User `json:"user"`
		Team *entities.Team `json:"team"`
	}

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
