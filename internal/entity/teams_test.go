// Copyright © 2022, Breu, Inc. <info@breu.io>. All rights reserved.
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

package entity_test

import (
	"testing"

	"github.com/gosimple/slug"

	"go.breu.io/ctrlplane/internal/entity"
	"go.breu.io/ctrlplane/internal/shared"
)

func TestTeam(t *testing.T) {
	team := &entity.Team{
		Name: "Team Name",
	}
	_ = team.PreCreate()

	opsTests := shared.TestFnMap{
		"Slug": shared.TestFn{Args: team, Want: nil, Run: testTeamSlug},
	}

	t.Run("GetTable", testEntityGetTable("teams", team))
	t.Run("EntityOps", testEntityOps(team, opsTests))
}

func testTeamSlug(args interface{}, want interface{}) func(*testing.T) {
	team := args.(*entity.Team)
	sluglen := len(slug.Make(team.Name)) + 1 + 22

	return func(t *testing.T) {
		if len(team.Slug) != sluglen {
			t.Errorf("slug length is not correct, got: %d, want: %d", len(team.Slug), sluglen)
		}
	}
}
