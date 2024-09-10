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

package mutex

import (
	"context"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/shared"
)

// PrepareMutexActivity either starts a new mutex workflow for the requested resource or signals the running mutex to schedule a new lock
// with the specified timeout.
//
// Parameters:
//   - ctx: The context for the activity execution.
//   - payload: A pointer to a Handler struct containing the resource ID and timeout for the mutex.
//
// Returns:
//   - *workflow.Execution: A pointer to a workflow.Execution struct containing the ID and RunID of the started or signaled workflow.
//   - error: An error if the operation fails, or nil if successful.
//
// The function performs the following steps:
// 1. Creates workflow options using the shared.Temporal() helper, setting the queue and workflow block details.
// 2. Creates a new MutexState instance with the necessary initial values.
// 3. Calls SignalWithStartWorkflow on the Temporal client to either start a new workflow or signal an existing one.
// 4. If an error occurs during the SignalWithStartWorkflow call, it returns an empty workflow.Execution and the error.
// 5. On success, it returns a workflow.Execution with the ID and RunID from the started or signaled workflow.
//
// This activity is typically used as part of the mutex preparation process in a distributed system,
// ensuring that mutex operations are properly coordinated across different workflows.
func PrepareMutexActivity(ctx context.Context, payload *Handler) (*workflow.Execution, error) {
	opts := shared.Temporal().Queue(shared.CoreQueue).WorkflowOptions(
		shared.WithWorkflowBlock("mutex"),
		shared.WithWorkflowBlockID(payload.ResourceID),
	)

	state := &MutexState{
		Status:  MutexStatusAcquiring,
		Handler: payload,
		Timeout: payload.Timeout,
		Persist: true,
	}

	exe, err := shared.Temporal().
		Client().
		SignalWithStartWorkflow(context.Background(), opts.ID, WorkflowSignalPrepare.String(), payload, opts, MutexWorkflow, state)

	if err != nil {
		return &workflow.Execution{}, err
	}

	return &workflow.Execution{ID: exe.GetID(), RunID: exe.GetRunID()}, nil
}
