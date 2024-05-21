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

	"go.temporal.io/sdk/activity"

	"go.breu.io/quantm/internal/core"
)

type (
	// Activities groups all the activities for the slack provider.
	Activities struct{}
)

func (a *Activities) SendStaleBranchMessage(ctx context.Context, teamID string, stale *core.LatestCommit) error {
	logger := activity.GetLogger(ctx)

	client, channelID, err := GetSlackClientAndChannelID(teamID)
	if err != nil {
		logger.Error("Error in GetSlackClientAndChannelID", "Error", err)
		return err
	}

	attachment := formatStaleBranchAttachment(stale)

	// call blockset to send the message to slack channel or sepecific workspace.
	if err := notify(client, channelID, attachment); err != nil {
		logger.Error("Failed to post message to channel", "Error", err)
		return err
	}

	logger.Info("Slack notification sent successfully")

	return nil
}

func (a *Activities) SendNumberOfLinesExceedMessage(ctx context.Context, payload *core.MessageIOLineExeededPayload) error {
	logger := activity.GetLogger(ctx)

	token, err := decodeAndDecryptToken(payload.BotToken, payload.WorkspaceID)
	if err != nil {
		logger.Error("Error in decodeAndDecryptToken", "Error", err)
		return err
	}

	client, err := instance.GetSlackClient(token)
	if err != nil {
		logger.Error("Error in GetSlackClient", "Error", err)
		return err
	}

	attachment := formatLineThresholdExceededAttachment(payload)

	// Call function to send the message to Slack channel or specific workspace.
	if err := notify(client, payload.ChannelID, attachment); err != nil {
		logger.Error("Failed to post message to channel", "Error", err)
		return err
	}

	logger.Info("Slack notification sent successfully")

	return nil
}

func (a *Activities) SendMergeConflictsMessage(ctx context.Context, teamID string, merge *core.LatestCommit) error {
	logger := activity.GetLogger(ctx)

	client, channelID, err := GetSlackClientAndChannelID(teamID)
	if err != nil {
		logger.Error("Error in GetSlackClientAndChannelID", "Error", err)
		return err
	}

	attachment := formatMergeConflictAttachment(merge)

	// call blockset to send the message to slack channel or sepecific workspace.
	if err := notify(client, channelID, attachment); err != nil {
		logger.Error("Failed to post message to channel", "Error", err)
		return err
	}

	logger.Info("Slack notification sent successfully")

	return nil
}
