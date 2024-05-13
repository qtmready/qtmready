// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

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
	gh "github.com/google/go-github/v53/github"
	"github.com/ilyakaznacheev/cleanenv"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/mutex"
)

type (
	Config struct {
		AppID              int64  `env:"GITHUB_APP_ID"`
		ClientID           string `env:"GITHUB_CLIENT_ID"`
		WebhookSecret      string `env:"GITHUB_WEBHOOK_SECRET"`
		PrivateKey         string `env:"GITHUB_PRIVATE_KEY"`
		PrivateKeyIsBase64 bool   `env:"GITHUB_PRIVATE_KEY_IS_BASE64" env-default:"false"` // If true, the private key is base64 encoded
		Activities         *Activities
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

func WithActivities(activities *Activities) ConfigOption {
	return func(config *Config) {
		config.Activities = activities
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
				WithActivities(&Activities{}),
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
			mutex.WithHandler(ctx),
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

func (config *Config) GetActivities() *Activities {
	return config.Activities
}

func (config *Config) GetClientFromInstallation(installationID int64) (*gh.Client, error) {
	transport, err := ghinstallation.New(http.DefaultTransport, config.AppID, installationID, []byte(config.PrivateKey))
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
