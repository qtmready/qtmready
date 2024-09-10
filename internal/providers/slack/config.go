// Copyright © 2024, Breu, Inc. <info@breu.io>
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
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/slack-go/slack"

	"go.breu.io/quantm/internal/shared"
)

type (
	config struct {
		ClientID     string `env:"SLACK_CLIENT_ID"`
		ClientSecret string `env:"SLACK_CLIENT_SECRET"`
		RedirectURL  string `env:"SLACK_REDIRECT_URL"`
	}

	integration struct {
		Config *config
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
	return &integration{
		Config: c,
	}
}

func (i *integration) GetSlackClient(accessToken string) (*slack.Client, error) {
	lgr := &logger{shared.Logger().WithGroup("slack")}
	client := slack.New(
		accessToken,
		slack.OptionDebug(shared.Service().GetDebug()),
		slack.OptionLog(lgr),
	)

	return client, nil
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
