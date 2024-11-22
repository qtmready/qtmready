package repos

import (
	"go.breu.io/quantm/internal/core/repos/activities"
	"go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/core/repos/workflows"
)

type (
	HypdratedRepo = defs.HypdratedRepo
)

var (
	RepoWorkflow         = workflows.Repo
	NewRepoWorkflowState = workflows.NewRepoState
	RepoWorkflowOptions  = defs.RepoWorkflowOptions
)

func NewActivities() *activities.Activity {
	return &activities.Activity{}
}
