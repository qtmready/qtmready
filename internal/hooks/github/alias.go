package github

import (
	"go.breu.io/quantm/internal/hooks/github/activities"
	"go.breu.io/quantm/internal/hooks/github/config"
	"go.breu.io/quantm/internal/hooks/github/nomad"
	"go.breu.io/quantm/internal/hooks/github/web"
	"go.breu.io/quantm/internal/hooks/github/workflows"
)

type (
	InstallActivity      = activities.Install
	InstallReposActivity = activities.InstallRepos
	PushActivity         = activities.Push
	RefActivity          = activities.Ref
	PullRequestActivity  = activities.PullRequest

	KernelImpl = activities.Kernel

	Config  = config.Config
	Webhook = web.Webhook
)

var (
	Configure  = config.Configure
	WithConfig = config.WithConfig
	Get        = config.Instance

	InstallWorkflow     = workflows.Install
	RefWorkflow         = workflows.Ref
	PushWorkflow        = workflows.Push
	PullRequestWorkflow = workflows.PullRequest
	SyncReposWorkflow   = workflows.SyncRepos

	NomadHandler = nomad.NewGithubServiceHandler
)
