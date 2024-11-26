package events

import (
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

// -- payloads --

type (
	// Payload represents all available event payloads.
	Payload interface {
		eventsv1.GitRef |
			eventsv1.Push |
			eventsv1.Diff
	}
)
