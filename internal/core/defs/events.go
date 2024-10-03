// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
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

	"go.breu.io/quantm/internal/db"
)

// Event paylaods.
type (
	// BranchOrTag represents a git branch or tag.
	BranchOrTag struct {
		Ref           string `json:"ref"`            // Ref is the name of the branch or tag.
		DefaultBranch string `json:"default_branch"` // DefaultBranch is the name of the default branch.
	}

	// TODO - set the event payload.
	LineChanges struct {
		Added     db.Int64 `json:"added"`     // Number of lines added in the commit.
		Removed   db.Int64 `json:"removed"`   // Number of lines removed in the commit.
		Threshold db.Int64 `json:"threshold"` // Set threshold for PR.
		Delta     db.Int64 `json:"delta"`     // Net change in lines (added - removed).
	}

	// Commit represents a git commit.
	Commit struct {
		SHA       string    `json:"sha"`       // SHA is the SHA of the commit.
		Message   string    `json:"message"`   // Message is the commit message.
		URL       string    `json:"url"`       // URL is the URL of the commit.
		Added     []string  `json:"added"`     // Added is a list of files added in the commit.
		Removed   []string  `json:"removed"`   // Removed is a list of files removed in the commit.
		Modified  []string  `json:"modified"`  // Modified is a list of files modified in the commit.
		Author    string    `json:"author"`    // Author is the author of the commit.
		Committer string    `json:"committer"` // Committer is the committer of the commit.
		Timestamp time.Time `json:"timestamp"` // Timestamp is the timestamp of the commit.
	}

	Commits []Commit

	// Push represents a git push.
	Push struct {
		Ref        string    `json:"ref"`        // Ref is the ref that was pushed to.
		Before     string    `json:"before"`     // Before is the SHA of the commit before the push.
		After      string    `json:"after"`      // After is the SHA of the commit after the push.
		Repository string    `json:"repository"` // Repository is the repository that was pushed to.
		SenderID   db.Int64  `json:"sender_id"`  // SenderID is the id of the user who pushed the changes.
		Commits    Commits   `json:"commits"`    // Commits is a list of commits that were pushed.
		Timestamp  time.Time `json:"timestamp"`  // Timestamp is the timestamp of the push.
	}

	// PullRequest represents a pull request.
	PullRequest struct {
		Number         db.Int64  `json:"number"`                     // Number is the pull request number.
		Title          string    `json:"title"`                      // Title is the pull request title.
		Body           string    `json:"body"`                       // Body is the pull request body.
		State          string    `json:"state"`                      // State is the pull request state.
		MergeCommitSHA *string   `json:"merge_commit_sha,omitempty"` // MergeCommitSHA is the SHA of the merge commit.
		AuthorID       db.Int64  `json:"author_id"`                  // AuthorID is the author_id of the pull request.
		HeadBranch     string    `json:"head_branch"`                // HeadBranch is the head branch of the pull request.
		BaseBranch     string    `json:"base_branch"`                // BaseBranch is the base branch of the pull request.
		Timestamp      time.Time `json:"timestamp"`                  // Timestamp is the timestamp when the pull request was created.
	}

	// PullRequestReview represents a pull request review.
	PullRequestReview struct {
		ID                db.Int64  `json:"id"`                  // ID is the pull request review ID.
		PullRequestNumber db.Int64  `json:"pull_request_number"` // PullRequestNumber is the pull request number.
		Branch            string    `json:"branch"`              // Branch is the branch the review belongs to.
		State             string    `json:"state"`               // State is the pull request review state.
		AuthorID          db.Int64  `json:"author_id"`           // AuthorID is the author of the review.
		Timestamp         time.Time `json:"submitted_at"`        // SubmittedAt is the timestamp when the review was submitted.
	}

	// PullRequestLabel represents a pull request label.
	PullRequestLabel struct {
		Name              string    `json:"name"`                // Name is the text of the label e.g. "ready", "fix" etc.
		PullRequestNumber db.Int64  `json:"pull_request_number"` // PullRequestNumber is the pull request number.
		Branch            string    `json:"branch"`              // Branch is the branch the label belongs to.
		Timestamp         time.Time `json:"timestamp"`           // Timestamp is the timestamp of the label.
	}

	// PullRequestComment represents a pull request comment.
	PullRequestComment struct {
		ID                db.Int64  `json:"id"`                    // ID is the pull request review comment ID.
		PullRequestNumber db.Int64  `json:"pull_request_number"`   // PullRequestNumber is the pull request number.
		Branch            string    `json:"branch"`                // Branch is the branch the comment belongs to.
		ReviewID          db.Int64  `json:"review_id"`             // ReviewID is the ID of the pull request review the comment belongs.
		InReplyTo         *db.Int64 `json:"in_reply_to,omitempty"` // InReplyTo is the ID of the parent comment.
		CommitSHA         string    `json:"commit_sha"`            // CommitSHA is the SHA of the commit associated with the comment.
		Path              string    `json:"path"`                  // Path is the path to the file where the comment was made.
		Position          db.Int64  `json:"position"`              // Position is the line number where the comment was made.
		AuthorID          db.Int64  `json:"author_id"`             // AuthorID is the author_id of the comment.
		Timestamp         time.Time `json:"timestamp"`             // Timestamp is the timestamp of the comment.
	}

	// PullRequestThread represents a pull request thread.
	PullRequestThread struct {
		ID                db.Int64   `json:"id"`                  // ID is the pull request thread ID.
		PullRequestNumber db.Int64   `json:"pull_request_number"` // PullRequestNumber is the pull request number.
		CommentIDs        []db.Int64 `json:"comment_ids"`         // CommentIDs is the list of comment IDs associated with the thread.
		Timestamp         time.Time  `json:"timestamp"`           // Timestamp is the timestamp of the thread.
	}

	// MergeConflict represents a git merge conflict.
	MergeConflict struct {
		HeadBranch string    `json:"head_branch"` // HeadBranch is the name of the head branch.
		HeadCommit Commit    `json:"head_commit"` // HeadCommit is the last commit on the head branch before rebasing.
		BaseBranch string    `json:"base_branch"` // BaseBranch is the name of the base branch.
		BaseCommit Commit    `json:"base_commit"` // BaseCommit is the last commit on the base branch before rebasing.
		Files      []string  `json:"files"`       // Files is the list of files with conflicts.
		Timestamp  time.Time `json:"timestamp"`   // Timestamp is the timestamp of the merge conflict.
	}

	LinesExceed struct {
		Branch    string      `json:"branch"`     // Branch is the name of the head or feature branch.
		Commit    Commit      `json:"commit"`     // Commit is the last commit on the head branch.
		LineStats LineChanges `json:"line_stats"` // LineStats contains details about lines added, removed, and the delta.
		Timestamp time.Time   `json:"timestamp"`  // Timestamp is the timestamp of the merge conflict.
	}

	// EventPayload represents all available event payloads.
	EventPayload interface {
		BranchOrTag |
			Push |
			PullRequest | PullRequestReview | PullRequestLabel | PullRequestComment | PullRequestThread |
			MergeConflict | LinesExceed
	}
)

// Event context.
type (
	EventVersion string // EventVersion represents the version of an event.
	EventScope   string // EventScope represents the scope of an event.
	EventAction  string // EventAction represents the action of an event.

	// EventProvider represents the origin of an event.
	EventProvider interface {
		RepoProvider | MessageProvider
	}

	// EventContext represents the contextual information surrounding an event.
	// This context is crucial for understanding and processing the event.
	EventContext[P EventProvider] struct {
		ParentID  gocql.UUID  `json:"parent_id"` // ParentID is the ID of preceding related event (tracing chains).
		Provider  P           `json:"provider"`  // Provider is the Event origin (e.g., GitHub, GitLab, GCP).
		Scope     EventScope  `json:"scope"`     // Scope is the Event category (e.g., branch, pull_request).
		Action    EventAction `json:"action"`    // Action is the Triggering action (e.g., created, updated, merged).
		Source    string      `json:"source"`    // Source is the Event source.
		Timestamp time.Time   `json:"timestamp"` // Timestamp is the Event occurrence time.
	}
)

type (
	// EventSubject represents the entity within the system that is the subject of an event.
	//
	// It encapsulates:
	//   - ID: The unique identifier of the entity i.e. the primary key within its respective database table.
	//   - Name: The name of the entity's corresponding database table. This provides context for the event's subject. For instance, an
	//     event related to a branch would have "repos" as the subject name, as branches are associated with repositories.
	//   - TeamID: The unique identifier of the team to which this entity belongs. This allows for team-based filtering and organizatio
	//     of events.
	EventSubject struct {
		ID     gocql.UUID `json:"id"`      // ID is the ID of the subject.
		Name   string     `json:"name"`    // Name of the database table.
		TeamID gocql.UUID `json:"team_id"` // TeamID is the ID of the team that the subject belongs to.
		UserID gocql.UUID `json:"user_id"` // UserID is the ID of the user that the subject belongs to. It can be null uuid.
	}
)

type (
	// Event represents an event.
	Event[T EventPayload, P EventProvider] struct {
		Version EventVersion    `json:"version"` // Version is the version of the event.
		ID      gocql.UUID      `json:"id"`      // ID is the ID of the event.
		Context EventContext[P] `json:"context"` // Context is the context of the event.
		Subject EventSubject    `json:"subject"` // Subject is the subject of the event.
		Payload T               `json:"payload"` // Payload is the payload of the event.
	}

	// FlatEvent is a flattened representation of an event.
	FlatEvent[P EventProvider] struct {
		Version  EventVersion `json:"version" cql:"version"`     // Version is the version of the event.
		ID       gocql.UUID   `json:"id" cql:"id"`               // ID is the ID of the event.
		ParentID gocql.UUID   `json:"parent_id" cql:"parent_id"` // ParentID is the ID of the parent event.
		Provider P            `json:"provider" cql:"provider"`   // Provider is the provider of the event.

		// TODO - replace the db tag to cql tag (custom mapper).
		Scope       EventScope  `json:"scope" db:"scope_"`               // Scope is the scope of the event.
		Action      EventAction `json:"action" cql:"action"`             // Action is the action of the event.
		Source      string      `json:"source" cql:"source"`             // Source is the source of the event.
		SubjectID   gocql.UUID  `json:"subject_id" cql:"subject_id"`     // SubjectID is the ID of the subject.
		SubjectName string      `json:"subject_name" cql:"subject_name"` // SubjectName is the name of the subject.
		Payload     []byte      `json:"payload" cql:"payload"`           // Payload is the payload of the event.
		TeamID      gocql.UUID  `json:"team_id" cql:"team_id"`           // TeamID is the ID of the team that the subject belongs to.
		UserID      gocql.UUID  `json:"user_id" cql:"user_id"`           // UserID is the ID of the user that the subject belongs to.
		CreatedAt   time.Time   `json:"created_at" cql:"created_at"`     // CreatedAt is the timestamp when the event was created.
		UpdatedAt   time.Time   `json:"updated_at" cql:"updated_at"`     // UpdatedAt is the timestamp when the event was last updated.
	}
)

const (
	EventVersionDefault EventVersion = "0.1.0" // EventVersionDefault is the default version of an event.
)

const (
	EventScopeBranch             EventScope = "branch"               // EventScopeBranch scopes branch event.
	EventScopeTag                EventScope = "tag"                  // EventScopeTag scopes tag event.
	EventScopeCommit             EventScope = "commit"               // EventScopeCommit scopes commit event.
	EventScopePush               EventScope = "push"                 // EventScopePush scopes push event.
	EventScopePullRequest        EventScope = "pull_request"         // EventScopePullRequest scopes PR event.
	EventScopePullRequestLabel   EventScope = "pull_request_label"   // EventScopePullRequestLabel scopes PR label event.
	EventScopePullRequestReview  EventScope = "pull_request_review"  // EventScopePullRequestReview scopes PR review event.
	EventScopePullRequestComment EventScope = "pull_request_comment" // EventScopePullRequestReviewComment scopes PR comment event.
	EventScopePullRequestThread  EventScope = "pull_request_thread"  // EventScopePullRequestThread scopes PR thread event.
	EventScopeMergeConflict      EventScope = "merge_conflict"       // EventScopeMergeCommit scopes merge commit event.
	EventScopeLineExceed         EventScope = "line_exceed"          // EventScopeLineExceed scopes line exceed event.
)

const (
	EventActionCreated   EventAction = "created"   // EventActionCreated is the action of a created event.
	EventActionDeleted   EventAction = "deleted"   // EventActionDeleted is the action of a deleted event.
	EventActionUpdated   EventAction = "updated"   // EventActionUpdated is the action of an updated event.
	EventActionForced    EventAction = "forced"    // EventActionForced is the action of a forced event.
	EventActionReopened  EventAction = "reopened"  // EventActionReopened is the action of a reopened event.
	EventActionClosed    EventAction = "closed"    // EventActionClosed is the action of a closed event.
	EventActionMerged    EventAction = "merged"    // EventActionMerged is the action of a merged event.
	EventActionStarted   EventAction = "started"   // EventActionStarted is the action of a started event.
	EventActionCompleted EventAction = "completed" // EventActionCompleted is the action of a completed event.
	EventActionDismissed EventAction = "dismissed" // EventActionDismissed is the action of a dismissed event.
	EventActionAbandoned EventAction = "abandoned" // EventActionAbandoned is the action of an abandoned event.
	EventActionAdded     EventAction = "added"     // EventActionAdded is the action of an added event.
)

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

// GetTable returns the table metadata for FlatEvent.
func (f *FlatEvent[P]) GetTable() itable.ITable {
	return flatEventTable
}

func (f *FlatEvent[P]) PreCreate() error { return nil }
func (f *FlatEvent[P]) PreUpdate() error { return nil }

// MarshalJSON customizes the JSON encoding for EventScope.
func (es EventScope) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(es))
}

// UnmarshalJSON customizes the JSON decoding for EventScope.
func (es *EventScope) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*es = EventScope(s)

	return nil
}

// String returns the string representation of the EventScope.
func (es EventScope) String() string {
	return string(es)
}

// MarshalJSON customizes the JSON encoding for EventAction.
func (ea EventAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(ea))
}

// UnmarshalJSON customizes the JSON decoding for EventAction.
func (ea *EventAction) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*ea = EventAction(s)

	return nil
}

// MarshalJSON customizes the JSON encoding for EventVersion.
func (ev EventVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(ev))
}

// UnmarshalJSON customizes the JSON decoding for EventVersion.
func (ev *EventVersion) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*ev = EventVersion(s)

	return nil
}

// String returns the string representation of the EventVersion.
func (ev EventVersion) String() string {
	return string(ev)
}

// String returns the string representation of the EventAction.
func (ea EventAction) String() string {
	return string(ea)
}

// UnmarshalJSON customizes the JSON decoding for Event.
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
	default:
		return fmt.Errorf("unsupported event scope: %s", e.Context.Scope)
	}

	if err := json.Unmarshal(aux.Payload, &payload); err != nil {
		return err
	}

	e.Payload = payload

	return nil
}

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

// SetSource sets the source field of the EventContext for the Event struct.
//
// The src parameter specifies the source of the event, such as the name of the repository.
//
// Example usage:
//
//	event := &Event[EventPayload, EventProvider]{}
//	event.SetSource("example/repo")
func (e *Event[T, P]) SetSource(src string) { e.Context.Source = src }

// SetParent sets the parentID field of the EventContext for the Event struct.
//
// The id parameter specifies the parent ID of the event, which can be used to trace the event chain.
//
// Example usage:
//
//	event := &Event[EventPayload, EventProvider]{}
//	event.SetParent(gocql.TimeUUID())
func (e *Event[T, P]) SetParent(id gocql.UUID) { e.Context.ParentID = id }

func (e *Event[T, P]) SetTimestamp(t time.Time) { e.Context.Timestamp = t }

func (e *Event[T, P]) SetUserID(id gocql.UUID) { e.Subject.UserID = id }

// SetActionCreated sets the action of the Event to EventActionCreated.
func (e *Event[T, P]) SetActionCreated() {
	e.Context.Action = EventActionCreated
}

// SetActionDeleted sets the action of the Event to EventActionDeleted.
func (e *Event[T, P]) SetActionDeleted() {
	e.Context.Action = EventActionDeleted
}

// SetActionUpdated sets the action of the Event to EventActionUpdated.
func (e *Event[T, P]) SetActionUpdated() {
	e.Context.Action = EventActionUpdated
}

// SetActionClosed sets the action of the Event to EventActionClosed.
func (e *Event[T, P]) SetActionClosed() {
	e.Context.Action = EventActionClosed
}

// SetActionMerged sets the action of the Event to EventActionMerged.
func (e *Event[T, P]) SetActionMerged() {
	e.Context.Action = EventActionMerged
}

// SetActionStarted sets the action of the Event to EventActionStarted.
func (e *Event[T, P]) SetActionStarted() {
	e.Context.Action = EventActionStarted
}

// SetActionCompleted sets the action of the Event to EventActionCompleted.
func (e *Event[T, P]) SetActionCompleted() {
	e.Context.Action = EventActionCompleted
}

// SetActionDismissed sets the action of the Event to EventActionDismissed.
func (e *Event[T, P]) SetActionDismissed() {
	e.Context.Action = EventActionDismissed
}

// SetActionAbandoned sets the action of the Event to EventActionAbandoned.
func (e *Event[T, P]) SetActionAbandoned() {
	e.Context.Action = EventActionAbandoned
}

// SetActionAdded sets the action of the Event to EventActionAdded.
func (e *Event[T, P]) SetActionAdded() {
	e.Context.Action = EventActionAdded
}

// SetScopeBranch sets the scope of the Event to EventScopeBranch.
func (e *Event[T, P]) SetScopeBranch() {
	e.Context.Scope = EventScopeBranch
}

// SetScopeTag sets the scope of the Event to EventScopeTag.
func (e *Event[T, P]) SetScopeTag() {
	e.Context.Scope = EventScopeTag
}

// SetScopePush sets the scope of the Event to EventScopePush.
func (e *Event[T, P]) SetScopePush() {
	e.Context.Scope = EventScopePush
}

// SetScopePullRequest sets the scope of the Event to EventScopePullRequest.
func (e *Event[T, P]) SetScopePullRequest() {
	e.Context.Scope = EventScopePullRequest
}

// SetScopePullRequestLabel sets the scope of the Event to EventScopePullRequestLabel.
func (e *Event[T, P]) SetScopePullRequestLabel() {
	e.Context.Scope = EventScopePullRequestLabel
}

// SetScopePullRequestReview sets the scope of the Event to EventScopePullRequestReview.
func (e *Event[T, P]) SetScopePullRequestReview() {
	e.Context.Scope = EventScopePullRequestReview
}

// SetScopePullRequestComment sets the scope of the Event to EventScopePullRequestComment.
func (e *Event[T, P]) SetScopePullRequestComment() {
	e.Context.Scope = EventScopePullRequestComment
}

// SetScopePullRequestThread sets the scope of the Event to EventScopePullRequestThread.
func (e *Event[T, P]) SetScopePullRequestThread() {
	e.Context.Scope = EventScopePullRequestThread
}

func (e *Event[T, P]) SetScopeMergeConflict() {
	e.Context.Scope = EventScopeMergeConflict
}

func (e *Event[T, P]) SetScopeLineExceed() {
	e.Context.Scope = EventScopeLineExceed
}

// Latest returns the latest commit based on the timestamp. It iterates through the Commits slice and returns the commit with the latest
// Timestamp. If the Commits slice is empty, it returns a zero-value Commit.
//
// Example:
//
//	commits := Commits{
//		{SHA: "sha1", Timestamp: time.Now().Add(-1 * time.Hour)},
//		{SHA: "sha2", Timestamp: time.Now()},
//	}
//
// latestCommit := commits.Latest()
// // latestCommit will be the commit with SHA "sha2".
func (c Commits) Latest() *Commit {
	if len(c) == 0 {
		return nil
	}

	latest := c[0]
	for _, commit := range c {
		if commit.Timestamp.After(latest.Timestamp) {
			latest = commit
		}
	}

	return &latest
}

// ToEvent converts a BranchOrTag struct to an Event struct.
//
// The provider parameter specifies the source of the event, such as "github" or "gitlab". The subject parameter specifies the subject of
// the event, such as "branches" or "tags". The action parameter specifies the type of the action, such as "created", "deleted", or
// "updated".
//
// The ID field of the EventContext is set to a new TimeUUID. The Scope field of the EventContext is set to "branch" or "tag". The Timestamp
// field of the EventContext is set to the current time.
//
// The method returns a pointer to the Event struct that is created.
//
// Example usage:
//
//	branch := &BranchOrTag{
//	  Ref: "main",
//	}
//	event := branch.ToEvent(
//	  RepoProviderGithub,
//	  EventSubject{
//	    ID:   repo_id,
//	    Name: "repos",
//	  },
//	  EventActionCreated,
//	  EventScopeBranch,
//	)
func (bt *BranchOrTag) ToEvent(
	provider RepoProvider, subject EventSubject, action EventAction, scope EventScope,
) *Event[BranchOrTag, RepoProvider] {
	event := &Event[BranchOrTag, RepoProvider]{
		Version: EventVersionDefault,
		ID:      gocql.TimeUUID(),
		Context: EventContext[RepoProvider]{
			Provider:  provider,
			Scope:     scope,
			Action:    action,
			Timestamp: time.Now(),
		},
		Subject: subject,
		Payload: *bt,
	}

	return event
}

// ToEvent converts a Push struct to an Event struct.
//
// The provider parameter specifies the source of the event, such as "github" or "gitlab". The subject parameter specifies the subject of
// the event, such as "pushes". The action parameter specifies the type of the action, such as "created" or "updated".
//
// The ID field of the EventContext is set to a new TimeUUID. The Scope field of the EventContext is set to "push". The Timestamp field of
// the EventContext is set to the Timestamp field of the Push struct.
//
// The method returns a pointer to the Event struct that is created.
//
// Example usage:
//
//	push := &Push{
//	  Ref: "main",
//	  Before: "old_sha",
//	  After: "new_sha",
//	  Repository: "example/repo",
//	  Pusher: "user",
//	  Commits: []Commit{},
//	  Timestamp: time.Now(),
//	}
//	event := push.ToEvent(
//	  RepoProviderGithub,
//	  EventSubject{
//	    ID:   repo_id,
//	    Name: "repos",
//	  },
//	  EventActionCreated,
//	)
func (p *Push) ToEvent(
	provider RepoProvider, subject EventSubject, action EventAction,
) *Event[Push, RepoProvider] {
	event := &Event[Push, RepoProvider]{
		Version: EventVersionDefault,
		ID:      gocql.TimeUUID(),
		Context: EventContext[RepoProvider]{
			Provider:  provider,
			Scope:     EventScopePush,
			Action:    action,
			Timestamp: p.Timestamp,
		},
		Subject: subject,
		Payload: *p,
	}

	return event
}

// ToEvent converts a PullRequest struct to an Event struct.
//
// The provider parameter specifies the source of the event, such as "github" or "gitlab". The subject parameter specifies the subject of
// the event, such as "pull_requests". The action parameter specifies the type of the action, such as "updated", "closed", or "merged".
//
// The ID field of the EventContext is set to a new TimeUUID. The Scope field of the EventContext is set to "pull_request". The Timestamp
// field of the EventContext is set to the UpdatedAt field of the PullRequest struct.
//
// The method returns a pointer to the Event struct that is created.
//
// Example usage:
//
//	pr := &PullRequest{
//	  Number: 1,
//	  Title:  "Test Pull Request",
//	  Body:   "This is a test pull request",
//	  State:  "open",
//	  UpdatedAt: time.Now(),
//	}
//	event := pr.ToEvent(
//	  RepoProviderGithub,
//	  EventSubject{
//	    ID:   repo_id,
//	    Name: "repos",
//	  },
//	  EventActionUpdated,
//	)
func (pr *PullRequest) ToEvent(
	provider RepoProvider, subject EventSubject, action EventAction,
) *Event[PullRequest, RepoProvider] {
	event := &Event[PullRequest, RepoProvider]{
		Version: EventVersionDefault,
		ID:      gocql.TimeUUID(),
		Context: EventContext[RepoProvider]{
			Provider:  provider,
			Scope:     EventScopePullRequest,
			Action:    action,
			Timestamp: pr.Timestamp,
		},
		Subject: subject,
		Payload: *pr,
	}

	return event
}

// ToEvent converts a PullRequestReview struct to an Event struct.
//
// The provider parameter specifies the source of the event, such as "github" or "gitlab". The subject parameter specifies the subject of
// the event, such as "pull_request_reviews". The action parameter specifies the type of the action, such as "submitted" or "edited".
//
// The ID field of the EventContext is set to a new TimeUUID. The Scope field of the EventContext is set to "pull_request_review". The
// Timestamp field of the EventContext is set to the SubmittedAt field of the PullRequestReview struct.
//
// The method returns a pointer to the Event struct that is created.
//
// Example usage:
//
//	review := &PullRequestReview{
//	  ID: 1,
//	  Body: "This is a review",
//	  State: "approved",
//	  SubmittedAt: time.Now(),
//	  Author: "user",
//	  PullRequestNumber: 1,
//	}
//	event := review.ToEvent(
//	  RepoProviderGithub,
//	  EventSubject{
//	    ID:   repo_id,
//	    Name: "repos",
//	  },
//	  EventActionSubmitted,
//	)
func (prr *PullRequestReview) ToEvent(
	provider RepoProvider, subject EventSubject, action EventAction,
) *Event[PullRequestReview, RepoProvider] {
	event := &Event[PullRequestReview, RepoProvider]{
		Version: EventVersionDefault,
		ID:      gocql.TimeUUID(),
		Context: EventContext[RepoProvider]{
			Provider:  provider,
			Scope:     EventScopePullRequestReview,
			Action:    action,
			Timestamp: prr.Timestamp,
		},
		Subject: subject,
		Payload: *prr,
	}

	return event
}

// ToEvent converts a PullRequestLabel struct to an Event struct.
//
// The provider parameter specifies the source of the event, such as "github" or "gitlab". The subject parameter specifies the subject of
// the event, such as "pull_request_labels". The action parameter specifies the type of the action, such as "added" or "removed".
//
// The ID field of the EventContext is set to a new TimeUUID. The Scope field of the EventContext is set to "pull_request_label". The
// Timestamp field of the EventContext is set to the UpdatedAt field of the PullRequestLabel struct.
//
// The method returns a pointer to the Event struct that is created.
//
// Example usage:
//
//	label := &PullRequestLabel{
//	  Name: "bug",
//	  Color: "red",
//	  Description: "This is a bug label",
//	  CreatedAt: time.Now(),
//	  UpdatedAt: time.Now(),
//	  PullRequestNumber: 1,
//	}
//	event := label.ToEvent(
//	  RepoProviderGithub,
//	  EventSubject{
//	    ID:   repo_id,
//	    Name: "repos",
//	  },
//	  EventActionAdded,
//	)
func (prl *PullRequestLabel) ToEvent(
	provider RepoProvider, subject EventSubject, action EventAction,
) *Event[PullRequestLabel, RepoProvider] {
	event := &Event[PullRequestLabel, RepoProvider]{
		Version: EventVersionDefault,
		ID:      gocql.TimeUUID(),
		Context: EventContext[RepoProvider]{
			Provider:  provider,
			Scope:     EventScopePullRequestLabel,
			Action:    action,
			Timestamp: prl.Timestamp,
		},
		Subject: subject,
		Payload: *prl,
	}

	return event
}

// ToEvent converts a PullRequestComment struct to an Event struct.
//
// The provider parameter specifies the source of the event, such as "github" or "gitlab". The subject parameter specifies the subject of
// the event, such as "pull_request_comments". The action parameter specifies the type of the action, such as "created" or "updated".
//
// The ID field of the EventContext is set to a new TimeUUID. The Scope field of the EventContext is set to "pull_request_comment". The
// Timestamp field of the EventContext is set to the Timestamp field of the PullRequestComment struct.
//
// The method returns a pointer to the Event struct that is created.
//
// Example usage:
//
//	comment := &PullRequestComment{
//	  ID: 1,
//	  Path: "path/to/file.go",
//	  Position: 15,
//	  Author: "user",
//	  PullRequestNumber: 1,
//	  ReviewID: 5,
//	  CommitSHA: "abcdef1234567890",
//	  InReplyTo: nil,
//	  Timestamp: time.Now(),
//	}
//	event := comment.ToEvent(
//	  RepoProviderGithub,
//	  EventSubject{
//	    ID:   repo_id,
//	    Name: "repos",
//	  },
//	  EventActionCreated,
//	)
func (prc *PullRequestComment) ToEvent(
	provider RepoProvider, subject EventSubject, action EventAction,
) *Event[PullRequestComment, RepoProvider] {
	event := &Event[PullRequestComment, RepoProvider]{
		Version: EventVersionDefault,
		ID:      gocql.TimeUUID(),
		Context: EventContext[RepoProvider]{
			Provider:  provider,
			Scope:     EventScopePullRequestComment,
			Action:    action,
			Timestamp: prc.Timestamp,
		},
		Subject: subject,
		Payload: *prc,
	}

	return event
}

// ToEvent converts a PullRequestThread struct to an Event struct.
//
// The provider parameter specifies the source of the event, such as "github" or "gitlab". The subject parameter specifies the subject of
// the event, such as "pull_request_threads". The action parameter specifies the type of the action, such as "created" or "updated".
//
// The ID field of the EventContext is set to a new TimeUUID. The Scope field of the EventContext is set to "pull_request_thread". The
// Timestamp field of the EventContext is set to the UpdatedAt field of the PullRequestThread struct.
//
// The method returns a pointer to the Event struct that is created.
//
// Example usage:
//
//	thread := &PullRequestThread{
//	  ID: 1,
//	  Title: "Question about implementation",
//	  Comments: []Comment{},
//	  CreatedAt: time.Now(),
//	  UpdatedAt: time.Now(),
//	  Path: "path/to/file.go",
//	  Position: 15,
//	}
//	event := thread.ToEvent(
//	  RepoProviderGithub,
//	  EventSubject{
//	    ID:   repo_id,
//	    Name: "repos",
//	  },
//	  EventActionCreated,
//	)
func (prt *PullRequestThread) ToEvent(
	provider RepoProvider, subject EventSubject, action EventAction,
) *Event[PullRequestThread, RepoProvider] {
	event := &Event[PullRequestThread, RepoProvider]{
		Version: EventVersionDefault,
		ID:      gocql.TimeUUID(),
		Context: EventContext[RepoProvider]{
			Provider:  provider,
			Scope:     EventScopePullRequestThread,
			Action:    action,
			Timestamp: prt.Timestamp,
		},
		Subject: subject,
		Payload: *prt,
	}

	return event
}

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
