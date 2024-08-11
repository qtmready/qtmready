package ws

import (
	"github.com/google/uuid"
	sdk_client "go.temporal.io/sdk/client"

	"go.breu.io/quantm/internal/shared"
	"go.breu.io/quantm/internal/shared/queue"
)

// _id creates an idempotent ID for a workflow element.
func _id() string {
	return uuid.NewString()
}

// opts_send returns StartWorkflowOptions for sending a message to a specific user.
func opts_send(q queue.Queue, user_id string) sdk_client.StartWorkflowOptions {
	return q.WorkflowOptions(
		queue.WithWorkflowBlock("user"),
		queue.WithWorkflowBlockID(user_id),
		queue.WithWorkflowElement("message"),
		queue.WithWorkflowElementID(_id()),
	)
}

// opts_broadcast returns StartWorkflowOptions for broadcasting a message to a team.
func opts_broadcast(q queue.Queue, team_id string) sdk_client.StartWorkflowOptions {
	return q.WorkflowOptions(
		queue.WithWorkflowBlock("team"),
		queue.WithWorkflowBlockID(team_id),
		queue.WithWorkflowElement("message"),
		queue.WithWorkflowElementID(_id()),
	)
}

// opts_hub returns StartWorkflowOptions for the WebSocket hub workflow.
func opts_hub() sdk_client.StartWorkflowOptions {
	return shared.Temporal().Queue(shared.WebSocketQueue).WorkflowOptions(
		queue.WithWorkflowBlock("hub"),
	)
}
