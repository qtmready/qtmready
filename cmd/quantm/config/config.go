package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"

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
		DB      *db.Config      `koanf:"DB" json:"db"`           // Configuration for the database.
		Durable *durable.Config `koanf:"DURABLE" json:"durable"` // Configuration for the durable.
		Pulse   *pulse.Config   `koanf:"PULSE" json:"pulse"`     // Configuration for the pulse.
		Nomad   *nomad.Config   `koanf:"NOMAD" json:"nomad"`     // Configuration for Nomad.
		Github  *github.Config  `koanf:"GITHUB" json:"github"`   // Configuration for the github.
		Slack   *slack.Config   `koanf:"SLACK" json:"slack"`     // Configuration for the slack.

		Secret  string `koanf:"SECRET" json:"secret"`   // Secret key for JWE.
		Debug   bool   `koanf:"DEBUG" json:"debug"`     // Flag to enable debug mode.
		Migrate bool   `koanf:"MIGRATE" json:"migrate"` // Flag to enable database migration.

		Mode Mode `koanf:"MODE" json:"mode"`
	}
)

const (
	ModeMigrate Mode = "migrate"
	ModeWebhook Mode = "webhook"
	ModeGRPC    Mode = "grpc"
	ModeWorkers Mode = "queues"
	ModeDefault Mode = "default"
)

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

	slog.Info("config loaded", "config.db", c.DB)
}

func New() *Config {
	return &Config{}
}

// Parse parses command-line flags and sets the application mode.
func (c *Config) Parse() {
	help := false
	count := 0
	selected := ""

	modes := map[string]Mode{
		"migrate": ModeMigrate,
		"webhook": ModeWebhook,
		"grpc":    ModeGRPC,
		"queues":  ModeWorkers,
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
