package workflows

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/core/repos/states"
)

func Trunk(ctx workflow.Context, state *states.Trunk) error {
	state.Init(ctx)

	selector := workflow.NewSelector(ctx)

	mq := workflow.GetSignalChannel(ctx, defs.SignalMergeQueue.String())
	selector.AddReceive(mq, state.OnMergeQueue(ctx))

	workflow.Go(ctx, state.StartQueueProcessor)

	for state.Continue() {
		selector.Select(ctx)
	}

	return workflow.NewContinueAsNewError(ctx, Trunk, state)
}
