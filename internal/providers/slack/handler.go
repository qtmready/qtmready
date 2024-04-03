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
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/db"
)

type (
	ServerHandler struct{ *auth.SecurityHandler }

	// TODO: move to openapi.
	SlackInfo struct {
		ChannelID string
		Workspace string
		Channels  []slack.Channel
	}
)

// NewServerHandler creates a new ServerHandler.
func NewServerHandler(middleware echo.MiddlewareFunc) *ServerHandler {
	return &ServerHandler{
		SecurityHandler: &auth.SecurityHandler{Middleware: middleware},
	}
}

func (e *ServerHandler) SlackLogin(ctx echo.Context) error {
	url := Instance().OauthConfig.AuthCodeURL("state")
	return ctx.Redirect(http.StatusFound, url)
}

// this for runs one time for one time.
func (e *ServerHandler) SlackOauth(ctx echo.Context) error {
	code := ctx.QueryParam("code")
	if code == "" {
		return errors.New("empty code")
	}

	// Exchange the authorization code for an access token
	token, err := Instance().OauthConfig.Exchange(ctx.Request().Context(), code)
	if err != nil {
		return err
	}

	// Create a Slack client using the obtained access token.
	client := slack.New(token.AccessToken)

	// Use auth.test method to get information about the authenticated user (bot)
	auth, err := client.AuthTestContext(ctx.Request().Context())
	if err != nil {
		return err
	}

	// NOTE: to send a message to workspace's channel need to requied the bot token returned in auth response.

	// Use conversations.list method to get a list of channels in the workspace
	channels, _, err := client.GetConversations(&slack.GetConversationsParameters{
		Types: []string{"public_channel"},
	})
	if err != nil {
		return err
	}

	channelID := findChannelIdForBot(client, channels, auth)

	// TODO: save this channel to our database

	// Construct response with workspace info, and channel details
	response := SlackInfo{
		ChannelID: channelID,
		Workspace: auth.Team,
		Channels:  channels,
	}

	return ctx.JSON(http.StatusOK, response)
}

// TODO: may return a list of channelIDs.
func findChannelIdForBot(client *slack.Client, channels []slack.Channel, auth *slack.AuthTestResponse) string {
	var channelID string

	for _, channel := range channels {
		input := &slack.GetConversationInfoInput{
			ChannelID:         channel.ID,
			IncludeLocale:     false,
			IncludeNumMembers: false,
		}

		// Use the GetConversationInfo method to get information about the channel
		channelInfo, _ := client.GetConversationInfo(input)

		// compare the auth.TeamID with channelInfo.SharedTeamIDs array
		for _, sharedTeamID := range channelInfo.SharedTeamIDs {
			if auth.TeamID == sharedTeamID {
				channelID = channel.ID
			}
		}
	}

	return channelID
}

func (e *ServerHandler) SlackIntegration(ctx echo.Context) error {
	// obtain query parameters
	workspace_name := ctx.QueryParam("workspace_name")

	// get from db
	slackIntegration := &SlackIntegration{}
	params := db.QueryParams{"workspace_name": workspace_name}

	if err := db.Get(slackIntegration, params); err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Error obtaining from database")
	}

	return ctx.JSON(http.StatusOK, slackIntegration)
}
