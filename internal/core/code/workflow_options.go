package code

import (
	"github.com/gocql/gocql"
	"go.breu.io/durex/workflows"
)

func RepoCtrlWorkflowOptions(team, repo string, id gocql.UUID) workflows.Options {
	opts, _ := workflows.NewOptions(
		workflows.WithBlock("team"),
		workflows.WithBlockID(team),
		workflows.WithElement("repo"),
		workflows.WithElementID(repo),
		workflows.WithMod("repo_id"),
		workflows.WithModID(id.String()),
	)

	return opts
}
