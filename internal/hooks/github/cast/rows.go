package githubcast

import (
	reposdefs "go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/db/entities"
)

func RowToFullRepo(row entities.GetRepoRow, user *entities.User) (*reposdefs.FullRepo, error) {
	meta := &reposdefs.FullRepo{}

	if user != nil {
		meta.User = user
	}

	meta.Repo = &row.Repo
	meta.Messaging = &row.Messaging
	meta.Org = &row.Org

	return meta, nil
}
