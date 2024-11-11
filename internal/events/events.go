package events

import (
	"time"

	"github.com/google/uuid"
)

type (
	// Event represents an event.
	Event[H Hook, P Payload] struct {
		Version   EventVersion `json:"version"`   // Version is the version of the event.
		ID        uuid.UUID    `json:"id"`        // ID is the ID of the event.
		Timestamp time.Time    `json:"timestamp"` // Timestamp is the Event occurrence time.
		Context   Context[H]   `json:"context"`   // Context is the context of the event.
		Subject   Subject      `json:"subject"`   // Subject is the subject of the event.
		Payload   *P           `json:"payload"`   // Payload is the payload of the event.
	}
)

func (e *Event[H, P]) SetParent(id uuid.UUID) *Event[H, P] {
	e.Context.ParentID = id

	return e
}

func (e *Event[H, P]) SettHook(hook H) *Event[H, P] {
	e.Context.Hook = hook

	return e
}

func (e *Event[H, P]) SetScope(scope Scope) *Event[H, P] {
	e.Context.Scope = scope

	return e
}

func (e *Event[H, P]) SetAction(action EventAction) *Event[H, P] {
	e.Context.Action = action

	return e
}

func (e *Event[H, P]) SetSource(source string) *Event[H, P] {
	e.Context.Source = source

	return e
}

func (e *Event[H, P]) SetSubjectID(id uuid.UUID) *Event[H, P] {
	e.Subject.ID = id

	return e
}

func (e *Event[H, P]) SetSubjectName(name string) *Event[H, P] {
	e.Subject.Name = name

	return e
}

func (e *Event[H, P]) SetOrg(id uuid.UUID) *Event[H, P] {
	e.Subject.OrgID = id

	return e
}

func (e *Event[H, P]) SetTeam(id uuid.UUID) *Event[H, P] {
	e.Subject.TeamID = id

	return e
}

func (e *Event[H, P]) SetUser(id uuid.UUID) *Event[H, P] {
	e.Subject.UserID = id

	return e
}

func (e *Event[H, P]) SetPayload(payload *P) *Event[H, P] {
	e.Payload = payload

	return e
}

func New[H Hook, P Payload]() *Event[H, P] {
	event := &Event[H, P]{
		Version:   EventVersionDefault,
		ID:        MustUUID(),
		Timestamp: time.Now(),
		Context:   Context[H]{},
		Subject:   Subject{},
	}

	return event
}
