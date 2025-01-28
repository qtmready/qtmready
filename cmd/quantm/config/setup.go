package config

import (
	"log/slog"
	"os"

	"go.breu.io/graceful"

	"go.breu.io/quantm/cmd/quantm/workers"
	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/hooks/github"
	"go.breu.io/quantm/internal/hooks/slack"
	"go.breu.io/quantm/internal/nomad"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
	"go.breu.io/quantm/internal/pulse"
)

const (
	ServiceGithub     = "github"
	ServiceSlack      = "slack"
	ServiceKernel     = "kernel"
	ServiceDB         = "db"
	ServicePulse      = "pulse"
	ServiceDurable    = "durable"
	ServiceWebhook    = "webhook"
	ServiceNomad      = "nomad"
	ServiceCoreQueue  = "core_queue"
	ServiceHooksQueue = "hooks_queue"
)

// Setup configures the application based on the provided config.
func (c *Config) Setup(app *graceful.Graceful) error {
	switch c.Mode {
	case ModeMigrate:
		c.SetupLogger()

		if err := c.SetupDB(); err != nil {
			return err
		}

	case ModeWebhook:
		if err := c.SetupServices(app); err != nil {
			return err
		}

		app.Add(ServiceWebhook, NewWebhookServer(), ServiceDurable)
	case ModeGRPC:
		if err := c.SetupServices(app); err != nil {
			return err
		}

		app.Add(ServiceNomad, nomad.New(nomad.WithConfig(c.Nomad)), ServiceKernel, ServiceDB, ServiceDurable, ServicePulse)
	case ModeWorkers:
		if err := c.SetupServices(app); err != nil {
			return err
		}

		workers.Core()
		workers.Hooks()

		app.Add(ServiceCoreQueue, durable.OnCore(), ServiceKernel, ServiceDB, ServiceDurable, ServicePulse)
		app.Add(ServiceHooksQueue, durable.OnHooks(), ServiceKernel, ServiceDB, ServiceDurable, ServicePulse)
	case ModeDefault:
		if err := c.SetupServices(app); err != nil {
			return err
		}

		workers.Core()
		workers.Hooks()

		app.Add(ServiceWebhook, NewWebhookServer(), ServiceDurable)
		app.Add(ServiceNomad, nomad.New(nomad.WithConfig(c.Nomad)), ServiceKernel, ServiceDB, ServiceDurable, ServicePulse)
		app.Add(ServiceCoreQueue, durable.OnCore(), ServiceKernel, ServiceDB, ServiceDurable, ServicePulse)
		app.Add(ServiceHooksQueue, durable.OnHooks(), ServiceKernel, ServiceDB, ServiceDurable, ServicePulse)
	default:
	}

	return nil
}

func (c *Config) SetupLogger() {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}

	handler = slog.NewJSONHandler(os.Stdout, opts)

	if c.Debug {
		opts.Level = slog.LevelDebug
		opts.AddSource = false
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}

// SetupServices configures common services.
func (c *Config) SetupServices(app *graceful.Graceful) error {
	c.SetupLogger()
	auth.SetSecret(c.Secret)

	if err := c.Github.Validate(); err != nil {
		return err
	}

	github.Configure(github.WithConfig(c.Github))

	if err := c.Slack.Validate(); err != nil {
		return err
	}

	slack.Configure(slack.WithConfig(c.Slack))

	kernel.Configure(
		kernel.WithRepoHook(eventsv1.RepoHook_REPO_HOOK_GITHUB, &github.KernelImpl{}),
		kernel.WithChatHook(eventsv1.ChatHook_CHAT_HOOK_SLACK, &slack.KernelImpl{}),
	)

	if err := c.SetupDB(); err != nil {
		return err
	}

	if err := c.SetupDurable(); err != nil {
		return err
	}

	if err := c.SetupPulse(); err != nil {
		return err
	}

	app.Add(ServiceGithub, github.Get())
	// app.Add(ServicesSlack, slack.Get())
	app.Add(ServiceKernel, kernel.Get(), ServiceGithub)
	app.Add(ServiceDB, db.Get())
	app.Add(ServicePulse, pulse.Get())
	app.Add(ServiceDurable, durable.Get())

	return nil
}

// SetupDB configures the database.
func (c *Config) SetupDB() error {
	if err := c.DB.Validate(); err != nil {
		return err
	}

	db.Get(db.WithConfig(c.DB))

	return nil
}

// SetupDurable configures the durable service.
func (c *Config) SetupDurable() error {
	if err := c.Durable.Validate(); err != nil {
		return err
	}

	durable.Get(durable.WithConfig(c.Durable))

	return nil
}

// SetupPulse configures the pulse service.
func (c *Config) SetupPulse() error {
	if err := c.Pulse.Validate(); err != nil {
		return err
	}

	pulse.Get(pulse.WithConfig(c.Pulse))

	return nil
}
