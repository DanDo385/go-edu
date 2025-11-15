# TCP Protocol Exercises

This document contains progressive exercises to deepen your understanding of TCP programming in Go.

## Exercise 1: Command Protocol ⭐⭐

Extend the echo server to support multiple commands instead of just echoing.

### Protocol Specification

```
Client → Server: <COMMAND> <args>\n
Server → Client: <response>\n

Commands:
  ECHO <text>    → ECHO: <text>
  UPPER <text>   → <TEXT IN UPPERCASE>
  LOWER <text>   → <text in lowercase>
  REVERSE <text> → <txet>
  TIME           → Current server time (RFC3339)
  PING           → PONG
  QUIT           → Goodbye! (then close connection)
```

### Example Session

```
C: ECHO Hello World
S: ECHO: Hello World
C: UPPER golang
S: GOLANG
C: REVERSE stressed
S: desserts
C: TIME
S: 2024-01-15T10:30:45Z
C: PING
S: PONG
C: QUIT
S: Goodbye!
[Connection closes]
```

### Implementation Tasks

1. Parse incoming messages to extract command and arguments
2. Implement command handlers
3. Return appropriate responses
4. Handle invalid commands with error messages
5. Implement QUIT command to gracefully close connection

### Starter Code

```go
// ParseCommand parses a line into command and arguments
func ParseCommand(line string) (command string, args string) {
	// TODO: implement
}

// HandleCommand processes a command and returns the response
func HandleCommand(command, args string) string {
	// TODO: implement command handlers
	switch command {
	case "ECHO":
		return "ECHO: " + args
	case "UPPER":
		// TODO
	case "LOWER":
		// TODO
	// ... more commands
	default:
		return "ERROR: Unknown command"
	}
}
```

### Test Cases

```go
func TestCommandProtocol(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ECHO test", "ECHO: test"},
		{"UPPER hello", "HELLO"},
		{"LOWER WORLD", "world"},
		{"REVERSE abc", "cba"},
		{"PING", "PONG"},
		{"INVALID", "ERROR: Unknown command"},
	}

	// TODO: implement tests
}
```

---

## Exercise 2: Connection Timeout ⭐⭐

Implement connection timeouts to prevent idle connections from consuming resources.

### Requirements

1. **Read timeout**: Close connection if no data received for 30 seconds
2. **Write timeout**: Fail write if it takes longer than 10 seconds
3. **Idle timeout**: Track last activity and close idle connections
4. **Graceful closure**: Send "Connection timeout" message before closing

### Implementation Hints

```go
// Set read deadline (resets with each read)
conn.SetReadDeadline(time.Now().Add(30 * time.Second))

// Set write deadline
conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

// Check for timeout errors
if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
	log.Println("Timeout occurred")
}
```

### Test Case

```go
func TestReadTimeout(t *testing.T) {
	listener, addr := startTestServer(t)
	defer listener.Close()

	conn, _ := net.Dial("tcp", addr)
	defer conn.Close()

	// Don't send anything for 35 seconds
	time.Sleep(35 * time.Second)

	// Connection should be closed by server
	_, err := conn.Write([]byte("test\n"))
	if err == nil {
		t.Error("Expected connection to be closed due to timeout")
	}
}
```

---

## Exercise 3: Binary Protocol ⭐⭐⭐

Implement a binary protocol instead of text-based protocol for better efficiency.

### Wire Format

```
[4 bytes: message length (uint32, big-endian)][message bytes]

Example:
0x00 0x00 0x00 0x05 'H' 'e' 'l' 'l' 'o'
^length = 5         ^5 bytes of data
```

### Implementation

```go
import "encoding/binary"

// SendBinaryMessage sends a length-prefixed binary message
func SendBinaryMessage(conn net.Conn, data []byte) error {
	// Write length prefix
	length := uint32(len(data))
	if err := binary.Write(conn, binary.BigEndian, length); err != nil {
		return err
	}

	// Write data
	_, err := conn.Write(data)
	return err
}

// ReceiveBinaryMessage receives a length-prefixed binary message
func ReceiveBinaryMessage(conn net.Conn) ([]byte, error) {
	// Read length prefix
	var length uint32
	if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
		return nil, err
	}

	// Validate length (prevent DoS)
	if length > 1024*1024 {  // Max 1MB
		return nil, fmt.Errorf("message too large: %d bytes", length)
	}

	// Read data
	data := make([]byte, length)
	if _, err := io.ReadFull(conn, data); err != nil {
		return nil, err
	}

	return data, nil
}
```

### Benefits Over Text Protocol

1. **No escaping needed**: Binary data can contain any bytes
2. **Efficient parsing**: No need to scan for delimiters
3. **Fixed overhead**: Always 4 bytes for length prefix
4. **Better for structured data**: Can encode Protocol Buffers, MessagePack, etc.

### Test Cases

```go
func TestBinaryProtocol(t *testing.T) {
	// Test various message sizes
	testCases := [][]byte{
		[]byte(""),
		[]byte("Hello"),
		[]byte(strings.Repeat("x", 1000)),
		[]byte{0x00, 0xFF, 0x01, 0x02},  // Binary data
	}

	for _, data := range testCases {
		// TODO: send and receive, verify data matches
	}
}
```

---

## Exercise 4: TLS/SSL Encryption ⭐⭐⭐

Add TLS encryption to secure the connection.

### Server Implementation

```go
import "crypto/tls"

func StartTLSServer(addr, certFile, keyFile string) error {
	// Load certificate and private key
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	// Configure TLS
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	// Create TLS listener
	listener, err := tls.Listen("tcp", addr, config)
	if err != nil {
		return err
	}
	defer listener.Close()

	// Accept loop (same as before)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		go handleClient(conn)
	}
}
```

### Client Implementation

```go
func TLSEchoClient(addr, message string, insecureSkipVerify bool) (string, error) {
	// Configure TLS
	config := &tls.Config{
		InsecureSkipVerify: insecureSkipVerify,  // For testing only!
	}

	// Connect with TLS
	conn, err := tls.Dial("tcp", addr, config)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// Send and receive (same as before)
	// ...
}
```

### Generate Test Certificate

```bash
# Generate private key
openssl genrsa -out server.key 2048

# Generate self-signed certificate
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365 \
  -subj "/C=US/ST=State/L=City/O=Organization/CN=localhost"
```

### Key Concepts

- **Certificate**: Proves server identity (signed by Certificate Authority)
- **Private key**: Used to decrypt data encrypted with certificate's public key
- **TLS handshake**: Negotiates encryption parameters, verifies certificate
- **Cipher suite**: Algorithm used for encryption (AES, ChaCha20, etc.)

---

## Exercise 5: Chat Server ⭐⭐⭐⭐

Build a multi-user chat server where messages are broadcast to all connected clients.

### Features

1. **User registration**: Clients send username on connect
2. **Broadcast**: Messages from one client sent to all others
3. **Join/leave notifications**: "Alice joined" / "Bob left"
4. **Private messages**: `/msg Bob Hello` sends to Bob only
5. **User list**: `/users` shows all connected users
6. **Nicknames**: `/nick NewName` changes username

### Protocol

```
Client → Server on connect: USERNAME <name>\n
Server → All:              JOINED <name>\n

Client → Server:           MSG <text>\n
Server → All (except sender): <name>: <text>\n

Client → Server:           PRIVMSG <recipient> <text>\n
Server → Recipient:        [PM from <sender>]: <text>\n

Client → Server:           USERS\n
Server → Client:           USERS <name1>,<name2>,<name3>\n
```

### Architecture

```go
type ChatServer struct {
	mu      sync.RWMutex
	clients map[string]net.Conn  // username -> connection
}

func (cs *ChatServer) Broadcast(message string, except string) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	for username, conn := range cs.clients {
		if username == except {
			continue
		}
		fmt.Fprintln(conn, message)
	}
}

func (cs *ChatServer) AddClient(username string, conn net.Conn) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if _, exists := cs.clients[username]; exists {
		return fmt.Errorf("username already taken")
	}

	cs.clients[username] = conn
	return nil
}
```

### Challenges

1. **Concurrency**: Multiple clients reading/writing simultaneously
2. **Synchronization**: Protect shared state (client map) with mutex
3. **Error handling**: What if sending to one client fails?
4. **Resource cleanup**: Remove disconnected clients from map

---

## Exercise 6: Connection Pooling ⭐⭐⭐

Implement a connection pool for efficient connection reuse.

### Why Connection Pooling?

Creating TCP connections is expensive:
1. Three-way handshake (1.5 round-trips)
2. TLS handshake if using encryption (2 additional round-trips)
3. TCP slow start (gradual ramp-up of throughput)

Reusing connections saves time and resources.

### Implementation

```go
type ConnectionPool struct {
	mu       sync.Mutex
	conns    []net.Conn
	addr     string
	maxConns int
}

func NewConnectionPool(addr string, maxConns int) *ConnectionPool {
	return &ConnectionPool{
		addr:     addr,
		maxConns: maxConns,
		conns:    make([]net.Conn, 0, maxConns),
	}
}

func (p *ConnectionPool) Get() (net.Conn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Try to reuse existing connection
	if len(p.conns) > 0 {
		conn := p.conns[len(p.conns)-1]
		p.conns = p.conns[:len(p.conns)-1]
		return conn, nil
	}

	// Create new connection
	return net.Dial("tcp", p.addr)
}

func (p *ConnectionPool) Put(conn net.Conn) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// If pool is full, close connection
	if len(p.conns) >= p.maxConns {
		return conn.Close()
	}

	// Return to pool
	p.conns = append(p.conns, conn)
	return nil
}

func (p *ConnectionPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, conn := range p.conns {
		conn.Close()
	}
	p.conns = nil
	return nil
}
```

### Usage

```go
pool := NewConnectionPool("localhost:8080", 10)
defer pool.Close()

// Get connection from pool
conn, err := pool.Get()
if err != nil {
	return err
}

// Use connection
SendMessage(conn, "Hello")

// Return to pool (don't close!)
pool.Put(conn)
```

### Enhancements

1. **Health checks**: Test connections before reuse
2. **Max idle time**: Close connections idle too long
3. **Connection lifecycle**: Track created/destroyed connections
4. **Metrics**: Monitor pool size, wait times, etc.

---

## Exercise 7: Load Balancer ⭐⭐⭐⭐

Build a TCP load balancer that distributes connections across multiple backend servers.

### Architecture

```
Client → Load Balancer → Backend Server 1
                      ↘ Backend Server 2
                      ↘ Backend Server 3
```

### Load Balancing Strategies

**1. Round Robin**: Cycle through servers
```go
type RoundRobin struct {
	mu       sync.Mutex
	backends []string
	current  int
}

func (rr *RoundRobin) Next() string {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	backend := rr.backends[rr.current]
	rr.current = (rr.current + 1) % len(rr.backends)
	return backend
}
```

**2. Least Connections**: Send to server with fewest active connections
```go
type LeastConn struct {
	mu    sync.Mutex
	conns map[string]int  // backend -> active connections
}

func (lc *LeastConn) Next() string {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	minConns := int(^uint(0) >> 1)  // Max int
	var backend string

	for addr, count := range lc.conns {
		if count < minConns {
			minConns = count
			backend = addr
		}
	}

	return backend
}
```

**3. Random**: Pick random server
```go
func (r *Random) Next() string {
	return r.backends[rand.Intn(len(r.backends))]
}
```

### Implementation

```go
func handleProxy(clientConn net.Conn, backendAddr string) {
	defer clientConn.Close()

	// Connect to backend
	backendConn, err := net.Dial("tcp", backendAddr)
	if err != nil {
		log.Printf("Failed to connect to backend: %v", err)
		return
	}
	defer backendConn.Close()

	// Bidirectional copy
	done := make(chan error, 2)

	go func() {
		_, err := io.Copy(backendConn, clientConn)
		done <- err
	}()

	go func() {
		_, err := io.Copy(clientConn, backendConn)
		done <- err
	}()

	// Wait for either direction to complete
	<-done
}
```

### Challenges

1. **Health checks**: Detect and skip unhealthy backends
2. **Graceful shutdown**: Don't interrupt active connections
3. **Metrics**: Track requests per backend
4. **Sticky sessions**: Route same client to same backend

---

## Exercise 8: Port Scanner ⭐⭐

Build a concurrent port scanner to check which ports are open on a host.

### Implementation

```go
func ScanPort(host string, port int, timeout time.Duration) bool {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func ScanPortsConcurrent(host string, ports []int, timeout time.Duration) []int {
	var mu sync.Mutex
	var openPorts []int
	var wg sync.WaitGroup

	// Limit concurrency with semaphore
	semaphore := make(chan struct{}, 100)

	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()

			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			if ScanPort(host, p, timeout) {
				mu.Lock()
				openPorts = append(openPorts, p)
				mu.Unlock()
			}
		}(port)
	}

	wg.Wait()
	sort.Ints(openPorts)
	return openPorts
}
```

### Usage

```go
// Scan common ports
commonPorts := []int{21, 22, 23, 25, 80, 443, 3306, 5432, 8080}
open := ScanPortsConcurrent("example.com", commonPorts, 2*time.Second)
fmt.Printf("Open ports: %v\n", open)

// Scan range of ports
ports := make([]int, 1000)
for i := range ports {
	ports[i] = i + 1
}
open = ScanPortsConcurrent("localhost", ports, 100*time.Millisecond)
```

### Enhancements

1. **Service detection**: Identify service running on port (HTTP, SSH, etc.)
2. **Banner grabbing**: Read initial server response
3. **Rate limiting**: Don't overwhelm target with too many connections
4. **Results formatting**: Pretty-print results

---

## Summary

These exercises progressively build your TCP programming skills:

1. **Command Protocol** → Protocol design and parsing
2. **Connection Timeout** → Resource management
3. **Binary Protocol** → Efficient serialization
4. **TLS Encryption** → Security
5. **Chat Server** → Multi-client coordination
6. **Connection Pooling** → Performance optimization
7. **Load Balancer** → Scalability
8. **Port Scanner** → Network analysis

Work through them in order, and you'll gain a deep understanding of TCP programming in Go!
