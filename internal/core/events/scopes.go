package events

import (
	"encoding/json"
)

type (
	EventScope string // EventScope is the scope of the event.
)

// MarshalJSON implements the json.Marshaler interface.
func (e EventScope) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(e))
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (e *EventScope) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	*e = EventScope(s)

	return nil
}

// String implements the fmt.Stringer interface.
func (e EventScope) String() string {
	return string(e)
}

const (
	EventScopeBranch             EventScope = "branch"               // EventScopeBranch scopes branch event.
	EventScopeTag                EventScope = "tag"                  // EventScopeTag scopes tag event.
	EventScopeCommit             EventScope = "commit"               // EventScopeCommit scopes commit event.
	EventScopePush               EventScope = "push"                 // EventScopePush scopes push event.
	EventScopeRebase             EventScope = "rebase"               // EventScopePush scopes push event.
	EventScopePullRequest        EventScope = "pull_request"         // EventScopePullRequest scopes PR event.
	EventScopePullRequestLabel   EventScope = "pull_request_label"   // EventScopePullRequestLabel scopes PR label event.
	EventScopePullRequestReview  EventScope = "pull_request_review"  // EventScopePullRequestReview scopes PR review event.
	EventScopePullRequestComment EventScope = "pull_request_comment" // EventScopePullRequestReviewComment scopes PR comment event.
	EventScopePullRequestThread  EventScope = "pull_request_thread"  // EventScopePullRequestThread scopes PR thread event.
	EventScopeMergeConflict      EventScope = "merge_conflict"       // EventScopeMergeCommit scopes merge commit event.
	EventScopeLineExceed         EventScope = "line_exceed"          // EventScopeLineExceed scopes line exceed event.
)
