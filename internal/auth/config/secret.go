package config

import (
	"log/slog"
	"sync/atomic"
)

const (
	_default string = "set me"
)

var (
	secret atomic.Value
)

func init() { secret.Store(_default) }

// Secret returns the configured secret value. It will log a warning if the secret is not set.
func Secret() string {
	if !IsValid() {
		slog.Warn("auth: secret is not set, configure it using the environment variable 'SECRET'")
	}

	return secret.Load().(string)
}

// SetSecret sets the secret value.
func SetSecret(val string) {
	secret.Store(val)
}

// IsValid returns true if the secret is valid, false otherwise.
func IsValid() bool {
	return secret.Load().(string) != _default
}
