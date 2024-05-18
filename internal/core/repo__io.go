package core

import (
	"context"
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
	// RepoIO is the interface that defines the operations that can be performed on a repository.
	RepoIO interface {
		// GetRepoData gets the name & default branch for the provider repo.
		GetRepoData(ctx context.Context, id string) (*RepoIORepoData, error)
		// SetEarlyWarning sets the early warning flag for the provider repo.
		SetEarlyWarning(ctx context.Context, id string, value bool) error
		GetAllBranches(ctx context.Context, payload *RepoIOGetAllBranchesPayload) ([]string, error)
		Clone(ctx context.Context, payload *RepoIOClonePayload) error

		// GetLatestCommit(ctx context.Context, payload *RepoIOGetLatestCommitPayload) (*LatestCommit, error)
		// DeployChangeset(ctx context.Context, payload *RepoIODeployChangesetPayload) error
		// TagCommit(ctx context.Context, payload *RepoIOTagCommitPayload) error
		// CreateBranch(ctx context.Context, payload *RepoIOCreateBranchPayload) error
		// DeleteBranch(ctx context.Context, payload *RepoIODeleteBranchPayload) error
		// MergeBranch(ctx context.Context, payload *RepoIOMergeBranchPayload) error
		// RebaseAndMerge(ctx context.Context, payload *RepoIORebaseAndMergePayload) (string, error)
		// DetectChange(ctx context.Context, payload *RepoIODetectChangePayload) (*BranchChanges, error)
		// GetAllBranches(ctx context.Context, payload *RepoIOGetAllBranchesPayload) ([]string, error)
		// TriggerCIAction(ctx context.Context, payload *RepoIOTriggerCIActionPayload) error
		// GetRepoTeamID(ctx context.Context, payload *RepoIOGetRepoTeamIDPayload) (string, error)
		// GetAllRelevantActions(ctx context.Context, payload *RepoIOGetAllRelevantActionsPayload) error
		// GetRepoByProviderID(ctx context.Context, payload *RepoIOGetRepoByProviderIDPayload) (*RepoProviderData, error)
		// UpdateRepoHasRarlyWarning(ctx context.Context, payload *RepoIOUpdateRepoHasRarlyWarningPayload) error
		// GetOrgUsers(ctx context.Context, payload *RepoIOGetOrgUsersPayload) error
		// RefreshDefaultBranches(ctx context.Context, payload *RepoIORefreshDefaultBranchesPayload) error
	}

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

	RepoIOClonePayload struct{}

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
)
