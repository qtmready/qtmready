package main

import (
	"go.breu.io/quantm/internal/core/code"
	"go.breu.io/quantm/internal/core/mutex"
	"go.breu.io/quantm/internal/providers/github"
	"go.breu.io/quantm/internal/providers/slack"
	"go.breu.io/quantm/internal/shared/queues"
)

func configure_core() {
	queues.Core().RegisterWorkflow(code.RepoCtrl)
	queues.Core().RegisterWorkflow(code.TrunkCtrl)
	queues.Core().RegisterWorkflow(code.BranchCtrl)
	queues.Core().RegisterWorkflow(code.QueueCtrl)

	queues.Core().RegisterActivity(&code.Activities{})

	queues.Core().RegisterActivity(&github.RepoIO{})
	queues.Core().RegisterActivity(&slack.Activities{})
	queues.Core().RegisterActivity(mutex.PrepareMutexActivity)
}

func configure_providers() {
	github_workflows := &github.Workflows{}

	queues.Providers().RegisterWorkflow(github_workflows.OnInstallationEvent)
	queues.Providers().RegisterWorkflow(github_workflows.OnInstallationRepositoriesEvent)
	queues.Providers().RegisterWorkflow(github_workflows.PostInstall)
	queues.Providers().RegisterWorkflow(github_workflows.OnPushEvent)
	queues.Providers().RegisterWorkflow(github_workflows.OnCreateOrDeleteEvent)
	queues.Providers().RegisterWorkflow(github_workflows.OnPullRequestEvent)
	queues.Providers().RegisterWorkflow(github_workflows.CollectRepoEventMetadata)
	queues.Providers().RegisterWorkflow(github_workflows.OnWorkflowRunEvent)

	queues.Providers().RegisterActivity(&github.Activities{})
	queues.Providers().RegisterActivity(&slack.Activities{})
}

func configure_mutex() {
	queues.Mutex().RegisterWorkflow(mutex.MutexWorkflow)
}
