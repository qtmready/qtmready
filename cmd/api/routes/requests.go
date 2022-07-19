package routes

import (
	"go.breu.io/ctrlplane/internal/common"
	"go.breu.io/ctrlplane/internal/db/models"
)

type RegRequest struct {
	Email            string `json:"email" validate:"required,email"`
	Password         string `json:"password" validate:"required"`
	FirstName        string `json:"first_name" validate:"required"`
	LastName         string `json:"last_name" validate:"required"`
	ConfirmPassword  string `json:"confirm_password" validate:"required,eqfield=Password"`
	OrganizationName string `json:"organization_name" validate:"required"`
}

func (r *RegRequest) validate() error {

	if err := common.Validator.Struct(r); err != nil {
		return err
	}

	return nil
}

func (r *RegRequest) save() error {
	if err := r.validate(); err != nil {
		return err
	}

	user := models.User{
		FirstName: r.FirstName,
		Email:     r.Email,
		Password:  r.Password,
		LastName:  r.LastName,
	}
	if err := user.Create(user); err != nil {
		return err
	}
	return nil
}
