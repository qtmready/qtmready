package entities

import (
	"testing"

	"go.breu.io/ctrlplane/internal/shared"
)

var (
	testKey string
)

func TestGuard(t *testing.T) {
	guard := &Guard{}
	testKey = guard.GenerateRandomValue()
	guard.Hashed = testKey
	_ = guard.PreCreate()

	preCreateTests := shared.TestFnMap{
		"SetHashed":    shared.TestFn{Args: guard, Want: nil, Fn: testGuardSetHashed},
		"VerifyHashed": shared.TestFn{Args: guard, Want: nil, Fn: testGuardVerifyHashed},
	}

	t.Run("GetTable", testTableName("guards", guard))
	t.Run("PreCreate", testPreCreate(guard, preCreateTests))
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
