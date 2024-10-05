package code

import (
	"github.com/gocql/gocql"
	"go.breu.io/durex/workflows"
)

// RepoCtrlWorkflowOptions returns workflow options for RepoCtrl, designed for use with the Core Queue.
// The workflow ID, when used with the Core Queue, is formatted as:
//
//	ai.ctrlplane.core.team.{team}.repo.{repo}.id.{repo_id}
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

// TrunkCtrlWorkflowOptions returns workflow options for TrunkCtrl, tailored for trunk-related workflows.
// The workflow ID, when used with the Core Queue, is formatted as:
//
//	ai.ctrlplane.core.team.{team}.repo.{repo}.id.{repo_id}.trunk
func TrunkCtrlWorkflowOptions(team, repo string, id gocql.UUID) workflows.Options {
	opts, _ := workflows.NewOptions(
		workflows.WithBlock("team"),
		workflows.WithBlockID(team),
		workflows.WithElement("repo"),
		workflows.WithElementID(repo),
		workflows.WithMod("repo_id"),
		workflows.WithModID(id.String()),
		workflows.WithProp("trunk", ""),
	)

	return opts
}

// BranchCtrlWorkflowOptions returns workflow options for BranchCtrl, specifying a branch.
// The workflow ID, when used with the Core Queue, is formatted as:
//
//	ai.ctrlplane.core.team.{team}.repo.{repo}.id.{repo_id}.branch.{branch}
func BranchCtrlWorkflowOptions(team, repo string, id gocql.UUID, branch string) workflows.Options {
	opts, _ := workflows.NewOptions(
		workflows.WithBlock("team"),
		workflows.WithBlockID(team),
		workflows.WithElement("repo"),
		workflows.WithElementID(repo),
		workflows.WithMod("repo_id"),
		workflows.WithModID(id.String()),
		workflows.WithProp("branch", branch),
	)

	return opts
}

// QueueCtrlWorkflowOptions returns workflow options for QueueCtrl, used for queue-related workflows.
// The workflow ID, when used with the Core Queue, is formatted as:
//
//	ai.ctrlplane.core.team.{team}.repo.{repo}.id.{repo_id}.queue
func QueueCtrlWorkflowOptions(team, repo string, id gocql.UUID) workflows.Options {
	opts, _ := workflows.NewOptions(
		workflows.WithBlock("team"),
		workflows.WithBlockID(team),
		workflows.WithElement("repo"),
		workflows.WithElementID(repo),
		workflows.WithMod("repo_id"),
		workflows.WithModID(id.String()),
		workflows.WithProp("queue", ""),
	)

	return opts
}
