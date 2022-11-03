package entities_test

import (
	"testing"

	"go.breu.io/ctrlplane/internal/entities"
)

func TestTeamUser(t *testing.T) {
	teamUser := &entities.TeamUser{}
	t.Run("GetTable", testEntityGetTable("team_users", teamUser))
}
