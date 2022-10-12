// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package entities_test

import (
	"testing"

	"go.breu.io/ctrlplane/internal/entities"
)

func TestTeam(t *testing.T) {
	team := &entities.Team{}
	t.Run("GetTable", testEntityGetTable("teams", team))
}
