package testutils

import (
	"testing"

	"go.breu.io/quantm/internal/db"
)

func TestEntityGetTable(expect string, entity db.Entity) func(*testing.T) {
	return func(t *testing.T) {
		if expect != entity.GetTable().Metadata().M.Name {
			t.Errorf("expected %s, got %s", expect, entity.GetTable().Metadata().M.Name)
		}
	}
}

func TestEntityOps(_ db.Entity, tests TestFnMap) func(*testing.T) {
	return func(t *testing.T) {
		for name, test := range tests {
			t.Run(name, test.Run(test.Args, test.Want))
		}
	}
}
