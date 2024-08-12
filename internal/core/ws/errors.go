package ws

import (
	"fmt"
)

type (
	// ErrorCode represents the type of error that occurred in the hub operations.
	ErrorCode int
)

const (
	// ErrorTypeUnknown represents an unknown error.
	ErrorTypeUnknown ErrorCode = iota
	// ErrorTypeLocalSendFailed represents a failure to send a message locally.
	ErrorTypeLocalSendFailed
	// ErrorTypeQueryFailed represents a failure to query the user's queue.
	ErrorTypeQueryFailed
	// ErrorTypeUserNotRegistered represents an error when the user is not registered to any queue.
	ErrorTypeUserNotRegistered
	// ErrorTypeWorkflowExecutionFailed represents a failure to execute the Temporal workflow.
	ErrorTypeWorkflowExecutionFailed
	// ErrorTypeBroadcastFailed represents a failure to broadcast a message to a team.
	ErrorTypeBroadcastFailed
)

type (
	// HubError represents an error that occurred during hub operations.
	HubError struct {
		Code    ErrorCode
		Message string
		Err     error
	}
)

// Error returns the string representation of the error.
func (e *HubError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}

	return e.Message
}

// Unwrap returns the underlying error.
func (e *HubError) Unwrap() error {
	return e.Err
}

// NewHubError creates a new HubError with the given type, message, and underlying error.
func NewHubError(errType ErrorCode, message string, err error) *HubError {
	return &HubError{
		Code:    errType,
		Message: message,
		Err:     err,
	}
}
