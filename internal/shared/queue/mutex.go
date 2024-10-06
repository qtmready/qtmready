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

package queue

import (
	"log/slog"
	"sync"

	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/shared"
)

var (
	mutex     queues.Queue // mutex queue instance.
	mutexonce sync.Once    // ensures initialization happens only once.
)

// Mutex returns a singleton instance of the mutex queue.
//
// Use this queue for all mutex-related operations.  This queue is intended for use with a worker.
//
// Example Usage:
//
//	worker := queues.CreateWorker(queue.Mutex())
//	// ... (RegisterWorkflows, RegisterActivities)
//	queue.Mutex().Start()
//
// NOTE: Ensure a worker is running to when using the package mutex.
func Mutex() queues.Queue {
	mutexonce.Do(func() {
		slog.Info("queues/mutex: initializing...")

		mutex = queues.New(
			queues.WithName("mutex"),                      // queue name.
			queues.WithClient(shared.Temporal().Client()), // temporal client for communication.
		)
	})

	return mutex
}
