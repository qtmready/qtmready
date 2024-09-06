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
		sha     string
	}
)

func (s *EventTestSuite) SetupSuite() {
	shared.InitServiceForTest()

	s.parent = gocql.TimeUUID()
	s.sha = "a1b2c3d4e5f678901234567890abcdef12345678"
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
  "version": "0.1.0",
  "id": "%s",
  "context": {
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
  "payload": {
    "ref": "test-branch",
    "default_branch": "main"
  }
}`,
		event.ID,
		event.Context.Timestamp.Format(time.RFC3339Nano),
		s.subject.ID,
		s.subject.Name,
		s.subject.TeamID,
	)

	// Test Marshal to JSON
	marshal, err := json.MarshalIndent(event, "", "  ")
	if s.NoError(err) {
		s.T().Log(string(marshal))
		s.Equal(defs.EventVersionDefault, event.Version)
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
		s.Equal(event.ID, unmarshal.ID)
		s.Equal(event.Context.Provider, unmarshal.Context.Provider)
		s.Equal(event.Context.Scope, unmarshal.Context.Scope)
		s.Equal(event.Context.Action, unmarshal.Context.Action)
		s.Equal(event.Context.Source, unmarshal.Context.Source)
		s.Equal(event.Context.Timestamp.Unix(), unmarshal.Context.Timestamp.Unix())
		s.Equal(event.Subject, unmarshal.Subject)
		s.Equal(event.Payload, unmarshal.Payload)
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
		s.Equal(event.ID, flat.ID)
		s.Equal(event.Context.Provider, flat.Provider)
		s.Equal(event.Context.Scope, flat.Scope)
		s.Equal(event.Context.Action, flat.Action)
		s.Equal(event.Context.Source, flat.Source)
		s.Equal(event.Context.Timestamp.Unix(), flat.CreatedAt.Unix())
		s.Equal(event.Context.Timestamp.Unix(), flat.UpdatedAt.Unix())
		s.Equal(event.Subject.ID, flat.SubjectID)
		s.Equal(event.Subject.Name, flat.SubjectName)
		s.Equal(bytes, flat.Payload)
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
		s.Equal(event.ID, deflate.ID)
		s.Equal(event.Context.Provider, deflate.Context.Provider)
		s.Equal(event.Context.Scope, deflate.Context.Scope)
		s.Equal(event.Context.Action, deflate.Context.Action)
		s.Equal(event.Context.Source, deflate.Context.Source)
		s.Equal(event.Context.Timestamp.Unix(), deflate.Context.Timestamp.Unix())
		s.Equal(event.Subject, deflate.Subject)
		s.Equal(event.Payload, deflate.Payload)
	}
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
  "version": "0.1.0",
  "id": "%s",
  "context": {
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
  "payload": {
    "ref": "test-branch",
    "default_branch": "main"
  }
}`,
		event.ID,
		event.Context.Timestamp.Format(time.RFC3339Nano),
		s.subject.ID,
		s.subject.Name,
		s.subject.TeamID,
	)

	// Test Marshal to JSON
	marshal, err := json.MarshalIndent(event, "", "  ")
	if s.NoError(err) {
		s.T().Log(string(marshal))
		s.Equal(defs.EventVersionDefault, event.Version)
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
		s.Equal(event.ID, unmarshal.ID)
		s.Equal(event.Context.Provider, unmarshal.Context.Provider)
		s.Equal(event.Context.Scope, unmarshal.Context.Scope)
		s.Equal(event.Context.Action, unmarshal.Context.Action)
		s.Equal(event.Context.Source, unmarshal.Context.Source)
		s.Equal(event.Context.Timestamp.Unix(), unmarshal.Context.Timestamp.Unix())
		s.Equal(event.Subject, unmarshal.Subject)
		s.Equal(event.Payload, unmarshal.Payload)
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
		s.Equal(event.ID, flat.ID)
		s.Equal(event.Context.Provider, flat.Provider)
		s.Equal(event.Context.Scope, flat.Scope)
		s.Equal(event.Context.Action, flat.Action)
		s.Equal(event.Context.Source, flat.Source)
		s.Equal(event.Context.Timestamp.Unix(), flat.CreatedAt.Unix())
		s.Equal(event.Context.Timestamp.Unix(), flat.UpdatedAt.Unix())
		s.Equal(event.Subject.ID, flat.SubjectID)
		s.Equal(event.Subject.Name, flat.SubjectName)
		s.Equal(bytes, flat.Payload)
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
		s.Equal(event.ID, deflate.ID)
		s.Equal(event.Context.Provider, deflate.Context.Provider)
		s.Equal(event.Context.Scope, deflate.Context.Scope)
		s.Equal(event.Context.Action, deflate.Context.Action)
		s.Equal(event.Context.Source, deflate.Context.Source)
		s.Equal(event.Context.Timestamp.Unix(), deflate.Context.Timestamp.Unix())
		s.Equal(event.Subject, deflate.Subject)
		s.Equal(event.Payload, deflate.Payload)
	}
}

func (s *EventTestSuite) Test_Tag_Create_MarshalJSON() {
	tag := &defs.BranchOrTag{
		Ref:           "v1.0.0",
		DefaultBranch: "main",
	}

	event := tag.ToEvent(
		defs.RepoProviderGithub,
		s.subject,
		defs.EventActionCreated,
		defs.EventScopeTag,
	)
	event.SetSource("test/test")

	expected := fmt.Sprintf(`{
  "version": "0.1.0",
  "id": "%s",
  "context": {
    "parent_id": "00000000-0000-0000-0000-000000000000",
    "provider": "github",
    "scope": "tag",
    "action": "created",
    "source": "test/test",
    "timestamp": "%s"
  },
  "subject": {
    "id": "%s",
    "name": "%s",
    "team_id": "%s"
  },
  "payload": {
    "ref": "v1.0.0",
    "default_branch": "main"
  }
}`,
		event.ID,
		event.Context.Timestamp.Format(time.RFC3339Nano),
		s.subject.ID,
		s.subject.Name,
		s.subject.TeamID,
	)

	// Test Marshal to JSON
	marshal, err := json.MarshalIndent(event, "", "  ")
	if s.NoError(err) {
		s.T().Log(string(marshal))
		s.Equal(defs.EventVersionDefault, event.Version)
		s.Equal(uuid.Nil.String(), event.Context.ParentID.String())
		s.Equal(expected, string(marshal))
	}
}

func (s *EventTestSuite) Test_Tag_Create_UnmarshalJSON() {
	tag := &defs.BranchOrTag{
		Ref:           "v1.0.0",
		DefaultBranch: "main",
	}

	event := tag.ToEvent(
		defs.RepoProviderGithub,
		s.subject,
		defs.EventActionCreated,
		defs.EventScopeTag,
	)
	event.SetSource("test/test")

	marshal, err := json.Marshal(event)
	s.Require().NoError(err)

	// Test Unmarshal from JSON
	var unmarshal defs.Event[defs.BranchOrTag, defs.RepoProvider]
	err = json.Unmarshal(marshal, &unmarshal)

	if s.NoError(err) {
		s.Equal(event.Version, unmarshal.Version)
		s.Equal(event.ID, unmarshal.ID)
		s.Equal(event.Context.Provider, unmarshal.Context.Provider)
		s.Equal(event.Context.Scope, unmarshal.Context.Scope)
		s.Equal(event.Context.Action, unmarshal.Context.Action)
		s.Equal(event.Context.Source, unmarshal.Context.Source)
		s.Equal(event.Context.Timestamp.Unix(), unmarshal.Context.Timestamp.Unix())
		s.Equal(event.Subject, unmarshal.Subject)
		s.Equal(event.Payload, unmarshal.Payload)
	}
}

func (s *EventTestSuite) Test_Tag_Create_Flatten() {
	tag := &defs.BranchOrTag{
		Ref:           "v1.0.0",
		DefaultBranch: "main",
	}

	bytes, _ := json.Marshal(tag)

	event := tag.ToEvent(
		defs.RepoProviderGithub,
		s.subject,
		defs.EventActionCreated,
		defs.EventScopeTag,
	)
	event.SetSource("test/test")

	// Test flattening the event
	flat, err := event.Flatten()
	if s.NoError(err) {
		s.Equal(event.Version, flat.Version)
		s.Equal(event.ID, flat.ID)
		s.Equal(event.Context.Provider, flat.Provider)
		s.Equal(event.Context.Scope, flat.Scope)
		s.Equal(event.Context.Action, flat.Action)
		s.Equal(event.Context.Source, flat.Source)
		s.Equal(event.Context.Timestamp.Unix(), flat.CreatedAt.Unix())
		s.Equal(event.Context.Timestamp.Unix(), flat.UpdatedAt.Unix())
		s.Equal(event.Subject.ID, flat.SubjectID)
		s.Equal(event.Subject.Name, flat.SubjectName)
		s.Equal(bytes, flat.Payload)
	}
}

func (s *EventTestSuite) Test_Tag_Create_Deflate() {
	tag := &defs.BranchOrTag{
		Ref:           "v1.0.0",
		DefaultBranch: "main",
	}

	event := tag.ToEvent(
		defs.RepoProviderGithub,
		s.subject,
		defs.EventActionCreated,
		defs.EventScopeTag,
	)
	event.SetSource("test/test")

	flat, err := event.Flatten()
	s.Require().NoError(err)

	var deflate defs.Event[defs.BranchOrTag, defs.RepoProvider]

	err = defs.Deflate(flat, &deflate)
	if s.NoError(err) {
		s.Equal(event.Version, deflate.Version)
		s.Equal(event.ID, deflate.ID)
		s.Equal(event.Context.Provider, deflate.Context.Provider)
		s.Equal(event.Context.Scope, deflate.Context.Scope)
		s.Equal(event.Context.Action, deflate.Context.Action)
		s.Equal(event.Context.Source, deflate.Context.Source)
		s.Equal(event.Context.Timestamp.Unix(), deflate.Context.Timestamp.Unix())
		s.Equal(event.Subject, deflate.Subject)
		s.Equal(event.Payload, deflate.Payload)
	}
}

func (s *EventTestSuite) Test_PullRequest_Create_MarshalJSON() {
	pr := &defs.PullRequest{
		Number:         1,
		Title:          "Test Pull Request",
		Body:           "This is a test pull request",
		State:          "open",
		MergeCommitSHA: &s.sha,
		Author:         "testuser",
		HeadBranch:     "test-branch",
		BaseBranch:     "main",
		Timestamp:      time.Now(),
	}

	event := pr.ToEvent(defs.RepoProviderGithub, s.subject, defs.EventActionCreated)
	event.SetSource("test/test")

	expected := fmt.Sprintf(`{
  "version": "0.1.0",
  "id": "%s",
  "context": {
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
  "payload": {
    "number": 1,
    "title": "Test Pull Request",
    "body": "This is a test pull request",
    "state": "open",
    "merge_commit_sha": "a1b2c3d4e5f678901234567890abcdef12345678",
    "author": "testuser",
    "head_branch": "test-branch",
    "base_branch": "main",
    "timestamp": "%s"
  }
}`,
		event.ID,
		event.Context.Timestamp.Format(time.RFC3339Nano),
		s.subject.ID,
		s.subject.Name,
		s.subject.TeamID,
		pr.Timestamp.Format(time.RFC3339Nano),
	)

	// Test Marshal to JSON
	marshal, err := json.MarshalIndent(event, "", "  ")
	if s.NoError(err) {
		s.T().Log(string(marshal))
		s.Equal(defs.EventVersionDefault, event.Version)
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
		MergeCommitSHA: &s.sha,
		Author:         "testuser",
		HeadBranch:     "test-branch",
		BaseBranch:     "main",
		Timestamp:      time.Now(),
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
		s.Equal(event.ID, unmarshal.ID)
		s.Equal(event.Context.ParentID, unmarshal.Context.ParentID)
		s.Equal(event.Context.Provider, unmarshal.Context.Provider)
		s.Equal(event.Context.Scope, unmarshal.Context.Scope)
		s.Equal(event.Context.Action, unmarshal.Context.Action)
		s.Equal(event.Context.Source, unmarshal.Context.Source)
		s.WithinDuration(event.Context.Timestamp, unmarshal.Context.Timestamp, time.Second)
		s.Equal(event.Subject.ID, unmarshal.Subject.ID)
		s.Equal(event.Subject.Name, unmarshal.Subject.Name)
		s.Equal(event.Subject.TeamID, unmarshal.Subject.TeamID)
		s.Equal(event.Payload.Number, unmarshal.Payload.Number)
		s.Equal(event.Payload.Title, unmarshal.Payload.Title)
		s.Equal(event.Payload.Body, unmarshal.Payload.Body)
		s.Equal(event.Payload.State, unmarshal.Payload.State)
		s.WithinDuration(event.Payload.Timestamp, unmarshal.Payload.Timestamp, time.Second)
		s.Equal(event.Payload.MergeCommitSHA, unmarshal.Payload.MergeCommitSHA)
		s.Equal(event.Payload.Author, unmarshal.Payload.Author)
		s.Equal(event.Payload.HeadBranch, unmarshal.Payload.HeadBranch)
		s.Equal(event.Payload.BaseBranch, unmarshal.Payload.BaseBranch)
	}
}

func (s *EventTestSuite) Test_PullRequest_Create_Deflate() {
	pr := &defs.PullRequest{
		Number:         1,
		Title:          "Test Pull Request",
		Body:           "This is a test pull request",
		State:          "open",
		MergeCommitSHA: &s.sha,
		Author:         "testuser",
		HeadBranch:     "test-branch",
		BaseBranch:     "main",
		Timestamp:      time.Now(),
	}

	event := pr.ToEvent(defs.RepoProviderGithub, s.subject, defs.EventActionCreated)
	event.SetSource("test/test")

	flat, err := event.Flatten()
	s.Require().NoError(err)

	var deflate defs.Event[defs.PullRequest, defs.RepoProvider]

	err = defs.Deflate(flat, &deflate)
	if s.NoError(err) {
		s.Equal(event.Version, deflate.Version)
		s.Equal(event.ID, deflate.ID)
		s.Equal(event.Context.ParentID, deflate.Context.ParentID)
		s.Equal(event.Context.Provider, deflate.Context.Provider)
		s.Equal(event.Context.Scope, deflate.Context.Scope)
		s.Equal(event.Context.Action, deflate.Context.Action)
		s.Equal(event.Context.Source, deflate.Context.Source)
		s.WithinDuration(event.Context.Timestamp, deflate.Context.Timestamp, time.Second)
		s.Equal(event.Subject.ID, deflate.Subject.ID)
		s.Equal(event.Subject.Name, deflate.Subject.Name)
		s.Equal(event.Subject.TeamID, deflate.Subject.TeamID)
		s.Equal(event.Payload.Number, deflate.Payload.Number)
		s.Equal(event.Payload.Title, deflate.Payload.Title)
		s.Equal(event.Payload.Body, deflate.Payload.Body)
		s.Equal(event.Payload.State, deflate.Payload.State)
		s.WithinDuration(event.Payload.Timestamp, deflate.Payload.Timestamp, time.Second)
		s.Equal(event.Payload.MergeCommitSHA, deflate.Payload.MergeCommitSHA)
		s.Equal(event.Payload.Author, deflate.Payload.Author)
		s.Equal(event.Payload.HeadBranch, deflate.Payload.HeadBranch)
		s.Equal(event.Payload.BaseBranch, deflate.Payload.BaseBranch)
	}
}

func (s *EventTestSuite) Test_PullRequestLabel_Create_UnmarshalJSON() {
	label := &defs.PullRequestLabel{
		Name:              "bug",
		PullRequestNumber: 123,
		Branch:            "main",
		Timestamp:         time.Now(),
	}

	event := label.ToEvent(defs.RepoProviderGithub, s.subject, defs.EventActionCreated)
	event.SetSource("test/test")

	marshal, err := json.Marshal(event)
	s.Require().NoError(err)

	// Test Unmarshal from JSON
	var unmarshal defs.Event[defs.PullRequestLabel, defs.RepoProvider]
	err = json.Unmarshal(marshal, &unmarshal)

	if s.NoError(err) {
		s.Equal(event.Version, unmarshal.Version)
		s.Equal(event.ID, unmarshal.ID)
		s.Equal(event.Context.ParentID, unmarshal.Context.ParentID)
		s.Equal(event.Context.Provider, unmarshal.Context.Provider)
		s.Equal(event.Context.Scope, unmarshal.Context.Scope)
		s.Equal(event.Context.Action, unmarshal.Context.Action)
		s.Equal(event.Context.Source, unmarshal.Context.Source)
		s.WithinDuration(event.Context.Timestamp, unmarshal.Context.Timestamp, time.Second)
		s.Equal(event.Subject.ID, unmarshal.Subject.ID)
		s.Equal(event.Subject.Name, unmarshal.Subject.Name)
		s.Equal(event.Subject.TeamID, unmarshal.Subject.TeamID)
		s.Equal(event.Payload.Name, unmarshal.Payload.Name)
		s.Equal(event.Payload.PullRequestNumber, unmarshal.Payload.PullRequestNumber)
		s.WithinDuration(event.Payload.Timestamp, unmarshal.Payload.Timestamp, time.Second)
	}
}

func (s *EventTestSuite) Test_PullRequestLabel_Create_Flatten() {
	label := &defs.PullRequestLabel{
		Name:              "bug",
		PullRequestNumber: 123,
		Branch:            "main",
		Timestamp:         time.Now(),
	}

	bytes, _ := json.Marshal(label)

	event := label.ToEvent(defs.RepoProviderGithub, s.subject, defs.EventActionCreated)
	event.SetSource("test/test")

	// Test flattening the event
	flat, err := event.Flatten()
	if s.NoError(err) {
		s.Equal(event.Version, flat.Version)
		s.Equal(event.ID, flat.ID)
		s.Equal(event.Context.Provider, flat.Provider)
		s.Equal(event.Context.Scope, flat.Scope)
		s.Equal(event.Context.Action, flat.Action)
		s.Equal(event.Context.Source, flat.Source)
		s.Equal(event.Context.Timestamp.Unix(), flat.CreatedAt.Unix())
		s.Equal(event.Context.Timestamp.Unix(), flat.UpdatedAt.Unix())
		s.Equal(event.Subject.ID, flat.SubjectID)
		s.Equal(event.Subject.Name, flat.SubjectName)
		s.Equal(bytes, flat.Payload)
	}
}

func (s *EventTestSuite) Test_PullRequestLabel_Create_Deflate() {
	label := &defs.PullRequestLabel{
		Name:              "bug",
		PullRequestNumber: 123,
		Branch:            "main",
		Timestamp:         time.Now(),
	}

	event := label.ToEvent(defs.RepoProviderGithub, s.subject, defs.EventActionCreated)
	event.SetSource("test/test")

	flat, err := event.Flatten()
	s.Require().NoError(err)

	var deflate defs.Event[defs.PullRequestLabel, defs.RepoProvider]

	err = defs.Deflate(flat, &deflate)
	if s.NoError(err) {
		s.Equal(event.Version, deflate.Version)
		s.Equal(event.ID, deflate.ID)
		s.Equal(event.Context.ParentID, deflate.Context.ParentID)
		s.Equal(event.Context.Provider, deflate.Context.Provider)
		s.Equal(event.Context.Scope, deflate.Context.Scope)
		s.Equal(event.Context.Action, deflate.Context.Action)
		s.Equal(event.Context.Source, deflate.Context.Source)
		s.Equal(event.Context.Timestamp.Unix(), deflate.Context.Timestamp.Unix())
		s.Equal(event.Subject, deflate.Subject)
	}
}

func (s *EventTestSuite) Test_PullRequestReview_Create_MarshalJSON() {
	review := &defs.PullRequestReview{
		ID:                1,
		State:             "approved",
		Author:            "testuser",
		PullRequestNumber: 123,
		Branch:            "main",
		Timestamp:         time.Now(),
	}

	event := review.ToEvent(defs.RepoProviderGithub, s.subject, defs.EventActionCreated)
	event.SetSource("test/test")

	expected := fmt.Sprintf(`{
  "version": "0.1.0",
  "id": "%s",
  "context": {
    "parent_id": "00000000-0000-0000-0000-000000000000",
    "provider": "github",
    "scope": "pull_request_review",
    "action": "created",
    "source": "test/test",
    "timestamp": "%s"
  },
  "subject": {
    "id": "%s",
    "name": "%s",
    "team_id": "%s"
  },
  "payload": {
    "id": 1,
    "pull_request_number": 123,
    "branch": "main",
    "state": "approved",
    "author": "testuser",
    "submitted_at": "%s"
  }
}`,
		event.ID,
		event.Context.Timestamp.Format(time.RFC3339Nano),
		s.subject.ID,
		s.subject.Name,
		s.subject.TeamID,
		review.Timestamp.Format(time.RFC3339Nano),
	)

	// Test Marshal to JSON
	marshal, err := json.MarshalIndent(event, "", "  ")
	if s.NoError(err) {
		s.T().Log(string(marshal))
		s.Equal(defs.EventVersionDefault, event.Version)
		s.Equal(uuid.Nil.String(), event.Context.ParentID.String())
		s.Equal(expected, string(marshal))
	}
}

func (s *EventTestSuite) Test_PullRequestReview_Create_UnmarshalJSON() {
	review := &defs.PullRequestReview{
		ID:                1,
		State:             "approved",
		Author:            "testuser",
		PullRequestNumber: 123,
		Timestamp:         time.Now(),
	}

	event := review.ToEvent(defs.RepoProviderGithub, s.subject, defs.EventActionCreated)
	event.SetSource("test/test")

	marshal, err := json.Marshal(event)
	s.Require().NoError(err)

	// Test Unmarshal from JSON
	var unmarshal defs.Event[defs.PullRequestReview, defs.RepoProvider]
	err = json.Unmarshal(marshal, &unmarshal)

	if s.NoError(err) {
		s.Equal(event.Version, unmarshal.Version)
		s.Equal(event.ID, unmarshal.ID)
		s.Equal(event.Context.ParentID, unmarshal.Context.ParentID)
		s.Equal(event.Context.Provider, unmarshal.Context.Provider)
		s.Equal(event.Context.Scope, unmarshal.Context.Scope)
		s.Equal(event.Context.Action, unmarshal.Context.Action)
		s.Equal(event.Context.Source, unmarshal.Context.Source)
		s.WithinDuration(event.Context.Timestamp, unmarshal.Context.Timestamp, time.Second)
		s.Equal(event.Subject.ID, unmarshal.Subject.ID)
		s.Equal(event.Subject.Name, unmarshal.Subject.Name)
		s.Equal(event.Subject.TeamID, unmarshal.Subject.TeamID)
		s.Equal(event.Payload.ID, unmarshal.Payload.ID)
		s.Equal(event.Payload.State, unmarshal.Payload.State)
		s.Equal(event.Payload.Author, unmarshal.Payload.Author)
		s.Equal(event.Payload.PullRequestNumber, unmarshal.Payload.PullRequestNumber)
		s.WithinDuration(event.Payload.Timestamp, unmarshal.Payload.Timestamp, time.Second)
	}
}

func (s *EventTestSuite) Test_PullRequestReview_Create_Flatten() {
	review := &defs.PullRequestReview{
		ID:                1,
		State:             "approved",
		Author:            "testuser",
		PullRequestNumber: 123,
		Timestamp:         time.Now(),
	}

	bytes, _ := json.Marshal(review)

	event := review.ToEvent(defs.RepoProviderGithub, s.subject, defs.EventActionCreated)
	event.SetSource("test/test")

	// Test flattening the event
	flat, err := event.Flatten()
	if s.NoError(err) {
		s.Equal(event.Version, flat.Version)
		s.Equal(event.ID, flat.ID)
		s.Equal(event.Context.Provider, flat.Provider)
		s.Equal(event.Context.Scope, flat.Scope)
		s.Equal(event.Context.Action, flat.Action)
		s.Equal(event.Context.Source, flat.Source)
		s.Equal(event.Context.Timestamp.Unix(), flat.CreatedAt.Unix())
		s.Equal(event.Context.Timestamp.Unix(), flat.UpdatedAt.Unix())
		s.Equal(event.Subject.ID, flat.SubjectID)
		s.Equal(event.Subject.Name, flat.SubjectName)
		s.Equal(bytes, flat.Payload)
	}
}

func (s *EventTestSuite) Test_PullRequestReview_Create_Deflate() {
	review := &defs.PullRequestReview{
		ID:                1,
		State:             "approved",
		Author:            "testuser",
		PullRequestNumber: 123,
		Timestamp:         time.Now(),
	}

	event := review.ToEvent(defs.RepoProviderGithub, s.subject, defs.EventActionCreated)
	event.SetSource("test/test")

	flat, err := event.Flatten()
	s.Require().NoError(err)

	var deflate defs.Event[defs.PullRequestReview, defs.RepoProvider]

	err = defs.Deflate(flat, &deflate)
	if s.NoError(err) {
		s.Equal(event.Version, deflate.Version)
		s.Equal(event.ID, deflate.ID)
		s.Equal(event.Context.ParentID, deflate.Context.ParentID)
		s.Equal(event.Context.Provider, deflate.Context.Provider)
		s.Equal(event.Context.Scope, deflate.Context.Scope)
		s.Equal(event.Context.Action, deflate.Context.Action)
		s.Equal(event.Context.Source, deflate.Context.Source)
		s.Equal(event.Context.Timestamp.Unix(), deflate.Context.Timestamp.Unix())
		s.Equal(event.Subject, deflate.Subject)
	}
}

func (s *EventTestSuite) Test_PullRequestReview_Create_SetParent() {
	review := &defs.PullRequestReview{
		ID:                1,
		State:             "approved",
		Author:            "testuser",
		PullRequestNumber: 123,
		Timestamp:         time.Now(),
	}

	event := review.ToEvent(defs.RepoProviderGithub, s.subject, defs.EventActionCreated)

	// Test setting the parent ID
	event.SetParent(s.parent)
	s.Equal(s.parent, event.Context.ParentID, "Parent ID should be set correctly")
}

func (s *EventTestSuite) Test_PullRequestLabel_Create_MarshalJSON() {
	label := &defs.PullRequestLabel{
		Name:              "bug",
		PullRequestNumber: 123,
		Branch:            "main",
		Timestamp:         time.Now(),
	}

	event := label.ToEvent(defs.RepoProviderGithub, s.subject, defs.EventActionCreated)
	event.SetSource("test/test")

	expected := fmt.Sprintf(`{
  "version": "0.1.0",
  "id": "%s",
  "context": {
    "parent_id": "00000000-0000-0000-0000-000000000000",
    "provider": "github",
    "scope": "pull_request_label",
    "action": "created",
    "source": "test/test",
    "timestamp": "%s"
  },
  "subject": {
    "id": "%s",
    "name": "%s",
    "team_id": "%s"
  },
  "payload": {
    "name": "bug",
    "pull_request_number": 123,
    "branch": "main",
    "timestamp": "%s"
  }
}`,
		event.ID,
		event.Context.Timestamp.Format(time.RFC3339Nano),
		s.subject.ID,
		s.subject.Name,
		s.subject.TeamID,
		label.Timestamp.Format(time.RFC3339Nano),
	)

	// Test Marshal to JSON
	marshal, err := json.MarshalIndent(event, "", "  ")
	if s.NoError(err) {
		s.T().Log(string(marshal))
		s.Equal(defs.EventVersionDefault, event.Version)
		s.Equal(uuid.Nil.String(), event.Context.ParentID.String())
		s.Equal(expected, string(marshal))
	}
}

func TestEventTestSuite(t *testing.T) {
	suite.Run(t, new(EventTestSuite))
}
