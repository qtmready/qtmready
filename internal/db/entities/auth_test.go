package entities

import (
	"testing"

	"go.breu.io/ctrlplane/internal/db"
)

type (
	fn struct {
		Args interface{} // Can be nil
		Want interface{} // Can be nil
		Fn   func(args interface{}, want interface{}) func(*testing.T)
	}

	testMap map[string]fn

	testVerifyPasswordArgs struct {
		User     *User
		Password string
	}
)

func TestUser(t *testing.T) {
	user := &User{Password: "password"}
	user.PreCreate()
	createTests := testMap{
		"SetPassword":    fn{user, nil, testUserSetPassword},
		"VerifyPassword": fn{testVerifyPasswordArgs{user, "password"}, nil, testUserVerifyPassword},
	}
	t.Run("GetTable", testTableName("users", user))
	t.Run("PreCreate", testPreCreate(user, createTests))
}
func TestTeam(t *testing.T) {
	team := &Team{}
	t.Run("GetTable", testTableName("teams", team))
}
func TestTeamUser(t *testing.T) {
	teamUser := &TeamUser{}
	t.Run("GetTable", testTableName("team_users", teamUser))
}

func testTableName(expect string, entity db.Entity) func(*testing.T) {
	return func(t *testing.T) {
		if expect != entity.GetTable().Metadata().Name {
			t.Errorf("expected %s, got %s", expect, entity.GetTable().Metadata().Name)
		}
	}
}

func testPreCreate(entity db.Entity, tests testMap) func(*testing.T) {
	return func(t *testing.T) {
		for name, test := range tests {
			t.Run(name, test.Fn(test.Args, test.Want))
		}
	}
}

func testUserSetPassword(args interface{}, want interface{}) func(*testing.T) {
	user := args.(*User)
	return func(t *testing.T) {
		if user.Password == "password" {
			t.Errorf("expected password to be encrypted")
		}
	}
}

func testUserVerifyPassword(args interface{}, want interface{}) func(*testing.T) {
	v := args.(testVerifyPasswordArgs)
	return func(t *testing.T) {
		if !v.User.VerifyPassword(v.Password) {
			t.Errorf("expected password to be verified")
		}
	}
}
