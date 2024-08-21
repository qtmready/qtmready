package defs

import (
	"encoding/json"
	"fmt"
	"time"

	itable "github.com/Guilospanck/igocqlx/table"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/table"
)

type (
	EventScope   string
	EventAction  string
	EventVersion string

	// EventContext is the context of an event.
	EventContext[P MessageProvider | RepoProvider] struct {
		ID        gocql.UUID  `json:"id"`
		ParentID  gocql.UUID  `json:"parent_id"`
		Provider  P           `json:"provider"`
		Scope     EventScope  `json:"scope"`
		Action    EventAction `json:"action"`
		Source    string      `json:"source"`
		Timestamp time.Time   `json:"timestamp"`
	}

	// EventSubject refers to the entity in quantm db that is the subject of an event.
	// ID is the primary key of the entity in the database while name is the name of the entity (db table).
	// For example, if the event is about a branch, since subject name represents the table, and since branch belong to repo,
	// in our case, subject name will be repos, while subject id will be the id for the repo to which the branch belongs.
	EventSubject struct {
		ID   gocql.UUID `json:"id"`
		Name string     `json:"name"`
	}

	Event[T any, P MessageProvider | RepoProvider] struct {
		Version EventVersion    `json:"version"`
		Context EventContext[P] `json:"context"`
		Subject EventSubject    `json:"subject"`
		Data    T               `json:"data"`
	}

	FlatEvent[P MessageProvider | RepoProvider] struct {
		Version     EventVersion `json:"version" cql:"version"`
		ID          gocql.UUID   `json:"id" cql:"id"`
		ParentID    gocql.UUID   `json:"parent_id" cql:"parent_id"`
		Provider    P            `json:"provider" cql:"provider"`
		Scope       EventScope   `json:"scope" cql:"scope"`
		Action      EventAction  `json:"action" cql:"action"`
		Source      string       `json:"source" cql:"source"`
		SubjectID   gocql.UUID   `json:"subject_id" cql:"subject_id"`
		SubjectName string       `json:"subject_name" cql:"subject_name"`
		Data        []byte       `json:"data" cql:"data"`
		CreatedAt   time.Time    `json:"created_at" cql:"created_at"`
		UpdatedAt   time.Time    `json:"updated_at" cql:"updated_at"`
	}

	BranchOrTag struct {
		Ref           string `json:"ref"`
		DefaultBranch string `json:"default_branch"`
	}

	Commit struct {
		SHA       string    `json:"sha"`
		Message   string    `json:"message"`
		Author    string    `json:"author"`
		Committer string    `json:"committer"`
		Timestamp time.Time `json:"timestamp"`
		URL       string    `json:"url"`
		Added     []string  `json:"added"`
		Removed   []string  `json:"removed"`
		Modified  []string  `json:"modified"`
	}

	PullRequest struct {
		Number         int       `json:"number"`
		Title          string    `json:"title"`
		Body           string    `json:"body"`
		State          string    `json:"state"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
		ClosedAt       time.Time `json:"closed_at"`
		MergedAt       time.Time `json:"merged_at"`
		MergeCommitSHA string    `json:"merge_commit_sha"`
		Author         string    `json:"author"`
		HeadBranch     string    `json:"head_branch"`
		BaseBranch     string    `json:"base_branch"`
	}

	Push struct {
		Ref        string    `json:"ref"`
		Before     string    `json:"before"`
		After      string    `json:"after"`
		Repository string    `json:"repository"`
		Pusher     string    `json:"pusher"`
		Commits    []Commit  `json:"commits"`
		Timestamp  time.Time `json:"timestamp"`
	}

	PullRequestReview struct {
		ID                int       `json:"id"`
		Body              string    `json:"body"`
		State             string    `json:"state"`
		SubmittedAt       time.Time `json:"submitted_at"`
		Author            string    `json:"author"`
		PullRequestNumber int       `json:"pull_request_number"`
	}

	PullRequestLabel struct {
		Name              string    `json:"name"`
		Color             string    `json:"color"`
		Description       string    `json:"description"`
		CreatedAt         time.Time `json:"created_at"`
		UpdatedAt         time.Time `json:"updated_at"`
		PullRequestNumber int       `json:"pull_request_number"`
	}
)

const (
	EventVersionV1 EventVersion = "v1"

	EventScopeBranch            EventScope = "branch"
	EventScopeTag               EventScope = "tag"
	EventScopePullRequest       EventScope = "pull_request"
	EventScopePush              EventScope = "push"
	EventScopePullRequestLabel  EventScope = "pull_request_label"
	EventScopePullRequestReview EventScope = "pull_request_review"

	EventTypeCreated   EventAction = "created"
	EventTypeDeleted   EventAction = "deleted"
	EventTypeUpdated   EventAction = "updated"
	EventTypeClosed    EventAction = "closed"
	EventTypeMerged    EventAction = "merged"
	EventTypeStarted   EventAction = "started"
	EventTypeCompleted EventAction = "completed"
	EventTypeAbandoned EventAction = "abandoned"
	EventTypeAdded     EventAction = "added"
)

var (
	// Metadata for FlatEvent table.
	flatEventMeta = itable.Metadata{
		M: &table.Metadata{
			Name: "flat_events",
			Columns: []string{
				"version",
				"id",
				"parent_id",
				"provider",
				"scope",
				"type",
				"source",
				"subject_id",
				"subject_name",
				"data",
				"created_at",
				"updated_at",
			},
			PartKey: []string{"id"},
		},
	}

	// Table instance for FlatEvent.
	flatEventTable = itable.New(*flatEventMeta.M)
)

// GetTable returns the table metadata for FlatEvent.
func (f *FlatEvent[P]) GetTable() itable.ITable {
	return flatEventTable
}

func (f *FlatEvent[P]) PreCreate() error { return nil }
func (f *FlatEvent[P]) PreUpdate() error { return nil }

func (f *FlatEvent[P]) ToEvent() (*Event[any, P], error) {
	var payload any

	switch f.Scope {
	case EventScopeBranch, EventScopeTag:
		payload = new(BranchOrTag)
	case EventScopePullRequest:
		payload = new(PullRequest)
	case EventScopePush:
		payload = new(Push)
	case EventScopePullRequestLabel:
		payload = new(PullRequestLabel)
	case EventScopePullRequestReview:
		payload = new(PullRequestReview)
	default:
		return &Event[any, P]{}, fmt.Errorf("unsupported event scope: %s", f.Scope)
	}

	err := json.Unmarshal(f.Data, payload)
	if err != nil {
		return &Event[any, P]{}, err
	}

	event := &Event[any, P]{
		Version: f.Version,
		Context: EventContext[P]{
			ID:        f.ID,
			ParentID:  f.ParentID,
			Provider:  f.Provider,
			Scope:     f.Scope,
			Action:    f.Action,
			Source:    f.Source,
			Timestamp: f.CreatedAt,
		},
		Subject: EventSubject{
			ID:   f.SubjectID,
			Name: f.SubjectName,
		},
		Data: payload,
	}

	return event, nil
}

// MarshalJSON customizes the JSON encoding for Event.
func (e Event[T, P]) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Version EventVersion    `json:"version"`
		Context EventContext[P] `json:"context"`
		Subject EventSubject    `json:"subject"`
		Payload T               `json:"payload"`
	}{
		Version: e.Version,
		Context: e.Context,
		Subject: e.Subject,
		Payload: e.Data,
	})
}

// UnmarshalJSON customizes the JSON decoding for Event.
func (e *Event[T, P]) UnmarshalJSON(data []byte) error {
	aux := struct {
		Version EventVersion    `json:"version"`
		Context EventContext[P] `json:"context"`
		Subject EventSubject    `json:"subject"`
		Payload json.RawMessage `json:"payload"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	e.Version = aux.Version
	e.Context = aux.Context
	e.Subject = aux.Subject

	var payload T

	switch e.Context.Scope {
	case EventScopeBranch, EventScopeTag:
		payload = any(new(BranchOrTag)).(T)
	case EventScopePullRequest:
		payload = any(new(PullRequest)).(T)
	case EventScopePush:
		payload = any(new(Push)).(T)
	case EventScopePullRequestLabel:
		payload = any(new(PullRequestLabel)).(T)
	case EventScopePullRequestReview:
		payload = any(new(PullRequestReview)).(T)
	default:
		return fmt.Errorf("unsupported event scope: %s", e.Context.Scope)
	}

	if err := json.Unmarshal(aux.Payload, &payload); err != nil {
		return err
	}

	e.Data = payload

	return nil
}

func (e *Event[T, P]) Flatten() (*FlatEvent[P], error) {
	payload, err := json.Marshal(e.Data)
	if err != nil {
		return nil, err
	}

	flat := &FlatEvent[P]{
		Version:     e.Version,
		ID:          e.Context.ID,
		ParentID:    e.Context.ParentID,
		Provider:    e.Context.Provider,
		Scope:       e.Context.Scope,
		Action:      e.Context.Action,
		Source:      e.Context.Source,
		SubjectID:   e.Subject.ID,
		SubjectName: e.Subject.Name,
		Data:        payload,
		CreatedAt:   e.Context.Timestamp,
		UpdatedAt:   e.Context.Timestamp,
	}

	return flat, nil
}

// SetSource sets the source field of the EventContext for the Event struct.
//
// The src parameter specifies the source of the event, such as the name of the repository.
func (e *Event[T, P]) SetSource(src string) {
	e.Context.Source = src
}

// SetParentID sets the parentID field of the EventContext for the Event struct.
//
// The id parameter specifies the parent ID of the event, which can be used to trace the event chain.
func (e *Event[T, P]) SetParentID(id gocql.UUID) {
	e.Context.ParentID = id
}

// ToEvent converts a PullRequest struct to an Event struct.
//
// The provider parameter specifies the source of the event, such as "github" or "gitlab".
// The subject parameter specifies the subject of the event, such as "pull_requests".
// The action parameter specifies the type of the action, such as "updated", "closed", or "merged".
//
// The ID field of the EventContext is set to a new TimeUUID.
// The Scope field of the EventContext is set to "pull_request".
// The Timestamp field of the EventContext is set to the UpdatedAt field of the PullRequest struct.
//
// The method returns a pointer to the Event struct that is created.
func (pr *PullRequest) ToEvent(provider RepoProvider, subject EventSubject, action EventAction) *Event[PullRequest, RepoProvider] {
	event := &Event[PullRequest, RepoProvider]{
		Version: EventVersionV1,
		Context: EventContext[RepoProvider]{
			ID:        gocql.TimeUUID(),
			Provider:  provider,
			Scope:     EventScopePullRequest,
			Action:    action,
			Timestamp: pr.UpdatedAt,
		},
		Subject: subject,
		Data:    *pr,
	}

	return event
}
