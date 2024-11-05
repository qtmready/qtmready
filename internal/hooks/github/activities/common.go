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

// AddRepo adds a new GitHub repository to the database.
// If the repository already exists, it will be activated.
func AddRepo(ctx context.Context, payload *githubdefs.SyncRepo) error {
	repo, err := db.Queries().GetGithubRepoByInstallationIDAndGithubID(ctx, entities.GetGithubRepoByInstallationIDAndGithubIDParams{
		InstallationID: payload.InstallationID,
		GithubID:       payload.Repo.ID,
	})

	if err == nil {
		err = db.Queries().ActivateGithubRepo(ctx, repo.ID)
		if err != nil {
			return err
		}

		return nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	create := entities.CreateGithubRepoParams{
		InstallationID: payload.InstallationID,
		GithubID:       payload.Repo.ID,
		Name:           payload.Repo.Name,
		FullName:       payload.Repo.FullName,
		Url:            fmt.Sprintf("https://github.com/%s", payload.Repo.FullName),
	}

	_, err = db.Queries().CreateGithubRepo(ctx, create)

	return err
}

// SuspendRepo suspends a GitHub repository from the database.
// If the repository does not exist, it will be ignored.
func SuspendRepo(ctx context.Context, payload *githubdefs.SyncRepo) error {
	repo, err := db.Queries().GetGithubRepoByInstallationIDAndGithubID(ctx, entities.GetGithubRepoByInstallationIDAndGithubIDParams{
		InstallationID: payload.InstallationID,
		GithubID:       payload.Repo.ID,
	})

	if err == nil {
		err = db.Queries().SuspendedGithubRepo(ctx, repo.ID)
		if err != nil {
			return err
		}

		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}

	return err
}
