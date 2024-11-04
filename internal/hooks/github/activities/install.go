package githubacts

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
)

type (
	// Install groups all the activities required for the Github Installation.
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

func (a *Install) GetOrCreateRepo(ctx context.Context, entity *entities.GithubRepo) error {
	_, err := db.Queries().GetGithubRepoByGithubID(ctx, entity.GithubID)
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		create := entities.CreateGithubRepoParams{
			RepoID:         entity.RepoID,
			InstallationID: entity.InstallationID,
			GithubID:       entity.GithubID,
			Name:           entity.Name,
			FullName:       entity.FullName,
			Url:            entity.Url,
			IsActive:       pgtype.Bool{Bool: true, Valid: true},
		}

		_, err = db.Queries().CreateGithubRepo(ctx, create)
		if err != nil {
			return err
		}

		return nil
	}

	return err
}
