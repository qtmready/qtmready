package defs_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/shared"
)

type (
	EventTestSuite struct {
		suite.Suite

		parent  gocql.UUID
		subject defs.EventSubject
	}
)

func (s *EventTestSuite) SetupSuite() {
	shared.InitServiceForTest()

	s.parent = gocql.TimeUUID()
	s.subject = defs.EventSubject{
		ID:   gocql.MustRandomUUID(),
		Name: "repos",
	}
}

func (s *EventTestSuite) TestBranchCreate() {
	branch := &defs.BranchOrTag{
		Ref:           "test-branch",
		DefaultBranch: "main",
	}

	bytes, _ := json.Marshal(branch)

	event := branch.ToEvent(
		defs.RepoProviderGithub,
		s.subject,
		defs.EventActionCreated,
		defs.EventScopeBranch,
	)
	event.SetSource("test/test")

	expected := fmt.Sprintf(`{
  "version": "v1",
  "context": {
    "id": "%s",
    "parent_id": "00000000-0000-0000-0000-000000000000",
    "provider": "github",
    "scope": "branch",
    "action": "created",
    "source": "test/test",
    "timestamp": "%s"
  },
  "subject": {
    "id": "%s",
    "name": "%s"
  },
  "payload": {
    "ref": "test-branch",
    "default_branch": "main"
  }
}`,
		event.Context.ID,
		event.Context.Timestamp.Format(time.RFC3339Nano),
		s.subject.ID, s.subject.Name,
	)

	// Test Marshal to JSON
	marshal, err := json.MarshalIndent(event, "", "  ")
	if s.NoError(err) {
		s.T().Log(string(marshal))
		s.Equal(defs.EventVersionV1, event.Version)
		s.Equal(uuid.Nil.String(), event.Context.ParentID.String())
		s.Equal(expected, string(marshal))
	}

	// Test Unmarshal from JSON
	var unmarshal defs.Event[defs.BranchOrTag, defs.RepoProvider]

	err = json.Unmarshal(marshal, &unmarshal)
	if s.NoError(err) {
		s.Equal(event.Version, unmarshal.Version)
		s.Equal(event.Context.ID, unmarshal.Context.ID)
		s.Equal(event.Context.Provider, unmarshal.Context.Provider)
		s.Equal(event.Context.Scope, unmarshal.Context.Scope)
		s.Equal(event.Context.Action, unmarshal.Context.Action)
		s.Equal(event.Context.Source, unmarshal.Context.Source)
		s.Equal(event.Context.Timestamp.Unix(), unmarshal.Context.Timestamp.Unix())
		s.Equal(event.Subject, unmarshal.Subject)
		s.Equal(event.Data, unmarshal.Data)
	}

	// Test flattening the event
	flat, err := event.Flatten()
	if s.NoError(err) {
		s.Equal(event.Version, flat.Version)
		s.Equal(event.Context.ID, flat.ID)
		s.Equal(event.Context.Provider, flat.Provider)
		s.Equal(event.Context.Scope, flat.Scope)
		s.Equal(event.Context.Action, flat.Action)
		s.Equal(event.Context.Source, flat.Source)
		s.Equal(event.Context.Timestamp.Unix(), flat.CreatedAt.Unix())
		s.Equal(event.Context.Timestamp.Unix(), flat.UpdatedAt.Unix())
		s.Equal(event.Subject.ID, flat.SubjectID)
		s.Equal(event.Subject.Name, flat.SubjectName)
		s.Equal(bytes, flat.Data)
	}

	var deflate defs.Event[defs.BranchOrTag, defs.RepoProvider]

	err = defs.Deflate(flat, &deflate)
	if s.NoError(err) {
		s.Equal(event.Version, deflate.Version)
		s.Equal(event.Context.ID, deflate.Context.ID)
		s.Equal(event.Context.Provider, deflate.Context.Provider)
		s.Equal(event.Context.Scope, deflate.Context.Scope)
		s.Equal(event.Context.Action, deflate.Context.Action)
		s.Equal(event.Context.Source, deflate.Context.Source)
		s.Equal(event.Context.Timestamp.Unix(), deflate.Context.Timestamp.Unix())
		s.Equal(event.Subject, deflate.Subject)
		s.Equal(event.Data, deflate.Data)
	}
}

func TestEventTestSuite(t *testing.T) {
	suite.Run(t, new(EventTestSuite))
}
