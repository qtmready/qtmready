package entities_test

import (
	"testing"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/shared"
)

func testEntityGetTable(expect string, entity db.Entity) func(*testing.T) {
	return func(t *testing.T) {
		if expect != entity.GetTable().Metadata().M.Name {
			t.Errorf("expected %s, got %s", expect, entity.GetTable().Metadata().M.Name)
		}
	}
}

func testEntityOps(entity db.Entity, tests shared.TestFnMap) func(*testing.T) {
	return func(t *testing.T) {
		for name, test := range tests {
			t.Run(name, test.Run(test.Args, test.Want))
		}
	}
}
