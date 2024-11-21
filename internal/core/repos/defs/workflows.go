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
	// HypdratedRepo is a full representation of a repository.
	HypdratedRepo struct {
		Repo      *entities.Repo      `json:"repo,omitempty"`
		Messaging *entities.Messaging `json:"messaging,omitempty"`
		Org       *entities.Org       `json:"org,omitempty"`
		User      *entities.User      `json:"user,omitempty"` // the user will change at every event. We don't need this.
	}
)
