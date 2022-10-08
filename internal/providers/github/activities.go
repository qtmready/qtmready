// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package github

import (
	"context"
	"strconv"

	"go.breu.io/ctrlplane/internal/entities"
	"go.temporal.io/sdk/activity"

	"go.breu.io/ctrlplane/internal/db"
)

type (
	Activities struct{}
)

// CreateOrUpdateInstallation creates or update the entities.GithubInstallation.
func (a *Activities) CreateOrUpdateInstallation(ctx context.Context, payload *entities.GithubInstallation) (*entities.GithubInstallation, error) {
	log := activity.GetLogger(ctx)
	installation, err := a.GetInstallation(ctx, payload.InstallationID)

	// if we get the installation, the error will be nil
	if err == nil {
		log.Info("installation found, updating status ...")

		installation.Status = payload.Status
	} else {
		log.Info("installation not found, creating ...", "payload", payload)

		installation = payload
	}

	log.Info("saving installation ...", "installation", installation)

	if err := db.Save(installation); err != nil {
		log.Error("error saving installation", "error", err)
		return installation, err
	}

	return installation, nil
}

// CreateOrUpdateRepo creates a single row for entities.GithubRepo.
func (a *Activities) CreateOrUpdateRepo(ctx context.Context, payload *entities.GithubRepo) error {
	log := activity.GetLogger(ctx)
	repo, err := a.GetRepo(ctx, payload)

	// if we get the repo, the error will be nil
	if err == nil {
		log.Info("repository found, updating ...")
	} else {
		log.Info("repository not found, creating ...", "payload", payload)
	}

	if err := db.Save(repo); err != nil {
		log.Error("error saving repository", "error", err)
		return err
	}

	return nil
}

// GetRepo gets entities.GithubRepo against given entities.GithubRepo.
func (a *Activities) GetRepo(ctx context.Context, payload *entities.GithubRepo) (*entities.GithubRepo, error) {
	repo := &entities.GithubRepo{}
	params := db.QueryParams{
		"name":      "'" + payload.Name + "'",
		"full_name": "'" + payload.FullName + "'",
		"github_id": strconv.FormatInt(payload.GithubID, 10),
		"team_id":   payload.TeamID.String(),
	}

	if err := db.Get(repo, params); err != nil {
		return payload, err
	}

	return repo, nil
}

// GetInstallation gets entities.GithubInstallation against given installation_id.
func (a *Activities) GetInstallation(ctx context.Context, id int64) (*entities.GithubInstallation, error) {
	installation := &entities.GithubInstallation{}

	if err := db.Get(installation, db.QueryParams{"installation_id": strconv.FormatInt(id, 10)}); err != nil {
		return installation, err
	}

	return installation, nil
}
