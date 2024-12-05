package slack

import (
	"go.breu.io/quantm/internal/hooks/slack/activities"
	"go.breu.io/quantm/internal/hooks/slack/config"
)

type (
	Config = config.Config

	KernelImpl = activities.Kernel
)

var (
	WithConfig = config.WithConfig
	Configure  = config.Instance
)
