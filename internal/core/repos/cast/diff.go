package cast

import (
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

// PushEventToDiffEvent converts a Push event to a diff event.
func PushEventToDiffEvent(
	push *events.Event[eventsv1.RepoHook, eventsv1.Push],
	hook int32,
	payload *eventsv1.Diff,
) *events.Event[eventsv1.ChatHook, eventsv1.Diff] {
	return events.NextWithHook[eventsv1.RepoHook, eventsv1.ChatHook, eventsv1.Push, eventsv1.Diff](
		push,
		eventsv1.ChatHook(hook),
		events.ScopeDiff,
		events.ActionRequested,
	).SetPayload(payload)
}
