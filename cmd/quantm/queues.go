package main

import (
	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/core/repos"
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/hooks/github"
	"go.breu.io/quantm/internal/pulse"
)

func q_core() {
	q := durable.OnCore()

	q.CreateWorker(
		queues.WithWorkerOptionEnableSessionWorker(true),
	)

	if q != nil {
		q.RegisterActivity(pulse.PersistRepoEvent)
		q.RegisterActivity(pulse.PersistChatEvent)

		q.RegisterWorkflow(repos.RepoWorkflow)
		q.RegisterActivity(repos.NewRepoActivities())

		q.RegisterWorkflow(repos.BranchWorkflow)
		q.RegisterActivity(repos.NewBranchActivities())
	}
}

func q_hooks() {
	q := durable.OnHooks()

	q.CreateWorker()

	if q != nil {
		q.RegisterActivity(pulse.PersistRepoEvent)
		q.RegisterActivity(pulse.PersistChatEvent)

		q.RegisterWorkflow(github.Install)
		q.RegisterActivity(&github.InstallActivity{})

		q.RegisterWorkflow(github.SyncRepos)
		q.RegisterActivity(&github.InstallReposActivity{})

		q.RegisterWorkflow(github.Push)
		q.RegisterActivity(&github.PushActivity{})
	}
}
