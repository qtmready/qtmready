package reposdefs

import (
	"github.com/google/uuid"
	"go.breu.io/durex/workflows"
)

// RepoWorkflowOptions returns workflow options for RepoCtrl, designed for use with the Core Queue.
// The workflow ID, when used with the Core Queue, is formatted as:
//
//	"ai.ctrlplane.core.org.{org}.repo.{repo}.id.{repo_id}"
func RepoWorkflowOptions(org, repo string, id uuid.UUID) workflows.Options {
	opts, _ := workflows.NewOptions(
		workflows.WithBlock("org"),
		workflows.WithBlockID(org),
		workflows.WithElement("repo"),
		workflows.WithElementID(repo),
		workflows.WithMod("repo_id"),
		workflows.WithModID(id.String()),
	)

	return opts
}
