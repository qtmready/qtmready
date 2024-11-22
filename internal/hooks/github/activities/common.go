package activities

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/events"
	"go.breu.io/quantm/internal/hooks/github/cast"
	"go.breu.io/quantm/internal/hooks/github/defs"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

func HydrateRepoEvent(ctx context.Context, payload *defs.HydrateRepoEventPayload) (*defs.HydratedRepoEvent, error) {
	install, err := db.Queries().GetGithubInstallationByInstallationID(ctx, payload.InstallationID)
	if err != nil {
		return nil, err
	}

	row, err := db.Queries().GetRepo(ctx, entities.GetRepoParams{InstallationID: install.ID, GithubID: payload.RepoID})
	if err != nil {
		return nil, err
	}

	hydrated := cast.RowToHydratedRepoEvent(row)

	// TODO: handle if no rows are fetched.
	if payload.Email != "" {
		u, _ := db.Queries().GetUserByEmail(ctx, payload.Email)
		hydrated.User = &u
	}

	// TODO: Query for Parent ID.

	return hydrated, nil
}

// AddRepo adds a new GitHub repository to the database. If the repository already exists, it will be activated.
func AddRepo(ctx context.Context, payload *defs.SyncRepoPayload) error {
	tx, qtx, err := db.Transaction(ctx)
	if err != nil {
		return err
	}

	defer func() { _ = tx.Rollback(ctx) }()

	repo, err := db.Queries().GetGithubRepoByInstallationIDAndGithubID(ctx, entities.GetGithubRepoByInstallationIDAndGithubIDParams{
		InstallationID: payload.InstallationID,
		GithubID:       payload.Repo.ID,
	})

	if err == nil {
		err = qtx.ActivateGithubRepo(ctx, repo.ID)
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

	created, err := qtx.CreateGithubRepo(ctx, create)
	if err != nil {
		return err
	}

	// create core repo
	reqst := entities.CreateRepoParams{
		OrgID:  payload.OrgID,
		Hook:   int32(eventsv1.RepoHook_REPO_HOOK_GITHUB),
		HookID: created.ID,
		Name:   payload.Repo.Name,
		Url:    fmt.Sprintf("https://github.com/%s", payload.Repo.FullName),
	}

	_, err = qtx.CreateRepo(ctx, reqst)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// SuspendRepo suspends a GitHub repository from the database.
// If the repository does not exist, it will be ignored.
func SuspendRepo(ctx context.Context, payload *defs.SyncRepoPayload) error {
	repo, err := db.Queries().
		GetGithubRepoByInstallationIDAndGithubID(
			ctx,
			entities.GetGithubRepoByInstallationIDAndGithubIDParams{
				InstallationID: payload.InstallationID,
				GithubID:       payload.Repo.ID,
			},
		)

	if err == nil {
		tx, qtx, err := db.Transaction(ctx)
		if err != nil {
			return err
		}

		defer func() { _ = tx.Rollback(ctx) }()

		err = qtx.SuspendedGithubRepo(ctx, repo.ID)
		if err != nil {
			return err
		}

		err = qtx.SuspendedRepoByHookID(ctx, repo.ID)
		if err != nil {
			return err
		}

		if err = tx.Commit(ctx); err != nil {
			return err
		}

		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}

	return err
}

func SignalToRepo[P events.Payload](ctx context.Context, hydrated *defs.HydratedQuantmEvent[P]) error {
	return nil
}
