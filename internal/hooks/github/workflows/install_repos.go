package githubwfs

import (
	"go.temporal.io/sdk/workflow"

	githubdefs "go.breu.io/quantm/internal/hooks/github/defs"
)

func InstallRepos(ctx workflow.Context, payload *githubdefs.WebhookInstallRepos) error {
	return nil
}
