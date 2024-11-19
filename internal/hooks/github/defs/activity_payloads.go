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
		RepoID         int64              `json:"repo_id"`
		InstallationID int64              `json:"installation_id"`
		Action         events.EventAction `json:"action"`
		Scope          events.EventScope  `json:"scope"`
		Email          string             `json:"email"` // get user by email
	}
)
