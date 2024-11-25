package github

import (
	"go.breu.io/quantm/internal/hooks/github/activities"
	"go.breu.io/quantm/internal/hooks/github/config"
	"go.breu.io/quantm/internal/hooks/github/web"
	"go.breu.io/quantm/internal/hooks/github/workflows"
)

type (
	InstallActivity      = activities.Install
	InstallReposActivity = activities.InstallRepos
	PushActivity         = activities.Push

	KernelImpl = activities.Kernel

	Config  = config.Config
	Webhook = web.Webhook
)

var (
	Configure  = config.Configure
	WithConfig = config.WithConfig
	Get        = config.Instance

	Install   = workflows.Install
	Push      = workflows.Push
	SyncRepos = workflows.SyncRepos
)
