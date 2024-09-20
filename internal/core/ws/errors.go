// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
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

package ws

import (
	"fmt"
)

type (
	// ErrorCode represents the type of error that occurred in the hub operations.
	ErrorCode int
)

const (
	// ErrorCodeUnknown represents an unknown error.
	ErrorCodeUnknown ErrorCode = iota
	// ErrorCodeLocalSendFailed represents a failure to send a message locally.
	ErrorCodeLocalSendFailed
	// ErrorCodeQueryFailed represents a failure to query the user's queue.
	ErrorCodeQueryFailed
	// ErrorCodeUserNotRegistered represents an error when the user is not registered to any queue.
	ErrorCodeUserNotRegistered
	// ErrorCodeWorkflowExecutionFailed represents a failure to execute the Temporal workflow.
	ErrorCodeWorkflowExecutionFailed
	// ErrorCodeBroadcastFailed represents a failure to broadcast a message to a team.
	ErrorCodeBroadcastFailed
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
