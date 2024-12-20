package workflows

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/repos/states"
)

func Trunk(ctx workflow.Context, state *states.Trunk) error {
	state.Init(ctx)

	selector := workflow.NewSelector(ctx)

	for !state.RestartRecommended(ctx) {
		selector.Select(ctx)
	}

	return workflow.NewContinueAsNewError(ctx, Trunk, state)
}
