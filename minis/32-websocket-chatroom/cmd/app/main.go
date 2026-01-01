package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow connections from any origin (adjust for production)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Message represents a chat message with metadata
type Message struct {
	Type      string    `json:"type"`      // "message", "join", "leave", "error"
	Username  string    `json:"username"`
	Room      string    `json:"room"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// Client represents a WebSocket client connection
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	username string
	room     string
}

// Room represents a chat room with multiple clients
type Room struct {
	name       string
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	hub        *Hub
}

// Hub maintains active rooms and coordinates message routing
type Hub struct {
	rooms      map[string]*Room
	mu         sync.RWMutex
	register   chan *Client
	unregister chan *Client
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]*Room),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// GetOrCreateRoom returns existing room or creates a new one
func (h *Hub) GetOrCreateRoom(name string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.rooms[name]
	if !exists {
		room = &Room{
			name:       name,
			clients:    make(map[*Client]bool),
			broadcast:  make(chan []byte, 256),
			register:   make(chan *Client),
			unregister: make(chan *Client),
			hub:        h,
		}
		h.rooms[name] = room
		go room.run()
		log.Printf("Created new room: %s", name)
	}
	return room
}

// Run starts the hub's main event loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			room := h.GetOrCreateRoom(client.room)
			room.register <- client

		case client := <-h.unregister:
			h.mu.RLock()
			room, exists := h.rooms[client.room]
			h.mu.RUnlock()

			if exists {
				room.unregister <- client
			}
		}
	}
}

// Shutdown gracefully closes all rooms and connections
func (h *Hub) Shutdown() {
	h.mu.Lock()
	defer h.mu.Unlock()

	log.Println("Shutting down all rooms...")
	for _, room := range h.rooms {
		close(room.register)
		close(room.unregister)

		// Send close message to all clients in room
		for client := range room.clients {
			client.conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server shutting down"))
			close(client.send)
		}

		close(room.broadcast)
	}
}

// run handles the room's message broadcasting and client management
func (r *Room) run() {
	for {
		select {
		case client, ok := <-r.register:
			if !ok {
				return
			}

			r.clients[client] = true
			log.Printf("Client %s joined room %s (total: %d)", client.username, r.name, len(r.clients))

			// Create join notification
			msg := Message{
				Type:      "join",
				Username:  client.username,
				Room:      r.name,
				Content:   fmt.Sprintf("%s joined the room", client.username),
				Timestamp: time.Now(),
			}

			// Broadcast join notification
			if data, err := json.Marshal(msg); err == nil {
				r.broadcastToAll(data)
			}

			// Send welcome message to new client
			welcome := Message{
				Type:      "message",
				Username:  "System",
				Room:      r.name,
				Content:   fmt.Sprintf("Welcome to room '%s', %s!", r.name, client.username),
				Timestamp: time.Now(),
			}
			if data, err := json.Marshal(welcome); err == nil {
				select {
				case client.send <- data:
				default:
				}
			}

		case client, ok := <-r.unregister:
			if !ok {
				return
			}

			if _, exists := r.clients[client]; exists {
				delete(r.clients, client)
				close(client.send)
				log.Printf("Client %s left room %s (remaining: %d)", client.username, r.name, len(r.clients))

				// Create leave notification
				msg := Message{
					Type:      "leave",
					Username:  client.username,
					Room:      r.name,
					Content:   fmt.Sprintf("%s left the room", client.username),
					Timestamp: time.Now(),
				}

				// Broadcast leave notification
				if data, err := json.Marshal(msg); err == nil {
					r.broadcastToAll(data)
				}
			}

		case message, ok := <-r.broadcast:
			if !ok {
				return
			}

			r.broadcastToAll(message)
		}
	}
}

// broadcastToAll sends a message to all clients in the room
func (r *Room) broadcastToAll(message []byte) {
	for client := range r.clients {
		select {
		case client.send <- message:
			// Message queued successfully
		default:
			// Client's send channel is full (slow client)
			log.Printf("Client %s is slow, disconnecting", client.username)
			close(client.send)
			delete(r.clients, client)
		}
	}
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
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
		_, msgBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for %s: %v", c.username, err)
			}
			break
		}

		// Parse incoming message
		var inMsg map[string]string
		if err := json.Unmarshal(msgBytes, &inMsg); err != nil {
			// If not JSON, treat as plain text
			msg := Message{
				Type:      "message",
				Username:  c.username,
				Room:      c.room,
				Content:   string(msgBytes),
				Timestamp: time.Now(),
			}

			if data, err := json.Marshal(msg); err == nil {
				c.hub.mu.RLock()
				room := c.hub.rooms[c.room]
				c.hub.mu.RUnlock()

				if room != nil {
					room.broadcast <- data
				}
			}
			continue
		}

		// Handle structured message
		content := inMsg["content"]
		if content == "" {
			content = inMsg["message"] // Support both "content" and "message" fields
		}

		msg := Message{
			Type:      "message",
			Username:  c.username,
			Room:      c.room,
			Content:   content,
			Timestamp: time.Now(),
		}

		if data, err := json.Marshal(msg); err == nil {
			c.hub.mu.RLock()
			room := c.hub.rooms[c.room]
			c.hub.mu.RUnlock()

			if room != nil {
				room.broadcast <- data
			}
		}
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
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
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current WebSocket frame
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

// serveWs handles WebSocket requests from clients
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Extract room and username from query parameters
	room := r.URL.Query().Get("room")
	username := r.URL.Query().Get("user")

	if room == "" {
		room = "general" // Default room
	}

	if username == "" {
		http.Error(w, "Missing 'user' parameter", http.StatusBadRequest)
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}

	// Create new client
	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		username: username,
		room:     room,
	}

	// Register client with hub
	hub.register <- client

	// Start client goroutines
	go client.writePump()
	go client.readPump()
}

// serveHome serves a simple HTML page for testing
func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, homeHTML)
}

// serveRooms returns the list of active rooms
func serveRooms(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hub.mu.RLock()
		defer hub.mu.RUnlock()

		type RoomInfo struct {
			Name        string `json:"name"`
			ClientCount int    `json:"client_count"`
		}

		rooms := make([]RoomInfo, 0, len(hub.rooms))
		for name, room := range hub.rooms {
			rooms = append(rooms, RoomInfo{
				Name:        name,
				ClientCount: len(room.clients),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rooms)
	}
}

const homeHTML = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>WebSocket Chatroom</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 50px auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background: white;
            border-radius: 8px;
            padding: 20px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        h1 {
            color: #333;
            margin-top: 0;
        }
        .input-group {
            margin: 15px 0;
        }
        label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
            color: #555;
        }
        input, button {
            padding: 10px;
            font-size: 14px;
            border: 1px solid #ddd;
            border-radius: 4px;
        }
        input {
            width: calc(100% - 22px);
            margin-bottom: 10px;
        }
        button {
            background-color: #4CAF50;
            color: white;
            border: none;
            cursor: pointer;
            margin-right: 10px;
        }
        button:hover {
            background-color: #45a049;
        }
        button:disabled {
            background-color: #ccc;
            cursor: not-allowed;
        }
        #disconnect {
            background-color: #f44336;
        }
        #disconnect:hover {
            background-color: #da190b;
        }
        #messages {
            border: 1px solid #ddd;
            height: 400px;
            overflow-y: auto;
            padding: 10px;
            margin: 20px 0;
            background-color: #fafafa;
            border-radius: 4px;
        }
        .message {
            margin: 8px 0;
            padding: 8px;
            border-radius: 4px;
        }
        .message-user {
            background-color: #e3f2fd;
        }
        .message-system {
            background-color: #fff3e0;
            font-style: italic;
        }
        .message-join {
            background-color: #e8f5e9;
            font-style: italic;
        }
        .message-leave {
            background-color: #ffebee;
            font-style: italic;
        }
        .timestamp {
            font-size: 11px;
            color: #999;
            margin-right: 8px;
        }
        .username {
            font-weight: bold;
            margin-right: 8px;
        }
        .status {
            padding: 10px;
            margin: 10px 0;
            border-radius: 4px;
            background-color: #e3f2fd;
            color: #1976d2;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>WebSocket Chatroom</h1>

        <div class="input-group">
            <label>Username:</label>
            <input type="text" id="username" placeholder="Enter your username" value="User_${Math.floor(Math.random()*1000)}">
        </div>

        <div class="input-group">
            <label>Room:</label>
            <input type="text" id="room" placeholder="Enter room name" value="general">
        </div>

        <div class="input-group">
            <button id="connect">Connect</button>
            <button id="disconnect" disabled>Disconnect</button>
        </div>

        <div id="status" class="status">Not connected</div>

        <div id="messages"></div>

        <div class="input-group">
            <input type="text" id="messageInput" placeholder="Type a message..." disabled>
            <button id="send" disabled>Send</button>
        </div>
    </div>

    <script>
        let ws = null;
        const messagesDiv = document.getElementById('messages');
        const statusDiv = document.getElementById('status');
        const connectBtn = document.getElementById('connect');
        const disconnectBtn = document.getElementById('disconnect');
        const sendBtn = document.getElementById('send');
        const messageInput = document.getElementById('messageInput');
        const usernameInput = document.getElementById('username');
        const roomInput = document.getElementById('room');

        connectBtn.onclick = () => {
            const username = usernameInput.value.trim();
            const room = roomInput.value.trim();

            if (!username) {
                alert('Please enter a username');
                return;
            }

            if (!room) {
                alert('Please enter a room name');
                return;
            }

            connect(username, room);
        };

        disconnectBtn.onclick = () => {
            if (ws) {
                ws.close();
            }
        };

        sendBtn.onclick = sendMessage;

        messageInput.onkeypress = (e) => {
            if (e.key === 'Enter') {
                sendMessage();
            }
        };

        function connect(username, room) {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = protocol + '//' + window.location.host + '/ws?user=' + encodeURIComponent(username) + '&room=' + encodeURIComponent(room);

            ws = new WebSocket(wsUrl);

            ws.onopen = () => {
                statusDiv.textContent = 'Connected to room: ' + room;
                statusDiv.style.backgroundColor = '#e8f5e9';
                statusDiv.style.color = '#2e7d32';
                connectBtn.disabled = true;
                disconnectBtn.disabled = false;
                sendBtn.disabled = false;
                messageInput.disabled = false;
                usernameInput.disabled = true;
                roomInput.disabled = true;
            };

            ws.onmessage = (event) => {
                try {
                    const msg = JSON.parse(event.data);
                    addMessage(msg);
                } catch (e) {
                    // Plain text message
                    addMessage({
                        type: 'message',
                        username: 'Unknown',
                        content: event.data,
                        timestamp: new Date().toISOString()
                    });
                }
            };

            ws.onclose = () => {
                statusDiv.textContent = 'Disconnected';
                statusDiv.style.backgroundColor = '#ffebee';
                statusDiv.style.color = '#c62828';
                connectBtn.disabled = false;
                disconnectBtn.disabled = true;
                sendBtn.disabled = true;
                messageInput.disabled = true;
                usernameInput.disabled = false;
                roomInput.disabled = false;
                ws = null;
            };

            ws.onerror = (error) => {
                console.error('WebSocket error:', error);
                statusDiv.textContent = 'Connection error';
                statusDiv.style.backgroundColor = '#ffebee';
                statusDiv.style.color = '#c62828';
            };
        }

        function sendMessage() {
            const content = messageInput.value.trim();
            if (!content || !ws) return;

            const msg = JSON.stringify({ content: content });
            ws.send(msg);
            messageInput.value = '';
        }

        function addMessage(msg) {
            const div = document.createElement('div');
            div.className = 'message';

            const timestamp = new Date(msg.timestamp).toLocaleTimeString();
            const timestampSpan = document.createElement('span');
            timestampSpan.className = 'timestamp';
            timestampSpan.textContent = timestamp;

            if (msg.type === 'join') {
                div.className += ' message-join';
                div.innerHTML = timestampSpan.outerHTML + msg.content;
            } else if (msg.type === 'leave') {
                div.className += ' message-leave';
                div.innerHTML = timestampSpan.outerHTML + msg.content;
            } else if (msg.username === 'System') {
                div.className += ' message-system';
                div.innerHTML = timestampSpan.outerHTML + msg.content;
            } else {
                div.className += ' message-user';
                const usernameSpan = document.createElement('span');
                usernameSpan.className = 'username';
                usernameSpan.textContent = msg.username + ':';
                div.innerHTML = timestampSpan.outerHTML + usernameSpan.outerHTML + ' ' + msg.content;
            }

            messagesDiv.appendChild(div);
            messagesDiv.scrollTop = messagesDiv.scrollHeight;
        }
    </script>
</body>
</html>`

func main() {
	// Create hub
	hub := NewHub()

	// Start hub
	go hub.Run()

	// Setup HTTP routes
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.HandleFunc("/rooms", serveRooms(hub))

	// HTTP server configuration
	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Println("WebSocket chat server starting on :8080")
		log.Println("Open http://localhost:8080 in your browser")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("\nShutting down server...")

	// Shutdown hub (closes all WebSocket connections)
	hub.Shutdown()

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped gracefully")
}
