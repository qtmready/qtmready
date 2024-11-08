// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2022, 2024.
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

package mutex

import (
	"errors"
	"fmt"
)

var (
	ErrNilContext   = errors.New("contexts not initialized")
	ErrNoResourceID = errors.New("no resource ID provided")
)

type (
	MutexError struct {
		id   string // the id of the mutex.
		kind string // kind of error. can be "acquire lock", "release lock", or "start workflow".
	}
)

func (e *MutexError) Error() string {
	return fmt.Sprintf("%s: failed to %s.", e.id, e.kind)
}

// NewAcquireLockError creates a new acquire lock error.
func NewAcquireLockError(id string) error {
	return &MutexError{id, "acquire lock"}
}

// NewReleaseLockError creates a new release lock error.
func NewReleaseLockError(id string) error {
	return &MutexError{id, "release lock"}
}

// NewPrepareMutexError creates a new start workflow error.
func NewPrepareMutexError(id string) error {
	return &MutexError{id, "prepare mutex"}
}

func NewCleanupMutexError(id string) error {
	return &MutexError{id, "cleanup mutex"}
}
