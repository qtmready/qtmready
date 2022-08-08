package serializers

import (
	"encoding/json"
	"errors"
	"io"

	"go.breu.io/ctrlplane/internal/common"
	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/db/entities"
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
		User entities.User `json:"user"`
		Team entities.Team `json:"team"`
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
		return RegistrationResponse{Team: entities.Team{}, User: entities.User{}}, err
	}

	team := entities.Team{Name: request.TeamName}
	user := entities.User{Email: request.Email, FirstName: request.FirstName, LastName: request.LastName, Password: request.Password}

	if err := common.Validator.Struct(request); err != nil {
		common.Logger.Error(err.Error())
		return RegistrationResponse{Team: team, User: user}, err
	}

	if err := common.Validator.Struct(&team); err != nil {
		return RegistrationResponse{Team: team, User: user}, err
	}

	if err := common.Validator.Struct(&user); err != nil {
		return RegistrationResponse{Team: team, User: user}, err
	}

	if err := db.Save(&team); err != nil {
		return RegistrationResponse{Team: team, User: user}, err
	}

	user.TeamID = team.ID

	if err := db.Save(&user); err != nil {
		return RegistrationResponse{Team: team, User: user}, err
	}

	return RegistrationResponse{Team: team, User: user}, nil
}

// Verifys the user's email and password
func (request *LoginRequest) Reply(body io.ReadCloser) (TokenResponse, error) {
	response := TokenResponse{}

	if err := json.NewDecoder(body).Decode(request); err != nil {
		return response, err
	}

	params := db.QueryParams{"email": request.Email}
	user, err := db.Get[entities.User](params)

	if err != nil {
		return response, err
	}

	if user.VerifyPassword(request.Password) {
		payload := db.QueryParams{"user_id": user.ID, "team_id": user.TeamID}
		_, response.Token, _ = common.JWT.Encode(payload)
		return response, nil
	}

	err = errors.New("invalid email or password")
	return response, err
}
