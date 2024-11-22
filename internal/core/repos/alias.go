package repos

import (
	"go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/core/repos/workflows"
)

type (
	HypdratedRepo = defs.HypdratedRepo
)

var (
	RepoWorkflow        = workflows.Repo
	RepoWorkflowOptions = defs.RepoWorkflowOptions
)
