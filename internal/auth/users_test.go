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
