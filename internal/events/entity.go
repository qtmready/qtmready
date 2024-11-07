package events

import (
	"time"

	"github.com/google/uuid"
)

type (
	// TODO - change name. points to questdb and clickhouse.
	QuantmEvent[H EventHook] struct {
		Version  EventVersion `json:"version" cql:"version"`     // Version is the version of the event.
		ID       uuid.UUID    `json:"id" cql:"id"`               // ID is the ID of the event.
		ParentID uuid.UUID    `json:"parent_id" cql:"parent_id"` // ParentID is the ID of the parent event.

		Hook        H           `json:"provider" cql:"provider"`         // Provider is the provider of the event.
		Scope       EventScope  `json:"scope" cql:"scope_"`              // Scope is the scope of the event.
		Action      EventAction `json:"action" cql:"action"`             // Action is the action of the event.
		Source      string      `json:"source" cql:"source"`             // Source is the source of the event.
		SubjectID   uuid.UUID   `json:"subject_id" cql:"subject_id"`     // SubjectID is the ID of the subject.
		SubjectName string      `json:"subject_name" cql:"subject_name"` // SubjectName is the name of the subject.
		TeamID      uuid.UUID   `json:"team_id" cql:"team_id"`           // TeamID is the ID of the team that the subject belongs to.
		UserID      uuid.UUID   `json:"user_id" cql:"user_id"`           // UserID is the ID of the user that the subject belongs to.
		CreatedAt   time.Time   `json:"created_at" cql:"created_at"`     // CreatedAt is the timestamp when the event was created.
		UpdatedAt   time.Time   `json:"updated_at" cql:"updated_at"`     // UpdatedAt is the timestamp when the event was last updated.

		Payload []byte `json:"payload" cql:"payload"` // Payload is the payload of the event.
	}
)
