package validations

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

func Unique(fl validator.FieldLevel) bool {
	params := map[string]interface{}{fl.FieldName(): fl.Field().Interface()}
	args := []reflect.Value{reflect.ValueOf(params)}
	if err := fl.Parent().Addr().MethodByName("Get").Call(args); err != nil {
		return false
	}
	return true
}
