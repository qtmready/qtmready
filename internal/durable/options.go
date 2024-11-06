package durable

import (
	"strings"

	"go.breu.io/durex/workflows"
)

type (
	// WorkflowOptions represents the options for a durable workflow. It extends on the interface provided by the durex
	// package to align with the naming strategy of ctrlplane.
	//
	//  - Hook: (Optional) The source of the hook (e.g., github, gitlab, etc.).
	//  - Subject: (Preferred) The primary entity the workflow operates on (e.g., repos, users, projects).
	//  - Scope: (Required) The broader context of the workflow (e.g., branch, pull_request).
	//  - Action: (Preferred) The primary operation performed (e.g., created, updated, review, merge).
	//  - Kind: (Optional) Subject kind, e.g. for repos, it will be branch name.
	//
	// The format of workflow IDs will be:
	//
	//  hooks.${hook}.${subject}.${subject_id}.${kind}.${scope}.${scope_id}.${action}.${action_id}.${...meta}
	//
	// For example, a push event on the `feature` branch of a repository with UUID `123e4567-e89b-12d3-a456-426655440000`
	// and commit SHA `deadbeef` triggered from GitHub might have an ID like:
	//
	//  hooks.github.repos.123e4567-e89b-12d3-a456-426655440000.feature.push.deadbeef.created.df94c890-980d-11ef-b864-0242ac120002
	//
	// Breakdown:
	//
	//  - hooks:  Indicates the workflow is processed through the hooks queue.
	//  - github:  Specifies the source of the webhook, which is GitHub.
	//  - repos: Specifies the subject of the workflow is a repository.
	//  - 123e4567-e89b-12d3-a456-426655440000:  The unique identifier (UUID) of the repository.
	//  - feature:  The kind of the subject, which is the branch name `feature`.
	//  - push:  The scope of the workflow, indicating the push event is happening on the branch.
	//  - deadbeef:  The scope ID, which is the commit SHA of the pushed commit.
	//  - created: The action performed, indicating a new branch was created.
	//  - df94c890-980d-11ef-b864-0242ac120002: The webhook event id from GitHub.
	//
	// Workflow IDs are constructed through helper functions specific to each workflow, rather than directly using this struct.
	WorkflowOptions struct {
		Hook      *string           `json:"hook,omitempty"`
		Subject   *string           `json:"subject,omitempty"`
		SubjectID *string           `json:"subject_id,omitempty"`
		Scope     *string           `json:"scope,omitempty"`
		ScopeID   *string           `json:"scope_id,omitempty"`
		Action    *string           `json:"action,omitempty"`
		ActionID  *string           `json:"action_id,omitempty"`
		Kind      *string           `json:"kind,omitempty"`
		Meta      map[string]string `json:"meta"`

		ParentID *string `json:"parent_id,omitempty"`

		MaximumAttempts *int32   `json:"maximum_attempts,omitempty"`
		IgnoreErrors    []string `json:"ignored_errors,omitempty"`
	}

	WorkflowOptionBuilder func(*WorkflowOptions)
)

// IsChild returns true if the workflow id is a child workflow id.
func (o *WorkflowOptions) IsChild() bool {
	return o.ParentID != nil
}

// ParentWorkflowID returns the parent workflow id.
func (o *WorkflowOptions) ParentWorkflowID() string {
	if o.ParentID == nil {
		return ""
	}

	return *o.ParentID
}

// IDSuffix returns the sanitized suffix of the workflow ID.
func (o *WorkflowOptions) IDSuffix() string {
	parts := []string{}
	if o.Hook != nil {
		parts = append(parts, *o.Hook)
	}

	// Subject
	if o.Subject != nil {
		parts = append(parts, *o.Subject)
	}

	if o.SubjectID != nil {
		parts = append(parts, *o.SubjectID)
	}

	// Kind (optional)
	if o.Kind != nil {
		parts = append(parts, *o.Kind)
	}

	// Scope
	if o.Scope != nil {
		parts = append(parts, *o.Scope)
	}

	if o.ScopeID != nil {
		parts = append(parts, *o.ScopeID)
	}

	// Action
	if o.Action != nil {
		parts = append(parts, *o.Action)
	}

	if o.ActionID != nil {
		parts = append(parts, *o.ActionID)
	}

	// Metadata
	for k, v := range o.Meta {
		parts = append(parts, k, v)
	}

	// Sanitization and joining
	sanitized := make([]string, 0)

	for _, part := range parts {
		trim := strings.TrimSpace(part)
		if trim != "" {
			sanitized = append(sanitized, trim)
		}
	}

	return strings.Join(sanitized, ".")
}

// MaxAttempts returns the max attempts for the workflow.
func (o *WorkflowOptions) MaxAttempts() int32 {
	if o.MaximumAttempts == nil {
		return workflows.RetryForever
	}

	return *o.MaximumAttempts
}

// IgnoredErrors returns the list of errors that are ok to ignore.
func (o *WorkflowOptions) IgnoredErrors() []string {
	return o.IgnoreErrors
}

// WithHook sets the hook for the workflow.
func WithHook(hook string) WorkflowOptionBuilder {
	return func(o *WorkflowOptions) {
		o.Hook = &hook
	}
}

// WithSubject sets the subject for the workflow.
func WithSubject(subject string) WorkflowOptionBuilder {
	return func(o *WorkflowOptions) {
		o.Subject = &subject
	}
}

// WithSubjectID sets the subject ID for the workflow.
func WithSubjectID(subjectID string) WorkflowOptionBuilder {
	return func(o *WorkflowOptions) {
		o.SubjectID = &subjectID
	}
}

// WithScope sets the scope for the workflow.
func WithScope(scope string) WorkflowOptionBuilder {
	return func(o *WorkflowOptions) {
		o.Scope = &scope
	}
}

// WithScopeID sets the scope ID for the workflow.
func WithScopeID(scopeID string) WorkflowOptionBuilder {
	return func(o *WorkflowOptions) {
		o.ScopeID = &scopeID
	}
}

// WithAction sets the action for the workflow.
func WithAction(action string) WorkflowOptionBuilder {
	return func(o *WorkflowOptions) {
		o.Action = &action
	}
}

// WithActionID sets the action ID for the workflow.
func WithActionID(actionID string) WorkflowOptionBuilder {
	return func(o *WorkflowOptions) {
		o.ActionID = &actionID
	}
}

// WithKind sets the kind for the workflow.
func WithKind(kind string) WorkflowOptionBuilder {
	return func(o *WorkflowOptions) {
		o.Kind = &kind
	}
}

// WithMeta sets the meta for the workflow.
func WithMeta(k, v string) WorkflowOptionBuilder {
	return func(o *WorkflowOptions) {
		if o.Meta == nil {
			o.Meta = make(map[string]string)
		}

		o.Meta[k] = v
	}
}

// WithParentID sets the parent ID for the workflow.
func WithParentID(parentID string) WorkflowOptionBuilder {
	return func(o *WorkflowOptions) {
		o.ParentID = &parentID
	}
}

// WithMaxAttempts sets the maximum attempts for the workflow.
func WithMaxAttempts(maxAttempts int32) WorkflowOptionBuilder {
	return func(o *WorkflowOptions) {
		o.MaximumAttempts = &maxAttempts
	}
}

// WithIgnoreErrors sets the ignored errors for the workflow.
func WithIgnoreErrors(ignoreErrors []string) WorkflowOptionBuilder {
	return func(o *WorkflowOptions) {
		o.IgnoreErrors = ignoreErrors
	}
}

// NewWorkflowOptions returns a new WorkflowOptions with the given options applied.
func NewWorkflowOptions(options ...WorkflowOptionBuilder) *WorkflowOptions {
	o := &WorkflowOptions{
		Meta: map[string]string{},
	}

	for _, option := range options {
		option(o)
	}

	return o
}
