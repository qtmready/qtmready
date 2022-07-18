package routes

import (
	"github.com/beego/beego/v2/core/validation"

	"go.breu.io/ctrlplane/internal/common"
	"go.breu.io/ctrlplane/internal/db/models"
)

type regRequest struct {
	Email            string `json:"email" valid:"Required;Email"`
	Password         string `json:"password valid:Required"`
	ConfirmPassword  string `json:"confirm_password" valid:"Required"`
	OrganizationName string `json:"organization_name" valid:"Required"`
}

// Validate validates the request.
// TODO: leverage the validator package. See https://github.com/beego/beego/tree/develop/core/validation
func (r *regRequest) validate() error {
	if r.Password != r.ConfirmPassword {
		return ErrorPasswordMismatch
	}

	validator := validation.Validation{}

	valid, err := validator.Valid(r)

	if err != nil {
		return err
	}

	if !valid {
		for _, err := range validator.Errors {
			common.Logger.Info(err.Key)
			common.Logger.Info(err.Message)
		}
	}

	user := models.User{Email: r.Email}
	if err := user.Get(user); err != nil {
		return nil
	}

	return ErrorEmailAlreadyExists
}

func (r *regRequest) save() error {
	if err := r.validate(); err != nil {
		return err
	}
	return nil
}
