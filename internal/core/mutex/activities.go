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

package mutex

import (
	"context"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/shared/queues"
)

// PrepareMutexActivity prepares a mutex for a given resource.
//
// It either stars a new mutex workflow or signals an existing one to schedule a new lock.  It prepares the mutex for a
// given resource by creating a MutexState with initial values and using SignalWithStartWorkflow to manage the
// workflow.  Errors indicate issues during workflow signal or start. Success returns a workflow.Execution with the
// workflow's ID and RunID.
func PrepareMutexActivity(ctx context.Context, payload *Handler) (*workflow.Execution, error) {
	state := &MutexState{
		Status:  MutexStatusAcquiring,
		Handler: payload,
		Timeout: payload.Timeout,
		Persist: true,
	}

	exe, err := queues.Mutex().SignalWithStartWorkflow(
		ctx,
		MutexWorkflowOptions(payload.ResourceID),
		WorkflowSignalPrepare,
		payload,
		MutexWorkflow,
		state,
	)

	return &workflow.Execution{ID: exe.GetID(), RunID: exe.GetRunID()}, err
}
