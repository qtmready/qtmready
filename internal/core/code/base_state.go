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
	"log/slog"
	"time"

	"github.com/gocql/gocql"
	"go.breu.io/durex/dispatch"
	"go.breu.io/durex/queues"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/core/kernel"
)

type (
	// BaseState represents the base state for repository operations. It provides common functionality for various
	// repository control types.
	BaseState struct {
		ActiveBranch string                   `json:"active_branch"` // active_branch is the branch currently being processed
		Kind         string                   `json:"kind"`          // kind identifies the type of control (e.g., "repo", "branch")
		Repo         *defs.Repo               `json:"repo"`          // repo is a reference to the repository
		Info         *defs.RepoIOProviderInfo `json:"info"`          // info stores provider-specific information
		Branches     []string                 `json:"branches"`      // branches is a list of branches in the repository
		Active       bool                     `json:"active"`        // active indicates if the control is still active

		activities *Activities    // activities holds the repository activities
		mutex      workflow.Mutex // mutex is used for thread-safe operations
	}

	// RepoEvent defines an interface for repository events. It simplifies working with repository events by
	// providing a common interface over various event types.
	//
	// This interface avoids using generics by leveraging the `Flatten` method, which is available on all
	// efs.Event[defs.EventPayload, defs.EventProvider] combinations, to streamline event operations.
	RepoEvent[P defs.RepoProvider] interface {
		// Flatten converts the RepoEvent into a defs.FlatEvent, which is a more efficient representation for storage and
		// retrieval.
		Flatten() (*defs.FlatEvent[P], error)

		// SetParent sets the parent event ID for the current event, enabling the reconstruction of the event lineage.
		SetParent(gocql.UUID)

		UnmarshalJSON(data []byte) error
	}
)

// needs_reset checks if the event count has reached the threshold for resetting.
func (base *BaseState) needs_reset(ctx workflow.Context) bool {
	return workflow.GetInfo(ctx).GetContinueAsNewSuggested()
}

// branch returns the branch name associated with this control.
func (base *BaseState) branch(_ workflow.Context) string {
	return base.ActiveBranch
}

// set_branch sets the active branch in the context.
func (base *BaseState) set_branch(ctx workflow.Context, branch string) workflow.Context {
	return workflow.WithValue(ctx, "active_branch", branch)
}

// set_info sets the provider-specific information for the control.
func (base *BaseState) set_info(ctx workflow.Context, info *defs.RepoIOProviderInfo) {
	_ = base.mutex.Lock(ctx)
	defer base.mutex.Unlock()

	base.Info = info
}

// set_branches sets the list of branches associated with the control.
func (base *BaseState) set_branches(ctx workflow.Context, branches []string) {
	_ = base.mutex.Lock(ctx)
	defer base.mutex.Unlock()

	base.Branches = branches
}

func (base *BaseState) is_active() bool {
	return base.Active
}

// set_done marks the control as inactive.
func (base *BaseState) set_done(ctx workflow.Context) {
	_ = base.mutex.Lock(ctx)
	defer base.mutex.Unlock()

	base.Active = false
}

// terminate marks the control as done and logs the termination.
func (base *BaseState) terminate(ctx workflow.Context) {
	base.set_done(ctx)
	base.log(ctx, "terminate").Info("state terminated")
}

// as_new continues the workflow as new with the given function and arguments.
func (base *BaseState) as_new(ctx workflow.Context, msg string, fn any, args ...any) error {
	base.log(ctx, "as_new").Warn(msg)
	return workflow.NewContinueAsNewError(ctx, fn, args...)
}

// add_branch adds a new branch to the list of branches.
func (base *BaseState) add_branch(ctx workflow.Context, branch string) {
	_ = base.mutex.Lock(ctx)
	defer base.mutex.Unlock()

	if branch != "" || branch != base.Repo.DefaultBranch {
		base.Branches = append(base.Branches, branch)
	}
}

// remove_branch removes a branch from the list of branches.
func (base *BaseState) remove_branch(ctx workflow.Context, branch string) {
	_ = base.mutex.Lock(ctx)
	defer base.mutex.Unlock()

	for i, b := range base.Branches {
		if b == branch {
			base.Branches = append(base.Branches[:i], base.Branches[i+1:]...)
			break
		}
	}
}

// signal_branch sends a signal to a specific branch.
func (base *BaseState) signal_branch(ctx workflow.Context, branch string, signal queues.Signal, payload any) {
	ctx = dispatch.WithDefaultActivityContext(ctx)

	next := &defs.RepoIOSignalBranchCtrlPayload{
		Repo:    base.Repo,
		Info:    base.Info,
		Branch:  branch,
		Signal:  signal,
		Payload: payload,
	}

	_ = base.do(
		ctx, "signal_branch_ctrl", base.activities.SignalBranch, next, nil,
		slog.String("signal", signal.String()),
		slog.String("branch", branch),
	)
}

// TODO - refine the logic.
// signal_branch sends a signal to a specific branch.
func (base *BaseState) signal_queue(ctx workflow.Context, branch string, signal queues.Signal, payload any) {
	ctx = dispatch.WithDefaultActivityContext(ctx)

	next := &defs.RepoIOSignalQueueCtrlPayload{
		Repo:    base.Repo,
		Branch:  branch,
		Signal:  signal,
		Payload: payload,
	}

	_ = base.do(
		ctx, "signal_queue_ctrl", base.activities.SignalQueue, next, nil,
		slog.String("signal", signal.String()),
		slog.String("branch", branch),
	)
}

// rx receives a message from a channel and logs the event.
func (base *BaseState) rx(ctx workflow.Context, channel workflow.ReceiveChannel, target any) {
	base.log(ctx, "rx").Info(channel.Name())

	channel.Receive(ctx, target)
}

// refresh_info updates the provider information for the repository.
func (base *BaseState) refresh_info(ctx workflow.Context) {
	ctx = dispatch.WithDefaultActivityContext(ctx)

	info := &defs.RepoIOProviderInfo{}
	io := kernel.Instance().RepoIO(base.Repo.Provider)

	_ = base.do(ctx, "get_repo_info", io.GetProviderInfo, base.Repo.CtrlID, info)
	base.set_info(ctx, info)
}

// refresh_branches updates the list of branches for the repository.
func (base *BaseState) refresh_branches(ctx workflow.Context) {
	if base.Info == nil {
		base.refresh_info(ctx)
	}

	io := kernel.Instance().RepoIO(base.Repo.Provider)
	branches := []string{}

	ctx = dispatch.WithDefaultActivityContext(ctx)

	_ = base.do(ctx, "refresh_branches", io.GetAllBranches, base.Info, &branches)
	base.set_branches(ctx, branches)
}

func (base *BaseState) persist(ctx workflow.Context, event RepoEvent[defs.RepoProvider]) {
	ctx = dispatch.WithDefaultActivityContext(ctx)

	flat, _ := event.Flatten()
	_ = base.do(ctx, "persist", base.activities.SaveRepoEvent, flat, nil)
}

// log creates a new logger for the current action.
func (base *BaseState) log(ctx workflow.Context, action string) *RepoIOWorkflowLogger {
	return NewRepoIOWorkflowLogger(ctx, base.Repo, base.Kind, base.ActiveBranch, action)
}

func (state *BaseState) query__parent_event_id(ctx workflow.Context, branch string) (gocql.UUID, bool) {
	ctx = dispatch.WithDefaultActivityContext(ctx)
	payload := &RepoCtrlQueryPayloadForBranchParent{Branch: branch, Repo: state.Repo}
	result := &RepoCtrlQueryResultForBranchParent{}

	err := state.do(ctx, "query__parent_event_id", state.activities.QueryRepoCtrlForBranchParent, payload, result)
	if err != nil {
		return result.ID, false
	}

	return result.ID, result.Found
}

func (state *BaseState) query__branch_triggers(ctx workflow.Context) BranchTriggers {
	ctx = dispatch.WithDefaultActivityContext(ctx)
	triggers := make(BranchTriggers)
	repo := state.Repo

	state.log(ctx, "query__branch_triggers").Info("querying ...", "repo", repo)

	_ = state.do(ctx, "hello", state.activities.QueryRepoCtrlForBranchTriggers, repo, triggers)

	// err := state.do(ctx, "query__branch_triggers", state.activities.QueryRepoCtrlForBranchTriggers, repo, &triggers)
	// if err != nil {
	// 	return triggers
	// }

	return triggers
}

// do is helper is an activity executor. It logs the activity execution and increments the operation counter.
func (base *BaseState) do(ctx workflow.Context, action string, activity, payload, result any, keyvals ...any) error {
	logger := base.log(ctx, action)
	logger.Info("init", keyvals...)

	if err := workflow.ExecuteActivity(ctx, activity, payload).Get(ctx, result); err != nil {
		logger.Warn("error", append(keyvals, "error", err)...)
		return err
	}

	logger.Info("success", keyvals...)

	return nil
}

// child executes a child workflow and logs the event.
//
// Deprecated: use signal with start to other workflows. implement.
func (base *BaseState) child(ctx workflow.Context, action, w_id string, fn, payload any, keyvals ...any) error {
	logger := base.log(ctx, action)
	logger.Info("init", keyvals...)

	opts := workflow.ChildWorkflowOptions{
		TaskQueue:                "quantm_queue", // TODO - queue name
		WorkflowExecutionTimeout: 10 * time.Minute,
		WorkflowID:               w_id,
	}
	ctx = workflow.WithChildOptions(ctx, opts)

	// Execute the child workflow
	err := workflow.ExecuteChildWorkflow(ctx, fn, payload).Get(ctx, nil)
	if err != nil {
		logger.Warn("error", append(keyvals, "error", err)...)
		return err
	}

	logger.Info("success", keyvals...)

	return nil
}

func (base *BaseState) restore(ctx workflow.Context) {
	base.mutex = workflow.NewMutex(ctx)
}

// NewBaseState creates a new base control instance. This is the preferred method to create a new base control instance.
func NewBaseState(ctx context.Context, kind string, repo *defs.Repo, info *defs.RepoIOProviderInfo, branch string) *BaseState {
	return &BaseState{
		activities:   &Activities{},
		ActiveBranch: branch,
		Kind:         kind,
		Repo:         repo,
		Info:         info,
		Branches:     make([]string, 0),
		Active:       true,
	}
}
