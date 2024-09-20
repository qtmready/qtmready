// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package auth_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
	"go.breu.io/quantm/internal/testutils"
)

type (
	containers struct {
		network *testcontainers.DockerNetwork
		db      *testutils.Container
	}

	ActivitiesTestSuite struct {
		suite.Suite
		ctx        context.Context
		ctrs       *containers
		activities *auth.Activities
	}
)

func (c *containers) shutdown(ctx context.Context) {
	shared.Logger().Info("graceful shutdown test environment ...")

	_ = c.db.DropKeyspace(db.TestKeyspace)
	_ = c.db.ShutdownCassandra()
	_ = c.network.Remove(ctx)

	db.DB().Session.Close()
	shared.Logger().Info("graceful shutdown complete.")
}

func (s *ActivitiesTestSuite) SetupSuite() {
	shared.InitServiceForTest()

	s.ctx = context.Background()
	s.setup_containers()
	s.activities = auth.NewActivities()
}

func (s *ActivitiesTestSuite) TearDownSuite() {
	s.ctrs.shutdown(s.ctx)
}

func (s *ActivitiesTestSuite) setup_containers() {
	shared.Logger().Info("setting up test environment ...")

	network, err := testutils.CreateTestNetwork(s.ctx)
	s.Require().NoError(err)

	dbctr, err := testutils.StartDBContainer(s.ctx)
	s.Require().NoError(err)

	err = dbctr.CreateKeyspace(db.TestKeyspace)
	s.Require().NoError(err)

	port, err := dbctr.Container.MappedPort(s.ctx, "9042")
	s.Require().NoError(err)

	db.NewE2ESession(port.Int(), "file://../db/migrations")

	s.ctrs = &containers{
		network: network,
		db:      dbctr,
	}
}

func (s *ActivitiesTestSuite) TestGetUser() {
	id, err := db.NewUUID()
	s.Require().NoError(err)

	user := &auth.User{
		ID:        id,
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	}
	s.Require().NoError(db.Save(user))

	fetched_user, err := s.activities.GetUser(s.ctx, db.QueryParams{"id": user.ID.String()})
	s.Require().NoError(err)
	s.Equal(user.ID, fetched_user.ID)
	s.Equal(user.Email, fetched_user.Email)
	s.Equal(user.FirstName, fetched_user.FirstName)
	s.Equal(user.LastName, fetched_user.LastName)
}

func (s *ActivitiesTestSuite) TestSaveUser() {
	id, err := db.NewUUID()
	s.Require().NoError(err)

	user := &auth.User{
		ID:        id,
		Email:     "save@example.com",
		FirstName: "Save",
		LastName:  "User",
	}

	saved_user, err := s.activities.SaveUser(s.ctx, user)
	s.Require().NoError(err)
	s.Equal(user.ID, saved_user.ID)

	fetched_user := &auth.User{}
	s.Require().NoError(db.Get(fetched_user, db.QueryParams{"id": user.ID.String()}))
	s.Equal(user.Email, fetched_user.Email)
	s.Equal(user.FirstName, fetched_user.FirstName)
	s.Equal(user.LastName, fetched_user.LastName)
}

func (s *ActivitiesTestSuite) TestCreateTeam() {
	id, err := db.NewUUID()
	s.Require().NoError(err)

	team := &auth.Team{
		ID:   id,
		Name: "Test Team",
		Slug: "test-team",
	}

	created_team, err := s.activities.CreateTeam(s.ctx, team)
	s.Require().NoError(err)
	s.Equal(team.ID, created_team.ID)

	fetched_team := &auth.Team{}
	s.Require().NoError(db.Get(fetched_team, db.QueryParams{"id": team.ID.String()}))
	s.Equal(team.Name, fetched_team.Name)
	s.Equal(team.Slug, fetched_team.Slug)
}

func (s *ActivitiesTestSuite) TestGetTeam() {
	id, err := db.NewUUID()
	s.Require().NoError(err)

	team := &auth.Team{
		ID:   id,
		Name: "Get Team",
		Slug: "get-team",
	}
	s.Require().NoError(db.Save(team))

	fetched_team, err := s.activities.GetTeam(s.ctx, db.QueryParams{"id": team.ID.String()})
	s.Require().NoError(err)
	s.Equal(team.ID, fetched_team.ID)
	s.Equal(team.Name, fetched_team.Name)
	s.Equal(team.Slug, fetched_team.Slug)
}

func (s *ActivitiesTestSuite) TestCreateOrUpdateTeamUser() {
	team_id, err := db.NewUUID()
	s.Require().NoError(err)
	user_id, err := db.NewUUID()
	s.Require().NoError(err)

	team := &auth.Team{ID: team_id, Name: "Team", Slug: "team"}
	user := &auth.User{ID: user_id, Email: "team@example.com"}

	s.Require().NoError(db.Save(team))
	s.Require().NoError(db.Save(user))

	team_user_id, err := db.NewUUID()

	s.Require().NoError(err)

	team_user := &auth.TeamUser{
		ID:       team_user_id,
		TeamID:   team.ID,
		UserID:   user.ID,
		IsAdmin:  true,
		IsActive: true,
		Role:     "Member",
	}

	updated_team_user, err := s.activities.CreateOrUpdateTeamUser(s.ctx, team_user)
	s.Require().NoError(err)
	s.Equal(team_user.ID, updated_team_user.ID)

	fetched_team_user := &auth.TeamUser{}
	s.Require().NoError(db.Get(fetched_team_user, db.QueryParams{"id": team_user.ID.String()}))
	s.Equal(team_user.TeamID, fetched_team_user.TeamID)
	s.Equal(team_user.UserID, fetched_team_user.UserID)
	s.Equal(team_user.IsAdmin, fetched_team_user.IsAdmin)
	s.Equal(team_user.IsActive, fetched_team_user.IsActive)
	s.Equal(team_user.Role, fetched_team_user.Role)
}

// func (s *ActivitiesTestSuite) TestGetTeamUser() {
// 	team_id, err := db.NewUUID()
// 	s.Require().NoError(err)
// 	user_id, err := db.NewUUID()
// 	s.Require().NoError(err)

// 	team := &auth.Team{ID: team_id, Name: "GetTeam", Slug: "get-team"}
// 	user := &auth.User{ID: user_id, Email: "getteam@example.com"}

// 	s.Require().NoError(db.Save(team))
// 	s.Require().NoError(db.Save(user))

// 	team_user_id, err := db.NewUUID()

// 	s.Require().NoError(err)

// 	team_user := &auth.TeamUser{
// 		ID:       team_user_id,
// 		TeamID:   team.ID,
// 		UserID:   user.ID,
// 		IsAdmin:  true,
// 		IsActive: true,
// 		Role:     "Admin",
// 	}
// 	s.Require().NoError(db.Save(team_user))

// 	fetched_team_user, err := s.activities.GetTeamUser(s.ctx, user.Email)
// 	s.Require().NoError(err)
// 	s.Equal(team_user.ID, fetched_team_user.ID)
// 	s.Equal(team_user.TeamID, fetched_team_user.TeamID)
// 	s.Equal(team_user.UserID, fetched_team_user.UserID)
// 	s.Equal(team_user.IsAdmin, fetched_team_user.IsAdmin)
// 	s.Equal(team_user.IsActive, fetched_team_user.IsActive)
// 	s.Equal(team_user.Role, fetched_team_user.Role)
// }

func TestActivitiesSuite(t *testing.T) {
	suite.Run(t, new(ActivitiesTestSuite))
}
