package github

import (
	"context"
	"go.breu.io/ctrlplane/internal/entities"

	"github.com/gocql/gocql"
	"go.breu.io/ctrlplane/internal/db"
)

type Activity struct{}

func (a *Activity) SaveInstallation(ctx context.Context, payload InstallationEventPayload) (*entities.GithubInstallation, error) {
	installation := &entities.GithubInstallation{}
	installation.InstallationID = payload.Installation.ID
	installation.InstallationLogin = payload.Installation.Account.Login
	installation.InstallationType = payload.Installation.Account.Type
	installation.SenderID = payload.Sender.ID
	installation.SenderLogin = payload.Sender.Login

	if err := db.Save(installation); err != nil {
		return installation, err
	}

	return installation, nil
}

func (a *Activity) GetInstallation(ctx context.Context, id string) (*entities.GithubInstallation, error) {
	uuid, err := gocql.ParseUUID(id)
	installation := &entities.GithubInstallation{}
	if err != nil {
		return installation, err
	}
	if err := db.Get(installation, db.QueryParams{"id": uuid}); err != nil {
		return installation, err
	}
	return installation, nil
}

func (a *Activity) CloneRepo(ctx context.Context, payload PushEventPayload) error {
	return nil
}
