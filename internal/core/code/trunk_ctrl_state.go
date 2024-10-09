// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
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
	"context"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/defs"
)

type (
	TrunkCtrlState struct {
		*BaseState
	}
)

// on_push handles push events for the trunk.
func (state *TrunkCtrlState) on_push(ctx workflow.Context) defs.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		push := &defs.Event[defs.Push, defs.RepoProvider]{} // Use Event type
		state.rx(ctx, rx, push)

		triggers := state.query__branch_triggers(ctx)

		state.log(ctx, "on_push").Info("triggers", triggers)

		for branch, parent_id := range triggers {
			rebase := ToRebaseEvent(ctx, push, branch, parent_id)
			state.signal_branch(ctx, branch, defs.RepoIOSignalRebase, rebase) // TODO: Fix the signal type
			state.persist(ctx, rebase)
		}
	}
}

func (state *TrunkCtrlState) restore(ctx workflow.Context) {
	state.BaseState.restore(ctx)
}

func NewTrunkCtrlState(ctx context.Context, repo *defs.Repo, info *defs.RepoIOProviderInfo) *TrunkCtrlState {
	return &TrunkCtrlState{
		BaseState: NewBaseState(ctx, "trunk_ctrl", repo, info, repo.DefaultBranch),
	}
}
