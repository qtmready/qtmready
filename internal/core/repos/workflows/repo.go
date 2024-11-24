package workflows

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/core/repos/states"
)

type (
	// RepoState defines the state for Repo Workflows. It embeds BaseState to inherit its functionality.
	RepoState struct {
		*BaseState `json:"base"` // Base workflow state.

		Triggers BranchTriggers `json:"triggers"` // Branch triggers.
	}
)

// Repo manages the event loop for a repository. It acts as a central router, orchestrating repository workflows.
//
// Repo uses Temporal's workflow selector to concurrently handle signals. The function initializes a RepoState and
// registers a signal handler for the reposdefs.RepoIOSignalPush signal.  Currently, the signal handler is a stub.
// Temporal workflow context is passed as input, along with hydrated repository data.  The function returns an error
// if one occurs during workflow execution; otherwise it returns nil.
func Repo(ctx workflow.Context, state *states.Repo) error {
	state.Init(ctx)

	selector := workflow.NewSelector(ctx)

	// - query handlers -
	if err := workflow.SetQueryHandler(ctx, defs.QueryRepoForEventParent.String(), state.QueryBranchTrigger); err != nil {
		return err
	}

	// - signal handlers -

	push := workflow.GetSignalChannel(ctx, defs.SignalPush.String())
	selector.AddReceive(push, state.OnPush(ctx))

	// - event loop -

	for !state.RestartRecommended(ctx) {
		selector.Select(ctx)
	}

	return workflow.NewContinueAsNewError(ctx, Repo, state)
}
