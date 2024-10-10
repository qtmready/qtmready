package events

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/gocql/gocql"

	"go.breu.io/quantm/internal/db"
)

type (
	// EventProvider represents a provider for events. It can be either a RepoProvider or a MessageProvider.
	EventProvider interface {
		RepoProvider | MessageProvider
	}

	// EventContext represents the contextual information surrounding an event.
	//
	// This context is crucial for understanding and processing the event.
	EventContext[P EventProvider] struct {
		ParentID  gocql.UUID  `json:"parent_id"` // ParentID is the ID of preceding related event (tracing chains).
		Provider  P           `json:"provider"`  // Provider is the Event origin (e.g., GitHub, GitLab, GCP).
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
		ID     gocql.UUID `json:"id"`      // ID is the ID of the subject.
		Name   string     `json:"name"`    // Name of the database table.
		TeamID gocql.UUID `json:"team_id"` // TeamID is the ID of the team that the subject belongs to.
		UserID gocql.UUID `json:"user_id"` // UserID is the ID of the user that the subject belongs to. It can be null uuid.
	}

	// Event represents an event.
	Event[T EventPayload, P EventProvider] struct {
		Version EventVersion    `json:"version"` // Version is the version of the event.
		ID      gocql.UUID      `json:"id"`      // ID is the ID of the event.
		Context EventContext[P] `json:"context"` // Context is the context of the event.
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
//	event := &Event[EventPayload, EventProvider]{}
//	event.SetSource("example/repo")
func (e *Event[T, P]) SetSource(src string) *Event[T, P] {
	e.Context.Source = src

	return e
}

// SetParent sets the parentID field of the EventContext for the Event struct and returns the event.
//
// The id parameter specifies the parent ID of the event, which can be used to trace the event chain.
//
// Example usage:
//
//	event := &Event[EventPayload, EventProvider]{}
//	event.SetParent(db.MustUUID())
func (e *Event[T, P]) SetParent(id gocql.UUID) *Event[T, P] {
	e.Context.ParentID = id

	return e
}

// SetTimestamp updates the timestamp field of the EventContext for the Event struct and returns the event.
func (e *Event[T, P]) SetTimestamp(t time.Time) *Event[T, P] {
	e.Context.Timestamp = t

	return e
}

// -- Subject --

func (e *Event[T, P]) SetSubjectID(id gocql.UUID) *Event[T, P] {
	e.Subject.ID = id

	return e
}

func (e *Event[T, P]) SetTeamID(id gocql.UUID) *Event[T, P] {
	e.Subject.TeamID = id

	return e
}

// SetUserID sets the UserID field of the EventSubject for the Event struct and returns the event.
func (e *Event[T, P]) SetUserID(id gocql.UUID) *Event[T, P] {
	e.Subject.UserID = id

	return e
}

// -- Action --

// SetActionCreated sets the action of the Event to EventActionCreated and returns the event.
func (e *Event[T, P]) SetActionCreated() *Event[T, P] {
	e.Context.Action = EventActionCreated

	return e
}

// SetActionDeleted sets the action of the Event to EventActionDeleted and returns the event.
func (e *Event[T, P]) SetActionDeleted() *Event[T, P] {
	e.Context.Action = EventActionDeleted

	return e
}

// SetActionUpdated sets the action of the Event to EventActionUpdated and returns the event.
func (e *Event[T, P]) SetActionUpdated() *Event[T, P] {
	e.Context.Action = EventActionUpdated

	return e
}

// SetActionClosed sets the action of the Event to EventActionClosed and returns the event.
func (e *Event[T, P]) SetActionClosed() *Event[T, P] {
	e.Context.Action = EventActionClosed

	return e
}

// SetActionMerged sets the action of the Event to EventActionMerged and returns the event.
func (e *Event[T, P]) SetActionMerged() *Event[T, P] {
	e.Context.Action = EventActionMerged

	return e
}

// SetActionStarted sets the action of the Event to EventActionStarted and returns the event.
func (e *Event[T, P]) SetActionStarted() *Event[T, P] {
	e.Context.Action = EventActionStarted

	return e
}

// SetActionCompleted sets the action of the Event to EventActionCompleted and returns the event.
func (e *Event[T, P]) SetActionCompleted() *Event[T, P] {
	e.Context.Action = EventActionCompleted

	return e
}

// SetActionDismissed sets the action of the Event to EventActionDismissed and returns the event.
func (e *Event[T, P]) SetActionDismissed() *Event[T, P] {
	e.Context.Action = EventActionDismissed

	return e
}

// SetActionAbandoned sets the action of the Event to EventActionAbandoned and returns the event.
func (e *Event[T, P]) SetActionAbandoned() *Event[T, P] {
	e.Context.Action = EventActionAbandoned

	return e
}

// SetActionAdded sets the action of the Event to EventActionAdded and returns the event.
func (e *Event[T, P]) SetActionAdded() *Event[T, P] {
	e.Context.Action = EventActionAdded

	return e
}

// SetActionUnknown sets the action of the Event to EventActionUnknown and returns the event.
func (e *Event[T, P]) SetActionUnknown(in string) *Event[T, P] {
	e.Context.Action = EventAction(in)

	return e
}

// -- Scope --

// SetScopeBranch sets the scope of the Event to EventScopeBranch and returns the event.
func (e *Event[T, P]) SetScopeBranch() *Event[T, P] {
	e.Context.Scope = EventScopeBranch

	return e
}

// SetScopeTag sets the scope of the Event to EventScopeTag and returns the event.
func (e *Event[T, P]) SetScopeTag() *Event[T, P] {
	e.Context.Scope = EventScopeTag

	return e
}

// SetScopePush sets the scope of the Event to EventScopePush and returns the event.
func (e *Event[T, P]) SetScopePush() *Event[T, P] {
	e.Context.Scope = EventScopePush

	return e
}

// SetScopePullRequest sets the scope of the Event to EventScopePullRequest and returns the event.
func (e *Event[T, P]) SetScopePullRequest() *Event[T, P] {
	e.Context.Scope = EventScopePullRequest

	return e
}

// SetScopePullRequestLabel sets the scope of the Event to EventScopePullRequestLabel and returns the event.
func (e *Event[T, P]) SetScopePullRequestLabel() *Event[T, P] {
	e.Context.Scope = EventScopePullRequestLabel

	return e
}

// SetScopePullRequestReview sets the scope of the Event to EventScopePullRequestReview and returns the event.
func (e *Event[T, P]) SetScopePullRequestReview() *Event[T, P] {
	e.Context.Scope = EventScopePullRequestReview

	return e
}

// SetScopePullRequestComment sets the scope of the Event to EventScopePullRequestComment and returns the event.
func (e *Event[T, P]) SetScopePullRequestComment() *Event[T, P] {
	e.Context.Scope = EventScopePullRequestComment

	return e
}

// SetScopePullRequestThread sets the scope of the Event to EventScopePullRequestThread and returns the event.
func (e *Event[T, P]) SetScopePullRequestThread() *Event[T, P] {
	e.Context.Scope = EventScopePullRequestThread

	return e
}

// SetScopeMergeConflict sets the scope of the Event to EventScopeMergeConflict and returns the event.
func (e *Event[T, P]) SetScopeMergeConflict() *Event[T, P] {
	e.Context.Scope = EventScopeMergeConflict

	return e
}

// SetScopeLineExceed sets the scope of the Event to EventScopeLineExceed and returns the event.
func (e *Event[T, P]) SetScopeLineExceed() *Event[T, P] {
	e.Context.Scope = EventScopeLineExceed

	return e
}

// SetScopeRebase sets the scope of the Event to EventScopeRebase and returns the event.
func (e *Event[T, P]) SetScopeRebase() *Event[T, P] {
	e.Context.Scope = EventScopeRebase

	return e
}

// -- Validation --

// Validate validates the Event.
func (e Event[T, P]) Validate() error {
	if err := e.must(); err != nil {
		return err
	}

	switch any(e.Payload).(type) {
	case BranchOrTag:
		return e.types(
			EventScopeBranch,
			EventActionCreated, EventActionDeleted,
		)
	case Push:
		return e.types(
			EventScopePush,
			EventActionCreated, EventActionCreated,
		)
	case PullRequest:
		return e.types(
			EventScopePullRequest,
			EventActionCreated, EventActionUpdated, EventActionReopened, EventActionClosed, EventActionMerged,
		)
	case PullRequestReview:
		return e.types(
			EventScopePullRequestReview,
			EventActionCreated, EventActionUpdated, EventActionDismissed, EventActionRequested,
		)
	case PullRequestLabel:
		return e.types(
			EventScopePullRequestLabel,
			EventActionAdded, EventActionDeleted,
		)
	case PullRequestComment:
		return e.types(
			EventScopePullRequestComment,
			EventActionCreated, EventActionUpdated, EventActionDeleted,
		)
	case PullRequestThread:
		return e.types(
			EventScopePullRequestThread,
			EventActionCreated, EventActionDeleted,
		)
	case RebaseRequest:
		return e.types(
			EventScopeRebase,
			EventActionCreated, EventActionAbandoned,
		)
	case MergeConflict:
		return e.types(
			EventScopeMergeConflict,
			EventActionCreated, EventActionCreated,
		)
	case LinesExceed:
		return e.types(
			EventScopeLineExceed,
			EventActionCreated, EventActionCreated,
		)
	default:
		return nil
	}
}

// must validates the required fields of the Event.
func (e Event[T, P]) must() error {
	if e.Version == "" {
		return NewVersionError("")
	}

	if e.ID.String() == db.NullString {
		return NewIDError(e.ID.String())
	}

	if e.Context.Action == "" {
		return NewActionError(e.Context.Action)
	}

	if e.Context.Scope == "" {
		return NewScopeError(e.Context.Scope)
	}

	if e.Context.Source == "" {
		return NewSourceError(e.Context.Source)
	}

	if e.Context.Timestamp.IsZero() {
		return NewTimestampError(e.Context.Timestamp)
	}

	if e.Subject.ID.String() == db.NullString {
		return NewSubjectIDError(e.Subject.ID.String())
	}

	if e.Subject.TeamID.String() == db.NullString {
		return NewTeamIDError(e.Subject.TeamID.String())
	}

	db_tables := []string{"repos", "stack"}
	if !slices.Contains(db_tables, e.Subject.Name) {
		return NewSubjectNameError(e.Subject.Name)
	}

	return nil
}

func (e Event[T, P]) types(scope EventScope, actions ...EventAction) error {
	if e.Context.Scope != scope {
		return NewScopeError(e.Context.Scope)
	}

	if !slices.Contains(actions, e.Context.Action) {
		return NewActionError(e.Context.Action)
	}

	return nil
}

// -- JSON --

// UnmarshalJSON customizes the JSON decoding for Event.
//
// This method handles the decoding of the Event payload based on the Event's scope.
// The specific payload structure is determined by the EventScope, which allows for different payload types.
//
// For example, an event with an EventScope of EventScopeBranch will have its payload decoded into a BranchOrTag struct.
func (e *Event[T, P]) UnmarshalJSON(data []byte) error {
	aux := struct {
		Version EventVersion    `json:"version"`
		ID      gocql.UUID      `json:"id"`
		Context EventContext[P] `json:"context"`
		Subject EventSubject    `json:"subject"`
		Payload json.RawMessage `json:"payload"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	e.Version = aux.Version
	e.ID = aux.ID
	e.Context = aux.Context
	e.Subject = aux.Subject

	var payload T

	switch e.Context.Scope {
	case EventScopeBranch, EventScopeTag:
		payload = any(BranchOrTag{}).(T)
	case EventScopePullRequest:
		payload = any(PullRequest{}).(T)
	case EventScopeCommit:
		payload = any(Commit{}).(T)
	case EventScopePush:
		payload = any(Push{}).(T)
	case EventScopePullRequestLabel:
		payload = any(PullRequestLabel{}).(T)
	case EventScopePullRequestReview:
		payload = any(PullRequestReview{}).(T)
	case EventScopePullRequestComment:
		payload = any(PullRequestComment{}).(T)
	case EventScopePullRequestThread:
		payload = any(PullRequestThread{}).(T)
	case EventScopeMergeConflict:
		payload = any(MergeConflict{}).(T)
	case EventScopeLineExceed:
		payload = any(LinesExceed{}).(T)
	case EventScopeRebase:
		payload = any(RebaseRequest{}).(T)
	default:
		return fmt.Errorf("unsupported event scope: %s", e.Context.Scope)
	}

	if err := json.Unmarshal(aux.Payload, &payload); err != nil {
		return err
	}

	e.Payload = payload

	return nil
}

// -- Flattening --

// Flatten converts an Event to a FlatEvent.
func (e *Event[T, P]) Flatten() (*FlatEvent[P], error) {
	payload, err := json.Marshal(e.Payload)
	if err != nil {
		return nil, err
	}

	flat := &FlatEvent[P]{
		Version:     e.Version,
		ID:          e.ID,
		ParentID:    e.Context.ParentID,
		Provider:    e.Context.Provider,
		Scope:       e.Context.Scope,
		Action:      e.Context.Action,
		Source:      e.Context.Source,
		SubjectID:   e.Subject.ID,
		SubjectName: e.Subject.Name,
		Payload:     payload,
		TeamID:      e.Subject.TeamID,
		UserID:      e.Subject.UserID,
		CreatedAt:   e.Context.Timestamp,
		UpdatedAt:   e.Context.Timestamp,
	}

	return flat, nil
}
