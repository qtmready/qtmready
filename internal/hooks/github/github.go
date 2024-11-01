package github

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/hooks/github/workflows"
)

func InstallWorkflow(ctx workflow.Context) error {
	return workflows.Install(ctx)
}
