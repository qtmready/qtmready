package serializers

import (
	"encoding/json"
	"io"

	"go.breu.io/ctrlplane/internal/common"
	"go.breu.io/ctrlplane/internal/db/models"
)

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
		User *models.User `json:"user"`
		Team *models.Team `json:"team"`
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

func (r *RegistrationRequest) Reply(body io.ReadCloser) (RegistrationResponse, error) {
	response := RegistrationResponse{
		Team: &models.Team{Name: r.TeamName},
		User: &models.User{FirstName: r.FirstName, Email: r.Email, Password: r.Password, LastName: r.LastName},
	}

	if err := json.NewDecoder(body).Decode(r); err != nil {
		return RegistrationResponse{}, err
	}

	if err := common.Validator.Struct(r); err != nil {
		return response, err
	}

	if err := response.Team.Save(); err != nil {
		return response, err
	}

	response.User.TeamID = response.Team.ID

	if err := response.User.Save(); err != nil {
		response.User.SendEmail()
		return response, err
	}

	return response, nil
}

func (r *LoginRequest) Reply(body io.ReadCloser) (TokenResponse, error) {
	response := TokenResponse{}
	return response, nil
}
