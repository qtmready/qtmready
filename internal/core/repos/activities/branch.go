package activities

import (
	"context"
	"log/slog"

	git "github.com/jeffwelling/git2go/v37"

	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/db/entities"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	ClonePayload struct {
		Repo   *entities.Repo    `json:"repo"`
		Hook   eventsv1.RepoHook `json:"hook"`
		Branch string            `json:"branch"`
		Path   string            `json:"path"`
	}

	Branch struct{}
)

func (a *Branch) Diff(ctx context.Context) error {
	_ = &git.Diff{}
	return nil
}

func (a *Branch) Clone(ctx context.Context, payload *ClonePayload) error {
	url, err := kernel.Get().RepoHook(payload.Hook).TokenizedCloneUrl(ctx, payload.Repo)
	if err != nil {
		return err
	}

	slog.Info("cloning ...", "url", url)

	return nil
}
