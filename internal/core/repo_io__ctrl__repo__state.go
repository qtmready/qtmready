package core

import (
	"log/slog"
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/shared"
)

type (
	RepoCtrlState struct {
		activties *RepoActivities // activities is the activities for the branch control state
		repo      *Repo           // repo is the db record of the repo
		branches  []string        // branches is the list of branches for the repo except the default branch.
		mutex     workflow.Mutex  // mutex is the mutex for the repo control state
		active    bool
		counter   int // counter is the number of activity executions
	}
)

// set_done marks the RepoCtrlState as inactive, releasing the mutex lock.
// This function should be called when the branch control state is no longer needed,
// such as when the branch is being deleted or merged.
func (state *RepoCtrlState) set_done(ctx workflow.Context) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.active = false
}

// is_active returns whether the RepoCtrlState is currently active.
// When the state is active, it means the branch control state is in use and the mutex is locked.
func (state *RepoCtrlState) is_active() bool {
	return state.active
}

func (state *RepoCtrlState) add_branch(ctx workflow.Context, branch string) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	if branch != "" || branch != state.repo.DefaultBranch {
		state.branches = append(state.branches, branch)
	}
}

func (state *RepoCtrlState) remove_branch(ctx workflow.Context, branch string) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	for i, b := range state.branches {
		if b == branch {
			state.branches = append(state.branches[:i], state.branches[i+1:]...)
			break
		}
	}
}

func (state *RepoCtrlState) refresh_info(ctx workflow.Context) {}

// signal_branch sends a signal to the branch control state for the specified branch.
func (state *RepoCtrlState) signal_branch(ctx workflow.Context, branch string, signal shared.WorkflowSignal, payload any) {
	opts := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
	ctx = workflow.WithActivityOptions(ctx, opts)

	next := &RepoIOSignalBranchCtrlPayload{
		Repo:    state.repo,
		Branch:  branch,
		Signal:  signal,
		Payload: payload,
	}

	_ = state.do(
		ctx, "signal_branch_ctrl", state.activties.SignalBranch, next, nil,
		slog.String("signal", signal.String()),
	)
}

// increment is a helper function that increments the steps counter in the RepoIOBranchCtrlState.
func (state *RepoCtrlState) increment(ctx workflow.Context) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.counter++
}

// terminate is called to mark the RepoIOBranchCtrlState as inactive and cancel any associated timers.
// This function should be called when the branch control state is no longer needed, such as branch is being deleted or merged.
func (state *RepoCtrlState) terminate(ctx workflow.Context) {
	state.set_done(ctx)
}

// on_push is a channel handler that processes push events for the repository.
func (state *RepoCtrlState) on_push(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		push := &RepoIOSignalPushPayload{}
		rx.Receive(ctx, push)

		state.signal_branch(ctx, BranchNameFromRef(push.BranchRef), RepoIOSignalPush, push)
	}
}

// on_create_delete is a channel handler that processes create or delete events for the repository.
//
// TODO: handle create and delete events for tags.
func (state *RepoCtrlState) on_create_delete(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		create_delete := &RepoIOSignalCreateOrDeletePayload{}
		rx.Receive(ctx, create_delete)

		if create_delete.RefType == "branch" {
			state.signal_branch(ctx, create_delete.Ref, RepoIOSignalCreateOrDelete, create_delete)
		}
	}
}

// on_pr is a channel handler that processes pull request events for the repository.
func (state *RepoCtrlState) on_pr(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		pr := &RepoIOSignalPullRequestPayload{}
		rx.Receive(ctx, pr)

		state.signal_branch(ctx, pr.HeadBranch, RepoIOSignalPullRequest, pr)
	}
}

// do is a helper function that executes an activity within the context of the RepoIOBranchCtrlState.
// It logs the start and success of the activity, and increments the steps counter in the state.
// If the activity fails, the function logs the error and returns it.
//
// NOTE: This assumes that workflow.Context has been updated with activity options.
func (state *RepoCtrlState) do(ctx workflow.Context, action string, activity, payload, result any, keyvals ...any) error {
	// logger := NewRepoIOWorkflowLogger(ctx, state.repo, "branch_ctrl", state.branch, action)
	// logger.Info("init", keyvals...)
	state.log(ctx, "info", action, "init", keyvals...)

	if err := workflow.ExecuteActivity(ctx, activity, payload).Get(ctx, result); err != nil {
		keyvals = append(keyvals, "error", err)
		state.log(ctx, "error", action, "error", keyvals...)

		return err
	}

	state.increment(ctx)

	state.log(ctx, "info", action, "success", keyvals...)

	return nil
}

// log is a helper function that logs a message with the given level, action, and key-value pairs.
// It uses the NewRepoIOWorkflowLogger to create a logger scoped to the repository, branch control, and the given action.
// The log levels supported are "info", "warn", "error", and "debug". If an unknown level is provided, it defaults to "info".
func (state *RepoCtrlState) log(ctx workflow.Context, level, action, msg string, keyvals ...any) {
	logger := NewRepoIOWorkflowLogger(ctx, state.repo, "branch_ctrl", "", action)

	switch level {
	case "info":
		logger.Info(msg, keyvals...)
	case "warn":
		logger.Warn(msg, keyvals...)
	case "error":
		logger.Error(msg, keyvals...)
	case "debug":
		logger.Debug(msg, keyvals...)
	default:
		logger.Info(msg, keyvals...)
	}
}

func NewRepoCtrlState(ctx workflow.Context, repo *Repo) *RepoCtrlState {
	return &RepoCtrlState{
		activties: &RepoActivities{},
		repo:      repo,
		mutex:     workflow.NewMutex(ctx),
		counter:   0,
	}
}
