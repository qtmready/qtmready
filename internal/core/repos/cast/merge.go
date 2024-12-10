package cast

import (
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

// MergeConflictEvent converts a Push event to a diff event.
func RebaseEventToMergeConflictEvent(
	rebase *events.Event[eventsv1.RepoHook, eventsv1.Rebase],
	hook int32,
	payload *eventsv1.Merge,
) *events.Event[eventsv1.ChatHook, eventsv1.Merge] {
	return events.NextWithHook[eventsv1.RepoHook, eventsv1.ChatHook, eventsv1.Rebase, eventsv1.Merge](
		rebase,
		eventsv1.ChatHook(hook),
		events.ScopeMerge,
		events.ActionMerge,
	).SetPayload(payload)
}
