package reposdefs

import (
	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/db/entities"
)

// signals.
const (
	RepoIOSignalPush queues.Signal = "push" // signals a push event.
)

type (
	// FullRepo is a full representation of a repository.
	FullRepo struct {
		Repo      *entities.Repo      `json:"repo,omitempty"`
		Messaging *entities.Messaging `json:"messaging,omitempty"`
		Org       *entities.Org       `json:"org,omitempty"`
		User      *entities.User      `json:"user,omitempty"`
	}
)
