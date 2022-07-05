package workflows

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/ctrlplane/internal/temporal/activities"
	"go.breu.io/ctrlplane/internal/types"
)

// Workflow for handling a Github App Installation event.
func OnGithubInstall(ctx workflow.Context, payload types.GithubInstallationEventPayload) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, options)
	logger := workflow.GetLogger(ctx)

	logger.Debug("Starting Workflow: OnGithubInstall")

	var result types.GithubInstallationEventPayload
	err := workflow.ExecuteActivity(ctx, activities.SaveGithubInstallation, payload).Get(ctx, &result)

	if err != nil {
		return err
	}

	logger.Debug("Finished Workflow: OnGithubInstall")
	return nil
}

func OnGithubPullRequest(ctx workflow.Context) {}
