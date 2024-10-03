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
	"encoding/json"

	"github.com/go-playground/validator/v10"

	"go.breu.io/quantm/internal/shared/queue"
)

// workflow related shared types and contants.
type (
	// WorkflowSignal is a type alias to define the name of the workflow signal.
	WorkflowSignal string

	WorkflowOption = queue.WorkflowOptions
)

var (
	WithWorkflowParent    = queue.WithWorkflowParent    // Sets the parent workflow ID.
	WithWorkflowBlock     = queue.WithWorkflowBlock     // Sets the block name for the workflow ID.
	WithWorkflowBlockID   = queue.WithWorkflowBlockID   // Sets the block value for the workflow ID.
	WithWorkflowElement   = queue.WithWorkflowElement   // Sets the element name for the workflow ID.
	WithWorkflowElementID = queue.WithWorkflowElementID // Sets the element value for the workflow ID.
	WithWorkflowMod       = queue.WithWorkflowMod       // Sets the modifier for the workflow ID.
	WithWorkflowModID     = queue.WithWorkflowModID     // Sets the modifier value for the workflow ID.
	WithWorkflowProp      = queue.WithWorkflowProp      // Sets the property for the workflow ID.
	NewWorkflowOptions    = queue.NewWorkflowOptions    // Creates a new workflow ID. (see queue.NewWorkflowOptions)
)

// queue definitions.
const (
	CoreQueue      queue.Name = "core"      // core queue
	ProvidersQueue queue.Name = "providers" // messaging related to providers
	MutexQueue     queue.Name = "mutex"     // mutex workflow queue
	WebSocketQueue queue.Name = "websocket" // websocket workflow queue
)

/*
 * Methods for WorkflowSignal.
 */
func (w WorkflowSignal) String() string { return string(w) }

// MarshalJSON implements the json.Marshaler interface for WorkflowSignal.
func (w WorkflowSignal) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(w))
}

// UnmarshalJSON implements the json.Unmarshaler interface for WorkflowSignal.
func (w *WorkflowSignal) UnmarshalJSON(data []byte) error {
	var signal string
	if err := json.Unmarshal(data, &signal); err != nil {
		return err
	}

	*w = WorkflowSignal(signal)

	return nil
}

type (
	// EchoValidator is a wrapper for the instantiated validator.
	EchoValidator struct {
		Validator *validator.Validate
	}
)

func (ev *EchoValidator) Validate(i any) error {
	return ev.Validator.Struct(i)
}
