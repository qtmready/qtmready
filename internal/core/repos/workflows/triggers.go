package reposwfs

import (
	"github.com/google/uuid"
)

type (
	BranchTriggers map[string]uuid.UUID
)

// add associates a branch with its triggering event ID.
func (b BranchTriggers) add(branch string, id uuid.UUID) {
	b[branch] = id
}

// clear removes the association between a branch and its triggering event ID.
func (b BranchTriggers) clear(branch string) {
	delete(b, branch)
}

// get retrieves the event ID associated with a branch.
//
// Returns the event ID and a boolean indicating whether the branch exists.
func (b BranchTriggers) get(branch string) (uuid.UUID, bool) {
	id, ok := b[branch]

	return id, ok
}
