# 32-websocket-chatroom

**WebSocket Chat Room**

Build a real-time chat room with WebSockets.

## What You'll Learn

- WebSocket protocol
- gorilla/websocket library
- Broadcasting to clients
- Connection management

## Functions to Implement

| Function | Description |
|----------|-------------|
| Implement WebSocket chat | Real-time messaging |

## Project Structure

```
32-websocket-chatroom/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness
├── internal/websocketchatroom/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/32-websocket-chatroom

# Start chat server
go run ./cmd/app/main.go --port 8080

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Start chat server
go run ./cmd/app/main.go --port 8080

# Connect with wscat (install: npm i -g wscat)
wscat -c ws://localhost:8080/ws

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **WebSocket Upgrade**: HTTP → WebSocket
2. **Full Duplex**: Two-way communication
3. **Hub Pattern**: Central message router
4. **Ping/Pong**: Connection keep-alive

## Next Steps

After completing this exercise, proceed to `minis/33-tcp-echo-server-client`.
