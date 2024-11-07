package events

type (
	// RepoHook defines model for RepoHook.
	RepoHook string

	// MessageHook defines model for MessageHook.
	MessageHook string
)

// String returns the string representation of the RepoHook.
func (rh RepoHook) String() string { return string(rh) }

// String returns the string representation of the MessageHook.
func (mh MessageHook) String() string { return string(mh) }
