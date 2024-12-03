package cast

import (
	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

// PushEventToDiffEvent converts a Push event to a diff event.
// TODO - refine or make a Next func which convet hook with payload.
func PushEventToDiffEvent(
	push *events.Event[eventsv1.RepoHook, eventsv1.Push],
	payload *eventsv1.Diff,
) *events.Event[eventsv1.ChatHook, eventsv1.Diff] {
	return events.NextHook[eventsv1.RepoHook, eventsv1.ChatHook, eventsv1.Push, eventsv1.Diff](
		push,
		events.ScopeDiff,
		events.ActionDiff,
		eventsv1.ChatHook_CHAT_HOOK_SLACK,
	).SetSubjectName(events.SubjectNameChat).SetPayload(payload)
}
