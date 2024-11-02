package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
	"go.breu.io/durex/queues"
	"go.breu.io/graceful"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/migrations"
	"go.breu.io/quantm/internal/durable"
	githubcfg "go.breu.io/quantm/internal/hooks/github/config"
	pkg_slack "go.breu.io/quantm/internal/hooks/slack/config"
	"go.breu.io/quantm/internal/nomad"
)

type (
	// Config defines the application's configuration.
	Config struct {
		DB      *db.Config        `koanf:"DB"`      // Configuration for the database.
		Durable *durable.Config   `koanf:"DURABLE"` // Configuration for the durable.
		Nomad   *nomad.Config     `koanf:"NOMAD"`   // Configuration for Nomad.
		Github  *githubcfg.Config `koanf:"GITHUB"`  // Configuration for the github.
		Slack   *pkg_slack.Config `koanf:"SLACK"`   // Configuration for the slack.
		Migrate bool              `koanf:"MIGRATE"` // Flag to enable database migration.
	}

	// Service is an interface for services that can be started and stopped.
	Service interface {
		Start(context.Context) error // Starts the service.
		Stop(context.Context) error  // Stops the service.
	}

	// Services represents a list of services.
	Services []Service
)

// main starts the application.
func main() {
	ctx := context.Background()
	release := make(chan any, 1)
	rx_errors := make(chan error)
	timeout := time.Second * 10

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	conf := read_env()                      // Read configuration from environment and flags.
	svcs := make(Services, 0)               // Initialize an empty list of services.
	queues := make(queues.Queues)           // Initialize an empty list of queues.
	cleanups := make([]graceful.Cleanup, 0) // Initialize an empty list of cleanup functions.

	if conf.Durable != nil {
		if err := configure_durable(conf.Durable); err != nil {
			slog.Error("main: unable to start durable service ...", "error", err.Error())

			os.Exit(1)
		}
	} else {
		slog.Warn("main: durable service not configured, this may cause issues ...")
	}

	q_prefix()
	q_hooks(queues)

	svcs = append(svcs, configure_db(conf.DB))
	svcs = append(svcs, configure_nomad(conf.Nomad))

	// If migration is enabled, append the migration service to the list.
	if conf.Migrate {
		if err := migrations.Run(ctx, conf.DB); err != nil {
			slog.Error("main: unable to run migrations, cannot continue ...", "error", err.Error())
		}

		os.Exit(0)
	}

	// Start each service in a goroutine, registering cleanup functions for graceful shutdown.
	for _, svc := range svcs {
		cleanups = append(cleanups, svc.Stop)
		graceful.Go(ctx, graceful.GrabAndGo(svc.Start, ctx), rx_errors)
	}

	for _, q := range queues {
		cleanups = append(cleanups, q.Shutdown)
		graceful.Go(ctx, q.Start, rx_errors)
	}

	// Wait for termination signal or service start error.
	select {
	case rx := <-terminate:
		slog.Info("main: shutdown requested ...", "signal", rx.String())
	case err := <-rx_errors:
		slog.Error("main: unable to start ...", "error", err.Error())
	}

	// Initiate graceful shutdown, waiting for cleanups to complete.
	code := graceful.Shutdown(ctx, cleanups, release, timeout, 0)
	if code == 1 {
		slog.Warn("main: failed to shutdown gracefully, exiting ...")
	} else {
		slog.Info("main: shutdown complete, exiting ...")
	}

	os.Exit(code)
}

// read_env reads configuration from environment variables and default values.
func read_env() *Config {
	conf := &Config{
		DB:      &db.DefaultConfig,
		Durable: &durable.DefaultConfig,
		Nomad:   &nomad.DefaultConfig,
		Migrate: false,
	}

	k := koanf.New("__")

	// Load default values from the Config struct.
	if err := k.Load(structs.Provider(conf, "__"), nil); err != nil {
		panic(err)
	}

	// Load environment variables with the "__" delimiter.
	if err := k.Load(env.Provider("", "__", nil), nil); err != nil {
		panic(err)
	}

	// Unmarshal configuration from the Koanf instance to the Config struct.
	if err := k.Unmarshal("", conf); err != nil {
		panic(err)
	}

	// Add -m or --migrate flag to enable database migration.
	if !conf.Migrate {
		flag.BoolVarP(&conf.Migrate, "migrate", "m", false, "run database migrations")
		flag.Parse()
	}

	return conf
}

// configure_db constructs a database service with the given configuration.
func configure_db(config *db.Config) Service {
	return db.Connection(db.WithConfig(config))
}

func configure_durable(config *durable.Config) error {
	return durable.Configure(durable.WithConfig(config))
}

// configure_nomad constructs a Nomad service with the given configuration.
func configure_nomad(config *nomad.Config) Service {
	return nomad.New(nomad.WithConfig(config))
}
