// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

// Package inspect provides functions to inspect the contents of type Context. You logger should be configured to output
// debug messages to see the output.
package inspect

import (
	"reflect"
	"unsafe"

	"github.com/labstack/echo/v4"

	"go.breu.io/ctrlplane/internal/shared"
)

// Context prints the contents of a context.
func Context(ctx any, inner bool) {
	contextValues := reflect.ValueOf(ctx).Elem()
	contextKeys := reflect.TypeOf(ctx).Elem()

	if !inner {
		shared.Logger.Debug("Fields For:", contextKeys.PkgPath(), contextKeys.Name())
	}

	if contextKeys.Kind() == reflect.Struct {
		for i := 0; i < contextValues.NumField(); i++ {
			reflectValue := contextValues.Field(i)
			reflectValue = reflect.NewAt(reflectValue.Type(), unsafe.Pointer(reflectValue.UnsafeAddr())).Elem()

			reflectField := contextKeys.Field(i)

			if reflectField.Name == "Context" {
				Context(reflectValue.Interface(), true)
			} else {
				shared.Logger.Debug("context", "name", reflectField.Name, "value", reflectValue.Interface())
			}
		}
	} else {
		shared.Logger.Debug("context is empty (int)\n")
	}
}

// EchoHeaders prints the headers of an echo context.
func EchoHeaders(ctx echo.Context) {
	shared.Logger.Debug("--- Headers ---")

	for k, v := range ctx.Request().Header {
		for _, vv := range v {
			shared.Logger.Debug(k, "val", vv)
		}
	}

	shared.Logger.Debug("--- End Headers ---")
}
