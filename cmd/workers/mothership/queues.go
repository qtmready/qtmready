package main

import (
	"go.breu.io/quantm/internal/core/code"
	"go.breu.io/quantm/internal/core/mutex"
	"go.breu.io/quantm/internal/providers/github"
	"go.breu.io/quantm/internal/providers/slack"
	queue "go.breu.io/quantm/internal/shared/queue"
)

func configure_core() {
	queue.Core().RegisterWorkflow(code.RepoCtrl)
	queue.Core().RegisterWorkflow(code.TrunkCtrl)
	queue.Core().RegisterWorkflow(code.BranchCtrl)
	queue.Core().RegisterWorkflow(code.QueueCtrl)

	queue.Core().RegisterActivity(&code.Activities{})

	queue.Core().RegisterActivity(&github.RepoIO{})
	queue.Core().RegisterActivity(&slack.Activities{})
	queue.Core().RegisterActivity(mutex.PrepareMutexActivity)
}

func configure_providers() {
	github_workflows := &github.Workflows{}

	queue.Providers().RegisterWorkflow(github_workflows.OnInstallationEvent)
	queue.Providers().RegisterWorkflow(github_workflows.OnInstallationRepositoriesEvent)
	queue.Providers().RegisterWorkflow(github_workflows.PostInstall)
	queue.Providers().RegisterWorkflow(github_workflows.OnPushEvent)
	queue.Providers().RegisterWorkflow(github_workflows.OnCreateOrDeleteEvent)
	queue.Providers().RegisterWorkflow(github_workflows.OnPullRequestEvent)
	queue.Providers().RegisterWorkflow(github_workflows.CollectRepoEventMetadata)
	queue.Providers().RegisterWorkflow(github_workflows.OnWorkflowRunEvent)

	queue.Providers().RegisterActivity(&github.Activities{})
	queue.Providers().RegisterActivity(&slack.Activities{})
}

func configure_mutex() {
	queue.Mutex().RegisterWorkflow(mutex.MutexWorkflow)
}
