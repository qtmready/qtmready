package core

import (
	"time"

	"go.breu.io/quantm/internal/shared"
)

// RepoIO signals.
const (
	RepoIOSignalPush        shared.WorkflowSignal = "repo__push"
	RepoIOPullRequestLabel  shared.WorkflowSignal = "repo__pull_request__label"
	RepoIOPullRequestMerged shared.WorkflowSignal = "repo__pull_request__merged"
)

// RepoIO signal payloads.
type (
	RepoSignalPushPayload struct {
		BranchRef      string         `json:"branch_ref"`
		Before         string         `json:"before"`
		After          string         `json:"after"`
		Name           string         `json:"name"`
		Owner          string         `json:"owner"`
		CtrlID         string         `json:"ctrl_id"` // ID is the repo ID in the quantm DB. Should be UUID
		InstallationID shared.Int64   `json:"installation_id"`
		ProviderID     string         `json:"provider_id"`
		Commits        []RepoIOCommit `json:"commits"`
	}

	RepoSignalPullRequestLabelPayload struct{}

	RepoSignalPullRequestMergedPayload struct{}
)

// RepoIO types.
type (
	RepoIOChanges struct {
		Added    []string `json:"added"`
		Removed  []string `json:"removed"`
		Modified []string `json:"modified"`
	}

	RepoIOCommit struct {
		SHA       string        `json:"sha"`
		Message   string        `json:"message"`
		Author    string        `json:"author"`
		Timestamp time.Time     `json:"timestamp"`
		Changes   RepoIOChanges `json:"changes"`
	}

	RepoIORepoData struct {
		Name          string `json:"name"`
		DefaultBranch string `json:"default_branch"`
		ProviderID    string `json:"provider_id"`
	}

	RepoIOGetAllBranchesPayload struct {
		InstallationID shared.Int64 `json:"installation_id"`
		RepoName       string       `json:"repo_name"`
		RepoOwner      string       `json:"repo_owner"`
	}
)
