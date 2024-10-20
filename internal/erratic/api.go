package erratic

import (
	"github.com/go-playground/validator/v10"

	"go.breu.io/quantm/internal/shared"
)

type (
	// APIError represents a generic API error.
	//
	// It includes the HTTP status code, a human-readable message, and additional information.
	APIError struct {
		ID      string    `json:"id"`      // Unique identifier for the error.
		Code    int       `json:"status"`  // HTTP status code of the error.
		Message string    `json:"message"` // Human-readable message describing the error.
		Info    ErrorInfo `json:"info"`    // Additional information about the error.
	}

	// ErrorInfo represents a map of key-value pairs providing additional information about an error.
	//
	// Example:
	//
	//     info := ErrorInfo{"field": "invalid value"}
	//     fmt.Println(info) // Output: map[string]string{"field": "invalid value"}
	ErrorInfo map[string]string

	// BadRequestError represents a 400 Bad Request error.
	BadRequestError APIError
	// UnauthorizedError represents a 401 Unauthorized error.
	UnauthorizedError APIError
	// NotFoundError represents a 404 Not Found error.
	NotFoundError APIError
)

// Error implements the error interface for APIError.
func (e *APIError) Error() string {
	return e.Message
}

// SetInfo sets the ErrorInformation field of the APIError.
//
// Example:
//
//	err := NewBadRequestError("field", "invalid value")
//	err.SetInfo(ErrorInformation{"field": "invalid value"})
//	fmt.Println(err.Information) // Output: map[string]string{"field": "invalid value"}
func (e *APIError) SetInfo(info ErrorInfo) *APIError {
	e.Info = info
	return e
}

// AddInfo adds a key-value pair to the ErrorInformation field of the APIError.
//
// Example:
//
//	err := NewBadRequestError()
//	err.AddInfo("field", "invalid value")
//	fmt.Println(err.Information) // Output: map[string]string{"field": "invalid value"}
func (e *APIError) AddInfo(key, value string) *APIError {
	if e.Info == nil {
		e.Info = make(ErrorInfo)
	}

	e.Info[key] = value

	return e
}

// FormatValidationError formats a validator.ValidationErrors object into the ErrorInformation field.
//
// Example:
//
//	err := validator.New().Struct(struct{}{})
//	apiErr := NewBadRequestError()
//	apiErr.FormatValidationError(err)
//	fmt.Println(apiErr.Information) // Output: (empty map)
func (e *APIError) FormatValidationError(err error) *APIError {
	valid, ok := err.(validator.ValidationErrors)
	if !ok {
		return e
	}

	for _, v := range valid {
		_ = e.AddInfo(v.Field(), v.Tag())
	}

	return e
}

func (e *APIError) NotLoggedIn() *APIError {
	return e.AddInfo("reason", "are you logged in?")
}

func (e *APIError) IllegalAccess() *APIError {
	return e.AddInfo("reason", "you are not allowed to access this resource")
}

func (e *APIError) DataBaseError(err error) *APIError {
	return e.AddInfo("reason", "database error").AddInfo("internal", err.Error())
}

func (e *APIError) SetInternal(err error) *APIError {
	return e.AddInfo("internal", err.Error())
}

// NewAPIError creates a new APIError instance.
//
// It takes the HTTP status code, a human-readable message, and optional key-value pairs for additional information.
//
// Example:
//
//	err := NewAPIError(400, "Bad Request", "field", "invalid value")
//	fmt.Println(err.Error()) // Output: Bad Request
//	fmt.Println(err.Information) // Output: map[string]string{"field": "invalid value"}
func NewAPIError(code int, message string, args ...string) *APIError {
	extra := false
	if len(args)%2 != 0 {
		extra = true
	}

	info := make(ErrorInfo)

	for i := 0; i < len(args); i += 2 {
		info[args[i]] = args[i+1]
	}

	if extra {
		info["unknown"] = args[len(args)-1]
	}

	return &APIError{
		ID:      shared.Idempotent(),
		Code:    code,
		Message: message,
		Info:    info,
	}
}

// NewBadRequestError creates a new 400 Bad Request error.
//
// Example:
//
//	err := NewBadRequestError("field", "invalid value")
//	fmt.Println(err.Error()) // Output: Bad Request
//	fmt.Println(err.Information) // Output: map[string]string{"field": "invalid value"}
func NewBadRequestError(args ...string) *APIError {
	return NewAPIError(400, "Bad Request", args...)
}

// NewUnauthorizedError creates a new 401 Unauthorized error.
//
// Example:
//
//	err := NewUnauthorizedError("user_id", "123")
//	fmt.Println(err.Error()) // Output: Unauthorized
//	fmt.Println(err.Information) // Output: map[string]string{"user_id": "123"}
func NewUnauthorizedError(args ...string) *APIError {
	return NewAPIError(401, "Unauthorized", args...)
}

// NewNotFoundError creates a new 404 Not Found error.
//
// Example:
//
//	err := NewNotFoundError("user_id", "123")
//	fmt.Println(err.Error()) // Output: Not Found
//	fmt.Println(err.Information) // Output: map[string]string{"user_id": "123"}
func NewNotFoundError(args ...string) *APIError {
	return NewAPIError(404, "Not Found", args...)
}

// NewInternalServerError creates a new 500 Internal Server Error.
//
// Example:
//
//	err := NewInternalServerError("reason", "database error")
//	fmt.Println(err.Error()) // Output: Internal Server Error
//	fmt.Println(err.Information) // Output: map[string]string{"reason": "database error"}
func NewInternalServerError(args ...string) *APIError {
	return NewAPIError(500, "Internal Server Error", args...)
}
