package code

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/comm"
	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/core/timers"
	"go.breu.io/quantm/internal/shared"
)

type (
	// RepoIOBranchCtrlState represents the state of a branch control workflow.
	RepoIOBranchCtrlState struct {
		*BaseCtrl                           // base_ctrl is the embedded struct with common functionality for repo controls.
		created_at  time.Time               // created_at is the time when the branch was created.
		last_commit *defs.RepoIOCommit      // last_commit is the most recent commit on the branch.
		pr          *defs.RepoIOPullRequest // pr is the pull request associated with the branch, if any.
		interval    timers.Interval         // interval is the stale check duration.
	}
)

// Event handlers

// on_push handles push events for the branch.
func (state *RepoIOBranchCtrlState) on_push(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		push := &defs.RepoIOSignalPushPayload{}
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
		push := &defs.RepoIOSignalPushPayload{}
		state.rx(ctx, rx, push) // Using base_ctrl.rx

		session := state.create_session(ctx)
		defer state.finish_session(session)

		cloned := state.clone_at_commit(session, push)
		if cloned == nil {
			return
		}

		state.fetch_default_branch(session, cloned)

		if err := state.rebase_at_commit(session, cloned); err != nil {
			state.warn_conflict(session, push)
			state.remove_cloned(session, cloned)

			return
		}

		state.push_branch(session, cloned)
		state.remove_cloned(session, cloned)
	}
}

// on_pr handles pull request events for the branch.
func (state *RepoIOBranchCtrlState) on_pr(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		pr := &defs.RepoIOSignalPullRequestPayload{}
		state.rx(ctx, rx, pr) // Using base_ctrl.rx

		switch pr.Action {
		case "opened":
			state.set_pr(ctx, &defs.RepoIOPullRequest{Number: pr.Number, HeadBranch: pr.HeadBranch, BaseBranch: pr.BaseBranch})
		case "closed":
			// when the pull request action is closed set it to nil.
			state.set_pr(ctx, nil)
		default:
			return
		}
	}
}

// on_create_delete handles branch creation and deletion events.
func (state *RepoIOBranchCtrlState) on_create_delete(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		payload := &defs.RepoIOSignalCreateOrDeletePayload{}
		state.rx(ctx, rx, payload) // Using base_ctrl.rx

		if payload.IsCreated {
			state.set_created_at(ctx, timers.Now(ctx))
		} else {
			state.set_done(ctx)
		}
	}
}

// Core methods

// set_created_at sets the creation time of the branch.
func (state *RepoIOBranchCtrlState) set_created_at(ctx workflow.Context, t time.Time) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.created_at = t
	state.increment(ctx, 1)
}

// set_commit updates the last commit of the branch.
func (state *RepoIOBranchCtrlState) set_commit(ctx workflow.Context, commit *defs.RepoIOCommit) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()
	state.last_commit = commit

	state.increment(ctx, 1)
}

// set_pr sets the pull request associated with the branch.
func (state *RepoIOBranchCtrlState) set_pr(ctx workflow.Context, pr *defs.RepoIOPullRequest) {
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
	state.log(ctx, "session").Info("init")

	opts := &workflow.SessionOptions{ExecutionTimeout: 60 * time.Minute, CreationTimeout: 60 * time.Minute}
	session, _ := workflow.CreateSession(ctx, opts)

	return session
}

func (state *RepoIOBranchCtrlState) finish_session(ctx workflow.Context) {
	workflow.CompleteSession(ctx)
	state.log(ctx, "session").Info("completed")
}

// clone_at_commit clones the repository at a specific commit.
func (state *RepoIOBranchCtrlState) clone_at_commit(ctx workflow.Context, push *defs.RepoIOSignalPushPayload) *defs.RepoIOClonePayload {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	cloned := &defs.RepoIOClonePayload{Repo: state.repo, Push: push, Branch: state.branch(ctx)}
	_ = workflow.SideEffect(ctx, func(ctx workflow.Context) any { return "/tmp/" + uuid.New().String() }).Get(&cloned.Path)

	_ = state.do(ctx, "clone_at_commit", state.activities.CloneBranch, cloned, nil)

	return cloned
}

// fetch_default_branch fetches the default branch for the cloned repository.
func (state *RepoIOBranchCtrlState) fetch_default_branch(ctx workflow.Context, cloned *defs.RepoIOClonePayload) {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	_ = state.do(ctx, "fetch_branch", state.activities.FetchBranch, cloned, nil)
}

// rebase_at_commit rebases the branch at a specific commit.
func (state *RepoIOBranchCtrlState) rebase_at_commit(ctx workflow.Context, cloned *defs.RepoIOClonePayload) error {
	retry_policy := &temporal.RetryPolicy{NonRetryableErrorTypes: []string{"RepoIORebaseError"}}
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second, RetryPolicy: retry_policy}
	ctx = workflow.WithActivityOptions(ctx, opts)

	response := &defs.RepoIORebaseAtCommitResponse{}

	if err := state.do(ctx, "rebase_at_commit", state.activities.RebaseAtCommit, cloned, response); err != nil {
		var apperr *temporal.ApplicationError
		if errors.As(err, &apperr) && apperr.Type() == "RepoIORebaseError" {
			return NewRebaseError(cloned.Push.After, "fetch the commit message here") // TODO: fill the right info
		}

		return nil
	}

	return nil
}

// push_branch pushes the rebased branch to the remote repository.
func (state *RepoIOBranchCtrlState) push_branch(ctx workflow.Context, cloned *defs.RepoIOClonePayload) {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)
	payload := &defs.RepoIOPushBranchPayload{Branch: cloned.Branch, Path: cloned.Path, Force: true}

	_ = state.do(ctx, "push_branch", state.activities.Push, payload, nil)
}

// remove_cloned removes the cloned repository from the local filesystem.
func (state *RepoIOBranchCtrlState) remove_cloned(ctx workflow.Context, cloned *defs.RepoIOClonePayload) {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	_ = state.do(ctx, "remove_cloned", state.activities.RemoveClonedAtPath, cloned.Path, nil)
}

// Complexity and warning methods

// calculate_complexity calculates the complexity of changes in a push event.
func (state *RepoIOBranchCtrlState) calculate_complexity(ctx workflow.Context, push *defs.RepoIOSignalPushPayload) *defs.RepoIOChanges {
	changes := &defs.RepoIOChanges{}
	detect := &defs.RepoIODetectChangesPayload{
		InstallationID: push.InstallationID,
		RepoName:       push.RepoName,
		RepoOwner:      push.RepoOwner,
		DefaultBranch:  state.repo.DefaultBranch,
		TargetBranch:   state.branch(ctx),
	}

	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	_ = state.do(ctx, "calculate_complexity", kernel.Instance().RepoIO(state.repo.Provider).DetectChanges, detect, changes)

	return changes
}

// warn_complexity sends a warning message if the complexity of changes exceeds the threshold.
func (state *RepoIOBranchCtrlState) warn_complexity(
	ctx workflow.Context, push *defs.RepoIOSignalPushPayload, complexity *defs.RepoIOChanges,
) {
	for_user := push.User != nil && push.User.IsMessageProviderLinked
	msg := comm.NewNumberOfLinesExceedMessage(push, state.repo, state.branch(ctx), complexity, for_user)
	io := kernel.Instance().MessageIO(state.repo.MessageProvider)
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	_ = state.do(ctx, "warn_complexity", io.SendNumberOfLinesExceedMessage, msg, nil)
}

// warn_stale sends a warning message if the branch is stale.
func (state *RepoIOBranchCtrlState) warn_stale(ctx workflow.Context) {
	msg := comm.NewStaleBranchMessage(state.info, state.repo, state.branch(ctx))
	io := kernel.Instance().MessageIO(state.repo.MessageProvider)

	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	_ = state.do(ctx, "warn_stale", io.SendStaleBranchMessage, msg, nil)
}

// warn_conflict sends a warning message if there's a merge conflict during rebase.
func (state *RepoIOBranchCtrlState) warn_conflict(ctx workflow.Context, push *defs.RepoIOSignalPushPayload) {
	opts := workflow.ActivityOptions{StartToCloseTimeout: 60 * time.Second}
	ctx = workflow.WithActivityOptions(ctx, opts)

	for_user := push.User != nil && push.User.IsMessageProviderLinked
	msg := comm.NewMergeConflictMessage(push, state.repo, state.branch(ctx), for_user)
	io := kernel.Instance().MessageIO(state.repo.MessageProvider)

	_ = state.do(ctx, "warn_merge_conflict", io.SendMergeConflictsMessage, msg, nil)
}

// NewBranchCtrlState creates a new RepoIOBranchCtrlState instance.
func NewBranchCtrlState(ctx workflow.Context, repo *defs.Repo, branch string) (workflow.Context, *RepoIOBranchCtrlState) {
	base := &RepoIOBranchCtrlState{
		BaseCtrl:   NewBaseCtrl(ctx, "branch_ctrl", repo),
		created_at: timers.Now(ctx),
		interval:   timers.NewInterval(ctx, repo.StaleDuration.Duration),
	}

	return base.set_branch(ctx, branch), base
}
