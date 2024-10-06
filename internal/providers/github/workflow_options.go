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

package github

import (
	"go.breu.io/durex/workflows"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/db"
)

// RepoWebhookWorkflowOptions returns workflow options for repository events from GitHub webhooks.
//
// When used with Providers Queue, the resulting workflow ID will be prefixed with "ai.ctrlplane.providers".
// The workflow ID produced by this function will be in the format:
//
//	"github.installation.{installation_id}.repo.{repo_name}.action.{action}.id.{event_id}"
//
// installation is the GitHub App installation ID,
// repo is the name of the repository,
// event is the type of event (e.g., push, pull_request),
// and event_id is the unique ID of the event.
//
// # Example
//
//	queues.Providers().
//	  ExecuteWorkflow(
//	    ctx,
//	    RepoWebhookWorkflowOptions(installation, repo, event, event_id),
//	    w.OnInstallationEvent,
//	    payload,
//	  )
func RepoWebhookWorkflowOptions(installation db.Int64, repo, action, event_id string) workflows.Options {
	opts, _ := workflows.NewOptions(
		workflows.WithBlock("github"),
		workflows.WithElement("installation"),
		workflows.WithElementID(installation.String()),
		workflows.WithMod("repo"),
		workflows.WithModID(repo),
		workflows.WithProp("action", action),
		workflows.WithProp("event_id", event_id),
	)

	return opts
}

// InstallationWebhookWorkflowOptions returns workflow options for installation events from GitHub webhooks.
//
// When used with Providers Queue, the resulting workflow ID will be prefixed with "ai.ctrlplane.providers".
// The workflow ID produced by this function will be in the format:
//
//	"github.installation.{installation_id}.action.{action}"
//
// installation is the GitHub App installation ID
// and action is the type of action (e.g., created, deleted).
//
// # Example
//
//	queues.Providers().
//	  ExecuteWorkflow(
//	    ctx,
//	    InstallationWebhookWorkflowOptions(installation, action),
//	    w.OnInstallationEvent,
//	    payload,
//	  )
func InstallationWebhookWorkflowOptions(installation db.Int64, action string) workflows.Options {
	opts, _ := workflows.NewOptions(
		workflows.WithBlock("github"),
		workflows.WithElement("installation"),
		workflows.WithElementID(installation.String()),
		workflows.WithMod("action"),
		workflows.WithModID(action),
	)

	return opts
}

func PrepareRepoEventChildWorkflowOptions(ctx workflow.Context) workflows.Options {
	opts, _ := workflows.NewOptions(
		workflows.WithParent(ctx),
		workflows.WithBlock("prepare"),
	)

	return opts
}
