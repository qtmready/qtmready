package cast

import (
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/hooks/github/defs"
)

// RowToHydratedRepoEvent converts a database row into a HydratedEvent.
func RowToHydratedRepoEvent(row entities.GetRepoRow) *defs.HydratedRepoEvent {
	return &defs.HydratedRepoEvent{
		Repo:      &row.Repo,
		Org:       &row.Org,
		Messaging: &defs.HydratedRepoEventMessaging{Org: &row.Messaging},
	}
}
