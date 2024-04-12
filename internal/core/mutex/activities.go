package mutex

import (
	"context"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/shared"
)

// PrepareMutex starts the mutex workflow, signals an existing workflow, for the given provider and resource,
// with the specified timeout.
func PrepareMutex(ctx context.Context, payload *Info) (*workflow.Execution, error) {
	opts := shared.Temporal().Queue(shared.CoreQueue).WorkflowOptions(
		shared.WithWorkflowBlock(payload.ID),
		shared.WithWorkflowBlockID("mutex"),
	)

	exe, err := shared.Temporal().
		Client().
		SignalWithStartWorkflow(ctx, opts.ID, WorkflowSignalPrepare.String(), payload, opts, Workflow, payload)

	if err != nil {
		return &workflow.Execution{}, err
	}

	return &workflow.Execution{ID: exe.GetID(), RunID: exe.GetRunID()}, nil
}
