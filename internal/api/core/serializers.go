// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package core

import (
	"github.com/gocql/gocql"

	"go.breu.io/ctrlplane/internal/entities"
)

type (
	AppCreateRequest struct {
		Name   string             `json:"name"`
		Config entities.AppConfig `json:"config"`
	}

	AppRepoCreateRequest struct {
		RepoID        gocql.UUID `json:"repo_id"`
		DefaultBranch string     `json:"default_branch"`
		IsMonorepo    bool       `json:"is_monorepo"`
		Provider      string     `json:"provider"`
	}
)
