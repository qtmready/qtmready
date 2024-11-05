package githubacts

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	githubdefs "go.breu.io/quantm/internal/hooks/github/defs"
)

type (
	InstallRepos struct{}
)

func (a *InstallRepos) RepoAdded(ctx context.Context, payload *githubdefs.SyncRepo) error {
	return AddRepo(ctx, payload)
}

func (a *InstallRepos) RepoRemoved(ctx context.Context, payload *githubdefs.SyncRepo) error {
	return SuspendRepo(ctx, payload)
}

func (a *InstallRepos) GetInstallationforInstallRepos(ctx context.Context, id int64) (*entities.GithubInstallation, error) {
	install, err := db.Queries().GetGithubInstallationByInstallationID(ctx, id)
	if err == nil {
		return &install, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}
