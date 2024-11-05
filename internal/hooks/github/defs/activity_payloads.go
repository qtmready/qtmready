package githubdefs

import (
	"github.com/google/uuid"
)

type (
	SyncRepoActivity struct {
		InstallationID uuid.UUID           `json:"installation_id"`
		Repos          []PartialRepository `json:"repos"`
		IsDeleted      bool                `json:"is_deleted"`
	}
)
