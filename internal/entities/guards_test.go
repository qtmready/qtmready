// Copyright © 2022, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLATING, DOWNLOADING, ACCESSING, USING OR DISTRUBTING ANY OF
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
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY  APPLICABLE LAW.

// This software is made available by Breu, Inc., under the terms of the Breu  Community License Agreement, Version 1.0 located at  http://www.breu.io/breu-community-license/v1. BY INSTALLING, DOWNLOADING,  ACCESSING, USING OR DISTRIBUTING ANY OF THE SOFTWARE, YOU AGREE TO THE TERMS  OF SUCH LICENSE AGREEMENT.

package entities_test

import (
	"context"
	"testing"

	"github.com/Guilospanck/gocqlxmock"
	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2/qb"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/entities"
	"go.breu.io/ctrlplane/internal/shared"
)

type (
	guardnkey struct {
		Key   string
		Guard *entities.Guard
	}
)

func TestTeamGuard(t *testing.T) {
	teamID, _ := gocql.RandomUUID()
	guard := &entities.Guard{}
	key := guard.NewForTeam(teamID)
	args := &guardnkey{Key: key, Guard: guard}

	opsTest := shared.TestFnMap{
		"TestGuardName":    shared.TestFn{Args: args, Want: nil, Run: testTeamGuardName},
		"TestSave":         shared.TestFn{Args: args, Want: nil, Run: testSave},
		"TestVerifyAPIKey": shared.TestFn{Args: args, Want: nil, Run: testVerifyAPIKey},
	}

	t.Run("GetTable", testEntityGetTable("guards", guard))
	t.Run("EntityOps", testEntityOps(guard, opsTest))
}

func TestUserGuard(t *testing.T) {
	userID, _ := gocql.RandomUUID()
	guard := &entities.Guard{}
	key := guard.NewForUser("test", userID)
	args := &guardnkey{Key: key, Guard: guard}

	opsTest := shared.TestFnMap{
		"TestGuardName":    shared.TestFn{Args: args, Want: nil, Run: testUserGuardName},
		"TestSave":         shared.TestFn{Args: args, Want: nil, Run: testSave},
		"TestVerifyAPIKey": shared.TestFn{Args: args, Want: nil, Run: testVerifyAPIKey},
	}

	t.Run("GetTable", testEntityGetTable("guards", guard))
	t.Run("EntityOps", testEntityOps(guard, opsTest))
}

func testTeamGuardName(args interface{}, want interface{}) func(*testing.T) {
	arg := args.(*guardnkey)

	return func(t *testing.T) {
		if arg.Guard.Name != "default" || arg.Guard.LookupType != "team" {
			t.Errorf("expected name to be 'default', got %s", arg.Guard.Name)
		}
	}
}

func testUserGuardName(args interface{}, want interface{}) func(*testing.T) {
	arg := args.(*guardnkey)

	return func(t *testing.T) {
		if arg.Guard.Name != "test" || arg.Guard.LookupType != "user" {
			t.Errorf("expected name to be 'test', got %s", arg.Guard.Name)
		}
	}
}

func testSave(args interface{}, want interface{}) func(*testing.T) {
	arg := args.(*guardnkey)

	return func(t *testing.T) {
		smock := &gocqlxmock.SessionxMock{}
		stmt, names := arg.Guard.GetTable().Insert()
		qmock := &gocqlxmock.QueryxMock{Ctx: context.Background(), Stmt: stmt, Names: names}

		db.DB.InitMockSession(smock)
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
			db.DB.Session.Close()
		})
	}
}

func testVerifyAPIKey(args interface{}, want interface{}) func(*testing.T) {
	arg := args.(*guardnkey)

	return func(t *testing.T) {
		smock := &gocqlxmock.SessionxMock{}
		stmt, names := db.SelectBuilder(arg.Guard.GetTable().Name()).
			Columns(arg.Guard.GetTable().Metadata().M.Columns...).
			AllowFiltering().
			Where(qb.EqLit("id", arg.Guard.ID.String())).
			ToCql()
		qmock := &gocqlxmock.QueryxMock{Ctx: context.Background(), Stmt: stmt, Names: names}

		db.DB.InitMockSession(smock)
		smock.On("Close").Return()
		smock.On("Query", stmt, names).Return(qmock)
		qmock.On("GetRelease", arg.Guard).Return(nil)

		valid, err := arg.Guard.VerifyAPIKey(arg.Key)

		if err != nil {
			t.Errorf("unable to verify api key: %v", err)
		}

		if !valid {
			t.Errorf("Unable to Verify API Key")
		}

		t.Cleanup(func() {
			db.DB.Session.Close()
		})
	}
}
