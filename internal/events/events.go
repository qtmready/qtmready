package events

import (
	"time"

	"github.com/google/uuid"
)

type (
	// Event represents an event.
	Event[H Hook, P Payload] struct {
		Version EventVersion `json:"version"` // Version is the version of the event.
		ID      uuid.UUID    `json:"id"`      // ID is the ID of the event.
		Context Context[H]   `json:"context"` // Context is the context of the event.
		Subject Subject      `json:"subject"` // Subject is the subject of the event.
		Payload P            `json:"payload"` // Payload is the payload of the event.
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
func (e *Event[H, P]) SetSource(src string) *Event[H, P] {
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
func (e *Event[H, P]) SetParent(id uuid.UUID) *Event[H, P] {
	e.Context.ParentID = id

	return e
}

// SetTimestamp updates the timestamp field of the EventContext for the Event struct and returns the event.
func (e *Event[H, P]) SetTimestamp(t time.Time) *Event[H, P] {
	e.Context.Timestamp = t

	return e
}

// -- Scope --

// SetScopeBranch sets the scope of the Event to EventScopeBranch and returns the event.
func (e *Event[H, P]) SetScopeBranch() *Event[H, P] {
	e.Context.Scope = EventScopeBranch

	return e
}

// SetScopeTag sets the scope of the Event to EventScopeTag and returns the event.
func (e *Event[H, P]) SetScopeTag() *Event[H, P] {
	e.Context.Scope = EventScopeTag

	return e
}

// SetScopePush sets the scope of the Event to EventScopePush and returns the event.
func (e *Event[H, P]) SetScopePush() *Event[H, P] {
	e.Context.Scope = EventScopePush

	return e
}
