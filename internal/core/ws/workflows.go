package ws

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/shared"
)

const (
	WorkflowSignalAddUser     shared.WorkflowSignal = "orchestrate__add_user"
	WorkflowSignalRemoveUser  shared.WorkflowSignal = "orchestrate__remove_user"
	WorkflowSignalFlushQueue  shared.WorkflowSignal = "orchestrate__flush_queue"
	WorkflowSignalWorkerAdded shared.WorkflowSignal = "orchestrate__worker_added"

	QueryGetUserQueue = "get_user_queue"
)

// ConnectionsHandlerWorkflow manages the state of WebSocket connections across multiple API containers in a Kubernetes cluster.
// It doesn't handle the actual WebSocket connections, but rather maintains a record of which users are connected to which containers.
//
// This workflow solves the challenge of routing WebSocket messages in a scalable API setup where:
// - The number of API containers can scale from a minimum of 3 to n.
// - Users can be connected to any of the available containers.
// - Events may be received by different containers than the one a user is connected to.
//
// The workflow:
// - Maintains a map of user connections and their associated API containers.
// - Handles signals for adding/removing users, & flushing queues if the container goes down.
// - Provides a query handler for retrieving queue which user is connected to, which is used by the Hub for message routing.
//
// Parameters:
//   - ctx: The workflow context
//   - conns: A pointer to the Connections struct that manages the connection state across containers
//
// Returns:
//   - error: Any error that occurs during the workflow execution
func ConnectionsHandlerWorkflow(ctx workflow.Context, conns *Connections) error {
	conns.Restore(ctx)
	selector := workflow.NewSelector(ctx)
	info := workflow.GetInfo(ctx)

	// Set up signal channels
	add := workflow.GetSignalChannel(ctx, WorkflowSignalAddUser.String())
	remove := workflow.GetSignalChannel(ctx, WorkflowSignalRemoveUser.String())
	flush := workflow.GetSignalChannel(ctx, WorkflowSignalFlushQueue.String())
	worker_added := workflow.GetSignalChannel(ctx, WorkflowSignalWorkerAdded.String())

	// Add signal handlers to the selector
	selector.AddReceive(add, conns.on_add(ctx))
	selector.AddReceive(remove, conns.on_remove(ctx))
	selector.AddReceive(flush, conns.on_flush(ctx))
	selector.AddReceive(worker_added, conns.on_worker_added(ctx))

	// Set up query handler for getting user queue
	_ = workflow.SetQueryHandler(ctx, QueryGetUserQueue, func(user_id string) string {
		q, ok := conns.GetQueueForUser(ctx, user_id)

		if ok {
			return q
		}

		return ""
	})

	// Main loop: continuously handle signals
	for {
		selector.Select(ctx)

		// Check if a new workflow instance should be created
		if info.GetContinueAsNewSuggested() {
			return workflow.NewContinueAsNewError(ctx, ConnectionsHandlerWorkflow, conns)
		}
	}
}

func SendMessageWorkflow(ctx workflow.Context, user_id string, message []byte) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	activities := &Activities{}

	err := workflow.ExecuteActivity(ctx, activities.SendMessage, user_id, message).Get(ctx, nil)
	if err != nil {
		shared.Logger().Error("Failed to send message", "error", err)
		return err
	}

	return nil
}

func BroadcastMessageWorkflow(ctx workflow.Context, team_id string, message []byte) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	activities := &Activities{}

	var userIDs []string

	err := workflow.ExecuteActivity(ctx, activities.GetTeamUsers, team_id).Get(ctx, &userIDs)
	if err != nil {
		shared.Logger().Error("Failed to get team users", "error", err)
		return err
	}

	for _, userID := range userIDs {
		err := workflow.ExecuteActivity(ctx, activities.SendMessage, userID, message).Get(ctx, nil)
		if err != nil {
			shared.Logger().Error("Failed to send message to user", "user_id", userID, "error", err)
		}
	}

	return nil
}
