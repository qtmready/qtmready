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
// 2. Calls SignalWithStartWorkflow on the Temporal client to either start a new workflow or signal an existing one.
// 3. If an error occurs during the SignalWithStartWorkflow call, it returns an empty workflow.Execution and the error.
// 4. On success, it returns a workflow.Execution with the ID and RunID from the started or signaled workflow.
//
// This activity is typically used as part of the mutex preparation process in a distributed system,
// ensuring that mutex operations are properly coordinated across different workflows.
func PrepareMutexActivity(ctx context.Context, payload *Handler) (*workflow.Execution, error) {
	opts := shared.Temporal().Queue(shared.CoreQueue).WorkflowOptions(
		shared.WithWorkflowBlock("mutex"),
		shared.WithWorkflowBlockID(payload.ResourceID),
	)

	exe, err := shared.Temporal().
		Client().
		SignalWithStartWorkflow(context.Background(), opts.ID, WorkflowSignalPrepare.String(), payload, opts, MutexWorkflow, payload)

	if err != nil {
		return &workflow.Execution{}, err
	}

	return &workflow.Execution{ID: exe.GetID(), RunID: exe.GetRunID()}, nil
}
