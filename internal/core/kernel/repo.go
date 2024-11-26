package kernel

import (
	"context"

	"go.breu.io/quantm/internal/db/entities"
)

type (
	Repo interface {
		// TokenizedCloneUrl returns the tokenized clone URL for the repository with the given ID.
		//
		// This method must not be called from the workflow.
		TokenizedCloneUrl(ctx context.Context, repo *entities.Repo) (string, error)
	}
)
