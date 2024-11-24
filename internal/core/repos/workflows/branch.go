package workflows

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/core/repos/states"
)

func Branch(ctx workflow.Context, state *states.Branch) error {
	state.Init(ctx)

	selector := workflow.NewSelector(ctx)

	push := workflow.GetSignalChannel(ctx, defs.SignalPush.String())
	selector.AddReceive(push, state.OnPush(ctx))

	for !state.RestartRecommended(ctx) {
		selector.Select(ctx)
	}

	return nil
}
