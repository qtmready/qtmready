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
	Install struct{}
)

func (a *Install) GetOrCreateInstallation(
	ctx context.Context, install *entities.GithubInstallation,
) (*entities.GithubInstallation, error) {
	response, err := db.Queries().GetGithubInstallationByInstallationID(ctx, install.InstallationID)
	if err == nil {
		return &response, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		create := entities.CreateGithubInstallationParams{
			OrgID:             install.OrgID,
			InstallationID:    install.InstallationID,
			InstallationLogin: install.InstallationLogin,
			InstallationType:  install.InstallationType,
			SenderID:          install.SenderID,
			SenderLogin:       install.SenderLogin,
		}

		response, err = db.Queries().CreateGithubInstallation(ctx, create)
		if err != nil {
			return nil, err
		}

		return &response, nil
	}

	return nil, err
}

func (a *Install) SyncRepos(ctx context.Context, repos []githubdefs.PartialRepository) error {
	tx, qtx, err := db.Transaction(ctx)
	if err != nil {
		return err
	}

	defer func() { _ = tx.Rollback(ctx) }()

	for _, repo := range repos {
		_, err := qtx.GetGithubRepoByGithubID(ctx, repo.ID)
		if err == nil {
			continue
		}

		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	}

	return nil
}
