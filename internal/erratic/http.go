package erratic

import (
	"net/http"
)

// -- QuantmError wraps an error with an HTTP status code.

type (
	// BadRequestError represents a 400 Bad Request error.
	//
	// Deprecated: Not to be used directly. Only meant to generate OpenAPI.
	BadRequestError = QuantmError

	// UnauthorizedError represents a 401 Unauthorized error.
	//
	// Deprecated: Not to be used directly. Only meant to generate OpenAPI.
	UnauthorizedError = QuantmError

	// ForbiddenError represents a 403 Forbidden Error.
	// Deprecated: Not to be used directly. Only meant to generate OpenAPI.
	ForbiddenError = QuantmError

	// NotFoundError represents a 404 Not Found error.
	//
	// Deprecated: Not to be used directly. Only meant to generate OpenAPI.
	NotFoundError = QuantmError

	InternalServerError = QuantmError
)

// -- New*Error functions create a new QuantmError with the given status code and message.

// NewBadRequestError creates a new 400 Bad Request error.
func NewBadRequestError(args ...string) *QuantmError {
	return New(http.StatusBadRequest, "Bad Request.", args...)
}

// NewUnauthorizedError creates a new 401 Unauthorized error.
func NewUnauthorizedError(args ...string) *QuantmError {
	return New(http.StatusUnauthorized, "Not Authorized", args...)
}

func NewForbiddenError(args ...string) *QuantmError {
	return New(http.StatusForbidden, "Permission Denied.", args...)
}

// NewNotFoundError creates a new 404 Not Found error.
func NewNotFoundError(args ...string) *QuantmError {
	return New(http.StatusNotFound, "Resource Not Found.", args...)
}

// NewInternalServerError creates a new 500 Internal Server Error.
func NewInternalServerError(args ...string) *QuantmError {
	return New(http.StatusInternalServerError, "Internal Server Error.", args...)
}
