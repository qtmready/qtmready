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
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/shared"
)

type (
	TrunkCtrlState struct {
		*BaseState
		active_branch string
	}
)

// on_push handles push events for the trunk.
func (state *TrunkCtrlState) on_push(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		push := &defs.Event[defs.Push, defs.RepoProvider]{} // Use Event type
		state.rx(ctx, rx, push)

		for _, branch := range state.branches {
			if branch == BranchNameFromRef(push.Payload.Ref) {
				continue
			}

			state.signal_branch(ctx, branch, defs.RepoIOSignalRebase, push) // TODO: Fix the signal type
		}
	}
}

func (state *TrunkCtrlState) on_create_delete(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		event := &defs.Event[defs.BranchOrTag, defs.RepoProvider]{}
		state.rx(ctx, rx, event)

		if event.Context.Scope == defs.EventScopeBranch {
			if event.Context.Action == defs.EventActionCreated {
				state.add_branch(ctx, event.Payload.Ref)
			} else if event.Context.Action == defs.EventActionDeleted {
				state.remove_branch(ctx, event.Payload.Ref)
			}
		}
	}
}

func NewTrunkState(ctx workflow.Context, repo *defs.Repo) *TrunkCtrlState {
	return &TrunkCtrlState{
		BaseState:     NewBaseCtrl(ctx, "trunk_ctrl", repo),
		active_branch: repo.DefaultBranch,
	}
}
