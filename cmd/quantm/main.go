package main

import (
	"log/slog"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"

	"go.breu.io/quantm/internal/nomad"
)

type (
	Config struct {
		Nomad *nomad.Config `koanf:"NOMAD"`
	}
)

func read_env() {
	conf := &Config{Nomad: &nomad.DefaultConfig}
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

	slog.Info("conf", slog.Any("nomad", conf.Nomad))
}

func main() {
	read_env()
}
