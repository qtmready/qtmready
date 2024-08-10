package ws

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"go.temporal.io/sdk/worker"

	"go.breu.io/quantm/internal/shared"
	"go.breu.io/quantm/internal/shared/queue"
)

type (
	// Hub interface defines the methods for managing WebSocket connections and messaging.
	Hub interface {
		HandleWebSocket(c echo.Context) error
		Send(ctx context.Context, user_id string, message []byte) error
		Broadcast(ctx context.Context, team_id string, message []byte) error
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
		broadcast  chan []byte
		mu         sync.RWMutex
		queue      queue.Queue
		stop       chan bool
	}
)

var (
	_h   *hub
	once sync.Once
)

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

// run is the main loop that handles client registration, unregistration, and broadcasting.
func (h *hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		case <-h.stop:
			return
		}
	}
}

// HandleWebSocket upgrades an HTTP connection to a WebSocket connection and manages it.
// It creates a new Client, registers it with the hub, and starts read and write pumps
// to handle incoming and outgoing messages.
//
// Example:
//
//	e.GET("/ws/:id", func(c echo.Context) error {
//	    return ws.Instance().HandleWebSocket(c)
//	})
func (h *hub) HandleWebSocket(c echo.Context) error {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // You may want to implement a more secure check
		},
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return fmt.Errorf("failed to upgrade connection: %w", err)
	}

	userID := c.Param("id")

	client := &client{
		user_id: userID,
		conn:    conn,
		send:    make(chan []byte, 256),
	}

	h.register <- client

	go h.readPump(client)
	go h.writePump(client)

	return nil
}

// readPump reads messages from the WebSocket connection.
func (h *hub) readPump(client *client) {
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

// writePump writes messages to the WebSocket connection.
func (h *hub) writePump(client *client) {
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
	h.mu.RLock()
	client, found := h.findClient(user_id)
	h.mu.RUnlock()

	if found {
		select {
		case client.send <- message:
			return nil // Message sent locally
		default:
			// Client's send buffer is full or client is disconnected
			h.mu.Lock()
			if _, stillConnected := h.clients[client]; stillConnected {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
		}
	}

	// If we couldn't send locally, start a Temporal workflow
	opts := h.queue.WorkflowOptions(
		queue.WithWorkflowBlock("send"),
		queue.WithWorkflowBlockID(user_id),
	)

	_, err := shared.Temporal().Client().ExecuteWorkflow(ctx, opts, SendMessageWorkflow, user_id, message)
	if err != nil {
		shared.Logger().Error("Failed to start SendMessageWorkflow", "error", err)
		return err
	}

	return nil
}

// findClient searches for a client with the given user_id.
func (h *hub) findClient(user_id string) (*client, bool) {
	for client := range h.clients {
		if client.user_id == user_id {
			return client, true
		}
	}

	return nil, false
}

// Broadcast sends a message to all members of a team.
// It starts a Temporal workflow to distribute the message to all team members,
// which allows for reliable delivery even if some team members are offline.
//
// Example:
//
//	ctx := context.Background()
//	err := hub.Broadcast(ctx, "team456", []byte("Team announcement"))
//	if err != nil {
//	    log.Printf("Failed to broadcast message: %v", err)
//	}
func (h *hub) Broadcast(ctx context.Context, team_id string, message []byte) error {
	// Start a Temporal workflow to distribute the message to all team members
	opts := h.queue.WorkflowOptions(
		queue.WithWorkflowBlock("broadcast"),
		queue.WithWorkflowBlockID(team_id),
	)

	_, err := shared.Temporal().Client().ExecuteWorkflow(ctx, opts, BroadcastMessageWorkflow, team_id, message)
	if err != nil {
		shared.Logger().Error("Failed to start BroadcastMessageWorkflow", "error", err)
		return err
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
	for client := range h.clients {
		close(client.send)
		client.conn.Close()
	}
	h.mu.Unlock()
	close(h.broadcast)
	close(h.register)
	close(h.unregister)
}

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
		queueName := queue.Name(fmt.Sprintf("ws_%s", uuid.New().String()))
		_h = &hub{
			clients:    make(map[*client]bool),
			register:   make(chan *client),
			unregister: make(chan *client),
			broadcast:  make(chan []byte),
			queue:      queue.NewQueue(queue.WithName(queueName)),
			stop:       make(chan bool, 1),
		}

		go _h.worker()
		go _h.run()
	})

	return _h
}
