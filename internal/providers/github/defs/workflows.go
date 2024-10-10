package defs

import (
	"go.breu.io/quantm/internal/auth"
	coredefs "go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/db"
)

// -- types required during workflows --

type (
	// RepoEventMetadataQuery contains details required to fetch metadata for a repository event.
	RepoEventMetadataQuery struct {
		RepoID         db.Int64 `json:"repo_id"`
		RepoName       string   `json:"repo_name"`
		InstallationID db.Int64 `json:"installation_id"`
		SenderID       string   `json:"sender_id"`
	}

	// RepoEventMetadata is the metadata for a repository event.
	RepoEventMetadata struct {
		CoreRepo *coredefs.Repo `json:"core_repo"`
		Repo     *Repo          `json:"repo"`
		User     *auth.TeamUser `json:"user"`
	}
)
