package activities

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/hooks/github/defs"
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

func (a *Install) AddRepoToInstall(ctx context.Context, payload *defs.SyncRepo) error {
	return AddRepo(ctx, payload)
}
