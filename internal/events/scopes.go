package events

type (
	EventScope string // EventScope is the scope of the event.
)

// String returns the string representation of the EventScope.
func (es EventScope) String() string { return string(es) }

const (
	EventScopeBranch EventScope = "branch" // EventScopeBranch scopes branch event.
	EventScopeTag    EventScope = "tag"    // EventScopeTag scopes tag event.
	EventScopePush   EventScope = "push"   // EventScopePush scopes push event.
)
