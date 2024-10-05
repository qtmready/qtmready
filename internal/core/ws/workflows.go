// Crafted with ❤ at Breu, Inc. <info@breu.io>, Copyright © 2024.
//
// Functional Source License, Version 1.1, Apache 2.0 Future License
//
// We hereby irrevocably grant you an additional license to use the Software under the Apache License, Version 2.0 that
// is effective on the second anniversary of the date we make the Software available. On or after that date, you may use
// the Software under the Apache License, Version 2.0, in which case the following will apply:
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
// the License.
//
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package ws

import (
	"go.breu.io/durex/dispatch"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/defs"
)

const (
	WorkflowSignalAddUser     defs.Signal = "orchestrate__add_user"
	WorkflowSignalRemoveUser  defs.Signal = "orchestrate__remove_user"
	WorkflowSignalFlushQueue  defs.Signal = "orchestrate__flush_queue"
	WorkflowSignalWorkerAdded defs.Signal = "orchestrate__worker_added"

	QueryGetUserQueue = "get_user_queue"
)

// ConnectionsHubWorkflow manages the state of WebSocket connections across multiple API containers in a Kubernetes
// cluster. It doesn't handle the actual WebSocket connections, but rather maintains a record of which users are
// connected to which containers.
//
// This workflow solves the challenge of routing WebSocket messages in a scalable API setup where: - The number of API
// containers can scale from a minimum of 3 to n. - Users can be connected to any of the available containers. - Events
// may be received by different containers than the one a user is connected to.
//
// The workflow: - Maintains a map of user connections and their associated API containers. - Handles signals for
// adding/removing users, & flushing queues if the container goes down. - Provides a query handler for retrieving queue
// which user is connected to, which is used by the Hub for message routing.
//
// Parameters:
//   - ctx: The workflow context
//   - conns: A pointer to the Connections struct that manages the connection state across containers
//
// Returns:
//   - error: Any error that occurs during the workflow execution
func ConnectionsHubWorkflow(ctx workflow.Context, conns *Connections) error {
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
			return workflow.NewContinueAsNewError(ctx, ConnectionsHubWorkflow, conns)
		}
	}
}

// SendMessageWorkflow is a Temporal workflow that routes messages to connected users.
//
// When the Hub identifies the container where a user is connected, it triggers this workflow. The workflow leverages
// the `SendMessage` activity and the Hub singleton to send the message directly to the intended user. If the message
// delivery fails, likely due to the user disconnecting during the routing process, a warning is logged.
func SendMessageWorkflow(ctx workflow.Context, user_id string, message []byte) error {
	logger := workflow.GetLogger(ctx)
	activities := &Activities{}
	sent := false

	ctx = dispatch.WithDefaultActivityContext(ctx)

	err := workflow.ExecuteActivity(ctx, activities.SendMessage, user_id, message).Get(ctx, &sent)
	if err != nil {
		logger.Error("ws/send: unable to execute activity ..", "error", err)
		return err
	}

	if !sent {
		logger.Warn("ws/send: unable to send locally, dropping ..", "user_id", user_id)
	}

	return nil
}

func BroadcastMessageWorkflow(ctx workflow.Context, team_id string, message []byte) error {
	logger := workflow.GetLogger(ctx)

	ctx = dispatch.WithDefaultActivityContext(ctx)

	activities := &Activities{}
	response := &TeamUsersReponse{IDs: make([]string, 0)}

	err := workflow.ExecuteActivity(ctx, activities.GetTeamUsers, team_id).Get(ctx, response)
	if err != nil {
		logger.Error("ws/broadcast: unable to fetch users ...", "error", err)
		return err
	}

	for _, id := range response.IDs {
		err := workflow.ExecuteActivity(ctx, activities.RouteMessage, id, message).Get(ctx, nil)
		if err != nil {
			logger.Error("ws/broadcast: unable to route message ...", "user_id", id, "error", err)
		}
	}

	return nil
}
