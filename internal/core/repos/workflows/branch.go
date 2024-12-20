package workflows

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/core/repos/states"
)

func Branch(ctx workflow.Context, state *states.Branch) error {
	state.Init(ctx)

	selector := workflow.NewSelector(ctx)

	// - activity monitors -

	state.PullRequestMonitor(ctx)
	state.StaleMonitor(ctx)

	// - signal handlers -

	push := workflow.GetSignalChannel(ctx, defs.SignalPush.String())
	selector.AddReceive(push, state.OnPush(ctx))

	rebase := workflow.GetSignalChannel(ctx, defs.SignalRebase.String())
	selector.AddReceive(rebase, state.OnRebase(ctx))

	label := workflow.GetSignalChannel(ctx, defs.SignalPullRequestLabel.String())
	selector.AddReceive(label, state.OnLabel(ctx))

	prr := workflow.GetSignalChannel(ctx, defs.SignalPullRequestReview.String())
	selector.AddReceive(prr, state.OnPrReview(ctx))

	// - event loop -

	for !state.ExitLoop(ctx) {
		selector.Select(ctx)
	}

	// - exit or continue -

	if state.RestartRecommended(ctx) {
		return workflow.NewContinueAsNewError(ctx, Branch, state)
	}

	return nil
}
