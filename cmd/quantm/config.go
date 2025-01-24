package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
	"go.breu.io/graceful"

	"go.breu.io/quantm/internal/core/kernel"
	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/hooks/github"
	"go.breu.io/quantm/internal/hooks/slack"
	"go.breu.io/quantm/internal/nomad"
	"go.breu.io/quantm/internal/pulse"
)

type (
	Mode string

	Config struct {
		DB      *db.Config      `koanf:"DB"`      // Configuration for the database.
		Durable *durable.Config `koanf:"DURABLE"` // Configuration for the durable.
		Pulse   *pulse.Config   `koanf:"PULSE"`   // Configuration for the pulse.
		Nomad   *nomad.Config   `koanf:"NOMAD"`   // Configuration for Nomad.
		Github  *github.Config  `koanf:"GITHUB"`  // Configuration for the github.
		Slack   *slack.Config   `koanf:"SLACK"`   // Configuration for the slack.

		Secret  string `koanf:"SECRET"`  // Secret key for JWE.
		Debug   bool   `koanf:"DEBUG"`   // Flag to enable debug mode.
		Migrate bool   `koanf:"MIGRATE"` // Flag to enable database migration.

		Mode Mode `koanf:"MODE"`
	}
)

const (
	ModeMigrate Mode = "migrate"
	ModeWebhook Mode = "webhook"
	ModeGRPC    Mode = "grpc"
	ModeQueues  Mode = "queues"
	ModeDefault Mode = "default"
)

// - Config Load -

// Load loads the configuration from different providers.
func (c *Config) Load() {
	c.DB = &db.DefaultConfig
	c.Durable = &durable.DefaultConfig
	c.Nomad = &nomad.DefaultConfig
	c.Pulse = &pulse.DefaultConfig
	c.Github = &github.Config{}
	c.Slack = &slack.Config{}

	k := koanf.New("__")

	// Load default values from the Config struct.
	if err := k.Load(structs.Provider(c, "__"), nil); err != nil {
		panic(err)
	}

	// Load environment variables with the "__" delimiter.
	if err := k.Load(env.Provider("", "__", nil), nil); err != nil {
		panic(err)
	}

	// Unmarshal configuration from the Koanf instance to the Config struct.
	if err := k.Unmarshal("", c); err != nil {
		panic(err)
	}

	slog.Info("config", "config", c)
}

// - Flags -

// Parse parses command-line flags and sets the application mode.
func (c *Config) Parse(app *graceful.Graceful) {
	help := false
	count := 0
	selected := ""

	modes := map[string]Mode{
		"migrate": ModeMigrate,
		"webhook": ModeWebhook,
		"grpc":    ModeGRPC,
		"queues":  ModeQueues,
	}

	flag.BoolVarP(&help, "help", "h", false, "show help message")

	flags := map[string]*bool{
		"migrate": flag.BoolP("migrate", "m", false, "run database migrations"),
		"webhook": flag.BoolP("webhook", "w", false, "start webhook server"),
		"grpc":    flag.BoolP("grpc", "g", false, "start gRPC server (nomad)"),
		"queues":  flag.BoolP("queues", "q", false, "start queues worker"),
	}

	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	for mode, ptr := range flags {
		if *ptr {
			count++
			selected = mode
		}
	}

	if count > 1 {
		panic("only one mode can be enabled at a time")
	}

	if selected != "" {
		if mode, ok := modes[selected]; ok {
			c.Mode = mode
		} else {
			panic(fmt.Sprintf("invalid mode selected: %s", selected))
		}
	} else {
		c.Mode = ModeDefault
	}
}

// - Dependencies -

// Dependencies sets up the application dependencies based on the selected mode.
func (c *Config) Dependencies(app *graceful.Graceful) {
	deps := map[Mode]func(app *graceful.Graceful){
		ModeMigrate: c.migrate,
		ModeWebhook: c.webhook,
		ModeGRPC:    c.nomad,
		ModeQueues:  c.queues,
		ModeDefault: c.all,
	}

	if fn, ok := deps[c.Mode]; ok {
		fn(app)
	} else {
		panic(fmt.Sprintf("invalid mode selected: %s", c.Mode))
	}
}

// - Setup -

// setup_common sets up the common dependencies for most modes.
func (c *Config) setup_common(app *graceful.Graceful) {
	app.Add(DB, db.Connection(db.WithConfig(c.DB)))
	app.Add(Pulse, pulse.Instance(pulse.WithConfig(c.Pulse)))
	app.Add(Github, github.Get())
}

// setup_durable sets up the durable dependency for relevant modes.
func (c *Config) setup_durable(app *graceful.Graceful) {
	app.Add(Durable, durable.Instance())
}

// setup_queues sets up the core and hooks queues.
func (c *Config) setup_queues(app *graceful.Graceful) {
	app.Add(CoreQ, durable.OnCore(), DB, Durable, Pulse, Github)
	app.Add(HooksQ, durable.OnHooks(), DB, Durable, Pulse, Github)
	app.Add(Kernel, kernel.Get(), Github)
}

// setup_nomad sets up the nomad service.
func (c *Config) setup_nomad(app *graceful.Graceful) {
	nmd := nomad.New(nomad.WithConfig(c.Nomad))
	app.Add(Nomad, nmd, DB, Durable, Pulse, Github)
}

// setup_webhook sets up the webhook server.
func (c *Config) setup_webhook(app *graceful.Graceful) {
	app.Add(Webhook, NewWebhookServer(), DB, Durable, Github)
}

// - Mode Handlers -

// migrate sets up dependencies for the migrate mode.
func (c *Config) migrate(app *graceful.Graceful) {
	c.setup_common(app)
}

// webhook sets up dependencies for the webhook mode.
func (c *Config) webhook(app *graceful.Graceful) {
	c.setup_common(app)
	c.setup_webhook(app)
}

// nomad sets up dependencies for the nomad mode.
func (c *Config) nomad(app *graceful.Graceful) {
	c.setup_common(app)
	c.setup_durable(app)
	c.setup_nomad(app)
}

// queues sets up dependencies for the queues mode.
func (c *Config) queues(app *graceful.Graceful) {
	c.setup_common(app)
	c.setup_durable(app)
	c.setup_queues(app)
}

// all sets up dependencies for the default (all) mode.
func (c *Config) all(app *graceful.Graceful) {
	c.setup_common(app)
	c.setup_durable(app)
	c.setup_nomad(app)
	c.setup_queues(app)
	c.setup_webhook(app)
}
