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

package entity_test

import (
	"testing"

	"go.breu.io/ctrlplane/internal/entity"
	"go.breu.io/ctrlplane/internal/shared"
)

var (
	password string
)

func TestUser(t *testing.T) {
	password = "password"
	user := &entity.User{Password: password}
	_ = user.PreCreate()

	opsTests := shared.TestFnMap{
		"SetPassword":    shared.TestFn{Args: user, Want: nil, Run: testUserSetPassword},
		"VerifyPassword": shared.TestFn{Args: user, Want: nil, Run: testUserVerifyPassword},
	}

	t.Run("GetTable", testEntityGetTable("users", user))
	t.Run("EntityOps", testEntityOps(user, opsTests))
}

func testUserSetPassword(args interface{}, want interface{}) func(*testing.T) {
	user := args.(*entity.User)

	return func(t *testing.T) {
		if user.Password == "password" {
			t.Errorf("expected password to be encrypted")
		}
	}
}

func testUserVerifyPassword(args interface{}, want interface{}) func(*testing.T) {
	v := args.(*entity.User)

	return func(t *testing.T) {
		if !v.VerifyPassword(password) {
			t.Errorf("expected password to be verified")
		}
	}
}
