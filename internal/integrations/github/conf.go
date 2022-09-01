package github

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"github.com/ilyakaznacheev/cleanenv"
	"go.breu.io/ctrlplane/internal/shared"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	gh "github.com/google/go-github/v45/github"
)

var (
	Github = &conf{}
)

type conf struct {
	AppID         int64  `env:"GITHUB_APP_ID"`
	ClientID      string `env:"GITHUB_CLIENT_ID"`
	WebhookSecret string `env:"GITHUB_WEBHOOK_SECRET"`
	PrivateKey    string `env:"GITHUB_PRIVATE_KEY"`
}

func (g *conf) ReadEnv() {
	if err := cleanenv.ReadEnv(g); err != nil {
		shared.Logger.Error("Failed to read environment variables ...", "error", err)
	}
}

func (g *conf) GetClientForInstallation(installationID int64) (*gh.Client, error) {
	transport, err := ghinstallation.New(http.DefaultTransport, g.AppID, installationID, []byte(g.PrivateKey))
	if err != nil {
		return nil, err
	}

	client := gh.NewClient(&http.Client{Transport: transport})
	return client, nil
}

func (g *conf) VerifyWebhookSignature(payload []byte, signature string) error {
	key := hmac.New(sha1.New, []byte(g.WebhookSecret))
	key.Write(payload)
	result := "sha1=" + hex.EncodeToString(key.Sum(nil))
	if result != signature {
		return ErrorVerifySignature
	}
	return nil
}
