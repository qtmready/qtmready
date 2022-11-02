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
