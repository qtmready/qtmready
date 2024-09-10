package auth_test

import (
	"testing"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/testutils"
)

var (
	password string
)

func TestUser(t *testing.T) {
	password = "password"
	user := &auth.User{Password: password}
	_ = user.PreCreate()

	opsTests := testutils.TestFnMap{
		"SetPassword":    testutils.TestFn{Args: user, Want: nil, Run: testUserSetPassword},
		"VerifyPassword": testutils.TestFn{Args: user, Want: nil, Run: testUserVerifyPassword},
	}

	t.Run("GetTable", testutils.TestEntityGetTable("users", user))
	t.Run("EntityOps", testutils.TestEntityOps(user, opsTests))
}

func testUserSetPassword(args, want any) func(*testing.T) {
	user := args.(*auth.User)

	return func(t *testing.T) {
		if user.Password == "password" {
			t.Errorf("expected password to be encrypted")
		}
	}
}

func testUserVerifyPassword(args, want any) func(*testing.T) {
	v := args.(*auth.User)

	return func(t *testing.T) {
		if !v.VerifyPassword(password) {
			t.Errorf("expected password to be verified")
		}
	}
}
