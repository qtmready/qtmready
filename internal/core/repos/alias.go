package repos

import (
	"go.breu.io/quantm/internal/core/repos/activities"
	"go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/core/repos/fns"
	"go.breu.io/quantm/internal/core/repos/nomad"
	"go.breu.io/quantm/internal/core/repos/states"
	"go.breu.io/quantm/internal/core/repos/workflows"
)

var (
	// BranchNameFromRef extracts the branch name from a full Git reference string.
	BranchNameFromRef = fns.BranchNameFromRef

	// BranchNameToRef constructs a full Git reference string from a branch name (e.g., "my-branch" becomes "refs/heads/my-branch").
	BranchNameToRef = fns.BranchNameToRef

	// CreateQuantmRef creates a Git reference string for a Quantm branch (e.g., "my-branch" becomes "refs/heads/qtm/my-branch").
	CreateQuantmRef = fns.CreateQuantmRef

	// IsQuantmRef checks if a Git reference string is a Quantm branch (starts with "refs/heads/qtm/").
	IsQuantmRef = fns.IsQuantmRef

	// IsQuantmBranch checks if a branch name belongs to the Quantm project (starts with "qtm/").
	IsQuantmBranch = fns.IsQuantmBranch
)

var (
	// RepoWorkflow is the main workflow function for managing repository events.
	RepoWorkflow = workflows.Repo

	// BranchWorkflow is the main workflow function for managing branch events.
	BranchWorkflow = workflows.Branch

	// TrunkWorkflow is the main workflow function for managing trunk events.
	TrunkWorkflow = workflows.Trunk

	// NewRepoWorkflowState creates a new state object for the repository workflow.
	NewRepoWorkflowState = states.NewRepo

	// RepoWorkflowOptions provides options for configuring the repository workflow.
	RepoWorkflowOptions = defs.RepoWorkflowOptions
)

var (
	NomadHandler = nomad.NewRepoServiceHandler
)

const (
	SignalPush                     = defs.SignalPush
	SignalRef                      = defs.SignalRef
	SignalPullRequest              = defs.SignalPullRequest
	SignalPullRequestReview        = defs.SignalPullRequestReview
	SignalPullRequestReviewComment = defs.SignalPullRequestReviewComment
	SignalMergeQueue               = defs.SignalMergeQueue
)

const (
	QueryRepoForEventParent = defs.QueryRepoForEventParent
)

const (
	LabelMerge    = defs.LabelMerge
	LabelPriority = defs.LabelPriority
)

// NewRepoActivities creates a new instance of the Activity struct, which handles repository-related actions.
func NewRepoActivities() *activities.Repo {
	return &activities.Repo{}
}

func NewBranchActivities() *activities.Branch {
	return &activities.Branch{}
}

func NewNotifyActivities() *activities.Notify {
	return &activities.Notify{}
}
