package slack

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"

	"go.breu.io/quantm/internal/shared"
)

type (
	config struct {
		BotToken  string `env:"SLACK_BOT_TOKEN"`
		UserToken string `env:"SLACK_USER_TOKEN"`
		AppToken  string `env:"SLACK_APP_TOKEN"`
	}

	integration struct {
		Config *config
		Client *slack.Client
		Socket *socketmode.Client
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
	client := slack.New(
		c.BotToken,
		slack.OptionDebug(shared.Service().GetDebug()),
		slack.OptionAppLevelToken(c.AppToken),
		slack.OptionLog(logger()),
	)

	socket := socketmode.New(
		client,
		socketmode.OptionDebug(shared.Service().GetDebug()),
		socketmode.OptionLog(logger()),
	)

	return &integration{
		Config: c,
		Client: client,
		Socket: socket,
	}
}

func Client() *slack.Client {
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
