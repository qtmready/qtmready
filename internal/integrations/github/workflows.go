package github

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// Workflow for handling a Github App Installation event.
func OnGithubInstallWorkflow(ctx workflow.Context, payload GithubInstallationEventPayload) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, options)
	logger := workflow.GetLogger(ctx)

	logger.Debug("Starting Workflow: OnGithubInstall")

	var result GithubInstallationEventPayload
	err := workflow.ExecuteActivity(ctx, SaveGithubInstallationActivity, payload).Get(ctx, &result)

	if err != nil {
		return err
	}

	return nil
}

func OnGithubPullRequest(ctx workflow.Context) {}
