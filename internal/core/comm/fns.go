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

	"go.breu.io/quantm/internal/core/defs"
)

// NewMergeConflictMessage creates a new MessageIOMergeConflictPayload instance.
//
// It takes a RepoIOSignalPushPayload, Repo information, branch name, and a flag
// indicating whether the message is for a user or a channel. The function
// constructs URLs for the repository and commit, and sets the appropriate
// MessageIOPayload based on the for_user flag.
//
// FIXME: this is generic to github. If we are using generic, should we create the url's depending upon the provider?
func NewMergeConflictMessage(
	payload *defs.RepoIOSignalPushPayload,
	repo *defs.Repo,
	branch string,
	for_user bool,
) *defs.MessageIOMergeConflictPayload {
	msg := &defs.MessageIOMergeConflictPayload{
		RepoUrl:   fmt.Sprintf("https://github.com/%s/%s", payload.RepoOwner, payload.RepoName),
		SHA:       payload.After,
		CommitUrl: fmt.Sprintf("https://github.com/%s/%s/commits/%s", payload.RepoOwner, payload.RepoName, payload.After),
	}

	// set the payload for user message provider
	if for_user {
		msg.MessageIOPayload = &defs.MessageIOPayload{
			WorkspaceID: payload.User.MessageProviderUserInfo.Slack.ProviderTeamID,
			ChannelID:   payload.User.MessageProviderUserInfo.Slack.ProviderUserID,
			BotToken:    payload.User.MessageProviderUserInfo.Slack.BotToken,
			RepoName:    repo.Name,
			BranchName:  branch,
			IsChannel:   false,
		}
	} else {
		// set the payload for channel message provider
		msg.MessageIOPayload = &defs.MessageIOPayload{
			WorkspaceID: repo.MessageProviderData.Slack.WorkspaceID,
			ChannelID:   repo.MessageProviderData.Slack.ChannelID,
			BotToken:    repo.MessageProviderData.Slack.BotToken,
			Author:      payload.Author,
			AuthorUrl:   fmt.Sprintf("https://github.com/%s", payload.Author),
			RepoName:    repo.Name,
			BranchName:  branch,
			IsChannel:   true,
		}
	}

	return msg
}

// NewNumberOfLinesExceedMessage creates a new MessageIOLineExeededPayload instance.
//
// It takes a RepoIOSignalPushPayload, Repo information, branch name, changes,
// and a flag indicating whether the message is for a user or a channel. The
// function sets the threshold and detected changes, and constructs the
// appropriate MessageIOPayload based on the for_user flag.
func NewNumberOfLinesExceedMessage(
	payload *defs.RepoIOSignalPushPayload,
	repo *defs.Repo,
	branch string,
	changes *defs.RepoIOChanges,
	for_user bool,
) *defs.MessageIOLineExeededPayload {
	msg := &defs.MessageIOLineExeededPayload{
		Threshold:     repo.Threshold,
		DetectChanges: changes,
	}

	// set the payload for user message provider
	if for_user {
		msg.MessageIOPayload = &defs.MessageIOPayload{
			WorkspaceID: payload.User.MessageProviderUserInfo.Slack.ProviderTeamID,
			ChannelID:   payload.User.MessageProviderUserInfo.Slack.ProviderUserID,
			BotToken:    payload.User.MessageProviderUserInfo.Slack.BotToken,
			RepoName:    repo.Name,
			BranchName:  branch,
			IsChannel:   false,
		}
	} else {
		// set the payload for channel message provider
		msg.MessageIOPayload = &defs.MessageIOPayload{
			WorkspaceID: repo.MessageProviderData.Slack.WorkspaceID,
			ChannelID:   repo.MessageProviderData.Slack.ChannelID,
			BotToken:    repo.MessageProviderData.Slack.BotToken,
			Author:      payload.Author,
			AuthorUrl:   fmt.Sprintf("https://github.com/%s", payload.Author),
			RepoName:    repo.Name,
			BranchName:  branch,
			IsChannel:   true,
		}
	}

	return msg
}

// NewStaleBranchMessage creates a new MessageIOStaleBranchPayload instance.
//
// It takes RepoIOProviderInfo, Repo information, and a branch name. The
// function constructs URLs for the commit and repository, and sets the
// MessageIOPayload for the channel. This function is only used for channel
// messages.
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
