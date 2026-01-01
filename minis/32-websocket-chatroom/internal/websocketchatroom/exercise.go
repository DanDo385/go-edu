//go:build !solution && !reference

package websocketchatroom

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
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	username string
	roomName string
}

// Room represents a chat room.
// It manages multiple clients and broadcasts messages to all of them.
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
	// TODO: Implement this function.
	// - Initialize the Hub struct with a non-nil `rooms` map and buffered `register` and `unregister` channels.
	// - A buffer for the channels is good practice to prevent a `ServeWS` call from blocking if the hub is busy.
	return &Hub{
		rooms:      make(map[string]*Room),
		register:   make(chan *Client, 256),
		unregister: make(chan *Client, 256),
	}
}

// Run starts the hub's main event loop.
// It should handle client registration and unregistration.
func (h *Hub) Run() {
	// TODO: Implement the main event loop for the hub.
	// - Use an infinite `for` loop and a `select` statement.
	// - `case client := <-h.register:`:
	//   - A new client has connected.
	//   - Use `h.GetOrCreateRoom(client.roomName)` to get or create the room for the client.
	//   - Send the client to that room's `register` channel: `room.register <- client`.
	// - `case client := <-h.unregister:`:
	//   - A client has disconnected.
	//   - Find the client's room. You'll need a read lock to safely access the `h.rooms` map.
	//   - If the room exists, send the client to that room's `unregister` channel.
}

// GetOrCreateRoom returns an existing room or creates a new one.
// This method must be thread-safe.
func (h *Hub) GetOrCreateRoom(name string) *Room {
	// TODO: Implement this thread-safe "get or create" logic.

	// This is a common pattern called "double-checked locking".

	// Step 1: Check for the room with a read lock.
	// - `h.mu.RLock()`
	// - `room, exists := h.rooms[name]`
	// - `h.mu.RUnlock()`
	// - If `exists`, you're done. Return the room.

	// Step 2: If it doesn't exist, acquire a full write lock.
	// - `h.mu.Lock()`
	// - `defer h.mu.Unlock()`

	// Step 3: Check *again* inside the write lock.
	// - Why? Because another goroutine might have created the room between your read unlock and your write lock.
	// - `room, exists = h.rooms[name]`
	// - If it now exists, return it.

	// Step 4: If it still doesn't exist, create it.
	// - Create a new `Room` instance with initialized channels and maps.
	// - Start its event loop in a new goroutine: `go room.Run()`.
	// - Store it in the hub's map: `h.rooms[name] = room`.
	// - Return the new room.
	return nil
}

// Shutdown gracefully closes all rooms and connections.
func (h *Hub) Shutdown() {
	// TODO: Implement graceful shutdown.
	// - Acquire a write lock on the hub.
	// - Loop through all the rooms.
	// - For each room, close its `register`, `unregister`, and `broadcast` channels.
	// - For each client in the room, send a WebSocket close message and close their `send` channel.
}

// Run starts the room's event loop.
// It handles client registration, unregistration, and message broadcasting.
func (r *Room) Run() {
	// TODO: Implement the main event loop for the room.
	// - Use an infinite `for` loop and a `select` statement.
	// - `case client := <-r.register:`
	//   - Add the client to the `r.clients` map.
	//   - Broadcast a "user joined" message to the room.
	// - `case client := <-r.unregister:`
	//   - Check if the client exists in the map.
	//   - If it does, delete them from the map and close their `send` channel.
	//   - Broadcast a "user left" message.
	// - `case message := <-r.broadcast:`
	//   - This is a message from a client.
	//   - Call `r.broadcastToAll(message)` to send it to every client in the room.
}

// broadcastToAll sends a message to all clients in the room.
// Slow clients should be disconnected rather than blocking the broadcast.
func (r *Room) broadcastToAll(message []byte) {
	// TODO: Implement the broadcast logic.
	// - Loop through all clients in the `r.clients` map.
	// - For each client, attempt to send the message to their `send` channel.
	// - **Important:** This send must be non-blocking. Use a `select` with a `default` case.
	//   ```
	//   select {
	//   case client.send <- message:
	//   default:
	//       // The client's send buffer is full. Assume they are disconnected or lagging.
	//       close(client.send)
	//       delete(r.clients, client)
	//   }
	//   ```
	// - This prevents one slow client from blocking message delivery for the entire room.
}

// NewClient creates a new Client instance.
func NewClient(hub *Hub, conn *websocket.Conn, username, roomName string) *Client {
	// TODO: Implement this function.
	// - Initialize and return a `*Client`.
	// - The `send` channel should be buffered. A size of 256 is reasonable.
	return &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		username: username,
		roomName: roomName,
	}
}

// ReadPump reads messages from the WebSocket connection and broadcasts them to the room.
func (c *Client) ReadPump() {
	// TODO: Implement the read loop for a client.

	// Step 1: Defer the cleanup.
	// - When this function exits (either by error or clean shutdown), the client needs to be unregistered from the hub and the connection needs to be closed.
	// - `defer func() { c.hub.unregister <- c; c.conn.Close() }()`

	// Step 2: Configure the connection.
	// - Set a read limit to prevent huge messages: `c.conn.SetReadLimit(maxMessageSize)`.
	// - Set an initial read deadline: `c.conn.SetReadDeadline(time.Now().Add(pongWait))`.
	// - Set a pong handler. The pong handler's job is to reset the read deadline, keeping the connection alive.
	//   - `c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })`

	// Step 3: Start the read loop.
	// - Use an infinite `for` loop.
	// - Inside, call `c.conn.ReadMessage()`. This will block until a message is received.
	// - If `ReadMessage` returns an error, break the loop. Log the error if it's unexpected.
	// - If a message is received, prepend the username to it and send it to the room's `broadcast` channel.
}

// WritePump sends messages from the send channel to the WebSocket connection.
func (c *Client) WritePump() {
	// TODO: Implement the write loop for a client.

	// Step 1: Create a ticker for sending pings.
	// - `ticker := time.NewTicker(pingPeriod)`

	// Step 2: Defer cleanup.
	// - Stop the ticker and close the connection.
	// - `defer func() { ticker.Stop(); c.conn.Close() }()`

	// Step 3: Start the write loop.
	// - Use an infinite `for` loop and a `select` statement.
	// - `case message, ok := <-c.send:`
	//   - A message is ready to be sent from our application.
	//   - Set a write deadline: `c.conn.SetWriteDeadline(...)`.
	//   - If `!ok`, the `send` channel was closed. Send a WebSocket close message and return.
	//   - Get a writer for the message: `w, err := c.conn.NextWriter(websocket.TextMessage)`.
	//   - Write the message.
	//   - (Optional) "Batch" multiple messages from the `send` channel into a single WebSocket frame for efficiency.
	//   - Close the writer.
	// - `case <-ticker.C:`
	//   - The ticker fired. Time to send a ping.
	//   - Set a write deadline.
	//   - Send a ping message: `c.conn.WriteMessage(websocket.PingMessage, nil)`.
}

// ServeWS handles WebSocket upgrade requests.
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// TODO: Implement the WebSocket handler.

	// This function is the entry point for a new client. It takes a standard HTTP request and "upgrades" it to a persistent WebSocket connection.

	// Step 1: Extract room and username from the request's query parameters.
	// - `room := r.URL.Query().Get("room")`
	// - `username := r.URL.Query().Get("user")`
	// - Add validation: username is required, and room can have a default value (e.g., "general").

	// Step 2: Upgrade the HTTP connection to a WebSocket connection.
	// - `conn, err := upgrader.Upgrade(w, r, nil)`
	// - Handle the error if the upgrade fails.

	// Step 3: Create a new client.
	// - `client := NewClient(hub, conn, username, room)`

	// Step 4: Register the new client with the hub.
	// - `hub.register <- client`

	// Step 5: Start the client's read and write pumps.
	// - These must run in their own goroutines to allow the HTTP handler to return while the connection persists.
	// - `go client.WritePump()`
	// - `go client.ReadPump()`
}
