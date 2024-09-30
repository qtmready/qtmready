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
)

// BranchCtrl is the event loop to process events during the lifecycle of a branch.
//
// It processes the following events:
//
//   - push
//   - rebase
//   - create_delete
//   - pr
func BranchCtrl(ctx workflow.Context, repo *defs.Repo, branch string) error {
	selector := workflow.NewSelector(ctx)
	ctx, state := NewBranchCtrlState(ctx, repo, branch)

	// start the stale check coroutine.
	state.check_stale(ctx)

	// setup signals

	// push event signal.
	// detect changes. if changes are greater than threshold, send early warning message.
	push := workflow.GetSignalChannel(ctx, defs.RepoIOSignalPush.String())
	selector.AddReceive(push, state.on_push(ctx))

	// rebase signal.
	// attempts to rebase the branch with the base branch. if there are merge conflicts, sends message.
	rebase := workflow.GetSignalChannel(ctx, defs.RepoIOSignalRebase.String())
	selector.AddReceive(rebase, state.on_rebase(ctx))

	// create_delete signal.
	// creates or deletes the branch.
	create_delete := workflow.GetSignalChannel(ctx, defs.RepoIOSignalCreateOrDelete.String())
	selector.AddReceive(create_delete, state.on_create_delete(ctx))

	// pr signal.
	pr := workflow.GetSignalChannel(ctx, defs.RepoIOSignalPullRequest.String())
	selector.AddReceive(pr, state.on_pr(ctx))

	// label signal.
	lebal := workflow.GetSignalChannel(ctx, defs.RepoIOSignalPullRequestLabel.String())
	selector.AddReceive(lebal, state.on_label(ctx))

	// main event loop
	for state.is_active() {
		selector.Select(ctx)

		// TODO - need to optimize
		// TODO - remove
		if state.pr != nil {
			_ctx, q_state := NewQueueCtrlState(ctx, repo, branch)
			q_state.push(_ctx, state.pr, false) // TODO - handle priority
		}

		if state.needs_reset() {
			return state.as_new(ctx, "event history exceeded threshold", BranchCtrl, repo, branch)
		}
	}

	// graceful shutdown
	state.terminate(ctx)

	return nil
}
