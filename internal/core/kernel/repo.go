package kernel

import (
	"context"

	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/events"
	commonv1 "go.breu.io/quantm/internal/proto/ctrlplane/common/v1"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	Repo interface {
		// TokenizedCloneUrl returns the tokenized clone URL for the repository with the given ID.
		//
		// This method must not be called from the workflow.
		TokenizedCloneUrl(ctx context.Context, repo *entities.Repo) string

		DetectChanges(ctx context.Context, event *events.Event[commonv1.RepoHook, eventsv1.Push])
	}
)
