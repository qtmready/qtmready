package events

import (
	"encoding/json"
)

type (
	EventAction string // EventAction is the action of the event.
)

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

// String returns the string representation of the EventAction.
func (ea EventAction) String() string {
	return string(ea)
}

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
	EventActionRequested EventAction = "requested" // EventActionRequested is the action of a requested event.
)
