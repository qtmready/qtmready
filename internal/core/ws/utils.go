// Copyright © 2024, Breu, Inc. <info@breu.io>
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
	"crypto/rand"
	"encoding/base32"

	"github.com/google/uuid"
	sdk_client "go.temporal.io/sdk/client"

	"go.breu.io/quantm/internal/shared"
	"go.breu.io/quantm/internal/shared/queue"
)

// idempotent creates an idempotent ID for a workflow element.
func idempotent() string {
	return uuid.NewString()
}

// opts_send returns StartWorkflowOptions for sending a message to a specific user.
func opts_send(q queue.Queue, user_id string) sdk_client.StartWorkflowOptions {
	return q.WorkflowOptions(
		queue.WithWorkflowBlock("user"),
		queue.WithWorkflowBlockID(user_id),
		queue.WithWorkflowElement("message"),
		queue.WithWorkflowElementID(idempotent()),
	)
}

// opts_broadcast returns StartWorkflowOptions for broadcasting a message to a team.
func opts_broadcast(q queue.Queue, team_id string) sdk_client.StartWorkflowOptions {
	return q.WorkflowOptions(
		queue.WithWorkflowBlock("team"),
		queue.WithWorkflowBlockID(team_id),
		queue.WithWorkflowElement("message"),
		queue.WithWorkflowElementID(idempotent()),
	)
}

// opts_hub returns StartWorkflowOptions for the WebSocket hub workflow.
func opts_hub() sdk_client.StartWorkflowOptions {
	return shared.Temporal().Queue(shared.WebSocketQueue).WorkflowOptions(
		queue.WithWorkflowBlock("hub"),
	)
}

// queue_name create a name for the temporal queue.
// It is used to create a unique name for the queue for each running container.
func queue_name() queue.Name {
	length := 8
	bytes := make([]byte, 5) // 5 bytes will give us at least 8 characters when base32 encoded

	_, _ = rand.Read(bytes)

	// Use base32 encoding to ensure we only get lowercase letters and numbers
	encoded := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes)

	// Trim to exactly 8 characters
	return queue.Name(encoded[:length])
}
