# Project 32: WebSocket Chatroom

## 1. What Is This About?

### Real-World Scenario

You're building a customer support chat application. With traditional HTTP:

**❌ What happens:**
1. Client sends AJAX request every second: "Any new messages?"
2. Server responds: "No" (99% of the time)
3. 3,600 HTTP requests per hour per user
4. High server load, high latency (500-1000ms delays)
5. Messages arrive in 1-second batches (feels slow)
6. Wasted bandwidth: HTTP headers (200-500 bytes) repeated every second

**✅ With WebSockets:**
1. Client opens one persistent connection
2. Server pushes messages instantly when they arrive
3. 1 connection per user (not 3,600 requests/hour)
4. Low latency (<50ms message delivery)
5. Messages arrive in real-time (feels instant)
6. Minimal bandwidth: Just message data (no repeated headers)

This project teaches you how to build **real-time applications** with:
- **WebSocket protocol**: Bidirectional, persistent connections
- **gorilla/websocket**: Production-ready WebSocket library
- **Chat rooms**: Multi-user message broadcasting
- **Connection management**: Handling joins, leaves, disconnects
- **Concurrency safety**: Thread-safe message handling

### What You'll Learn

1. **WebSocket protocol**: Upgrade handshake, frames, control messages
2. **gorilla/websocket library**: Upgrader, Conn, message types
3. **Real-time patterns**: Broadcasting, pub/sub, presence
4. **Connection lifecycle**: Connect, authenticate, disconnect, cleanup
5. **Concurrent broadcasting**: Safe message distribution to multiple clients
6. **Production patterns**: Heartbeats, timeouts, error handling

### The Challenge

Build a WebSocket chat server with:
- Multiple chat rooms (users can join specific rooms)
- Real-time message broadcasting within rooms
- User join/leave notifications
- Connection management with proper cleanup
- Concurrent client handling
- Graceful shutdown support

---

## 2. First Principles: WebSocket Protocol

### What is WebSocket?

**WebSocket** is a protocol that provides full-duplex communication over a single TCP connection.

**HTTP vs WebSocket**:
```
HTTP (Request-Response):
Client              Server
  |--- Request --→   |
  |                  |
  |←-- Response ---|  |
  |                  |
  |--- Request --→   |  (New connection each time)

WebSocket (Bidirectional):
Client              Server
  |--- Upgrade --→   |
  |←-- 101 OK -------|  (Connection established)
  |                  |
  |←-- Message -------|  (Server can send anytime)
  |--- Message --→   |  (Client can send anytime)
  |←-- Message -------|
  |--- Message --→   |
  |                  |  (Same connection!)
```

**Key differences**:
| Feature | HTTP | WebSocket |
|---------|------|-----------|
| Direction | Request → Response | Bidirectional |
| Connection | New per request | Persistent |
| Server push | Impossible | Native |
| Overhead | ~200-500 bytes/req | ~2-14 bytes/frame |
| Latency | Higher (TCP handshake) | Lower (persistent) |

### WebSocket Handshake

**Step 1: Client initiates upgrade**:
```http
GET /chat HTTP/1.1
Host: example.com
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
Sec-WebSocket-Version: 13
```

**Step 2: Server accepts**:
```http
HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
```

**Step 3: Connection established**:
- HTTP protocol "upgraded" to WebSocket
- Same TCP connection now speaks WebSocket protocol
- Both sides can send messages anytime

**In Go**:
```go
upgrader := websocket.Upgrader{}
conn, err := upgrader.Upgrade(w, r, nil)
// conn is now a WebSocket connection
```

### WebSocket Frames

**Messages are sent in frames**:

```
Frame structure:
┌─────┬─────┬──────────┬──────────┬─────────────┐
│ FIN │ Op  │ Mask bit │ Length   │ Payload     │
├─────┼─────┼──────────┼──────────┼─────────────┤
│ 1   │ 1   │ 1        │ 7/16/64  │ 0-2^63      │
└─────┴─────┴──────────┴──────────┴─────────────┘
```

**Frame types (opcodes)**:
- `0x1`: Text frame (UTF-8 text)
- `0x2`: Binary frame (arbitrary data)
- `0x8`: Close frame (connection closing)
- `0x9`: Ping frame (heartbeat check)
- `0xA`: Pong frame (heartbeat response)

**Example**:
```go
// Send text message
conn.WriteMessage(websocket.TextMessage, []byte("Hello"))

// Send binary data
conn.WriteMessage(websocket.BinaryMessage, imageData)

// Send ping (heartbeat)
conn.WriteMessage(websocket.PingMessage, nil)
```

### Why WebSocket Matters

**Before WebSocket (polling)**:
```
Every 1 second:
Client → Server: "Any updates?"
Server → Client: "No"
Server → Client: "No"
Server → Client: "No"
Server → Client: "Yes: New message!"

Cost: 3,600 requests/hour per user
```

**With WebSocket**:
```
Once:
Client ↔ Server: Connection established

When events occur:
Server → Client: "New message!"
Client → Server: "User typing..."

Cost: 1 connection per user
```

**Real-world impact**:
- **Slack**: 10M+ concurrent WebSocket connections
- **WhatsApp**: Handles billions of messages via persistent connections
- **Trading platforms**: Sub-millisecond price updates
- **Multiplayer games**: Real-time position synchronization

---

## 3. Breaking Down the Solution

### Architecture Overview

```
                    WebSocket Chat Server
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
    Room: "general"    Room: "tech"     Room: "random"
        │                  │                  │
    ┌───┴───┐          ┌───┴───┐         ┌───┴───┐
    │ User1 │          │ User3 │         │ User5 │
    │ User2 │          │ User4 │         └───────┘
    └───────┘          └───────┘

Message flow:
1. User1 sends: "Hello" to room "general"
2. Server broadcasts to all users in "general"
3. User1, User2 receive message
4. User3, User4, User5 do NOT receive (different rooms)
```

### Core Components

**1. Client (Connection)**:
```go
type Client struct {
    conn     *websocket.Conn  // WebSocket connection
    send     chan []byte       // Outbound message channel
    room     *Room             // Which room this client is in
    username string            // User identifier
}
```

**Why `send` channel?**
- Decouples message production from sending
- Prevents blocking when client is slow
- Enables graceful shutdown

**2. Room (Broadcast Hub)**:
```go
type Room struct {
    name       string                  // Room identifier
    clients    map[*Client]bool       // Connected clients
    broadcast  chan []byte             // Inbound messages to broadcast
    register   chan *Client            // Register new client
    unregister chan *Client            // Remove client
}
```

**Why channels?**
- Centralizes concurrent access to `clients` map
- Avoids data races
- Provides clean synchronization

**3. Hub (Room Manager)**:
```go
type Hub struct {
    rooms      map[string]*Room       // All available rooms
    mu         sync.RWMutex           // Protects rooms map
}
```

### Connection Lifecycle

**Phase 1: Connection establishment**:
```go
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
    // 1. Upgrade HTTP → WebSocket
    conn, err := upgrader.Upgrade(w, r, nil)

    // 2. Extract room from URL
    roomName := r.URL.Query().Get("room")

    // 3. Create client
    client := &Client{
        conn: conn,
        send: make(chan []byte, 256),
        room: hub.getRoom(roomName),
    }

    // 4. Register with room
    client.room.register <- client

    // 5. Start goroutines
    go client.writePump()  // Send messages
    go client.readPump()   // Receive messages
}
```

**Phase 2: Active communication**:
```
Client goroutines:

writePump():                    readPump():
┌─────────────────┐            ┌─────────────────┐
│ for msg := range│            │ for {           │
│   client.send { │            │   msg := conn.  │
│                 │            │     ReadMessage │
│   conn.Write(msg)│           │                 │
│ }               │            │   room.broadcast│
└─────────────────┘            │     <- msg      │
     ↑                         └─────────────────┘
     │                              │
     └──────────────────────────────┘
          Message loop
```

**Phase 3: Disconnection**:
```go
// Triggered by: client disconnect, error, or shutdown
defer func() {
    client.room.unregister <- client
    client.conn.Close()
}()
```

### Broadcasting Pattern

**Room's run() method (central message hub)**:
```go
func (r *Room) run() {
    for {
        select {
        case client := <-r.register:
            r.clients[client] = true

        case client := <-r.unregister:
            if _, ok := r.clients[client]; ok {
                delete(r.clients, client)
                close(client.send)
            }

        case message := <-r.broadcast:
            // Send to ALL clients in room
            for client := range r.clients {
                select {
                case client.send <- message:
                    // Message queued successfully
                default:
                    // Client's send buffer full, disconnect
                    close(client.send)
                    delete(r.clients, client)
                }
            }
        }
    }
}
```

**Why this pattern?**
- **Single goroutine** manages all room state (no races)
- **Non-blocking sends**: `select` with `default` prevents slow clients from blocking others
- **Automatic cleanup**: Slow/dead clients are disconnected

**Visual flow**:
```
User1 types: "Hello"
     │
     ↓
readPump() receives message
     │
     ↓
room.broadcast <- message
     │
     ↓
Room's run() receives on broadcast channel
     │
     ↓
for client := range room.clients {
     │
     ├→ client1.send <- message
     ├→ client2.send <- message
     └→ client3.send <- message
}
     │
     ↓
Each writePump() receives and sends to WebSocket
     │
     ↓
All users see: "Hello"
```

---

## 4. Complete Solution Walkthrough

### Step 1: Client Structure

```go
type Client struct {
    room     *Room
    conn     *websocket.Conn
    send     chan []byte
    username string
}
```

**Field purposes**:
- `room`: Which room to broadcast messages to
- `conn`: WebSocket connection for reading/writing
- `send`: Buffered channel to queue outbound messages
- `username`: Display name for this user

**Why buffered `send` channel (256)?**
- If server broadcasts 300 messages instantly, client might not send fast enough
- Buffer absorbs bursts
- If buffer fills, client is too slow → disconnect

### Step 2: Reading Messages (readPump)

```go
func (c *Client) readPump() {
    defer func() {
        c.room.unregister <- c
        c.conn.Close()
    }()

    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }

        // Prepend username
        fullMessage := []byte(c.username + ": " + string(message))
        c.room.broadcast <- fullMessage
    }
}
```

**Key mechanisms**:

**1. Read deadline**:
```go
c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
```
- If no message received within 60 seconds, connection errors
- Prevents zombie connections
- Reset by pong messages (heartbeats)

**2. Pong handler**:
```go
c.conn.SetPongHandler(func(string) error {
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    return nil
})
```
- Server sends ping every 54 seconds
- Client automatically responds with pong
- Pong resets read deadline → connection stays alive

**3. Message broadcasting**:
```go
c.room.broadcast <- fullMessage
```
- Non-blocking send to room's broadcast channel
- Room's run() goroutine handles actual distribution

### Step 3: Writing Messages (writePump)

```go
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
                // Room closed channel
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
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
```

**Key mechanisms**:

**1. Ping ticker**:
```go
ticker := time.NewTicker(54 * time.Second)
```
- Sends ping every 54 seconds
- Keeps connection alive through proxies/firewalls
- Detects dead connections (if client doesn't pong back)

**2. Write deadline**:
```go
c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
```
- If write takes >10 seconds, give up
- Prevents blocking on slow clients

**3. Channel closure detection**:
```go
message, ok := <-c.send
if !ok {
    // Channel closed by room
    c.conn.WriteMessage(websocket.CloseMessage, []byte{})
    return
}
```
- Room closes `send` channel when unregistering client
- Triggers graceful WebSocket close

### Step 4: Room Management

```go
type Room struct {
    name       string
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

func (r *Room) run() {
    for {
        select {
        case client := <-r.register:
            r.clients[client] = true

            // Broadcast join notification
            msg := []byte(client.username + " joined the room")
            for c := range r.clients {
                select {
                case c.send <- msg:
                default:
                    close(c.send)
                    delete(r.clients, c)
                }
            }

        case client := <-r.unregister:
            if _, ok := r.clients[client]; ok {
                delete(r.clients, client)
                close(client.send)

                // Broadcast leave notification
                msg := []byte(client.username + " left the room")
                for c := range r.clients {
                    select {
                    case c.send <- msg:
                    default:
                        close(c.send)
                        delete(r.clients, c)
                    }
                }
            }

        case message := <-r.broadcast:
            for client := range r.clients {
                select {
                case client.send <- message:
                default:
                    // Send buffer full, client too slow
                    close(client.send)
                    delete(r.clients, client)
                }
            }
        }
    }
}
```

**Why single goroutine?**
- All room state (`clients` map) accessed in one goroutine
- No races, no mutexes needed
- Channels provide synchronization

**Slow client handling**:
```go
select {
case client.send <- message:
    // Success
default:
    // Channel full, client too slow
    close(client.send)
    delete(r.clients, client)
}
```
- Prevents one slow client from blocking others
- Automatically disconnects clients that can't keep up

### Step 5: HTTP Upgrade Handler

```go
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true  // Allow all origins (adjust for production)
    },
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
    // Extract parameters
    roomName := r.URL.Query().Get("room")
    username := r.URL.Query().Get("user")

    if roomName == "" || username == "" {
        http.Error(w, "Missing room or user parameter", http.StatusBadRequest)
        return
    }

    // Upgrade connection
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrade error:", err)
        return
    }

    // Get or create room
    room := hub.getOrCreateRoom(roomName)

    // Create client
    client := &Client{
        room:     room,
        conn:     conn,
        send:     make(chan []byte, 256),
        username: username,
    }

    // Register client
    client.room.register <- client

    // Start goroutines
    go client.writePump()
    go client.readPump()
}
```

**Connection URL format**:
```
ws://localhost:8080/ws?room=general&user=Alice
                          └─────┬─────┘ └────┬────┘
                            Room name      Username
```

**Upgrader configuration**:
- `ReadBufferSize`: Size of read buffer per connection
- `WriteBufferSize`: Size of write buffer per connection
- `CheckOrigin`: CORS validation (set properly for production!)

---

## 5. Key Concepts Explained

### Concept 1: Why Two Goroutines Per Client?

**Design**:
```go
go client.readPump()   // Goroutine 1: Read from WebSocket
go client.writePump()  // Goroutine 2: Write to WebSocket
```

**Why not one goroutine?**

**❌ Single goroutine attempt**:
```go
for {
    // Read message
    msg, err := conn.ReadMessage()

    // Process and broadcast
    room.broadcast <- msg

    // Try to send queued messages
    select {
    case outMsg := <-client.send:
        conn.WriteMessage(websocket.TextMessage, outMsg)
    default:
    }
}
```

**Problem**: Can't send while waiting for read!

**✅ Two goroutines**:
- `readPump()`: Blocks on `conn.ReadMessage()` (waiting for incoming)
- `writePump()`: Blocks on `<-client.send` (waiting for outgoing)
- Both can run concurrently!

### Concept 2: Ping/Pong Heartbeats

**Purpose**: Detect dead connections

**Scenario without heartbeats**:
```
1. Client loses internet connection
2. TCP doesn't immediately notice (no data being sent)
3. Server thinks client is still connected
4. Server keeps broadcasting messages to dead client
5. Memory leak: Dead client never cleaned up
```

**With heartbeats**:
```
Server writePump():
  Every 54s: Send ping

Client (browser):
  Auto-responds with pong

Server readPump():
  SetReadDeadline(60s)
  When pong received: Reset deadline

If no pong within 60s:
  → ReadDeadline exceeded
  → Connection errors
  → Cleanup triggered
```

**Timing**:
```go
const (
    writeWait  = 10 * time.Second      // Max write duration
    pongWait   = 60 * time.Second      // Max time between pongs
    pingPeriod = (pongWait * 9) / 10   // Send pings at 54s (before 60s timeout)
)
```

**Why `pingPeriod = 90% of pongWait`?**
- Accounts for network latency
- If ping sent at 54s, pong arrives by 56s
- Still 4 seconds before 60s timeout

### Concept 3: Buffered Send Channel

```go
send: make(chan []byte, 256)
```

**Why buffered?**

**Scenario 1: Burst of messages**:
```
Room receives 100 messages in 100ms
  ↓
for client := range clients {
    client.send <- message  // Must not block!
}
  ↓
Buffer absorbs burst
  ↓
writePump() sends at network speed
```

**Scenario 2: Slow client**:
```
Client network: 10 KB/s
Server broadcasts: 100 KB/s
  ↓
Buffer fills (256 messages)
  ↓
select {
    case client.send <- message:
    default:
        // Buffer full!
        close(client.send)
        delete(clients, client)
}
  ↓
Slow client disconnected, others unaffected
```

**Choosing buffer size**:
- Too small (1): Any burst causes disconnects
- Too large (10000): Slow clients stay connected, wasting memory
- Sweet spot (256): Handles bursts, disconnects legitimately slow clients

### Concept 4: Thread-Safe Room Access

**Problem**: Multiple HTTP requests hit `/ws` concurrently

```go
// RACE CONDITION!
type Hub struct {
    rooms map[string]*Room  // Multiple goroutines access this!
}

func (h *Hub) getRoom(name string) *Room {
    return h.rooms[name]  // Data race if another goroutine modifies map
}
```

**Solution 1: Mutex**:
```go
type Hub struct {
    rooms map[string]*Room
    mu    sync.RWMutex
}

func (h *Hub) getRoom(name string) *Room {
    h.mu.RLock()
    defer h.mu.RUnlock()
    return h.rooms[name]
}

func (h *Hub) createRoom(name string) *Room {
    h.mu.Lock()
    defer h.mu.Unlock()
    room := &Room{name: name, ...}
    h.rooms[name] = room
    return room
}
```

**Solution 2: sync.Map** (for high concurrency):
```go
type Hub struct {
    rooms sync.Map  // Optimized for concurrent access
}

func (h *Hub) getRoom(name string) *Room {
    val, _ := h.rooms.Load(name)
    return val.(*Room)
}
```

### Concept 5: Graceful Shutdown

**Challenge**: Server shutting down, 100 active WebSocket connections

**Without graceful shutdown**:
```go
os.Exit(0)
// All connections instantly closed
// Messages in transit lost
// Clients see unexpected disconnects
```

**With graceful shutdown**:
```go
func (h *Hub) Shutdown() {
    // 1. Stop accepting new connections
    listener.Close()

    // 2. Close all rooms
    for _, room := range h.rooms {
        close(room.register)    // No new clients
        close(room.unregister)  // No new unregisters

        // Send close message to all clients
        for client := range room.clients {
            client.conn.WriteMessage(websocket.CloseMessage, []byte{})
            close(client.send)
        }

        close(room.broadcast)
    }

    // 3. Wait for clients to disconnect
    time.Sleep(5 * time.Second)
}
```

**Client-side handling**:
```javascript
ws.onclose = (event) => {
    if (event.code === 1000) {
        // Normal closure (server shutdown)
        console.log("Server is shutting down, reconnecting in 5s...");
        setTimeout(reconnect, 5000);
    } else {
        // Abnormal closure (error)
        console.error("Connection error, reconnecting...");
        reconnect();
    }
};
```

---

## 6. Common Patterns You Can Reuse

### Pattern 1: Private Messaging

```go
type Client struct {
    id   string  // Unique client ID
    send chan []byte
}

type Hub struct {
    clients map[string]*Client  // ID → Client lookup
}

func (h *Hub) sendToClient(recipientID string, message []byte) {
    h.mu.RLock()
    client, ok := h.clients[recipientID]
    h.mu.RUnlock()

    if ok {
        select {
        case client.send <- message:
        default:
            // Client disconnected
        }
    }
}
```

**Message format**:
```json
{
    "type": "private",
    "to": "user123",
    "message": "Secret message"
}
```

### Pattern 2: Typing Indicators

```go
type TypingStatus struct {
    username string
    typing   bool
    lastSeen time.Time
}

type Room struct {
    typing map[string]*TypingStatus
}

func (r *Room) handleTyping(username string, isTyping bool) {
    r.typing[username] = &TypingStatus{
        username: username,
        typing:   isTyping,
        lastSeen: time.Now(),
    }

    // Broadcast typing status
    status := map[string]interface{}{
        "type":     "typing",
        "username": username,
        "typing":   isTyping,
    }

    data, _ := json.Marshal(status)
    r.broadcast <- data
}
```

### Pattern 3: Message History

```go
type Room struct {
    history    []Message
    historyMu  sync.RWMutex
    maxHistory int
}

func (r *Room) addToHistory(msg Message) {
    r.historyMu.Lock()
    defer r.historyMu.Unlock()

    r.history = append(r.history, msg)

    // Keep only last N messages
    if len(r.history) > r.maxHistory {
        r.history = r.history[1:]
    }
}

func (r *Room) getHistory() []Message {
    r.historyMu.RLock()
    defer r.historyMu.RUnlock()

    return append([]Message{}, r.history...)  // Return copy
}
```

**Send history on join**:
```go
case client := <-r.register:
    r.clients[client] = true

    // Send recent history
    for _, msg := range r.getHistory() {
        data, _ := json.Marshal(msg)
        client.send <- data
    }
```

### Pattern 4: Presence Tracking

```go
type Presence struct {
    userID     string
    status     string  // "online", "away", "offline"
    lastActive time.Time
}

func (r *Room) updatePresence(userID, status string) {
    presence := Presence{
        userID:     userID,
        status:     status,
        lastActive: time.Now(),
    }

    // Broadcast presence update
    data, _ := json.Marshal(map[string]interface{}{
        "type":     "presence",
        "user":     userID,
        "status":   status,
        "lastSeen": presence.lastActive,
    })

    r.broadcast <- data
}
```

### Pattern 5: Rate Limiting

```go
import "golang.org/x/time/rate"

type Client struct {
    limiter *rate.Limiter
}

func newClient() *Client {
    return &Client{
        limiter: rate.NewLimiter(rate.Every(100*time.Millisecond), 10),
    }
}

func (c *Client) readPump() {
    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }

        // Rate limit
        if !c.limiter.Allow() {
            c.send <- []byte("Error: Rate limit exceeded")
            continue
        }

        c.room.broadcast <- message
    }
}
```

---

## 7. Real-World Applications

### Live Chat Support

Customer support platforms like Intercom, Zendesk.

```go
type SupportHub struct {
    agents    map[string]*Client  // Support agents
    customers map[string]*Client  // Customers
    queue     chan *Client        // Waiting customers
}

func (h *SupportHub) assignAgent(customer *Client) {
    // Find available agent
    for agentID, agent := range h.agents {
        if agent.available {
            // Connect customer to agent
            room := createPrivateRoom(customer, agent)
            return
        }
    }

    // No agents available, queue customer
    h.queue <- customer
}
```

**Companies**: Intercom, Drift, Zendesk

### Collaborative Editing

Real-time document editing like Google Docs.

```go
type Document struct {
    content   string
    version   int
    broadcast chan Operation
}

type Operation struct {
    Type   string  // "insert", "delete"
    Pos    int
    Text   string
    UserID string
}

func (d *Document) applyOperation(op Operation) {
    // Apply operational transformation
    // Broadcast to all editors
    d.broadcast <- op
}
```

**Companies**: Google Docs, Notion, Figma

### Live Dashboards

Real-time metrics and monitoring.

```go
type MetricsHub struct {
    subscribers map[string]*Client
    metrics     chan Metric
}

func (h *MetricsHub) publishMetric(m Metric) {
    data, _ := json.Marshal(m)

    for _, client := range h.subscribers {
        select {
        case client.send <- data:
        default:
        }
    }
}
```

**Companies**: Grafana, Datadog, New Relic

### Multiplayer Games

Real-time game state synchronization.

```go
type GameRoom struct {
    players  map[string]*Player
    gameState *State
    tickRate time.Duration
}

func (r *GameRoom) gameLoop() {
    ticker := time.NewTicker(r.tickRate)

    for range ticker.C {
        // Update game state
        r.gameState.Update()

        // Broadcast to all players
        state, _ := json.Marshal(r.gameState)
        r.broadcast <- state
    }
}
```

**Companies**: Fortnite, Among Us, Agar.io

### Financial Trading

Real-time price feeds and order updates.

```go
type TradingHub struct {
    priceFeeds map[string]chan PriceUpdate
    orders     map[string]chan Order
}

func (h *TradingHub) publishPrice(symbol string, price float64) {
    update := PriceUpdate{
        Symbol:    symbol,
        Price:     price,
        Timestamp: time.Now(),
    }

    if feed, ok := h.priceFeeds[symbol]; ok {
        select {
        case feed <- update:
        default:
        }
    }
}
```

**Companies**: Coinbase, Robinhood, Bloomberg Terminal

---

## 8. Common Mistakes to Avoid

### Mistake 1: Not Setting Timeouts

**❌ Wrong**:
```go
conn.ReadMessage()  // Can block forever
```

**Problem**: Dead connection never times out, goroutine leaks.

**✅ Correct**:
```go
conn.SetReadDeadline(time.Now().Add(60 * time.Second))
conn.ReadMessage()
```

### Mistake 2: Blocking Broadcasts

**❌ Wrong**:
```go
for client := range room.clients {
    client.send <- message  // Blocks if channel full!
}
```

**Problem**: One slow client blocks entire broadcast loop.

**✅ Correct**:
```go
for client := range room.clients {
    select {
    case client.send <- message:
    default:
        // Client too slow, disconnect
        close(client.send)
        delete(room.clients, client)
    }
}
```

### Mistake 3: Concurrent Map Access

**❌ Wrong**:
```go
type Hub struct {
    rooms map[string]*Room
}

// Goroutine 1:
hub.rooms["general"] = newRoom()

// Goroutine 2:
room := hub.rooms["general"]  // RACE!
```

**✅ Correct**:
```go
type Hub struct {
    rooms map[string]*Room
    mu    sync.RWMutex
}

func (h *Hub) setRoom(name string, room *Room) {
    h.mu.Lock()
    defer h.mu.Unlock()
    h.rooms[name] = room
}

func (h *Hub) getRoom(name string) *Room {
    h.mu.RLock()
    defer h.mu.RUnlock()
    return h.rooms[name]
}
```

### Mistake 4: Not Handling Close Frames

**❌ Wrong**:
```go
conn.ReadMessage()  // Ignores close frames
```

**✅ Correct**:
```go
messageType, message, err := conn.ReadMessage()

if messageType == websocket.CloseMessage {
    // Client initiated close
    return
}

if err != nil {
    if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
        // Normal close
        return
    }
    // Unexpected error
    log.Println("Error:", err)
    return
}
```

### Mistake 5: Forgetting CheckOrigin

**❌ Wrong**:
```go
upgrader := websocket.Upgrader{}
// CheckOrigin defaults to same-origin only
```

**Problem**: Can't connect from different origins (CORS).

**✅ Correct (development)**:
```go
upgrader := websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true  // Allow all origins
    },
}
```

**✅ Correct (production)**:
```go
upgrader := websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        origin := r.Header.Get("Origin")
        return origin == "https://yourdomain.com"
    },
}
```

### Mistake 6: Not Closing Channels

**❌ Wrong**:
```go
case client := <-r.unregister:
    delete(r.clients, client)
    // Forgot to close client.send!
```

**Problem**: `writePump()` blocks forever on `<-client.send`.

**✅ Correct**:
```go
case client := <-r.unregister:
    if _, ok := r.clients[client]; ok {
        delete(r.clients, client)
        close(client.send)  // Signals writePump to exit
    }
```

---

## 9. Stretch Goals

### Goal 1: Add Authentication ⭐

Require JWT token for WebSocket connections.

**Hint**:
```go
func serveWs(w http.ResponseWriter, r *http.Request) {
    // Verify token
    token := r.URL.Query().Get("token")
    claims, err := verifyJWT(token)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    client.username = claims.Username
    client.userID = claims.UserID
}
```

### Goal 2: Persist Messages ⭐⭐

Store messages in database (PostgreSQL).

**Hint**:
```go
func (r *Room) handleMessage(msg Message) {
    // Save to database
    if err := db.SaveMessage(msg); err != nil {
        log.Println("Failed to save message:", err)
    }

    // Broadcast
    r.broadcast <- msg.Bytes()
}
```

### Goal 3: Scale with Redis Pub/Sub ⭐⭐⭐

Multiple server instances sharing messages via Redis.

**Hint**:
```go
func (r *Room) subscribeRedis() {
    pubsub := redisClient.Subscribe(ctx, "room:"+r.name)

    for msg := range pubsub.Channel() {
        // Broadcast to local clients
        r.broadcast <- []byte(msg.Payload)
    }
}

func (r *Room) publishMessage(msg []byte) {
    // Publish to Redis (reaches all servers)
    redisClient.Publish(ctx, "room:"+r.name, msg)
}
```

### Goal 4: Add Video/Voice Calling ⭐⭐⭐

Integrate WebRTC for peer-to-peer calls.

**Hint**:
```go
type SignalingMessage struct {
    Type   string  // "offer", "answer", "ice-candidate"
    From   string
    To     string
    Payload interface{}
}

func (h *Hub) handleSignaling(msg SignalingMessage) {
    // Forward signaling messages between peers
    recipient := h.clients[msg.To]
    data, _ := json.Marshal(msg)
    recipient.send <- data
}
```

### Goal 5: Add End-to-End Encryption ⭐⭐⭐

Encrypt messages client-to-client (Signal protocol).

**Hint**:
```go
// Client side (JavaScript):
const encrypted = await crypto.subtle.encrypt(
    { name: "AES-GCM", iv: iv },
    sharedKey,
    message
);

// Server just forwards encrypted data
// Cannot read message content (zero-knowledge)
```

---

## How to Run

```bash
# Install dependencies
cd /home/user/go-edu/minis/32-websocket-chatroom
go get github.com/gorilla/websocket

# Run the server
go run ./cmd/chatroom/main.go

# In browser console:
const ws = new WebSocket("ws://localhost:8080/ws?room=general&user=Alice");
ws.onmessage = (e) => console.log(e.data);
ws.send("Hello, world!");

# Or use websocat:
websocat ws://localhost:8080/ws?room=general&user=Bob
# Type messages and press Enter
```

---

## Summary

**What you learned**:
- ✅ WebSocket protocol: Handshake, frames, persistent connections
- ✅ gorilla/websocket: Upgrader, Conn, message types
- ✅ Broadcasting pattern: Hub-and-spoke architecture
- ✅ Connection management: Join, leave, disconnect, cleanup
- ✅ Concurrency safety: Channels, mutexes, goroutines per client
- ✅ Production patterns: Heartbeats, timeouts, graceful shutdown

**Why this matters**:
Real-time communication is essential for modern apps (chat, notifications, live updates). WebSockets enable instant bidirectional communication without polling overhead. Proper connection management ensures scalability and reliability.

**Key takeaway**:
WebSockets = Persistent bidirectional channels for real-time data

**Next steps**:
- Explore WebRTC for peer-to-peer video/audio
- Study operational transformation for collaborative editing
- Learn NATS/Redis for distributed messaging

Build in real-time! 🚀
