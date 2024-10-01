// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2022, 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
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

// Copyright © 2022, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRUBTING ANY OF
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
	"context"
	"testing"

	"github.com/Guilospanck/gocqlxmock"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/qb"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/testutils"
)

type (
	guardnkey struct {
		Key   string
		Guard *auth.Guard
	}
)

func TestTeamGuard(t *testing.T) {
	teamID, _ := gocql.RandomUUID()
	guard := &auth.Guard{}
	key := guard.NewForTeam(teamID)
	args := &guardnkey{Key: key, Guard: guard}

	opsTest := testutils.TestFnMap{
		"TestGuardName":    testutils.TestFn{Args: args, Want: nil, Run: testTeamGuardName},
		"TestSave":         testutils.TestFn{Args: args, Want: nil, Run: testSave},
		"TestVerifyAPIKey": testutils.TestFn{Args: args, Want: nil, Run: testVerifyAPIKey},
	}

	t.Run("GetTable", testutils.TestEntityGetTable("guards", guard))
	t.Run("EntityOps", testutils.TestEntityOps(guard, opsTest))
}

func TestUserGuard(t *testing.T) {
	userID, _ := gocql.RandomUUID()
	guard := &auth.Guard{}
	key := guard.NewForUser("test", userID)
	args := &guardnkey{Key: key, Guard: guard}

	opsTest := testutils.TestFnMap{
		"TestGuardName":    testutils.TestFn{Args: args, Want: nil, Run: testUserGuardName},
		"TestSave":         testutils.TestFn{Args: args, Want: nil, Run: testSave},
		"TestVerifyAPIKey": testutils.TestFn{Args: args, Want: nil, Run: testVerifyAPIKey},
	}

	t.Run("GetTable", testutils.TestEntityGetTable("guards", guard))
	t.Run("EntityOps", testutils.TestEntityOps(guard, opsTest))
}

func testTeamGuardName(args, want any) func(*testing.T) {
	arg := args.(*guardnkey)

	return func(t *testing.T) {
		if arg.Guard.Name != "default" || arg.Guard.LookupType != "team" {
			t.Errorf("expected name to be 'default', got %s", arg.Guard.Name)
		}
	}
}

func testUserGuardName(args, want any) func(*testing.T) {
	arg := args.(*guardnkey)

	return func(t *testing.T) {
		if arg.Guard.Name != "test" || arg.Guard.LookupType != "user" {
			t.Errorf("expected name to be 'test', got %s", arg.Guard.Name)
		}
	}
}

func testSave(args, want any) func(*testing.T) {
	arg := args.(*guardnkey)

	return func(t *testing.T) {
		smock := &gocqlxmock.SessionxMock{}
		stmt, names := arg.Guard.GetTable().Insert()
		qmock := &gocqlxmock.QueryxMock{Ctx: context.Background(), Stmt: stmt, Names: names}

		db.NewMockSession(smock)
		smock.On("Query", stmt, names).Return(qmock)
		qmock.On("BindStruct", arg.Guard).Return(qmock)
		smock.On("Close").Return()
		qmock.On("ExecRelease").Return(nil)

		if err := arg.Guard.Save(); err != nil {
			t.Errorf("unable to save guard: %v", err)
		}

		if arg.Key == arg.Guard.Hashed {
			t.Errorf("expected hashed to be different, got %s", arg.Guard.Hashed)
		}

		t.Cleanup(func() {
			db.DB().Session.Close()
		})
	}
}

func testVerifyAPIKey(args, want any) func(*testing.T) {
	arg := args.(*guardnkey)

	return func(t *testing.T) {
		smock := &gocqlxmock.SessionxMock{}
		stmt, names := db.SelectBuilder(arg.Guard.GetTable().Name()).
			Columns(arg.Guard.GetTable().Metadata().M.Columns...).
			AllowFiltering().
			Where(qb.EqLit("id", arg.Guard.ID.String())).
			ToCql()
		qmock := &gocqlxmock.QueryxMock{Ctx: context.Background(), Stmt: stmt, Names: names}

		db.NewMockSession(smock)
		smock.On("Close").Return()
		smock.On("Query", stmt, names).Return(qmock)
		qmock.On("GetRelease", arg.Guard).Return(nil)

		err := arg.Guard.VerifyAPIKey(arg.Key)

		if err != nil {
			t.Errorf("unable to verify api key: %v", err)
		}

		t.Cleanup(func() {
			db.DB().Session.Close()
		})
	}
}
