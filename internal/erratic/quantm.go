package erratic

import (
	"log/slog"
	"runtime"

	"connectrpc.com/connect"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/types/known/anypb"

	"go.breu.io/quantm/internal/utils"
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

		internal error // internal error, if any.
	}
)

// Error implements the error interface for APIError.
func (e *QuantmError) Error() string {
	return e.Message
}

// Wrap sets the internal error field of QuantmError.
func (e *QuantmError) Wrap(err error) *QuantmError {
	e.internal = err

	return e
}

// Unwrap returns the internal error of the QuantmError.
func (e *QuantmError) Unwrap() error {
	return e.internal
}

// SetHintsWith sets the ErrorDetails field of the APIError.
func (e *QuantmError) SetHintsWith(hints Hints) *QuantmError {
	if e.Hints == nil {
		e.Hints = make(Hints)
	}

	for k, v := range hints {
		e.Hints[k] = v
	}

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

// WithReason adds a "reason" hint to the error.
func (e *QuantmError) WithReason(reason string) *QuantmError {
	if e.Hints == nil {
		e.Hints = make(Hints)
	}

	e.Hints["reason"] = reason

	return e
}

// WithResource adds a "resource" hint to the error.
func (e *QuantmError) WithResource(resource string) *QuantmError {
	if e.Hints == nil {
		e.Hints = make(Hints)
	}

	e.Hints["resource"] = resource

	return e
}

// WithStack adds a "stack" hint to the error containing the current stack trace.
func (e *QuantmError) WithStack(stack string) *QuantmError {
	if e.Hints == nil {
		e.Hints = make(Hints)
	}

	buf := make([]byte, 1024)
	buf = buf[:runtime.Stack(buf, false)]
	e.Hints["stack"] = string(buf)

	return e
}

// WithHint adds a hint to the error.
func (e *QuantmError) WithHint(key, value string) *QuantmError {
	if e.Hints == nil {
		e.Hints = make(Hints)
	}

	e.Hints[key] = value

	return e
}

// WithHints adds multiple hints to the error.
func (e *QuantmError) WithHints(hints Hints) *QuantmError {
	if e.Hints == nil {
		e.Hints = make(Hints)
	}

	for k, v := range hints {
		e.Hints[k] = v
	}

	return e
}

// ToProto converts the QuantmError to a gRPC error.
//
// It maps the HTTP status code to a corresponding gRPC error code, sets the error message,
// and attaches additional information as error details.
func (e *QuantmError) ToProto() *status.Status {
	code := CodeToProto(e.Code)
	grpc := status.New(code, e.Message)

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

	if stack, ok := e.Hints["stack"]; ok {
		trace := &errdetails.DebugInfo{
			StackEntries: []string{stack},
			Detail:       "See stack entries for internal details.",
		}

		anytrace, err := anypb.New(trace)
		if err != nil {
			slog.Warn("Error creating Any proto", "error", err.Error())
		}

		details = append(details, anytrace)
	}

	detailed, err := grpc.WithDetails(details...)
	if err != nil {
		return grpc
	}

	return detailed
}

// ToConnectError converts the QuantmError to a Connect error.
//
// It maps the HTTP status code to a corresponding Connect error code, sets the error message,
// and attaches additional information as error details.
func (e *QuantmError) ToConnectError() *connect.Error {
	code := CodeToConnect(e.Code)
	err := connect.NewError(code, e)

	for key, val := range e.Hints {
		err.Meta().Add(key, val)
	}

	return err
}

// New creates a new QuantmError instance.
//
// For developer convenience, especially when dealing with http or grpc handlers, it is recommended to use
// the following helper functions to create new errors:
//
//   - NewBadRequestError
//   - NewUnauthorizedError
//   - NewForbiddenError
//   - NewNotFoundError
//   - NewInternalServerError
//
// The function receives an error code, a human-readable message, and optional key-value pairs for additional
// information.
func New(module, code int, message string, fields ...string) *QuantmError {
	return &QuantmError{
		ID:      utils.Idempotent(),
		Code:    compose(module, code),
		Message: message,
		Hints:   NewHints(fields...),
	}
}
