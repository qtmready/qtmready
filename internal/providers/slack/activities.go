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
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/slack-go/slack"
	"go.temporal.io/sdk/activity"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core/code"
	"go.breu.io/quantm/internal/core/defs"
	"go.breu.io/quantm/internal/db"
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
	codeacts := &code.Activities{}

	repo, err := codeacts.GetCoreRepoByID(ctx, event.Subject.ID.String()) // FIXME: we should directly get the token and channel id.
	if err != nil {
		return err
	}

	target := repo.MessageProviderData.Slack.WorkspaceID

	token, err := reveal(repo.MessageProviderData.Slack.BotToken, target)
	if err != nil {
		return err
	}

	fields := a.conflict_fields(event)

	if event.Subject.UserID.String() != db.NullUUID {
		// user, err := authacts.GetTeamUser(ctx, event.Subject.UserID.String())
		user, err := auth.TeamUserIO().Get(ctx, event.Subject.UserID.String(), event.Subject.TeamID.String())
		if user == nil || err != nil { // This should never error out.
			return err
		}
	}

	// by default we send the message to the channel, but if the user is a member of the team, we send the message to the user.
	// if we are sending on the channel, we introduce the owner of the branch to the channel.
	if event.Subject.UserID.String() != db.NullUUID {
		token, target, err = a.user_tokens(ctx, event)
		if err != nil {
			logger.Error("Error in getTokenAndChannelID", "Error", err)
			return err
		}
	} else {
		owner := event.Payload.BaseCommit.Author
		url := fmt.Sprintf("https://github.com/%s", owner)
		fields = append(fields, slack.AttachmentField{
			Title: "Branch Owner",
			Value: fmt.Sprintf("[%s](%s)", owner, url),
			Short: false,
		})
	}

	attachment := slack.Attachment{
		Color: "warning",
		Pretext: `We've detected a merge conflict in your feature branch, [Branch Name]. 
		This means there are changes in your branch that clash with recent updates on the main branch (trunk).`,
		Fallback:   "Merge Conflict Detected",
		MarkdownIn: []string{"fields"},
		Footer:     footer,
		Fields:     fields,
		Ts:         json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}

	client, err := instance.GetSlackClient(token)
	if err != nil {
		logger.Error("Error in GetSlackClient", "Error", err)
		return err
	}

	return notify(client, target, attachment)
}

func (a *Activities) conflict_fields(event *defs.Event[defs.MergeConflict, defs.RepoProvider]) []slack.AttachmentField {
	fields := []slack.AttachmentField{
		{
			Title: "*Branch*",
			Value: fmt.Sprintf("%s", event.Payload.BaseBranch),
			Short: true,
		}, {
			Title: "Current HEAD",
			Value: fmt.Sprintf("%s", event.Payload.BaseBranch),
			Short: true,
		}, {
			Title: "Conflict HEAD",
			Value: fmt.Sprintf("%s", event.Payload.HeadCommit.SHA),
			Short: true,
		},
		{
			Title: "Affected Files",
			Value: fmt.Sprintf("%s", event.Payload.Files), // TODO:
			Short: true,
		},
	}

	return fields
}

func (a *Activities) user_tokens(ctx context.Context, event *defs.Event[defs.MergeConflict, defs.RepoProvider]) (string, string, error) {
	tuser, err := githubacts.GetTeamUserByLoginID(ctx, event.Subject.UserID.String())

	if err != nil || tuser == nil {
		return "", "", nil // We should never arrive here.
	}

	token, err := reveal(tuser.MessageProviderUserInfo.Slack.BotToken, tuser.MessageProviderUserInfo.Slack.ProviderTeamID)
	if err != nil {
		return "", "", err
	}

	return token, tuser.MessageProviderUserInfo.Slack.ProviderUserID, nil
}

func (a *Activities) team_tokens(ctx context.Context, event *defs.Event[defs.MergeConflict, defs.RepoProvider]) (string, string, error) {
	return "", "", nil
}
