package serializers

import (
	"encoding/json"
	"errors"
	"io"

	"go.breu.io/ctrlplane/internal/common"
	"go.breu.io/ctrlplane/internal/db/models"
	"go.uber.org/zap"
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

// Composing reply for RegistrationRequest
func (request *RegistrationRequest) Reply(body io.ReadCloser) (RegistrationResponse, error) {
	if err := json.NewDecoder(body).Decode(request); err != nil {
		return RegistrationResponse{}, err
	}

	response := RegistrationResponse{
		Team: &models.Team{Name: request.TeamName},
		User: &models.User{FirstName: request.FirstName, Email: request.Email, Password: request.Password, LastName: request.LastName},
	}

	if err := common.Validator.Struct(request); err != nil {
		common.Logger.Error(err.Error())
		return response, err
	}

	if err := response.Team.Save(); err != nil {
		common.Logger.Error(err.Error())
		return response, err
	}

	response.User.TeamID = response.Team.ID

	if err := response.User.Save(); err != nil {
		common.Logger.Info("User ...", zap.Any("user", response.User))
		common.Logger.Error(err.Error())
		return response, err
	}

	return response, nil
}

// Composing Request for LoginRequest
func (request *LoginRequest) Reply(body io.ReadCloser) (TokenResponse, error) {
	response := TokenResponse{}

	if err := json.NewDecoder(body).Decode(request); err != nil {
		return response, err
	}

	if err := common.Validator.Struct(request); err != nil {
		return response, err
	}
	params := map[string]interface{}{"email": request.Email}
	user := models.User{}
	if err := user.Get(params); err != nil {
		return response, err
	}

	if user.VerifyPassword(request.Password) {
		payload := map[string]interface{}{"user_id": user.ID, "team_id": user.TeamID}
		_, response.Token, _ = common.JWT.Encode(payload)
		return response, nil
	}

	err := errors.New("invalid credentials")
	return response, err
}
