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

package shared

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrValidation          = errors.New("validation error")
	ErrInvalidRolloutState = errors.New("invalid rollout state")
)

type (
	ErrorMap map[string]string
)

func (e *APIError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// SetInternal sets the internal error.
func (e *APIError) SetInternal(err error) {
	e.Internal = err
}

func (e *APIError) WithInternal(err error) *APIError {
	return &APIError{
		Code:     e.Code,
		Message:  e.Message,
		Internal: err,
	}
}

func (e *APIError) format() *APIError {
	e.Message = e.Message.(error).Error()
	if e.Internal != nil && e.Errors == nil {
		errs := ErrorMap{}
		errs["internal"] = e.Internal.Error()
		e.Errors = &errs
	}

	return e
}

func (e *APIError) Unwrap() error {
	return e.Internal
}

func (e *ErrorMap) Get(key string) (string, bool) {
	val, ok := (*e)[key]
	return val, ok
}

// NewAPIError replaces echo.NewHTTPError.
func NewAPIError(code int, message error) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// EchoAPIErrorHandler adds syntax sugar to the default echo error handler.
func EchoAPIErrorHandler(err error, ctx echo.Context) {
	var apierr *APIError

	if ctx.Response().Committed {
		return
	}

	// We create an APIError from the error if it is not already one.
	apierr, ok := err.(*APIError)
	if !ok {
		apierr = NewAPIError(http.StatusInternalServerError, ErrInternalServerError).WithInternal(err)
	}

	// Now we check if the internal error is a validator.ValidationErrors.
	validerr, ok := apierr.Message.(validator.ValidationErrors)
	if ok {
		errs := ErrorMap{}
		for _, fe := range validerr {
			errs[fe.Field()] = TagMessage(fe.Tag())
		}
		// We set the error map to the APIError and set the error to ErrValidation.
		apierr = NewAPIError(apierr.Code, ErrValidation).WithInternal(validerr)
		apierr.Errors = &errs
	}

	// We set the status code and return the error.
	if ctx.Request().Method == http.MethodHead {
		err = ctx.NoContent(apierr.Code)
	} else {
		err = ctx.JSON(apierr.Code, apierr.format())
	}

	if err != nil {
		ctx.Logger().Error(err)
	}
}

func TagMessage(tag string) string {
	switch tag {
	case "required":
		return "required"
	case "email":
		return "invalid format"
	case "db_unique":
		return "already exists"
	default:
		return fmt.Sprintf("%s, validation error", tag)
	}
}
