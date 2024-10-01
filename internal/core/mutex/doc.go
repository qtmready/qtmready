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

// Package mutex provides a distributed, durable mutex implementation for Temporal workflows.
//
// This package offers a custom mutex solution that extends beyond Temporal's built-in mutex capabilities. While
// Temporal's native mutex is local to a specific workflow, this implementation provides global and durable locks that
// can persist across multiple workflows.
//
// Features:
//
//   - Global Locking: Allows locking resources across different workflows and activities.
//   - Durability: Locks persist even if the original locking workflow terminates unexpectedly.
//   - Timeout Handling: Supports automatic lock release after a specified timeout.
//   - Orphan Tracking: Keeps track of timed-out locks for potential recovery or cleanup.
//   - Cleanup Mechanism: Provides a way to clean up and shut down mutex workflows when no longer needed.
//   - Flexible Resource Identification: Supports a hierarchical resource ID system for precise locking.
//
// Global and durable locks are necessary in distributed systems for several reasons:
//
//   - Cross-Workflow Coordination: Ensures only one workflow can access a resource at a time.
//   - Long-Running Operations: Protects resources during extended operations, even if workflows crash.
//   - Consistency in Distributed State: Maintains consistency by serializing access to shared resources.
//   - Workflow Independence: Allows for flexible system design with runtime coordination.
//   - Fault Tolerance: Prevents conflicts during partial system failures and recovery.
//   - Complex Resource Hierarchies: Manages access to interrelated resources across workflows.
//
// The mutex provides four operations, all of which must be used during the lifecycle of usage:
//
//   - Prepare: Gets the reference for the lock. If not found, creates a new global reference.
//   - Acquire: Attempts to acquire the lock, blocking until successful or timeout occurs.
//   - Release: Releases the held lock, allowing other workflows to acquire it.
//   - Cleanup: Attempts to shut down the mutex workflow if it's no longer needed.
//
// Usage:
//
//	m := mutex.New(
//		ctx,
//		mutex.WithResourceID("io.quantm.stack.123.mutex"),
//		mutex.WithTimeout(30*time.Minute),
//	)
//	if err := m.Prepare(ctx); err != nil {
//		// handle error
//	}
//	if err := m.Acquire(ctx); err != nil {
//		// handle error
//	}
//	if err := m.Release(ctx); err != nil {
//		// handle error
//	}
//	if err := m.Cleanup(ctx); err != nil {
//		// handle error
//	}
//
// This mutex implementation relies on Temporal workflows and should be used within a Temporal workflow context.
package mutex
