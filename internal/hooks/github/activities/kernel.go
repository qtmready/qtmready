package activities

import (
	"context"

	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	Kernel struct{}
)

func (k *Kernel) TokenizedCloneUrl(ctx context.Context, repo *entities.Repo) string {
	return ""
}

func (k *Kernel) DetectChanges(ctx context.Context, event *events.Event[eventsv1.RepoHook, eventsv1.Push]) error {
	return nil
}
