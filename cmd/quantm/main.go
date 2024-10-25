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
	"go.breu.io/graceful"

	"go.breu.io/quantm/internal/db/config"
	pkg_nomad "go.breu.io/quantm/internal/nomad"
)

type (
	Config struct {
		Nomad *pkg_nomad.Config  `koanf:"NOMAD"`
		DB    *config.Connection `koanf:"DB"`
	}

	Service interface {
		Start(context.Context) error
		Stop(context.Context) error
	}

	Services []Service
)

func nomad(config *pkg_nomad.Config) Service {
	return pkg_nomad.New(pkg_nomad.WithConfig(config))
}

func read_env() *Config {
	conf := &Config{
		Nomad: &pkg_nomad.DefaultConfig,
		DB:    &config.DefaultConn,
	}
	k := koanf.New("__")

	if err := k.Load(structs.Provider(conf, "__"), nil); err != nil {
		panic(err)
	}

	if err := k.Load(env.Provider("", "__", nil), nil); err != nil {
		panic(err)
	}

	if err := k.Unmarshal("", conf); err != nil {
		panic(err)
	}

	return conf
}

func main() {
	conf := read_env()

	svcs := make(Services, 0)
	cleanups := make([]graceful.Cleanup, 0)

	ctx := context.Background()
	release := make(chan any, 1)
	rx_errors := make(chan error)
	timeout := time.Second * 10

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	svcs = append(svcs, nomad(conf.Nomad))

	for _, svc := range svcs {
		cleanups = append(cleanups, svc.Stop)
		graceful.Go(ctx, graceful.GrabAndGo(svc.Start, ctx), rx_errors)
	}

	select {
	case rx := <-terminate:
		slog.Info("main: shutdown requested ...", "signal", rx.String())
	case err := <-rx_errors:
		slog.Error("main: unable to start ...", "error", err.Error())
	}

	code := graceful.Shutdown(ctx, cleanups, release, timeout, 0)
	if code == 1 {
		slog.Warn("main: failed to shutdown gracefully, exiting ...")
	} else {
		slog.Info("main: shutdown complete, exiting ...")
	}

	os.Exit(code)
}
