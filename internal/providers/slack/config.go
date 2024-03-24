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
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"golang.org/x/oauth2"

	"go.breu.io/quantm/internal/shared"
)

type (
	config struct {
		BotToken     string `env:"SLACK_BOT_TOKEN"`
		UserToken    string `env:"SLACK_USER_TOKEN"`
		AppToken     string `env:"SLACK_APP_TOKEN"`
		ClientID     string `env:"SLACK_CLIENT_ID"`
		ClientSecret string `env:"SLACK_CLIENT_SECRET"`
		RedirectURL  string `env:"SLACK_REDIRECT_URL"`
	}

	integration struct {
		Config      *config
		Client      *slack.Client
		Socket      *socketmode.Client
		OauthConfig *oauth2.Config
	}
)

var (
	once     sync.Once
	instance *integration
)

func Instance() *integration {
	once.Do(func() {
		c := &config{}

		if err := cleanenv.ReadEnv(c); err != nil {
			panic("Failed to load slack configuration from environment variables: " + err.Error())
		}

		instance = connect(c)
	})

	return instance
}

func connect(c *config) *integration {
	lgr := &logger{shared.Logger().WithGroup("slack")}
	client := slack.New(
		c.BotToken,
		slack.OptionDebug(shared.Service().GetDebug()),
		slack.OptionAppLevelToken(c.AppToken),
		slack.OptionLog(lgr),
	)

	socket := socketmode.New(
		client,
		socketmode.OptionDebug(shared.Service().GetDebug()),
		socketmode.OptionLog(lgr),
	)

	oauthConfig := &oauth2.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
		Scopes:       []string{"channels:read", "chat:write"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://slack.com/oauth/v2/authorize",
			TokenURL: "https://slack.com/api/oauth.v2.access",
		},
	}

	return &integration{
		Config:      c,
		Client:      client,
		Socket:      socket,
		OauthConfig: oauthConfig,
	}
}

func SlackClient() *slack.Client {
	return Instance().Client
}

func Socket() *socketmode.Client {
	return Instance().Socket
}

func BotToken() string {
	return Instance().Config.BotToken
}

func AppToken() string {
	return Instance().Config.AppToken
}

func UserToken() string {
	return Instance().Config.UserToken
}

func ClientID() string {
	return Instance().Config.ClientID
}

func ClientSecret() string {
	return Instance().Config.ClientSecret
}

func ClientRedirectURL() string {
	return Instance().Config.RedirectURL
}
