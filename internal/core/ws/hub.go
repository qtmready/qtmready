package ws

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"go.temporal.io/sdk/worker"

	"go.breu.io/quantm/internal/shared"
	"go.breu.io/quantm/internal/shared/queue"
)

type (
	// Hub interface defines the methods for managing WebSocket connections and messaging.
	Hub interface {
		// HandleWebSocket upgrades an HTTP connection to a WebSocket connection and manages it.
		// It creates a new Client, registers it with the hub, and starts read and write pumps.
		//
		// Example usage:
		//     e.GET("/ws/:id", ws.Instance().HandleWebSocket)
		HandleWebSocket(ctx echo.Context) error

		// Send sends a message to a specific user. It first checks if the user is local to the container.
		// If the user is found, the message is sent directly. If the user is not local, it checks which
		// queue the user is connected to. If the user is connected to a queue, the message is routed via
		// that queue. If the user is not connected to any queue, the message is dropped, and an
		// informational log is generated.
		// The method returns nil if the message is dropped or sent loally.
		// For all other errors, HubError is returned.
		//
		// Example usage:
		//     ctx := context.Background()
		//     err := hub.Send(ctx, "user123", []byte("Hello, user!"))
		//     if err != nil {
		//         log.Printf("Failed to send message: %v", err)
		//     }
		Send(ctx context.Context, user_id string, message []byte) error

		// Signal is a shorthand for signaling the ConnectionsHubWorkflow. It takes a signal type and
		// payload as parameters.
		//
		// Example usage:
		//     ctx := context.Background()
		//     err := hub.Signal(ctx, shared.WorkflowSignalStart, payload)
		//     if err != nil {
		//         log.Printf("Failed to send signal: %v", err)
		//     }
		Signal(ctx context.Context, signal shared.WorkflowSignal, payload any) error

		// Stop gracefully shuts down the hub and closes all client connections. It should be called
		// when the application is shutting down to ensure all resources are properly released and
		// all connections are closed.
		//
		// Example usage:
		//     hub.Stop()
		Stop()
	}

	// client represents a WebSocket client connection.
	client struct {
		user_id string
		conn    *websocket.Conn
		send    chan []byte
	}

	// hub manages WebSocket connections and message broadcasting.
	hub struct {
		clients    map[*client]bool
		register   chan *client
		unregister chan *client
		queue      queue.Queue
		mu         sync.RWMutex
		stop       chan bool
	}
)

var (
	instance *hub
	once     sync.Once
)

// HandleWebSocket upgrades an HTTP connection to a WebSocket connection and manages it.
func (h *hub) HandleWebSocket(ctx echo.Context) error {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // You may want to implement a more secure check
		},
	}

	conn, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return fmt.Errorf("failed to upgrade connection: %w", err)
	}

	user_id := ctx.Param("id")

	c := &client{
		user_id: user_id,
		conn:    conn,
		send:    make(chan []byte, 256),
	}

	go h.read(c)
	go h.write(c)

	return nil
}

// Send sends a message to a specific user.
func (h *hub) Send(ctx context.Context, user_id string, message []byte) error {
	if h.send_local(user_id, message) {
		return nil
	}

	name, err := h.query(ctx, user_id)
	if err != nil {
		var hubErr *HubError
		if errors.As(err, &hubErr) && hubErr.Code == ErrorCodeUserNotRegistered {
			shared.Logger().Warn("ws: user not registered", "user_id", user_id)
			return nil
		}

		return err
	}

	// Call send_global to send the message using the queue name
	err = h.route_message(ctx, queue.Name(name), user_id, message)
	if err != nil {
		return err
	}

	return nil
}

// Signal is a shorthand for signaling the ConnectionsHubWorkflow.
// It takes a signal type and payload as parameters.
func (h *hub) Signal(ctx context.Context, signal shared.WorkflowSignal, payload any) error {
	opts := opts_hub()
	_, err := shared.Temporal().
		Client().
		SignalWithStartWorkflow(
			ctx, opts.ID, signal.String(), payload, opts, ConnectionsHubWorkflow, NewConnections(),
		)

	return err
}

// Stop gracefully shuts down the hub and closes all client connections.
func (h *hub) Stop() {
	h.stop <- true
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		close(client.send)
		client.conn.Close()
	}

	// Signal to flush queues or perform any necessary cleanup
	if err := h.Signal(context.Background(), WorkflowSignalFlushQueue, RegisterOrFlush{Queue: h.queue.Name()}); err != nil {
		shared.Logger().Warn("ws: failed to signal flush", "error", err.Error())
	}

	close(h.register)
	close(h.unregister)
}

// send_local attempts to send a message to a client locally.
//
// It returns true if the message was sent successfully, false otherwise.
// If the client's send buffer is full or the client is disconnected,
// it removes the client from the hub.
func (h *hub) send_local(user_id string, message []byte) bool {
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

// route_message sends a message to a specific user using the provided queue name.
// It returns an error if the workflow execution fails.
//
// Example usage:
// err := h.route_message(ctx, queue.Name("userQueueName"), user_id, message)
//
//	if err != nil {
//	    // handle error
//	}
func (h *hub) route_message(ctx context.Context, q queue.Name, user_id string, message []byte) error {
	// Use the retrieved queue name to create workflow options
	opts := opts_send(queue.NewQueue(queue.WithName(q)), user_id)

	_, err := shared.Temporal().Client().ExecuteWorkflow(ctx, opts, SendMessageWorkflow, user_id, message)
	if err != nil {
		return NewHubError(ErrorCodeWorkflowExecutionFailed, "failed to send message", err)
	}

	return nil
}

// query queries the ConnectionsHandlerWorkflow to get the user's queue name.
// It returns the queue name and an error if the query fails.
//
// Example usage:
// name, err := h.query(ctx, user_id)
//
//	if err != nil {
//	    // handle error
//	}
func (h *hub) query(ctx context.Context, user_id string) (string, error) {
	response, err := shared.Temporal().Client().QueryWorkflow(ctx, opts_hub().ID, "", QueryGetUserQueue, user_id)
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

// worker sets up and runs the temporal worker for handling the hub's queue.
// On receiving the signal stop, it stops the worker.
func (h *hub) worker() {
	worker := h.queue.Worker(shared.Temporal().Client())

	// Register workflows
	worker.RegisterWorkflow(SendMessageWorkflow)
	worker.RegisterWorkflow(BroadcastMessageWorkflow)

	// Register activities
	worker.RegisterActivity(&Activities{})

	if err := worker.Start(); err != nil {
		shared.Logger().Error("Failed to start worker", "error", err)
		panic(err)
	}

	<-h.stop

	// Graceful shutdown
	worker.Stop()
}

// Helper functions

// client searches for a client with the given user_id.
func (h *hub) client(user_id string) (*client, bool) {
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
func (h *hub) read(client *client) {
	defer func() {
		h.unregister <- client
		client.conn.Close()

		// Signal that a user has disconnected
		if err := h.Signal(context.Background(), WorkflowSignalRemoveUser, User{UserID: client.user_id}); err != nil {
			shared.Logger().Warn("Failed to signal user disconnection", "user_id", client.user_id, "error", err)
			return
		}
	}()

	// Register the client
	h.register <- client

	// Signal that a user has connected
	if err := h.Signal(context.Background(), WorkflowSignalAddUser, QueueUser{UserID: client.user_id, Queue: h.queue.Name()}); err != nil {
		shared.Logger().Warn("Failed to signal user connection", "user_id", client.user_id, "error", err)
		return
	}

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				shared.Logger().Warn("Unexpected close error", "error", err)
			}

			break
		}

		// Handle incoming messages
		shared.Logger().Info("Received message from client", "userID", client.user_id, "message", string(message))
	}
}

// write writes messages to the WebSocket connection.
func (h *hub) write(client *client) {
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

// Additional functions

// ConnectionsHubWorker creates and configures a Temporal worker for handling WebSocket connections.
// It sets up a new queue with the WebSocketQueue name and registers the ConnectionsHandlerWorkflow.
//
// Example usage:
//
//	worker := ws.ConnectionsHubWorker()
//	err := worker.Start()
//	if err != nil {
//	    log.Fatalf("Failed to start worker: %v", err)
//	}
//	defer worker.Stop()
func ConnectionsHubWorker() worker.Worker {
	q := queue.NewQueue(queue.WithName(shared.WebSocketQueue))

	worker := q.Worker(shared.Temporal().Client())

	worker.RegisterWorkflow(ConnectionsHubWorkflow)

	return worker
}

// Instance returns a singleton instance of the Hub. It initializes the hub if it hasn't been created yet,
// setting up necessary channels, registering workflows and activities, and starting the worker and run loops.
// This method ensures that only one hub instance is created and used throughout the application.
//
// Example usage:
//
//	hub := ws.Instance()
//	// Use hub methods...
func Instance() Hub {
	once.Do(func() {
		instance = &hub{
			clients:    make(map[*client]bool),
			register:   make(chan *client),
			unregister: make(chan *client),
			queue:      queue.NewQueue(queue.WithName(queue_name())),
			stop:       make(chan bool, 1),
		}

		go instance.run()
		go instance.worker()
	})

	return instance
}
