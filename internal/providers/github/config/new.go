package config

type (
	// ConfigOption represents a function that modifies the GitHub configuration.
	ConfigOption func(*config)
)

// WithAppID sets the GitHub App ID in the configuration.
func WithAppID(id int64) ConfigOption {
	return func(config *config) {
		config.AppID = id
	}
}

// WithClientID sets the GitHub Client ID in the configuration.
func WithClientID(id string) ConfigOption {
	return func(config *config) {
		config.ClientID = id
	}
}

// WithWebhookSecret sets the GitHub Webhook Secret in the configuration.
func WithWebhookSecret(secret string) ConfigOption {
	return func(config *config) {
		config.WebhookSecret = secret
	}
}

// WithPrivateKey sets the GitHub Private Key in the configuration.
func WithPrivateKey(key string) ConfigOption {
	return func(config *config) {
		config.PrivateKey = key
	}
}

// New creates a new GitHub configuration with the given options.
func New(opts ...ConfigOption) *config {
	cfg := &config{}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}
