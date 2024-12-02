package events

type (
	Scope string // EventScope is the scope of the event.
)

// String returns the string representation of the EventScope.
func (es Scope) String() string { return string(es) }

const (
	ScopeBranch Scope = "branch" // ScopeBranch scopes branch event.
	ScopeTag    Scope = "tag"    // ScopeTag scopes tag event.
	ScopePush   Scope = "push"   // ScopePush scopes push event.
	ScopeRebase Scope = "rebase" // ScopeRebase scopes rebase event.
	ScopeDiff   Scope = "diff"   // ScopeRebase scopes diff event.
	ScopePr     Scope = "pr"     // ScopeRebase scopes pull request event.
)
