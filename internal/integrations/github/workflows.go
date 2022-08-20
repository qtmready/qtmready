package github

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type Workflows struct{}

var activities *Activities

// OnInstall is a workflow that is executed when an installation is created.
// NOTE: This workflow will only partially update the database. We would need to handle the complete event
// from Github to assign the team id.
func (w *Workflows) OnInstall(ctx workflow.Context, payload InstallationEventPayload) error {
	opts := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, opts)
	logger := workflow.GetLogger(ctx)

	logger.Debug("Starting Workflow: OnGithubInstall")

	var result InstallationEventPayload
	err := workflow.
		ExecuteActivity(ctx, activities.GetOrCreateInstallation, payload).
		Get(ctx, &result)

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
