// Copyright © 2022, 2024, Breu, Inc. <info@breu.io>
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

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
