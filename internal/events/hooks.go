package events

import (
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// Hook represents a hook for events. It can be either a RepoHook or a MessageHook.
	Hook interface {
		eventsv1.RepoHook | eventsv1.MessagingHook
	}
)
