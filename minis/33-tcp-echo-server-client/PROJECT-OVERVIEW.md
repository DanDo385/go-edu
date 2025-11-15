# Project 33: TCP Echo Server/Client - Project Overview

## What You'll Build

A complete TCP echo server/client application that demonstrates low-level network programming in Go. The server accepts multiple concurrent client connections and echoes back any messages it receives.

## Learning Objectives

After completing this project, you will understand:

1. **TCP Fundamentals**
   - How TCP connections work (three-way handshake)
   - Difference between TCP and UDP
   - TCP state machine and connection lifecycle
   - Byte streams vs message boundaries

2. **Go's net Package**
   - Creating TCP listeners with `net.Listen()`
   - Accepting connections with `listener.Accept()`
   - Connecting to servers with `net.Dial()`
   - Reading and writing data to connections

3. **Buffered I/O**
   - Using `bufio.Scanner` for line-based reading
   - Using `bufio.Writer` for efficient writing
   - When and why to flush buffers
   - Performance implications of buffering

4. **Protocol Design**
   - Line-based protocols with delimiters
   - Binary protocols with length prefixes
   - Protocol framing and message boundaries
   - Error handling in custom protocols

5. **Concurrent Programming**
   - Handling multiple clients with goroutines
   - Synchronizing access to shared state
   - Graceful shutdown with signal handling
   - Resource cleanup and connection management

## Project Structure

```
33-tcp-echo-server-client/
├── README.md                     # Comprehensive tutorial (1,481 lines)
├── DEMO.md                       # How to run and test the project
├── QUICK-REFERENCE.md            # Quick reference for TCP programming
├── PROJECT-OVERVIEW.md           # This file
│
├── cmd/
│   ├── tcp-server/
│   │   └── main.go              # Complete TCP echo server (128 lines)
│   │                            # - Listens on configurable port
│   │                            # - Handles multiple clients concurrently
│   │                            # - Graceful shutdown support
│   │                            # - Comprehensive logging
│   │
│   └── tcp-client/
│       └── main.go              # Interactive TCP client (89 lines)
│                                # - Connects to server
│                                # - Reads from stdin, sends to server
│                                # - Displays server responses
│                                # - Handles errors gracefully
│
└── exercise/
    ├── exercise.go              # Exercise skeleton with TODOs (150 lines)
    ├── solution.go              # Complete solution (136 lines)
    ├── exercise_test.go         # Comprehensive tests (397 lines)
    ├── example_test.go          # Usage examples (162 lines)
    └── EXERCISES.md             # 8 progressive exercises (672 lines)
                                 # 1. Command Protocol
                                 # 2. Connection Timeout
                                 # 3. Binary Protocol
                                 # 4. TLS/SSL Encryption
                                 # 5. Chat Server
                                 # 6. Connection Pooling
                                 # 7. Load Balancer
                                 # 8. Port Scanner
```

## File Descriptions

### Core Documentation

**README.md** (1,481 lines)
- Complete first-principles explanation of TCP
- Network stack fundamentals (OSI layers)
- TCP vs UDP comparison
- Three-way handshake explained
- Sockets and port numbers
- Byte streams and protocol framing
- Complete solution walkthrough
- Real-world applications
- Common mistakes and how to avoid them

**DEMO.md**
- Quick start guide
- Terminal commands to run server/client
- Testing with netcat and telnet
- Running tests and benchmarks
- Troubleshooting common issues
- Performance testing examples

**QUICK-REFERENCE.md**
- Cheat sheet for TCP programming
- Common patterns (echo server, request-response client)
- Error handling examples
- Binary protocol helpers
- TLS/SSL snippets
- Debugging tips

### Implementation Files

**cmd/tcp-server/main.go** (128 lines)
Complete TCP echo server with:
- Configurable listening address
- Concurrent client handling (goroutines)
- Graceful shutdown (signal handling)
- Per-client logging
- Error recovery
- Welcome message to clients

**cmd/tcp-client/main.go** (89 lines)
Interactive TCP client with:
- Connection to server
- Reading from stdin
- Concurrent response handling
- Error reporting
- Clean disconnection

### Exercise Files

**exercise/exercise.go** (150 lines)
Skeleton implementation with:
- `StartEchoServer()` - TODO
- `EchoClient()` - TODO
- `SendMessage()` - TODO
- `ReadResponse()` - TODO
- Detailed documentation for each function
- Implementation hints in comments

**exercise/solution.go** (136 lines)
Complete reference implementation using `//go:build solution` tag

**exercise/exercise_test.go** (397 lines)
Comprehensive test suite:
- Basic echo functionality
- Multiple messages
- Empty messages
- Unicode/emoji support
- Concurrent clients (10 simultaneous)
- Long messages (1000 chars)
- Persistent connections
- Connection closure
- Error cases
- Benchmarks (single client, concurrent, persistent)

**exercise/example_test.go** (162 lines)
Runnable examples:
- Basic client usage
- Persistent connection
- Interactive client
- Concurrent clients
- Error handling
- Timeout handling

**exercise/EXERCISES.md** (672 lines)
8 progressive exercises with increasing difficulty:
1. ⭐⭐ Command Protocol - Multi-command server
2. ⭐⭐ Connection Timeout - Idle connection management
3. ⭐⭐⭐ Binary Protocol - Length-prefixed messages
4. ⭐⭐⭐ TLS/SSL Encryption - Secure connections
5. ⭐⭐⭐⭐ Chat Server - Multi-user broadcast
6. ⭐⭐⭐ Connection Pooling - Efficient reuse
7. ⭐⭐⭐⭐ Load Balancer - Distribute across backends
8. ⭐⭐ Port Scanner - Concurrent port scanning

## How to Use This Project

### 1. Read the Tutorial

Start with `README.md` for a comprehensive explanation:
```bash
cat README.md | less
# or open in your editor
```

### 2. Run the Demo

Follow `DEMO.md` to run the server and client:
```bash
# Terminal 1
go run cmd/tcp-server/main.go

# Terminal 2
go run cmd/tcp-client/main.go
```

### 3. Complete the Exercises

Work through `exercise/exercise.go`:
```bash
cd exercise
# Read the TODOs
cat exercise.go

# Run tests (they should fail)
go test -v

# Implement the functions
# ... edit exercise.go ...

# Test again
go test -v

# When all tests pass, you're done!
```

### 4. Stretch Goals

Work through `exercise/EXERCISES.md` for advanced challenges:
```bash
cat exercise/EXERCISES.md | less
```

### 5. Reference

Use `QUICK-REFERENCE.md` as a cheat sheet while coding.

## Testing

### Run All Tests
```bash
cd exercise
go test -v
```

### Run Specific Test
```bash
go test -run TestEchoClient -v
```

### Run Benchmarks
```bash
go test -tags=solution -bench=. -benchmem
```

### Test Coverage
```bash
go test -tags=solution -cover
```

### With Solution
```bash
go test -tags=solution -v
```

## Key Concepts Covered

### Networking Concepts
- OSI Model and TCP/IP stack
- IP addresses and port numbers
- Sockets and file descriptors
- Three-way handshake
- Connection states (ESTABLISHED, FIN_WAIT, TIME_WAIT, etc.)
- Flow control and congestion control
- Nagle's algorithm

### Protocol Design
- Message framing (delimiters, length-prefix, self-describing)
- Text vs binary protocols
- Error handling in protocols
- Protocol versioning
- Backward compatibility

### Go Programming Patterns
- Accepting connections in a loop
- One goroutine per connection
- Graceful shutdown with channels and WaitGroups
- Buffered I/O for performance
- Error handling and recovery
- Resource cleanup with defer

### Performance Considerations
- Connection pooling vs new connections
- Buffering to reduce system calls
- Setting socket buffer sizes
- TCP_NODELAY for low latency
- Concurrent connections limits

## Common Patterns Demonstrated

1. **Server Accept Loop**
   ```go
   for {
       conn, err := listener.Accept()
       if err != nil { /* handle */ }
       go handleClient(conn)
   }
   ```

2. **Graceful Shutdown**
   ```go
   sigCh := make(chan os.Signal, 1)
   signal.Notify(sigCh, os.Interrupt)
   <-sigCh
   listener.Close()
   wg.Wait()
   ```

3. **Buffered Line Reading**
   ```go
   scanner := bufio.NewScanner(conn)
   for scanner.Scan() {
       line := scanner.Text()
   }
   ```

4. **Buffered Writing**
   ```go
   writer := bufio.NewWriter(conn)
   fmt.Fprintln(writer, "response")
   writer.Flush()
   ```

## Real-World Applications

This project teaches patterns used in:

- **Databases**: PostgreSQL, MySQL client-server communication
- **Caches**: Redis, Memcached protocols
- **Message Queues**: RabbitMQ, Apache Kafka
- **Remote Procedure Calls**: gRPC, Thrift
- **Chat Systems**: Slack, Discord
- **Game Servers**: Multiplayer game networking
- **IoT**: Device-to-cloud communication
- **Proxies**: Load balancers, reverse proxies

## Troubleshooting

### Server won't start
- Port already in use → Change port or kill other process
- Permission denied → Use port > 1024 or run with sudo

### Client can't connect
- Server not running → Start server first
- Wrong address → Check host:port
- Firewall blocking → Check firewall rules

### Tests failing
- Make sure server is NOT running during tests (tests start their own)
- Check you're in the exercise/ directory
- Try running with -v flag for more details

### Build errors
- Run `go mod tidy` to fix dependencies
- Check Go version (requires 1.18+ for generics)

## Next Steps

After completing this project:

1. **Implement the stretch goals** in EXERCISES.md
2. **Build a real application**:
   - Multi-user chat server
   - Simple database (key-value store)
   - Distributed cache
   - Load balancer

3. **Explore related topics**:
   - UDP programming (Project 34+)
   - HTTP/2 and WebSockets (Project 32)
   - gRPC (Project 10)
   - TLS/SSL encryption

4. **Study production systems**:
   - Redis protocol (RESP)
   - PostgreSQL wire protocol
   - HTTP/1.1 and HTTP/2
   - MQTT for IoT

## Resources

- [Go net package](https://pkg.go.dev/net)
- [Go bufio package](https://pkg.go.dev/bufio)
- [TCP RFC 793](https://www.rfc-editor.org/rfc/rfc793)
- [Effective Go - Concurrency](https://go.dev/doc/effective_go#concurrency)

## Summary Statistics

- **Total lines of code**: 3,215
- **Documentation**: 1,481 lines (README.md)
- **Implementation**: 353 lines (server + client + exercise)
- **Tests**: 559 lines
- **Exercises**: 672 lines

This project provides a complete, production-quality example of TCP programming in Go with extensive documentation and progressive exercises to master network programming fundamentals.

Happy coding! 🚀
