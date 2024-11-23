package defs

import (
	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/db/entities"
)

// signals.
const (
	SignalPush queues.Signal = "push" // signals a push event.
)

const (
	QueryRepoForEventParent queues.Query = "event_parent"
)

type (
	// HypdratedRepo is a full representation of a repository.
	// NOTE: I think we should keep this github package. it will make our lives easier. easpecially after state.
	HypdratedRepo struct {
		Repo      *entities.Repo      `json:"repo,omitempty"`
		Messaging *entities.Messaging `json:"messaging,omitempty"`
		Org       *entities.Org       `json:"org,omitempty"`
		User      *entities.User      `json:"user,omitempty"`
	}
)
