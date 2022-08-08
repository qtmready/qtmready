package entities_test

import (
	"testing"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/db/entities"
)

type fn struct {
	Provide interface{} // Can be nil
	Expect  interface{} // Can be nil
	Fn      func(args interface{}) func(*testing.T)
}

type fnMap map[string]fn

func TestUser(t *testing.T) {
	user := &entities.User{Password: "password"}
	user.PreCreate()
	tests := fnMap{
		"SetPassword":    fn{Provide: user, Expect: nil, Fn: userSetPassword},
		"VerifyPassword": fn{},
	}
	t.Run("GetTable", testTableName("users", user))
	t.Run("PreCreate", testPreCreate(user, tests))
}
func TestTeam(t *testing.T) {
	team := &entities.Team{}
	t.Run("GetTable", testTableName("teams", team))
}
func TestTeamUser(t *testing.T) {
	teamUser := &entities.TeamUser{}
	t.Run("GetTable", testTableName("team_users", teamUser))
}

func testTableName(expect string, entity db.Entity) func(*testing.T) {
	return func(t *testing.T) {
		if expect != entity.GetTable().Metadata().Name {
			t.Errorf("expected %s, got %s", expect, entity.GetTable().Metadata().Name)
		}
	}
}

func testPreCreate(entity db.Entity, args fnMap) func(*testing.T) {
	return func(t *testing.T) {
		for name, run := range args {
			t.Run(name, run.Fn(run.Provide))
		}
	}
}

func userSetPassword(args interface{}) func(*testing.T) {
	user := args.(*entities.User)
	return func(t *testing.T) {
		if user.Password == "password" {
			t.Errorf("expected password to be encrypted")
		}
	}
}
