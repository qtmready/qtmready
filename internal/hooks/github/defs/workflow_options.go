package defs

import (
	"strings"

	"go.breu.io/quantm/internal/durable"
	githubv1 "go.breu.io/quantm/internal/proto/hooks/github/v1"
	"go.breu.io/quantm/internal/utils"
)

// NewInstallWorkflowOptions standardize the workflow options for Install Workflow.
//
//	io.ctrlplane.hooks.github.install.${installation_id}.{action}
//
// for exmple, for an installation id of 1234 and action of "CREATED", the resulting
// workflow options would be:
//
//	io.ctrlplane.hooks.github.install.1234.created
func NewInstallWorkflowOptions(install_id int64, action githubv1.SetupAction) *durable.WorkflowOptions {
	return durable.NewWorkflowOptions(
		durable.WithHook("github"),
		durable.WithSubject("install"),
		durable.WithSubjectID(utils.Int64ToString(install_id)),
		durable.WithAction(strings.ToLower(action.String())),
	)
}

// NewSyncReposWorkflows standardize the workflow options for InstallRepos Workflow.
//
//	io.ctrlplane.hooks.github.install.${installation_id}.repos.${action}.${action_id}
//
// for exmple, for an installation id of 1234 and action of "CREATED", the resulting
// workflow options would be:
//
//	io.ctrlplane.hooks.github.install.1234.repos.created.abcdef123-4567-8901-2345-678901234567
func NewSyncReposWorkflows(install_id int64, action string, action_id string) *durable.WorkflowOptions {
	return durable.NewWorkflowOptions(
		durable.WithHook("github"),
		durable.WithSubject("install"),
		durable.WithSubjectID(utils.Int64ToString(install_id)),
		durable.WithScope("sync-repos"),
		durable.WithAction(strings.ToLower(action)),
		durable.WithActionID(action_id),
	)
}

// NewPushWorkflowOptions standardize the workflow options for Push Workflow.
//
//	io.ctrlplane.hooks.github.repo.${repo_name}.${install_id}.${action}.${action_id}
func NewPushWorkflowOptions(repo_id int64, repo, event_id string) *durable.WorkflowOptions {
	return durable.NewWorkflowOptions(
		durable.WithHook("github"),
		durable.WithSubject("repo"),
		durable.WithSubjectID(repo),
		durable.WithScopeID(utils.Int64ToString(repo_id)),
		durable.WithAction("push"),
		durable.WithActionID(event_id),
	)
}

// NewCreateOrDeleteWorkflowOptions standardize the workflow options for CreateOrDelete Workflow.
//
//	io.ctrlplane.hooks.github.${install_id}.repo.${repo_name}.${action}.${event_id}
func NewCreateOrDeleteWorkflowOptions(install_id int64, name, action, event_id string) *durable.WorkflowOptions {
	return durable.NewWorkflowOptions(
		durable.WithHook("github"),
		durable.WithSubject("repo"),
		durable.WithSubjectID(name),
		durable.WithScopeID(utils.Int64ToString(install_id)),
		durable.WithAction(action),
		durable.WithActionID(event_id),
	)
}

// NewRefWorkflowOptions generates a workflow ID for various GitHub repository events
// based on a Git ref, a specified scope, and action.
//
// This function is designed to create workflow IDs that adhere to the
// following format:
//
//	io.ctrlplane.hooks.github.repo.${repo_id}.${ref}.${scope}.${scope_id}.${action}.${event_id}
//
// Where:
//
//   - repo_id is the GitHub repository ID.
//   - ref is the full Git ref that was pushed or involved in the event (e.g., "refs/heads/main").
//   - scope represents the broader context of the event (e.g., "push", "pull_request").
//   - scope_id is the specific ID associated with the scope (e.g., `after` SHA for push, pull request number for pull request).
//   - action is the action performed within the event (e.g., "created", "updated", "deleted").
//   - event_id is the unique UUID associated with the event.
//
// # Push Event, Deleted
//
// For a push event, the workflow ID would be constructed using the following parameters:
//
//   - repo_id: 683467348
//   - ref: refs/heads/di
//   - after: deadbeef
//   - action: deleted
//   - event_id: abcdef123-4567-8901-2345-678901234567
//
// This would result in the following Workflow ID:
//
//	io.ctrlplane.hooks.github.repo.683467348.refs/heads/di.push.deadbeef.deleted.abcdef123-4567-8901-2345-678901234567
//
// # Pull Request Event: Pull Request Opened
//
// For a pull request event (specifically, a pull request being opened), the workflow ID would be constructed with these parameters:
//
//   - repo_id: 683467348
//   - ref: refs/heads/di
//   - pr: 123
//   - action: opened
//   - event_id: abcdef123-4567-8901-2345-678901234567
//
// This would result in the following Workflow ID:
//
//	io.ctrlplane.hooks.github.repo.683467348.refs/heads/di.pull_request.123.opened.abcdef123-4567-8901-2345-678901234567
func NewRefWorkflowOptions(repo_id int64, ref, scope, scope_id, action, event_id string) *durable.WorkflowOptions {
	return durable.NewWorkflowOptions(
		durable.WithHook("github"),
		durable.WithSubject("repo"),
		durable.WithSubjectID(utils.Int64ToString(repo_id)),
		durable.WithKind(ref),
		durable.WithScope(scope),
		durable.WithScopeID(scope_id),
		durable.WithAction(action),
		durable.WithActionID(event_id),
	)
}
