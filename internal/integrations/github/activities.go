package github

import (
	"context"

	"go.breu.io/ctrlplane/internal/entities"
	"go.temporal.io/sdk/activity"
	"go.uber.org/zap"

	"go.breu.io/ctrlplane/internal/db"
)

type Activities struct{}

func (a *Activities) GetOrCreateInstallation(ctx context.Context, payload *entities.GithubInstallation) (*entities.GithubInstallation, error) {
	log := activity.GetLogger(ctx)
	installation, err := a.GetInstallation(ctx, payload.InstallationID)

	// if we get the installation, the error will be nil
	if err == nil {
		log.Info("installation found, updating status", zap.Any("installation", installation))
		installation.Status = payload.Status
	} else {
		log.Info("installation not found, preparing create", zap.Any("installation", payload))
		installation = payload
	}

	log.Info("saving installation", zap.Any("installation", installation))
	if err := db.Save(installation); err != nil {
		log.Error("error saving installation", zap.Error(err))
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
