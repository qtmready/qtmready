package reposdefs

import (
	"time"

	"github.com/google/uuid"
	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/db/entities"
)

// signals.
const (
	RepoIOSignalPush queues.Signal = "push" // signals a push event.
)

type (
	CoreRepo struct {
		ID            uuid.UUID           `json:"id"`
		OrgID         uuid.UUID           `json:"org_id"`
		Name          string              `json:"name"`
		Hook          string              `json:"hook"`
		HookID        uuid.UUID           `json:"hook_id"`
		DefaultBranch string              `json:"default_branch"`
		IsMonorepo    bool                `json:"is_monorepo"`
		Threshold     int32               `json:"threshold"`
		StaleDuration time.Duration       `json:"stale_duration"`
		Url           string              `json:"url"`
		IsActive      bool                `json:"is_active"`
		User          *entities.Messaging `json:"user,omitempty"`
		Org           *entities.Org       `json:"org,omitempty"`
	}
)
