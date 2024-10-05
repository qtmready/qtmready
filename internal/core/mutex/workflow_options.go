package mutex

import (
	"go.breu.io/durex/workflows"
)

// MutexWorkflowOptions returns workflow options for mutex operations. When used with the mutex queue, the resulting
// workflow ID will be
//
//	"ai.ctrlplane.mutex.resource.{resource_id}"
func MutexWorkflowOptions(resource_id string) workflows.Options {
	opts, _ := workflows.NewOptions(
		workflows.WithBlock("resource"),
		workflows.WithBlockID(resource_id),
	)

	return opts
}
