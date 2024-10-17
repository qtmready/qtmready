package events

import (
	"slices"

	"go.breu.io/quantm/internal/db"
)

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

// types validates against the EventScope and EventAction.
func (e Event[T, P]) types(scope EventScope, actions ...EventAction) error {
	if e.Context.Scope != scope {
		return NewScopeError(e.Context.Scope)
	}

	if !slices.Contains(actions, e.Context.Action) {
		return NewActionError(e.Context.Action)
	}

	return nil
}
