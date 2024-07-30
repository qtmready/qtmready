// Package core provides core functionality for repository operations and workflows.
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
	// RepoIOBranchCtrlState represents the state of a branch control workflow.
	RepoIOBranchCtrlState struct {
		*base_ctrl                       // base_ctrl is the embedded struct with common functionality for repo controls.
		active_branch string             // active_branch is the name of the branch associated with this control.
		created_at    time.Time          // created_at is the time when the branch was created.
		last_commit   *RepoIOCommit      // last_commit is the most recent commit on the branch.
		pr            *RepoIOPullRequest // pr is the pull request associated with the branch, if any.
		interval      timers.Interval    // interval is the stale check duration.
	}
)

// Event handlers

// on_push handles push events for the branch.
func (state *RepoIOBranchCtrlState) on_push(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		push := &RepoIOSignalPushPayload{}
		state.rx(ctx, rx, push) // Using base_ctrl.rx

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

// on_rebase handles rebase events for the branch.
func (state *RepoIOBranchCtrlState) on_rebase(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		push := &RepoIOSignalPushPayload{}
		state.rx(ctx, rx, push) // Using base_ctrl.rx

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

// on_pr handles pull request events for the branch.
func (state *RepoIOBranchCtrlState) on_pr(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		pr := &RepoIOSignalPullRequestPayload{}
		state.rx(ctx, rx, pr) // Using base_ctrl.rx

		switch pr.Action {
		case "opened":
			state.set_pr(ctx, &RepoIOPullRequest{Number: pr.Number, HeadBranch: pr.HeadBranch, BaseBranch: pr.BaseBranch})
		default:
			return
		}
	}
}

// on_create_delete handles branch creation and deletion events.
func (state *RepoIOBranchCtrlState) on_create_delete(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		payload := &RepoIOSignalCreateOrDeletePayload{}
		state.rx(ctx, rx, payload) // Using base_ctrl.rx

		if payload.IsCreated {
			state.set_created_at(ctx, timers.Now(ctx))
		} else {
			state.set_done(ctx)
		}
	}
}

// Core methods

func (state *RepoIOBranchCtrlState) branch() string {
	return state.active_branch
}

// set_created_at sets the creation time of the branch.
func (state *RepoIOBranchCtrlState) set_created_at(ctx workflow.Context, t time.Time) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.created_at = t
	state.increment(ctx, 1)
}

// set_commit updates the last commit of the branch.
func (state *RepoIOBranchCtrlState) set_commit(ctx workflow.Context, commit *RepoIOCommit) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()
	state.last_commit = commit

	state.increment(ctx, 1)
}

// set_pr sets the pull request associated with the branch.
func (state *RepoIOBranchCtrlState) set_pr(ctx workflow.Context, pr *RepoIOPullRequest) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()
	state.pr = pr

	state.increment(ctx, 1)
}

// has_pr checks if the branch has an associated pull request.
func (state *RepoIOBranchCtrlState) has_pr() bool {
	return state.pr != nil
}

// last_active returns the timestamp of the last activity on the branch.
func (state *RepoIOBranchCtrlState) last_active() time.Time {
	if state.last_commit == nil {
		return state.created_at
	}

	return state.last_commit.Timestamp
}

// check_stale periodically checks if the branch is stale and sends warnings.
func (state *RepoIOBranchCtrlState) check_stale(ctx workflow.Context) {
	workflow.Go(ctx, func(ctx workflow.Context) {
		for state.is_active() {
			state.interval.Next(ctx)
			state.warn_stale(ctx)
		}
	})
}

// Git operations

// create_session creates a new workflow session for Git operations.
func (state *RepoIOBranchCtrlState) create_session(ctx workflow.Context) workflow.Context {
	opts := &workflow.SessionOptions{ExecutionTimeout: 60 * time.Minute, CreationTimeout: 60 * time.Minute}
	ctx, _ = workflow.CreateSession(ctx, opts)

	return ctx
}

// clone_at_commit clones the repository at a specific commit.
func (state *RepoIOBranchCtrlState) clone_at_commit(ctx workflow.Context, push *RepoIOSignalPushPayload) *RepoIOClonePayload {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	cloned := &RepoIOClonePayload{Repo: state.repo, Push: push, Branch: state.branch()}
	_ = workflow.SideEffect(ctx, func(ctx workflow.Context) any { return "/tmp/" + uuid.New().String() }).Get(&cloned.Path)

	_ = state.do(ctx, "clone_at_commit", state.activities.CloneBranch, cloned, nil)

	return cloned
}

// fetch_default_branch fetches the default branch for the cloned repository.
func (state *RepoIOBranchCtrlState) fetch_default_branch(ctx workflow.Context, cloned *RepoIOClonePayload) {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	_ = state.do(ctx, "fetch_branch", state.activities.FetchBranch, cloned, nil)
}

// rebase_at_commit rebases the branch at a specific commit.
func (state *RepoIOBranchCtrlState) rebase_at_commit(ctx workflow.Context, cloned *RepoIOClonePayload) error {
	retry_policy := &temporal.RetryPolicy{NonRetryableErrorTypes: []string{"RepoIORebaseError"}}
	opts := workflow.ActivityOptions{StartToCloseTimeout: time.Minute, RetryPolicy: retry_policy}
	ctx = workflow.WithActivityOptions(ctx, opts)

	response := &RepoIORebaseAtCommitResponse{}

	if err := state.do(ctx, "rebase_at_commit", state.activities.RebaseAtCommit, cloned, response); err != nil {
		var apperr *temporal.ApplicationError
		if errors.As(err, &apperr) && apperr.Type() == "RepoIORebaseError" {
			return NewRepoIORebaseError(cloned.Push.After, "fetch the commit message here")
		}

		return nil
	}

	return nil
}

// push_branch pushes the rebased branch to the remote repository.
func (state *RepoIOBranchCtrlState) push_branch(ctx workflow.Context, cloned *RepoIOClonePayload) {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)
	payload := &RepoIOPushBranchPayload{Branch: cloned.Branch, Path: cloned.Path, Force: true}

	_ = state.do(ctx, "push_branch", state.activities.Push, payload, nil)
}

// remove_cloned removes the cloned repository from the local filesystem.
func (state *RepoIOBranchCtrlState) remove_cloned(ctx workflow.Context, cloned *RepoIOClonePayload) {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	_ = state.do(ctx, "remove_cloned", state.activities.RemoveClonedAtPath, cloned.Path, nil)
}

// Complexity and warning methods

// calculate_complexity calculates the complexity of changes in a push event.
func (state *RepoIOBranchCtrlState) calculate_complexity(ctx workflow.Context, push *RepoIOSignalPushPayload) *RepoIOChanges {
	changes := &RepoIOChanges{}
	detect := &RepoIODetectChangesPayload{
		InstallationID: push.InstallationID,
		RepoName:       push.RepoName,
		RepoOwner:      push.RepoOwner,
		DefaultBranch:  state.repo.DefaultBranch,
		TargetBranch:   state.branch(),
	}

	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	_ = state.do(ctx, "calculate_complexity", Instance().RepoIO(state.repo.Provider).DetectChanges, detect, changes)

	return changes
}

// warn_complexity sends a warning message if the complexity of changes exceeds the threshold.
func (state *RepoIOBranchCtrlState) warn_complexity(ctx workflow.Context, push *RepoIOSignalPushPayload, complexity *RepoIOChanges) {
	for_user := push.User != nil && push.User.IsMessageProviderLinked
	msg := NewNumberOfLinesExceedMessage(push, state.repo, state.branch(), complexity, for_user)
	io := Instance().MessageIO(state.repo.MessageProvider)
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	_ = state.do(ctx, "warn_complexity", io.SendNumberOfLinesExceedMessage, msg, nil)
}

// warn_stale sends a warning message if the branch is stale.
func (state *RepoIOBranchCtrlState) warn_stale(ctx workflow.Context) {
	msg := NewStaleBranchMessage(state.info, state.repo, state.branch())
	io := Instance().MessageIO(state.repo.MessageProvider)

	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	_ = state.do(ctx, "warn_stale", io.SendStaleBranchMessage, msg, nil)
}

// warn_conflict sends a warning message if there's a merge conflict during rebase.
func (state *RepoIOBranchCtrlState) warn_conflict(ctx workflow.Context, push *RepoIOSignalPushPayload) {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	for_user := push.User != nil && push.User.IsMessageProviderLinked
	msg := NewMergeConflictMessage(push, state.repo, state.branch(), for_user)
	io := Instance().MessageIO(state.repo.MessageProvider)

	_ = state.do(ctx, "warn_merge_conflict", io.SendMergeConflictsMessage, msg, nil)
}

// NewBranchCtrlState creates a new RepoIOBranchCtrlState instance.
func NewBranchCtrlState(ctx workflow.Context, repo *Repo, branch string) *RepoIOBranchCtrlState {
	return &RepoIOBranchCtrlState{
		base_ctrl:     NewBaseCtrl(ctx, "branch_ctrl", repo),
		active_branch: branch,
		created_at:    timers.Now(ctx),
		interval:      timers.NewInterval(ctx, repo.StaleDuration.Duration),
	}
}
