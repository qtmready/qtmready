// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 

package entities

import (
	"testing"

	"go.breu.io/ctrlplane/internal/db"
	"go.breu.io/ctrlplane/internal/shared"
)

func testTableName(expect string, entity db.Entity) func(*testing.T) {
	return func(t *testing.T) {
		if expect != entity.GetTable().Metadata().Name {
			t.Errorf("expected %s, got %s", expect, entity.GetTable().Metadata().Name)
		}
	}
}

func testPreCreate(entity db.Entity, tests shared.TestFnMap) func(*testing.T) {
	return func(t *testing.T) {
		for name, test := range tests {
			t.Run(name, test.Fn(test.Args, test.Want))
		}
	}
}
