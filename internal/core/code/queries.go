package code

import (
	"go.breu.io/durex/queues"
)

const (
	QueryRepoCtrlForBranchTriggers      queues.Query = "repo_ctrl__branch_triggers"        // get all the branches for a repo.
	QueryRepoCtrlForBranchParentEventID queues.Query = "repo_ctrl__branch_parent_event_id" // get the parent event ID for a branch.
)
