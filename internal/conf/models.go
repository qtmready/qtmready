package conf

import (
	_tclient "go.temporal.io/sdk/client"
)

type githubConf struct {
	AppID         string `env:"GITHUB_APP_ID"`
	ClinetID      string `env:"GITHUB_CLIENT_ID"`
	WebhookSecret string `env:"GITHUB_WEBHOOK_SECRET"`
	PrivateKey    string `env:"GITHUB_PRIVATE_KEY"`
	// PrivateKey Base64EncodedValue `env:"GITHUB_PRIVATE_KEY"`
}

type kratosConf struct {
	ServerUrl string `env:"KRATOS_SERVER_URL"`
}

type temporal struct {
	Client _tclient.Client
	Queues struct {
		Webhooks string `env-default:"webhooks"`
	}
}

type service struct {
	Name    string
	Debug   bool   `env:"DEBUG" env-default:"false"`
	Version string `env:"VERSION" env-default:"0.0.0-dev"`
}
