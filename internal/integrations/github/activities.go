package github

import (
	"context"
	"strconv"

	"go.breu.io/ctrlplane/internal/entities"
	"go.temporal.io/sdk/activity"

	"go.breu.io/ctrlplane/internal/db"
)

type Activities struct{}

// CreateOrUpdateInstallation creates or update the entities.GithubInstallation
func (a *Activities) CreateOrUpdateInstallation(ctx context.Context, payload *entities.GithubInstallation) (*entities.GithubInstallation, error) {
	log := activity.GetLogger(ctx)
	installation, err := a.GetInstallation(ctx, payload.InstallationID)

	// if we get the installation, the error will be nil
	if err == nil {
		log.Info("installation found, updating status ...", "installation", installation)
		installation.Status = payload.Status
	} else {
		log.Info("installation not found, preparing create ...", "payload", payload)
		installation = payload
	}

	log.Info("saving installation", "installation", installation)
	if err := db.Save(installation); err != nil {
		log.Error("error saving installation", "error", err)
		return installation, err
	}

	return installation, nil
}

// CreateRepo creates a single row for entities.GithubRepo
func (a *Activities) CreateRepo(ctx context.Context, payload *entities.GithubRepo) error {
	log := activity.GetLogger(ctx)

	log.Info("saving repository", "repository", payload)
	if err := db.Save(payload); err != nil {
		log.Error("error saving repository", "error", err)
		return err
	}

	return nil
}

// GetInstallation gets entities.GithubInstallation against given installation_id
func (a *Activities) GetInstallation(ctx context.Context, id int64) (*entities.GithubInstallation, error) {
	installation := &entities.GithubInstallation{}

	if err := db.Get(installation, db.QueryParams{"installation_id": strconv.FormatInt(id, 10)}); err != nil {
		return installation, err
	}

	return installation, nil
}

func (a *Activities) CloneRepo(ctx context.Context, payload PushEventPayload) error {
	return nil
}
