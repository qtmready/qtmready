package githubacts

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	githubdefs "go.breu.io/quantm/internal/hooks/github/defs"
)

type (
	Install struct{}
)

// GetOrCreateInstallation retrieves a Github installation from the database by installation ID.
// If the installation does not exist, it creates a new one.
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

// SyncRepos synchronizes repositories associated with a Github installation. It retrieves all repositories from the
// database and compares them to the repositories in the webhook payload. If a repository is missing from the database,
// it's created. This function is designed to be called when a new installation is created. Github provides an
// installation_repositories` webhook event that is used to sync repositories for existing installations.
func (a *Install) SyncRepos(ctx context.Context, webhook *githubdefs.WebhookInstall) error {
	tx, qtx, err := db.Transaction(ctx)
	if err != nil {
		return err
	}

	defer func() { _ = tx.Rollback(ctx) }()

	install, err := qtx.GetGithubInstallationByInstallationID(ctx, webhook.Installation.ID)
	if err != nil {
		return err
	}

	for _, repo := range webhook.Repositories {
		_, err := qtx.GetGithubRepoByInstallationIDAndGithubID(ctx, entities.GetGithubRepoByInstallationIDAndGithubIDParams{
			InstallationID: install.ID,
			GithubID:       repo.ID,
		})

		if err == nil {
			continue
		}

		if errors.Is(err, pgx.ErrNoRows) {
			create := entities.CreateGithubRepoParams{
				InstallationID: install.ID,
				GithubID:       repo.ID,
				Name:           repo.Name,
				FullName:       repo.FullName,
				Url:            fmt.Sprintf("https://github.com/%s", repo.FullName),
			}

			_, err = qtx.CreateGithubRepo(ctx, create)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return tx.Commit(ctx)
}
