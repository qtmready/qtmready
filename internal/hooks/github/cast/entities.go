package cast

import (
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/hooks/github/defs"
)

// RepoForGithubToHydratedRepoEvent converts a database row into a HydratedEvent.
func RepoForGithubToHydratedRepoEvent(row entities.GetRepoForGithubRow) *defs.HydratedRepoEvent {
	return &defs.HydratedRepoEvent{
		Repo:      &row.Repo,
		Org:       &row.Org,
		ChatLinks: &defs.ChatLinks{},
	}
}
