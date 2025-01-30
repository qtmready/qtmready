package workers

import (
	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/core/repos"
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/pulse"
)

// Core registers the core activities and workflows at the core queue.
//
// In future we might split this into multiple workers depending on how the system scales.
func Core() {
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

		// Register trunk workflows and activities
		q.RegisterWorkflow(repos.TrunkWorkflow)

		// Register notify activities
		q.RegisterActivity(repos.NewNotifyActivities())
	}
}
