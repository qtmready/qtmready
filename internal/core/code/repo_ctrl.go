package code

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/defs"
)

// RepoCtrl is the event loop to process events during the lifecycle of a repository.
//
// It processes the following events:
//
//   - push
//   - create_delete
//   - pr
func RepoCtrl(ctx workflow.Context, repo *defs.Repo) error {
	state := NewRepoCtrlState(ctx, repo)
	selector := workflow.NewSelector(ctx)

	// setup
	state.refresh_info(ctx)
	state.refresh_branches(ctx)

	// channels
	// push event
	push := workflow.GetSignalChannel(ctx, defs.RepoIOSignalPush.String())
	selector.AddReceive(push, state.on_push(ctx))

	// create_delete event
	create_delete := workflow.GetSignalChannel(ctx, defs.RepoIOSignalCreateOrDelete.String())
	selector.AddReceive(create_delete, state.on_create_delete(ctx))

	// pull request event
	pr := workflow.GetSignalChannel(ctx, defs.RepoIOSignalPullRequest.String())
	selector.AddReceive(pr, state.on_pr(ctx))

	// main event loop
	for state.is_active() {
		selector.Select(ctx)

		if state.needs_reset() {
			return state.as_new(ctx, "event history exceeded threshold", RepoCtrl, repo)
		}
	}

	// graceful shutdown
	state.terminate(ctx)

	return nil
}
