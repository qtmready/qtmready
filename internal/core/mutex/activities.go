package mutex

import (
	"context"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/shared"
)

// PrepareMutexActivity either starts a new mutex workflow for the requested resource or signals the running mutex to schedule a new lock.
// with the specified timeout.
func PrepareMutexActivity(ctx context.Context, payload *Info) (*workflow.Execution, error) {
	opts := shared.Temporal().Queue(shared.CoreQueue).WorkflowOptions(
		shared.WithWorkflowBlock("mutex"),
		shared.WithWorkflowBlockID(payload.ResourceID),
	)

	exe, err := shared.Temporal().
		Client().
		SignalWithStartWorkflow(ctx, opts.ID, WorkflowSignalPrepare.String(), payload, opts, Workflow, payload)

	if err != nil {
		return &workflow.Execution{}, err
	}

	return &workflow.Execution{ID: exe.GetID(), RunID: exe.GetRunID()}, nil
}
