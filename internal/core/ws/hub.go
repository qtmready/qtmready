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
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"go.breu.io/durex/queues"

	"go.breu.io/quantm/internal/shared"
)

type (
	// Hub interface defines the methods for managing WebSocket connections and messaging.
	Hub interface {
		// HandleWebSocket upgrades an HTTP connection to a WebSocket connection and manages it.
		// It creates a new Client, registers it with the hub, and starts read and write pumps.
		//
		// Example:
		//  e.GET("/ws/:id", ws.Instance().HandleWebSocket)
		HandleWebSocket(ctx echo.Context) error

		// Send sends a message to a specific user.
		//
		// If the user is local to the container, the message is sent directly.
		// Otherwise, the message is routed via the queue the user is connected to.
		// If the user is not connected to any queue, the message is dropped.
		//
		// Returns nil if the message is dropped or sent locally.
		// For all other errors, HubError is returned.
		//
		// Example:
		//  ctx := context.Background()
		//  err := hub.Send(ctx, "user123", []byte("Hello, user!"))
		//  if err != nil {
		//      log.Printf("Failed to send message: %v", err)
		//  }
		Send(ctx context.Context, user_id string, message []byte) error

		// Signal sends a signal to the ConnectionsHubWorkflow.
		//
		// Example:
		//  ctx := context.Background()
		//  err := hub.Signal(ctx, defs.WorkflowSignalStart, payload)
		//  if err != nil {
		//      log.Printf("Failed to send signal: %v", err)
		//  }
		Signal(ctx context.Context, signal queues.Signal, payload any) error

		// Stop gracefully shuts down the hub and closes all client connections.
		//
		// It should be called when the application is shutting down.
		//
		// Example:
		//  hub.Stop()
		Stop(ctx context.Context) error

		// SetAuthFn sets the authentication function for the hub.
		//
		// This function is used to configure the authentication process
		// for WebSocket connections.
		//
		// Example:
		//  verify = func (ctx context.Context, token string) (string, error) {
		//    parsed, err := parse(token)
		//    if err != nil {
		//      return "", err
		//    }
		//
		//    return parsed.User.ID, nil
		//  }
		//
		//  hub.SetAuthFn(verify)
		SetAuthFn(fn AuthFn)

		// SetQueryParam sets the name of query parameter to retrieve the token from the request. The default is "token".
		//
		// Example:
		//  hub.SetQueryParam("auth_token") // Token is in ?auth_token=...
		SetQueryParam(param string)
	}

	// connection holds information about a user's WebSocket connection.
	connection struct {
		user_id string
		conn    *websocket.Conn
		send    chan []byte
	}

	// hub manages WebSocket connections and message broadcasting.
	hub struct {
		clients    map[*connection]bool
		register   chan *connection
		unregister chan *connection
		queue      queues.Queue
		mu         sync.RWMutex
		stop       chan bool

		// auth handler
		auth  AuthFn // auth handler.
		param string // query param for auth e.g. for ?token=..., param is "token"
	}
)

var (
	instance *hub
	once     sync.Once
)

func (h *hub) HandleWebSocket(ctx echo.Context) error {
	token := ctx.QueryParam(h.param)

	user_id, err := h.auth(ctx.Request().Context(), token)
	if err != nil {
		return err
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return fmt.Errorf("failed to upgrade connection: %w", err)
	}

	c := &connection{
		user_id: user_id,
		conn:    conn,
		send:    make(chan []byte, 256),
	}

	go h.read(c)
	go h.write(c)

	return nil
}

func (h *hub) Send(ctx context.Context, user_id string, message []byte) error {
	if h.local(user_id, message) {
		return nil
	}

	q, err := h.query(ctx, user_id)
	if err != nil {
		var hubErr *HubError
		if errors.As(err, &hubErr) && hubErr.Code == ErrorCodeUserNotRegistered {
			slog.Warn("ws/hub: user not registered", "user_id", user_id)
			return nil
		}

		return err
	}

	err = h.route(ctx, q, user_id, message)
	if err != nil {
		return err
	}

	return nil
}

func (h *hub) Signal(ctx context.Context, signal queues.Signal, payload any) error {
	_, err := Queue().SignalWithStartWorkflow(ctx, opts_hub(), signal, payload, ConnectionsHubWorkflow, NewConnections())

	return err
}

func (h *hub) Stop(ctx context.Context) error {
	h.stop <- true
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		close(client.send)
		client.conn.Close()
	}

	if err := h.Signal(ctx, SingalContainerDisconnected, RegisterOrFlush{Queue: h.queue.String()}); err != nil {
		slog.Warn("ws/hub: failed to signal flush", "error", err.Error())
	}

	close(h.register)
	close(h.unregister)

	return nil
}

func (h *hub) SetAuthFn(fn AuthFn) {
	slog.Info("ws/hub: setting auth handler ...")

	h.auth = fn
}

func (h *hub) SetQueryParam(param string) {
	h.param = param
}

// local attempts to send a message to a client locally.
//
// It returns true if the message was sent successfully, false otherwise. If the client's send buffer is full or the
// client is disconnected, it removes the client from the hub.
func (h *hub) local(user_id string, message []byte) bool {
	client, found := h.client(user_id)
	if found {
		select {
		case client.send <- message:
			return true

		default:
			h.mu.Lock()
			defer h.mu.Unlock()

			if _, connected := h.clients[client]; connected {
				delete(h.clients, client)
				close(client.send)
			}

			return false
		}
	}

	return false
}

// route sends a message to a specific user using the provided queue name.
// It returns an error if the workflow execution fails.
//
// Example usage: err := h.route(ctx, queue.Name("userQueueName"), user_id, message)
//
//	if err != nil {
//	    // handle error
//	}
func (h *hub) route(ctx context.Context, _queue, user_id string, message []byte) error {
	q := queues.New(queues.WithName(_queue), queues.WithClient(shared.Temporal().Client()))

	_, err := q.ExecuteWorkflow(ctx, opts_send(user_id), SendMessageWorkflow, user_id, message)
	if err != nil {
		return NewHubError(ErrorCodeWorkflowExecutionFailed, "failed to send message", err)
	}

	return nil
}

// query queries the ConnectionsHandlerWorkflow to get the user's queue name. It returns the queue name and an error if
// the query fails.
//
// Example usage: name, err := h.query(ctx, user_id)
//
//	if err != nil {
//	    // handle error
//	}
func (h *hub) query(ctx context.Context, user_id string) (string, error) {
	response, err := Queue().QueryWorkflow(ctx, opts_hub(), QueryGetUserQueue, user_id)
	if err != nil {
		return "", NewHubError(ErrorCodeQueryFailed, "failed to query user queue", err)
	}

	var name string
	if err := response.Get(&name); err != nil {
		return "", NewHubError(ErrorCodeQueryFailed, "failed to decode user queue response", err)
	}

	if name == "" {
		return "", NewHubError(ErrorCodeUserNotRegistered, "user not registered to any queue", nil)
	}

	return name, nil
}

// run is the main loop that handles client registration, unregistration, and broadcasting.
func (h *hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			defer h.mu.Unlock()

			h.clients[client] = true

		case client := <-h.unregister:
			h.mu.Lock()
			defer h.mu.Unlock()

			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case <-h.stop:
			return
		}
	}
}

// worker sets up and runs the Temporal worker for handling the hub's queue. It starts the worker, registers workflows
// and activities, and listens for the stop signal. On receiving the stop signal, it gracefully shuts down the worker.
func (h *hub) worker() {
	slog.Info("ws/queue: starting worker ...")

	h.queue.CreateWorker()

	h.queue.RegisterWorkflow(SendMessageWorkflow)
	h.queue.RegisterWorkflow(BroadcastMessageWorkflow)

	h.queue.RegisterActivity(&Activities{})

	if err := h.queue.Start(); err != nil {
		h.stop <- true
	}

	err := h.Signal(context.Background(), SignalContainerConnected, &RegisterOrFlush{Queue: h.queue.String()})
	if err != nil {
		slog.Warn("ws/queue: failed to signal worker addition", "error", err)
		panic(err)
	}

	<-h.stop

	_ = h.queue.Shutdown(context.Background())
}

// client searches for a client with the given user_id.
//
// It returns the connection and a boolean indicating whether the client was found.
func (h *hub) client(user_id string) (*connection, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		if client.user_id == user_id {
			return client, true
		}
	}

	return nil, false
}

// read reads messages from the WebSocket connection.
func (h *hub) read(client *connection) {
	defer func() {
		h.unregister <- client
		client.conn.Close()

		// Signal that a user has disconnected
		if err := h.Signal(context.Background(), SignalUserDisconnected, User{UserID: client.user_id}); err != nil {
			slog.Warn("ws/hub: failed to remove client", "user_id", client.user_id, "error", err)
			return
		}
	}()

	// Register the client
	h.register <- client

	// Signal that a user has connected
	if err := h.Signal(context.Background(), SignalUserConnected, User{UserID: client.user_id}); err != nil {
		slog.Warn("ws/hub: failed to add client", "user_id", client.user_id, "error", err)
		return
	}

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Warn("ws/hub: websocket read error", "error", err)
			}

			break
		}

		// FIXME: Implement message handling.
		slog.Debug("ws: recieved on websocket", "user_id", client.user_id, "message", string(message))
	}
}

// write writes messages to the WebSocket connection.
func (h *hub) write(client *connection) {
	defer func() {
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			if _, err := w.Write(message); err != nil {
				return
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

// Instance returns a singleton instance of the Hub. It initializes the hub if it hasn't been created yet, setting up
// necessary channels, registering workflows and activities, and starting the worker and run loops. This method ensures
// that only one hub instance is created and used throughout the application.
//
// Example usage:
//
//	hub := ws.Instance()
//	// Use hub methods...
func Instance() Hub {
	once.Do(func() {
		instance = &hub{
			// websocket connections
			clients:    make(map[*connection]bool),
			register:   make(chan *connection),
			unregister: make(chan *connection),
			// temporal: to provide better visibility, we prefix the container id with websockets queue, so that our id for
			// tasks on the queue worker will look like "io.ctrlplane.{websockets_queue_name}.{container_id}""
			queue: queues.New(
				queues.WithName(queue_name(container_id())),
				queues.WithClient(shared.Temporal().Client()),
			),
			stop: make(chan bool, 1),
			// authentication
			auth:  noop,
			param: "token",
		}

		go instance.run()
		go instance.worker()
	})

	return instance
}
