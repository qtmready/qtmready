package core

import (
	"go.temporal.io/sdk/workflow"
)

// RepoCtrl is the main control loop for managing a repository. It sets up signal channels
// for various repository events, processes those events, and gracefully shuts down when
// the context is canceled.
func RepoCtrl(ctx workflow.Context, repo *Repo) error {
	state := NewRepoCtrlState(ctx, repo)
	selector := workflow.NewSelector(ctx)

	// setup
	state.refresh_info(ctx)
	state.refresh_branches(ctx)

	// channels
	// push event
	push := workflow.GetSignalChannel(ctx, RepoIOSignalPush.String())
	selector.AddReceive(push, state.on_push(ctx))

	// create_delete event
	create_delete := workflow.GetSignalChannel(ctx, RepoIOSignalCreateOrDelete.String())
	selector.AddReceive(create_delete, state.on_create_delete(ctx))

	// pull request event
	pr := workflow.GetSignalChannel(ctx, RepoIOSignalPullRequest.String())
	selector.AddReceive(pr, state.on_pr(ctx))

	// processing signals
	for state.is_active() {
		selector.Select(ctx)
	}

	// graceful shutdown
	state.terminate(ctx)

	return nil
}
