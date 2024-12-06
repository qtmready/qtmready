package slack

import (
	"go.breu.io/quantm/internal/hooks/slack/activities"
	"go.breu.io/quantm/internal/hooks/slack/config"
	"go.breu.io/quantm/internal/hooks/slack/nomad"
)

type (
	Config = config.Config

	KernelImpl = activities.Kernel
)

var (
	WithConfig = config.WithConfig
	Configure  = config.Instance

	NomadHandler = nomad.NewSlackServiceHandler
)
