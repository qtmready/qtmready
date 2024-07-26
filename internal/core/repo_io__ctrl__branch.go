package core

import (
	"go.temporal.io/sdk/workflow"
)

func BranchCtrl(ctx workflow.Context, repo *Repo, branch string) error {
	selector := workflow.NewSelector(ctx)
	state := NewBranchCtrlState(ctx, repo, branch)

	// start the stale check coroutine.
	state.check_stale(ctx)

	// setup signals

	// push event signal.
	// detect changes. if changes are greater than threshold, send early warning message.
	push := workflow.GetSignalChannel(ctx, RepoIOSignalPush.String())
	selector.AddReceive(push, state.on_push(ctx))

	// rebase signal.
	// attempts to rebase the branch with the base branch. if there are merge conflicts, sends message.
	rebase := workflow.GetSignalChannel(ctx, ReopIOSignalRebase.String())
	selector.AddReceive(rebase, state.on_rebase(ctx))

	// create_delete signal.
	// creates or deletes the branch.
	create_delete := workflow.GetSignalChannel(ctx, RepoIOSignalCreateOrDelete.String())
	selector.AddReceive(create_delete, state.on_create_delete(ctx))

	// pr signal.
	pr := workflow.GetSignalChannel(ctx, RepoIOSignalPullRequest.String())
	selector.AddReceive(pr, state.on_pr(ctx))

	// listen to all signals

	for state.is_active() {
		selector.Select(ctx)
	}

	return nil
}
