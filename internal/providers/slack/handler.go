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
	"net/http"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack"

	"go.breu.io/quantm/internal/auth"
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

// this for runs one time for one time.
func (e *ServerHandler) SlackOauth(ctx echo.Context) error {
	var c HTTPClient

	code := ctx.QueryParam("code")
	if code == "" {
		return shared.NewAPIError(http.StatusNotFound, ErrCodeEmpty)
	}

	teamID, _ := gocql.ParseUUID(ctx.Get("team_id").(string))

	response, err := slack.GetOAuthV2Response(&c, ClientID(), ClientSecret(), code, ClientRedirectURL())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	slack := &Slack{
		ChannelID:         response.IncomingWebhook.ChannelID,
		ChannelName:       response.IncomingWebhook.Channel,
		WorkspaceName:     response.Team.Name,
		WorkspaceID:       response.Team.ID,
		WorkspaceBotToken: response.AccessToken,
		TeamID:            teamID,
	}

	if err := db.Save(slack); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, slack)
}
