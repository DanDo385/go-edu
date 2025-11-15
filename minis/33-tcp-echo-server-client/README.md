# Project 33: TCP Echo Server/Client

## 1. What Is This About?

### Real-World Scenario

You're building a distributed system where services need to communicate over a network. HTTP is great for request-response patterns, but what if you need:
- **Persistent connections**: Keep a connection open for continuous data exchange
- **Low overhead**: Minimal protocol overhead (HTTP headers can be bulky)
- **Full-duplex communication**: Both sides can send data simultaneously
- **Custom protocols**: Design your own communication format

**Examples of TCP-based systems**:
- Redis: High-performance key-value store
- PostgreSQL: Database client-server communication
- SSH: Secure remote shell
- SMTP: Email transmission
- Custom microservice communication

This project teaches you the fundamentals of **TCP network programming** in Go by building an echo server and client that communicate using raw TCP sockets.

### What You'll Learn

1. **TCP fundamentals**: How TCP connections work at the socket level
2. **net package**: Go's standard library for network programming
3. **Connection handling**: Accept connections, read/write data, handle errors
4. **bufio on sockets**: Buffered I/O for efficient network communication
5. **Protocol design**: Create a simple line-based protocol
6. **Graceful shutdown**: Clean connection closure and resource cleanup
7. **Concurrent servers**: Handle multiple clients simultaneously

### The Challenge

Build a TCP echo server that:
- Listens on a specified port
- Accepts multiple concurrent client connections
- Reads lines of text from clients
- Echoes each line back to the client
- Handles connection errors gracefully
- Shuts down cleanly on signal

Build a TCP client that:
- Connects to the server
- Sends user input to the server
- Receives and displays server responses
- Handles network errors

---

## 2. First Principles: What Is TCP?

### The Network Stack (Bottom-Up)

Computer networks are organized in layers. Each layer builds on the layer below:

**Layer 1: Physical** - Electrical signals over wires
```
01001010 11010101 → Voltage changes on a wire
```

**Layer 2: Data Link** - Frames between directly connected devices
```
[Ethernet Frame: Source MAC | Dest MAC | Data | Checksum]
```

**Layer 3: Network (IP)** - Packets routed between networks
```
[IP Packet: Source IP | Dest IP | Data]
Example: 192.168.1.100 → 93.184.216.34 (example.com)
```

**Layer 4: Transport (TCP/UDP)** - End-to-end communication between processes
```
[TCP Segment: Source Port | Dest Port | Sequence# | Data]
Example: localhost:3000 → localhost:8080
```

**Layer 5-7: Application** - High-level protocols (HTTP, SSH, etc.)
```
GET /index.html HTTP/1.1
```

**Where TCP fits**: Layer 4 (Transport). It provides reliable, ordered, error-checked delivery of a stream of bytes between applications.

### TCP vs UDP: The Fundamental Difference

**TCP (Transmission Control Protocol)**: Reliable, ordered, connection-oriented
```
Client: "Hello"
Server: "I received 'Hello'" (acknowledgment)
Client: "World"
Server: "I received 'World'"
```
- **Reliable**: Guarantees delivery (retransmits lost packets)
- **Ordered**: Data arrives in the same order it was sent
- **Connection-oriented**: Establishes a connection before data transfer
- **Flow control**: Prevents overwhelming the receiver
- **Use when**: You need reliability (file transfer, HTTP, database queries)

**UDP (User Datagram Protocol)**: Unreliable, connectionless, lightweight
```
Client: "Hello" → (packet might be lost)
Client: "World" → (might arrive before "Hello")
```
- **Unreliable**: No delivery guarantees (fire and forget)
- **Unordered**: Packets may arrive out of order
- **Connectionless**: No connection setup overhead
- **Low latency**: No waiting for acknowledgments
- **Use when**: You need speed over reliability (video streaming, gaming, DNS)

**Example - Why TCP for chat, UDP for video calls**:
- **Chat**: Missing a message is bad → use TCP
- **Video call**: Old video frames are useless, better to drop them → use UDP

### How TCP Works: The Three-Way Handshake

Before data can be sent, TCP establishes a connection through a "three-way handshake":

```
Client                                Server
   |                                     |
   |  1. SYN (Synchronize)              |
   |  "Let's establish a connection"    |
   | ----------------------------------> |
   |                                     |
   |  2. SYN-ACK (Synchronize-Acknowledge)
   |  "OK, I acknowledge your request"  |
   | <---------------------------------- |
   |                                     |
   |  3. ACK (Acknowledge)              |
   |  "Great, connection established"   |
   | ----------------------------------> |
   |                                     |
   |  [Connection established]          |
   |  Now data can flow both ways       |
   |                                     |
```

**Step-by-step**:
1. **Client → Server (SYN)**: "I want to connect. My initial sequence number is X."
2. **Server → Client (SYN-ACK)**: "OK! My initial sequence number is Y. I got your X."
3. **Client → Server (ACK)**: "Got it! We're connected."

**Why three steps?**
- Synchronizes sequence numbers (needed for ordering packets)
- Confirms both sides are ready to communicate
- Prevents old duplicate packets from interfering

**In Go**:
```go
// Client side
conn, err := net.Dial("tcp", "localhost:8080")
// Three-way handshake happens automatically!

// Server side
listener, _ := net.Listen("tcp", ":8080")
conn, err := listener.Accept()
// Handshake is already complete when Accept() returns
```

### Sockets: The Programming Interface

A **socket** is an endpoint for network communication. Think of it as a "mailbox" for sending/receiving data.

**Socket address** = IP address + port number
```
192.168.1.100:8080
  ^IP address  ^port
```

**Port numbers**:
- Range: 0-65535
- **Well-known ports** (0-1023): Reserved for standard services
  - 80: HTTP
  - 443: HTTPS
  - 22: SSH
  - 25: SMTP (email)
- **Registered ports** (1024-49151): Assigned by IANA for specific applications
- **Dynamic/Private ports** (49152-65535): Temporary ports for client connections

**Server socket workflow**:
```
1. Create socket        → socket()
2. Bind to address      → bind(0.0.0.0:8080)
3. Listen for connections → listen()
4. Accept connection    → accept() → new socket for this client
5. Read/Write data      → read()/write()
6. Close connection     → close()
```

**Client socket workflow**:
```
1. Create socket        → socket()
2. Connect to server    → connect(192.168.1.100:8080)
3. Read/Write data      → read()/write()
4. Close connection     → close()
```

**In Go, this is abstracted**:
```go
// Server
listener, _ := net.Listen("tcp", ":8080")  // socket() + bind() + listen()
conn, _ := listener.Accept()                // accept()
conn.Write([]byte("Hello"))                 // write()
conn.Close()                                // close()

// Client
conn, _ := net.Dial("tcp", "localhost:8080") // socket() + connect()
buf := make([]byte, 1024)
conn.Read(buf)                                // read()
conn.Close()                                  // close()
```

### TCP Data Stream: Byte Stream, Not Messages

**Critical insight**: TCP provides a **byte stream**, not a message stream.

**What this means**:
```
Client sends: "Hello" then "World"

❌ WRONG assumption:
Server receives: "Hello" (first read), "World" (second read)

✅ REALITY (any of these is possible):
Server receives: "HelloWorld" (one read)
Server receives: "Hel" (first read), "loWorld" (second read)
Server receives: "H" (first), "elloWor" (second), "ld" (third)
```

**Why?** TCP doesn't preserve message boundaries. It:
- Breaks data into segments (based on network conditions)
- May combine multiple writes into one segment
- May split one write into multiple segments
- Delivers a continuous stream of bytes in order

**Solution: Protocol framing**
You must define how messages are delimited:

**Option 1: Fixed-length messages**
```
Every message is exactly 10 bytes
"Hello     " "World     "
```

**Option 2: Length prefix**
```
[4 bytes: length][message]
0005Hello0005World
```

**Option 3: Delimiter (what we'll use)**
```
Messages end with newline
"Hello\n" "World\n"
```

**Option 4: Self-describing (like HTTP)**
```
GET / HTTP/1.1\r\n
Content-Length: 13\r\n
\r\n
Hello, World!
```

---

## 3. Breaking Down the Solution

### Part 1: TCP Server Architecture

**High-level flow**:
```
1. Create listener socket on port 8080
2. Loop forever:
   a. Accept new client connection (blocks until client connects)
   b. Spawn goroutine to handle this client
   c. Go back to step 2 (accept next client)

3. Per-client goroutine:
   a. Read line from client
   b. Echo line back to client
   c. Repeat until client disconnects or error
   d. Close connection
```

**Code structure**:
```go
func main() {
    // 1. Create listener
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()

    // 2. Accept loop
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println("Accept error:", err)
            continue
        }

        // 3. Handle client concurrently
        go handleClient(conn)
    }
}

func handleClient(conn net.Conn) {
    defer conn.Close()

    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        line := scanner.Text()
        fmt.Fprintf(conn, "ECHO: %s\n", line)
    }

    if err := scanner.Err(); err != nil {
        log.Println("Read error:", err)
    }
}
```

### Part 2: Understanding net.Listen

```go
listener, err := net.Listen("tcp", ":8080")
```

**What this does**:
1. Creates a socket
2. Binds it to all network interfaces (`:8080` means `0.0.0.0:8080`)
3. Marks it as a listening socket
4. Returns a `net.Listener` interface

**Bind address options**:
```go
net.Listen("tcp", ":8080")              // All interfaces: 0.0.0.0:8080
net.Listen("tcp", "127.0.0.1:8080")     // Localhost only
net.Listen("tcp", "192.168.1.100:8080") // Specific interface
net.Listen("tcp", "localhost:8080")     // Resolves to 127.0.0.1:8080
```

**Port 0 (OS chooses port)**:
```go
listener, _ := net.Listen("tcp", ":0")
addr := listener.Addr().(*net.TCPAddr)
fmt.Println("Listening on port:", addr.Port) // e.g., 54321
```

**Common errors**:
- **`bind: address already in use`**: Another process is using port 8080
- **`bind: permission denied`**: Ports < 1024 require root/admin on Unix

### Part 3: Understanding listener.Accept()

```go
conn, err := listener.Accept()
```

**What this does**:
1. **Blocks** until a client initiates a connection
2. Completes the three-way handshake
3. Returns a `net.Conn` representing the established connection
4. **Each connection gets its own `net.Conn` object**

**Why it blocks**: There's nothing to do until a client connects. The OS wakes up the thread when a client arrives.

**Handling errors**:
```go
for {
    conn, err := listener.Accept()
    if err != nil {
        // Temporary network errors (recoverable)
        if ne, ok := err.(net.Error); ok && ne.Temporary() {
            log.Printf("Temporary accept error: %v; retrying", err)
            time.Sleep(10 * time.Millisecond)
            continue
        }
        // Permanent error (listener closed, fatal error)
        log.Fatal("Accept error:", err)
    }

    go handleClient(conn)
}
```

### Part 4: Reading Data with bufio.Scanner

**Problem**: TCP is a byte stream. We need to read line-by-line.

**Solution 1: Read into buffer (manual)**
```go
buf := make([]byte, 4096)
n, err := conn.Read(buf)
// buf now contains n bytes
// But where does one line end? Manual parsing needed!
```

**Solution 2: bufio.Scanner (recommended)**
```go
scanner := bufio.NewScanner(conn)
for scanner.Scan() {
    line := scanner.Text()  // One complete line (newline removed)
    fmt.Println("Received:", line)
}

if err := scanner.Err(); err != nil {
    log.Println("Error:", err)
}
```

**How bufio.Scanner works**:
1. Maintains an internal buffer
2. Reads chunks from the network into the buffer
3. Scans for delimiter (default: newline `\n`)
4. Returns complete lines via `Text()`
5. Handles partial lines (waits for more data)

**What `scanner.Scan()` returns**:
- **`true`**: A line was successfully read (access via `scanner.Text()`)
- **`false`**: No more data (check `scanner.Err()` for errors vs EOF)

**Custom delimiters**:
```go
scanner := bufio.NewScanner(conn)
scanner.Split(bufio.ScanWords)  // Split by whitespace instead of lines

// Custom split function
scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
    // Find delimiter, return token
})
```

### Part 5: Writing Data to Connection

**Writing strings**:
```go
fmt.Fprintf(conn, "ECHO: %s\n", line)
```

**Why `fmt.Fprintf` works**: `net.Conn` implements `io.Writer`, so any function that writes to an `io.Writer` works:
```go
io.WriteString(conn, "Hello\n")
conn.Write([]byte("Hello\n"))
fmt.Fprintln(conn, "Hello")
```

**Buffered writing (more efficient)**:
```go
writer := bufio.NewWriter(conn)
fmt.Fprintf(writer, "Line 1\n")
fmt.Fprintf(writer, "Line 2\n")
writer.Flush()  // Send buffered data to network
```

**When to flush**:
- After complete message (so client receives it immediately)
- Before reading (prevent deadlock if both sides are waiting)
- Periodically (if buffering for efficiency)

**Deadlock example**:
```go
// Server
scanner := bufio.NewScanner(conn)
writer := bufio.NewWriter(conn)

scanner.Scan()
writer.WriteString("Response\n")
// ❌ BUG: No Flush! Client never receives response!
```

### Part 6: Closing Connections

**Always close connections**:
```go
defer conn.Close()
```

**What happens when you close**:
1. Sends FIN (finish) packet to peer
2. No more data can be sent
3. Can still receive data until peer closes
4. Releases OS resources (file descriptors)

**Graceful vs abrupt close**:
```go
// Graceful close
conn.Close()
// Sends FIN, waits for peer to acknowledge

// Abrupt close (force)
if tcpConn, ok := conn.(*net.TCPConn); ok {
    tcpConn.SetLinger(0)  // RST instead of FIN
    tcpConn.Close()
}
```

**Connection states after close**:
```
Before close: ESTABLISHED
After Close(): FIN_WAIT_1 → FIN_WAIT_2 → TIME_WAIT → CLOSED
```

### Part 7: TCP Client Implementation

```go
func main() {
    // 1. Connect to server
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        log.Fatal("Dial error:", err)
    }
    defer conn.Close()

    // 2. Read user input
    scanner := bufio.NewScanner(os.Stdin)
    writer := bufio.NewWriter(conn)
    reader := bufio.NewReader(conn)

    for {
        fmt.Print("> ")

        // Read from stdin
        if !scanner.Scan() {
            break
        }
        line := scanner.Text()

        // Send to server
        fmt.Fprintf(writer, "%s\n", line)
        writer.Flush()

        // Read response
        response, err := reader.ReadString('\n')
        if err != nil {
            log.Fatal("Read error:", err)
        }

        fmt.Print("< " + response)
    }
}
```

**Key points**:
- Use `bufio.Writer` for efficient writes (buffer multiple small writes)
- **Always Flush()** after writing a complete message
- Use `bufio.Reader` for efficient reads
- `ReadString('\n')` reads until delimiter (includes delimiter in result)

### Part 8: Concurrent Client Handling

**Why goroutines for each client?**

Without goroutines:
```go
for {
    conn, _ := listener.Accept()
    handleClient(conn)  // ❌ Blocks! Other clients must wait!
}
```

With goroutines:
```go
for {
    conn, _ := listener.Accept()
    go handleClient(conn)  // ✅ Non-blocking! Handles clients concurrently!
}
```

**Scalability**:
- Each goroutine costs ~2-4 KB of stack memory
- Can handle thousands of concurrent connections
- OS limits: File descriptors (typically 1024-65535)

**Connection state tracking**:
```go
var (
    mu      sync.Mutex
    clients = make(map[net.Conn]bool)
)

func addClient(conn net.Conn) {
    mu.Lock()
    clients[conn] = true
    mu.Unlock()
}

func removeClient(conn net.Conn) {
    mu.Lock()
    delete(clients, conn)
    mu.Unlock()
}

func handleClient(conn net.Conn) {
    addClient(conn)
    defer removeClient(conn)
    defer conn.Close()

    // ... handle client ...
}
```

---

## 4. Protocol Design: Line-Based Echo Protocol

### Our Simple Protocol Specification

**Protocol**: Line-based echo protocol
**Framing**: Lines delimited by `\n` (newline)
**Encoding**: UTF-8 text

**Message format**:
```
Client → Server: <text>\n
Server → Client: ECHO: <text>\n
```

**Example session**:
```
C: Hello\n
S: ECHO: Hello\n
C: How are you?\n
S: ECHO: How are you?\n
C: Goodbye\n
S: ECHO: Goodbye\n
[Client closes connection]
```

**Protocol rules**:
1. Each message ends with `\n`
2. Server echoes back with "ECHO: " prefix
3. Maximum line length: 64 KB (bufio.Scanner default)
4. Connection closes when client disconnects

### Protocol Evolution: Adding Features

**Feature 1: Commands**
```
Client → Server: ECHO <text>\n
Client → Server: UPPER <text>\n
Client → Server: LOWER <text>\n
Client → Server: QUIT\n

Server → Client: <result>\n
```

**Feature 2: Error handling**
```
Client → Server: INVALID\n
Server → Client: ERROR: Unknown command\n
```

**Feature 3: Binary protocol**
Instead of text, use binary framing:
```
[4 bytes: message length][message bytes]
```

**Feature 4: Multiplexing**
Multiple logical streams over one connection:
```
[4 bytes: stream ID][4 bytes: length][data]
```

### Real-World Protocol Examples

**Redis protocol (RESP)**:
```
*2\r\n
$3\r\n
GET\r\n
$3\r\n
key\r\n
```
- Lines end with `\r\n`
- `*N` = array of N elements
- `$N` = bulk string of N bytes

**HTTP/1.1**:
```
GET /index.html HTTP/1.1\r\n
Host: example.com\r\n
\r\n
```
- Text-based protocol
- Headers end with `\r\n\r\n`
- Body length specified by `Content-Length` header

**PostgreSQL wire protocol**:
```
[1 byte: message type]['Q'][4 bytes: length][SQL query]
```
- Binary protocol
- Type byte indicates message kind
- Length-prefixed messages

---

## 5. Complete Solution Walkthrough

### Server Implementation

```go
package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os"
    "os/signal"
    "sync"
    "syscall"
)

func main() {
    // Configuration
    addr := ":8080"
    if len(os.Args) > 1 {
        addr = os.Args[1]
    }

    // Create listener
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        log.Fatal("Listen error:", err)
    }

    log.Printf("TCP echo server listening on %s", listener.Addr())

    // Track active connections
    var wg sync.WaitGroup
    done := make(chan struct{})

    // Handle graceful shutdown
    go func() {
        sigCh := make(chan os.Signal, 1)
        signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
        <-sigCh

        log.Println("Shutting down server...")
        close(done)
        listener.Close()
    }()

    // Accept loop
    for {
        conn, err := listener.Accept()
        if err != nil {
            select {
            case <-done:
                // Shutdown initiated
                wg.Wait()
                log.Println("Server stopped")
                return
            default:
                log.Println("Accept error:", err)
                continue
            }
        }

        wg.Add(1)
        go handleClient(conn, &wg, done)
    }
}

func handleClient(conn net.Conn, wg *sync.WaitGroup, done <-chan struct{}) {
    defer wg.Done()
    defer conn.Close()

    clientAddr := conn.RemoteAddr().String()
    log.Printf("Client connected: %s", clientAddr)
    defer log.Printf("Client disconnected: %s", clientAddr)

    scanner := bufio.NewScanner(conn)
    writer := bufio.NewWriter(conn)

    for scanner.Scan() {
        select {
        case <-done:
            // Server shutting down
            return
        default:
        }

        line := scanner.Text()
        log.Printf("[%s] Received: %q", clientAddr, line)

        // Echo back
        response := fmt.Sprintf("ECHO: %s\n", line)
        if _, err := writer.WriteString(response); err != nil {
            log.Printf("Write error: %v", err)
            return
        }

        if err := writer.Flush(); err != nil {
            log.Printf("Flush error: %v", err)
            return
        }
    }

    if err := scanner.Err(); err != nil {
        log.Printf("Read error: %v", err)
    }
}
```

**Key features**:
1. **Configurable address**: Pass address as command-line argument
2. **Logging**: Log connections, disconnections, and messages
3. **Graceful shutdown**: Wait for active connections to complete
4. **Error handling**: Handle temporary vs permanent errors
5. **Concurrent clients**: Each client in its own goroutine

### Client Implementation

```go
package main

import (
    "bufio"
    "fmt"
    "io"
    "log"
    "net"
    "os"
)

func main() {
    addr := "localhost:8080"
    if len(os.Args) > 1 {
        addr = os.Args[1]
    }

    // Connect to server
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        log.Fatal("Dial error:", err)
    }
    defer conn.Close()

    log.Printf("Connected to %s", conn.RemoteAddr())

    // Channel for errors from reader goroutine
    errCh := make(chan error, 1)

    // Read server responses concurrently
    go func() {
        reader := bufio.NewReader(conn)
        for {
            line, err := reader.ReadString('\n')
            if err != nil {
                if err != io.EOF {
                    errCh <- err
                }
                return
            }
            fmt.Print("< " + line)
        }
    }()

    // Read user input and send to server
    scanner := bufio.NewScanner(os.Stdin)
    writer := bufio.NewWriter(conn)

    for {
        fmt.Print("> ")

        // Check for errors from reader goroutine
        select {
        case err := <-errCh:
            log.Fatal("Read error:", err)
        default:
        }

        // Read from stdin
        if !scanner.Scan() {
            break
        }

        line := scanner.Text()

        // Send to server
        if _, err := fmt.Fprintf(writer, "%s\n", line); err != nil {
            log.Fatal("Write error:", err)
        }

        if err := writer.Flush(); err != nil {
            log.Fatal("Flush error:", err)
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatal("Input error:", err)
    }
}
```

**Key features**:
1. **Concurrent read/write**: Separate goroutine for reading server responses
2. **Interactive**: Read from stdin, prompt user with `>`
3. **Error handling**: Check for errors from both reading and writing
4. **Buffered I/O**: Use bufio for efficiency

---

## 6. Key Concepts Explained

### Concept 1: TCP State Machine

TCP connections go through various states:

```
CLOSED → LISTEN → SYN_SENT → SYN_RECEIVED → ESTABLISHED
         (server)  (client)
```

**Active (client-initiated) connection**:
```
CLOSED → SYN_SENT → ESTABLISHED → FIN_WAIT_1 → FIN_WAIT_2 → TIME_WAIT → CLOSED
```

**Passive (server-side) connection**:
```
CLOSED → LISTEN → SYN_RECEIVED → ESTABLISHED → CLOSE_WAIT → LAST_ACK → CLOSED
```

**TIME_WAIT state**:
- Lasts 2 × MSL (Maximum Segment Lifetime) ≈ 60-240 seconds
- Prevents port reuse immediately after close
- Can cause "address already in use" errors

**Solution for rapid restarts**:
```go
listener, err := net.Listen("tcp", ":8080")
// Enable SO_REUSEADDR (Go does this automatically)
```

### Concept 2: TCP Flow Control

**Problem**: Sender can overwhelm receiver with data

**Solution**: Sliding window protocol
```
Receiver advertises: "I have 64KB of buffer space"
Sender can send up to 64KB before waiting for acknowledgment
As receiver processes data, it advertises more space
```

**In Go, this is automatic**. But you can observe it:
```go
if tcpConn, ok := conn.(*net.TCPConn); ok {
    tcpConn.SetReadBuffer(256 * 1024)  // 256KB receive buffer
    tcpConn.SetWriteBuffer(256 * 1024) // 256KB send buffer
}
```

### Concept 3: Nagle's Algorithm

**Problem**: Sending many small packets is inefficient
```
Send "H" → [20 bytes TCP header + 1 byte data] = 21 bytes
Send "e" → [20 bytes TCP header + 1 byte data] = 21 bytes
Send "l" → [20 bytes TCP header + 1 byte data] = 21 bytes
Total: 63 bytes sent for 3 bytes of data (95% overhead!)
```

**Nagle's algorithm**: Buffer small writes until:
- Enough data to fill a packet (MSS ≈ 1460 bytes)
- OR acknowledgment received for previous packet

**Trade-off**: Reduces overhead but increases latency

**Disable for low-latency apps**:
```go
if tcpConn, ok := conn.(*net.TCPConn); ok {
    tcpConn.SetNoDelay(true)  // Disable Nagle's algorithm
}
```

**When to disable**:
- Interactive protocols (SSH, telnet, gaming)
- Request-response patterns (small messages)

**When to keep enabled**:
- Bulk data transfer (file upload/download)
- Many small writes that can be batched

### Concept 4: Connection Pooling

**Problem**: Creating new TCP connections is expensive
- Three-way handshake
- TLS handshake (if using encryption)
- Slow start (TCP congestion control)

**Solution**: Reuse connections (connection pooling)

```go
type Pool struct {
    mu    sync.Mutex
    conns []net.Conn
    addr  string
}

func (p *Pool) Get() (net.Conn, error) {
    p.mu.Lock()
    defer p.mu.Unlock()

    if len(p.conns) > 0 {
        conn := p.conns[len(p.conns)-1]
        p.conns = p.conns[:len(p.conns)-1]
        return conn, nil
    }

    return net.Dial("tcp", p.addr)
}

func (p *Pool) Put(conn net.Conn) {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.conns = append(p.conns, conn)
}
```

**Used by**: HTTP clients, database drivers, Redis clients

### Concept 5: Timeouts

**Read timeout**: Maximum time to wait for data
```go
conn.SetReadDeadline(time.Now().Add(30 * time.Second))
```

**Write timeout**: Maximum time to wait for write to complete
```go
conn.SetWriteDeadline(time.Now().Add(30 * time.Second))
```

**Overall timeout**: Both read and write
```go
conn.SetDeadline(time.Now().Add(30 * time.Second))
```

**Clear timeout**:
```go
conn.SetDeadline(time.Time{})  // Zero value = no timeout
```

**Per-operation timeout**:
```go
for {
    conn.SetReadDeadline(time.Now().Add(30 * time.Second))
    data, err := reader.ReadString('\n')
    // ...
}
```

### Concept 6: Half-Close

**Half-close**: Close write side but keep read side open

```go
if tcpConn, ok := conn.(*net.TCPConn); ok {
    tcpConn.CloseWrite()  // Send FIN, but can still read
}

// Read remaining data from peer
io.Copy(os.Stdout, conn)

conn.Close()  // Close read side
```

**Use case**: HTTP client
1. Send request
2. Close write side (signal end of request)
3. Read response
4. Close connection

---

## 7. Real-World Applications

### Database Client-Server (PostgreSQL)

PostgreSQL uses a custom binary protocol over TCP:
```go
// Simplified PostgreSQL connection
conn, _ := net.Dial("tcp", "localhost:5432")

// Send startup message
msg := []byte{'Q', 0, 0, 0, 20, 'S', 'E', 'L', 'E', 'C', 'T', ' ', '1', ';', 0}
conn.Write(msg)

// Read response
buf := make([]byte, 4096)
n, _ := conn.Read(buf)
// Parse binary protocol...
```

### Redis Client

Redis uses RESP (REdis Serialization Protocol), a text-based protocol:
```go
conn, _ := net.Dial("tcp", "localhost:6379")

// Send command
fmt.Fprintf(conn, "*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n")

// Read response
reader := bufio.NewReader(conn)
line, _ := reader.ReadString('\n')  // e.g., "$5\r\n"
// Parse response...
```

### Load Balancer

Proxy TCP connections to backend servers:
```go
func handleConnection(clientConn net.Conn) {
    // Connect to backend
    backendConn, err := net.Dial("tcp", "backend:8080")
    if err != nil {
        clientConn.Close()
        return
    }

    // Bidirectional copy
    go io.Copy(backendConn, clientConn)  // Client → Backend
    io.Copy(clientConn, backendConn)     // Backend → Client

    clientConn.Close()
    backendConn.Close()
}
```

### Chat Server

Broadcast messages to all connected clients:
```go
var (
    mu      sync.Mutex
    clients = make(map[net.Conn]bool)
)

func broadcast(message string) {
    mu.Lock()
    defer mu.Unlock()

    for conn := range clients {
        fmt.Fprintln(conn, message)
    }
}

func handleClient(conn net.Conn) {
    mu.Lock()
    clients[conn] = true
    mu.Unlock()

    defer func() {
        mu.Lock()
        delete(clients, conn)
        mu.Unlock()
        conn.Close()
    }()

    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        broadcast(scanner.Text())
    }
}
```

### Port Scanner

Check which ports are open on a host:
```go
func scanPort(host string, port int, timeout time.Duration) bool {
    addr := fmt.Sprintf("%s:%d", host, port)
    conn, err := net.DialTimeout("tcp", addr, timeout)
    if err != nil {
        return false
    }
    conn.Close()
    return true
}

// Scan ports 1-1024
for port := 1; port <= 1024; port++ {
    if scanPort("example.com", port, 1*time.Second) {
        fmt.Printf("Port %d is open\n", port)
    }
}
```

---

## 8. Common Mistakes to Avoid

### Mistake 1: Not Closing Connections

**❌ Wrong**:
```go
conn, _ := net.Dial("tcp", "localhost:8080")
// Forgot to close!
```

**Problem**: File descriptor leak. Eventually runs out of file descriptors.

**✅ Correct**:
```go
conn, err := net.Dial("tcp", "localhost:8080")
if err != nil {
    return err
}
defer conn.Close()
```

### Mistake 2: Assuming Complete Reads

**❌ Wrong**:
```go
buf := make([]byte, 1024)
n, _ := conn.Read(buf)
// Assume buf[:n] contains complete message
```

**Problem**: TCP is a stream. One Read() might return partial data.

**✅ Correct**:
```go
// For fixed-length messages
io.ReadFull(conn, buf)

// For delimited messages
scanner := bufio.NewScanner(conn)
scanner.Scan()
line := scanner.Text()
```

### Mistake 3: Forgetting to Flush Buffered Writers

**❌ Wrong**:
```go
writer := bufio.NewWriter(conn)
fmt.Fprintln(writer, "Hello")
// Data stuck in buffer!
```

**✅ Correct**:
```go
writer := bufio.NewWriter(conn)
fmt.Fprintln(writer, "Hello")
writer.Flush()  // Send data immediately
```

### Mistake 4: Blocking Accept in Main Goroutine

**❌ Wrong**:
```go
conn1, _ := listener.Accept()
handleClient(conn1)  // Blocks! Other clients can't connect!
conn2, _ := listener.Accept()
```

**✅ Correct**:
```go
for {
    conn, _ := listener.Accept()
    go handleClient(conn)  // Non-blocking
}
```

### Mistake 5: Ignoring Errors

**❌ Wrong**:
```go
conn.Write([]byte("Hello"))
// What if write failed?
```

**✅ Correct**:
```go
if _, err := conn.Write([]byte("Hello")); err != nil {
    log.Printf("Write error: %v", err)
    return err
}
```

### Mistake 6: Not Setting Timeouts

**❌ Wrong**:
```go
scanner := bufio.NewScanner(conn)
scanner.Scan()  // Blocks forever if client doesn't send data!
```

**✅ Correct**:
```go
conn.SetReadDeadline(time.Now().Add(30 * time.Second))
scanner := bufio.NewScanner(conn)
if !scanner.Scan() {
    if err := scanner.Err(); err != nil {
        // Handle timeout or other error
    }
}
```

### Mistake 7: Reusing Buffers Incorrectly

**❌ Wrong**:
```go
buf := make([]byte, 1024)
for {
    n, _ := conn.Read(buf)
    go process(buf[:n])  // BUG: buf is reused!
}
```

**✅ Correct**:
```go
for {
    buf := make([]byte, 1024)  // New buffer each iteration
    n, _ := conn.Read(buf)
    go process(buf[:n])
}

// Or copy data:
for {
    buf := make([]byte, 1024)
    n, _ := conn.Read(buf)
    data := make([]byte, n)
    copy(data, buf[:n])
    go process(data)
}
```

---

## 9. Stretch Goals

### Goal 1: Add Protocol Commands ⭐⭐

Extend the protocol to support multiple commands:
```
ECHO <text>    → ECHO: <text>
UPPER <text>   → <TEXT IN UPPERCASE>
LOWER <text>   → <text in lowercase>
REVERSE <text> → <txet>
TIME           → Current server time
QUIT           → Close connection
```

### Goal 2: Implement Connection Timeout ⭐⭐

Close connections that are idle for more than N seconds:
```go
conn.SetReadDeadline(time.Now().Add(30 * time.Second))
```

### Goal 3: Add TLS/SSL Encryption ⭐⭐⭐

Secure the connection with TLS:
```go
// Server
cert, _ := tls.LoadX509KeyPair("server.crt", "server.key")
config := &tls.Config{Certificates: []tls.Certificate{cert}}
listener, _ := tls.Listen("tcp", ":8080", config)

// Client
config := &tls.Config{InsecureSkipVerify: true}  // For testing
conn, _ := tls.Dial("tcp", "localhost:8080", config)
```

### Goal 4: Implement Binary Protocol ⭐⭐⭐

Use length-prefixed binary messages instead of newline-delimited text:
```go
// Wire format: [4 bytes: length][message bytes]

// Write
msg := []byte("Hello")
length := uint32(len(msg))
binary.Write(conn, binary.BigEndian, length)
conn.Write(msg)

// Read
var length uint32
binary.Read(conn, binary.BigEndian, &length)
msg := make([]byte, length)
io.ReadFull(conn, msg)
```

### Goal 5: Build a Chat Server ⭐⭐⭐

Multiple clients can send messages that are broadcast to all:
- Assign usernames to clients
- Broadcast messages to all connected clients
- Show join/leave notifications
- Implement private messages

### Goal 6: Add Metrics and Monitoring ⭐⭐

Track server statistics:
- Total connections
- Active connections
- Bytes sent/received
- Messages processed
- Average response time

### Goal 7: Implement Connection Pooling ⭐⭐⭐

Client-side connection pool for reusing connections:
- Pool of connections to server
- Acquire/release connections
- Health checks for stale connections
- Maximum pool size

---

## How to Run

```bash
# Run the server
cd /home/user/go-edu/minis/33-tcp-echo-server-client
go run cmd/tcp-server/main.go

# In another terminal, run the client
go run cmd/tcp-client/main.go

# Or with custom address
go run cmd/tcp-server/main.go :9000
go run cmd/tcp-client/main.go localhost:9000

# Test with netcat
echo "Hello" | nc localhost 8080

# Run tests
go test ./exercise/...

# Run with verbose output
go test -v ./exercise/...
```

---

## Summary

**What you learned**:
- ✅ TCP fundamentals: Three-way handshake, byte streams, connection states
- ✅ net package: Listen, Accept, Dial, Read, Write
- ✅ bufio: Efficient buffered I/O for network connections
- ✅ Protocol design: Line-based framing with delimiters
- ✅ Concurrency: Handle multiple clients with goroutines
- ✅ Error handling: Graceful connection closure and error recovery
- ✅ Real-world patterns: Connection pooling, timeouts, graceful shutdown

**Why this matters**:
TCP is the foundation of most network communication. Understanding TCP programming is essential for:
- Building custom network protocols
- Debugging network issues
- Understanding how HTTP, databases, and RPC work under the hood
- Building high-performance network services

**Key takeaways**:
- TCP provides a reliable byte stream, not a message stream
- Always frame your protocol (delimiters, length-prefix, etc.)
- Use buffered I/O for efficiency
- Handle each connection in its own goroutine
- Always close connections and handle errors
- Set timeouts to prevent resource leaks

**Next steps**:
- Project 34: Build a rate limiter with token bucket algorithm
- Project 35: Implement JWT authentication middleware
- Explore: Build your own Redis-like database with TCP

Happy networking! 🌐
