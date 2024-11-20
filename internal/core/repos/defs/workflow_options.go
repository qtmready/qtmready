package reposdefs

import (
	"github.com/google/uuid"
	"go.breu.io/durex/workflows"

	"go.breu.io/quantm/internal/durable"
)

// RepoWorkflowOptions returns workflow options for RepoCtrl, designed for use with the Core Queue.
// The workflow ID, when used with the Core Queue, is formatted as:
//
//	"ai.ctrlplane.core.org.{org}.repos.{id}.name.{name}"
func RepoWorkflowOptions(org, id uuid.UUID, name string) workflows.Options {
	opts := durable.NewWorkflowOptions(
		durable.WithOrg(org.String()),
		durable.WithSubject("repos"),
		durable.WithSubjectID(id.String()),
		durable.WithMeta("name", name),
	)

	return opts
}
