package code

import (
	"log/slog"
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/shared"
)

const (
	event_threshold = 4000
)

// DoFn represents the signature of the do function.
type (
	CallAsync func(workflow.Context)

	// BaseCtrl represents the base control structure for repository operations.
	// It provides common functionality for various repository control types.
	BaseCtrl struct {
		kind       string                   // kind identifies the type of control (e.g., "repo", "branch")
		activities *Activities              // activities holds the repository activities
		repo       *defs.Repo               // repo is a reference to the repository
		info       *defs.RepoIOProviderInfo // info stores provider-specific information
		branches   []string                 // branches is a list of branches in the repository
		mutex      workflow.Mutex           // mutex is used for thread-safe operations
		active     bool                     // active indicates if the control is active
		counter    int                      // counter counts the number of operations performed
	}
)

// is_active returns the active status of the control.
func (base *BaseCtrl) is_active() bool {
	return base.active
}

// needs_reset checks if the event count has reached the threshold for resetting.
func (base *BaseCtrl) needs_reset() bool {
	return base.counter >= event_threshold
}

// branch returns the branch name associated with this control.
func (base *BaseCtrl) branch(ctx workflow.Context) string {
	if branch, ok := ctx.Value("active_branch").(string); ok {
		return branch
	}

	return ""
}

// set_branch sets the active branch in the context.
func (base *BaseCtrl) set_branch(ctx workflow.Context, branch string) workflow.Context {
	return workflow.WithValue(ctx, "active_branch", branch)
}

// set_info sets the provider-specific information for the control.
func (base *BaseCtrl) set_info(ctx workflow.Context, info *defs.RepoIOProviderInfo) {
	_ = base.mutex.Lock(ctx)
	defer base.mutex.Unlock()

	base.info = info
}

// set_branches sets the list of branches associated with the control.
func (base *BaseCtrl) set_branches(ctx workflow.Context, branches []string) {
	_ = base.mutex.Lock(ctx)
	defer base.mutex.Unlock()

	base.branches = branches
}

// set_done marks the control as inactive.
func (base *BaseCtrl) set_done(ctx workflow.Context) {
	_ = base.mutex.Lock(ctx)
	defer base.mutex.Unlock()

	base.active = false
}

// terminate marks the control as done and logs the termination.
func (base *BaseCtrl) terminate(ctx workflow.Context) {
	base.set_done(ctx)
	base.log(ctx, "terminate").Info("state terminated")
}

// as_new continues the workflow as new with the given function and arguments.
func (base *BaseCtrl) as_new(ctx workflow.Context, msg string, fn any, args ...any) error {
	base.log(ctx, "as_new").Warn(msg)
	return workflow.NewContinueAsNewError(ctx, fn, args...)
}

// increment increases the operation counter by the specified number of steps.
func (base *BaseCtrl) increment(ctx workflow.Context, steps int) {
	_ = base.mutex.Lock(ctx)
	defer base.mutex.Unlock()

	base.counter += steps
}

// add_branch adds a new branch to the list of branches.
func (base *BaseCtrl) add_branch(ctx workflow.Context, branch string) {
	_ = base.mutex.Lock(ctx)
	defer base.mutex.Unlock()

	if branch != "" || branch != base.repo.DefaultBranch {
		base.branches = append(base.branches, branch)
	}
}

// remove_branch removes a branch from the list of branches.
func (base *BaseCtrl) remove_branch(ctx workflow.Context, branch string) {
	_ = base.mutex.Lock(ctx)
	defer base.mutex.Unlock()

	for i, b := range base.branches {
		if b == branch {
			base.branches = append(base.branches[:i], base.branches[i+1:]...)
			break
		}
	}
}

// signal_branch sends a signal to a specific branch.
func (base *BaseCtrl) signal_branch(ctx workflow.Context, branch string, signal shared.WorkflowSignal, payload any) {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	next := &defs.RepoIOSignalBranchCtrlPayload{
		Repo:    base.repo,
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

// rx receives a message from a channel and logs the event.
func (base *BaseCtrl) rx(ctx workflow.Context, channel workflow.ReceiveChannel, target any) {
	base.log(ctx, "rx").Info(channel.Name())

	channel.Receive(ctx, target)
}

// refresh_info updates the provider information for the repository.
func (base *BaseCtrl) refresh_info(ctx workflow.Context) {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	info := &defs.RepoIOProviderInfo{}
	io := kernel.Instance().RepoIO(base.repo.Provider)

	_ = base.do(ctx, "get_repo_info", io.GetProviderInfo, base.repo.CtrlID, info)
	base.set_info(ctx, info)
}

// refresh_branches updates the list of branches for the repository.
func (base *BaseCtrl) refresh_branches(ctx workflow.Context) {
	if base.info == nil {
		base.refresh_info(ctx)
	}

	io := kernel.Instance().RepoIO(base.repo.Provider)
	branches := []string{}

	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	_ = base.do(ctx, "refresh_branches", io.GetAllBranches, base.info, &branches)
	base.set_branches(ctx, branches)
}

// log creates a new logger for the current action.
func (base *BaseCtrl) log(ctx workflow.Context, action string) *RepoIOWorkflowLogger {
	return NewRepoIOWorkflowLogger(ctx, base.repo, base.kind, base.branch(ctx), action)
}

// do is helper is an activity executor. It logs the activity execution and increments the operation counter.
func (base *BaseCtrl) do(ctx workflow.Context, action string, activity, payload, result any, keyvals ...any) error {
	logger := base.log(ctx, action)
	logger.Info("init", keyvals...)

	if err := workflow.ExecuteActivity(ctx, activity, payload).Get(ctx, result); err != nil {
		logger.Warn("error", append(keyvals, "error", err)...)
		return err
	}

	logger.Info("success", keyvals...)

	base.increment(ctx, 10)

	return nil
}

// TODO: @asgr - need tp refine.
func (base *BaseCtrl) do_child(ctx workflow.Context, action, w_id string, fn, payload any, keyvals ...any) error {
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

	base.increment(ctx, 3)

	return nil
}

// call_async executes an activity asynchronously and returns a Future.
// If a WaitGroup is provided, it will be decremented when the operation completes.
func (base *BaseCtrl) call_async(ctx workflow.Context, action string, fn CallAsync, wg workflow.WaitGroup) workflow.Future {
	logger := base.log(ctx, action)

	future, setable := workflow.NewFuture(ctx)
	workflow.Go(ctx, func(ctx workflow.Context) {
		logger.Info("calling async ...")

		if wg != nil {
			defer wg.Done()
		}

		fn(ctx)
		setable.Set(nil, nil)
	})

	return future
}

// NewBaseCtrl creates a new base control instance and refreshes repository information and branches.
func NewBaseCtrl(ctx workflow.Context, kind string, repo *defs.Repo) *BaseCtrl {
	base := &BaseCtrl{
		kind:       kind,
		activities: &Activities{},
		info:       &defs.RepoIOProviderInfo{},
		repo:       repo,
		mutex:      workflow.NewMutex(ctx),
		active:     true,
		counter:    0,
	}

	base.refresh_info(ctx)
	base.refresh_branches(ctx)

	return base
}
