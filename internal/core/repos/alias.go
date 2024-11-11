package repos

import (
	"go.breu.io/quantm/internal/core/repos/activities"
	"go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/core/repos/fns"
	"go.breu.io/quantm/internal/core/repos/workflows"
)

// HypdratedRepo represents a fully hydrated repository, including its associated messaging, organization, and user data.
type (
	HypdratedRepo = defs.HypdratedRepo
)

var (
	// BranchNameFromRef extracts the branch name from a full Git reference string.
	//
	// Example:
	//
	//  BranchNameToRef("refs/head/name") // "name"
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

	// NewRepoWorkflowState creates a new state object for the repository workflow.
	NewRepoWorkflowState = workflows.NewRepoState

	// RepoWorkflowOptions provides options for configuring the repository workflow.
	RepoWorkflowOptions = defs.RepoWorkflowOptions
)

const (
	// QueryRepoForEventParent is a query used to find the parent event of a given branch.
	QueryRepoForEventParent = defs.QueryRepoForEventParent
)

// NewActivities creates a new instance of the Activity struct, which handles repository-related actions.
func NewActivities() *activities.Activity {
	return &activities.Activity{}
}
