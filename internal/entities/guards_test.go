// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package entities_test

import (
	"testing"

	"github.com/gocql/gocql"
	"go.breu.io/ctrlplane/internal/entities"
	"go.breu.io/ctrlplane/internal/shared"
)

type (
	GuardWithKey struct {
		Key    string
		Token  string
		Prefix string
		Guard  *entities.Guard
	}
)

func TestTeamGuard(t *testing.T) {
	id, _ := gocql.RandomUUID()

	guard := &entities.Guard{}
	key := guard.NewForTeam(id)
	_ = guard.PreCreate()
	prefix, token, _ := guard.SplitAPIKey(key)
	args := &GuardWithKey{Key: key, Prefix: prefix, Token: token, Guard: guard}

	opsTest := shared.TestFnMap{
		"TokenEncryption": shared.TestFn{Args: args, Want: nil, Run: testTokenEncryption},
		"PrefixToID":      shared.TestFn{Args: args, Want: nil, Run: testPrefixToID},
		"VerifyToken":     shared.TestFn{Args: args, Want: nil, Run: testVerifyToken},
		"VerifyAPIKey":    shared.TestFn{Args: args, Want: nil, Run: testVerifyAPIKey},
		"TestGuardName":   shared.TestFn{Args: args, Want: nil, Run: testTeamGuardName},
	}

	t.Run("GetTable", testEntityGetTable("guards", guard))
	t.Run("EntityOps", testEntityOps(guard, opsTest))
}

func TestUserGuard(t *testing.T) {
	id, _ := gocql.RandomUUID()

	guard := &entities.Guard{}
	key := guard.NewForUser("test", id)
	_ = guard.PreCreate()
	prefix, token, _ := guard.SplitAPIKey(key)
	args := &GuardWithKey{Key: key, Prefix: prefix, Token: token, Guard: guard}

	opsTest := shared.TestFnMap{
		"TokenEncryption": shared.TestFn{Args: args, Want: nil, Run: testTokenEncryption},
		"PrefixToID":      shared.TestFn{Args: args, Want: nil, Run: testPrefixToID},
		"VerifyToken":     shared.TestFn{Args: args, Want: nil, Run: testVerifyToken},
		"VerifyAPIKey":    shared.TestFn{Args: args, Want: nil, Run: testVerifyAPIKey},
		"TestGuardName":   shared.TestFn{Args: args, Want: nil, Run: testUserGuardName},
	}

	t.Run("GetTable", testEntityGetTable("guards", guard))
	t.Run("EntityOps", testEntityOps(guard, opsTest))
}

func testTokenEncryption(args interface{}, want interface{}) func(*testing.T) {
	guard := args.(*GuardWithKey)

	return func(t *testing.T) {
		if guard.Token == guard.Guard.Hashed {
			t.Errorf("Expected token to be hashed")
		}
	}
}

func testPrefixToID(args interface{}, want interface{}) func(*testing.T) {
	guard := args.(*GuardWithKey)
	id, _ := guard.Guard.PrefixToID(guard.Prefix)

	return func(t *testing.T) {
		if id.String() != guard.Guard.LookupID.String() {
			t.Errorf("prefix mismatch when verifying")
		}
	}
}

func testVerifyToken(args interface{}, want interface{}) func(*testing.T) {
	guard := args.(*GuardWithKey)

	return func(t *testing.T) {
		if !guard.Guard.VerifyToken(guard.Token) {
			t.Errorf("unable to verify token")
		}
	}
}

func testVerifyAPIKey(args interface{}, want interface{}) func(*testing.T) {
	guard := args.(*GuardWithKey)

	return func(t *testing.T) {
		v, _ := guard.Guard.VerifyAPIKey(guard.Key)
		if !v {
			t.Errorf("unable to verify APIKey")
		}
	}
}

func testTeamGuardName(args interface{}, want interface{}) func(*testing.T) {
	guard := args.(*GuardWithKey)

	return func(t *testing.T) {
		if guard.Guard.Name != "default" {
			t.Errorf("team guard should always be named default")
		}
	}
}

func testUserGuardName(args interface{}, want interface{}) func(*testing.T) {
	guard := args.(*GuardWithKey)

	return func(t *testing.T) {
		if guard.Guard.Name != "test" {
			t.Errorf("user guard name not setting correctly")
		}
	}
}
