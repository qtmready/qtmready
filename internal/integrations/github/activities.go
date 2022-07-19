package github

import (
	"context"

	"github.com/gocql/gocql"
	"go.breu.io/ctrlplane/internal/db/models"
)

type Activity struct{}

func (a *Activity) SaveInstallation(ctx context.Context, payload InstallationEventPayload) (models.GithubInstallation, error) {
	installation := models.GithubInstallation{}
	installation.InstallationID = payload.Installation.ID
	installation.InstallationLogin = payload.Installation.Account.Login
	installation.InstallationType = payload.Installation.Account.Type
	installation.SenderID = payload.Sender.ID
	installation.SenderLogin = payload.Sender.Login

	if err := installation.Create(); err != nil {
		return installation, err
	}

	return installation, nil
}

func (a *Activity) GetInstallation(ctx context.Context, id string) (models.GithubInstallation, error) {
	installation := models.GithubInstallation{}
	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		return installation, err
	}

	params := models.GithubInstallation{ID: uuid}
	installation.ID = uuid
	if err = installation.Get(params); err != nil {
		return installation, err
	}

	return installation, nil
}

func (a *Activity) CloneRepo(ctx context.Context, payload PushEventPayload) error {
	return nil
}
