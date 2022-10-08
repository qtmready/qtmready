// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package entities

import (
	"testing"

	"go.breu.io/ctrlplane/internal/shared"
)

var (
	testPassword string
)

func TestUser(t *testing.T) {
	testPassword = "password"
	user := &User{Password: testPassword}
	_ = user.PreCreate()

	preCreateTests := shared.TestFnMap{
		"SetPassword":    shared.TestFn{Args: user, Want: nil, Run: testUserSetPassword},
		"VerifyPassword": shared.TestFn{Args: user, Want: nil, Run: testUserVerifyPassword},
	}

	t.Run("GetTable", testEntityGetTable("users", user))
	t.Run("PreCreate", testEntityPreCreate(user, preCreateTests))
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
	v := args.(*User)

	return func(t *testing.T) {
		if !v.VerifyPassword(testPassword) {
			t.Errorf("expected password to be verified")
		}
	}
}
