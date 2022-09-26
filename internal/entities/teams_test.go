// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved.

package entities

import (
	"testing"
)

func TestTeam(t *testing.T) {
	team := &Team{}
	t.Run("GetTable", testTableName("teams", team))
}
