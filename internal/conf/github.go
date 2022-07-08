package conf

import (
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	gh "github.com/google/go-github/v45/github"
	"github.com/ilyakaznacheev/cleanenv"
)

var Github = &github{}

type github struct {
	AppID         int64  `env:"GITHUB_APP_ID"`
	ClinetID      string `env:"GITHUB_CLIENT_ID"`
	WebhookSecret string `env:"GITHUB_WEBHOOK_SECRET"`
	PrivateKey    string `env:"GITHUB_PRIVATE_KEY"`
	// PrivateKey Base64EncodedValue `env:"GITHUB_PRIVATE_KEY"`
}

func (g *github) ReadConf() {
	cleanenv.ReadEnv(g)
}

func (g *github) GetClientForInstallation(installationID int64) (*gh.Client, error) {
	transport, err := ghinstallation.New(http.DefaultTransport, g.AppID, installationID, []byte(g.PrivateKey))

	if err != nil {
		return nil, err
	}

	client := gh.NewClient(&http.Client{Transport: transport})
	return client, nil
}
