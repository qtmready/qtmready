package events

import (
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// EventHook represents a hook for events. It can be either a RepoHook or a MessageHook.
	EventHook interface {
		eventsv1.RepoHook | eventsv1.MessagingHook
	}
)
