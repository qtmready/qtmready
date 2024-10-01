// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2022, 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
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

package github

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"

	"github.com/bradleyfalzon/ghinstallation/v2"
	gh "github.com/google/go-github/v62/github"
	"github.com/ilyakaznacheev/cleanenv"

	"go.breu.io/quantm/internal/shared"
)

type (
	// Config holds configuration settings for the GitHub integration.
	Config struct {
		AppID              int64  `env:"GITHUB_APP_ID"`
		ClientID           string `env:"GITHUB_CLIENT_ID"`
		WebhookSecret      string `env:"GITHUB_WEBHOOK_SECRET"`
		PrivateKey         string `env:"GITHUB_PRIVATE_KEY"`
		PrivateKeyIsBase64 bool   `env:"GITHUB_PRIVATE_KEY_IS_BASE64" env-default:"false"` // If true, the private key is base64 encoded

	}

	// ConfigOption represents a function that modifies the GitHub configuration.
	ConfigOption func(*Config)
)

var (
	// instance stores the singleton instance of the GitHub configuration.
	instance               *Config
	once                   sync.Once
	actionWorkflowStatuses map[string]map[string]string
)

// NewGithub creates a new instance of the GitHub configuration with the provided options.
//
// The function initializes the `actionWorkflowStatuses` map and returns the configured instance.
func NewGithub(options ...ConfigOption) *Config {
	g := &Config{}

	for _, option := range options {
		option(g)
	}

	actionWorkflowStatuses = make(map[string]map[string]string)

	return g
}

// WithAppID sets the GitHub App ID in the configuration.
func WithAppID(id int64) ConfigOption {
	return func(config *Config) {
		config.AppID = id
	}
}

// WithClientID sets the GitHub Client ID in the configuration.
func WithClientID(id string) ConfigOption {
	return func(config *Config) {
		config.ClientID = id
	}
}

// WithWebhookSecret sets the GitHub Webhook Secret in the configuration.
func WithWebhookSecret(secret string) ConfigOption {
	return func(config *Config) {
		config.WebhookSecret = secret
	}
}

// WithPrivateKey sets the GitHub Private Key in the configuration.
func WithPrivateKey(key string) ConfigOption {
	return func(config *Config) {
		config.PrivateKey = key
	}
}

// WithConfigFromEnv reads the GitHub configuration from environment variables.
//
// This function reads the configuration from the environment variables and sets them in the provided config struct. It also
// handles base64 decoding of the private key if it is encoded in base64.
func WithConfigFromEnv() ConfigOption {
	return func(config *Config) {
		if err := cleanenv.ReadEnv(config); err != nil {
			panic(fmt.Errorf("failed to read environment variables: %w", err))
		}

		if config.PrivateKeyIsBase64 {
			key, err := base64.StdEncoding.DecodeString(config.PrivateKey)
			if err != nil {
				panic(fmt.Errorf("failed to decode base64 private key: %w", err))
			}

			config.PrivateKey = string(key)
		}
	}
}

// Instance returns the singleton instance of the GitHub configuration.
//
// The function uses a `sync.Once` to ensure that the configuration is initialized only once. It initializes the
// instance with the `WithConfigFromEnv` option, which reads the configuration from environment variables.
func Instance() *Config {
	if instance == nil {
		once.Do(func() {
			instance = NewGithub(
				WithConfigFromEnv(),
			)
		})
	}

	return instance
}

// GetClientForInstallationID retrieves a GitHub client for a specific installation ID.
//
// The function uses the GitHub Installation API to create a new client for the specified installation ID. It uses the
// private key provided in the configuration to authenticate with the GitHub API.
func (config *Config) GetClientForInstallationID(installationID shared.Int64) (*gh.Client, error) {
	transport, err := ghinstallation.New(http.DefaultTransport, config.AppID, installationID.Int64(), []byte(config.PrivateKey))
	if err != nil {
		return nil, err
	}

	client := gh.NewClient(&http.Client{Transport: transport})

	return client, nil
}

// VerifyWebhookSignature verifies the signature of a webhook payload.
//
// The function verifies that the provided signature matches the signature generated by signing the payload with the
// webhook secret. It returns an error if the signatures don't match.
func (config *Config) VerifyWebhookSignature(payload []byte, signature string) error {
	result := config.SignPayload(payload)

	if result != signature {
		return ErrVerifySignature
	}

	return nil
}

// SignPayload generates a signature for a given payload.
//
// The function calculates the HMAC-SHA256 hash of the payload using the webhook secret provided in the configuration. It
// returns the base64 encoded signature in the format "sha256=<hash>".
func (config *Config) SignPayload(payload []byte) string {
	key := hmac.New(sha256.New, []byte(config.WebhookSecret))
	key.Write(payload)
	result := "sha256=" + hex.EncodeToString(key.Sum(nil))

	return result
}
