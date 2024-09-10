// Copyright © 2023, Breu, Inc. <info@breu.io>
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
	"go.breu.io/quantm/internal/shared"
)

type (
	// NOTE - this base struct need for any type of message. getting from core repo.
	MessageIOPayload struct {
		WorkspaceID string `json:"workspace_id"`
		ChannelID   string `json:"channel_id"`
		BotToken    string `json:"bot_token"`
		RepoName    string `json:"repo_name"`
		BranchName  string `json:"branch_name"`
		Author      string `json:"author"`
		AuthorUrl   string `json:"author_url"`
		IsChannel   bool   `json:"is_channel"`
	}

	// TODO: need to refine.
	MessageIOLineExeededPayload struct {
		MessageIOPayload *MessageIOPayload `json:"message_io_payload"`
		Threshold        shared.Int64      `json:"threshold"`
		DetectChanges    *RepoIOChanges    `json:"detect_changes"`
	}

	// TODO: need to refine.
	MessageIOMergeConflictPayload struct {
		MessageIOPayload *MessageIOPayload `json:"message_io_payload"`
		CommitUrl        string            `json:"commit_url"`
		RepoUrl          string            `json:"repo_url"`
		SHA              string            `json:"sha"`
	}

	// TODO: need to refine.
	MessageIOStaleBranchPayload struct {
		MessageIOPayload *MessageIOPayload `json:"message_io_payload"`
		CommitUrl        string            `json:"commit_url"`
		RepoUrl          string            `json:"repo_url"`
	}
)
