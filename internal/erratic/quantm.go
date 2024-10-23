package erratic

import (
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.breu.io/quantm/internal/shared"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/types/known/anypb"
)

type (
	// QuantmError is the standard error type used within the application.
	//
	// It encapsulates a unique identifier, an error code (initially HTTP status codes), a human-readable
	// message, and additional error details. The initial implementation uses HTTP status codes for
	// convenience, scheme may be adopted in the future. This is subject to how the application evolves.
	QuantmError struct {
		ID      string `json:"id"`      // Unique identifier for the error.
		Code    int    `json:"code"`    // HTTP status code of the error.
		Message string `json:"message"` // Human-readable message describing the error.
		Hints   Hints  `json:"hints"`   // Additional information about the error.
	}
)

// Error implements the error interface for APIError.
func (e *QuantmError) Error() string {
	return e.Message
}

// SetHintsWith sets the ErrorDetails field of the APIError.
func (e *QuantmError) SetHintsWith(hints Hints) *QuantmError {
	e.Hints = hints
	return e
}

// AddHint adds a key-value pair to the ErrorDetails field of the APIError.
func (e *QuantmError) AddHint(key, value string) *QuantmError {
	if e.Hints == nil {
		e.Hints = make(Hints)
	}

	e.Hints[key] = value

	return e
}

// SetVaidationErrors formats a validator.ValidationErrors object into the ErrorDetails field.
func (e *QuantmError) SetVaidationErrors(err error) *QuantmError {
	valid, ok := err.(validator.ValidationErrors)
	if !ok {
		return e
	}

	for _, v := range valid {
		_ = e.AddHint(v.Field(), v.Tag())
	}

	return e
}

// DataBaseError sets the ErrorDetails field of the APIError with information related to a database error.
func (e *QuantmError) DataBaseError(err error) *QuantmError {
	return e.AddHint("reason", "database error").AddHint("internal", err.Error())
}

// SetInternal sets the ErrorDetails field of the APIError with an internal error message.
func (e *QuantmError) SetInternal(err error) *QuantmError {
	return e.AddHint("internal", err.Error())
}

// ToProto converts the QuantmError to a gRPC error.
//
// It maps the HTTP status code to a corresponding gRPC error code, sets the error message,
// and attaches additional information as error details.
func (e *QuantmError) ToProto() *status.Status {
	code := codes.Unknown

	// Map HTTP status code to gRPC error code and create a new grpc status.

	switch e.Code {
	case http.StatusBadRequest:
		code = codes.InvalidArgument
	case http.StatusUnauthorized:
		code = codes.Unauthenticated
	case http.StatusForbidden:
		code = codes.PermissionDenied
	case http.StatusNotFound:
		code = codes.NotFound
	case http.StatusInternalServerError:
		code = codes.Internal
	}

	sts := status.New(code, e.Message)

	// Creating error details from the hints. See
	//
	// - https://grpc.io/docs/guides/error/#richer-error-model
	// - https://cloud.google.com/apis/design/errors#error_model

	details := make([]protoiface.MessageV1, 0)

	info := &errdetails.ErrorInfo{
		Reason:   e.Message,
		Domain:   "quantm",
		Metadata: make(map[string]string),
	}

	for key, val := range e.Hints {
		info.Metadata[key] = val
	}

	anyinfo, err := anypb.New(info)
	if err != nil {
		slog.Warn("Error creating Any proto", "error", err.Error())
	}

	details = append(details, anyinfo)

	if internal, ok := e.Hints["internal"]; ok {
		trace := &errdetails.DebugInfo{
			StackEntries: []string{internal},
			Detail:       "See stack entries for internal details.",
		}

		anytrace, err := anypb.New(trace)
		if err != nil {
			slog.Warn("Error creating Any proto", "error", err.Error())
		}

		details = append(details, anytrace)

		delete(e.Hints, "internal")
	}

	// Finally, attach the error details to the status.

	detailed, err := sts.WithDetails(details...)
	if err != nil {
		return sts
	}

	return detailed
}

// New creates a new QuantmError instance.
//
// This function should never be called directly. Use the following functions instead:
//
//   - NewBadRequestError
//   - NewUnauthorizedError
//   - NewForbiddenError
//   - NewNotFoundError
//   - NewInternalServerError
//
// The function receives an error code, a human-readable message, and optional key-value pairs for additional
// information.  Currently, HTTP status codes are used, but this may be revised in the future.
func New(code int, message string, args ...string) *QuantmError {
	return &QuantmError{
		ID:      shared.Idempotent(),
		Code:    code,
		Message: message,
		Hints:   NewErrorDetails(args...),
	}
}
