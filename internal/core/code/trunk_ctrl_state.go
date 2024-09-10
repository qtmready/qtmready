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

type (
	TrunkState struct {
		*BaseCtrl
		active_branch string
	}
)

func (state *TrunkState) on_push(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		push := &defs.RepoIOSignalPushPayload{}
		state.rx(ctx, rx, push)

		for _, branch := range state.branches {
			if branch == BranchNameFromRef(push.BranchRef) {
				continue
			}

			state.signal_branch(ctx, branch, defs.RepoIOSignalRebase, push)
		}
	}
}

func (state *TrunkState) on_create_delete(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		create_delete := &defs.RepoIOSignalCreateOrDeletePayload{}
		state.rx(ctx, rx, create_delete)

		if create_delete.ForBranch(ctx) {
			if create_delete.IsCreated {
				state.add_branch(ctx, create_delete.Ref)
			} else {
				state.remove_branch(ctx, create_delete.Ref)
			}
		}
	}
}

func NewTrunkState(ctx workflow.Context, repo *defs.Repo) *TrunkState {
	return &TrunkState{
		BaseCtrl:      NewBaseCtrl(ctx, "trunk_ctrl", repo),
		active_branch: repo.DefaultBranch,
	}
}
