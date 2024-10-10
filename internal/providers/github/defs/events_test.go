package defs_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/gocql/gocql"
	"github.com/stretchr/testify/suite"

	coredefs "go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/providers/github/defs"
)

type (
	EventTestSuite struct {
		suite.Suite

		Provider coredefs.RepoProvider

		ParentID gocql.UUID
		RepoID   gocql.UUID

		TeamID     gocql.UUID
		UserID     gocql.UUID
		CoreRepoID gocql.UUID
	}
)

func (s *EventTestSuite) SetupSuite() {
	s.ParentID = db.MustUUID()
	s.TeamID = db.MustUUID()
	s.UserID = db.MustUUID()
}

func (s *EventTestSuite) TestCreateBranch() {
	f, err := os.ReadFile("testdata/create-branch.json")
	if s.NoError(err) {
		evt := &defs.CreateOrDeleteEvent{}

		if s.NoError(json.Unmarshal(f, evt)) {
			payload := evt.Payload()

			s.Assert().Equal(evt.Ref, payload.Ref)
			s.Assert().Equal(evt.RefType, payload.Kind)
		}
	}
}

func TestEventTestSuite(t *testing.T) {
	suite.Run(t, new(EventTestSuite))
}
