// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package github

import (
	"context"
	"strconv"

	"go.temporal.io/sdk/activity"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/entity"
)

type (
	// Activities groups all the activities for the github provider.
	Activities struct{}
)

// CreateOrUpdateInstallation creates or update the entity.GithubInstallation.
func (a *Activities) CreateOrUpdateInstallation(ctx context.Context, payload *entity.GithubInstallation) (*entity.GithubInstallation, error) {
	log := activity.GetLogger(ctx)
	installation, err := a.GetInstallation(ctx, payload.InstallationID)

	// if we get the installation, the error will be nil
	if err == nil {
		log.Info("installation found, updating status ...")

		installation.Status = payload.Status
	} else {
		log.Info("installation not found, creating ...")
		log.Debug("payload", "payload", payload)

		installation = payload
	}

	log.Info("saving installation ...")

	if err := db.Save(installation); err != nil {
		log.Error("error saving installation", "error", err)
		return installation, err
	}

	return installation, nil
}

// CreateOrUpdateGithubRepo creates a single row for entity.GithubRepo.
func (a *Activities) CreateOrUpdateGithubRepo(ctx context.Context, payload *entity.GithubRepo) error {
	log := activity.GetLogger(ctx)
	repo, err := a.GetGithubRepo(ctx, payload)

	// if we get the repo, the error will be nil
	if err == nil {
		log.Info("repository found, updating ...")
	} else {
		log.Info("repository not found, creating ...")
		log.Debug("payload", "payload", payload)
	}

	if err := db.Save(repo); err != nil {
		log.Error("error saving repository ...", "error", err)
		return err
	}

	return nil
}

// GetGithubRepo gets entity.GithubRepo against given entity.GithubRepo.
func (a *Activities) GetGithubRepo(ctx context.Context, payload *entity.GithubRepo) (*entity.GithubRepo, error) {
	repo := &entity.GithubRepo{}
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

// GetInstallation gets entity.GithubInstallation against given installation_id.
func (a *Activities) GetInstallation(ctx context.Context, id int64) (*entity.GithubInstallation, error) {
	installation := &entity.GithubInstallation{}

	if err := db.Get(installation, db.QueryParams{"installation_id": strconv.FormatInt(id, 10)}); err != nil {
		return installation, err
	}

	return installation, nil
}

func (a *Activities) GetRepo(ctx context.Context, repo *entity.GithubRepo) (*entity.Repo, error) {
	r := &entity.Repo{}

	if err := db.Get(r, db.QueryParams{"github_id": strconv.FormatInt(repo.GithubID, 10)}); err != nil {
		return r, err
	}

	return r, nil
}
