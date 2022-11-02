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

	"go.breu.io/ctrlplane/internal/entities"
	"go.breu.io/ctrlplane/internal/shared"
)

var (
	password string
)

func TestUser(t *testing.T) {
	password = "password"
	user := &entities.User{Password: password}
	_ = user.PreCreate()

	opsTests := shared.TestFnMap{
		"SetPassword":    shared.TestFn{Args: user, Want: nil, Run: testUserSetPassword},
		"VerifyPassword": shared.TestFn{Args: user, Want: nil, Run: testUserVerifyPassword},
	}

	t.Run("GetTable", testEntityGetTable("users", user))
	t.Run("EntityOps", testEntityOps(user, opsTests))
}

func testUserSetPassword(args interface{}, want interface{}) func(*testing.T) {
	user := args.(*entities.User)

	return func(t *testing.T) {
		if user.Password == "password" {
			t.Errorf("expected password to be encrypted")
		}
	}
}

func testUserVerifyPassword(args interface{}, want interface{}) func(*testing.T) {
	v := args.(*entities.User)

	return func(t *testing.T) {
		if !v.VerifyPassword(password) {
			t.Errorf("expected password to be verified")
		}
	}
}
