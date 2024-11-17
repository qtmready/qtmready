package githubacts

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.breu.io/durex/queues"

	reposdefs "go.breu.io/quantm/internal/core/repos/defs"
	reposwfs "go.breu.io/quantm/internal/core/repos/workflows"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/events"
	githubcast "go.breu.io/quantm/internal/hooks/github/cast"
	githubdefs "go.breu.io/quantm/internal/hooks/github/defs"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

func PopulateRepoEvent[H eventsv1.RepoHook, P events.Payload](
	ctx context.Context, params *githubdefs.RepoEventPayload,
) (*githubdefs.RepoEvent[H, P], error) {
	var event *events.Event[H, P]

	install, err := db.Queries().GetGithubInstallationByInstallationID(ctx, params.InstallationID)
	if err != nil {
		return nil, err
	}

	// get the core repo from hook_repo (join)
	// TODO - may change the query and get the user and team info
	row, err := db.Queries().GetRepo(ctx, entities.GetRepoParams{
		InstallationID: install.ID,
		GithubID:       params.RepoID,
	})
	if err != nil {
		return nil, err
	}

	// convert the messaging and org byte into entity
	repo, err := githubcast.ConvertGetRepoRowToCoreRepo(row)
	if err != nil {
		return nil, err
	}

	id := events.MustUUID()
	event = &events.Event[H, P]{
		ID:      id,
		Version: events.EventVersionDefault,
		Context: events.Context[H]{
			ParentID:  id,
			Hook:      H(eventsv1.RepoHook_REPO_HOOK_GITHUB), // FIXME: the should come from core repo, (ysf)
			Scope:     params.Scope,
			Action:    params.Action,
			Source:    repo.Url,
			Timestamp: time.Now(),
		},
		Subject: events.Subject{
			ID:     repo.ID,
			Name:   repo.Name,
			OrgID:  install.OrgID,
			TeamID: uuid.Nil, // TODO - need to set after github oauth flow is done? why? (ysf)
			UserID: uuid.Nil, // TODO - need to set after github oauth flow is done? why? we discussed this and it comes from email (ysf)
		},
	}

	return &githubdefs.RepoEvent[H, P]{Event: event, Repo: repo}, nil
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
	ctx context.Context, repo *reposdefs.CoreRepo, signal queues.Signal, payload any,
) error {
	_, err := durable.OnCore().SignalWithStartWorkflow(
		ctx,
		reposdefs.RepoWorkflowOptions("", repo.Name, repo.ID),
		signal,
		payload,
		reposwfs.Repo,
		repo,
	)

	return err
}
