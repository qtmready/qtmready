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

	// ChatLinks contains the possible messaging channels for a HydratedRepoEvent.
	ChatLinks struct {
		Org  *entities.ChatLink `json:"org"`
		Team *entities.ChatLink `json:"team"`
		User *entities.ChatLink `json:"user"`
		Repo *entities.ChatLink `json:"repo"`
	}

	// HydratedRepoEvent contains the hydrated event data.
	HydratedRepoEvent struct {
		ParentID  uuid.UUID      `json:"parent_id"`
		Repo      *entities.Repo `json:"repo"`
		Org       *entities.Org  `json:"org"`
		Team      *entities.Team `json:"team"`
		User      *entities.User `json:"user"`
		ChatLinks *ChatLinks     `json:"chat_links"`
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

func (hr *HydratedRepoEvent) GetRepoID() uuid.UUID {
	return hr.Repo.ID
}

func (hr *HydratedRepoEvent) GetOrgID() uuid.UUID {
	return hr.Repo.ID
}

func (hr *HydratedRepoEvent) GetRepoUrl() string {
	return hr.Repo.Url
}

func (hr *HydratedRepoEvent) GetParentID() uuid.UUID {
	return hr.ParentID
}

func (hr *HydratedRepoEvent) GetTeamID() uuid.UUID {
	return hr.Team.ID
}

func (hr *HydratedRepoEvent) GetUserID() uuid.UUID {
	return hr.User.ID
}

func (hr *HydratedRepoEvent) GetRepo() *entities.Repo {
	return hr.Repo
}

func (hr *HydratedRepoEvent) GetTeam() *entities.Team {
	return hr.Team
}

func (hr *HydratedRepoEvent) GetUser() *entities.User {
	return hr.User
}

func (hr *HydratedRepoEvent) GetRepoChatLink() *entities.ChatLink {
	return hr.ChatLinks.Repo
}
