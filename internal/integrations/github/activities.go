package github

import (
	"context"

	"go.breu.io/ctrlplane/internal/db/models"
)

type Activities struct{}

func (w *Activities) SaveInstallation(ctx context.Context, payload InstallationEventPayload) (models.GithubInstallation, error) {
	g := models.GithubInstallation{}
	g.InstallationID = payload.Installation.ID
	g.InstallationLogin = payload.Installation.Account.Login
	g.InstallationType = payload.Installation.Account.Type
	g.SenderID = payload.Sender.ID
	g.SenderLogin = payload.Sender.Login

	if err := g.Create(); err != nil {
		return g, err
	}

	return g, nil
}
