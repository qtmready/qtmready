package githubdefs

import (
	"strings"

	"go.breu.io/quantm/internal/cast"
	"go.breu.io/quantm/internal/durable"
	githubv1 "go.breu.io/quantm/internal/proto/hooks/github/v1"
)

// NewInstallWorkflowOptions standardize the workflow options for Install Workflow.
//
//	io.ctrlplane.hooks.github.install.${installation_id}
func NewInstallWorkflowOptions(id int64, action githubv1.SetupAction) *durable.WorkflowOptions {
	return durable.NewWorkflowOptions(
		durable.WithHook("github"),
		durable.WithSubject("install"),
		durable.WithSubjectID(cast.Int64ToString(id)),
		durable.WithAction(strings.ToLower(action.String())),
	)
}
