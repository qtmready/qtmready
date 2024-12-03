package activities

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"go.breu.io/quantm/internal/core/repos"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/events"
	"go.breu.io/quantm/internal/hooks/github/cast"
	"go.breu.io/quantm/internal/hooks/github/defs"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

// HydrateRepoEvent enriches a repository event using database data. It fetches GitHub installation and repository
// details, optionally adding user information if an email is provided. For non-default branches, it retrieves the
// parent event ID from the core workflow, accounting for potential asynchronous delays.
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

	if payload.Email != "" {
		user, _ := db.Queries().GetUserByEmail(ctx, payload.Email)
		hydrated.User = &user
	}

	time.Sleep(500 * time.Second)

	if payload.Branch != "" || payload.Branch != hydrated.Repo.DefaultBranch || payload.ShouldFetchParent {
		parent, err := durable.
			OnCore().
			QueryWorkflow(ctx, hydrated.RepoWorkflowOptions(), repos.QueryRepoForEventParent, payload.Branch)
		if err == nil {
			_ = parent.Get(hydrated.ParentID)
		}
	}

	return hydrated, nil
}

// AddRepo adds a GitHub repository or activates an existing one using a database transaction.  It retrieves the
// repository; if found, it activates it. Otherwise, it creates database entries for both the GitHub and core
// repositories.
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
		return qtx.ActivateGithubRepo(ctx, repo.ID)
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

	return tx.Commit(ctx)
}

// SuspendRepo suspends a GitHub repository, handling cases where it doesn't exist.  It retrieves the repository and, if
// found, suspends both its GitHub and core repository entries using a database transaction.
func SuspendRepo(ctx context.Context, payload *defs.SyncRepoPayload) error {
	repo, err := db.Queries().
		GetGithubRepoByInstallationIDAndGithubID(
			ctx,
			entities.GetGithubRepoByInstallationIDAndGithubIDParams{
				InstallationID: payload.InstallationID,
				GithubID:       payload.Repo.ID,
			},
		)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}

		return err
	}

	tx, qtx, err := db.Transaction(ctx)
	if err != nil {
		return err
	}

	defer func() { _ = tx.Rollback(ctx) }()

	if err := qtx.SuspendedGithubRepo(ctx, repo.ID); err != nil {
		return err
	}

	if err := qtx.SuspendedRepoByHookID(ctx, repo.ID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// SignalRepo signals a GitHub repository event to the core workflow.  Error handling is included.
func SignalRepo[P events.Payload](ctx context.Context, hydrated *defs.HydratedQuantmEvent[P]) error {
	_, err := durable.OnCore().SignalWithStartWorkflow(
		ctx,
		hydrated.Meta.RepoWorkflowOptions(),
		repos.SignalPush,
		hydrated.Event,
		repos.RepoWorkflow,
		repos.NewRepoWorkflowState(hydrated.Meta.Repo, hydrated.Meta.Messaging.Repo),
	)

	return err
}
