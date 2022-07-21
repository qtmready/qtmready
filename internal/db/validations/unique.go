package validations

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

// Unique validates that the value of the field is unique in the database.
func Unique(fl validator.FieldLevel) bool {
	params := map[string]interface{}{fl.FieldName(): fl.Field().Interface()}
	args := []reflect.Value{reflect.ValueOf(params)}
	// if there is any kind of error, the validation fails.
	result := fl.Parent().Addr().MethodByName("Get").Call(args)
	err := result[0].Interface()
	return err != nil
}
