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
	"go.breu.io/graceful"

	pkg_db "go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/db/migrations"
	pkg_durable "go.breu.io/quantm/internal/durable"
	pkg_github "go.breu.io/quantm/internal/hooks/github/config"
	pkg_slack "go.breu.io/quantm/internal/hooks/slack/config"
	pkg_nomad "go.breu.io/quantm/internal/nomad"
)

type (
	// Config defines the application's configuration.
	Config struct {
		DB      *pkg_db.Config      `koanf:"DB"`      // Configuration for the database.
		Durable *pkg_durable.Config `koanf:"DURABLE"` // Configuration for the durable.
		Nomad   *pkg_nomad.Config   `koanf:"NOMAD"`   // Configuration for Nomad.
		Github  *pkg_github.Config  `koanf:"GITHUB"`  // Configuration for the github.
		Slack   *pkg_slack.Config   `koanf:"SLACK"`   // Configuration for the slack.
		Migrate bool                `koanf:"MIGRATE"` // Flag to enable database migration.
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
	// Initialize context, channels, and timeout.
	ctx := context.Background()
	release := make(chan any, 1)
	rx_errors := make(chan error)
	timeout := time.Second * 10

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	conf := configure()                     // Read configuration from environment and flags.
	svcs := make(Services, 0)               // Initialize an empty list of services.
	cleanups := make([]graceful.Cleanup, 0) // Initialize an empty list of cleanup functions.

	if conf.Durable != nil {
		if err := durable(conf.Durable); err != nil {
			slog.Error("main: unable to start durable service ...", "error", err.Error())

			os.Exit(1)
		}
	} else {
		slog.Warn("main: durable service not configured, this may cause issues ...")
	}

	svcs = append(svcs, db(conf.DB))
	// Append Nomad and database services to the list.
	svcs = append(svcs, nomad(conf.Nomad))

	// If migration is enabled, append the migration service to the list.
	if conf.Migrate {
		if err := migrations.Run(ctx, conf.DB); err != nil {
			slog.Error("main: unable to run migrations, cannot continue ...", "error", err.Error())
		}

		os.Exit(1)
	}

	// Start each service in a goroutine, registering cleanup functions for graceful shutdown.
	for _, svc := range svcs {
		cleanups = append(cleanups, svc.Stop)
		graceful.Go(ctx, graceful.GrabAndGo(svc.Start, ctx), rx_errors)
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

// nomad constructs a Nomad service with the given configuration.
func nomad(config *pkg_nomad.Config) Service {
	return pkg_nomad.New(pkg_nomad.WithConfig(config))
}

// db constructs a database service with the given configuration.
func db(config *pkg_db.Config) Service {
	return pkg_db.Connection(pkg_db.WithConfig(config))
}

func durable(config *pkg_durable.Config) error {
	return pkg_durable.Configure(pkg_durable.WithConfig(config))
}

// configure reads configuration from environment variables and default values.
func configure() *Config {
	conf := &Config{
		DB:      &pkg_db.DefaultConfig,
		Durable: &pkg_durable.DefaultConfig,
		Nomad:   &pkg_nomad.DefaultConfig,
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
