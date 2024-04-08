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
	"context"
	"log/slog"

	"github.com/slack-go/slack"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

type (
	// Activities groups all the activities for the slack provider.
	Activities struct{}
)

// get the slack record by team id
// get the access token for workspace and channel ID
// create a slack instance
// post message the specific channel
// call the notify function to post the message to slack workspace channel.
func (a *Activities) SendChannelMessage(ctx context.Context, teamID, message string) error {
	// Get the slack info from database
	s := &Slack{}
	params := db.QueryParams{"team_id": teamID}

	if err := db.Get(s, params); err != nil {
		slog.Info("Failed to get the slack record", slog.Any("e", err))
		return err
	}

	// Create a Slack client using the obtained access token.
	client := slack.New(s.WorkspaceBotToken)

	// call blockset to send the message to slack channel or sepecific workspace.
	if err := notify(client, s.ChannelID, message); err != nil {
		slog.Info("Failed to post message to channel", slog.Any("e", err))
		return err
	}

	shared.Logger().Info("Early-Detection", "slack notify", message)

	return nil
}
