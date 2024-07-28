package core

import (
	"log/slog"
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/shared"
)

type (
	TrunkState struct {
		activties *RepoActivities     // activities is the activities for the branch control state
		repo      *Repo               // repo is the db record of the repo
		info      *RepoIOProviderInfo // info is the provider info for the repo
		branches  []string            // branches is the list of branches for the repo except the default branch.
		mutex     workflow.Mutex      // mutex is the mutex for the repo control state
		active    bool                // active is whether the branch control state is currently active
		counter   int                 // counter is the number of activity executions
	}
)

func (state *TrunkState) on_push(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		push := &RepoIOSignalPushPayload{}
		rx.Receive(ctx, push)

		for _, branch := range state.branches {
			state.signal_branch(ctx, branch, RepoIOSignalRebase, push)
		}
	}
}

func (state *TrunkState) on_create_delete(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		create_delete := &RepoIOSignalCreateOrDeletePayload{}
		rx.Receive(ctx, create_delete)

		if create_delete.RefType == "branch" {
			if create_delete.IsCreated {
				state.add_branch(ctx, create_delete.Ref)
			} else {
				state.remove_branch(ctx, create_delete.Ref)
			}
		}
	}
}

// set_done marks the RepoCtrlState as inactive, releasing the mutex lock.
// This function should be called when the branch control state is no longer needed,
// such as when the branch is being deleted or merged.
func (state *TrunkState) set_done(ctx workflow.Context) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.active = false
}

// is_active returns whether the RepoCtrlState is currently active.
// When the state is active, it means the branch control state is in use and the mutex is locked.
func (state *TrunkState) is_active() bool {
	return state.active
}

func (state *TrunkState) refresh_info(ctx workflow.Context) {
	io := Instance().RepoIO(state.repo.Provider)
	info := &RepoIOProviderInfo{}

	_ = state.do(ctx, "refresh_provider_info", io.GetProviderInfo, state.repo.CtrlID.String(), info, nil)

	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.info = info
}

func (state *TrunkState) refresh_branches(ctx workflow.Context) {
	if state.info == nil {
		state.refresh_info(ctx)
	}

	branches := make([]string, 0)
	io := Instance().RepoIO(state.repo.Provider)

	_ = state.do(ctx, "refresh_branches", io.GetAllBranches, state.repo, &branches)

	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.branches = branches
}

func (state *TrunkState) add_branch(ctx workflow.Context, branch string) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	if branch != "" || branch != state.repo.DefaultBranch {
		state.branches = append(state.branches, branch)
	}
}

func (state *TrunkState) remove_branch(ctx workflow.Context, branch string) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	for i, b := range state.branches {
		if b == branch {
			state.branches = append(state.branches[:i], state.branches[i+1:]...)
			break
		}
	}
}

// signal_branch sends a signal to the branch control state for the specified branch.
func (state *TrunkState) signal_branch(ctx workflow.Context, branch string, signal shared.WorkflowSignal, payload any) {
	opts := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
	ctx = workflow.WithActivityOptions(ctx, opts)

	next := &RepoIOSignalBranchCtrlPayload{state.repo, branch, signal, payload}

	_ = state.do(
		ctx, "signal_branch_ctrl", state.activties.SignalBranch, next, nil,
		slog.String("signal", signal.String()),
	)
}

// increment is a helper function that increments the steps counter in the RepoIOBranchCtrlState.
func (state *TrunkState) increment(ctx workflow.Context) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.counter++
}

// terminate is called to mark the RepoIOBranchCtrlState as inactive and cancel any associated timers.
// This function should be called when the branch control state is no longer needed, such as branch is being deleted or merged.
func (state *TrunkState) terminate(ctx workflow.Context) {
	state.set_done(ctx)
}

// do is a helper function that executes an activity within the context of the RepoIOBranchCtrlState.
// It logs the start and success of the activity, and increments the steps counter in the state.
// If the activity fails, the function logs the error and returns it.
//
// NOTE: This assumes that workflow.Context has been updated with activity options.
func (state *TrunkState) do(ctx workflow.Context, action string, activity, payload, result any, keyvals ...any) error {
	return _do(ctx, state.repo, state.repo.DefaultBranch, "trunk_ctrl", action, activity, payload, result, keyvals...)
}

// NewTrunkState creates a new TrunkState with the specified repo and activities.
func NewTrunkState(ctx workflow.Context, repo *Repo) *TrunkState {
	return &TrunkState{
		activties: &RepoActivities{},
		repo:      repo,
		branches:  make([]string, 0),
		mutex:     workflow.NewMutex(ctx),
		active:    true,
		counter:   0,
	}
}
