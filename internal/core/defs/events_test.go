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
		ID:     gocql.MustRandomUUID(),
		TeamID: gocql.MustRandomUUID(),
		Name:   "repos",
	}
}

func (s *EventTestSuite) Test_Branch_Create_MarshalJSON() {
	branch := &defs.BranchOrTag{
		Ref:           "test-branch",
		DefaultBranch: "main",
	}

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
    "name": "%s",
    "team_id": "%s"
  },
  "data": {
    "ref": "test-branch",
    "default_branch": "main"
  }
}`,
		event.Context.ID,
		event.Context.Timestamp.Format(time.RFC3339Nano),
		s.subject.ID,
		s.subject.Name,
		s.subject.TeamID,
	)

	// Test Marshal to JSON
	marshal, err := json.MarshalIndent(event, "", "  ")
	if s.NoError(err) {
		s.T().Log(string(marshal))
		s.Equal(defs.EventVersionV1, event.Version)
		s.Equal(uuid.Nil.String(), event.Context.ParentID.String())
		s.Equal(expected, string(marshal))
	}
}

func (s *EventTestSuite) Test_Branch_Create_UnmarshalJSON() {
	branch := &defs.BranchOrTag{
		Ref:           "test-branch",
		DefaultBranch: "main",
	}

	event := branch.ToEvent(
		defs.RepoProviderGithub,
		s.subject,
		defs.EventActionCreated,
		defs.EventScopeBranch,
	)
	event.SetSource("test/test")

	marshal, err := json.Marshal(event)
	s.Require().NoError(err)

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
}

func (s *EventTestSuite) Test_Branch_Create_Flatten() {
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
}

func (s *EventTestSuite) Test_Branch_Create_Deflate() {
	branch := &defs.BranchOrTag{
		Ref:           "test-branch",
		DefaultBranch: "main",
	}

	event := branch.ToEvent(
		defs.RepoProviderGithub,
		s.subject,
		defs.EventActionCreated,
		defs.EventScopeBranch,
	)
	event.SetSource("test/test")

	flat, err := event.Flatten()
	s.Require().NoError(err)

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

func (s *EventTestSuite) Test_Branch_Create_SetParent() {
	branch := &defs.BranchOrTag{
		Ref:           "test-branch",
		DefaultBranch: "main",
	}

	event := branch.ToEvent(
		defs.RepoProviderGithub,
		s.subject,
		defs.EventActionCreated,
		defs.EventScopeBranch,
	)

	// Test setting the parent ID
	event.SetParent(s.parent)
	s.Equal(s.parent, event.Context.ParentID, "Parent ID should be set correctly")
}

func (s *EventTestSuite) Test_Branch_Delete_MarshalJSON() {
	branch := &defs.BranchOrTag{
		Ref:           "test-branch",
		DefaultBranch: "main",
	}

	event := branch.ToEvent(
		defs.RepoProviderGithub,
		s.subject,
		defs.EventActionDeleted,
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
    "action": "deleted",
    "source": "test/test",
    "timestamp": "%s"
  },
  "subject": {
    "id": "%s",
    "name": "%s",
    "team_id": "%s"
  },
  "data": {
    "ref": "test-branch",
    "default_branch": "main"
  }
}`,
		event.Context.ID,
		event.Context.Timestamp.Format(time.RFC3339Nano),
		s.subject.ID,
		s.subject.Name,
		s.subject.TeamID,
	)

	// Test Marshal to JSON
	marshal, err := json.MarshalIndent(event, "", "  ")
	if s.NoError(err) {
		s.T().Log(string(marshal))
		s.Equal(defs.EventVersionV1, event.Version)
		s.Equal(uuid.Nil.String(), event.Context.ParentID.String())
		s.Equal(expected, string(marshal))
	}
}

func (s *EventTestSuite) Test_Branch_Delete_UnmarshalJSON() {
	branch := &defs.BranchOrTag{
		Ref:           "test-branch",
		DefaultBranch: "main",
	}

	event := branch.ToEvent(
		defs.RepoProviderGithub,
		s.subject,
		defs.EventActionDeleted,
		defs.EventScopeBranch,
	)
	event.SetSource("test/test")

	marshal, err := json.Marshal(event)
	s.Require().NoError(err)

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
}

func (s *EventTestSuite) Test_Branch_Delete_Flatten() {
	branch := &defs.BranchOrTag{
		Ref:           "test-branch",
		DefaultBranch: "main",
	}

	bytes, _ := json.Marshal(branch)

	event := branch.ToEvent(
		defs.RepoProviderGithub,
		s.subject,
		defs.EventActionDeleted,
		defs.EventScopeBranch,
	)
	event.SetSource("test/test")

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
}

func (s *EventTestSuite) Test_Branch_Delete_Deflate() {
	branch := &defs.BranchOrTag{
		Ref:           "test-branch",
		DefaultBranch: "main",
	}

	event := branch.ToEvent(
		defs.RepoProviderGithub,
		s.subject,
		defs.EventActionDeleted,
		defs.EventScopeBranch,
	)
	event.SetSource("test/test")

	flat, err := event.Flatten()
	s.Require().NoError(err)

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

func (s *EventTestSuite) Test_Branch_Delete_SetParent() {
	branch := &defs.BranchOrTag{
		Ref:           "test-branch",
		DefaultBranch: "main",
	}

	event := branch.ToEvent(
		defs.RepoProviderGithub,
		s.subject,
		defs.EventActionDeleted,
		defs.EventScopeBranch,
	)

	// Test setting the parent ID
	event.SetParent(s.parent)
	s.Equal(s.parent, event.Context.ParentID, "Parent ID should be set correctly")
}

func (s *EventTestSuite) Test_Push_Create_MarshalJSON() {
	push := &defs.Push{
		Ref:        "refs/heads/test-branch",
		Before:     "old-commit-sha",
		After:      "a1b2c3d4e5f678901234567890abcdef12345678", // Valid Git hash
		Repository: "test/test",
		Pusher:     "testuser",
		Commits: []defs.Commit{
			{
				SHA:       "a1b2c3d4e5f678901234567890abcdef12345678", // Valid Git hash
				Message:   "Test commit message",
				Author:    "testuser",
				Committer: "testuser",
				Timestamp: time.Now(),
				URL:       "https://github.com/test/test/commit/a1b2c3d4e5f678901234567890abcdef12345678",
				Added:     []string{"new-file.txt"},
				Removed:   []string{},
				Modified:  []string{"modified-file.txt"},
			},
		},
		Timestamp: time.Now(),
	}

	event := push.ToEvent(defs.RepoProviderGithub, s.subject, defs.EventActionCreated)
	event.SetSource("test/test")

	expected := fmt.Sprintf(`{
  "version": "v1",
  "context": {
    "id": "%s",
    "parent_id": "00000000-0000-0000-0000-000000000000",
    "provider": "github",
    "scope": "push",
    "action": "created",
    "source": "test/test",
    "timestamp": "%s"
  },
  "subject": {
    "id": "%s",
    "name": "%s",
    "team_id": "%s"
  },
  "data": {
    "ref": "refs/heads/test-branch",
    "before": "old-commit-sha",
    "after": "a1b2c3d4e5f678901234567890abcdef12345678",
    "repository": "test/test",
    "pusher": "testuser",
    "commits": [
      {
        "sha": "a1b2c3d4e5f678901234567890abcdef12345678",
        "message": "Test commit message",
        "author": "testuser",
        "committer": "testuser",
        "timestamp": "%s",
        "url": "https://github.com/test/test/commit/a1b2c3d4e5f678901234567890abcdef12345678",
        "added": [
          "new-file.txt"
        ],
        "removed": [],
        "modified": [
          "modified-file.txt"
        ]
      }
    ],
    "timestamp": "%s"
  }
}`,
		event.Context.ID,
		event.Context.Timestamp.Format(time.RFC3339Nano),
		s.subject.ID,
		s.subject.Name,
		s.subject.TeamID,
		push.Commits[0].Timestamp.Format(time.RFC3339Nano),
		push.Timestamp.Format(time.RFC3339Nano),
	)

	// Test Marshal to JSON
	marshal, err := json.MarshalIndent(event, "", "  ")
	if s.NoError(err) {
		s.T().Log(string(marshal))
		s.Equal(defs.EventVersionV1, event.Version)
		s.Equal(uuid.Nil.String(), event.Context.ParentID.String())
		s.Equal(expected, string(marshal))
	}
}

func (s *EventTestSuite) Test_Push_Create_UnmarshalJSON() {
	push := &defs.Push{
		Ref:        "refs/heads/test-branch",
		Before:     "old-commit-sha",
		After:      "a1b2c3d4e5f678901234567890abcdef12345678", // Valid Git hash
		Repository: "test/test",
		Pusher:     "testuser",
		Commits: []defs.Commit{
			{
				SHA:       "a1b2c3d4e5f678901234567890abcdef12345678", // Valid Git hash
				Message:   "Test commit message",
				Author:    "testuser",
				Committer: "testuser",
				Timestamp: time.Now(),
				URL:       "https://github.com/test/test/commit/a1b2c3d4e5f678901234567890abcdef12345678",
				Added:     []string{"new-file.txt"},
				Removed:   []string{},
				Modified:  []string{"modified-file.txt"},
			},
		},
		Timestamp: time.Now(),
	}

	event := push.ToEvent(defs.RepoProviderGithub, s.subject, defs.EventActionCreated)
	event.SetSource("test/test")

	marshal, err := json.Marshal(event)
	s.Require().NoError(err)

	// Test Unmarshal from JSON
	var unmarshal defs.Event[defs.Push, defs.RepoProvider]
	err = json.Unmarshal(marshal, &unmarshal)

	if s.NoError(err) {
		s.Equal(event.Version, unmarshal.Version)
		s.Equal(event.Context.ID, unmarshal.Context.ID)
		s.Equal(event.Context.ParentID, unmarshal.Context.ParentID)
		s.Equal(event.Context.Provider, unmarshal.Context.Provider)
		s.Equal(event.Context.Scope, unmarshal.Context.Scope)
		s.Equal(event.Context.Action, unmarshal.Context.Action)
		s.Equal(event.Context.Source, unmarshal.Context.Source)
		s.WithinDuration(event.Context.Timestamp, unmarshal.Context.Timestamp, time.Second)
		s.Equal(event.Subject.ID, unmarshal.Subject.ID)
		s.Equal(event.Subject.Name, unmarshal.Subject.Name)
		s.Equal(event.Subject.TeamID, unmarshal.Subject.TeamID)
		s.Equal(event.Data.Ref, unmarshal.Data.Ref)
		s.Equal(event.Data.Before, unmarshal.Data.Before)
		s.Equal(event.Data.After, unmarshal.Data.After)
		s.Equal(event.Data.Repository, unmarshal.Data.Repository)
		s.Equal(event.Data.Pusher, unmarshal.Data.Pusher)
		s.Equal(len(event.Data.Commits), len(unmarshal.Data.Commits))

		for i := range event.Data.Commits {
			s.Equal(event.Data.Commits[i].SHA, unmarshal.Data.Commits[i].SHA)
			s.Equal(event.Data.Commits[i].Message, unmarshal.Data.Commits[i].Message)
			s.Equal(event.Data.Commits[i].Author, unmarshal.Data.Commits[i].Author)
			s.Equal(event.Data.Commits[i].Committer, unmarshal.Data.Commits[i].Committer)
			s.WithinDuration(event.Data.Commits[i].Timestamp, unmarshal.Data.Commits[i].Timestamp, time.Second)
			s.Equal(event.Data.Commits[i].URL, unmarshal.Data.Commits[i].URL)
			s.Equal(event.Data.Commits[i].Added, unmarshal.Data.Commits[i].Added)
			s.Equal(event.Data.Commits[i].Removed, unmarshal.Data.Commits[i].Removed)
			s.Equal(event.Data.Commits[i].Modified, unmarshal.Data.Commits[i].Modified)
		}

		s.WithinDuration(event.Data.Timestamp, unmarshal.Data.Timestamp, time.Second)
	}
}

func (s *EventTestSuite) Test_Push_Create_Flatten() {
	push := &defs.Push{
		Ref:        "refs/heads/test-branch",
		Before:     "old-commit-sha",
		After:      "a1b2c3d4e5f678901234567890abcdef12345678", // Valid Git hash
		Repository: "test/test",
		Pusher:     "testuser",
		Commits: []defs.Commit{
			{
				SHA:       "a1b2c3d4e5f678901234567890abcdef12345678", // Valid Git hash
				Message:   "Test commit message",
				Author:    "testuser",
				Committer: "testuser",
				Timestamp: time.Now(),
				URL:       "https://github.com/test/test/commit/a1b2c3d4e5f678901234567890abcdef12345678",
				Added:     []string{"new-file.txt"},
				Removed:   []string{},
				Modified:  []string{"modified-file.txt"},
			},
		},
		Timestamp: time.Now(),
	}

	bytes, _ := json.Marshal(push)

	event := push.ToEvent(defs.RepoProviderGithub, s.subject, defs.EventActionCreated)
	event.SetSource("test/test")

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
}

func (s *EventTestSuite) Test_Push_Create_Deflate() {
	push := &defs.Push{
		Ref:        "refs/heads/test-branch",
		Before:     "old-commit-sha",
		After:      "a1b2c3d4e5f678901234567890abcdef12345678", // Valid Git hash
		Repository: "test/test",
		Pusher:     "testuser",
		Commits: []defs.Commit{
			{
				SHA:       "a1b2c3d4e5f678901234567890abcdef12345678", // Valid Git hash
				Message:   "Test commit message",
				Author:    "testuser",
				Committer: "testuser",
				Timestamp: time.Now(),
				URL:       "https://github.com/test/test/commit/a1b2c3d4e5f678901234567890abcdef12345678",
				Added:     []string{"new-file.txt"},
				Removed:   []string{},
				Modified:  []string{"modified-file.txt"},
			},
		},
		Timestamp: time.Now(),
	}

	event := push.ToEvent(defs.RepoProviderGithub, s.subject, defs.EventActionCreated)
	event.SetSource("test/test")

	flat, err := event.Flatten()
	s.Require().NoError(err)

	var deflate defs.Event[defs.Push, defs.RepoProvider]

	err = defs.Deflate(flat, &deflate)
	if s.NoError(err) {
		s.Equal(event.Version, deflate.Version)
		s.Equal(event.Context.ID, deflate.Context.ID)
		s.Equal(event.Context.ParentID, deflate.Context.ParentID)
		s.Equal(event.Context.Provider, deflate.Context.Provider)
		s.Equal(event.Context.Scope, deflate.Context.Scope)
		s.Equal(event.Context.Action, deflate.Context.Action)
		s.Equal(event.Context.Source, deflate.Context.Source)
		s.WithinDuration(event.Context.Timestamp, deflate.Context.Timestamp, time.Second)
		s.Equal(event.Subject.ID, deflate.Subject.ID)
		s.Equal(event.Subject.Name, deflate.Subject.Name)
		s.Equal(event.Subject.TeamID, deflate.Subject.TeamID)

		// Compare Push data fields individually
		s.Equal(event.Data.Ref, deflate.Data.Ref)
		s.Equal(event.Data.Before, deflate.Data.Before)
		s.Equal(event.Data.After, deflate.Data.After)
		s.Equal(event.Data.Repository, deflate.Data.Repository)
		s.Equal(event.Data.Pusher, deflate.Data.Pusher)
		s.Equal(len(event.Data.Commits), len(deflate.Data.Commits)) // Compare number of commits

		// Compare Commits within the array
		for i := range event.Data.Commits {
			s.Equal(event.Data.Commits[i].SHA, deflate.Data.Commits[i].SHA)
			s.Equal(event.Data.Commits[i].Message, deflate.Data.Commits[i].Message)
			s.Equal(event.Data.Commits[i].Author, deflate.Data.Commits[i].Author)
			s.Equal(event.Data.Commits[i].Committer, deflate.Data.Commits[i].Committer)
			s.WithinDuration(event.Data.Commits[i].Timestamp, deflate.Data.Commits[i].Timestamp, time.Second)
			s.Equal(event.Data.Commits[i].URL, deflate.Data.Commits[i].URL)
			s.Equal(event.Data.Commits[i].Added, deflate.Data.Commits[i].Added)
			s.Equal(event.Data.Commits[i].Removed, deflate.Data.Commits[i].Removed)
			s.Equal(event.Data.Commits[i].Modified, deflate.Data.Commits[i].Modified)
		}

		// Compare timestamps using WithinDuration for flexibility
		s.WithinDuration(event.Data.Timestamp, deflate.Data.Timestamp, time.Second)
	}
}

func (s *EventTestSuite) Test_Push_Create_SetParent() {
	push := &defs.Push{
		Ref:        "refs/heads/test-branch",
		Before:     "old-commit-sha",
		After:      "a1b2c3d4e5f678901234567890abcdef12345678", // Valid Git hash
		Repository: "test/test",
		Pusher:     "testuser",
		Commits: []defs.Commit{
			{
				SHA:       "a1b2c3d4e5f678901234567890abcdef12345678", // Valid Git hash
				Message:   "Test commit message",
				Author:    "testuser",
				Committer: "testuser",
				Timestamp: time.Now(),
				URL:       "https://github.com/test/test/commit/a1b2c3d4e5f678901234567890abcdef12345678",
				Added:     []string{"new-file.txt"},
				Removed:   []string{},
				Modified:  []string{"modified-file.txt"},
			},
		},
		Timestamp: time.Now(),
	}

	event := push.ToEvent(defs.RepoProviderGithub, s.subject, defs.EventActionCreated)

	// Test setting the parent ID
	event.SetParent(s.parent)
	s.Equal(s.parent, event.Context.ParentID, "Parent ID should be set correctly")
}

func (s *EventTestSuite) Test_PullRequest_Create_MarshalJSON() {
	pr := &defs.PullRequest{
		Number:         1,
		Title:          "Test Pull Request",
		Body:           "This is a test pull request",
		State:          "open",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ClosedAt:       time.Time{},
		MergedAt:       time.Time{},
		MergeCommitSHA: "a1b2c3d4e5f678901234567890abcdef12345678", // Valid Git hash
		Author:         "testuser",
		HeadBranch:     "test-branch",
		BaseBranch:     "main",
	}

	event := pr.ToEvent(defs.RepoProviderGithub, s.subject, defs.EventActionCreated)
	event.SetSource("test/test")

	expected := fmt.Sprintf(`{
  "version": "v1",
  "context": {
    "id": "%s",
    "parent_id": "00000000-0000-0000-0000-000000000000",
    "provider": "github",
    "scope": "pull_request",
    "action": "created",
    "source": "test/test",
    "timestamp": "%s"
  },
  "subject": {
    "id": "%s",
    "name": "%s",
    "team_id": "%s"
  },
  "data": {
    "number": 1,
    "title": "Test Pull Request",
    "body": "This is a test pull request",
    "state": "open",
    "created_at": "%s",
    "updated_at": "%s",
    "closed_at": "0001-01-01T00:00:00Z",
    "merged_at": "0001-01-01T00:00:00Z",
    "merge_commit_sha": "a1b2c3d4e5f678901234567890abcdef12345678",
    "author": "testuser",
    "head_branch": "test-branch",
    "base_branch": "main"
  }
}`,
		event.Context.ID,
		event.Context.Timestamp.Format(time.RFC3339Nano),
		s.subject.ID,
		s.subject.Name,
		s.subject.TeamID,
		pr.CreatedAt.Format(time.RFC3339Nano),
		pr.UpdatedAt.Format(time.RFC3339Nano),
	)

	// Test Marshal to JSON
	marshal, err := json.MarshalIndent(event, "", "  ")
	if s.NoError(err) {
		s.T().Log(string(marshal))
		s.Equal(defs.EventVersionV1, event.Version)
		s.Equal(uuid.Nil.String(), event.Context.ParentID.String())
		s.Equal(expected, string(marshal))
	}
}

func (s *EventTestSuite) Test_PullRequest_Create_UnmarshalJSON() {
	pr := &defs.PullRequest{
		Number:         1,
		Title:          "Test Pull Request",
		Body:           "This is a test pull request",
		State:          "open",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ClosedAt:       time.Time{},
		MergedAt:       time.Time{},
		MergeCommitSHA: "a1b2c3d4e5f678901234567890abcdef12345678", // Valid Git hash
		Author:         "testuser",
		HeadBranch:     "test-branch",
		BaseBranch:     "main",
	}

	event := pr.ToEvent(defs.RepoProviderGithub, s.subject, defs.EventActionCreated)
	event.SetSource("test/test")

	marshal, err := json.Marshal(event)
	s.Require().NoError(err)

	// Test Unmarshal from JSON
	var unmarshal defs.Event[defs.PullRequest, defs.RepoProvider]
	err = json.Unmarshal(marshal, &unmarshal)

	if s.NoError(err) {
		s.Equal(event.Version, unmarshal.Version)
		s.Equal(event.Context.ID, unmarshal.Context.ID)
		s.Equal(event.Context.ParentID, unmarshal.Context.ParentID)
		s.Equal(event.Context.Provider, unmarshal.Context.Provider)
		s.Equal(event.Context.Scope, unmarshal.Context.Scope)
		s.Equal(event.Context.Action, unmarshal.Context.Action)
		s.Equal(event.Context.Source, unmarshal.Context.Source)
		s.WithinDuration(event.Context.Timestamp, unmarshal.Context.Timestamp, time.Second)
		s.Equal(event.Subject.ID, unmarshal.Subject.ID)
		s.Equal(event.Subject.Name, unmarshal.Subject.Name)
		s.Equal(event.Subject.TeamID, unmarshal.Subject.TeamID)
		s.Equal(event.Data.Number, unmarshal.Data.Number)
		s.Equal(event.Data.Title, unmarshal.Data.Title)
		s.Equal(event.Data.Body, unmarshal.Data.Body)
		s.Equal(event.Data.State, unmarshal.Data.State)
		s.WithinDuration(event.Data.CreatedAt, unmarshal.Data.CreatedAt, time.Second)
		s.WithinDuration(event.Data.UpdatedAt, unmarshal.Data.UpdatedAt, time.Second)
		s.Equal(event.Data.ClosedAt.Unix(), unmarshal.Data.ClosedAt.Unix())
		s.Equal(event.Data.MergedAt.Unix(), unmarshal.Data.MergedAt.Unix())
		s.Equal(event.Data.MergeCommitSHA, unmarshal.Data.MergeCommitSHA)
		s.Equal(event.Data.Author, unmarshal.Data.Author)
		s.Equal(event.Data.HeadBranch, unmarshal.Data.HeadBranch)
		s.Equal(event.Data.BaseBranch, unmarshal.Data.BaseBranch)
	}
}

func (s *EventTestSuite) Test_PullRequest_Create_Deflate() {
	pr := &defs.PullRequest{
		Number:         1,
		Title:          "Test Pull Request",
		Body:           "This is a test pull request",
		State:          "open",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		ClosedAt:       time.Time{},
		MergedAt:       time.Time{},
		MergeCommitSHA: "a1b2c3d4e5f678901234567890abcdef12345678", // Valid Git hash
		Author:         "testuser",
		HeadBranch:     "test-branch",
		BaseBranch:     "main",
	}

	event := pr.ToEvent(defs.RepoProviderGithub, s.subject, defs.EventActionCreated)
	event.SetSource("test/test")

	flat, err := event.Flatten()
	s.Require().NoError(err)

	var deflate defs.Event[defs.PullRequest, defs.RepoProvider]

	err = defs.Deflate(flat, &deflate)
	if s.NoError(err) {
		s.Equal(event.Version, deflate.Version)
		s.Equal(event.Context.ID, deflate.Context.ID)
		s.Equal(event.Context.ParentID, deflate.Context.ParentID)
		s.Equal(event.Context.Provider, deflate.Context.Provider)
		s.Equal(event.Context.Scope, deflate.Context.Scope)
		s.Equal(event.Context.Action, deflate.Context.Action)
		s.Equal(event.Context.Source, deflate.Context.Source)
		s.WithinDuration(event.Context.Timestamp, deflate.Context.Timestamp, time.Second)
		s.Equal(event.Subject.ID, deflate.Subject.ID)
		s.Equal(event.Subject.Name, deflate.Subject.Name)
		s.Equal(event.Subject.TeamID, deflate.Subject.TeamID)
		s.Equal(event.Data.Number, deflate.Data.Number)
		s.Equal(event.Data.Title, deflate.Data.Title)
		s.Equal(event.Data.Body, deflate.Data.Body)
		s.Equal(event.Data.State, deflate.Data.State)
		s.WithinDuration(event.Data.CreatedAt, deflate.Data.CreatedAt, time.Second)
		s.WithinDuration(event.Data.UpdatedAt, deflate.Data.UpdatedAt, time.Second)
		s.Equal(event.Data.ClosedAt.Unix(), deflate.Data.ClosedAt.Unix())
		s.Equal(event.Data.MergedAt.Unix(), deflate.Data.MergedAt.Unix())
		s.Equal(event.Data.MergeCommitSHA, deflate.Data.MergeCommitSHA)
		s.Equal(event.Data.Author, deflate.Data.Author)
		s.Equal(event.Data.HeadBranch, deflate.Data.HeadBranch)
		s.Equal(event.Data.BaseBranch, deflate.Data.BaseBranch)
	}
}

func TestEventTestSuite(t *testing.T) {
	suite.Run(t, new(EventTestSuite))
}
