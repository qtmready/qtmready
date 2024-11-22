package defs

import (
	"github.com/google/uuid"

	githubv1 "go.breu.io/quantm/internal/proto/hooks/github/v1"
)

type (
	RequestInstall struct {
		InstallationID int64                `json:"installation_id"`
		SetupAction    githubv1.SetupAction `json:"setup_action"`
		OrgID          uuid.UUID            `json:"org_id"`
	}
)
