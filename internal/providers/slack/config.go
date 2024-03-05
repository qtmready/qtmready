package slack

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"golang.org/x/oauth2"
	oauth "golang.org/x/oauth2/slack"

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
		Endpoint:     oauth.Endpoint,
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
