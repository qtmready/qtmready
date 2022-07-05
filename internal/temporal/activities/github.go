package activities

import (
	"context"

	"go.breu.io/ctrlplane/internal/db/models"
	"go.breu.io/ctrlplane/internal/types"
)

func SaveGithubInstallation(ctx context.Context, payload types.GithubInstallationEventPayload) (models.GithubInstallation, error) {
	gi := models.GithubInstallation{}
	gi.GithubInstallationID = payload.Installation.ID
	gi.GithubInstallationLogin = payload.Installation.Account.Login
	gi.GithubInstallationType = payload.Installation.Account.Type
	gi.GithubSenderID = payload.Sender.ID
	gi.GithubSenderLogin = payload.Sender.Login

	if err := gi.Create(); err != nil {
		return gi, err
	}

	return gi, nil
}
