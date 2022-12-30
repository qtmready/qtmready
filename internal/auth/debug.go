package auth

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/labstack/echo/v4"

	"go.breu.io/ctrlplane/internal/shared"
)

func printContext(ctx interface{}, inner bool) {
	contextValues := reflect.ValueOf(ctx).Elem()
	contextKeys := reflect.TypeOf(ctx).Elem()

	if !inner {
		fmt.Printf("\nFields for %s.%s\n", contextKeys.PkgPath(), contextKeys.Name())
	}

	if contextKeys.Kind() == reflect.Struct {
		for i := 0; i < contextValues.NumField(); i++ {
			reflectValue := contextValues.Field(i)
			reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).Elem()

			reflectField := contextKeys.Field(i)

			if reflectField.Name == "Context" {
				printContext(reflectValue.Interface(), true)
			} else {
				// fmt.Printf("field name: %+v\n", reflectField.Name)
				shared.Logger.Debug("context", "name", reflectField.Name, "value", reflectValue.Interface())
				// fmt.Printf("value: %+v\n", reflectValue.Interface())
			}
		}
	} else {
		fmt.Printf("context is empty (int)\n")
	}
}

func printHeaders(ctx echo.Context) {
	shared.Logger.Debug("--- Headers ---")

	for k, v := range ctx.Request().Header {
		for _, vv := range v {
			shared.Logger.Debug(k, "val", vv)
		}
	}

	shared.Logger.Debug("--- End Headers ---")
}
