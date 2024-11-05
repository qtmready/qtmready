package githubdefs

import (
	"github.com/google/uuid"
)

type (
	SyncRepo struct {
		InstallationID uuid.UUID         `json:"installation_id"`
		Repo           PartialRepository `json:"repo"`
		IsDeleted      bool              `json:"is_deleted"`
	}
)
