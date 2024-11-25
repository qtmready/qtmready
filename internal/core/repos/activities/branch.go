package activities

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	git "github.com/jeffwelling/git2go/v37"

	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	Branch struct{}
)

func (a *Branch) CloneRepo(ctx context.Context, event events.Event[eventsv1.RepoHook, eventsv1.Push]) error {
	repo, err := db.Queries().GetRepoByID(ctx, event.Subject.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}

		return err
	}

	url := kernel.Get().RepoHook(event.Context.Hook).TokenizedCloneUrl(ctx, &repo)
	_, err = git.Clone(url, "", nil)

	return err
}
