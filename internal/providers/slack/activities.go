// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2022, 2024.
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


package slack

import (
	"context"

	"go.temporal.io/sdk/activity"

	"go.breu.io/quantm/internal/core/code"
	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/providers/github"
)

type (
	// Activities groups all the activities for the slack provider.
	Activities struct{}
)

var (
	coreacts   *code.Activities
	githubacts *github.Activities
)

func (a *Activities) SendStaleBranchMessage(ctx context.Context, payload *defs.MessageIOStaleBranchPayload) error {
	logger := activity.GetLogger(ctx)

	token, err := reveal(payload.MessageIOPayload.BotToken, payload.MessageIOPayload.WorkspaceID)
	if err != nil {
		logger.Error("Error in reveal", "Error", err)
		return err
	}

	client, err := instance.GetSlackClient(token)
	if err != nil {
		logger.Error("Error in GetSlackClient", "Error", err)
		return err
	}

	attachment := formatStaleBranchAttachment(payload)

	// call blockset to send the message to slack channel or sepecific workspace.
	if err := notify(client, payload.MessageIOPayload.ChannelID, attachment); err != nil {
		logger.Error("Failed to post message to channel", "Error", err)
		return err
	}

	logger.Info("Slack notification sent successfully")

	return nil
}

func (a *Activities) SendNumberOfLinesExceedMessage(ctx context.Context, payload *defs.MessageIOLineExeededPayload) error {
	logger := activity.GetLogger(ctx)

	token, err := reveal(payload.MessageIOPayload.BotToken, payload.MessageIOPayload.WorkspaceID)
	if err != nil {
		logger.Error("Error in reveal", "Error", err)
		return err
	}

	client, err := instance.GetSlackClient(token)
	if err != nil {
		logger.Error("Error in GetSlackClient", "Error", err)
		return err
	}

	attachment := formatLineThresholdExceededAttachment(payload)

	// Call function to send the message to Slack channel or specific workspace.
	if err := notify(client, payload.MessageIOPayload.ChannelID, attachment); err != nil {
		logger.Error("Failed to post message to channel", "Error", err)
		return err
	}

	logger.Info("Slack notification sent successfully")

	return nil
}

func (a *Activities) SendMergeConflictsMessage(ctx context.Context, payload *defs.MergeConflictMessage) error {
	logger := activity.GetLogger(ctx)

	token, err := reveal(payload.MessageIOPayload.BotToken, payload.MessageIOPayload.WorkspaceID)
	if err != nil {
		logger.Error("Error in reveal", "Error", err)
		return err
	}

	client, err := instance.GetSlackClient(token)
	if err != nil {
		logger.Error("Error in GetSlackClient", "Error", err)
		return err
	}

	attachment := formatMergeConflictAttachment(payload)

	// call blockset to send the message to slack channel or sepecific workspace.
	if err := notify(client, payload.MessageIOPayload.ChannelID, attachment); err != nil {
		logger.Error("Failed to post message to channel", "Error", err)
		return err
	}

	logger.Info("Slack notification sent successfully")

	return nil
}

// TODO - move the uint functions.
func (a *Activities) NotifyMergeConflict(ctx context.Context, event *defs.Event[defs.MergeConflict, defs.RepoProvider]) error {
	logger := activity.GetLogger(ctx)

	token, channelID, err := derive(ctx, event)
	if err != nil {
		logger.Error("Error in getTokenAndChannelID", "Error", err)
		return err
	}

	client, err := instance.GetSlackClient(token)
	if err != nil {
		logger.Error("Error in GetSlackClient", "Error", err)
		return err
	}

	attachment := compose_merge_conflict(event)

	// call blockset to send the message to slack channel or sepecific workspace.
	if err := notify(client, channelID, attachment); err != nil {
		logger.Error("Failed to post message to channel", "Error", err)
		return err
	}

	logger.Info("Slack notification sent successfully")

	return nil
}
