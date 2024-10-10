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

	"go.breu.io/durex/workflows"
)

// opts_send constructs workflow options for sending a message to a user.
//
// It creates a block named "send" with two elements: "user" and "message". The user element uses the provided user_id.
// The message element uses a unique container ID.
//
// For example, for a user with ID "123e4567-e89b-12d3-a456-426614174000", and coupled with Queue(), the resulting
// workflow ID would be:
//
//	"io.ctrlplane.websocket.send.user.123e4567-e89b-12d3-a456-426614174000.message.89ab2c3d"
func opts_send(user_id string) workflows.Options {
	opts, _ := workflows.NewOptions(
		workflows.WithBlock("send"),
		workflows.WithElement("user"),
		workflows.WithElementID(user_id),
		workflows.WithElement("message"),
		workflows.WithElementID(idempotent()),
	)

	return opts
}

// opts_broadcast constructs workflow options for broadcasting a message to a team.
//
// It creates a block named "broadcast" with two elements: "team" and "message". The team element uses the provided
// team_id. The message element uses a unique container ID.
//
// For example, for a team with ID "456789ab-cdef-1234-5678-901234567890", the resulting workflow ID, when coupled with
// Queue(), would be:
//
//	"io.ctrlplane.websocket.broadcast.team.456789ab-cdef-1234-5678-901234567890.message.fay7cxkg"
func opts_broadcast(team_id string) workflows.Options {
	opts, _ := workflows.NewOptions(
		workflows.WithBlock("broadcast"),
		workflows.WithElement("team"),
		workflows.WithElementID(team_id),
		workflows.WithElement("message"),
		workflows.WithElementID(idempotent()),
	)

	return opts
}

// opts_hub constructs workflow options for the hub block.
//
// It creates a block named "hub".
func opts_hub() workflows.Options {
	opts, _ := workflows.NewOptions(
		workflows.WithBlock("hub"),
	)

	return opts
}

// encode converts a byte slice into a base-33 encoded string.
//
// Each byte is encoded as three characters using the character set, (lowercase a-z, digits 2-3, 4-5, 6-8, 9).
// "abcdefghijklmnopqrstuvwxyz2345689". The encoding process converts each byte to an integer, calculates its modulo 33,
// and uses the corresponding character from the character set. This is repeated three times, dividing the integer by 33
// each time. The resulting encoded string has a length of 3 * len(bytes).
//
// Example:
//
//	encode([]byte("hello")) == "2c3h6w82b2b3h6w82b2b3h6w82b"  // length: 5 x 3 = 15
func encode(bytes []byte) string {
	chars := "abcdefghijklmnopqrstuvwxyz2345689"
	encoded := ""

	for i := 0; i < len(bytes); i++ {
		value := int(bytes[i])

		encoded += string(chars[value%33])
		value /= 33
		encoded += string(chars[value%33])
		value /= 33
		encoded += string(chars[value%33])
	}

	return encoded
}

// idempotent generates a unique 8-character ID using base-33 encoding of 8 random bytes, providing a
// namespace of approximately 14.3 trillion unique IDs (33^8). These IDs are short-lived and require unique generation
// within their lifespan. While collisions are unlikely given the large namespace, we need to avoid them to ensure
// proper functionality.
//
// The implementation anticipates sufficient capacity for current needs. However, scalability beyond a few hundred
// containers may necessitate a more robust solution to minimize the risk of collisions. This is because the birthday
// paradox dictates that the probability of collisions increases with the number of IDs generated, even if not all are
// in use simultaneously.
//
// Future improvements could include increasing the namespace size by using a larger character set or encoding more
// bytes. Alternatively, exploring sophisticated collision-resistant algorithms could further minimize the probability
// of collisions, even with a large number of generated IDs.
func idempotent() string {
	length := 8
	bytes := make([]byte, 8)

	_, _ = rand.Read(bytes)

	encoded := encode(bytes)

	return encoded[:length]
}

// queue_name constructs a queue name by concatenating the queue name and the provided ID.
//
// The queue name is retrieved using Queue().String(), whereas Queue() is a package local singleton.
func queue_name(id string) string {
	return fmt.Sprintf("%s.%s", Queue().String(), id)
}
