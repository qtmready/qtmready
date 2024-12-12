package events

type (
	Scope string // EventScope is the scope of the event.
)

// String returns the string representation of the EventScope.
func (es Scope) String() string { return string(es) }

const (
	ScopeBranch     Scope = "branch"      // ScopeBranch scopes branch event.
	ScopeTag        Scope = "tag"         // ScopeTag scopes tag event.
	ScopePush       Scope = "push"        // ScopePush scopes push event.
	ScopeRebase     Scope = "rebase"      // ScopeRebase scopes rebase event.
	ScopeDiff       Scope = "diff"        // ScopeDiff scopes diff event.
	ScopePr         Scope = "pr"          // ScopePr scopes pull request event.
	ScopePrLabel    Scope = "pr_label"    // ScopePrLabel scopes pull request label event.
	ScopeMerge      Scope = "merge"       // ScopeMerge scopes merge event.
	ScopeMergeQueue Scope = "merge_queue" // ScopeMergeQueue scopes merge queue event.
)
