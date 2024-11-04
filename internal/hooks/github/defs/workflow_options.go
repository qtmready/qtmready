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

// NewPushWorkflowOptions standardize the workflow options for Push Workflow.
//
//	io.ctrlplane.hooks.github.repo.${repo_name}.push.${installation_id}.${action}.${actionID}
func NewPushWorkflowOptions(id int64, repo, action, event_id string) *durable.WorkflowOptions {
	return durable.NewWorkflowOptions(
		durable.WithHook("github"),
		durable.WithSubject("repo"),
		durable.WithSubjectID(repo),
		durable.WithScope("push"),
		durable.WithScopeID(cast.Int64ToString(id)),
		durable.WithAction(action),
		durable.WithActionID(event_id),
	)
}
