package core

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/timers"
	"go.breu.io/quantm/internal/shared"
)

type (
	// BranchCtrlState defines the state for RepoWorkflows.BranchCtrl.
	//
	// NOTE: This state is local to the workflow and all the members are private. It cannot be passed to child workflows.
	BranchCtrlState struct {
		activties   *RepoActivities    // activities is the activities for the branch control state
		repo        *Repo              // Repo is the db record of the repo
		branch      string             // Branch is the name of the branch
		created_at  time.Time          // created_at is the time when the branch was created
		last_commit *RepoIOCommit      // last_commit is the last commit on the branch
		pr          *RepoIOPullRequest // pr is the pull request associated with the branch
		interval    timers.Interval    // interval is the interval at which the branch is checked for staleness
		mutex       workflow.Mutex     // mutex to provide thread safe access to the state
		active      bool               // active is the flag to indicate if the branch is active
		counter     int                // counter is the number of steps taken by the branch
	}
)

func (state *BranchCtrlState) on_push(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		push := &RepoIOSignalPushPayload{}
		rx.Receive(ctx, push)

		latest := push.Commits.Latest()
		if latest != nil {
			state.set_commit(ctx, latest)
		}

		complexity := state.calculate_complexity(ctx, push)
		if complexity.Delta > state.repo.Threshold {
			state.warn_complexity(ctx, push, complexity)
		}

		state.interval.Restart(ctx)
	}
}

// on_rebase is a shared.ChannelHandler that is called when a branch needs to be rebased. It handles the logic for
// cloning the repository, fetching the default branch, rebasing the branch at the latest commit, and pushing the rebased
// branch back to the repository.
func (state *BranchCtrlState) on_rebase(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		push := &RepoIOSignalPushPayload{}

		rx.Receive(ctx, push)

		session := state.create_session(ctx)
		defer workflow.CompleteSession(session)

		cloned := state.clone_at_commit(session, push)
		if cloned == nil {
			return
		}

		state.fetch_default_branch(session, cloned)

		if err := state.rebase_at_commit(session, cloned); err != nil {
			state.warn_conflict(ctx, push)
		}

		state.push_branch(session, cloned)
		state.remove_cloned(ctx, cloned)
	}
}

func (state *BranchCtrlState) on_pr(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		pr := &RepoIOSignalPullRequestPayload{}
		rx.Receive(ctx, pr)

		switch pr.Action {
		case "opened": //nolint
			state.set_pr(ctx, &RepoIOPullRequest{Number: pr.Number, HeadBranch: pr.HeadBranch, BaseBranch: pr.BaseBranch})
		default:
			return
		}
	}
}

// on_create_delete is a shared.ChannelHandler that is called when a branch is created or deleted. It handles the logic for
// updating the state of the branch control when a create or delete event is received.
//
// If the payload indicates the branch was created, the function sets the created timestamp in the state.
// If the payload indicates the branch was deleted, the function terminates the state.
func (state *BranchCtrlState) on_create_delete(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		payload := &RepoIOSignalCreateOrDeletePayload{}
		rx.Receive(ctx, payload)

		if payload.IsCreated {
			state.set_created_at(ctx, timers.Now(ctx))
		} else {
			state.set_done(ctx)
		}
	}
}

// set_created_at sets the created_at timestamp for the RepoIOBranchCtrlState.
// This method is thread-safe and locks the state's mutex before updating the created_at field.
func (state *BranchCtrlState) set_created_at(ctx workflow.Context, t time.Time) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.created_at = t
}

// set_commit sets the last_commit field of the RepoIOBranchCtrlState.
// This method is thread-safe and locks the state's mutex before updating the last_commit field.
func (state *BranchCtrlState) set_commit(ctx workflow.Context, commit *RepoIOCommit) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.last_commit = commit
}

// set_pr sets the pr field of the RepoIOBranchCtrlState.
// This method is thread-safe and locks the state's mutex before updating the pr field.
func (state *BranchCtrlState) set_pr(ctx workflow.Context, pr *RepoIOPullRequest) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.pr = pr
}

func (state *BranchCtrlState) set_done(ctx workflow.Context) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.active = false
}

func (state *BranchCtrlState) is_active() bool {
	return state.active
}

func (state *BranchCtrlState) has_pr() bool {
	return state.pr != nil
}

func (state *BranchCtrlState) last_active() time.Time {
	if state.last_commit == nil {
		return state.created_at
	}

	return state.last_commit.Timestamp
}

// check_stale runs a background goroutine that periodically checks if the branch is stale and sends
// a warning message if it is.
func (state *BranchCtrlState) check_stale(ctx workflow.Context) {
	data := state.get_repo_data(ctx)

	workflow.Go(ctx, func(ctx workflow.Context) {
		for state.is_active() {
			state.interval.Next(ctx)
			state.warn_stale(ctx, data)
		}
	})
}

// get_repo_data returns the core repo by provider repo id.
func (state *BranchCtrlState) get_repo_data(ctx workflow.Context) *RepoIOProviderInfo {
	info := &RepoIOProviderInfo{}
	io := Instance().RepoIO(state.repo.Provider)

	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	_ = state.do(ctx, "get_repo_data", io.GetProviderInfo, state.repo.CtrlID, info)

	return info
}

// calculate_complexity checks the complexity of the changes pushed on the current branch.
func (state *BranchCtrlState) calculate_complexity(ctx workflow.Context, push *RepoIOSignalPushPayload) *RepoIOChanges {
	changes := &RepoIOChanges{}
	detect := &RepoIODetectChangesPayload{
		InstallationID: push.InstallationID,
		RepoName:       push.RepoName,
		RepoOwner:      push.RepoOwner,
		DefaultBranch:  state.repo.DefaultBranch,
		TargetBranch:   state.branch,
	}

	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	_ = state.do(ctx, "calculate_complexity", Instance().RepoIO(state.repo.Provider).DetectChanges, detect, changes)

	return changes
}

func (state *BranchCtrlState) create_session(ctx workflow.Context) workflow.Context {
	opts := &workflow.SessionOptions{ExecutionTimeout: 60 * time.Minute, CreationTimeout: 60 * time.Minute}
	ctx, _ = workflow.CreateSession(ctx, opts)

	return ctx
}

// clone_at_commit clones the branch at the specified commit with depth = 0. The cloned repository is stored in the cloned.Path field.
func (state *BranchCtrlState) clone_at_commit(ctx workflow.Context, push *RepoIOSignalPushPayload) *RepoIOClonePayload {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	cloned := &RepoIOClonePayload{Repo: state.repo, Push: push, Branch: state.branch}
	_ = workflow.SideEffect(ctx, func(ctx workflow.Context) any { return "/tmp/" + uuid.New().String() }).Get(&cloned.Path)

	_ = state.do(ctx, "clone_at_commit", state.activties.CloneBranch, cloned, nil)

	return cloned
}

// fetch_default_branch fetches the default branch for the cloned repository.
func (state *BranchCtrlState) fetch_default_branch(ctx workflow.Context, cloned *RepoIOClonePayload) {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	_ = state.do(ctx, "fetch_branch", state.activties.FetchBranch, cloned, nil)
}

// rebase_at_commit rebases the branch at the specified commit.
func (state *BranchCtrlState) rebase_at_commit(ctx workflow.Context, cloned *RepoIOClonePayload) error {
	retry_policy := &temporal.RetryPolicy{NonRetryableErrorTypes: []string{"RepoIORebaseError"}}
	opts := workflow.ActivityOptions{StartToCloseTimeout: time.Minute, RetryPolicy: retry_policy}
	ctx = workflow.WithActivityOptions(ctx, opts)

	response := &RepoIORebaseAtCommitResponse{}

	if err := state.do(ctx, "rebase_at_commit", state.activties.RebaseAtCommit, cloned, response); err != nil {
		var apperr *temporal.ApplicationError
		if errors.As(err, &apperr) && apperr.Type() == "RepoIORebaseError" {
			return NewRepoIORebaseError(cloned.Push.After, "fetch the commit message here")
		}

		return nil
	}

	return nil
}

// push_branch pushes the branch from the cloned repository to the remote repository.
func (state *BranchCtrlState) push_branch(ctx workflow.Context, cloned *RepoIOClonePayload) {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)
	payload := &RepoIOPushBranchPayload{Branch: cloned.Branch, Path: cloned.Path, Force: true}

	_ = state.do(ctx, "push_branch", state.activties.Push, payload, nil)
}

// remove_cloned removes the cloned repository at the specified path.
func (state *BranchCtrlState) remove_cloned(ctx workflow.Context, cloned *RepoIOClonePayload) {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	_ = state.do(ctx, "remove_cloned", state.activties.RemoveClonedAtPath, cloned.Path, nil)
}

// warn_complexity sends a warning message to the linked message provider if the complexity of the changes exceeds the threshold.
// it sends a message to the user if the git user is linked to the quantm user and the linked quantm user also has connected the message
// provider, otherwise it sends a message to the linked channel of the repo.
func (state *BranchCtrlState) warn_complexity(ctx workflow.Context, push *RepoIOSignalPushPayload, complexity *RepoIOChanges) {
	for_user := push.User != nil && push.User.IsMessageProviderLinked
	msg := NewNumberOfLinesExceedMessage(push, state.repo, state.branch, complexity, for_user)
	io := Instance().MessageIO(state.repo.MessageProvider)
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	_ = state.do(ctx, "warn_complexity", io.SendNumberOfLinesExceedMessage, msg, nil)
}

// warn_stale sends a warning message to the linked message provider if the branch is stale.
// It sends a message to the user if the git user is linked to the quantm user and the linked quantm user also has connected the message
// provider, otherwise it sends a message to the linked channel of the repo.
func (state *BranchCtrlState) warn_stale(ctx workflow.Context, data *RepoIOProviderInfo) {
	msg := NewStaleBranchMessage(data, state.repo, state.branch)
	io := Instance().MessageIO(state.repo.MessageProvider)

	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	_ = state.do(ctx, "warn_stale", io.SendStaleBranchMessage, msg, nil)
}

// warn_conflict sends a warning message to the linked message provider if there is a merge conflict.
// It sends a message to the user if the git user is linked to the quantm user and the linked quantm user also has connected the message
// provider, otherwise it sends a message to the linked channel of the repo.
func (state *BranchCtrlState) warn_conflict(ctx workflow.Context, push *RepoIOSignalPushPayload) {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	for_user := push.User != nil && push.User.IsMessageProviderLinked
	msg := NewMergeConflictMessage(push, state.repo, state.branch, for_user)
	io := Instance().MessageIO(state.repo.MessageProvider)

	_ = state.do(ctx, "warn_merge_conflict", io.SendMergeConflictsMessage, msg, nil)
}

// increment is a helper function that increments the steps counter in the RepoIOBranchCtrlState.
func (state *BranchCtrlState) increment(ctx workflow.Context) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.counter++
}

// terminate is called to mark the RepoIOBranchCtrlState as inactive and cancel any associated timers.
// This function should be called when the branch control state is no longer needed, such as branch is being deleted or merged.
func (state *BranchCtrlState) terminate(ctx workflow.Context) {
	state.set_done(ctx)
	state.interval.Cancel(ctx)
}

// do is a helper function that executes an activity within the context of the RepoIOBranchCtrlState.
// It logs the start and success of the activity, and increments the steps counter in the state.
// If the activity fails, the function logs the error and returns it.
//
// NOTE: This assumes that workflow.Context has been updated with activity options.
func (state *BranchCtrlState) do(ctx workflow.Context, action string, activity, payload, result any, keyvals ...any) error {
	return _do(ctx, state.repo, state.branch, "branch_ctrl", action, activity, payload, result, keyvals...)
}

func NewBranchCtrlState(ctx workflow.Context, repo *Repo, branch string) *BranchCtrlState {
	return &BranchCtrlState{
		activties:  &RepoActivities{},
		repo:       repo,
		branch:     branch,
		created_at: timers.Now(ctx),
		interval:   timers.NewInterval(ctx, repo.StaleDuration.Duration),
		mutex:      workflow.NewMutex(ctx),
		active:     true,
	}
}
