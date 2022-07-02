package workflows

import (
	"net/http"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v45/github"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/ctrlplane/internal/conf"
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

	client, err := createGithubClient(payload.Installation.ID)

	if err != nil {
		return err
	}

	var result types.GithubInstallationEventPayload
	err = workflow.ExecuteActivity(ctx, activities.SaveGithubInstallation, client, payload).Get(ctx, &result)

	if err != nil {
		return err
	}

	logger.Debug("Finished Workflow: OnGithubInstall")
	return nil
}

func OnGithubPullRequest(ctx workflow.Context) {}

// Creates a github client for the given installation ID.

func createGithubClient(installationID int64) (*github.Client, error) {
	transport, err := ghinstallation.New(
		http.DefaultTransport,
		conf.Github.AppID,
		installationID,
		[]byte(conf.Github.PrivateKey),
	)

	if err != nil {
		return nil, err
	}

	client := github.NewClient(&http.Client{Transport: transport})
	return client, nil
}
