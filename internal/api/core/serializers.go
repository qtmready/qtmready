// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
//
// This software is made available by Breu, Inc., under the terms of the Breu  
// Community License Agreement, Version 1.0 located at  
// http://www.breu.io/breu-community-license/v1. BY INSTALLING, DOWNLOADING,  
// ACCESSING, USING OR DISTRIBUTING ANY OF THE SOFTWARE, YOU AGREE TO THE TERMS  
// OF SUCH LICENSE AGREEMENT. 

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
