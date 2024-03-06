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
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack"

	"go.breu.io/quantm/internal/auth"
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

func (e *ServerHandler) Login(ctx echo.Context) error {
	url := Instance().OauthConfig.AuthCodeURL("")
	return ctx.Redirect(http.StatusFound, url)
}

func (e *ServerHandler) SlackOauth(ctx echo.Context) error {
	code := ctx.QueryParam("code")
	// token, err := Instance().OauthConfig.Exchange(ctx.Request().Context(), code)

	token, _, err := slack.GetOAuthToken(http.DefaultClient, ClientID(), ClientSecret(), code, ClientRedirectURL())
	if err != nil {
		return err
	}

	btoken, _, bot, err := slack.GetBotOAuthToken(http.DefaultClient, ClientID(), ClientSecret(), code, ClientRedirectURL())
	if err != nil {
		return err
	}

	oauth, err := slack.GetOAuthV2Response(http.DefaultClient, ClientID(), ClientSecret(), code, ClientRedirectURL())
	if err != nil {
		return err
	}

	log.Println("token", token)
	log.Println("btoken", btoken)
	log.Println("bot", bot)

	return ctx.JSON(http.StatusOK, map[string]any{
		"oauth":  oauth,
		"token":  token,
		"btoken": btoken,
		"bot":    bot,
	})
}
