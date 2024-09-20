// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2022, 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.


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
