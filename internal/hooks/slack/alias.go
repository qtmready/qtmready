package slack

import (
	"go.breu.io/quantm/internal/hooks/slack/config"
)

type (
	Config = config.Config
)

var (
	WithConfig = config.WithConfig
	Configure  = config.Instance
)
