package auth_test

import (
	"testing"

	"github.com/gosimple/slug"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/testutils"
)

func TestTeam(t *testing.T) {
	team := &auth.Team{
		Name: "Team Name",
	}
	_ = team.PreCreate()

	opsTests := testutils.TestFnMap{
		"Slug": testutils.TestFn{Args: team, Want: nil, Run: testTeamSlug},
	}

	t.Run("GetTable", testutils.TestEntityGetTable("teams", team))
	t.Run("EntityOps", testutils.TestEntityOps(team, opsTests))
}

func testTeamSlug(args, want any) func(*testing.T) {
	team := args.(*auth.Team)
	sluglen := len(slug.Make(team.Name)) + 1 + 4

	return func(t *testing.T) {
		if len(team.Slug) != sluglen {
			t.Errorf("slug length is not correct, got: %d, want: %d", len(team.Slug), sluglen)
		}
	}
}
