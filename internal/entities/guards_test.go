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
		Key             string
		Token           string
		EncodedID       string
		EncodedLookupID string
		Guard           *entities.Guard
	}
)

func TestTeamGuard(t *testing.T) {
	id, _ := gocql.RandomUUID()

	guard := &entities.Guard{}
	key := guard.NewForTeam(id)
	_ = guard.PreCreate()
	encodedID, encodedLookupID, token, err := guard.SplitAPIKey(key)

	if err != nil {
		t.Errorf("unable to split api key")
	}

	args := &GuardWithKey{Key: key, EncodedID: encodedID, EncodedLookupID: encodedLookupID, Token: token, Guard: guard}

	opsTest := shared.TestFnMap{
		"TokenEncryption": shared.TestFn{Args: args, Want: nil, Run: testTokenEncryption},
		"DecodeUUID":      shared.TestFn{Args: args, Want: nil, Run: testDecodeUUID},
		"VerifyToken":     shared.TestFn{Args: args, Want: nil, Run: testVerifyToken},
		"TestGuardName":   shared.TestFn{Args: args, Want: nil, Run: testTeamGuardName},
		"VerifyAPIKey":    shared.TestFn{Args: args, Want: nil, Run: testVerifyAPIKey},
	}

	t.Run("GetTable", testEntityGetTable("guards", guard))
	t.Run("EntityOps", testEntityOps(guard, opsTest))
}

func TestUserGuard(t *testing.T) {
	id, _ := gocql.RandomUUID()

	guard := &entities.Guard{}
	key := guard.NewForUser("test", id)
	_ = guard.Save()
	encodedID, encodedLookupID, token, err := guard.SplitAPIKey(key)

	if err != nil {
		t.Errorf("unable to split api key")
	}

	args := &GuardWithKey{Key: key, EncodedID: encodedID, EncodedLookupID: encodedLookupID, Token: token, Guard: guard}

	opsTest := shared.TestFnMap{
		"TokenEncryption": shared.TestFn{Args: args, Want: nil, Run: testTokenEncryption},
		"DecodeUUID":      shared.TestFn{Args: args, Want: nil, Run: testDecodeUUID},
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

func testDecodeUUID(args interface{}, want interface{}) func(*testing.T) {
	guard := args.(*GuardWithKey)
	id, err := guard.Guard.DecodeUUID(guard.EncodedLookupID)

	return func(t *testing.T) {
		if err != nil {
			t.Errorf("unable to decode prefix to uuid")
		}

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

func testVerifyAPIKey(args interface{}, want interface{}) func(*testing.T) {
	guard := args.(*GuardWithKey)
	valid, err := guard.Guard.VerifyAPIKey(guard.Key)

	return func(t *testing.T) {
		if err != nil {
			t.Errorf("unable to verify api key")
		}

		if !valid {
			t.Errorf("api key is not valid")
		}
	}
}
