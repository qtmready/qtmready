package workflows

import (
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/durable"
)

// NewInstallWorkflowOptions standardize the workflow options for Install Workflow.
//
//	io.ctrlplane.hooks.github.install.${installation_id}
func NewInstallWorkflowOptions(id db.Int64, action string) *durable.WorkflowOptions {
	return durable.NewWorkflowOptions(
		durable.WithHook("github"),
		durable.WithSubject("install"),
		durable.WithSubjectID(id.String()),
		durable.WithAction(action),
	)
}
