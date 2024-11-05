package githubacts

import (
	"context"

	githubdefs "go.breu.io/quantm/internal/hooks/github/defs"
)

type (
	InstallRepos struct{}
)

// AddRepos synchronizes repositories associated with a Github installation. It retrieves all repositories from the.
func (a *InstallRepos) AddRepo(ctx context.Context, webhook *githubdefs.WebhookInstallRepos) error {
	return nil
}
