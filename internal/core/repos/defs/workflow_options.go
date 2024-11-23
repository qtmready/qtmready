package defs

import (
	"go.breu.io/durex/workflows"

	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/durable"
)

// RepoWorkflowOptions returns workflow options for RepoCtrl, designed for use with the Core Queue.
// The workflow ID, when used with the Core Queue, is formatted as:
//
//	"ai.ctrlplane.core.org.{org}.repos.{id}.name.{name}"
func RepoWorkflowOptions(repo *entities.Repo) workflows.Options {
	opts := durable.NewWorkflowOptions(
		durable.WithOrg(repo.OrgID.String()),
		durable.WithSubject("repos"),
		durable.WithSubjectID(repo.ID.String()),
		durable.WithMeta("name", repo.Name),
	)

	return opts
}
