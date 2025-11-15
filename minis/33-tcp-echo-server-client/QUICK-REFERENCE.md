# TCP Programming Quick Reference

## Go net Package Cheat Sheet

### Creating a TCP Server

```go
// Listen on a TCP port
listener, err := net.Listen("tcp", ":8080")
if err != nil {
    log.Fatal(err)
}
defer listener.Close()

// Accept connections in a loop
for {
    conn, err := listener.Accept()
    if err != nil {
        log.Println("Accept error:", err)
        continue
    }

    // Handle connection (usually in goroutine)
    go handleConnection(conn)
}
```

### Creating a TCP Client

```go
// Connect to a TCP server
conn, err := net.Dial("tcp", "localhost:8080")
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

// Use the connection
conn.Write([]byte("Hello"))
```

### Reading Data

```go
// Read into buffer
buf := make([]byte, 4096)
n, err := conn.Read(buf)
if err != nil {
    if err == io.EOF {
        // Connection closed
    }
    return err
}
data := buf[:n]

// Read exact number of bytes
buf := make([]byte, 100)
_, err := io.ReadFull(conn, buf)

// Read until delimiter (buffered)
reader := bufio.NewReader(conn)
line, err := reader.ReadString('\n')

// Scan lines
scanner := bufio.NewScanner(conn)
for scanner.Scan() {
    line := scanner.Text()
    // Process line
}
```

### Writing Data

```go
// Write bytes
n, err := conn.Write([]byte("Hello\n"))

// Write string
io.WriteString(conn, "Hello\n")

// Formatted write
fmt.Fprintf(conn, "Hello %s\n", name)

// Buffered write (more efficient)
writer := bufio.NewWriter(conn)
fmt.Fprintf(writer, "Line 1\n")
fmt.Fprintf(writer, "Line 2\n")
writer.Flush()  // Send buffered data
```

### Timeouts

```go
// Set read deadline
conn.SetReadDeadline(time.Now().Add(30 * time.Second))

// Set write deadline
conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

// Set overall deadline
conn.SetDeadline(time.Now().Add(1 * time.Minute))

// Clear deadline
conn.SetDeadline(time.Time{})

// Check for timeout error
if err != nil {
    if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
        // Timeout occurred
    }
}
```

### Connection Info

```go
// Get remote address
remoteAddr := conn.RemoteAddr().String()  // "192.168.1.100:54321"

// Get local address
localAddr := conn.LocalAddr().String()  // "192.168.1.1:8080"

// Get listener address
addr := listener.Addr().String()  // "[::]:8080"

// Type assert to get more details
if tcpAddr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
    ip := tcpAddr.IP      // IP address
    port := tcpAddr.Port  // Port number
}
```

### TCP-Specific Options

```go
// Type assert to TCPConn
tcpConn, ok := conn.(*net.TCPConn)
if !ok {
    // Not a TCP connection
}

// Disable Nagle's algorithm (reduce latency)
tcpConn.SetNoDelay(true)

// Set keep-alive
tcpConn.SetKeepAlive(true)
tcpConn.SetKeepAlivePeriod(30 * time.Second)

// Set linger (how long to wait for unsent data on close)
tcpConn.SetLinger(0)  // Close immediately (send RST)
tcpConn.SetLinger(10) // Wait up to 10 seconds

// Set buffer sizes
tcpConn.SetReadBuffer(256 * 1024)   // 256KB
tcpConn.SetWriteBuffer(256 * 1024)  // 256KB
```

### Binary Protocol

```go
import "encoding/binary"

// Write binary data
length := uint32(len(data))
binary.Write(conn, binary.BigEndian, length)
conn.Write(data)

// Read binary data
var length uint32
binary.Read(conn, binary.BigEndian, &length)
data := make([]byte, length)
io.ReadFull(conn, data)
```

### TLS/SSL

```go
import "crypto/tls"

// Server
cert, _ := tls.LoadX509KeyPair("server.crt", "server.key")
config := &tls.Config{Certificates: []tls.Certificate{cert}}
listener, _ := tls.Listen("tcp", ":8080", config)

// Client
config := &tls.Config{InsecureSkipVerify: true}  // Testing only!
conn, _ := tls.Dial("tcp", "localhost:8080", config)
```

### Error Handling

```go
// Check for specific errors
if err == io.EOF {
    // Connection closed normally
}

// Check for network errors
if netErr, ok := err.(net.Error); ok {
    if netErr.Timeout() {
        // Timeout
    }
    if netErr.Temporary() {
        // Temporary error, can retry
    }
}

// Check for specific error types
var opErr *net.OpError
if errors.As(err, &opErr) {
    // Network operation error
    fmt.Println("Operation:", opErr.Op)
    fmt.Println("Network:", opErr.Net)
    fmt.Println("Address:", opErr.Addr)
}
```

### Common Patterns

#### Echo Server Handler

```go
func handleClient(conn net.Conn) {
    defer conn.Close()

    scanner := bufio.NewScanner(conn)
    writer := bufio.NewWriter(conn)

    for scanner.Scan() {
        line := scanner.Text()
        fmt.Fprintf(writer, "ECHO: %s\n", line)
        writer.Flush()
    }

    if err := scanner.Err(); err != nil {
        log.Println("Error:", err)
    }
}
```

#### Request-Response Client

```go
func sendRequest(addr, request string) (string, error) {
    conn, err := net.Dial("tcp", addr)
    if err != nil {
        return "", err
    }
    defer conn.Close()

    // Send request
    writer := bufio.NewWriter(conn)
    fmt.Fprintf(writer, "%s\n", request)
    writer.Flush()

    // Read response
    reader := bufio.NewReader(conn)
    response, err := reader.ReadString('\n')
    if err != nil {
        return "", err
    }

    return strings.TrimRight(response, "\n"), nil
}
```

#### Graceful Shutdown

```go
func main() {
    listener, _ := net.Listen("tcp", ":8080")

    var wg sync.WaitGroup
    done := make(chan struct{})

    // Shutdown handler
    go func() {
        sigCh := make(chan os.Signal, 1)
        signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
        <-sigCh

        close(done)
        listener.Close()
    }()

    // Accept loop
    for {
        conn, err := listener.Accept()
        if err != nil {
            select {
            case <-done:
                wg.Wait()
                return
            default:
                continue
            }
        }

        wg.Add(1)
        go func() {
            defer wg.Done()
            handleConnection(conn, done)
        }()
    }
}
```

#### Bidirectional Proxy

```go
func proxy(client, backend net.Conn) {
    done := make(chan error, 2)

    go func() {
        _, err := io.Copy(backend, client)
        done <- err
    }()

    go func() {
        _, err := io.Copy(client, backend)
        done <- err
    }()

    <-done  // Wait for one direction to complete
    client.Close()
    backend.Close()
}
```

## Protocol Design Patterns

### Line-Based Protocol

```
Client → Server: COMMAND arg1 arg2\n
Server → Client: RESPONSE data\n
```

Pros: Human-readable, easy to debug
Cons: Need to escape newlines in data

### Length-Prefixed Protocol

```
[4 bytes: length][data]
```

Pros: No escaping, efficient parsing
Cons: Binary, harder to debug

### Self-Describing Protocol

```
[1 byte: message type][4 bytes: length][data]
```

Pros: Flexible, can extend easily
Cons: More complex parsing

### Frame-Based Protocol

```
[frame header: sync, type, length, checksum][payload]
```

Pros: Robust, error detection
Cons: Most complex

## Performance Tips

1. **Use buffered I/O**: `bufio.Reader` and `bufio.Writer`
2. **Reuse connections**: Connection pooling for clients
3. **Set buffer sizes**: Increase for high-throughput apps
4. **Disable Nagle's**: For low-latency interactive protocols
5. **Use goroutines**: One per connection for concurrency
6. **Limit concurrency**: Semaphore to prevent resource exhaustion
7. **Set timeouts**: Prevent resource leaks from stalled connections

## Debugging Tips

### View TCP connections

```bash
# Linux/Mac
netstat -an | grep :8080
lsof -i :8080

# Windows
netstat -an | findstr :8080
```

### Capture packets

```bash
# tcpdump (requires root)
sudo tcpdump -i lo -A port 8080

# Wireshark
# GUI tool, filter: tcp.port == 8080
```

### Test with command-line tools

```bash
# netcat
nc localhost 8080

# telnet
telnet localhost 8080

# OpenSSL (for TLS)
openssl s_client -connect localhost:8080
```

### Enable logging

```go
import "log"

log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
log.Printf("[%s] Message: %q", conn.RemoteAddr(), msg)
```

## Common Pitfalls

1. ❌ Forgetting `defer conn.Close()`
2. ❌ Not checking errors from `Write()`
3. ❌ Forgetting `Flush()` on buffered writers
4. ❌ Assuming one `Read()` gets complete message
5. ❌ Not handling `io.EOF` separately
6. ❌ Blocking `Accept()` in main goroutine without shutdown handler
7. ❌ Not setting timeouts (connections leak if client disappears)
8. ❌ Reusing buffers across goroutines

## Further Reading

- [Go net package docs](https://pkg.go.dev/net)
- [Go bufio package docs](https://pkg.go.dev/bufio)
- [TCP RFC 793](https://www.rfc-editor.org/rfc/rfc793)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
