package githubdefs

import (
	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/events"
)

type (
	Trigger[H events.EventHook, P events.EventPayload] struct {
		Event *events.Event[H, P]
		Repo  *entities.GetRepoRow
	}
)
