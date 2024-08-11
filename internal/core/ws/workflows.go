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

func ConnectionsHandlerWorkflow(ctx workflow.Context, conns *Connections) error {
	conns.Restore(ctx)
	selector := workflow.NewSelector(ctx)
	info := workflow.GetInfo(ctx)

	add := workflow.GetSignalChannel(ctx, WorkflowSignalAddUser.String())
	remove := workflow.GetSignalChannel(ctx, WorkflowSignalRemoveUser.String())
	flush := workflow.GetSignalChannel(ctx, WorkflowSignalFlushQueue.String())
	worker_added := workflow.GetSignalChannel(ctx, WorkflowSignalWorkerAdded.String())

	selector.AddReceive(add, conns.on_add(ctx))
	selector.AddReceive(remove, conns.on_remove(ctx))
	selector.AddReceive(flush, conns.on_flush(ctx))
	selector.AddReceive(worker_added, conns.on_worker_added(ctx))

	_ = workflow.SetQueryHandler(ctx, QueryGetUserQueue, func(user_id string) string {
		q, ok := conns.GetQueueForUser(ctx, user_id)

		if ok {
			return q
		}

		return ""
	})

	for {
		selector.Select(ctx)

		if info.GetContinueAsNewSuggested() {
			if err := workflow.NewContinueAsNewError(ctx, ConnectionsHandlerWorkflow, conns); err != nil {
				return err
			}

			break
		}
	}

	return nil
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
