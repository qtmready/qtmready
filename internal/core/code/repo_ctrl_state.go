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

// RepoCtrlState defines the state for RepoWorkflows.RepoCtrl. It embeds base_ctrl to inherit common functionality.
type (
	RepoCtrlState struct {
		*BaseState
		triggers BranchTriggers
		stash    StashedEvents[defs.RepoProvider]
	}
)

// on_push is a channel handler that processes push events for the repository. It receives a RepoIOSignalPushPayload and
// signals the corresponding branch.
func (state *RepoCtrlState) on_push(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		event := &defs.Event[defs.Push, defs.RepoProvider]{}
		state.rx(ctx, rx, event)

		branch := BranchNameFromRef(event.Payload.Ref)
		id, ok := state.triggers.get(event.Payload.Ref)

		if ok {
			event.SetParent(id)
			state.signal_branch(ctx, branch, defs.RepoIOSignalPush, event)
		} else {
			state.log(ctx, "on_push").Warn("no parent id found, stashing ...", "event_id", event.ID)
			state.stash.push(branch, event)
		}
	}
}

// on_create_delete is a channel handler that processes create or delete events for the repository. It receives a
// defs.Event[defs.BranchOrTag, defs.RepoProvider], signals the corresponding branch, and updates the branch list in the
// state.
func (state *RepoCtrlState) on_create_delete(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		event := &defs.Event[defs.BranchOrTag, defs.RepoProvider]{}
		state.rx(ctx, rx, event)

		if event.Context.Scope == defs.EventScopeBranch {
			// Assuming signal_branch is updated to accept the Event:
			state.signal_branch(ctx, event.Payload.Ref, defs.RepoIOSignalCreateOrDelete, event) // TODO: fix the payload

			if event.Context.Action == defs.EventActionCreated {
				state.add_branch(ctx, event.Payload.Ref)
				state.triggers.add(event.Payload.Ref, event.ID)

				events, ok := state.stash.all(event.Payload.Ref)
				if ok {
					for _, each := range events {
						each.SetParent(event.ID)
						state.signal_branch(ctx, event.Payload.Ref, defs.RepoIOSignalPush, each)
					}
				} else {
					state.log(ctx, "on_create_delete").Warn("no stashed events found.")
				}
			} else if event.Context.Action == defs.EventActionDeleted {
				state.remove_branch(ctx, event.Payload.Ref)
				state.triggers.del(event.Payload.Ref)
			}
		}
	}
}

// on_pr is a channel handler that processes pull request events for the repository. It receives a
// RepoIOSignalPullRequestPayload and signals the corresponding branch.
func (state *RepoCtrlState) on_pr(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		event := &defs.Event[defs.PullRequest, defs.RepoProvider]{}
		state.rx(ctx, rx, event)

		branch := event.Payload.HeadBranch
		id, ok := state.triggers.get(branch)

		if ok {
			event.SetParent(id)
			state.signal_branch(ctx, branch, defs.RepoIOSignalPullRequest, event)
		} else {
			state.log(ctx, "on_pr").Warn("no parent id found, stashing ...", "event_id", event.ID)
			state.stash.push(branch, event)
		}
	}
}

func (state *RepoCtrlState) on_label(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		label := &defs.Event[defs.PullRequestLabel, defs.RepoProvider]{}
		state.rx(ctx, rx, label)

		branch := label.Payload.Branch
		id, ok := state.triggers.get(branch)

		if ok {
			label.SetParent(id)
			state.signal_branch(ctx, branch, defs.RepoIOSignalPullRequestLabel, label)
		} else {
			state.log(ctx, "on_label").Warn("no parent id found, stashing ...", "event_id", label.ID)
			state.stash.push(branch, label)
		}
	}
}

// NewRepoCtrlState creates a new RepoCtrlState with the specified repo. It initializes the embedded base_ctrl using
// NewBaseCtrl.
func NewRepoCtrlState(ctx workflow.Context, repo *defs.Repo) *RepoCtrlState {
	return &RepoCtrlState{
		BaseState: NewBaseCtrl(ctx, "repo_ctrl", repo),
		triggers:  make(BranchTriggers),
		stash:     make(StashedEvents[defs.RepoProvider]),
	}
}
