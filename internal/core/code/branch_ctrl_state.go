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
	"errors"
	"time"

	"github.com/google/uuid"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core/comm"
	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/core/timers"
	"go.breu.io/quantm/internal/shared"
)

type (
	// RepoIOBranchCtrlState represents the state of a branch control workflow.
	RepoIOBranchCtrlState struct {
		*BaseCtrl                     // base_ctrl is the embedded struct with common functionality for repo controls.
		created_at  time.Time         // created_at is the time when the branch was created.
		last_commit *defs.Commit      // last_commit is the most recent commit on the branch.
		pr          *defs.PullRequest // pr is the pull request associated with the branch, if any.
		interval    timers.Interval   // interval is the stale check duration.
		author      *auth.TeamUser    // owner is owner of the pull request.
	}
)

// Event handlers

// on_push handles push events for the branch.
func (state *RepoIOBranchCtrlState) on_push(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		event := &defs.Event[defs.Push, defs.RepoProvider]{}
		state.rx(ctx, rx, event)

		latest := event.Payload.Commits.Latest()
		if latest != nil {
			state.set_commit(ctx, latest)
		}

		complexity := state.calculate_complexity(ctx)
		if complexity.Delta > state.repo.Threshold {
			// state.warn_complexity(ctx, event, complexity) // TODO: handle this.
		}

		state.interval.Restart(ctx)
	}
}

// on_rebase handles rebase events for the branch.
func (state *RepoIOBranchCtrlState) on_rebase(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		event := &defs.Event[defs.Push, defs.RepoProvider]{}
		state.rx(ctx, rx, event)

		ctx = shared.WithDefaultActivityContext(ctx)

		session := state.create_session(ctx)
		defer state.finish_session(session)

		cloned := state.clone_at_commit(session, &event.Payload)
		if cloned == nil {
			return
		}

		state.fetch_default_branch(session, cloned)

		if err := state.rebase_at_commit(session, cloned); err != nil {
			state.warn_conflict(session, event)
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
		event := &defs.Event[defs.PullRequest, defs.RepoProvider]{}
		state.rx(ctx, rx, event)

		switch event.Context.Action { // nolint
		case defs.EventActionCreated, defs.EventActionReopened: // Or defs.EventActionOpened, if more appropriate
			state.set_pr(ctx, &event.Payload)
			state.refresh_author(ctx, event.Payload.AuthorID)
		case defs.EventActionClosed:
			state.set_pr(ctx, nil)
			state.set_author(ctx, nil)
		default:
			return
		}
	}
}

// TODO - refine the logic.
// on_label handles pull request label events.
func (state *RepoIOBranchCtrlState) on_label(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		label := &defs.RepoIOSignalPullRequestPayload{}
		state.rx(ctx, rx, label)

		switch label.Action {
		case "labeled":
			// TODO - need to finalize the Lable Names and logic.
			switch *label.LabelName {
			case "qmerge":
				state.signal_queue(ctx, label.HeadBranch, defs.RepoIOSignalQueueAdd, label)
			case "priority-qmerge":
				state.signal_queue(ctx, label.HeadBranch, defs.RepoIOSignalQueueAddPriority, label)
			case "remove":
				state.signal_queue(ctx, label.HeadBranch, defs.RepoIOSignalQueueRemove, label)
			default:
				return
			}
		case "unlabeled":
			state.signal_queue(ctx, label.HeadBranch, defs.RepoIOSignalQueueRemove, label)
		default:
			return
		}
	}
}

// on_create_delete handles branch creation and deletion events.
func (state *RepoIOBranchCtrlState) on_create_delete(ctx workflow.Context) shared.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		event := &defs.Event[defs.BranchOrTag, defs.RepoProvider]{}
		state.rx(ctx, rx, event)

		if event.Context.Action == defs.EventActionCreated {
			state.set_created_at(ctx, timers.Now(ctx))
		} else if event.Context.Action == defs.EventActionDeleted {
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
}

// set_commit updates the last commit of the branch.
func (state *RepoIOBranchCtrlState) set_commit(ctx workflow.Context, commit *defs.Commit) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()
	state.last_commit = commit
}

// set_pr sets the pull request associated with the branch.
func (state *RepoIOBranchCtrlState) set_pr(ctx workflow.Context, pr *defs.PullRequest) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()
	state.pr = pr
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

func (state *RepoIOBranchCtrlState) has_author() bool {
	return state.author != nil
}

func (state *RepoIOBranchCtrlState) set_author(ctx workflow.Context, owner *auth.TeamUser) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.author = owner
}

func (state *RepoIOBranchCtrlState) refresh_author(ctx workflow.Context, id shared.Int64) {
	ctx = shared.WithDefaultActivityContext(ctx)

	user := &auth.TeamUser{}

	_ = state.do(ctx, "refresh_author", state.activities.GetByLogin, id.String(), user)
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
func (state *RepoIOBranchCtrlState) clone_at_commit(ctx workflow.Context, push *defs.Push) *defs.RepoIOClonePayload {
	ctx = shared.WithDefaultActivityContext(ctx)

	cloned := &defs.RepoIOClonePayload{Repo: state.repo, Push: push, Info: state.info, Branch: state.branch(ctx)}
	_ = workflow.SideEffect(ctx, func(ctx workflow.Context) any { return "/tmp/" + uuid.New().String() }).Get(&cloned.Path)

	_ = state.do(ctx, "clone_at_commit", state.activities.CloneBranch, cloned, nil)

	return cloned
}

// fetch_default_branch fetches the default branch for the cloned repository.
func (state *RepoIOBranchCtrlState) fetch_default_branch(ctx workflow.Context, cloned *defs.RepoIOClonePayload) {
	ctx = shared.WithDefaultActivityContext(ctx)

	_ = state.do(ctx, "fetch_branch", state.activities.FetchBranch, cloned, nil)
}

// rebase_at_commit rebases the branch at a specific commit.
func (state *RepoIOBranchCtrlState) rebase_at_commit(ctx workflow.Context, cloned *defs.RepoIOClonePayload) error {
	ctx = shared.WithIgnoredErrorsContext(ctx, "RebaseError")

	response := &defs.RepoIORebaseAtCommitResponse{}

	if err := state.do(ctx, "rebase_at_commit", state.activities.RebaseAtCommit, cloned, response); err != nil {
		var apperr *temporal.ApplicationError
		if errors.As(err, &apperr) && apperr.Type() == "RebaseError" {
			return NewRebaseError(cloned.Push.After, "fetch the commit message here") // TODO: fill the right info
		}

		return nil
	}

	return nil
}

// push_branch pushes the rebased branch to the remote repository.
func (state *RepoIOBranchCtrlState) push_branch(ctx workflow.Context, cloned *defs.RepoIOClonePayload) {
	ctx = shared.WithDefaultActivityContext(ctx)
	payload := &defs.RepoIOPushBranchPayload{Branch: cloned.Branch, Path: cloned.Path, Force: true}

	_ = state.do(ctx, "push_branch", state.activities.Push, payload, nil)
}

// remove_cloned removes the cloned repository from the local filesystem.
func (state *RepoIOBranchCtrlState) remove_cloned(ctx workflow.Context, cloned *defs.RepoIOClonePayload) {
	ctx = shared.WithDefaultActivityContext(ctx)

	_ = state.do(ctx, "remove_cloned", state.activities.RemoveClonedAtPath, cloned.Path, nil)
}

// Complexity and warning methods

// calculate_complexity calculates the complexity of changes in a push event.
//
// TODO: we should compare the default branch head commit and the push's latest commit.
func (state *RepoIOBranchCtrlState) calculate_complexity(ctx workflow.Context) *defs.RepoIOChanges {
	changes := &defs.RepoIOChanges{}
	detect := &defs.RepoIODetectChangesPayload{
		InstallationID: state.info.InstallationID,
		RepoName:       state.info.RepoName,
		RepoOwner:      state.info.RepoOwner,
		DefaultBranch:  state.repo.DefaultBranch,
		TargetBranch:   state.branch(ctx),
	}

	ctx = shared.WithDefaultActivityContext(ctx)

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
	ctx = shared.WithDefaultActivityContext(ctx)

	_ = state.do(ctx, "warn_complexity", io.SendNumberOfLinesExceedMessage, msg, nil)
}

// warn_stale sends a warning message if the branch is stale.
func (state *RepoIOBranchCtrlState) warn_stale(ctx workflow.Context) {
	msg := comm.NewStaleBranchMessage(state.info, state.repo, state.branch(ctx))
	io := kernel.Instance().MessageIO(state.repo.MessageProvider)

	ctx = shared.WithDefaultActivityContext(ctx)

	_ = state.do(ctx, "warn_stale", io.SendStaleBranchMessage, msg, nil)
}

// warn_conflict sends a warning message if there's a merge conflict during rebase.
func (state *RepoIOBranchCtrlState) warn_conflict(ctx workflow.Context, event *defs.Event[defs.Push, defs.RepoProvider]) {
	ctx = shared.WithDefaultActivityContext(ctx)

	// msg := comm.NewMergeConflictMessage(event, state.repo, state.author, state.branch(ctx))

	conflict := comm.NewMergeConflictEvent(event, state.pr.HeadBranch, state.pr.BaseBranch, state.last_commit)
	conflict.SetUserID(state.author.ID)

	io := kernel.Instance().MessageIO(state.repo.MessageProvider)

	state.log(ctx, "warn_conflict").Info("message", "payload", conflict, "event", event)

	_ = state.do(ctx, "warn_merge_conflict", io.NotifyMergeConflict, conflict, nil)
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
