package code // Package name updated to "code"

import (
	"github.com/gocql/gocql"
)

type (
	// BranchTriggers maintains a mapping of branch names to the event IDs that triggered their creation.
	// This data structure allows us to determine the root event for any given branch, enabling event
	// lineage tracing.
	//
	// Example:
	//   events := make(BranchTriggers)
	//   events.Add("branch", event.ID)
	BranchTriggers map[string]gocql.UUID
)

// add associates a new branch with its corresponding event ID.
func (b BranchTriggers) add(branch string, id gocql.UUID) {
	b[branch] = id
}

// del disassociates a branch from its stored event ID.
func (b BranchTriggers) del(branch string) {
	delete(b, branch)
}

// get retrieves the event ID associated with the specified branch.
//
// Example:
//
//	branch, ok := events.get("branch")
func (b BranchTriggers) get(branch string) (gocql.UUID, bool) {
	id, ok := b[branch]

	return id, ok
}
