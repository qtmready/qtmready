// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package core

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
