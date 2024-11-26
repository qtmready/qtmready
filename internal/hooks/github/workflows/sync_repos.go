package workflows

import (
	"go.breu.io/durex/dispatch"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/db/entities"
	"go.breu.io/quantm/internal/hooks/github/activities"
	"go.breu.io/quantm/internal/hooks/github/defs"
)

// // SyncRepos synchronizes repositories for a given installation.
//
// This workflow handles the addition and removal of repositories
// from a GitHub installation. It retrieves the installation details,
// then iterates through the added and removed repositories, executing
// activities to handle the synchronization process for each repository.
//
// The workflow uses a selector to manage concurrent execution of
// activities for added and removed repositories. It waits for all
// activities to complete before returning.
func SyncRepos(ctx workflow.Context, payload *defs.WebhookInstallRepos) error {
	selector := workflow.NewSelector(ctx)
	acts := &activities.InstallRepos{}
	total := make([]string, len(payload.RepositoriesAdded)+len(payload.RepositoriesRemoved))
	install := &entities.GithubInstallation{}

	ctx = dispatch.WithDefaultActivityContext(ctx)

	if err := workflow.
		ExecuteActivity(ctx, acts.GetInstallationForSync, payload.Installation.ID).
		Get(ctx, install); err != nil {
		return err
	}

	for _, repo := range payload.RepositoriesAdded {
		payload := &defs.SyncRepoPayload{InstallationID: install.ID, Repo: repo, OrgID: install.OrgID}

		selector.AddFuture(workflow.ExecuteActivity(ctx, acts.RepoAdded, payload), func(f workflow.Future) {})
	}

	for _, repo := range payload.RepositoriesRemoved {
		payload := &defs.SyncRepoPayload{InstallationID: install.ID, Repo: repo}

		selector.AddFuture(workflow.ExecuteActivity(ctx, acts.RepoRemoved, payload), func(f workflow.Future) {})
	}

	for range total {
		selector.Select(ctx)
	}

	return nil
}
