// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2023, 2024.
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

type (
	// APIError defines the structure of an API error response.
	APIError struct {
		Code     int       `json:"code"`
		Errors   *ErrorMap `json:"errors,omitempty"`
		internal error     `json:"-"`
		Message  string    `json:"message"`
	}
	// ErrorMap is a map of error messages keyed by field names.
	ErrorMap map[string]string

	// BadRequest defines the structure of an API error response.
	BadRequest = APIError

	// InternalServerError defines the structure of an API error response.
	InternalServerError = APIError

	// NotFound defines the structure of an API error response.
	NotFound = APIError

	// Unauthorized defines the structure of an API error response.
	Unauthorized = APIError
)

// Get retrieves a value from an ErrorMap by key.
func (e *ErrorMap) Get(key string) (string, bool) {
	val, ok := (*e)[key]
	return val, ok
}

// format formats the APIError for JSON serialization.
func (e *APIError) format() *APIError {
	// Message is always an error, so no need to convert.
	if e.internal != nil && e.Errors == nil {
		e.Errors = &ErrorMap{"internal": e.internal.Error()}
	}

	return e
}

// Error returns a string representation of the APIError.
func (e *APIError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

// SetInternal sets the internal error for this APIError.
func (e *APIError) SetInternal(err error) {
	e.internal = err
}

// WithInternal returns a new APIError with the given internal error.
func (e *APIError) WithInternal(err error) *APIError {
	e.internal = err
	return e
}

// WithErrors returns a new APIError with the given validation errors.
func (e *APIError) WithErrors(errs *ErrorMap) *APIError {
	e.Errors = errs
	return e
}

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

// format_echo_error formats an Echo HTTPError for use in the APIError.
func format_echo_error(err *echo.HTTPError) error {
	return fmt.Errorf("%v", err.Message)
}

// to_api_error converts any error to an APIError.
func to_api_error(err error) *APIError {
	if apierr, ok := err.(*APIError); ok {
		return apierr // Return the existing APIError.
	}

	if httpErr, ok := err.(*echo.HTTPError); ok {
		// Create a new APIError based on the echo.HTTPError.
		return NewAPIError(httpErr.Code, format_echo_error(httpErr)).WithInternal(httpErr.Internal)
	}

	if validerr, ok := err.(validator.ValidationErrors); ok {
		// Handle validation errors directly.
		errs := ErrorMap{}
		for _, fe := range validerr {
			errs[fe.Field()] = tag_msg(fe.Tag())
		}

		return NewAPIError(http.StatusBadRequest, ErrValidation).WithErrors(&errs)
	}

	// Wrap any other error type in an APIError.
	return NewAPIError(http.StatusInternalServerError, err)
}

// EchoAPIErrorHandler is an Echo error handler that standardizes error responses.
//
// It converts any error to an APIError, handles validation errors, and formats the error response before sending
// it back to the client.
func EchoAPIErrorHandler(err error, ctx echo.Context) {
	if ctx.Response().Committed {
		return
	}

	// Convert to APIError.
	apierr := to_api_error(err)

	// Refactored response handling.
	if err_ := respond(ctx, apierr); err_ != nil {
		ctx.Logger().Error(err_)
	}
}

// respond handles the API response based on the request method.
func respond(ctx echo.Context, err *APIError) error {
	if ctx.Request().Method == http.MethodHead {
		return ctx.NoContent(err.Code)
	}

	return ctx.JSON(err.Code, err.format()) // Format the APIError before sending.
}

// NewAPIError creates a new APIError instance.
func NewAPIError(code int, message error) *APIError {
	return &APIError{
		// Message:  echo.NewHTTPError(code, message),
		internal: nil,
		Code:     code, // Set the Code
	}
}
