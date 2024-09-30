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

package comm

import (
	"fmt"
	"time"

	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/db"
)

// NewMergeConflictEvent creates a new defs.Event instance for a merge conflict.
func NewMergeConflictEvent(
	event *defs.Event[defs.Push, defs.RepoProvider], head, base string, base_commit *defs.Commit,
) *defs.Event[defs.MergeConflict, defs.RepoProvider] {
	id, _ := db.NewUUID()
	now := time.Now()

	// creating payload
	conflict := defs.MergeConflict{
		HeadBranch: head,
		HeadCommit: *event.Payload.Commits.Latest(),
		BaseBranch: base,
		BaseCommit: *base_commit,
		Files:      make([]string, 0),
		Timestamp:  now,
	}

	// creating new event
	reply := &defs.Event[defs.MergeConflict, defs.RepoProvider]{
		Version: event.Version,
		ID:      id,
		Context: event.Context,
		Subject: event.Subject,
		Payload: conflict,
	}

	// updating event
	reply.SetParent(event.ID)
	reply.SetScopeMergeConflict()
	reply.SetActionCreated()
	reply.SetTimestamp(now)

	return reply
}

// NewLineExceedEvent creates a new defs.Event instance for a line exceed.
func NewLineExceedEvent(
	event *defs.Event[defs.Push, defs.RepoProvider], head string, lc *defs.LineChanges,
) *defs.Event[defs.LinesExceed, defs.RepoProvider] {
	id, _ := db.NewUUID()
	now := time.Now()

	// creating payload
	exceed := defs.LinesExceed{
		Branch:    head,
		Commit:    *event.Payload.Commits.Latest(),
		LineStats: *lc,
		Timestamp: now,
	}

	// creating new event
	reply := &defs.Event[defs.LinesExceed, defs.RepoProvider]{
		Version: event.Version,
		ID:      id,
		Context: event.Context,
		Subject: event.Subject,
		Payload: exceed,
	}

	// updating event
	reply.SetParent(event.ID)
	reply.SetScopeLineExceed()
	reply.SetActionCreated()
	reply.SetTimestamp(now)

	return reply
}

// NewStaleBranchMessage creates a new MessageIOStaleBranchPayload instance.
//
// It takes RepoIOProviderInfo, Repo information, and a branch name. The
// function constructs URLs for the commit and repository, and sets the
// MessageIOPayload for the channel. This function is only used for channel
// messages.
// TODO - handle using event.
func NewStaleBranchMessage(data *defs.RepoIOProviderInfo, repo *defs.Repo, branch string) *defs.MessageIOStaleBranchPayload {
	return &defs.MessageIOStaleBranchPayload{
		CommitUrl: fmt.Sprintf("https://github.com/%s/%s/tree/%s",
			data.RepoOwner, data.RepoName, branch),
		RepoUrl: fmt.Sprintf("https://github.com/%s/%s", data.RepoOwner, data.RepoName),
		MessageIOPayload: &defs.MessageIOPayload{
			WorkspaceID: repo.MessageProviderData.Slack.WorkspaceID,
			ChannelID:   repo.MessageProviderData.Slack.ChannelID,
			BotToken:    repo.MessageProviderData.Slack.BotToken,
			RepoName:    repo.Name,
			BranchName:  branch,
		},
	}
}
