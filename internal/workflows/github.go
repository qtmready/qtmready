package workflows

import (
	"go.breu.io/ctrlplane/internal/models"
	_workflow "go.temporal.io/sdk/workflow"
)

func OnGithubInstall(ctx _workflow.Context, payload models.GithubInstallationEventPayload) error {
	return nil
}
