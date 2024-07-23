// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package core

import (
	"context"
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core/timers"
	"go.breu.io/quantm/internal/shared"
)

type (
	PullRequestAction string
)

// RepoIO signals.
const (
	RepoIOSignalPush           shared.WorkflowSignal = "repo_io__push"
	RepoIOSignalCreateOrDelete shared.WorkflowSignal = "repo_io__create_or_delete"
	ReopIOSignalRebase         shared.WorkflowSignal = "repo_io__rebase"
	RepoIOSignalPullRequest    shared.WorkflowSignal = "repo_io__pull_request"
)

const (
	PullRequestActionCreated       PullRequestAction = "created"
	PullRequestActionLabeled       PullRequestAction = "labeled"
	PullRequestActionClosed        PullRequestAction = "closed"
	PullRequestActionMerged        PullRequestAction = "merged"
	PullRequestActionReviewRequest PullRequestAction = "review_request"
	PullRequestActionApproved      PullRequestAction = "approved"
)

// RepoIO signal payloads.
type (
	// RepoIO is the interface that defines the operations that can be performed on a repository.
	RepoIO interface {
		// GetRepoData gets the name & default branch for the provider repo.
		GetRepoData(ctx context.Context, id string) (*RepoIORepoData, error)

		// SetEarlyWarning sets the early warning flag for the provider repo.
		SetEarlyWarning(ctx context.Context, id string, value bool) error

		// GetAllBranches gets all the branches for the provider repo.
		GetAllBranches(ctx context.Context, payload *RepoIOInfoPayload) ([]string, error)

		DetectChanges(ctx context.Context, payload *RepoIODetectChangesPayload) (*RepoIOChanges, error)

		// TokenizedCloneURL returns the url with oauth token in it.
		//
		// NOTE - Since the url contains oauth token, it is best not to call this as activity.
		// LINK - https://github.com/orgs/community/discussions/24575#discussioncomment-3244524
		TokenizedCloneURL(ctx context.Context, payload *RepoIOInfoPayload) (string, error)
	}

	RepoIOSignalPushPayload struct {
		BranchRef      string         `json:"branch_ref"`
		Before         string         `json:"before"`
		After          string         `json:"after"`
		RepoName       string         `json:"repo_name"`
		RepoOwner      string         `json:"repo_owner"`
		CtrlID         string         `json:"ctrl_id"` // CtrlID represents the id of the provider repo in the quantm DB. Should be UUID.
		InstallationID shared.Int64   `json:"installation_id"`
		ProviderID     string         `json:"provider_id"`
		Commits        RepoIOCommits  `json:"commits"`
		User           *auth.TeamUser `json:"user"`
		Author         string         `json:"author"`
	}

	RepoIOSignalCreatePayload struct {
		IsCreated      bool           `json:"is_created"`
		Ref            string         `json:"ref"`
		RefType        string         `json:"ref_type"`
		DefaultBranch  string         `json:"default_branch"`
		RepoName       string         `json:"repo_name"`
		RepoOwner      string         `json:"repo_owner"`
		CtrlID         string         `json:"ctrl_id"` // CtrlID represents the id of the provider repo in the quantm DB. Should be UUID.
		InstallationID shared.Int64   `json:"installation_id"`
		ProviderID     string         `json:"provider_id"`
		User           *auth.TeamUser `json:"user"`
	}

	RepoIOSignalPullRequestPayload struct {
		Action         string         `json:"action"`
		Number         shared.Int64   `json:"number"`
		RepoName       string         `json:"repo_name"`
		RepoOwner      string         `json:"repo_owner"`
		BaseBranch     string         `json:"base_branch"`
		HeadBranch     string         `json:"head_branch"`
		CtrlID         string         `json:"ctrl_id"`
		InstallationID shared.Int64   `json:"installation_id"`
		ProviderID     string         `json:"provider_id"`
		User           *auth.TeamUser `json:"user"` // TODO: need to find more optimze way
	}
)

// RepoIO types.
type (
	RepoIORepoData struct {
		Name          string `json:"name"`
		DefaultBranch string `json:"default_branch"`
		ProviderID    string `json:"provider_id"`
		Owner         string `json:"owner"`
	}

	RepoIOInfoPayload struct {
		InstallationID shared.Int64 `json:"installation_id"`
		RepoName       string       `json:"repo_name"`
		RepoOwner      string       `json:"repo_owner"`
		DefaultBranch  string       `json:"defualt_branch"`
	}

	RepoIOClonePayload struct {
		Repo   *Repo                    `json:"repo"`   // Repo is the db record of the repo
		Push   *RepoIOSignalPushPayload `json:"push"`   // Push event payload
		Branch string                   `json:"branch"` // Branch to clone
		Path   string                   `json:"path"`   // Path to clone to
	}

	RepoIODetectChangesPayload struct {
		InstallationID shared.Int64 `json:"installation_id"`
		RepoName       string       `json:"repo_name"`
		RepoOwner      string       `json:"repo_owner"`
		DefaultBranch  string       `json:"defualt_branch"`
		TargetBranch   string       `json:"target_branch"`
	}

	RepoIOChanges struct {
		Added      shared.Int64 `json:"added"`
		Removed    shared.Int64 `json:"removed"`
		Modified   []string     `json:"modified"`
		Delta      shared.Int64 `json:"delta"`
		CompareUrl string       `json:"compare_url"`
		RepoUrl    string       `json:"repo_url"`
	}

	RepoIOCommit struct {
		SHA       string        `json:"sha"`
		Message   string        `json:"message"`
		Author    string        `json:"author"`
		Timestamp time.Time     `json:"timestamp"`
		Changes   RepoIOChanges `json:"changes"`
	}

	RepoIOCommits []RepoIOCommit

	RepoIOPullRequest struct {
		Number     shared.Int64 `json:"number"`
		HeadBranch string       `json:"head_branch"`
		BaseBranch string       `json:"base_branch"`
	}

	RepoIORebaseAtCommitResponse struct {
		SHA        string `json:"sha"`
		Message    string `json:"message"`
		InProgress bool   `json:"in_progress"`
	}

	RepoIOGetRepoByProviderIDPayload struct {
		ProviderID string `json:"provider_id"`
	}
)

// Workflow States.
type (
	RepoIOBranchCtrlState struct {
		repo        *Repo                 // Repo is the db record of the repo
		branch      string                // Branch is the name of the branch
		created_at  time.Time             // created_at is the time when the branch was created
		last_commit *RepoIOCommit         // last_commit is the last commit on the branch
		pr          *RepoIOPullRequest    // pr is the pull request associated with the branch
		interval    timers.Interval       // interval is the interval at which the branch is checked for staleness
		mutex       workflow.Mutex        // mutex is the mutex for the state
		active      bool                  // active is the flag to indicate if the branch is active
		counter     int                   // counter is the number of steps taken by the branch
		logger      *RepoIOWorkflowLogger // log is the logger for the branch
	}
)

func (commits RepoIOCommits) Size() int {
	return len(commits)
}

func (commits RepoIOCommits) Latest() *RepoIOCommit {
	if len(commits) == 0 {
		return nil
	}

	return &commits[0]
}

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

// execute is a helper function that executes an activity within the context of the RepoIOBranchCtrlState.
// It logs the start and success of the activity, and increments the steps counter in the state.
//
// The function takes the following parameters:
// - ctx: the workflow context
// - action: a string describing the action being performed
// - activity: the activity function to execute
// - result: a pointer to a variable to receive the result of the activity
// - args: any additional arguments to pass to the activity function
//
// If the activity fails, the function logs the error and returns it.
func (state *RepoIOBranchCtrlState) execute(ctx workflow.Context, action string, activity any, result any, args ...any) error {
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

// check_complexity checks the complexity of the changes pushed on the current branch.
func (state *RepoIOBranchCtrlState) check_complexity(ctx workflow.Context, signal *RepoIOSignalPushPayload) *RepoIOChanges {
	changes := &RepoIOChanges{}
	detect := &RepoIODetectChangesPayload{
		InstallationID: signal.InstallationID,
		RepoName:       signal.RepoName,
		RepoOwner:      signal.RepoOwner,
		DefaultBranch:  state.repo.DefaultBranch,
		TargetBranch:   state.branch,
	}

	_ = state.execute(ctx, "detect_changes", Instance().RepoIO(state.repo.Provider).DetectChanges, changes, detect)

	return changes
}

// increment is a helper function that increments the steps counter in the RepoIOBranchCtrlState.
func (state *RepoIOBranchCtrlState) increment(ctx workflow.Context) {
	_ = state.mutex.Lock(ctx)
	defer state.mutex.Unlock()

	state.counter++
}

// warn_for_complexity sends a warning message to the linked message provider if the complexity of the changes exceeds the threshold.
// it sends a message to the user if the git user is linked to the quantm user and the linked quantm user also has connected the message
// provider, otherwise it sends a message to the linked channel of the repo.
func (state *RepoIOBranchCtrlState) warn_for_complexity(ctx workflow.Context, signal *RepoIOSignalPushPayload, complexity *RepoIOChanges) {
	for_user := signal.User != nil && signal.User.IsMessageProviderLinked
	msg := NewNumberOfLinesExceedMessage(signal, state.repo, state.branch, complexity, for_user)
	io := Instance().MessageIO(state.repo.MessageProvider)

	_ = state.execute(ctx, "send_complexity_warning", io.SendNumberOfLinesExceedMessage, nil, msg)
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
