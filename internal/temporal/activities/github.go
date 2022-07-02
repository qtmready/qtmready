package activities

import (
	"context"

	"github.com/google/go-github/v45/github"

	"go.breu.io/ctrlplane/internal/types"
)

func SaveGithubInstallation(ctx context.Context, client *github.Client, payload types.GithubInstallationEventPayload) (types.GithubInstallationEventPayload, error) {
	return payload, nil
}
