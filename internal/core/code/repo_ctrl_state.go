package code

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/shared"
)

// RepoCtrlState defines the state for RepoWorkflows.RepoCtrl.
// It embeds base_ctrl to inherit common functionality.
type (
	RepoCtrlState struct {
		*BaseCtrl
	}
)

// on_push is a channel handler that processes push events for the repository.
// It receives a RepoIOSignalPushPayload and signals the corresponding branch.
func (state *RepoCtrlState) on_push(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		push := &defs.RepoIOSignalPushPayload{}
		state.rx(ctx, rx, push)
		state.signal_branch(ctx, BranchNameFromRef(push.BranchRef), defs.RepoIOSignalPush, push)
	}
}

// on_create_delete is a channel handler that processes create or delete events for the repository.
// It receives a RepoIOSignalCreateOrDeletePayload, signals the corresponding branch,
// and updates the branch list in the state.
func (state *RepoCtrlState) on_create_delete(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		create_delete := &defs.RepoIOSignalCreateOrDeletePayload{}
		state.rx(ctx, rx, create_delete)

		if create_delete.ForBranch(ctx) {
			state.signal_branch(ctx, create_delete.Ref, defs.RepoIOSignalCreateOrDelete, create_delete)

			if create_delete.IsCreated {
				state.add_branch(ctx, create_delete.Ref)
			} else {
				state.remove_branch(ctx, create_delete.Ref)
			}
		}
	}
}

// on_pr is a channel handler that processes pull request events for the repository.
// It receives a RepoIOSignalPullRequestPayload and signals the corresponding branch.
func (state *RepoCtrlState) on_pr(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		pr := &defs.RepoIOSignalPullRequestPayload{}
		state.rx(ctx, rx, pr)
		state.signal_branch(ctx, pr.HeadBranch, defs.RepoIOSignalPullRequest, pr)
	}
}

// NewRepoCtrlState creates a new RepoCtrlState with the specified repo.
// It initializes the embedded base_ctrl using NewBaseCtrl.
func NewRepoCtrlState(ctx workflow.Context, repo *defs.Repo) *RepoCtrlState {
	return &RepoCtrlState{
		BaseCtrl: NewBaseCtrl(ctx, "repo_ctrl", repo),
	}
}
