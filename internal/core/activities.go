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

	"github.com/gocql/gocql"
	"go.temporal.io/sdk/activity"

	"go.breu.io/ctrlplane/internal/db"
)

var (
	activities *Activities
)

type (
	Activities struct{}
)

// GetResources gets resources from DB against a stack.
func (a *Activities) GetResources(ctx context.Context, stackID string) (*SlicedResult[Resource], error) {
	log := activity.GetLogger(ctx)
	resources := make([]Resource, 0)
	err := db.Filter(&Resource{}, &resources, db.QueryParams{"stack_id": stackID})

	if err != nil {
		log.Error("GetResources Error", "error", err)
	}

	return &SlicedResult[Resource]{Data: resources}, err
}

// GetWorkloads gets workloads from DB against a stack.
func (a *Activities) GetWorkloads(ctx context.Context, stackID string) (*SlicedResult[Workload], error) {
	log := activity.GetLogger(ctx)
	workloads := make([]Workload, 0)
	err := db.Filter(&Workload{}, &workloads, db.QueryParams{"stack_id": stackID})

	if err != nil {
		log.Error("GetWorkloads Error", "error", err)
	}

	return &SlicedResult[Workload]{Data: workloads}, err
}

// GetWorkloads gets workloads from DB against a stack.
func (a *Activities) GetRepos(ctx context.Context, stackID string) (*SlicedResult[Repo], error) {
	log := activity.GetLogger(ctx)
	repos := make([]Repo, 0)
	err := db.Filter(&Repo{}, &repos, db.QueryParams{"stack_id": stackID})

	if err != nil {
		log.Error("GetRepos Error", "error", err)
	}

	return &SlicedResult[Repo]{Data: repos}, err
}

// GetBluePrint gets blueprint from DB against a stack.
func (a *Activities) GetBluePrint(ctx context.Context, stackID string) (*Blueprint, error) {
	log := activity.GetLogger(ctx)
	blueprint := &Blueprint{}
	params := db.QueryParams{"stack_id": stackID}

	if err := db.Get(blueprint, params); err != nil {
		log.Error("GetBlueprint Error", "error", err)
		return blueprint, err
	}

	return blueprint, nil
}

func (a *Activities) CreateChangeset(ctx context.Context, changeSet *ChangeSet, ID gocql.UUID) error {
	err := db.CreateWithID(changeSet, ID)
	return err
}

func (a *Activities) GetChangeset(ctx context.Context, changeSetID gocql.UUID) (*ChangeSet, error) {
	c := new(ChangeSet)
	if err := db.Get(c, db.QueryParams{"id": changeSetID.String()}); err != nil {
		return c, err
	}

	return c, nil
}
