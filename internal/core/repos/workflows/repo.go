package workflows

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/core/repos/states"
)

// Repo orchestrates repository workflows, routing incoming events. It initializes RepoState, registers query/signal
// handlers, and enters an event loop for workflow event processing. Workflow persistence spans the repo lifecycle,
// leveraging Temporal's continue-as-new feature to mitigate history size limitations.
func Repo(ctx workflow.Context, state *states.Repo) error {
	state.Init(ctx)

	selector := workflow.NewSelector(ctx)

	// - query handlers -
	if err := workflow.SetQueryHandler(ctx, defs.QueryRepoForEventParent.String(), state.QueryBranchTrigger); err != nil {
		return err
	}

	// - signal handlers -

	ref := workflow.GetSignalChannel(ctx, defs.SignalRef.String())
	selector.AddReceive(ref, state.OnRef(ctx))

	push := workflow.GetSignalChannel(ctx, defs.SignalPush.String())
	selector.AddReceive(push, state.OnPush(ctx))

	// - event loop -

	for !state.RestartRecommended(ctx) {
		selector.Select(ctx)
	}

	return workflow.NewContinueAsNewError(ctx, Repo, state)
}
