package events

import (
	"encoding/json"
	"fmt"
	"time"

	itable "github.com/Guilospanck/igocqlx/table"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/table"

	"go.breu.io/quantm/internal/db"
)

type (

	// FlatEvent is a flattened representation of an event.
	FlatEvent[P EventProvider] struct {
		Version  EventVersion `json:"version" cql:"version"`     // Version is the version of the event.
		ID       gocql.UUID   `json:"id" cql:"id"`               // ID is the ID of the event.
		ParentID gocql.UUID   `json:"parent_id" cql:"parent_id"` // ParentID is the ID of the parent event.
		Provider P            `json:"provider" cql:"provider"`   // Provider is the provider of the event.

		Scope       EventScope  `json:"scope" cql:"scope_"`              // Scope is the scope of the event.
		Action      EventAction `json:"action" cql:"action"`             // Action is the action of the event.
		Source      string      `json:"source" cql:"source"`             // Source is the source of the event.
		SubjectID   gocql.UUID  `json:"subject_id" cql:"subject_id"`     // SubjectID is the ID of the subject.
		SubjectName string      `json:"subject_name" cql:"subject_name"` // SubjectName is the name of the subject.
		TeamID      gocql.UUID  `json:"team_id" cql:"team_id"`           // TeamID is the ID of the team that the subject belongs to.
		UserID      gocql.UUID  `json:"user_id" cql:"user_id"`           // UserID is the ID of the user that the subject belongs to.
		CreatedAt   time.Time   `json:"created_at" cql:"created_at"`     // CreatedAt is the timestamp when the event was created.
		UpdatedAt   time.Time   `json:"updated_at" cql:"updated_at"`     // UpdatedAt is the timestamp when the event was last updated.

		Payload []byte `json:"payload" cql:"payload"` // Payload is the payload of the event.
	}
)

// -- db.Entity implementation --

var (
	// Metadata for FlatEvent table.
	flatEventMeta = itable.Metadata{
		M: &table.Metadata{
			Name: "flat_events__v_0_1",
			Columns: []string{
				"version",
				"id",
				"parent_id",
				"provider",
				"scope_",
				"action",
				"source",
				"subject_id",
				"subject_name",
				"payload",
				"team_id",
				"user_id",
				"created_at",
				"updated_at",
			},
			PartKey: []string{"subject_id", "team_id"},
		},
	}

	// Table instance for FlatEvent.
	flatEventTable = itable.New(*flatEventMeta.M)
)

func (f *FlatEvent[P]) GetTable() itable.ITable {
	return flatEventTable
}

func (f *FlatEvent[P]) PreCreate() error { return nil }
func (f *FlatEvent[P]) PreUpdate() error { return nil }

// -- Database operations --

func (f *FlatEvent[P]) Persist() error {
	return db.CreateWithID(f, f.ID)
}

// -- Deflate --

// Deflate converts a FlatEvent to an Event.
func Deflate[T EventPayload, P EventProvider](flat *FlatEvent[P], event *Event[T, P]) error {
	valid := false

	switch any(event).(type) {
	case *Event[BranchOrTag, P]:
		valid = flat.Scope == EventScopeBranch || flat.Scope == EventScopeTag

	case *Event[PullRequest, P]:
		valid = flat.Scope == EventScopePullRequest

	case *Event[Push, P]:
		valid = flat.Scope == EventScopePush

	case *Event[PullRequestLabel, P]:
		valid = flat.Scope == EventScopePullRequestLabel

	case *Event[PullRequestReview, P]:
		valid = flat.Scope == EventScopePullRequestReview

	case *Event[PullRequestComment, P]:
		valid = flat.Scope == EventScopePullRequestComment

	case *Event[PullRequestThread, P]:
		valid = flat.Scope == EventScopePullRequestThread

	default:
		return fmt.Errorf("unsupported event type: %T", event)
	}

	if !valid {
		return fmt.Errorf("mismatch between event type and scope: %T and scope %s", event, flat.Scope)
	}

	var payload T

	err := json.Unmarshal(flat.Payload, &payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	event.Version = flat.Version
	event.ID = flat.ID
	event.Context = EventContext[P]{
		ParentID:  flat.ParentID,
		Provider:  flat.Provider,
		Scope:     flat.Scope,
		Action:    flat.Action,
		Source:    flat.Source,
		Timestamp: flat.CreatedAt,
	}
	event.Subject = EventSubject{
		ID:     flat.SubjectID,
		Name:   flat.SubjectName,
		TeamID: flat.TeamID,
		UserID: flat.UserID,
	}
	event.Payload = payload

	return nil
}
