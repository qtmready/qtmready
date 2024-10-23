package erratic

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/types/known/anypb"

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
func (e *QuantmError) ResetDetailsWith(info ErrorDetails) *QuantmError {
	e.Details = info
	return e
}

// AddDetail adds a key-value pair to the ErrorDetails field of the APIError.
func (e *QuantmError) AddDetail(key, value string) *QuantmError {
	if e.Details == nil {
		e.Details = make(ErrorDetails)
	}

	e.Details[key] = value

	return e
}

// SetVaidationErrors formats a validator.ValidationErrors object into the ErrorDetails field.
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

func (e *QuantmError) Unauthorized() *QuantmError {
	return e.AddDetail("reason", "are you logged in?")
}

func (e *QuantmError) Forbidden() *QuantmError {
	return e.AddDetail("reason", "you are not allowed to access this resource")
}

func (e *QuantmError) DataBaseError(err error) *QuantmError {
	return e.AddDetail("reason", "database error").AddDetail("internal", err.Error())
}

func (e *QuantmError) SetInternal(err error) *QuantmError {
	return e.AddDetail("internal", err.Error())
}

func (e *QuantmError) ToProto() error {
	code := codes.Code(codes.Unknown)

	switch e.Status {
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

	st := status.New(code, e.Message)

	details := make([]protoiface.MessageV1, 0) // Create an empty slice of protoiface.MessageV1

	info := &errdetails.ErrorInfo{
		Reason:   e.ID,
		Domain:   "quantm",
		Metadata: make(map[string]string),
	}

	for key, val := range e.Details {
		info.Metadata[key] = val
	}

	anydtl, err := anypb.New(info)
	if err != nil {
		fmt.Println("Error creating Any proto:", err)
	}
	// Convert to protoiface.MessageV1
	details = append(details, anydtl)

	// DebugInfo details
	if internal, ok := e.Details["internal"]; ok {
		dbg := &errdetails.DebugInfo{
			StackEntries: []string{internal},
			Detail:       "See stack entries for internal details.",
		}

		anyDetail, err := anypb.New(dbg)
		if err != nil {
			fmt.Println("Error creating Any proto:", err)
		}

		// Convert to protoiface.MessageV1
		details = append(details, anyDetail)
		delete(e.Details, "internal")
	}

	st, err = st.WithDetails(details...)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to attach error details: %v", err)
	}

	return st.Err()
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
