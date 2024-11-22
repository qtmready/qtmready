package workflows

import (
	"go.breu.io/durex/dispatch"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/db/entities"
	githubacts "go.breu.io/quantm/internal/hooks/github/activities"
	githubdefs "go.breu.io/quantm/internal/hooks/github/defs"
)

func SyncRepos(ctx workflow.Context, payload *githubdefs.WebhookInstallRepos) error {
	selector := workflow.NewSelector(ctx)
	acts := &githubacts.InstallRepos{}
	total := make([]string, len(payload.RepositoriesAdded)+len(payload.RepositoriesRemoved))
	install := &entities.GithubInstallation{}

	ctx = dispatch.WithDefaultActivityContext(ctx)

	if err := workflow.
		ExecuteActivity(ctx, acts.GetInstallationForSync, payload.Installation.ID).
		Get(ctx, install); err != nil {
		return err
	}

	for _, repo := range payload.RepositoriesAdded {
		payload := &githubdefs.SyncRepo{InstallationID: install.ID, Repo: repo, OrgID: install.OrgID}

		selector.AddFuture(workflow.ExecuteActivity(ctx, acts.RepoAdded, payload), func(f workflow.Future) {})
	}

	for _, repo := range payload.RepositoriesRemoved {
		payload := &githubdefs.SyncRepo{InstallationID: install.ID, Repo: repo}

		selector.AddFuture(workflow.ExecuteActivity(ctx, acts.RepoRemoved, payload), func(f workflow.Future) {})
	}

	for range total {
		selector.Select(ctx)
	}

	return nil
}
