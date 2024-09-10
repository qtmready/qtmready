// Package inspect provides functions to inspect the contents of type Context. You logger should be configured to output
// debug messages to see the output.
package inspect

import (
	"reflect"
	"unsafe"

	"github.com/labstack/echo/v4"

	"go.breu.io/quantm/internal/shared"
)

// Context prints the contents of a context.
func Context(ctx any, inner bool) {
	contextValues := reflect.ValueOf(ctx).Elem()
	contextKeys := reflect.TypeOf(ctx).Elem()

	if !inner {
		shared.Logger().Debug("Fields For:", contextKeys.PkgPath(), contextKeys.Name())
	}

	if contextKeys.Kind() == reflect.Struct {
		for i := 0; i < contextValues.NumField(); i++ {
			reflectValue := contextValues.Field(i)
			reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).Elem()

			reflectField := contextKeys.Field(i)

			if reflectField.Name == "Context" {
				Context(reflectValue.Interface(), true)
			} else {
				shared.Logger().Debug("context", "name", reflectField.Name, "value", reflectValue.Interface())
			}
		}
	} else {
		shared.Logger().Debug("context is empty (int)\n")
	}
}

// EchoHeaders prints the headers of an echo context.
func EchoHeaders(ctx echo.Context) {
	shared.Logger().Debug("--- Headers ---")

	for k, v := range ctx.Request().Header {
		for _, vv := range v {
			shared.Logger().Debug(k, "val", vv)
		}
	}

	shared.Logger().Debug("--- End Headers ---")
}
