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
	"crypto/rand"
	"fmt"

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

// opts_hub returns StartWorkflowOptions for the ConnectionsHubWorkflow.
func opts_hub() sdk_client.StartWorkflowOptions {
	return shared.Temporal().Queue(shared.WebSocketQueue).WorkflowOptions(
		queue.WithWorkflowBlock("hub"),
	)
}

// encode converts a byte slice into a 24-character string using base-32 encoding.
//
// Each byte is encoded as three characters, using a 32-character set (lowercase a-z, digits 2-3, 4-5, 6-8, 9). The
// algorithm converts each byte to an integer, calculates its modulo 32, and uses the corresponding character from the
// character set. This process is repeated three times, dividing the integer by 32 each time.
func encode(bytes []byte) string {
	chars := "abcdefghijklmnopqrstuvwxyz2345689"
	encoded := ""

	for i := 0; i < len(bytes); i++ {
		value := int(bytes[i])

		encoded += string(chars[value%32])
		value /= 32
		encoded += string(chars[value%32])
		value /= 32
		encoded += string(chars[value%32])
	}

	return encoded
}

// container_id generates a cryptographically secure, collision-resistant container ID.
//
// It uses base-32 encoding of 8 random bytes to produce a 24-character string. The first 8 characters are used.
//
// This provides a namespace of 1 billion unique IDs (32^8), sufficient for current scale. However, the birthday paradox
// dictates that collision probability increases as the number of containers grows.
//
// For scalability beyond a few hundred containers at a given time, a more complex algorithm or a larger namespace may
// be necessary to
// minimize collision risk.
func container_id() queue.Name {
	length := 8
	bytes := make([]byte, 8)

	_, _ = rand.Read(bytes)

	encoded := encode(bytes)

	return queue.Name(fmt.Sprintf("%s.%s", shared.WebSocketQueue.String(), encoded[:length]))
}
