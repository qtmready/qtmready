package activities

import (
	"context"

	gh "github.com/google/go-github/v45/github"
	tac "go.temporal.io/sdk/activity"

	"go.breu.io/ctrlplane/internal/temporal/common"
)

func GetOrCreateGithubInstallation(ctx context.Context, client *gh.Client, payload common.GithubInstallationEventPayload) error {
	logger := tac.GetLogger(ctx)
	logger.Info("Github installation event received")
	return nil
}
