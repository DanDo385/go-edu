# 33-tcp-echo-server-client

**TCP Echo Server**

Build a TCP echo server and client.

## What You'll Learn

- net.Listen and net.Dial
- Connection handling
- Buffer management
- Protocol design

## Functions to Implement

| Function | Description |
|----------|-------------|
| Implement TCP echo | Server echoes client messages |

## Project Structure

```
33-tcp-echo-server-client/
├── cmd/
│   ├── app/main.go      # CLI with custom arguments
│   └── dev/main.go      # Debug harness
├── internal/tcpechoserverclient/
│   ├── exercise.go      # YOUR CODE HERE
│   ├── exercise_test.go # Tests
│   └── solution.reference.go # Reference solution
└── README.md
```

## CLI Usage

```bash
cd minis/33-tcp-echo-server-client

# Start server
go run ./cmd/app/main.go server --port 9000

# In another terminal, run client
go run ./cmd/app/main.go client --addr localhost:9000

# Debug harness
go run ./cmd/dev/main.go
```

## Quick Copy & Paste

```bash
# Start server
go run ./cmd/app/main.go server --port 9000

# Connect client (in another terminal)
go run ./cmd/app/main.go client --addr localhost:9000

# Test with netcat
echo "hello" | nc localhost 9000

# Debug harness
go run ./cmd/dev/main.go
```

## Key Concepts

1. **net.Listen("tcp", addr)**: Start server
2. **net.Dial("tcp", addr)**: Connect client
3. **io.Copy**: Efficient data copying
4. **Connection Pooling**: Reuse connections

## Next Steps

After completing this exercise, proceed to `minis/34-rate-limiter-token-bucket`.
