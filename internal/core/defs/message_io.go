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

package defs

import (
	"go.breu.io/quantm/internal/db"
)

type (
	// MessageIOPayload represents the base payload for all message-related events.
	//
	// This struct contains common information shared across different message types, such as workspace ID, channel ID,
	// bot token, repository name, branch name, author, author URL, and whether the event is triggered by a channel.
	//
	// NOTE: This struct is a shared resource, sourced from the core repository.
	MessageIOPayload struct {
		WorkspaceID string `json:"workspace_id"` // ID of the workspace.
		ChannelID   string `json:"channel_id"`   // ID of the channel.
		BotToken    string `json:"bot_token"`    // Bot's authentication token.
		RepoName    string `json:"repo_name"`    // Name of the repository.
		BranchName  string `json:"branch_name"`  // Name of the branch.
		Author      string `json:"author"`       // Author of the event.
		AuthorURL   string `json:"author_url"`   // URL of the author's profile.
		IsChannel   bool   `json:"is_channel"`   // Indicates whether the event originates from a channel.
	}

	// MessageIOLineExeededPayload contains information about a line limit exceeding event.
	//
	// TODO: Revise this structure for improved clarity and functionality.
	MessageIOLineExeededPayload struct {
		MessageIOPayload *MessageIOPayload `json:"message_io_payload"` // Base payload for message events.
		Threshold        db.Int64          `json:"threshold"`          // Line limit threshold.
		DetectChanges    *RepoIOChanges    `json:"detect_changes"`     // Details of the detected changes.
	}

	// MergeConflictMessage encapsulates information related to a merge conflict event.
	//
	// TODO: Revise this structure for improved clarity and functionality.
	MergeConflictMessage struct {
		MessageIOPayload *MessageIOPayload `json:"message_io_payload"` // Base payload for message events.
		CommitUrl        string            `json:"commit_url"`         // URL of the commit related to the conflict.
		RepoUrl          string            `json:"repo_url"`           // URL of the repository.
		SHA              string            `json:"sha"`                // SHA hash of the commit.
	}

	// MessageIOStaleBranchPayload provides information about a stale branch event.
	//
	// TODO: Revise this structure for improved clarity and functionality.
	MessageIOStaleBranchPayload struct {
		MessageIOPayload *MessageIOPayload `json:"message_io_payload"` // Base payload for message events.
		CommitUrl        string            `json:"commit_url"`         // URL of the commit related to the stale branch.
		RepoUrl          string            `json:"repo_url"`           // URL of the repository.
	}
)
