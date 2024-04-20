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

package slack

import (
	"encoding/base64"
	"log/slog"

	"github.com/slack-go/slack"

	"go.breu.io/quantm/internal/db"
)

func GetSlackClientAndChannelID(teamID string) (*slack.Client, string, error) {
	// Get the slack info from the database
	s := &Slack{}
	params := db.QueryParams{"team_id": teamID}

	if err := db.Get(s, params); err != nil {
		slog.Info("Failed to get the slack record", slog.Any("e", err))
		return nil, "", err
	}

	// Decode the base64-encoded encrypted token.
	decode, err := base64.StdEncoding.DecodeString(s.WorkspaceBotToken)
	if err != nil {
		slog.Info("Failed to decode the token", slog.Any("e", err))
		return nil, "", err
	}

	// Generate the same key used for encryption.
	key := generateKey(s.WorkspaceID)

	// Decrypt the token.
	decryptedToken, err := decrypt(decode, key)
	if err != nil {
		slog.Info("Failed to decrypt the token", slog.Any("e", err))
		return nil, "", err
	}

	// Create a Slack client using the decrypted access token.
	client := slack.New(string(decryptedToken))

	return client, s.ChannelID, nil
}
