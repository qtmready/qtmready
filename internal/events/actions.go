package events

type (
	// Action is the action of the event.
	Action string
)

const (
	ActionCreated    Action = "created"   // EventActionCreated is the action of a created event.
	ActionDeleted    Action = "deleted"   // EventActionDeleted is the action of a deleted event.
	ActionUpdated    Action = "updated"   // EventActionUpdated is the action of an updated event.
	ActionForced     Action = "forced"    // ActionForced is the action of a forced event.
	ActionReopened   Action = "reopened"  // ActionReopened is the action of a reopened event.
	ActionClosed     Action = "closed"    // ActionClosed is the action of a closed event.
	ActionMerged     Action = "merged"    // ActionMerged is the action of a merged event.
	ActionStarted    Action = "started"   // ActionStarted is the action of a started event.
	ActionCompleted  Action = "completed" // ActionCompleted is the action of a completed event.
	ActionDismissed  Action = "dismissed" // ActionDismissed is the action of a dismissed event.
	ActionAbandoned  Action = "abandoned" // ActionAbandoned is the action of an abandoned event.
	EventActionAdded Action = "added"     // EventActionAdded is the action of an added event.
	ActionRequested  Action = "requested" // EventActionRequested is the action of a requested event.
	ActionDiff       Action = "diff"      // ActionDiff is the action of a diff change pushed event.
	ActionMerge      Action = "merge"     // ActionMerge is the action of a merge change pr event.
)

// String returns the string representation of the EventAction.
func (ea Action) String() string { return string(ea) }
