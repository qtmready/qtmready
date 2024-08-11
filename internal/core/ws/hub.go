package ws

import (
	"context"
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
		HandleWebSocket(ctx echo.Context) error
		Send(ctx context.Context, user_id string, message []byte) error
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

// Hub interface methods

// HandleWebSocket upgrades an HTTP connection to a WebSocket connection and manages it.
// It creates a new Client, registers it with the hub, and starts read and write pumps
// to handle incoming and outgoing messages.
//
// Example:
//
//	e.GET("/ws/:id", ws.Instance().HandleWebSocket)
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

	h.register <- c

	go h.read(c)
	go h.write(c)

	return nil
}

// send_local attempts to send a message to a client locally.
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

// Send sends a message to a specific user.
// If the user is connected, it attempts to send the message directly.
// If the user is not connected or the send fails, it starts a Temporal workflow
// to handle the message delivery.
//
// Example:
//
//	ctx := context.Background()
//	err := hub.Send(ctx, "user123", []byte("Hello, user!"))
//	if err != nil {
//	    log.Printf("Failed to send message: %v", err)
//	}
func (h *hub) Send(ctx context.Context, user_id string, message []byte) error {
	if h.send_local(user_id, message) {
		return nil
	}

	// Query the ConnectionsHandlerWorkflow to get the user's queue
	response, err := shared.Temporal().Client().QueryWorkflow(ctx, opts_hub().ID, "", QueryGetUserQueue, user_id)
	if err != nil {
		return NewHubError(ErrorTypeQueryFailed, "failed to query user queue", err)
	}

	var name string
	if err := response.Get(&name); err != nil {
		return NewHubError(ErrorTypeQueryFailed, "failed to decode user queue response", err)
	}

	if name == "" {
		return NewHubError(ErrorTypeUserNotRegistered, "user not registered to any queue", nil)
	}

	// Use the retrieved queue name to create workflow options
	opts := opts_send(queue.NewQueue(queue.WithName(queue.Name(name))), user_id)

	_, err = shared.Temporal().Client().ExecuteWorkflow(ctx, opts, SendMessageWorkflow, user_id, message)
	if err != nil {
		return NewHubError(ErrorTypeWorkflowExecutionFailed, "failed to send message", err)
	}

	return nil
}

// Stop gracefully shuts down the hub and closes all client connections.
// It should be called when the application is shutting down to ensure
// all resources are properly released and all connections are closed.
//
// Example:
//
//	hub := ws.Instance()
//	// ... use hub
//	defer hub.Stop() // Ensure hub is stopped when main function exits
func (h *hub) Stop() {
	h.stop <- true
	h.mu.Lock()
	defer h.mu.Unlock()

	for client := range h.clients {
		close(client.send)
		client.conn.Close()
	}

	close(h.register)
	close(h.unregister)
}

// Internal hub operations

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

// worker sets up and manages the Temporal worker running against the queue,
// for handling workflows and activities.
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

// WebSocket connection handling

// read reads messages from the WebSocket connection.
func (h *hub) read(client *client) {
	defer func() {
		h.unregister <- client
		client.conn.Close()
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				shared.Logger().Error("Unexpected close error", "error", err)
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

// Additional functions

// ConnectionHandlerWorker creates and configures a Temporal worker for handling WebSocket connections.
// It sets up a new queue with the WebSocketQueue name and registers the ConnectionsHandlerWorkflow.
//
// Example:
//
//	worker := ws.ConnectionHandlerWorker()
//	err := worker.Start()
//	if err != nil {
//	    log.Fatalf("Failed to start worker: %v", err)
//	}
//	defer worker.Stop()
func ConnectionHandlerWorker() worker.Worker {
	q := queue.NewQueue(queue.WithName(shared.WebSocketQueue))

	worker := q.Worker(shared.Temporal().Client())

	worker.RegisterWorkflow(ConnectionsHandlerWorkflow)

	return worker
}

// Instance returns a singleton instance of the Hub.
// It initializes the hub if it hasn't been created yet, setting up necessary
// channels, registering workflows and activities, and starting the worker and run loops.
// This method ensures that only one hub instance is created and used throughout the application.
//
// Example:
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

		go instance.worker()
		go instance.run()
	})

	return instance
}
