package core

import (
	"go.temporal.io/sdk/workflow"
)

// QueueCtrl is the main workflow for managing the queue.
func QueueCtrl(ctx workflow.Context, repo *Repo, branch string) error {
	ctx, state := NewQueueCtrlState(ctx, repo, branch)
	selector := workflow.NewSelector(ctx)

	// goroutine to handle signals, enabling uninterrupted addition and reordering of prs in the queue
	workflow.Go(ctx, func(ctx workflow.Context) {
		// setting up add signal for adding prs to the primary queue
		add := workflow.GetSignalChannel(ctx, RepoIOSignalQueueAdd.String())
		selector.AddReceive(add, state.on_add(ctx))

		// setting up add_priority signal for adding prs to the priority queue
		add_priority := workflow.GetSignalChannel(ctx, RepoIOSignalQueueAddPriority.String())
		selector.AddReceive(add_priority, state.on_add_priority(ctx))

		// setting up promote signal for moving a pr up in the queue
		promote := workflow.GetSignalChannel(ctx, RepoIOSignalQueuePromote.String())
		selector.AddReceive(promote, state.on_promote(ctx))

		// setting up demote signal for moving a pr down in the queue
		demote := workflow.GetSignalChannel(ctx, RepoIOSignalQueueDemote.String())
		selector.AddReceive(demote, state.on_demote(ctx))

		for state.is_active() {
			selector.Select(ctx)
		}
	})

	// the main event loop is not implemented yet
	for state.is_active() {
		pr := state.pop(ctx)
		if pr != nil {
			// process the pr
			err := state.process(ctx, pr) // TODO - handle error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// PrCtrl is a child worklow handle the pr activities and handle the pull request.
func PrCtrl(ctx workflow.Context, pr *RepoIOPullRequest) error {
	return nil
}
