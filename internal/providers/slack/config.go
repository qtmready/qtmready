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
