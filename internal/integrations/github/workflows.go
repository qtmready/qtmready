package github

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type Workflows struct{}

var activity *Activity

func (w *Workflows) OnInstall(ctx workflow.Context, payload InstallationEventPayload) error {
	opts := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, opts)
	logger := workflow.GetLogger(ctx)

	logger.Debug("Starting Workflow: OnGithubInstall")

	var result InstallationEventPayload
	err := workflow.ExecuteActivity(ctx, activity.SaveInstallation, payload).Get(ctx, &result)

	if err != nil {
		return err
	}

	return nil
}

func (w *Workflows) OnPush(ctx workflow.Context, payload PushEventPayload) error {
	return nil
}

func (w *Workflows) OnPullRequest(ctx workflow.Context, payload PullRequestEventPayload) error {
	return nil
}
