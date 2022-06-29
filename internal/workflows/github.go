package workflows

import (
	"go.breu.io/ctrlplane/internal/models"
	_twrkflow "go.temporal.io/sdk/workflow"
)

func OnGithubInstall(ctx _twrkflow.Context, payload models.GithubInstallationEventPayload) error {
	logger := _twrkflow.GetLogger(ctx)
	logger.Info("Github installation event received")
	logger.Info("Installation: %v", payload.Installation)
	return nil
}
