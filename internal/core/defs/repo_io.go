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

package defs

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

type (
	PullRequestAction string
)

// RepoIO signals.
const (
	RepoIOSignalPush               shared.WorkflowSignal = "repo_io__push"                 // signals a push event.
	RepoIOSignalCreateOrDelete     shared.WorkflowSignal = "repo_io__create_or_delete"     // signals creation or deletion of a repo.
	RepoIOSignalRebase             shared.WorkflowSignal = "repo_io__rebase"               // signals a rebase operation.
	RepoIOSignalPullRequest        shared.WorkflowSignal = "repo_io__pull_request"         // signals a pull request event.
	RepoIOSignalPullRequestLabel   shared.WorkflowSignal = "repo_io__pull_request_label"   // signals a pull request label event.
	RepoIOSignalPullRequestComment shared.WorkflowSignal = "repo_io__pull_request_comment" // signals a pull request comment event.
	RepoIOSignalPullRequestReview  shared.WorkflowSignal = "repo_io__pull_request_review"  // signals a pull request review event.
	RepoIOSignalQueueAdd           shared.WorkflowSignal = "repo_io__queue__add"           // signals an add to the repo queue.
	RepoIOSignalQueueRemove        shared.WorkflowSignal = "repo_io__queue__remove"        // signals a removal from the repo queue.
	RepoIOSignalQueueAddPriority   shared.WorkflowSignal = "repo_io__queue__add__priority" // signals an add with priority to the queue.
	RepoIOSignalQueuePromote       shared.WorkflowSignal = "repo_io__queue__promote"       // signals a promotion in the repo queue.
	RepoIOSignalQueueDemote        shared.WorkflowSignal = "repo_io__queue__demote"        // signals a demotion in the repo queue.
)

type (
	// RepoIOSignalPushPayload represents the payload for the `RepoIOSignalPush` signal.
	RepoIOSignalPushPayload struct {
		BranchRef      string         `json:"branch_ref"`      // BranchRef represents the branch ref.
		Before         string         `json:"before"`          // Before represents the commit SHA before the push.
		After          string         `json:"after"`           // After represents the commit SHA after the push.
		RepoName       string         `json:"repo_name"`       // RepoName represents the name of the repository.
		RepoOwner      string         `json:"repo_owner"`      // RepoOwner represents the owner of the repository.
		CtrlID         string         `json:"ctrl_id"`         // CtrlID represents the id of the provider repo in the quantm DB.
		InstallationID db.Int64       `json:"installation_id"` // InstallationID represents the GitHub installation ID.
		ProviderID     string         `json:"provider_id"`     // ProviderID represents the provider ID.
		Commits        RepoIOCommits  `json:"commits"`         // Commits represents a list of commits in the push.
		User           *auth.TeamUser `json:"user"`            // User represents the user who triggered the push.
		Author         string         `json:"author"`          // Author represents the author of the commits.
	}

	// RepoIOSignalCreateOrDeletePayload represents the payload for the `RepoIOSignalCreateOrDelete` signal.
	RepoIOSignalCreateOrDeletePayload struct {
		IsCreated      bool           `json:"is_created"`      // IsCreated indicates whether the repository was created or deleted.
		Ref            string         `json:"ref"`             // Ref represents the ref of the repository (branch or tag).
		RefType        string         `json:"ref_type"`        // RefType represents the type of the ref (branch or tag).
		DefaultBranch  string         `json:"default_branch"`  // DefaultBranch represents the default branch of the repository.
		RepoName       string         `json:"repo_name"`       // RepoName represents the name of the repository.
		RepoOwner      string         `json:"repo_owner"`      // RepoOwner represents the owner of the repository.
		CtrlID         string         `json:"ctrl_id"`         // CtrlID represents the id of the provider repo in the quantm DB.
		InstallationID db.Int64       `json:"installation_id"` // InstallationID represents the GitHub installation ID.
		ProviderID     string         `json:"provider_id"`     // ProviderID represents the provider ID.
		User           *auth.TeamUser `json:"user"`            // User represents the user who triggered the event.
	}

	// RepoIOSignalPullRequestPayload represents the payload for the `RepoIOSignalPullRequest` signal.
	RepoIOSignalPullRequestPayload struct {
		Action         string         `json:"action"`          // Action represents the action taken on the pull request.
		Number         db.Int64       `json:"number"`          // Number represents the pull request number.
		RepoName       string         `json:"repo_name"`       // RepoName represents the name of the repository.
		RepoOwner      string         `json:"repo_owner"`      // RepoOwner represents the owner of the repository.
		BaseBranch     string         `json:"base_branch"`     // BaseBranch represents the base branch of the pull request.
		HeadBranch     string         `json:"head_branch"`     // HeadBranch represents the head branch of the pull request.
		CtrlID         string         `json:"ctrl_id"`         // CtrlID represents the id of the provider repo in the quantm DB.
		InstallationID db.Int64       `json:"installation_id"` // InstallationID represents the GitHub installation ID.
		ProviderID     string         `json:"provider_id"`     // ProviderID represents the provider ID.
		User           *auth.TeamUser `json:"user"`            // User represents the user who triggered the event.
		LabelName      *string        `json:"label_name"`      // LabelName represents the name of the label.
	}

	// RepoIOSignalPullRequestReviewPayload represents the payload for the `RepoIOSignalPullRequestReview` signal.
	RepoIOSignalPullRequestReviewPayload struct {
		Action         string         `json:"action"`          // Action represents the action taken on the pull request review.
		Number         db.Int64       `json:"number"`          // Number represents the pull request number.
		RepoName       string         `json:"repo_name"`       // RepoName represents the name of the repository.
		RepoOwner      string         `json:"repo_owner"`      // RepoOwner represents the owner of the repository.
		BaseBranch     string         `json:"base_branch"`     // BaseBranch represents the base branch of the pull request.
		HeadBranch     string         `json:"head_branch"`     // HeadBranch represents the head branch of the pull request.
		CtrlID         string         `json:"ctrl_id"`         // CtrlID represents the id of the provider repo in the quantm DB.
		InstallationID db.Int64       `json:"installation_id"` // InstallationID represents the GitHub installation ID.
		ProviderID     string         `json:"provider_id"`     // ProviderID represents the provider ID.
		User           *auth.TeamUser `json:"user"`            // User represents the user who triggered the event.
	}

	// RepoIOSignalPullRequestReviewCommentPayload represents the payload for the `RepoIOSignalPullRequestReviewComment` signal.
	RepoIOSignalPullRequestReviewCommentPayload struct {
		Action         string         `json:"action"`          // Action represents the action taken on the pull request review comment.
		Number         db.Int64       `json:"number"`          // Number represents the pull request number.
		RepoName       string         `json:"repo_name"`       // RepoName represents the name of the repository.
		RepoOwner      string         `json:"repo_owner"`      // RepoOwner represents the owner of the repository.
		BaseBranch     string         `json:"base_branch"`     // BaseBranch represents the base branch of the pull request.
		HeadBranch     string         `json:"head_branch"`     // HeadBranch represents the head branch of the pull request.
		CtrlID         string         `json:"ctrl_id"`         // CtrlID represents the id of the provider repo in the quantm DB.
		InstallationID db.Int64       `json:"installation_id"` // InstallationID represents the GitHub installation ID.
		ProviderID     string         `json:"provider_id"`     // ProviderID represents the provider ID.
		User           *auth.TeamUser `json:"user"`            // User represents the user who triggered the event.
	}

	// RepoIOSignalLabelPayload represents the payload for the `RepoIOSignalLabel` signal.
	RepoIOSignalLabelPayload struct {
		Action         string         `json:"action"`          // Action represents the action taken on the pull request label.
		Number         db.Int64       `json:"number"`          // Number represents the pull request number.
		RepoName       string         `json:"repo_name"`       // RepoName represents the name of the repository.
		RepoOwner      string         `json:"repo_owner"`      // RepoOwner represents the owner of the repository.
		CtrlID         string         `json:"ctrl_id"`         // CtrlID represents the id of the provider repo in the quantm DB.
		InstallationID db.Int64       `json:"installation_id"` // InstallationID represents the GitHub installation ID.
		ProviderID     string         `json:"provider_id"`     // ProviderID represents the provider ID.
		User           *auth.TeamUser `json:"user"`            // User represents the user who triggered the event.
	}

	// RepoIOSignalWorkflowRunPayload represents the payload for the `RepoIOSignalWorkflowRun` signal.
	RepoIOSignalWorkflowRunPayload struct {
		Action         string              `json:"action"`          // Action represents the action taken on the workflow run.
		Number         db.Int64            `json:"number"`          // Number represents the workflow run number.
		RepoName       string              `json:"repo_name"`       // RepoName represents the name of the repository.
		RepoOwner      string              `json:"repo_owner"`      // RepoOwner represents the owner of the repository.
		CtrlID         string              `json:"ctrl_id"`         // CtrlID represents the id of the provider repo in the quantm DB.
		InstallationID db.Int64            `json:"installation_id"` // InstallationID represents the GitHub installation ID.
		ProviderID     string              `json:"provider_id"`     // ProviderID represents the provider ID.
		WorkflowInfo   *RepoIOWorkflowInfo `json:"workflow_info"`   // WorkflowInfo represents the workflow run information.
		User           *auth.TeamUser      `json:"user"`            // User represents the user who triggered the event.
	}
)

// RepoIO types.
type (
	// RepoIOProviderInfo represents the information about a repository from a provider.
	RepoIOProviderInfo struct {
		RepoName       string   `json:"repo_name"`       // RepoName represents the name of the repository.
		RepoOwner      string   `json:"owner"`           // RepoOwner represents the owner of the repository.
		DefaultBranch  string   `json:"default_branch"`  // DefaultBranch represents the default branch of the repository.
		ProviderID     string   `json:"provider_id"`     // ProviderID represents the provider ID.
		InstallationID db.Int64 `json:"installation_id"` // InstallationID represents the GitHub installation ID.
	}

	// RepoIOClonePayload represents the payload for cloning a repository.
	RepoIOClonePayload struct {
		Repo   *Repo               `json:"repo"`   // Repo represents the database record of the repository.
		Push   *Push               `json:"push"`   // Push represents the push event payload.
		Info   *RepoIOProviderInfo `json:"info"`   // Info represents the repository information from the provider.
		Branch string              `json:"branch"` // Branch represents the branch to clone.
		Path   string              `json:"path"`   // Path represents the path to clone the repository to.
	}

	// RepoIOPushBranchPayload represents the payload for pushing a branch to a repository.
	RepoIOPushBranchPayload struct {
		Branch string `json:"branch"` // Branch represents the branch to push.
		Path   string `json:"path"`   // Path represents the path of the repository.
		Force  bool   `json:"force"`  // Force indicates whether to force the push.
	}

	// RepoIODetectChangesPayload represents the payload for detecting changes in a repository.
	RepoIODetectChangesPayload struct {
		InstallationID db.Int64 `json:"installation_id"` // InstallationID represents the GitHub installation ID.
		RepoName       string   `json:"repo_name"`       // RepoName represents the name of the repository.
		RepoOwner      string   `json:"repo_owner"`      // RepoOwner represents the owner of the repository.
		DefaultBranch  string   `json:"defualt_branch"`  // DefaultBranch represents the default branch of the repository.
		TargetBranch   string   `json:"target_branch"`   // TargetBranch represents the target branch for comparison.
	}

	// RepoIOChanges represents the changes detected in a repository.
	RepoIOChanges struct {
		Added      db.Int64 `json:"added"`       // Added represents the number of files added.
		Removed    db.Int64 `json:"removed"`     // Removed represents the number of files removed.
		Modified   []string `json:"modified"`    // Modified represents a list of files modified.
		Delta      db.Int64 `json:"delta"`       // Delta represents the total number of changes.
		CompareUrl string   `json:"compare_url"` // CompareUrl represents the URL to compare the branches.
		RepoUrl    string   `json:"repo_url"`    // RepoUrl represents the URL of the repository.
	}

	// RepoIOCommit represents a commit in a repository.
	RepoIOCommit struct {
		SHA       string        `json:"sha"`       // SHA represents the commit SHA.
		Message   string        `json:"message"`   // Message represents the commit message.
		Author    string        `json:"author"`    // Author represents the author of the commit.
		Timestamp time.Time     `json:"timestamp"` // Timestamp represents the timestamp of the commit.
		Changes   RepoIOChanges `json:"changes"`   // Changes represents the changes made in the commit.
	}

	// RepoIOMergePRPayload represents the payload for merging a pull request.
	RepoIOMergePRPayload struct {
		RepoName       string   `json:"repo_name"`       // RepoName represents the name of the repository.
		RepoOwner      string   `json:"owner"`           // RepoOwner represents the owner of the repository.
		DefaultBranch  string   `json:"defualt_branch"`  // DefaultBranch represents the default branch of the repository.
		TargetBranch   string   `json:"target_branch"`   // TargetBranch represents the target branch for merging.
		InstallationID db.Int64 `json:"installation_id"` // InstallationID represents the GitHub installation ID.
	}

	// RepoIOWorkflowActionPayload represents the payload for performing an action on a workflow.
	RepoIOWorkflowActionPayload struct {
		RepoName       string   `json:"repo_name"`       // RepoName represents the name of the repository.
		RepoOwner      string   `json:"owner"`           // RepoOwner represents the owner of the repository.
		InstallationID db.Int64 `json:"installation_id"` // InstallationID represents the GitHub installation ID.
	}

	// RepoIOWorkflowInfo represents information about workflow runs in a repository.
	RepoIOWorkflowInfo struct {
		TotalCount db.Int64         `json:"total_count"` // TotalCount represents the total number of workflow runs.
		Workflows  []*RepIOWorkflow `json:"workflows"`   // Workflows represents a list of workflow runs.
	}

	// RepIOWorkflow represents a single workflow run.
	RepIOWorkflow struct {
		ID      db.Int64 `json:"id"`       // ID represents the workflow run ID.
		NodeID  string   `json:"node_id"`  // NodeID represents the workflow run node ID.
		Name    string   `json:"name"`     // Name represents the workflow run name.
		Path    string   `json:"path"`     // Path represents the workflow run path.
		State   string   `json:"state"`    // State represents the workflow run state.
		HTMLURL string   `json:"html_url"` // HTMLURL represents the workflow run URL.
	}

	// RepoIOCommits represents a slice of RepoIOCommit.
	RepoIOCommits []RepoIOCommit

	// RepoIOPullRequest represents a pull request.
	RepoIOPullRequest struct {
		Number     db.Int64 `json:"number"`      // Number represents the pull request number.
		HeadBranch string   `json:"head_branch"` // HeadBranch represents the head branch of the pull request.
		BaseBranch string   `json:"base_branch"` // BaseBranch represents the base branch of the pull request.
	}

	// RepoIORebaseAtCommitResponse represents the response for a rebase operation.
	RepoIORebaseAtCommitResponse struct {
		SHA        string `json:"sha"`         // SHA represents the commit SHA after the rebase.
		Message    string `json:"message"`     // Message represents the commit message after the rebase.
		InProgress bool   `json:"in_progress"` // InProgress indicates whether the rebase is in progress.
	}

	// RepoIOGetRepoByProviderIDPayload represents the payload for retrieving a repository by provider ID.
	RepoIOGetRepoByProviderIDPayload struct {
		ProviderID string `json:"provider_id"` // ProviderID represents the provider ID.
	}

	// RepoIOSignalBranchCtrlPayload represents the payload for signaling a branch.
	RepoIOSignalBranchCtrlPayload struct {
		Repo    *Repo                 `json:"repo"`    // Repo represents the database record of the repository.
		Branch  string                `json:"branch"`  // Branch represents the branch to signal.
		Signal  shared.WorkflowSignal `json:"signal"`  // Signal represents the signal to send.
		Payload any                   `json:"payload"` // Payload represents the payload to send.
	}

	// RepoIOSignalQueueCtrlPayload represents the payload for signaling a queue.
	RepoIOSignalQueueCtrlPayload struct {
		Repo    *Repo                 `json:"repo"`    // Repo represents the database record of the repository.
		Branch  string                `json:"branch"`  // Branch represents the branch to signal.
		Signal  shared.WorkflowSignal `json:"signal"`  // Signal represents the signal to send.
		Payload any                   `json:"payload"` // Payload represents the payload to send.
	}
)

// Size returns the number of elements in the RepoIOCommits slice.
//
// Time complexity: O(1).
func (commits RepoIOCommits) Size() int {
	return len(commits)
}

// Latest returns the most recent RepoIOCommit from the provided slice of commits.
//
// If the slice is empty, it returns nil.
//
// This function currently returns the first commit in the slice, which is not necessarily the latest commit based on
// the timestamp. It should be updated to iterate over the commits and return the most recent commit based on the
// timestamp.
func (commits RepoIOCommits) Latest() *RepoIOCommit {
	if len(commits) == 0 {
		return nil
	}

	return &commits[0]
}

// ForBranch returns true if the RefType field of the RepoIOSignalCreateOrDeletePayload is "branch".
//
// This method is used to determine if the signal represents a branch creation or deletion event.
//
// Example:
//
//	signal := &defs.RepoIOSignalCreateOrDeletePayload{
//		RefType: "branch",
//	}
//
//	if signal.ForBranch(ctx) {
//		// Handle branch event
//	}
func (signal *RepoIOSignalCreateOrDeletePayload) ForBranch(ctx workflow.Context) bool {
	return signal.RefType == "branch"
}

// ForTag returns true if the RefType field of the RepoIOSignalCreateOrDeletePayload is "tag".
//
// This method is used to determine if the signal represents a tag creation or deletion event.
//
// Example:
//
//	signal := &defs.RepoIOSignalCreateOrDeletePayload{
//		RefType: "tag",
//	}
//
//	if signal.ForTag(ctx) {
//		// Handle tag event
//	}
func (signal *RepoIOSignalCreateOrDeletePayload) ForTag(ctx workflow.Context) bool {
	return signal.RefType == "tag"
}
