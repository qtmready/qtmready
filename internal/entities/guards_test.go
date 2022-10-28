// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package entities_test

import (
	"context"
	"testing"

	"github.com/Guilospanck/gocqlxmock"
	"github.com/gocql/gocql"

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
		stmt := "INSERT INTO guards (id,name,hashed,lookup_id,lookup_type,created_at,updated_at) VALUES (?,?,?,?,?,?,?) "
		names := []string{"id", "name", "hashed", "lookup_id", "lookup_type", "created_at", "updated_at"}
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
		stmt := "SELECT id,name,hashed,lookup_id,lookup_type,created_at,updated_at FROM guards WHERE id=" + arg.Guard.ID.String() + " ALLOW FILTERING "
		names := []string(nil)
		qmock := &gocqlxmock.QueryxMock{Ctx: context.Background(), Stmt: stmt, Names: names}

		db.DB.InitMockSession(smock)
		smock.On("Query", stmt, names).Return(qmock)
		qmock.On("GetRelease", arg.Guard).Return(nil)

		valid, err := arg.Guard.VerifyAPIKey(arg.Key)

		if err != nil {
			t.Errorf("unable to verify api key: %v", err)
		}

		if !valid {
			t.Errorf("Unable to Verify API Key")
		}
	}
}
