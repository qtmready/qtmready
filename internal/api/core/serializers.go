// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
//
// This software is made available by Breu, Inc., under the terms of the 
// Breu Community License Agreement ("BCL Agreement"), version 1.0, found at  
// https://www.breu.io/license/community. By installating, downloading, 
// accessing, using or distrubting any of the software, you agree to the  
// terms of the license agreement. 
//
// The above copyright notice and the subsequent license agreement shall be 
// included in all copies or substantial portions of the software. 
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, 
// IMPLIED, STATUTORY, OR OTHERWISE, AND SPECIFICALLY DISCLAIMS ANY WARRANTY OF 
// MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE 
// SOFTWARE. 
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT 
// LIMITED TO, LOST PROFITS OR ANY CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, 
// OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, ARISING 
// OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY  
// APPLICABLE LAW. 

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
