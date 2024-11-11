package githubdefs

import (
	"github.com/google/uuid"

	"go.breu.io/quantm/internal/events"
)

type (
	SyncRepo struct {
		InstallationID uuid.UUID         `json:"installation_id"`
		Repo           PartialRepository `json:"repo"`
		IsDeleted      bool              `json:"is_deleted"`
		OrgID          uuid.UUID         `json:"org_id"`
	}

	RepoEventPayload struct {
		RepoID         int64
		InstallationID int64
		Action         events.EventAction
		Scope          events.EventScope

		// TODO - may add the senderID for user.
	}
)
