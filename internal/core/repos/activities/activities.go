package activities

import (
	"context"

	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/durable"
)

type (
	SignalBranchPayload struct {
		Signal queues.Signal  `json:"signal"`
		Repo   *entities.Repo `json:"repo"`
		Branch string         `json:"branch"`
	}

	SignalTrunkPayload struct{}

	SignalQueuePayload struct{}

	Repo struct{}
)

const (
	WorkflowBranch = "Branch" // WorkflowBranch is string representation of workflows.Branch
	WorkflowTrunk  = "Trunk"  // WorkflowTrunk is string representation of workflows.Trunk
)

func (a *Repo) SignalBranch(ctx context.Context, payload *SignalBranchPayload, event, state any) error {
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

func (a *Repo) SignalTrunk(ctx context.Context, payload *SignalTrunkPayload, event, state any) error {
	return nil
}

func (a *Repo) SignalQueue(ctx context.Context, payload *SignalQueuePayload, event, state any) error {
	return nil
}
