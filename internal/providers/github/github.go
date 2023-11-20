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
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"

	"github.com/bradleyfalzon/ghinstallation/v2"
	gh "github.com/google/go-github/v53/github"
	"github.com/ilyakaznacheev/cleanenv"
)

var (
	instance *Config
	once     sync.Once
)

func NewGithub(options ...ConfigOption) *Config {
	g := &Config{}

	for _, option := range options {
		option(g)
	}

	return g
}

func WithAppID(id int64) ConfigOption {
	return func(g *Config) {
		g.AppID = id
	}
}

func WithClientID(id string) ConfigOption {
	return func(g *Config) {
		g.ClientID = id
	}
}

func WithWebhookSecret(secret string) ConfigOption {
	return func(g *Config) {
		g.WebhookSecret = secret
	}
}

func WithPrivateKey(key string) ConfigOption {
	return func(g *Config) {
		g.PrivateKey = key
	}
}

func WithActivities(activities *Activities) ConfigOption {
	return func(g *Config) {
		g.Activities = activities
	}
}

func WithConfigFromEnv() ConfigOption {
	return func(g *Config) {
		if err := cleanenv.ReadEnv(g); err != nil {
			panic(fmt.Errorf("failed to read environment variables: %w", err))
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

type (
	Config struct {
		AppID         int64  `env:"GITHUB_APP_ID"`
		ClientID      string `env:"GITHUB_CLIENT_ID"`
		WebhookSecret string `env:"GITHUB_WEBHOOK_SECRET"`
		PrivateKey    string `env:"GITHUB_PRIVATE_KEY"`
		Activities    *Activities
	}

	ConfigOption func(*Config)
)

func (g *Config) GetActivities() *Activities {
	return g.Activities
}

func (g *Config) GetClientForInstallation(installationID int64) (*gh.Client, error) {
	transport, err := ghinstallation.New(http.DefaultTransport, g.AppID, installationID, []byte(g.PrivateKey))
	if err != nil {
		return nil, err
	}

	client := gh.NewClient(&http.Client{Transport: transport})

	return client, nil
}

func (g *Config) VerifyWebhookSignature(payload []byte, signature string) error {
	result := g.SignPayload(payload)

	if result != signature {
		return ErrVerifySignature
	}

	return nil
}

func (g *Config) SignPayload(payload []byte) string {
	key := hmac.New(sha256.New, []byte(g.WebhookSecret))
	key.Write(payload)
	result := "sha256=" + hex.EncodeToString(key.Sum(nil))

	return result
}

// func (g *github) CloneRepo(repo string, branch string, ref string) {}
