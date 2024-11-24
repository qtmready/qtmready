package activities

import (
	"context"

	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/db/entities"
)

type (
	Repo struct{}

	SignalBranchPayload struct{}

	SignalTrunkPayload struct{}
)

const (
	WorkflowBranch = "Branch" // WorkflowBranch is string representation of workflows.Branch
	WorkflowTrunk  = "Trunk"  // WorkflowTrunk is string representation of workflows.Trunk
)

func (a *Repo) SignalBranch(ctx context.Context, signal queues.Signal, repo *entities.Repo, branch string, event any) error {
	return nil
}

func (a *Repo) SignalTrunk(ctx context.Context, signal queues.Signal, repo *entities.Repo, event any) error {
	return nil
}
