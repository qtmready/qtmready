package events

import (
	commonv1 "go.breu.io/quantm/internal/proto/ctrlplane/common/v1"
)

type (
	// EventHook represents a hook for events. It can be either a RepoHook or a MessageHook.
	EventHook interface {
		commonv1.RepoHook | commonv1.MessagingHook
	}
)
