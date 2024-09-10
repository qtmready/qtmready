// Copyright © 2023, Breu, Inc. <info@breu.io>
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

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

		state.signal_branch(ctx, pr.HeadBranch, defs.RepoIOSignalPullRequestOpenedOrClosedOrReopened, pr)
	}
}

func (state *RepoCtrlState) on_label(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		label := &defs.RepoIOSignalPullRequestPayload{}
		state.rx(ctx, rx, label)

		state.signal_branch(ctx, label.HeadBranch, defs.RepoIOSignalPullRequestLabeledOrUnlabeled, label)
	}
}

// NewRepoCtrlState creates a new RepoCtrlState with the specified repo.
// It initializes the embedded base_ctrl using NewBaseCtrl.
func NewRepoCtrlState(ctx workflow.Context, repo *defs.Repo) *RepoCtrlState {
	return &RepoCtrlState{
		BaseCtrl: NewBaseCtrl(ctx, "repo_ctrl", repo),
	}
}
