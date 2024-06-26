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
	"net/http"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/shared"
)

type (
	ServerHandler struct{ *auth.SecurityHandler }
)

// NewServerHandler creates a new ServerHandler.
func NewServerHandler(middleware echo.MiddlewareFunc) *ServerHandler {
	return &ServerHandler{
		SecurityHandler: &auth.SecurityHandler{Middleware: middleware},
	}
}

func (e *ServerHandler) SlackOauth(ctx echo.Context, params SlackOauthParams) error {
	var c HTTPClient

	code := ctx.QueryParam("code")
	if code == "" {
		return shared.NewAPIError(http.StatusNotFound, ErrCodeEmpty)
	}

	response, err := slack.GetOAuthV2Response(&c, ClientID(), ClientSecret(), code, ClientRedirectURL())
	if err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	if response.AuthedUser.AccessToken != "" {
		// linked message provider user (slack) to quantm user
		_, err := _user(ctx, response)
		if err != nil {
			return shared.NewAPIError(http.StatusBadRequest, err)
		}

		return ctx.JSON(http.StatusOK, nil)
	}

	// connect slack bot user with channel info
	resp, err := _bot(response)
	if err != nil {
		return shared.NewAPIError(http.StatusBadRequest, err)
	}

	return ctx.JSON(http.StatusOK, resp)
}

func _user(ctx echo.Context, response *slack.OAuthV2Response) (*auth.MessageProviderUserInfo, error) {
	userID, _ := gocql.ParseUUID(ctx.Get("user_id").(string))

	teamuser := &auth.TeamUser{}
	query := db.QueryParams{"user_id": userID.String()}

	if err := db.Get(teamuser, query); err != nil {
		return nil, err
	}

	client, _ := instance.GetSlackClient(response.AuthedUser.AccessToken)

	identity, err := client.GetUserIdentity()
	if err != nil {
		shared.Logger().Error("SlackOauth/identity", "error", err.Error())
		return nil, err
	}

	// Generate a key for AES-256.
	key := generateKey(response.Team.ID)

	// Encrypt the access token.
	encryptedToken, err := encrypt([]byte(response.AuthedUser.AccessToken), key)
	if err != nil {
		return nil, err
	}

	userinfo := auth.MessageProviderUserInfo{
		Slack: &auth.MessageProviderSlackUserInfo{
			UserToken:      base64.StdEncoding.EncodeToString(encryptedToken),
			ProviderUserID: identity.User.ID,
			ProviderTeamID: identity.Team.ID,
		},
	}

	teamuser.IsMessageProviderLinked = true
	teamuser.MessageProvider = auth.MessageProviderSlack
	teamuser.MessageProviderUserInfo = userinfo

	if err := db.Save(teamuser); err != nil {
		return nil, err
	}

	return &userinfo, nil
}

func _bot(response *slack.OAuthV2Response) (*core.MessageProviderData, error) {
	// Generate a key for AES-256.
	key := generateKey(response.Team.ID)

	// Encrypt the access token.
	encryptedToken, err := encrypt([]byte(response.AccessToken), key)
	if err != nil {
		return nil, err
	}

	resp := &core.MessageProviderData{
		Slack: &core.MessageProviderSlackData{
			ChannelID:     response.IncomingWebhook.ChannelID,
			ChannelName:   response.IncomingWebhook.Channel,
			WorkspaceName: response.Team.Name,
			WorkspaceID:   response.Team.ID,
			BotToken:      base64.StdEncoding.EncodeToString(encryptedToken), // Store the base64-encoded encrypted token
		},
	}

	return resp, nil
}
