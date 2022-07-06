package activities

import (
	"context"

	"go.breu.io/ctrlplane/internal/db/models"
	"go.breu.io/ctrlplane/internal/types"
)

func SaveGithubInstallation(ctx context.Context, payload types.GithubInstallationEventPayload) (models.GithubInstallation, error) {
	gi := models.GithubInstallation{}
	gi.InstallationID = payload.Installation.ID
	gi.InstallationLogin = payload.Installation.Account.Login
	gi.InstallationType = payload.Installation.Account.Type
	gi.SenderID = payload.Sender.ID
	gi.SenderLogin = payload.Sender.Login

	if err := gi.Create(); err != nil {
		return gi, err
	}

	return gi, nil
}
