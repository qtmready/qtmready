package githubacts

import (
	"context"

	githubdefs "go.breu.io/quantm/internal/hooks/github/defs"
)

type (
	InstallRepos struct{}
)

func (a *InstallRepos) AddRepo(ctx context.Context, payload *githubdefs.SyncRepo) error {
	return AddRepo(ctx, payload)
}

func (a *InstallRepos) SuspendRepo(ctx context.Context, payload *githubdefs.SyncRepo) error {
	return SuspendRepo(ctx, payload)
}
