// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package entities

import (
	"testing"

	"go.breu.io/ctrlplane/internal/shared"
)

var (
	testKey string
)

func TestGuard(t *testing.T) {
	args := &Guard{}
	testKey = args.GenerateRandomValue()
	args.Hashed = testKey
	_ = args.PreCreate()

	preCreateTests := shared.TestFnMap{
		"SetHashed":    shared.TestFn{Args: args, Want: nil, Run: testGuardSetHashed},
		"VerifyHashed": shared.TestFn{Args: args, Want: nil, Run: testGuardVerifyHashed},
	}

	t.Run("GetTable", testEntityGetTable("guards", args))
	t.Run("PreCreate", testEntityPreCreate(args, preCreateTests))
}

func testGuardSetHashed(args interface{}, want interface{}) func(*testing.T) {
	guard := args.(*Guard)

	return func(t *testing.T) {
		if guard.Hashed == testKey {
			t.Errorf("expected hashed to be encrypted")
		}
	}
}

func testGuardVerifyHashed(args interface{}, want interface{}) func(*testing.T) {
	v := args.(*Guard)

	return func(t *testing.T) {
		if !v.VerifyHashed(testKey) {
			t.Errorf("expected hashed to be verified")
		}
	}
}
