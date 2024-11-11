package events

import (
	"github.com/google/uuid"
)

type (
	// Context represents the contextual information surrounding an event.
	//
	// This context is crucial for understanding and processing the event.
	Context[H Hook] struct {
		ParentID uuid.UUID   `json:"parent_id"` // ParentID is the ID of preceding related event (tracing chains).
		Hook     H           `json:"hook"`      // Hook is the Event origin (e.g., GitHub, GitLab, GCP).
		Scope    Scope       `json:"scope"`     // Scope is the Event category (e.g., branch, pull_request).
		Action   EventAction `json:"action"`    // Action is the Triggering action (e.g., created, updated, merged).
		Source   string      `json:"source"`    // Source is the Event source.
	}
)
