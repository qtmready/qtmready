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

// Copyright © 2022, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRUBTING ANY OF
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

package entity_test

import (
	"testing"

	"github.com/gocql/gocql"
	"github.com/gosimple/slug"

	"go.breu.io/ctrlplane/internal/entity"
	"go.breu.io/ctrlplane/internal/shared"
)

func TestApp(t *testing.T) {
	stack := &entity.Stack{
		ID:     gocql.MustRandomUUID(),
		Name:   "Test Stack",
		Config: entity.StackConfig{},
		TeamID: gocql.MustRandomUUID(),
	}
	_ = stack.PreCreate()

	opsTests := shared.TestFnMap{
		"Slug": shared.TestFn{Args: stack, Want: nil, Run: testAppSlug},
	}

	t.Run("GetTable", testEntityGetTable("stacks", stack))
	t.Run("EntityOps", testEntityOps(stack, opsTests))
}

func TestRepo(t *testing.T) {
	repo := &entity.Repo{}
	t.Run("GetTable", testEntityGetTable("repos", repo))
}

func TestWorkload(t *testing.T) {
	workload := &entity.Workload{}
	t.Run("GetTable", testEntityGetTable("workloads", workload))
}

func TestResource(t *testing.T) {
	resource := &entity.Resource{}
	t.Run("GetTable", testEntityGetTable("resources", resource))
}

func TestBlueprint(t *testing.T) {
	blueprint := &entity.Blueprint{}
	t.Run("GetTable", testEntityGetTable("blueprints", blueprint))
}

func TestRollout(t *testing.T) {
	rollout := &entity.Rollout{}
	t.Run("GetTable", testEntityGetTable("rollouts", rollout))
}

func testAppSlug(args interface{}, want interface{}) func(*testing.T) {
	app := args.(*entity.Stack)
	sluglen := len(slug.Make(app.Name)) + 1 + 22

	return func(t *testing.T) {
		if len(app.Slug) != sluglen {
			t.Errorf("slug length is not correct, got: %d, want: %d", len(app.Slug), sluglen)
		}
	}
}
