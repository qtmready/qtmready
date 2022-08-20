package github

import (
	"context"

	"go.breu.io/ctrlplane/internal/cmn"
	"go.breu.io/ctrlplane/internal/entities"

	"go.breu.io/ctrlplane/internal/db"
)

type Activities struct{}

func (a *Activities) GetOrCreateInstallation(ctx context.Context, payload InstallationEventPayload) (*entities.GithubInstallation, error) {
	installation, err := a.GetInstallation(ctx, payload.Installation.ID)
	// if we get the installation, the error will be nil
	if err == nil {
		return installation, nil
	}

	installation.InstallationLogin = payload.Installation.Account.Login
	installation.InstallationType = payload.Installation.Account.Type
	installation.SenderID = payload.Sender.ID
	installation.SenderLogin = payload.Sender.Login

	if err := cmn.Validate.Struct(installation); err != nil {
		return installation, err
	}

	if err := db.Save(installation); err != nil {
		return installation, err
	}

	return installation, nil
}

func (a *Activities) GetInstallation(ctx context.Context, id int64) (*entities.GithubInstallation, error) {
	installation := &entities.GithubInstallation{}

	if err := db.Get(installation, db.QueryParams{"installation_id": id}); err != nil {
		return installation, err
	}

	return installation, nil
}

func (a *Activities) CloneRepo(ctx context.Context, payload PushEventPayload) error {
	return nil
}
