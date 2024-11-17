package events

import (
	"time"

	"github.com/google/uuid"
)

type (
	// QuantmEvent is the flat structure of an event for time series databases. It does not
	// contain payload data, only metadata.
	QuantmEvent[H Hook] struct {
		Version  EventVersion `json:"version"`   // Version is the version of the event.
		ID       uuid.UUID    `json:"id"`        // ID is the ID of the event.
		ParentID uuid.UUID    `json:"parent_id"` // ParentID is the ID of the parent event.

		Hook        H           `json:"provider"`     // Provider is the provider of the event.
		Scope       EventScope  `json:"scope"`        // Scope is the scope of the event.
		Action      EventAction `json:"action"`       // Action is the action of the event.
		Source      string      `json:"source"`       // Source is the source of the event. For every hook it will be in different format.
		SubjectID   uuid.UUID   `json:"subject_id"`   // SubjectID is the ID of the subject.
		SubjectName string      `json:"subject_name"` // SubjectName is the name of the subject.
		UserID      uuid.UUID   `json:"user_id"`      // UserID is the ID of the user that the subject belongs to. Can be empty.
		TeamID      uuid.UUID   `json:"team_id"`      // TeamID is the ID of the team that the subject belongs to. Can be empty.
		OrgID       uuid.UUID   `json:"org_id"`       // OrgID is the ID of the organization that the subject belongs to.
		Timestamp   time.Time   `json:"timestamp"`    // Timestamp is the timestamp of the event.
	}
)
