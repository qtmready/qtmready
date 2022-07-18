package routes

// TODO: move this to seperate package or utils

import (
	"github.com/beego/beego/v2/core/validation"
)

// Looks up the db and return if the value already exists.
// Assumption being the obj provided has the same field name as the model field.
// TODO: complete this
func IsUnique(v *validation.Validation, obj interface{}, model interface{}) {}
