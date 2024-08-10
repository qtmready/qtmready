package ws

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"go.breu.io/quantm/internal/shared"
	"go.breu.io/quantm/internal/shared/queue"
)

type (
	Client struct {
		UserID string
		Conn   *websocket.Conn
		Send   chan []byte
	}

	hub struct {
		clients    map[*Client]bool
		register   chan *Client
		unregister chan *Client
		broadcast  chan []byte
		mu         sync.RWMutex
		queue      queue.Queue
		stop       chan bool
	}

	Hub interface {
		HandleWebSocket(c echo.Context) error
		Send(ctx context.Context, user_id string, message []byte) error
		Broadcast(ctx context.Context, team_id string, message []byte) error
		Stop()
	}
)

var (
	_h   *hub
	once sync.Once
)

func Instance() Hub {
	once.Do(func() {
		queueName := queue.Name(fmt.Sprintf("ws_%s", uuid.New().String()))
		_h = &hub{
			clients:    make(map[*Client]bool),
			register:   make(chan *Client),
			unregister: make(chan *Client),
			broadcast:  make(chan []byte),
			queue:      queue.NewQueue(queue.WithName(queueName)),
			stop:       make(chan bool, 1),
		}

		// Create and start the worker
		worker := _h.queue.Worker(shared.Temporal().Client())

		// Register workflows
		worker.RegisterWorkflow(SendMessageWorkflow)
		worker.RegisterWorkflow(BroadcastMessageWorkflow)

		// Register activities
		activities := NewActivities(_h)
		worker.RegisterActivity(activities.SendMessage)
		worker.RegisterActivity(activities.GetTeamUsers)
		worker.RegisterActivity(activities.BroadcastMessage)

		// Start the worker
		go func() {
			err := worker.Start()
			if err != nil {
				shared.Logger().Error("Failed to start worker", "error", err)
				panic(err)
			}

			<-_h.stop

			// Graceful shutdown
			worker.Stop()
		}()

		go _h.run()
	})

	return _h
}

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
				close(client.Send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		case <-h.stop:
			return
		}
	}
}

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

	client := &Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	h.register <- client

	go h.readPump(client)
	go h.writePump(client)

	return nil
}

func (h *hub) readPump(client *Client) {
	defer func() {
		h.unregister <- client
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				shared.Logger().Error("Unexpected close error", "error", err)
			}

			break
		}

		// Handle incoming messages
		shared.Logger().Info("Received message from client", "userID", client.UserID, "message", string(message))
	}
}

func (h *hub) writePump(client *Client) {
	defer func() {
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				return
			}

			w, err := client.Conn.NextWriter(websocket.TextMessage)
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

func (h *hub) Send(ctx context.Context, user_id string, message []byte) error {
	h.mu.RLock()
	client, found := h.findClient(user_id)
	h.mu.RUnlock()

	if found {
		select {
		case client.Send <- message:
			return nil // Message sent locally
		default:
			// Client's send buffer is full or client is disconnected
			h.mu.Lock()
			if _, stillConnected := h.clients[client]; stillConnected {
				delete(h.clients, client)
				close(client.Send)
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

func (h *hub) findClient(user_id string) (*Client, bool) {
	for client := range h.clients {
		if client.UserID == user_id {
			return client, true
		}
	}

	return nil, false
}

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

func (h *hub) Stop() {
	h.stop <- true
	h.mu.Lock()
	for client := range h.clients {
		close(client.Send)
		client.Conn.Close()
	}
	h.mu.Unlock()
	close(h.broadcast)
	close(h.register)
	close(h.unregister)
}
