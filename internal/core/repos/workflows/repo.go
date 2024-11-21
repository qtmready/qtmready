package reposwfs

import (
	"go.temporal.io/sdk/workflow"

	reposdefs "go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/durable"
)

type (
	// RepoState defines the state for Repo Workflows. It embeds BaseState to inherit its functionality.
	RepoState struct {
		*BaseState `json:"base"` // Base workflow state.
	}
)

// Repo manages the event loop for a repository. It acts as a central router, orchestrating repository workflows.
//
// Repo uses Temporal's workflow selector to concurrently handle signals. The function initializes a RepoState and
// registers a signal handler for the reposdefs.RepoIOSignalPush signal.  Currently, the signal handler is a stub.
// Temporal workflow context is passed as input, along with hydrated repository data.  The function returns an error
// if one occurs during workflow execution; otherwise it returns nil.
//
// FIXME: start the function with the state instead of HydratedRepo.If we start from a state, It will be easier to
// restart the workflow as new.
func Repo(ctx workflow.Context, hydrated_repo *reposdefs.HypdratedRepo) error {
	selector := workflow.NewSelector(ctx)

	state := NewRepoState(ctx, hydrated_repo)

	push := workflow.GetSignalChannel(ctx, reposdefs.RepoIOSignalPush.String())
	selector.AddReceive(push, state.on_push(ctx))

	for !state.restart(ctx) {
		selector.Select(ctx)
	}

	return workflow.NewContinueAsNewError(ctx, Repo, hydrated_repo) //
}

// - signal handlers -

// on_push is a signal handler for the push signal.
func (state *RepoState) on_push(ctx workflow.Context) durable.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		state.rx(ctx, rx, nil)
	}
}

// - state managers -

func (state *RepoState) restart(ctx workflow.Context) bool {
	return workflow.GetInfo(ctx).GetContinueAsNewSuggested()
}

// NewRepoState creates a new RepoState instance. It initializes BaseState using the provided context and
// hydrated repository data.
func NewRepoState(ctx workflow.Context, hydrated *reposdefs.HypdratedRepo) *RepoState {
	base := NewBaseState(ctx, hydrated) // Initialize BaseState.
	return &RepoState{base}             // Return new RepoState instance.
}
