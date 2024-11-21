package githubcast

import (
	reposdefs "go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/db/entities"
)

func RowToHydratedRepo(row entities.GetRepoRow, user *entities.User) (*reposdefs.HypdratedRepo, error) {
	meta := &reposdefs.HypdratedRepo{}

	if user != nil {
		meta.User = user
	}

	meta.Repo = &row.Repo
	meta.Messaging = &row.Messaging
	meta.Org = &row.Org

	return meta, nil
}
