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

// Package shared provides a set of utilities for handling errors in an API. It
// defines a standard error structure (APIError) and provides helper functions
// for converting and formatting errors.
package shared

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Error types.
var (
	// ErrInternalServerError represents an internal server error.
	ErrInternalServerError = errors.New("internal server error")
	// ErrValidation represents a validation error.
	ErrValidation = errors.New("validation error")
	// ErrInvalidRolloutState represents an invalid rollout state error.
	ErrInvalidRolloutState = errors.New("invalid rollout state")
)

// Type definitions.
type (
	// ErrorMap is a map of error messages keyed by field names.
	ErrorMap map[string]string
)

// Get retrieves a value from an ErrorMap by key.
func (e *ErrorMap) Get(key string) (string, bool) {
	val, ok := (*e)[key]
	return val, ok
}

// Private receiver methods

// format formats the APIError for JSON serialization.
func (e *APIError) format() *APIError {
	// Message is always an error, so no need to convert
	if e.Internal != nil && e.Errors == nil {
		e.Errors = &ErrorMap{"internal": e.Internal.Error()}
	}

	return e
}

// Public receiver methods

// Error returns a string representation of the APIError.
func (e *APIError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message.Error())
}

// SetInternal sets the internal error for this APIError.
func (e *APIError) SetInternal(err error) {
	e.Internal = err
}

// WithInternal returns a new APIError with the given internal error.
func (e *APIError) WithInternal(err error) *APIError {
	e.Internal = err
	return e
}

// WithErrors returns a new APIError with the given validation errors.
func (e *APIError) WithErrors(errs *ErrorMap) *APIError {
	e.Errors = errs
	return e
}

// Private methods

// tag_msg returns a user-friendly error message for a given validation tag.
func tag_msg(tag string) string {
	switch tag {
	case "required":
		return "This field is required."
	case "email":
		return "Please enter a valid email address."
	case "db_unique":
		return "This value already exists."
	default:
		return fmt.Sprintf("%s, validation error", tag)
	}
}

// handle_validation handles validation errors from validator.
func handle_validation(apiErr *APIError) *APIError {
	if validerr, ok := apiErr.Message.(validator.ValidationErrors); ok {
		errs := ErrorMap{}
		for _, fe := range validerr {
			errs[fe.Field()] = tag_msg(fe.Tag())
		}

		// Return a new APIError with status 400 and validation message
		return NewAPIError(http.StatusBadRequest, ErrValidation).WithInternal(validerr).WithErrors(&errs)
	}

	return apiErr // Return the original APIError if no validation errors
}

// to_api_error converts any error to an APIError.
func to_api_error(err error) *APIError {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr
	}

	// We should never reach this point because all errors are wrapped in APIError
	panic("Unexpected error type: " + err.Error())
}

// Public methods

// EchoAPIErrorHandler is an Echo error handler that standardizes error responses.
// It converts any error to an APIError, handles validation errors, and formats
// the error response before sending it back to the client.
func EchoAPIErrorHandler(err error, ctx echo.Context) {
	if ctx.Response().Committed {
		return
	}

	// Convert to APIError
	apierr := to_api_error(err)

	// Handle validation errors
	apierr = handle_validation(apierr)

	// Refactored response handling
	if err_ := respond(ctx, apierr); err_ != nil {
		ctx.Logger().Error(err_)
	}
}

// respond handles the API response based on the request method.
func respond(ctx echo.Context, err *APIError) error {
	if ctx.Request().Method == http.MethodHead {
		return ctx.NoContent(err.Code)
	}

	return ctx.JSON(err.Code, err)
}

// NewAPIError creates a new APIError instance.
func NewAPIError(code int, message error) *APIError {
	return &APIError{
		Message:  echo.NewHTTPError(code, message),
		Internal: nil,
		Code:     code, // Set the Code
	}
}
