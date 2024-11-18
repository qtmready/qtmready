package slackcfg

import (
	"log/slog"
	"sync"

	"github.com/slack-go/slack"
)

var (
	_once sync.Once
	_c    *Config
)

// Config holds the configuration for the Slack client.
type (
	Config struct {
		ClientID     string `koanf:"CLIENT_ID"`
		ClientSecret string `koanf:"CLIENT_SECRET"`
		RedirectURL  string `koanf:"REDIRECT_URL"`
		Debug        bool   `koanf:"DEBUG"`
	}

	ConfigOption func(*Config)
)

// GetSlackClient creates a new Slack client using the token.
func GetSlackClient(token string) (*slack.Client, error) {
	lgr := &logger{slog.Default().WithGroup("slack")}
	client := slack.New(
		token,
		slack.OptionDebug(_c.Debug),
		slack.OptionLog(lgr),
	)

	return client, nil
}

// ClientID returns the Slack client ID.
func ClientID() string {
	return Instance().ClientID
}

// ClientSecret returns the Slack client secret.
func ClientSecret() string {
	return Instance().ClientSecret
}

// ClientRedirectURL returns the Slack redirect URL.
func ClientRedirectURL() string {
	return Instance().RedirectURL
}

func WithConfig(options ...ConfigOption) *Config {
	c := Instance()

	for _, option := range options {
		option(c)
	}

	return c
}

func Instance(opts ...ConfigOption) *Config {
	_once.Do(func() {
		_c = &Config{}

		for _, opt := range opts {
			opt(_c)
		}
	})

	return _c
}
