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
	"encoding/base64"
	"log/slog"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core/code"
	"go.breu.io/quantm/internal/core/defs"
)

// reveal decodes a base64-encoded encrypted token and decrypts it using a generated key.
func reveal(botToken, workspaceID string) (string, error) {
	// Decode the base64-encoded encrypted token.
	decoded, err := base64.StdEncoding.DecodeString(botToken)
	if err != nil {
		slog.Error("Failed to decode the token", slog.Any("e", err))
		return "", err
	}

	// Generate the same key used for encryption.
	key := generate(workspaceID)

	// Decrypt the token.
	decrypted, err := decrypt(decoded, key)
	if err != nil {
		slog.Error("Failed to decrypt the token", slog.Any("e", err))
		return "", err
	}

	return string(decrypted), nil
}

// NOTE - may be move the core or code
// derive determines whether to retrieve token/channel data from user or repo.
func derive(ctx context.Context, event *defs.Event[defs.MergeConflict, defs.RepoProvider]) (string, string, error) {
	tuser, err := githubacts.GetTeamUserByLoginID(ctx, event.Subject.UserID.String())
	if err != nil {
		return "", "", err
	}

	if tuser != nil {
		return user_data(tuser)
	}

	return repo_data(ctx, event.Subject.ID.String())
}

// user_data extracts token and channel ID from user-specific message provider data.
func user_data(tuser *auth.TeamUser) (string, string, error) {
	token, err := reveal(tuser.MessageProviderUserInfo.Slack.BotToken, tuser.MessageProviderUserInfo.Slack.ProviderTeamID)
	if err != nil {
		return "", "", err
	}

	return token, tuser.MessageProviderUserInfo.Slack.ProviderUserID, nil
}

// repo_data extracts token and channel ID from repo-specific message provider data.
func repo_data(ctx context.Context, repoID string) (string, string, error) {
	repo, err := code.RepoIO().GetByID(ctx, repoID)
	if err != nil {
		return "", "", err
	}

	token, err := reveal(repo.MessageProviderData.Slack.BotToken, repo.MessageProviderData.Slack.WorkspaceID)
	if err != nil {
		return "", "", err
	}

	return token, repo.MessageProviderData.Slack.ChannelID, nil
}
