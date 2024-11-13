package githubacts

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.breu.io/durex/queues"

	coredefs "go.breu.io/quantm/internal/core/repos/defs"
	corewfs "go.breu.io/quantm/internal/core/repos/workflows"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/events"
	githubdefs "go.breu.io/quantm/internal/hooks/github/defs"
	commonv1 "go.breu.io/quantm/internal/proto/ctrlplane/common/v1"
)

func PopulateRepoEvent[H events.EventHook, P events.EventPayload](
	ctx context.Context, params *githubdefs.RepoEventPayload,
) (*githubdefs.Eventory[H, P], error) {
	var event events.Event[H, P]

	install, err := db.Queries().GetGithubInstallationByInstallationID(ctx, params.InstallationID)
	if err != nil {
		return nil, nil
	}

	// get the core repo from hook_repo (join)
	// TODO - may change the query and get the user and team info
	// TODO - convert the messaging byte into entity
	repo, err := db.Queries().GetRepo(ctx, entities.GetRepoParams{
		InstallationID: install.ID,
		GithubID:       params.RepoID,
	})
	if err != nil {
		return nil, nil
	}

	id := uuid.New()

	event = events.Event[H, P]{
		ID:      id,
		Version: events.EventVersionDefault,
		Context: events.EventContext[H]{
			ParentID:  id,
			Hook:      H(commonv1.RepoHook_REPO_HOOK_GITHUB),
			Scope:     params.Scope,
			Action:    params.Action,
			Source:    repo.Url,
			Timestamp: time.Now(),
		},
		Subject: events.EventSubject{
			ID:     repo.ID,
			Name:   repo.Name,
			OrgID:  install.OrgID,
			TeamID: uuid.Nil, // TODO - need to set after github oauth flow is done
			UserID: uuid.Nil, // TODO - need to set after github oauth flow is done
		},
	}

	tr := &githubdefs.Eventory[H, P]{
		Event: &event,
		Repo:  &repo,
	}

	return tr, nil
}

// AddRepo adds a new GitHub repository to the database.
// If the repository already exists, it will be activated.
func AddRepo(ctx context.Context, payload *githubdefs.SyncRepo) error {
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
		Hook:   "github",
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
func SuspendRepo(ctx context.Context, payload *githubdefs.SyncRepo) error {
	repo, err := db.Queries().GetGithubRepoByInstallationIDAndGithubID(ctx, entities.GetGithubRepoByInstallationIDAndGithubIDParams{
		InstallationID: payload.InstallationID,
		GithubID:       payload.Repo.ID,
	})

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

// SignalCoreRepo signals the core repository control workflow with the given signal and payload.
func SignalCoreRepo(
	ctx context.Context, repo *entities.GetRepoRow, signal queues.Signal, payload any,
) error {
	_, err := durable.OnCore().SignalWithStartWorkflow(
		ctx,
		coredefs.RepoWorkflowOptions("", repo.Name, repo.ID),
		signal,
		payload,
		corewfs.Repo,
		repo,
	)

	return err
}
