package core

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/timers"
)

// Workflow States.
type (
	RepoIOBranchCtrlState struct {
		repo        *Repo              // Repo is the db record of the repo
		branch      string             // Branch is the name of the branch
		created_at  time.Time          // created_at is the time when the branch was created
		last_commit *RepoIOCommit      // last_commit is the last commit on the branch
		pr          *RepoIOPullRequest // pr is the pull request associated with the branch
		interval    timers.Interval    // interval is the interval at which the branch is checked for staleness
		mutex       workflow.Mutex     // mutex is the mutex for the state
		active      bool               // active is the flag to indicate if the branch is active
		counter     int                // counter is the number of steps taken by the branch
	}
)

// set_created_at sets the created_at timestamp for the RepoIOBranchCtrlState.
// This method is thread-safe and locks the state's mutex before updating the created_at field.
func (state *RepoIOBranchCtrlState) set_created_at(ctx workflow.Context, t time.Time) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.created_at = t
}

// set_commit sets the last_commit field of the RepoIOBranchCtrlState.
// This method is thread-safe and locks the state's mutex before updating the last_commit field.
func (state *RepoIOBranchCtrlState) set_commit(ctx workflow.Context, commit *RepoIOCommit) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.last_commit = commit
}

// set_pr sets the pr field of the RepoIOBranchCtrlState.
// This method is thread-safe and locks the state's mutex before updating the pr field.
func (state *RepoIOBranchCtrlState) set_pr(ctx workflow.Context, pr *RepoIOPullRequest) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.pr = pr
}

func (state *RepoIOBranchCtrlState) set_done(ctx workflow.Context) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.active = false
}

func (state *RepoIOBranchCtrlState) is_active(ctx workflow.Context) bool {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	return state.active
}

func (state *RepoIOBranchCtrlState) has_pr(ctx workflow.Context) bool {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	return state.pr != nil
}

func (state *RepoIOBranchCtrlState) last_active(ctx workflow.Context) time.Time {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	if state.last_commit == nil {
		return state.created_at
	}

	return state.last_commit.Timestamp
}

// run_coroutine_state_check runs a background goroutine that periodically checks if the branch is stale and sends
// a warning message if it is.
func (state *RepoIOBranchCtrlState) run_coroutine_state_check(ctx workflow.Context) {
	data := state.get_repo_data(ctx)

	workflow.Go(ctx, func(ctx workflow.Context) {
		for state.is_active(ctx) {
			state.interval.Next(ctx)
			state.warn_stale(ctx, data)
		}
	})
}

func (state *RepoIOBranchCtrlState) get_repo_data(ctx workflow.Context) *RepoIORepoData {
	data := &RepoIORepoData{}
	io := Instance().RepoIO(state.repo.Provider)

	opts := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
	ctx = workflow.WithActivityOptions(ctx, opts)

	_ = state.run(ctx, "get_repo_data", io.GetRepoData, data, state.repo.CtrlID)

	return data
}

// calculate_complexity checks the complexity of the changes pushed on the current branch.
func (state *RepoIOBranchCtrlState) calculate_complexity(ctx workflow.Context, push *RepoIOSignalPushPayload) *RepoIOChanges {
	changes := &RepoIOChanges{}
	detect := &RepoIODetectChangesPayload{
		InstallationID: push.InstallationID,
		RepoName:       push.RepoName,
		RepoOwner:      push.RepoOwner,
		DefaultBranch:  state.repo.DefaultBranch,
		TargetBranch:   state.branch,
	}

	opts := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
	ctx = workflow.WithActivityOptions(ctx, opts)

	_ = state.run(ctx, "calculate_complexity", Instance().RepoIO(state.repo.Provider).DetectChanges, changes, detect)

	return changes
}

// warn_complexity sends a warning message to the linked message provider if the complexity of the changes exceeds the threshold.
// it sends a message to the user if the git user is linked to the quantm user and the linked quantm user also has connected the message
// provider, otherwise it sends a message to the linked channel of the repo.
func (state *RepoIOBranchCtrlState) warn_complexity(ctx workflow.Context, push *RepoIOSignalPushPayload, complexity *RepoIOChanges) {
	for_user := push.User != nil && push.User.IsMessageProviderLinked
	msg := NewNumberOfLinesExceedMessage(push, state.repo, state.branch, complexity, for_user)
	io := Instance().MessageIO(state.repo.MessageProvider)
	opts := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
	ctx = workflow.WithActivityOptions(ctx, opts)

	_ = state.run(ctx, "warn_complexity", io.SendNumberOfLinesExceedMessage, nil, msg)
}

func (state *RepoIOBranchCtrlState) warn_stale(ctx workflow.Context, data *RepoIORepoData) {
	msg := NewStaleBranchMessage(data, state.repo, state.branch)
	io := Instance().MessageIO(state.repo.MessageProvider)

	opts := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
	ctx = workflow.WithActivityOptions(ctx, opts)

	_ = state.run(ctx, "warn_stale", io.SendStaleBranchMessage, nil, msg)
}

// increment is a helper function that increments the steps counter in the RepoIOBranchCtrlState.
func (state *RepoIOBranchCtrlState) increment(ctx workflow.Context) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.counter++
}

func (state *RepoIOBranchCtrlState) create_session(ctx workflow.Context) workflow.Context {
	opts := &workflow.SessionOptions{ExecutionTimeout: 30 * time.Minute, CreationTimeout: 60 * time.Minute}
	ctx, _ = workflow.CreateSession(ctx, opts)

	return ctx
}

func (state *RepoIOBranchCtrlState) clone_at_commit(ctx workflow.Context, push *RepoIOSignalPushPayload) *RepoIOClonePayload {
	return nil
}

func (state *RepoIOBranchCtrlState) fetch_branch(ctx workflow.Context, cloned *RepoIOClonePayload) {}

func (state *RepoIOBranchCtrlState) rebase_at_commit(ctx workflow.Context, cloned *RepoIOClonePayload) error {
	return nil
}

func (state *RepoIOBranchCtrlState) warn_merge_conflict(ctx workflow.Context, push *RepoIOSignalPushPayload) {
}

func (state *RepoIOBranchCtrlState) remove_cloned(ctx workflow.Context, cloned *RepoIOClonePayload) {}

// run is a helper function that executes an activity within the context of the RepoIOBranchCtrlState.
// It logs the start and success of the activity, and increments the steps counter in the state.
// If the activity fails, the function logs the error and returns it.
//
// NOTE: This assumes that workflow.Context has been updated with activity options.
func (state *RepoIOBranchCtrlState) run(ctx workflow.Context, action string, activity any, result any, args ...any) error {
	logger := NewRepoIOWorkflowLogger(ctx, state.repo, "branch_ctrl", state.branch, action)
	logger.Info("init")

	if err := workflow.ExecuteActivity(ctx, activity, args...).Get(ctx, result); err != nil {
		logger.Error("failed", "error", err)
		return err
	}

	state.increment(ctx)

	logger.Info("success")

	return nil
}

// shutdown is called to mark the RepoIOBranchCtrlState as inactive and cancel any associated timers.
// This function should be called when the branch control state is no longer needed, such as branch is being deleted or merged.
func (state *RepoIOBranchCtrlState) shutdown(ctx workflow.Context) {
	state.set_done(ctx)
	state.interval.Cancel(ctx)
}

func NewBranchCtrlState(ctx workflow.Context, repo *Repo, branch string) *RepoIOBranchCtrlState {
	return &RepoIOBranchCtrlState{
		repo:       repo,
		branch:     branch,
		created_at: timers.Now(ctx),
		interval:   timers.NewInterval(ctx, repo.StaleDuration.Duration),
		mutex:      workflow.NewMutex(ctx),
		active:     true,
	}
}
