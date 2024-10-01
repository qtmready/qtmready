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
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	gh "github.com/google/go-github/v62/github"
	"github.com/ilyakaznacheev/cleanenv"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/mutex"
	"go.breu.io/quantm/internal/shared"
)

type (
	Config struct {
		AppID              int64  `env:"GITHUB_APP_ID"`
		ClientID           string `env:"GITHUB_CLIENT_ID"`
		WebhookSecret      string `env:"GITHUB_WEBHOOK_SECRET"`
		PrivateKey         string `env:"GITHUB_PRIVATE_KEY"`
		PrivateKeyIsBase64 bool   `env:"GITHUB_PRIVATE_KEY_IS_BASE64" env-default:"false"` // If true, the private key is base64 encoded

	}

	ConfigOption func(*Config)
)

var (
	instance               *Config
	once                   sync.Once
	lockRepo               map[string]mutex.Mutex
	actionWorkflowStatuses map[string]map[string]string // github repo -> workflow file -> status (idle, requested, in_progress, completed)
)

func NewGithub(options ...ConfigOption) *Config {
	g := &Config{}

	for _, option := range options {
		option(g)
	}

	actionWorkflowStatuses = make(map[string]map[string]string)

	return g
}

func WithAppID(id int64) ConfigOption {
	return func(config *Config) {
		config.AppID = id
	}
}

func WithClientID(id string) ConfigOption {
	return func(config *Config) {
		config.ClientID = id
	}
}

func WithWebhookSecret(secret string) ConfigOption {
	return func(config *Config) {
		config.WebhookSecret = secret
	}
}

func WithPrivateKey(key string) ConfigOption {
	return func(config *Config) {
		config.PrivateKey = key
	}
}

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

func LockInstance(ctx workflow.Context, repoID string) (mutex.Mutex, error) {
	lockID := "repo." + repoID

	lock, exists := lockRepo[lockID]
	if !exists {
		lock = mutex.New(
			ctx,
			mutex.WithTimeout(10*time.Second),
			mutex.WithResourceID(lockID),
		)

		if err := lock.Prepare(ctx); err != nil {
			return nil, err
		}

		if err := lock.Acquire(ctx); err != nil {
			return nil, err
		}
	}

	return lock, nil
}

func (config *Config) GetClientForInstallationID(installationID shared.Int64) (*gh.Client, error) {
	transport, err := ghinstallation.New(http.DefaultTransport, config.AppID, installationID.Int64(), []byte(config.PrivateKey))
	if err != nil {
		return nil, err
	}

	client := gh.NewClient(&http.Client{Transport: transport})

	return client, nil
}

func (config *Config) VerifyWebhookSignature(payload []byte, signature string) error {
	result := config.SignPayload(payload)

	if result != signature {
		return ErrVerifySignature
	}

	return nil
}

func (config *Config) SignPayload(payload []byte) string {
	key := hmac.New(sha256.New, []byte(config.WebhookSecret))
	key.Write(payload)
	result := "sha256=" + hex.EncodeToString(key.Sum(nil))

	return result
}

// func (g *github) CloneRepo(repo string, branch string, ref string) {}
