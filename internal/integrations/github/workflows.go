package github

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type Workflows struct{}

func (w *Workflows) OnInstallationEvent(ctx workflow.Context, payload InstallationEventPayload) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)

	logger.Debug("Starting Workflow: OnGithubInstall")

	var a *Activities
	var result InstallationEventPayload
	err := workflow.ExecuteActivity(ctx, a.SaveInstallation, payload).Get(ctx, &result)

	if err != nil {
		return err
	}

	return nil
}
