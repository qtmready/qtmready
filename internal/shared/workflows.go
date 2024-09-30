// Copyright © 2024, Breu, Inc. <info@breu.io>
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
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/shared/queue"
)

// WithDefaultActivityContext returns a workflow.Context with the default activity options applied.
// The default options include a StartToCloseTimeout of 60 seconds.
//
// Example:
//
//	ctx = shared.WithDefaultActivityContext(ctx)
func WithDefaultActivityContext(ctx workflow.Context) workflow.Context {
	return workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
	})
}

// WithIgnoredErrorsContext returns a workflow.Context with activity options configured with a
// StartToCloseTimeout of 60 seconds and a RetryPolicy that allows a single attempt and ignores
// specified error types.
//
// Example:
//
//	ignored := []string{"CustomErrorType"}
//	ctx = shared.WithIgnoredErrorsContext(ctx, ignored...)
func WithIgnoredErrorsContext(ctx workflow.Context, args ...string) workflow.Context {
	return workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts:        1,
			NonRetryableErrorTypes: args,
		},
	})
}

// WithLongRunningContext returns a workflow.Context with activity options configured for long-running activities.
// It sets the StartToCloseTimeout to 60 minutes and the HeartbeatTimeout to 30 seconds.
//
// Example:
//
//	ctx = shared.WithLongRunningContext(ctx)
func WithLongRunningContext(ctx workflow.Context) workflow.Context {
	return workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Minute,
		HeartbeatTimeout:    30 * time.Second,
	})
}

// WithCustomQueueContext returns a workflow.Context with activity options configured with a
// StartToCloseTimeout of 60 seconds and a dedicated task queue. This allows scheduling activities
// on a different queue than the one the workflow is running on.
//
// Example:
//
//	ctx = shared.WithCustomQueueContext(ctx, queues.MyTaskQueue)
func WithCustomQueueContext(ctx workflow.Context, q queue.Queue) workflow.Context {
	return workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 60 * time.Second,
		TaskQueue:           q.Name(),
	})
}
