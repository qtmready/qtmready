package activities

import (
	"context"

	"go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/durable"
)

type (
	Repo struct{}
)

const (
	WorkflowBranch = "Branch" // WorkflowBranch is string representation of workflows.Branch
	WorkflowTrunk  = "Trunk"  // WorkflowTrunk is string representation of workflows.Trunk
)

func (a *Repo) ForwardToBranch(ctx context.Context, payload *defs.SignalBranchPayload, event, state any) error {
	_, err := durable.OnCore().SignalWithStartWorkflow(
		ctx,
		defs.BranchWorkflowOptions(payload.Repo, payload.Branch),
		payload.Signal,
		event,
		WorkflowBranch,
		state,
	)

	return err
}

func (a *Repo) ForwardToTrunk(ctx context.Context, payload *defs.SignalTrunkPayload, event, state any) error {
	return nil
}

func (a *Repo) ForwardToQueue(ctx context.Context, payload *defs.SignalQueuePayload, event, state any) error {
	return nil
}
