// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
//
// This software is made available by Breu, Inc., under the terms of the Breu  
// Community License Agreement, Version 1.0 located at  
// http://www.breu.io/breu-community-license/v1. BY INSTALLING, DOWNLOADING,  
// ACCESSING, USING OR DISTRIBUTING ANY OF THE SOFTWARE, YOU AGREE TO THE TERMS  
// OF SUCH LICENSE AGREEMENT. 

package entities_test

import (
	"testing"

	"github.com/gosimple/slug"

	"go.breu.io/ctrlplane/internal/entities"
	"go.breu.io/ctrlplane/internal/shared"
)

func TestTeam(t *testing.T) {
	team := &entities.Team{
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
	team := args.(*entities.Team)
	sluglen := len(slug.Make(team.Name)) + 1 + 22

	return func(t *testing.T) {
		if len(team.Slug) != sluglen {
			t.Errorf("slug length is not correct, got: %d, want: %d", len(team.Slug), sluglen)
		}
	}
}
