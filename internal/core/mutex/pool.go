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

package mutex

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type (
	// Pool represents a collection of Resource IDs with associated timeouts
	// for distributed locks. It can be used to track waiting, orphaned,
	// or processed resources in a distributed system.
	Pool struct {
		Data  map[string]time.Duration `json:"data"`
		mutex workflow.Mutex
	}
)

// add inserts or updates a Resource ID in the pool with the given timeout.
// It ensures thread-safe access to the Data map using a mutex.
func (p *Pool) add(ctx workflow.Context, resourceID string, timeout time.Duration) {
	_ = p.mutex.Lock(ctx)
	defer p.mutex.Unlock()
	p.Data[resourceID] = timeout
}

// remove deletes a Resource ID from the pool.
// It ensures thread-safe access to the Data map using a mutex.
func (p *Pool) remove(ctx workflow.Context, resourceID string) {
	_ = p.mutex.Lock(ctx)
	defer p.mutex.Unlock()
	delete(p.Data, resourceID)
}

// get retrieves the timeout for a Resource ID.
// It returns the timeout and a boolean indicating if the Resource ID was found.
func (p *Pool) get(resourceID string) (time.Duration, bool) {
	timeout, ok := p.Data[resourceID]
	return timeout, ok
}

// size returns the number of Resource IDs currently in the pool.
func (p *Pool) size() int {
	return len(p.Data)
}

// restore initializes the mutex and ensures the Data map is not nil.
// It should be called after deserializing a Pool instance or when
// creating a Pool instance from existing data.
func (p *Pool) restore(ctx workflow.Context) {
	p.mutex = workflow.NewMutex(ctx)
	if p.Data == nil {
		p.Data = make(map[string]time.Duration)
	}
}

// NewPool creates and returns a new Pool instance with an initialized
// Data map and mutex for managing Resource IDs and their timeouts.
func NewPool(ctx workflow.Context) *Pool {
	return &Pool{
		Data:  make(map[string]time.Duration),
		mutex: workflow.NewMutex(ctx),
	}
}
