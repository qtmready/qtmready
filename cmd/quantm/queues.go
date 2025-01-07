package main

import (
	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/core/repos"
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/hooks/github"
	"go.breu.io/quantm/internal/pulse"
)

// q_core sets up the core queue.
func q_core() {
	q := durable.OnCore()

	q.CreateWorker(
		queues.WithWorkerOptionEnableSessionWorker(true),
	)

	if q != nil {
		// Register core activities
		q.RegisterActivity(pulse.PersistRepoEvent)
		q.RegisterActivity(pulse.PersistChatEvent)

		// Register repo workflows and activities
		q.RegisterWorkflow(repos.RepoWorkflow)
		q.RegisterActivity(repos.NewRepoActivities())

		// Register branch workflows and activities
		q.RegisterWorkflow(repos.BranchWorkflow)
		q.RegisterActivity(repos.NewBranchActivities())
	}
}

// q_hooks sets up the hooks queue.
func q_hooks() {
	q := durable.OnHooks()

	q.CreateWorker()

	if q != nil {
		// Register pulse activities
		q.RegisterActivity(pulse.PersistRepoEvent)
		q.RegisterActivity(pulse.PersistChatEvent)

		// Register github install workflow and activity
		q.RegisterWorkflow(github.InstallWorkflow)
		q.RegisterActivity(&github.InstallActivity{})

		// Register github sync repos workflow and activity
		q.RegisterWorkflow(github.SyncReposWorkflow)
		q.RegisterActivity(&github.InstallReposActivity{})

		// Register github push workflow and activity
		q.RegisterWorkflow(github.PushWorkflow)
		q.RegisterActivity(&github.PushActivity{})

		// Register github ref workflow and activity
		q.RegisterWorkflow(github.RefWorkflow)
		q.RegisterActivity(&github.RefActivity{})

		// Register github pull request workflow and activity
		q.RegisterWorkflow(github.PullRequestWorkflow)
		q.RegisterActivity(&github.PullRequestActivity{})
	}
}
