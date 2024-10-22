package erratic

import (
	"github.com/go-playground/validator/v10"

	"go.breu.io/quantm/internal/shared"
)

type (
	// QuantmError is the standard error rising from application.
	//
	// It includes the HTTP status code, a human-readable message, and additional information.
	QuantmError struct {
		ID      string       `json:"id"`      // Unique identifier for the error.
		Status  int          `json:"status"`  // HTTP status code of the error.
		Message string       `json:"message"` // Human-readable message describing the error.
		Details ErrorDetails `json:"details"` // Additional information about the error.
	}
)

// Error implements the error interface for APIError.
func (e *QuantmError) Error() string {
	return e.Message
}

// ResetDetailsWith sets the ErrorDetails field of the APIError.
//
// Example:
//
//	err := NewBadRequestError("field", "invalid value")
//	err.ResetDetailsWith(ErrorDetails{"field": "invalid value"})
//	fmt.Println(err.Information) // Output: map[string]string{"field": "invalid value"}
func (e *QuantmError) ResetDetailsWith(info ErrorDetails) *QuantmError {
	e.Details = info
	return e
}

// AddDetail adds a key-value pair to the ErrorDetails field of the APIError.
//
// Example:
//
//	err := NewBadRequestError()
//	err.AddDetail("field", "invalid value")
//	fmt.Println(err.Information) // Output: map[string]string{"field": "invalid value"}
func (e *QuantmError) AddDetail(key, value string) *QuantmError {
	if e.Details == nil {
		e.Details = make(ErrorDetails)
	}

	e.Details[key] = value

	return e
}

// SetVaidationErrors formats a validator.ValidationErrors object into the ErrorDetails field.
//
// Example:
//
//	err := validator.New().Struct(struct{}{})
//	apiErr := NewBadRequestError()
//	apiErr.SetVaidationErrors(err)
//	fmt.Println(apiErr.Information) // Output: (empty map)
func (e *QuantmError) SetVaidationErrors(err error) *QuantmError {
	valid, ok := err.(validator.ValidationErrors)
	if !ok {
		return e
	}

	for _, v := range valid {
		_ = e.AddDetail(v.Field(), v.Tag())
	}

	return e
}

func (e *QuantmError) NotLoggedIn() *QuantmError {
	return e.AddDetail("reason", "are you logged in?")
}

func (e *QuantmError) IllegalAccess() *QuantmError {
	return e.AddDetail("reason", "you are not allowed to access this resource")
}

func (e *QuantmError) DataBaseError(err error) *QuantmError {
	return e.AddDetail("reason", "database error").AddDetail("internal", err.Error())
}

func (e *QuantmError) SetInternal(err error) *QuantmError {
	return e.AddDetail("internal", err.Error())
}

// New creates a new QuantmError instance.
//
// It takes the HTTP status code, a human-readable message, and optional key-value pairs for additional information.
//
// Example:
//
//	err := New(400, "Bad Request", "field", "invalid value")
//	fmt.Println(err.New()) // Output: Bad Request
//	fmt.Println(err.Information) // Output: map[string]string{"field": "invalid value"}
func New(code int, message string, args ...string) *QuantmError {
	return &QuantmError{
		ID:      shared.Idempotent(),
		Status:  code,
		Message: message,
		Details: NewErrorDetails(args...),
	}
}
