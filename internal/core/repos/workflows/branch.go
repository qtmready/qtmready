package workflows

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/repos/states"
)

func Branch(ctx workflow.Context, state *states.Branch) error {
	return nil
}
