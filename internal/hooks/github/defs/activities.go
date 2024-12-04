package defs

import (
	"github.com/google/uuid"
	"go.breu.io/durex/queues"
	"go.breu.io/durex/workflows"

	"go.breu.io/quantm/internal/core/repos"
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// SyncRepoPayload is the payload for the SyncRepo activity.
	SyncRepoPayload struct {
		InstallationID uuid.UUID         `json:"installation_id"`
		Repo           PartialRepository `json:"repo"`
		IsDeleted      bool              `json:"is_deleted"`
		OrgID          uuid.UUID         `json:"org_id"`
	}

	// HydrateRepoEventPayload is the payload for the HydrateRepoEvent activity.
	HydrateRepoEventPayload struct {
		RepoID            int64  `json:"repo_id"`
		InstallationID    int64  `json:"installation_id"`
		Email             string `json:"email"`
		Branch            string `json:"branch"`
		ShouldFetchParent bool   `json:"should_fetch_parent"`
	}

	// HydratedRepoEventMessaging contains the possible messaging channels for a HydratedRepoEvent.
	HydratedRepoEventMessaging struct {
		Org  *entities.Messaging `json:"org"`
		Team *entities.Messaging `json:"team"`
		User *entities.Messaging `json:"user"`
		Repo *entities.Messaging `json:"repo"`
	}

	// HydratedRepoEvent contains the hydrated event data.
	HydratedRepoEvent struct {
		ParentID  uuid.UUID                   `json:"parent_id"`
		Repo      *entities.Repo              `json:"repo"`
		Org       *entities.Org               `json:"org"`
		Team      *entities.Team              `json:"team"`
		User      *entities.User              `json:"user"`
		Messaging *HydratedRepoEventMessaging `json:"messaging"`
	}

	// HydratedQuantmEvent is the hydrated event data for a Quantm event.
	HydratedQuantmEvent[P events.Payload] struct {
		Event  *events.Event[eventsv1.RepoHook, P] `json:"event"`
		Meta   *HydratedRepoEvent                  `json:"meta"`
		Signal queues.Signal                       `json:"signal"`
	}
)

func (h *HydratedRepoEvent) RepoWorkflowOptions() workflows.Options {
	return repos.RepoWorkflowOptions(h.Repo)
}
