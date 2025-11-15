//go:build solution
// +build solution

package exercise

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocket timeouts and limits
const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Upgrader configures the WebSocket upgrade from HTTP
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// Client represents a WebSocket client connection.
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	username string
	roomName string
}

// Room represents a chat room.
type Room struct {
	name       string
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// Hub maintains all active rooms and coordinates client connections.
type Hub struct {
	rooms      map[string]*Room
	mu         sync.RWMutex
	register   chan *Client
	unregister chan *Client
}

// NewHub creates and initializes a new Hub.
func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]*Room),
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
	}
}

// Run starts the hub's main event loop.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			room := h.GetOrCreateRoom(client.roomName)
			room.register <- client

		case client := <-h.unregister:
			h.mu.RLock()
			room, exists := h.rooms[client.roomName]
			h.mu.RUnlock()

			if exists {
				room.unregister <- client
			}
		}
	}
}

// GetOrCreateRoom returns an existing room or creates a new one.
func (h *Hub) GetOrCreateRoom(name string) *Room {
	// First try with read lock
	h.mu.RLock()
	room, exists := h.rooms[name]
	h.mu.RUnlock()

	if exists {
		return room
	}

	// Need to create room, acquire write lock
	h.mu.Lock()
	defer h.mu.Unlock()

	// Check again in case another goroutine created it
	room, exists = h.rooms[name]
	if exists {
		return room
	}

	// Create new room
	room = &Room{
		name:       name,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	h.rooms[name] = room
	go room.Run()

	return room
}

// Shutdown gracefully closes all rooms and connections.
func (h *Hub) Shutdown() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, room := range h.rooms {
		// Close room channels
		close(room.register)
		close(room.unregister)

		// Close all clients in room
		for client := range room.clients {
			client.conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server shutting down"))
			close(client.send)
		}

		close(room.broadcast)
	}
}

// Run starts the room's event loop.
func (r *Room) Run() {
	for {
		select {
		case client, ok := <-r.register:
			if !ok {
				return
			}

			r.clients[client] = true

			// Broadcast join notification
			joinMsg := []byte(client.username + " joined room '" + r.name + "'")
			r.broadcastToAll(joinMsg)

		case client, ok := <-r.unregister:
			if !ok {
				return
			}

			if _, exists := r.clients[client]; exists {
				delete(r.clients, client)
				close(client.send)

				// Broadcast leave notification
				leaveMsg := []byte(client.username + " left room '" + r.name + "'")
				r.broadcastToAll(leaveMsg)
			}

		case message, ok := <-r.broadcast:
			if !ok {
				return
			}

			r.broadcastToAll(message)
		}
	}
}

// broadcastToAll sends a message to all clients in the room.
func (r *Room) broadcastToAll(message []byte) {
	for client := range r.clients {
		select {
		case client.send <- message:
			// Message queued successfully
		default:
			// Client's send buffer is full, disconnect them
			close(client.send)
			delete(r.clients, client)
		}
	}
}

// NewClient creates a new Client instance.
func NewClient(hub *Hub, conn *websocket.Conn, username, roomName string) *Client {
	return &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		username: username,
		roomName: roomName,
	}
}

// ReadPump reads messages from the WebSocket connection.
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Prepend username to message
		fullMessage := []byte(c.username + ": " + string(message))

		// Get room and broadcast
		c.hub.mu.RLock()
		room := c.hub.rooms[c.roomName]
		c.hub.mu.RUnlock()

		if room != nil {
			room.broadcast <- fullMessage
		}
	}
}

// WritePump sends messages to the WebSocket connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Room closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to current websocket frame
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWS handles WebSocket upgrade requests.
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Extract parameters
	room := r.URL.Query().Get("room")
	username := r.URL.Query().Get("user")

	if room == "" {
		room = "general" // Default room
	}

	if username == "" {
		http.Error(w, "Missing 'user' parameter", http.StatusBadRequest)
		return
	}

	// Upgrade connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}

	// Create client
	client := NewClient(hub, conn, username, room)

	// Register with hub
	hub.register <- client

	// Start goroutines
	go client.WritePump()
	go client.ReadPump()
}
