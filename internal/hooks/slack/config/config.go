package config

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

func ClientID() string {
	return _c.ClientID
}

func ClientSecret() string {
	return _c.ClientSecret
}

func ClientRedirectURL() string {
	return _c.RedirectURL
}

func WithConfig(cfg *Config) ConfigOption {
	return func(config *Config) {
		config.ClientID = cfg.ClientID
		config.ClientSecret = cfg.ClientSecret
		config.RedirectURL = cfg.RedirectURL
	}
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
