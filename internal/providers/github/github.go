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
	gh "github.com/google/go-github/v50/github"
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

const (
	KEY = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEA0UG2NPBaqqY+XSZQfFpQm8jUVX8KOUTJhnSpwb02tq7VlugQ
FQEb1+gXIb+/aNRi0KWR/cW5VUlpyRaMoiGzg63FkLxlL5EQCjn3dAfVCiQsWvzb
3hm1NUh1QtDpOoVCXL0veVoJYgEh+ioA5FGYuKBPuqp6Bp1NdGTq+eNBf35JSM6A
Um23l0AUrlyK+qHDFHYBa9Qupn/vk4FDN2vxH1unELhHXTTBFRS1aDYjjGd/wRy6
2HJD8wXqcAZ8F9OY08tetw3+n1gqSsL3xtAmtFbqvhRNxaijMwr8l2CmNZm5vw7V
SFIc9WWtS5rROmzemS3uCp+/qHYZqFlNzTuYJQIDAQABAoIBAAmKYXBQdRHKupUs
pgbFZ19y7JtpS2IJDNcggozev5vcpMhYlEMg5dAWONfFEkkJRegVZG6ZkTWeP0B3
0rmhp7mdNqC+ti5RAtY0hl+367Kmq48KcEvUCDsBrrb5J2kPolLwHTX/MOZS/uWU
/K1sOvZP+NKd6ypaCaoA3+W8wsO5PRrtVq4pxPEVbkfIh7md+d5ogV6fVK/Byah7
LCMAypKY+jkIBzv7BzqZnDFlS5knleUgll0QuFz9p7FJm5RyX3djbXtQuSgZSAlS
fSuCxeEqGMpaXb34tt+0GaeICF44+UNSiM6Agh5Uu2FJHJESZGfsNmRRxxJxjQSd
B2TrJgECgYEA+h38yW2fb8V33QjGwm5bQnvkyNk5AJW/RoWiLMjU3zaxNNYbqOzK
fQAh4EQ8AQQhQhsSHVw4W3iolnk5xIFZcR+nQL1G/j0xzTarw4hZUyFb/or5t2gS
EWwKrjjxT1/64f2DyVQtcxHBzNTW2DCA/OSn4cNQlnm+bMgh6cTX0IUCgYEA1i2x
tImboxO4gejbZPrStl4aqhUdzzbsZ776c3YM0z8Q2C7EWY4JWTCivkO4l960rEft
h2kc/6Khd7IlvvOGf/hvGEHmJarv3K8eh7FGYtRwTX0zQQDY+0tfdid1LH0WfZl2
0Ot2hEqMxN/u/q17AWDZz+rDu3+K08I5yf3kCyECgYEAoCQdGzcGE0FiynH5GLoh
0kKTLInwdlBqxJOBT51StoxFD6ha02CxETHJftcReDEVvkao5YWLS/3IK3f4pbmP
898pbkkCMHwr69GqTip5zsEYLrT6yBRpJSCBAiXRU1oHvzRbccdkxj1DUYug95Cu
tb0NRH6SlZXjd7D4Db4L1CUCgYEAiHN6KNQWtPHGdfV9eTsXbZpMkJl9cVvDh2Ez
vMWz7A3c1G4PKCMGr6z9sgwBGbiIEM6OdNux3uekyVZVF++cfAEx/hlV4B+kS0vC
Pp7hgetoVOXz9nDszES739HJo/tZjdFs0jOBQU0hm/gzEkxB9qHWgtFFvDnIn5q5
KIg5diECgYEA+RieA4ydYCKPboEgMJ7up1Wx3zPmYLCeUnY1cpMZ6rynvGhZFJw6
rh1U79mLewaxDgTehNbyQuDjwK9mYWZ86h/a6iGoqrUPoua8rJCflJJAd8iAtAZ7
CXrNxr5r1sMKe2XPdQ1WMh/xvBkLpbWFajtzttTPPYJf6j8cdEjwqgU=
-----END RSA PRIVATE KEY-----`
)

func (g *Config) GetActivities() *Activities {
	return g.Activities
}

func (g *Config) GetClientForInstallation(installationID int64) (*gh.Client, error) {
	transport, err := ghinstallation.New(http.DefaultTransport, g.AppID, installationID, []byte(KEY))
	if err != nil {
		return nil, err
	}

	client := gh.NewClient(&http.Client{Transport: transport})

	return client, nil
}

func (g *Config) VerifyWebhookSignature(payload []byte, signature string) error {
	// result := g.SignPayload(payload)

	// if result != signature {
	// 	return ErrVerifySignature
	// }

	return nil
}

func (g *Config) SignPayload(payload []byte) string {
	key := hmac.New(sha256.New, []byte(g.WebhookSecret))
	key.Write(payload)
	result := "sha256=" + hex.EncodeToString(key.Sum(nil))

	return result
}

// func (g *github) CloneRepo(repo string, branch string, ref string) {}
