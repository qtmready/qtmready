package events

import (
	"time"

	"github.com/google/uuid"
)

type (
	// EventHook represents a hook for events. It can be either a RepoHook or a MessageHook.
	EventHook interface {
		RepoHook | MessageHook
	}

	// EventContext represents the contextual information surrounding an event.
	//
	// This context is crucial for understanding and processing the event.
	EventContext[H EventHook] struct {
		ParentID  uuid.UUID   `json:"parent_id"` // ParentID is the ID of preceding related event (tracing chains).
		Hook      H           `json:"hook"`      // Hook is the Event origin (e.g., GitHub, GitLab, GCP).
		Scope     EventScope  `json:"scope"`     // Scope is the Event category (e.g., branch, pull_request).
		Action    EventAction `json:"action"`    // Action is the Triggering action (e.g., created, updated, merged).
		Source    string      `json:"source"`    // Source is the Event source.
		Timestamp time.Time   `json:"timestamp"` // Timestamp is the Event occurrence time.
	}

	// EventSubject represents the entity within the system that is the subject of an event.
	//
	// It encapsulates:
	//   - ID: The unique identifier of the entity i.e. the primary key within its respective database table.
	//   - Name: The name of the entity's corresponding database table. This provides context for the event's subject.
	//     For instance, an event related to a branch would have "repos" as the subject name, as branches are associated
	//     with repositories.
	//   - TeamID: The unique identifier of the team to which this entity belongs. This allows for team-based filtering
	//     and organization
	//     of events.
	EventSubject struct {
		ID     uuid.UUID `json:"id"`      // ID is the ID of the subject.
		Name   string    `json:"name"`    // Name of the database table.
		TeamID uuid.UUID `json:"team_id"` // TeamID is the ID of the team that the subject belongs to.
		UserID uuid.UUID `json:"user_id"` // UserID is the ID of the user that the subject belongs to. It can be null uuid.
	}

	// Event represents an event.
	Event[T EventPayload, H EventHook] struct {
		Version EventVersion    `json:"version"` // Version is the version of the event.
		ID      uuid.UUID       `json:"id"`      // ID is the ID of the event.
		Context EventContext[H] `json:"context"` // Context is the context of the event.
		Subject EventSubject    `json:"subject"` // Subject is the subject of the event.
		Payload T               `json:"payload"` // Payload is the payload of the event.
	}
)

// -- Context --

// SetSource sets the source field of the EventContext for the Event struct and returns the event.
//
// The src parameter specifies the source of the event, such as the name of the repository.
//
// Example usage:
//
//	event := &Event[EventPayload, EventHook]{}
//	event.SetSource("example/repo")
func (e *Event[T, H]) SetSource(src string) *Event[T, H] {
	e.Context.Source = src

	return e
}

// SetParent sets the parentID field of the EventContext for the Event struct and returns the event.
//
// The id parameter specifies the parent ID of the event, which can be used to trace the event chain.
//
// Example usage:
//
//	event := &Event[EventPayload, EventHook]{}
//	event.SetParent(id)
func (e *Event[T, H]) SetParent(id uuid.UUID) *Event[T, H] {
	e.Context.ParentID = id

	return e
}

// SetTimestamp updates the timestamp field of the EventContext for the Event struct and returns the event.
func (e *Event[T, H]) SetTimestamp(t time.Time) *Event[T, H] {
	e.Context.Timestamp = t

	return e
}

// -- Scope --

// SetScopeBranch sets the scope of the Event to EventScopeBranch and returns the event.
func (e *Event[T, H]) SetScopeBranch() *Event[T, H] {
	e.Context.Scope = EventScopeBranch

	return e
}

// SetScopeTag sets the scope of the Event to EventScopeTag and returns the event.
func (e *Event[T, H]) SetScopeTag() *Event[T, H] {
	e.Context.Scope = EventScopeTag

	return e
}

// SetScopePush sets the scope of the Event to EventScopePush and returns the event.
func (e *Event[T, H]) SetScopePush() *Event[T, H] {
	e.Context.Scope = EventScopePush

	return e
}
