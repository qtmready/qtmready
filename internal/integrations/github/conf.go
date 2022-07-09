package github

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	gh "github.com/google/go-github/v45/github"
	"github.com/ilyakaznacheev/cleanenv"
)

var Github githubconf

type githubconf struct {
	AppID         int64  `env:"GITHUB_APP_ID"`
	ClinetID      string `env:"GITHUB_CLIENT_ID"`
	WebhookSecret string `env:"GITHUB_WEBHOOK_SECRET"`
	PrivateKey    string `env:"GITHUB_PRIVATE_KEY"`
}

func (g *githubconf) ReadEnv() {
	cleanenv.ReadEnv(g)
}

func (g *githubconf) GetClientForInstallation(installationID int64) (*gh.Client, error) {
	transport, err := ghinstallation.New(http.DefaultTransport, g.AppID, installationID, []byte(g.PrivateKey))
	if err != nil {
		return nil, err
	}

	client := gh.NewClient(&http.Client{Transport: transport})
	return client, nil
}

func (g *githubconf) VerifyWebhookSignature(payload []byte, signature string) error {
	key := hmac.New(sha1.New, []byte(g.WebhookSecret))
	key.Write(payload)
	result := "sha1=" + hex.EncodeToString(key.Sum(nil))
	if result != signature {
		return ErrorVerifySignature
	}
	return nil
}
