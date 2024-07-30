package core

import (
	"go.temporal.io/sdk/workflow"
)

// TrunkCtrl is a workflow function that controls the trunk.
func TrunkCtrl(ctx workflow.Context, repo *Repo) error {
	state := NewTrunkState(ctx, repo)
	selector := workflow.NewSelector(ctx)

	// channels
	// push event
	push := workflow.GetSignalChannel(ctx, RepoIOSignalPush.String())
	selector.AddReceive(push, state.on_push(ctx))

	// create_delete
	create_delete := workflow.GetSignalChannel(ctx, RepoIOSignalCreateOrDelete.String())
	selector.AddReceive(create_delete, state.on_create_delete(ctx))

	for state.is_active() {
		selector.Select(ctx)
	}

	return nil
}
