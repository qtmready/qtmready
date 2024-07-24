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

// calculate_complexity checks the complexity of the changes pushed on the current branch.
func (state *RepoIOBranchCtrlState) calculate_complexity(ctx workflow.Context, signal *RepoIOSignalPushPayload) *RepoIOChanges {
	changes := &RepoIOChanges{}
	detect := &RepoIODetectChangesPayload{
		InstallationID: signal.InstallationID,
		RepoName:       signal.RepoName,
		RepoOwner:      signal.RepoOwner,
		DefaultBranch:  state.repo.DefaultBranch,
		TargetBranch:   state.branch,
	}

	_ = state._execute(ctx, "detect_changes", Instance().RepoIO(state.repo.Provider).DetectChanges, changes, detect)

	return changes
}

// warn_complexity sends a warning message to the linked message provider if the complexity of the changes exceeds the threshold.
// it sends a message to the user if the git user is linked to the quantm user and the linked quantm user also has connected the message
// provider, otherwise it sends a message to the linked channel of the repo.
func (state *RepoIOBranchCtrlState) warn_complexity(ctx workflow.Context, signal *RepoIOSignalPushPayload, complexity *RepoIOChanges) {
	for_user := signal.User != nil && signal.User.IsMessageProviderLinked
	msg := NewNumberOfLinesExceedMessage(signal, state.repo, state.branch, complexity, for_user)
	io := Instance().MessageIO(state.repo.MessageProvider)

	_ = state._execute(ctx, "send_complexity_warning", io.SendNumberOfLinesExceedMessage, nil, msg)
}

// increment is a helper function that increments the steps counter in the RepoIOBranchCtrlState.
func (state *RepoIOBranchCtrlState) increment(ctx workflow.Context) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.counter++
}

// _execute is a helper function that executes an activity within the context of the RepoIOBranchCtrlState.
// It logs the start and success of the activity, and increments the steps counter in the state.
//
// The function takes the following parameters:
// - ctx: the workflow context
// - action: a string describing the action being performed
// - activity: the activity function to _execute
// - result: a pointer to a variable to receive the result of the activity
// - args: any additional arguments to pass to the activity function
//
// If the activity fails, the function logs the error and returns it.
func (state *RepoIOBranchCtrlState) _execute(ctx workflow.Context, action string, activity any, result any, args ...any) error {
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
