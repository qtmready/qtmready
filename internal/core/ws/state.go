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
	"fmt"

	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"

	"go.breu.io/quantm/internal/core/defs"
)

type (
	// Connections represents the state of websocket connections, managing the relationship between users and containers.
	Connections struct {
		UserContainer        map[string]string              `json:"user_container"`
		ContainerConnections map[string]map[string]struct{} `json:"container_connections"`
		mu                   workflow.Mutex
		logger               log.Logger
	}

	// UserContainer represents a signal containing user ID and container name.
	UserContainer struct {
		UserID    string
		Container string
	}

	// User represents a signal containing user ID.
	User struct {
		UserID string
	}

	// RegisterOrFlush represents a signal containing a container name.
	RegisterOrFlush struct {
		Container string
	}
)

// -- Channel Handlers --

// on_add is a channel handler triggered when a new user connection is registered.
func (con *Connections) on_add(ctx workflow.Context) defs.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		var signal UserContainer

		rx.Receive(ctx, &signal)

		if err := con.ConnectUser(ctx, signal.UserID, signal.Container); err != nil {
			con.error("connection registration failed", "user_id", signal.UserID, "container", signal.Container, "error", err)
		} else {
			con.info("connection registered", "user_id", signal.UserID, "container", signal.Container)
		}
	}
}

// on_remove is a channel handler triggered when a user connection is dropped.
func (con *Connections) on_remove(ctx workflow.Context) defs.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		var signal User

		rx.Receive(ctx, &signal)

		if err := con.DropUserConnection(ctx, signal.UserID); err != nil {
			con.error("failed to drop connection", "user_id", signal.UserID, "error", err)
		} else {
			con.info("connection dropped", "user_id", signal.UserID)
		}
	}
}

// on_drop is a channel handler triggered when a container is disconnected.
func (con *Connections) on_drop(ctx workflow.Context) defs.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		var signal RegisterOrFlush

		rx.Receive(ctx, &signal)

		if err := con.DropContainer(ctx, signal.Container); err != nil {
			con.error("failed to drop container", "container", signal.Container, "error", err)
		} else {
			con.info("container disconnected", "container", signal.Container)
		}
	}
}

// on_container_connected is a channel handler triggered when a container is connected.
func (con *Connections) on_container_connected(ctx workflow.Context) defs.ChannelHandler {
	return func(rx workflow.ReceiveChannel, more bool) {
		var signal RegisterOrFlush

		rx.Receive(ctx, &signal)

		con.info("container connected", "container", signal.Container)
	}
}

// -- Connection Management --

// GetContainerForUser returns the container name for a given user ID. It returns `false` if the user is not connected
// to any container.
func (con *Connections) GetContainerForUser(ctx workflow.Context, user_id string) (string, bool) {
	con.info("get container for user", "user_id", user_id)

	if err := con.mu.Lock(ctx); err != nil {
		return "", false
	}
	defer con.mu.Unlock()

	container, exists := con.UserContainer[user_id]

	if exists {
		con.info("container found", "user_id", user_id, "container", container)
	} else {
		con.warn("no container found for user", "user_id", user_id)
	}

	return container, exists
}

// ConnectUser registers a user to a specific container. If the user is already connected to a different container, the
// connection is updated.
func (con *Connections) ConnectUser(ctx workflow.Context, user_id, container string) error {
	if err := con.mu.Lock(ctx); err != nil {
		return err
	}
	defer con.mu.Unlock()

	if old, exists := con.UserContainer[user_id]; exists {
		delete(con.ContainerConnections[old], user_id)

		if len(con.ContainerConnections[old]) == 0 {
			delete(con.ContainerConnections, old)
		}
	}

	con.UserContainer[user_id] = container
	if _, exists := con.ContainerConnections[container]; !exists {
		con.ContainerConnections[container] = make(map[string]struct{})
	}

	con.ContainerConnections[container][user_id] = struct{}{}

	return nil
}

// DropUserConnection removes a user from the connection list. If the user is not connected, no action is taken.
func (con *Connections) DropUserConnection(ctx workflow.Context, user_id string) error {
	if err := con.mu.Lock(ctx); err != nil {
		return err
	}
	defer con.mu.Unlock()

	if container, exists := con.UserContainer[user_id]; exists {
		delete(con.UserContainer, user_id)
		delete(con.ContainerConnections[container], user_id)

		if len(con.ContainerConnections[container]) == 0 {
			delete(con.ContainerConnections, container)
		}
	}

	return nil
}

// GetContainerUser returns a list of user IDs connected to a specific container.
func (con *Connections) GetContainerUser(ctx workflow.Context, container string) ([]string, error) {
	if err := con.mu.Lock(ctx); err != nil {
		return nil, err
	}
	defer con.mu.Unlock()

	users := make([]string, 0, len(con.ContainerConnections[container]))
	for user_id := range con.ContainerConnections[container] {
		users = append(users, user_id)
	}

	return users, nil
}

// DropContainer removes all users connected to a container and the container itself from the connection list.
func (con *Connections) DropContainer(ctx workflow.Context, container string) error {
	if err := con.mu.Lock(ctx); err != nil {
		return err
	}
	defer con.mu.Unlock()

	if users, exists := con.ContainerConnections[container]; exists {
		for user_id := range users {
			delete(con.UserContainer, user_id)
		}

		delete(con.ContainerConnections, container)
	}

	return nil
}

// -- Utility Functions --

// Restore reinitializes the mutex and logger. This should be called when deserializing the Connections struct.
func (con *Connections) Restore(ctx workflow.Context) {
	con.mu = workflow.NewMutex(ctx)
	con.logger = workflow.GetLogger(ctx)
}

func (con *Connections) prefixed(msg string) string {
	return fmt.Sprintf("ws: %s", msg)
}

func (con *Connections) info(msg string, keyvals ...any) {
	con.logger.Info(con.prefixed(msg), keyvals...)
}

func (con *Connections) debug(msg string, keyvals ...any) {
	con.logger.Debug(con.prefixed(msg), keyvals...)
}

func (con *Connections) warn(msg string, keyvals ...any) {
	con.logger.Warn(con.prefixed(msg), keyvals...)
}

func (con *Connections) error(msg string, keyvals ...any) {
	con.logger.Error(con.prefixed(msg), keyvals...)
}

// -- Construction --

// NewConnections creates a new Connections instance.
func NewConnections() *Connections {
	return &Connections{
		UserContainer:        make(map[string]string),
		ContainerConnections: make(map[string]map[string]struct{}),
	}
}
