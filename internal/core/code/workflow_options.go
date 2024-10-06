// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package code

import (
	"github.com/gocql/gocql"
	"go.breu.io/durex/workflows"
)

// RepoCtrlWorkflowOptions returns workflow options for RepoCtrl, designed for use with the Core Queue.
// The workflow ID, when used with the Core Queue, is formatted as:
//
//	"ai.ctrlplane.core.team.{team}.repo.{repo}.id.{repo_id}"
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
//	"ai.ctrlplane.core.team.{team}.repo.{repo}.id.{repo_id}.trunk"
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
//	"ai.ctrlplane.core.team.{team}.repo.{repo}.id.{repo_id}.branch.{branch}"
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
//	"ai.ctrlplane.core.team.{team}.repo.{repo}.id.{repo_id}.queue"
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
