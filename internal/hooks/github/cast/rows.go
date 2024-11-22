package cast

import (
	"go.breu.io/quantm/internal/core/repos/defs"
	"go.breu.io/quantm/internal/db/entities"
)

// RowToHydratedRepo converts a database row and user entity into a hydrated repository object.
//
// It takes a GetRepoRow from the database and a pointer to a User entity.  If the user pointer is not nil,
// the hydrated repository will include the user information.  The function returns a pointer to a
// HypdratedRepo struct and an error.  The error will be nil if the conversion is successful.
func RowToHydratedRepo(row entities.GetRepoRow, user *entities.User) (*defs.HypdratedRepo, error) {
	hydrated := &defs.HypdratedRepo{}

	if user != nil {
		hydrated.User = user
	}

	hydrated.Repo = &row.Repo
	hydrated.Messaging = &row.Messaging
	hydrated.Org = &row.Org

	return hydrated, nil
}
