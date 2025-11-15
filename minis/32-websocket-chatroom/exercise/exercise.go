//go:build !solution
// +build !solution

package exercise

import (
	"net/http"
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
// Each client belongs to one room and can send/receive messages.
type Client struct {
	// TODO: Add fields:
	// - hub: pointer to Hub
	// - conn: WebSocket connection
	// - send: buffered channel for outbound messages
	// - username: client's display name
	// - roomName: name of the room this client is in
}

// Room represents a chat room.
// It manages multiple clients and broadcasts messages to all of them.
type Room struct {
	// TODO: Add fields:
	// - name: room identifier
	// - clients: map of active clients
	// - broadcast: channel for messages to broadcast
	// - register: channel for registering new clients
	// - unregister: channel for removing clients
}

// Hub maintains all active rooms and coordinates client connections.
type Hub struct {
	// TODO: Add fields:
	// - rooms: map of room name to Room
	// - mu: mutex to protect rooms map
	// - register: channel for registering clients
	// - unregister: channel for unregistering clients
}

// NewHub creates and initializes a new Hub.
func NewHub() *Hub {
	// TODO: Initialize Hub with:
	// - Empty rooms map
	// - Initialized mutex
	// - Buffered channels for register/unregister
	return nil
}

// Run starts the hub's main event loop.
// It should handle client registration and unregistration.
func (h *Hub) Run() {
	// TODO: Implement event loop that:
	// - Listens on register channel
	// - Listens on unregister channel
	// - Routes clients to appropriate rooms
}

// GetOrCreateRoom returns an existing room or creates a new one.
// This method must be thread-safe.
func (h *Hub) GetOrCreateRoom(name string) *Room {
	// TODO: Implement with proper locking:
	// - Check if room exists (read lock)
	// - If not, create new room (write lock)
	// - Start room's Run() goroutine for new rooms
	// - Return room
	return nil
}

// Shutdown gracefully closes all rooms and connections.
func (h *Hub) Shutdown() {
	// TODO: Implement graceful shutdown:
	// - Close all room channels
	// - Send close message to all clients
	// - Close client send channels
}

// Run starts the room's event loop.
// It handles client registration, unregistration, and message broadcasting.
func (r *Room) Run() {
	// TODO: Implement event loop using select:
	// - Handle register: add client to map, broadcast join notification
	// - Handle unregister: remove client, close send channel, broadcast leave notification
	// - Handle broadcast: send message to all clients (non-blocking)
	// Hint: Use select with default case to prevent blocking on slow clients
}

// broadcastToAll sends a message to all clients in the room.
// Slow clients should be disconnected rather than blocking the broadcast.
func (r *Room) broadcastToAll(message []byte) {
	// TODO: Implement broadcasting:
	// - Iterate over all clients
	// - Try to send message (use select with default)
	// - If send would block, close client and remove from map
}

// NewClient creates a new Client instance.
func NewClient(hub *Hub, conn *websocket.Conn, username, roomName string) *Client {
	// TODO: Initialize Client with:
	// - hub reference
	// - WebSocket connection
	// - Buffered send channel (size 256)
	// - username and roomName
	return nil
}

// ReadPump reads messages from the WebSocket connection and broadcasts them to the room.
// It runs in its own goroutine and handles:
// - Setting read deadline and pong handler
// - Reading messages in a loop
// - Broadcasting messages to the room
// - Cleaning up on disconnect
func (c *Client) ReadPump() {
	// TODO: Implement read loop:
	// 1. Set up defer to unregister client and close connection
	// 2. Configure read limit and read deadline
	// 3. Set pong handler that updates read deadline
	// 4. Loop: read messages and send to room's broadcast channel
	// 5. Handle errors and close gracefully
}

// WritePump sends messages from the send channel to the WebSocket connection.
// It runs in its own goroutine and handles:
// - Sending queued messages
// - Sending periodic pings
// - Handling channel closure
func (c *Client) WritePump() {
	// TODO: Implement write loop:
	// 1. Create ticker for ping period
	// 2. Set up defer to stop ticker and close connection
	// 3. Loop using select:
	//    - Read from send channel: write message to WebSocket
	//    - Ticker fires: send ping message
	// 4. Handle send channel closure (send close frame)
	// 5. Set write deadlines before each write
}

// ServeWS handles WebSocket upgrade requests.
// It should:
// - Extract room and username from query parameters
// - Upgrade HTTP connection to WebSocket
// - Create a new client
// - Register client with hub
// - Start ReadPump and WritePump goroutines
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// TODO: Implement WebSocket handler:
	// 1. Get "room" and "user" from query parameters
	// 2. Validate parameters (user is required, room defaults to "general")
	// 3. Upgrade connection using upgrader.Upgrade()
	// 4. Create new Client
	// 5. Register client with hub
	// 6. Start WritePump and ReadPump goroutines
}
