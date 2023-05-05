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

package core

import (
	"context"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/shared"
)

type (
	Activities     struct{}
	ResourceResult struct {
		Resources []Resource `json:"resources"`
	}
)

// GetResources gets resources from DB against a stack.
func (a *Activities) GetResources(ctx context.Context, stackID string) (*ResourceResult, error) {
	resources := make([]Resource, 0)
	params := db.QueryParams{"stack_id": stackID}

	if err := db.Filter(&Resource{}, &resources, params); err != nil {
		return &ResourceResult{Resources: resources}, err
	}

	shared.Logger.Debug("GetResources", "resources", resources)

	return &ResourceResult{Resources: resources}, nil
}

// GetWorkloads gets workloads from DB against a stack.
func (a *Activities) GetWorkloads(ctx context.Context, stackID string) ([]Workload, error) {
	wl := make([]Workload, 0)
	params := db.QueryParams{"stack_id": stackID}

	if err := db.Filter(&Workload{}, &wl, params); err != nil {
		return wl, err
	}

	shared.Logger.Debug("GetWorkloads", "workloads", wl)

	return wl, nil
}

// GetWorkloads gets workloads from DB against a stack.
func (a *Activities) GetRepos(ctx context.Context, stackID string) ([]Repo, error) {
	repos := make([]Repo, 0)
	params := db.QueryParams{"stack_id": stackID}

	if err := db.Filter(&Repo{}, &repos, params); err != nil {
		return repos, err
	}

	shared.Logger.Debug("GetRepos", "repos", repos)

	return repos, nil
}

func (a *Activities) GetBluePrint(ctx context.Context, stackID string) (*Blueprint, error) {
	blueprint := &Blueprint{}
	params := db.QueryParams{"stack_id": stackID}

	if err := db.Get(blueprint, params); err != nil {
		return blueprint, err
	}

	shared.Logger.Debug("GetBluePrint", "blueprint", blueprint)

	return blueprint, nil
}

func (a *Activities) CreateChangeset(ctx context.Context, changeSet *ChangeSet) error {
	err := db.Save(changeSet)
	return err
}
