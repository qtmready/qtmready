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
	"encoding/base64"

	"github.com/slack-go/slack"
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

func (a *Activities) SendNumberOfLinesExceedMessage(
	ctx context.Context, teamID, repoName, branchName string, threshold int, branchChanges *core.BranchChanges,
) error {
	logger := activity.GetLogger(ctx)

	client, channelID, err := GetSlackClientAndChannelID(teamID)
	if err != nil {
		logger.Error("Error in GetSlackClientAndChannelID", "Error", err)
		return err
	}

	attachment := formatLineThresholdExceededAttachment(repoName, branchName, threshold, branchChanges)

	// Call function to send the message to Slack channel or specific workspace.
	if err := notify(client, channelID, attachment); err != nil {
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

func (a *Activities) CompleteOauthResponse(ctx context.Context, code string) (*core.MessageProviderData, error) {
	logger := activity.GetLogger(ctx)

	var c HTTPClient

	response, err := slack.GetOAuthV2Response(&c, ClientID(), ClientSecret(), code, ClientRedirectURL())
	if err != nil {
		logger.Error("Failed get response from slack oauth", "Error", err)
		return nil, err
	}

	// Generate a key for AES-256.
	key := generateKey(response.Team.ID)

	// Encrypt the access token.
	encryptedToken, err := encrypt([]byte(response.AccessToken), key)
	if err != nil {
		logger.Error("Failed to encrypt bot token", "Error", err)
		return nil, err
	}

	mpsd := core.MessageProviderData{
		Slack: &core.MessageProviderSlackData{
			BotToken:      base64.StdEncoding.EncodeToString(encryptedToken), // Store the base64-encoded encrypted token
			ChannelID:     response.IncomingWebhook.ChannelID,
			ChannelName:   response.IncomingWebhook.Channel,
			WorkspaceName: response.Team.Name,
			WorkspaceID:   response.Team.ID,
		},
	}

	return &mpsd, nil
}
