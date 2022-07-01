package workflows

import (
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	gh "github.com/google/go-github/v45/github"
	twf "go.temporal.io/sdk/workflow"

	"go.breu.io/ctrlplane/internal/conf"
	"go.breu.io/ctrlplane/internal/temporal/activities"
	"go.breu.io/ctrlplane/internal/temporal/common"
)

func OnGithubInstall(ctx twf.Context, payload common.GithubInstallationEventPayload) error {
	logger := twf.GetLogger(ctx)
	logger.Info("Github installation event received")
	logger.Info("Creating Github Client")

	client, err := createGithubClient(payload.Installation.ID)

	if err != nil {
		return err
	}

	twf.ExecuteActivity(ctx, activities.GetOrCreateGithubInstallation, client, payload)

	return nil
}

func OnGithubPullRequest(ctx twf.Context) {}

func createGithubClient(installationID int64) (*gh.Client, error) {
	transport, err := ghinstallation.New(http.DefaultTransport, conf.Github.AppID, installationID, conf.Github.PrivateKey)

	if err != nil {
		return nil, err
	}

	client := gh.NewClient(&http.Client{Transport: transport})
	return client, nil
}
