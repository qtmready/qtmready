package events

import (
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

// -- payloads --

type (
	// EventPayload represents all available event payloads.
	EventPayload interface {
		eventsv1.GitRef |
			eventsv1.Push |
			eventsv1.DetectRepoChange | eventsv1.RepoChange
	}
)
