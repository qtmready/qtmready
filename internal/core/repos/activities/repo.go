package activities

import (
	"context"
	"log/slog"

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
	id := defs.BranchWorkflowOptions(payload.Repo, payload.Branch)
	run, err := durable.
		OnCore().
		SignalWithStartWorkflow(ctx, id, payload.Signal, event, WorkflowBranch, state)

	if err != nil {
		slog.Warn("fwd_to_branch: unable to signal", "id", id.IDSuffix(), "error", err.Error())
	} else {
		slog.Info("fwd_to_branch: signaled", "id", id.IDSuffix(), "run_id", run.GetRunID())
	}

	return err
}

func (a *Repo) ForwardToTrunk(ctx context.Context, payload *defs.SignalTrunkPayload, event, state any) error {
	_, err := durable.
		OnCore().
		SignalWithStartWorkflow(ctx, defs.TrunkWorkflowOptions(payload.Repo), payload.Signal, event, WorkflowTrunk, state)

	return err
}

func (a *Repo) ForwardToQueue(ctx context.Context, payload *defs.SignalQueuePayload, event, state any) error {
	return nil
}
