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

// Package ws provides a robust and scalable WebSocket connection management and messaging system built on Temporal.
//
// # Challenge
//
// Our application, running on Cloud Run, needed to provide real-time updates to the frontend for features like
// visibility into merge queues and near real-time logs for build and orchestration tasks. This required reliable
// message delivery to specific users.  We faced two key challenges:
//
//   - Events of interest to a user could occur on a container different from the one where the user is connected.
//     This necessitated a solution for reliably routing messages between containers, even when the event-generating
//     container differed from the one hosting the user's connection.
//   - WebSockets' stateless nature demanded a mechanism to determine the container a user is connected to.
//
// # Design:
//
// Our application's existing reliance on Temporal makes it a natural choice for building our WebSocket management
// system. This approach minimizes new components, reduces complexity, and leverages Temporal's inherent durability for
// scalability.
//
// The design hinges on temporal isolation of workers, each listening to a single queue. By leveraging this isolation,
// we assign a unique queue to each container on startup. This queue name is generated using a UUID at startup and acts
// as the container's unique identifier, ensuring messages are delivered to the correct container. After the container
// is assigned a queue, we need to make this information available to
//
//   - the container itself.
//   - the entire applicaiton, which is a collection of containers.
//
// To achieve this, we employ a two-pronged approach: a singleton called "Hub" for local coordination and a
// long-running, always-available workflow called "ConnectionsHubWorkflow" for persistent connection management and
// message routing.
//
// The ConnectionsHubWorkflow acts as the central state manager for all WebSocket connections across all containers in
// the application. It:
//
//   - Maintains a map of users and their connected container.
//   - Provides signals for adding or removing connection information for a user.
//   - Provides a query handler for retrieving the queue a user is connected to.
//   - Flushes queues on container shutdown.
//
// The Hub manages local WebSocket connections, handling connection establishment, disconnection, and message sending.
// It interacts with the ConnectionsHubWorkflow, using signals to add or remove users from the global state and queries
// to retrieve a user's connected queue.
//
// When sending a message, the Hub first checks if the user is connected to the current container. If so, it uses the
// Hub's WebSocket handler to send the message directly. Otherwise, it queries the ConnectionsHubWorkflow to find the
// user's connected queue. If a queue is found, the message is routed through Temporal using the SendMessageWorkflow
// with the retrieved queue information. If the user is not found, the message is dropped.  This entire process is
// hidden from the developer, who simply calls the Send function on the Hub.
//
// # Limitations:
//
//   - The current design does not handle users having multiple connections.
//   - The websocket connection handler is tied to Echo Framework. This can be abstracted out to support other
//     frameworks.
//
// # Open Questions:
//
// We have only one ConnectionsHubWorkflow running in the system. We are not worried about the unavialability of this
// because Temporal guarantees workflow completion. We have taken care of workflow event history by taking care of
// the restart of the workflow. It works for now, but what happens when we scale? Can ConnectionsHubWorkflow keep up
// with handling several hundred queries and signals per second?  We need to test this.
//
// # Future Work:
//
//   - Implement broadcast functionality to send messages to multiple users or teams.
//   - Implement a mechanism to handle users connected to multiple containers.
//   - Implement robust monitoring and alerting for the WebSocket management system, tracking key metrics like
//     connection counts, message throughput, and latency.
//
// # Example:
//
//	import (
//	  "context"
//	  "fmt"
//	  "github.com/labstack/echo/v4"
//	  "go.breu.io/quantm/internal/ws"
//	)
//
//	func main() {
//	  e := echo.New()
//
//	  // Set up a WebSocket endpoint
//	  e.GET("/ws/:token", ws.Instance().HandleWebSocket)
//
//	  // Send a message to a user
//	  ctx := context.Background()
//	  err := ws.Instance().Send(ctx, "user123", []byte("Hello, user!"))
//	  if err != nil {
//	    fmt.Printf("Failed to send message: %v", err)
//	  }
//	}
package ws
