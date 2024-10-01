package main

import (
	"go.temporal.io/sdk/client"

	"go.breu.io/quantm/internal/core/code"
	"go.breu.io/quantm/internal/core/mutex"
	"go.breu.io/quantm/internal/providers/github"
	"go.breu.io/quantm/internal/providers/slack"
	"go.breu.io/quantm/internal/shared/queue"
)

func configure_core(q queue.Queue, c client.Client) {
	worker := q.Worker(c)

	worker.RegisterWorkflow(code.RepoCtrl)
	worker.RegisterWorkflow(code.TrunkCtrl)
	worker.RegisterWorkflow(code.BranchCtrl)
	worker.RegisterWorkflow(code.QueueCtrl)

	worker.RegisterActivity(&code.Activities{})

	// TODO: this will not work if we have more than one provider.
	worker.RegisterActivity(&github.RepoIO{})
	worker.RegisterActivity(&slack.Activities{})

	worker.RegisterActivity(mutex.PrepareMutexActivity)
}

func configure_provider(q queue.Queue, c client.Client) {
	worker := q.Worker(c)

	github_workflows := &github.Workflows{}

	worker.RegisterWorkflow(github_workflows.OnInstallationEvent)
	worker.RegisterWorkflow(github_workflows.OnInstallationRepositoriesEvent)
	worker.RegisterWorkflow(github_workflows.PostInstall)
	worker.RegisterWorkflow(github_workflows.OnPushEvent)
	worker.RegisterWorkflow(github_workflows.OnCreateOrDeleteEvent)
	worker.RegisterWorkflow(github_workflows.OnPullRequestEvent)
	worker.RegisterWorkflow(github_workflows.OnWorkflowRunEvent)

	worker.RegisterActivity(&github.Activities{})
	worker.RegisterActivity(&slack.Activities{})
}
