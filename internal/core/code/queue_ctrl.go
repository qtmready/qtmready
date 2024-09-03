package code

import (
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/defs"
)

// QueueCtrl processes PRs sequentially, ensuring only one PR is handled at a time.
//
// Queue modifications (add, promote, demote) are handled concurrently via signals,
// allowing for uninterrupted queue management during PR processing.
func QueueCtrl(ctx workflow.Context, repo *defs.Repo, branch string, queues *QueueCtrlSerializedState) error {
	ctx, state := NewQueueCtrlState(ctx, repo, branch)

	// Deserialize the state if provided
	if queues != nil {
		state.deserialize(ctx, queues)
	}

	selector := workflow.NewSelector(ctx)

	// goroutine to handle signals, enabling uninterrupted addition and reordering of prs in the queue
	workflow.Go(ctx, func(ctx workflow.Context) {
		// setting up add signal for adding prs to the primary queue
		add := workflow.GetSignalChannel(ctx, defs.RepoIOSignalQueueAdd.String())
		selector.AddReceive(add, state.on_add(ctx))

		priority := workflow.GetSignalChannel(ctx, defs.RepoIOSignalQueueAddPriority.String())
		selector.AddReceive(priority, state.on_add_priority(ctx))

		remove := workflow.GetSignalChannel(ctx, defs.RepoIOSignalQueueRemove.String())
		selector.AddReceive(remove, state.on_remove(ctx))

		// setting up add_priority signal for adding prs to the priority queue
		add_priority := workflow.GetSignalChannel(ctx, defs.RepoIOSignalQueueAddPriority.String())
		selector.AddReceive(add_priority, state.on_add_priority(ctx))

		// setting up promote signal for moving a pr up in the queue
		promote := workflow.GetSignalChannel(ctx, defs.RepoIOSignalQueuePromote.String())
		selector.AddReceive(promote, state.on_promote(ctx))

		// setting up demote signal for moving a pr down in the queue
		demote := workflow.GetSignalChannel(ctx, defs.RepoIOSignalQueueDemote.String())
		selector.AddReceive(demote, state.on_demote(ctx))

		for state.is_active() {
			selector.Select(ctx)
		}
	})

	for state.is_active() {
		err := state.next(ctx)
		if err != nil {
			state.log(ctx, "next").Warn("context error")
			continue
		}

		pr := state.pop(ctx)
		if pr != nil {
			// process the pr
			err := state.process(ctx, pr)
			if err != nil {
				state.log(ctx, "process").Warn("processing error", "error", err)
				// Push the PR back into the queue
				state.push(ctx, pr, false) // Assuming it goes back to the primary queue
			}
		}

		// check if reset is needed
		if state.needs_reset() {
			queues := state.serialize(ctx)
			return state.as_new(ctx, "resetting due to event threshold", QueueCtrl, repo, branch, queues)
		}
	}

	return nil
}

// PrCtrl is a child worklow handle the pr activities and handle the pull request.
func PrCtrl(ctx workflow.Context, pr *defs.RepoIOPullRequest) error {
	return nil
}
