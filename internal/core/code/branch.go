package code // Package name updated to "code"

import (
	"github.com/gocql/gocql"
)

type (
	// BranchRootEvents maintains a mapping of branch names to the event IDs that triggered their creation.
	// This data structure allows us to determine the root event for any given branch, enabling event
	// lineage tracing.
	//
	// Example:
	//   events := make(BranchRootEvents)
	//   events.Add("branch", event.ID)
	BranchRootEvents map[string]gocql.UUID
)

// Add associates a new branch with its corresponding event ID.
func (b BranchRootEvents) Add(branch string, id gocql.UUID) {
	b[branch] = id
}

// Remove disassociates a branch from its stored event ID.
func (b BranchRootEvents) Remove(branch string) {
	delete(b, branch)
}

// Get retrieves the event ID associated with the specified branch.
//
// Example:
//
//	branch, ok := events.Get("branch")
func (b BranchRootEvents) Get(branch string) (gocql.UUID, bool) {
	id, ok := b[branch]

	return id, ok
}
