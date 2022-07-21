package serializers

import (
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
	RegisterationResponse struct {
		User models.User `json:"user"`
	}
)

func (r *RegistrationRequest) Validate() error {

	if err := common.Validator.Struct(r); err != nil {
		return err
	}

	return nil
}

func (r *RegistrationRequest) Save() (models.User, error) {
	user := models.User{
		FirstName: r.FirstName,
		Email:     r.Email,
		Password:  r.Password,
		LastName:  r.LastName,
	}

	if err := r.Validate(); err != nil {
		return user, err
	}

	if err := user.Create(); err != nil {
		return user, err
	}

	return user, nil
}
