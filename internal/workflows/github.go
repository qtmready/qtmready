package workflows

import (
	tworkflow "go.temporal.io/sdk/workflow"

	"go.breu.io/ctrlplane/internal/models"
)

func OnGithubInstall(ctx tworkflow.Context, payload models.GithubInstallationEventPayload) error {
	logger := tworkflow.GetLogger(ctx)
	logger.Info("Github installation event received")
	logger.Info("Installation: %v", payload.Installation)
	return nil
}
