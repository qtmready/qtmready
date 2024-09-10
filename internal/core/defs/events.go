// Copyright © 2024, Breu, Inc. <info@breu.io>
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

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

	EventScope   string
	EventAction  string
	EventVersion string

	EventProvider interface {
		RepoProvider | MessageProvider
	}

	EventPayload interface {
		BranchOrTag | PullRequest | Push | PullRequestReview | PullRequestLabel
	}

	// EventContext is the context of an event.
	EventContext[P EventProvider] struct {
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
		ID     gocql.UUID `json:"id"`
		Name   string     `json:"name"`
		TeamID gocql.UUID `json:"team_id"`
	}

	Event[T EventPayload, P EventProvider] struct {
		Version EventVersion    `json:"version"`
		Context EventContext[P] `json:"context"`
		Subject EventSubject    `json:"subject"`
		Data    T               `json:"data"`
	}

	FlatEvent[P EventProvider] struct {
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
		TeamID      gocql.UUID   `json:"team_id" cql:"team_id"`
		CreatedAt   time.Time    `json:"created_at" cql:"created_at"`
		UpdatedAt   time.Time    `json:"updated_at" cql:"updated_at"`
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

	EventActionCreated   EventAction = "created"
	EventActionDeleted   EventAction = "deleted"
	EventActionUpdated   EventAction = "updated"
	EventActionClosed    EventAction = "closed"
	EventActionMerged    EventAction = "merged"
	EventActionStarted   EventAction = "started"
	EventActionCompleted EventAction = "completed"
	EventActionAbandoned EventAction = "abandoned"
	EventActionAdded     EventAction = "added"
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

func (es EventScope) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(es))
}

func (es *EventScope) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*es = EventScope(s)

	return nil
}

func (es EventScope) String() string {
	return string(es)
}

func (ea EventAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(ea))
}

func (ea *EventAction) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*ea = EventAction(s)

	return nil
}

func (ev EventVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(ev))
}

func (ev *EventVersion) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*ev = EventVersion(s)

	return nil
}

func (ev EventVersion) String() string {
	return string(ev)
}

func (ea EventAction) String() string {
	return string(ea)
}

// MarshalJSON customizes the JSON encoding for Event.
func (e Event[T, P]) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Version EventVersion    `json:"version"`
		Context EventContext[P] `json:"context"`
		Subject EventSubject    `json:"subject"`
		Data    T               `json:"data"`
	}{
		Version: e.Version,
		Context: e.Context,
		Subject: e.Subject,
		Data:    e.Data,
	})
}

// UnmarshalJSON customizes the JSON decoding for Event.
func (e *Event[T, P]) UnmarshalJSON(data []byte) error {
	aux := struct {
		Version EventVersion    `json:"version"`
		Context EventContext[P] `json:"context"`
		Subject EventSubject    `json:"subject"`
		Data    json.RawMessage `json:"data"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	e.Version = aux.Version
	e.Context = aux.Context
	e.Subject = aux.Subject

	var data_ T

	switch e.Context.Scope {
	case EventScopeBranch, EventScopeTag:
		data_ = any(BranchOrTag{}).(T)
	case EventScopePullRequest:
		data_ = any(PullRequest{}).(T)
	case EventScopePush:
		data_ = any(Push{}).(T)
	case EventScopePullRequestLabel:
		data_ = any(PullRequestLabel{}).(T)
	case EventScopePullRequestReview:
		data_ = any(PullRequestReview{}).(T)
	default:
		return fmt.Errorf("unsupported event scope: %s", e.Context.Scope)
	}

	if err := json.Unmarshal(aux.Data, &data_); err != nil {
		return err
	}

	e.Data = data_

	return nil
}

func (e *Event[T, P]) Flatten() (*FlatEvent[P], error) {
	data, err := json.Marshal(e.Data)
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
		Data:        data,
		TeamID:      e.Subject.TeamID,
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

// SetParent sets the parentID field of the EventContext for the Event struct.
//
// The id parameter specifies the parent ID of the event, which can be used to trace the event chain.
func (e *Event[T, P]) SetParent(id gocql.UUID) {
	e.Context.ParentID = id
}

// ToEvent converts a BranchOrTag struct to an Event struct.
//
// The provider parameter specifies the source of the event, such as "github" or "gitlab".
// The subject parameter specifies the subject of the event, such as "branches" or "tags".
// The action parameter specifies the type of the action, such as "created", "deleted", or "updated".
//
// The ID field of the EventContext is set to a new TimeUUID.
// The Scope field of the EventContext is set to "branch" or "tag".
// The Timestamp field of the EventContext is set to the current time.
//
// The method returns a pointer to the Event struct that is created.
func (bt *BranchOrTag) ToEvent(
	provider RepoProvider, subject EventSubject, action EventAction, scope EventScope,
) *Event[BranchOrTag, RepoProvider] {
	event := &Event[BranchOrTag, RepoProvider]{
		Version: EventVersionV1,
		Context: EventContext[RepoProvider]{
			ID:        gocql.TimeUUID(),
			Provider:  provider,
			Scope:     scope,
			Action:    action,
			Timestamp: time.Now(),
		},
		Subject: subject,
		Data:    *bt,
	}

	return event
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
func (pr *PullRequest) ToEvent(
	provider RepoProvider, subject EventSubject, action EventAction,
) *Event[PullRequest, RepoProvider] {
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

// ToEvent converts a Push struct to an Event struct.
//
// The provider parameter specifies the source of the event, such as "github" or "gitlab".
// The subject parameter specifies the subject of the event, such as "pushes".
// The action parameter specifies the type of the action, such as "created" or "updated".
//
// The ID field of the EventContext is set to a new TimeUUID.
// The Scope field of the EventContext is set to "push".
// The Timestamp field of the EventContext is set to the Timestamp field of the Push struct.
//
// The method returns a pointer to the Event struct that is created.
func (p *Push) ToEvent(provider RepoProvider, subject EventSubject, action EventAction) *Event[Push, RepoProvider] {
	event := &Event[Push, RepoProvider]{
		Version: EventVersionV1,
		Context: EventContext[RepoProvider]{
			ID:        gocql.TimeUUID(),
			Provider:  provider,
			Scope:     EventScopePush,
			Action:    action,
			Timestamp: p.Timestamp,
		},
		Subject: subject,
		Data:    *p,
	}

	return event
}

// ToEvent converts a PullRequestReview struct to an Event struct.
//
// The provider parameter specifies the source of the event, such as "github" or "gitlab".
// The subject parameter specifies the subject of the event, such as "pull_request_reviews".
// The action parameter specifies the type of the action, such as "submitted" or "edited".
//
// The ID field of the EventContext is set to a new TimeUUID.
// The Scope field of the EventContext is set to "pull_request_review".
// The Timestamp field of the EventContext is set to the SubmittedAt field of the PullRequestReview struct.
//
// The method returns a pointer to the Event struct that is created.
func (prr *PullRequestReview) ToEvent(
	provider RepoProvider, subject EventSubject, action EventAction,
) *Event[PullRequestReview, RepoProvider] {
	event := &Event[PullRequestReview, RepoProvider]{
		Version: EventVersionV1,
		Context: EventContext[RepoProvider]{
			ID:        gocql.TimeUUID(),
			Provider:  provider,
			Scope:     EventScopePullRequestReview,
			Action:    action,
			Timestamp: prr.SubmittedAt,
		},
		Subject: subject,
		Data:    *prr,
	}

	return event
}

// ToEvent converts a PullRequestLabel struct to an Event struct.
//
// The provider parameter specifies the source of the event, such as "github" or "gitlab".
// The subject parameter specifies the subject of the event, such as "pull_request_labels".
// The action parameter specifies the type of the action, such as "added" or "removed".
//
// The ID field of the EventContext is set to a new TimeUUID.
// The Scope field of the EventContext is set to "pull_request_label".
// The Timestamp field of the EventContext is set to the UpdatedAt field of the PullRequestLabel struct.
//
// The method returns a pointer to the Event struct that is created.
func (prl *PullRequestLabel) ToEvent(
	provider RepoProvider, subject EventSubject, action EventAction,
) *Event[PullRequestLabel, RepoProvider] {
	event := &Event[PullRequestLabel, RepoProvider]{
		Version: EventVersionV1,
		Context: EventContext[RepoProvider]{
			ID:        gocql.TimeUUID(),
			Provider:  provider,
			Scope:     EventScopePullRequestLabel,
			Action:    action,
			Timestamp: prl.UpdatedAt,
		},
		Subject: subject,
		Data:    *prl,
	}

	return event
}

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

	default:
		return fmt.Errorf("unsupported event type: %T", event)
	}

	if !valid {
		return fmt.Errorf("mismatch between event type and scope: %T and scope %s", event, flat.Scope)
	}

	var payload T

	err := json.Unmarshal(flat.Data, &payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	event.Version = flat.Version
	event.Context = EventContext[P]{
		ID:        flat.ID,
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
	}
	event.Data = payload

	return nil
}
