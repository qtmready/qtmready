package workers

import (
	"go.breu.io/quantm/internal/durable"
	"go.breu.io/quantm/internal/hooks/github"
	"go.breu.io/quantm/internal/pulse"
)

// Hooks registers the activites and workflows for the hooks queue.
func Hooks() {
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
	}
}
