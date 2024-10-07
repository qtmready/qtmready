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
	"github.com/gocql/gocql"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/defs"
)

// RepoCtrlState defines the state for RepoWorkflows.RepoCtrl. It embeds base_ctrl to inherit common functionality.
type (
	RepoCtrlState struct {
		*BaseState                                  // Embedded base state for common workflow logic.
		stash      StashedEvents[defs.RepoProvider] // Storage for events that have no parent yet.
		triggers   BranchTriggers                   // Map of branch names to event IDs, used to track event dependencies.
	}
)

// on_push is a channel handler that processes push events for the repository. It receives a RepoIOSignalPushPayload and
// signals the corresponding branch.
func (state *RepoCtrlState) on_push(ctx workflow.Context) defs.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		event := &defs.Event[defs.Push, defs.RepoProvider]{}
		state.rx(ctx, rx, event)

		if !state.on_trunk(BranchNameFromRef(event.Payload.Ref)) {
			if parent, ok := state.triggers.get(event.Payload.Ref); ok {
				event.SetParent(parent)
				state.signal_branch(ctx, event.Payload.Ref, defs.RepoIOSignalPush, event)
				state.persist(ctx, event)

				return
			}

			state.stash.push(event.Payload.Ref, event)
		}

		state.signal_branch(ctx, state.Repo.DefaultBranch, defs.RepoIOSignalPush, event)
		state.persist(ctx, event)
	}
}

// on_create_delete is a channel handler that processes create or delete events for the repository. It receives a
// defs.Event[defs.BranchOrTag, defs.RepoProvider], signals the corresponding branch, and updates the branch list in the
// state.
func (state *RepoCtrlState) on_create_delete(ctx workflow.Context) defs.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		event := &defs.Event[defs.BranchOrTag, defs.RepoProvider]{}
		state.rx(ctx, rx, event)

		if event.Context.Scope == defs.EventScopeBranch && event.Context.Action == defs.EventActionCreated {
			state.on_branch_create(ctx, event)

			return
		}

		if event.Context.Scope == defs.EventScopeBranch && event.Context.Action == defs.EventActionDeleted {
			state.on_branch_delete(ctx, event)

			return
		}

		state.log(ctx, "on_create_delete").Warn("unhandled event", "kind", event.Context.Scope, "action", event.Context.Action)
	}
}

// on_pr is a channel handler that processes pull request events for the repository. It receives a
// RepoIOSignalPullRequestPayload and signals the corresponding branch.
func (state *RepoCtrlState) on_pr(ctx workflow.Context) defs.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		event := &defs.Event[defs.PullRequest, defs.RepoProvider]{}
		state.rx(ctx, rx, event)

		parent, ok := state.triggers.get(event.Payload.HeadBranch)
		if !ok {
			state.log(ctx, "on_pr").Warn("attempting to process pr without parent, this may result in orphaned events")
		} else {
			event.SetParent(parent)
		}

		state.persist(ctx, event)
		state.signal_branch(ctx, event.Payload.HeadBranch, defs.RepoIOSignalPullRequest, event)
	}
}

// on_label is a channel handler that processes label events for the repository. It receives a
// RepoIOSignalPullRequestLabelPayload and signals the corresponding branch.
func (state *RepoCtrlState) on_label(ctx workflow.Context) defs.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		event := &defs.Event[defs.PullRequestLabel, defs.RepoProvider]{}
		state.rx(ctx, rx, event)

		parent, ok := state.triggers.get(event.Payload.Branch)
		if !ok {
			state.log(ctx, "on_label").Warn("attempting to process label without parent, this may result in orphaned events")
		} else {
			event.SetParent(parent)
		}

		state.persist(ctx, event)
	}
}

func (state *RepoCtrlState) on_branch_create(ctx workflow.Context, event *defs.Event[defs.BranchOrTag, defs.RepoProvider]) {
	state.add_branch(ctx, event.Payload.Ref)
	state.triggers.add(event.Payload.Ref, event.ID)

	events, ok := state.stash.all(event.Payload.Ref)
	if ok {
		for _, each := range events {
			each.SetParent(event.ID)
			state.signal_branch(ctx, event.Payload.Ref, defs.RepoIOSignalPush, each)
			state.persist(ctx, each)
		}
	}
}

func (state *RepoCtrlState) on_branch_delete(ctx workflow.Context, event *defs.Event[defs.BranchOrTag, defs.RepoProvider]) {
	state.remove_branch(ctx, event.Payload.Ref)

	parent, ok := state.triggers.get(event.Payload.Ref)
	if ok {
		event.SetParent(parent)
		state.persist(ctx, event)
	}

	state.triggers.del(event.Payload.Ref)
}

func (state *RepoCtrlState) on_trunk(branch string) bool {
	return branch == state.Repo.DefaultBranch
}

func (state *RepoCtrlState) setup_query__get_parent(ctx workflow.Context) error {
	logger := state.log(ctx, "query")

	return workflow.SetQueryHandler(ctx, QueryRepoGetParentForBranch.String(), func(branch string) gocql.UUID {
		logger.Info("getting parent ...", "branch", branch)

		parent, ok := state.triggers.get(branch)
		if !ok {
			logger.Warn("no parent found for branch", "branch", branch)

			return parent
		}

		logger.Info("parent found", "branch", branch, "parent", parent)

		return parent
	})
}

// NewRepoCtrlState creates a new RepoCtrlState with the specified repo. Embedded BaseState is initialized using NewBaseState.
func NewRepoCtrlState(ctx workflow.Context, repo *defs.Repo) *RepoCtrlState {
	return &RepoCtrlState{
		BaseState: NewBaseState(ctx, "repo_ctrl", repo),
		stash:     make(StashedEvents[defs.RepoProvider]),
		triggers:  make(BranchTriggers),
	}
}
