package github

import (
	"go.breu.io/durex/workflows"

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
func RepoWebhookWorkflowOptions(installation db.Int64, repo, event, event_id string) workflows.Options {
	opts, _ := workflows.NewOptions(
		workflows.WithBlock("github"),
		workflows.WithElement("installation"),
		workflows.WithElementID(installation.String()),
		workflows.WithMod("repo"),
		workflows.WithModID(repo),
		workflows.WithProp("event", event),
		workflows.WithProp("eventid", event_id),
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
