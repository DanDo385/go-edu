# TCP Echo Server/Client Demo

This guide shows you how to run and test the TCP echo server and client.

## Quick Start

### Terminal 1: Start the Server

```bash
cd /home/user/go-edu/minis/33-tcp-echo-server-client
go run cmd/tcp-server/main.go
```

You should see:
```
2024/01/15 10:30:45 TCP echo server listening on [::]:8080
2024/01/15 10:30:45 Press Ctrl+C to stop the server
```

### Terminal 2: Run the Client

```bash
cd /home/user/go-edu/minis/33-tcp-echo-server-client
go run cmd/tcp-client/main.go
```

You should see:
```
2024/01/15 10:30:50 Connected to 127.0.0.1:8080
Type messages to send to the server (Ctrl+D or Ctrl+C to quit)

>
```

Now type messages and press Enter:

```
> Hello, World!
< ECHO: Hello, World!
> This is a test
< ECHO: This is a test
> Goodbye
< ECHO: Goodbye
```

Press `Ctrl+D` (or `Ctrl+C`) to quit.

## Custom Port

### Server on port 9000

```bash
go run cmd/tcp-server/main.go :9000
```

### Client connecting to port 9000

```bash
go run cmd/tcp-client/main.go localhost:9000
```

## Testing with netcat

You can also use `netcat` (nc) to test the server:

```bash
# In one terminal, start the server
go run cmd/tcp-server/main.go

# In another terminal, use netcat
echo "Hello from netcat" | nc localhost 8080
```

You should see:
```
Welcome to TCP Echo Server! Type messages and they will be echoed back.
ECHO: Hello from netcat
```

## Interactive Testing with netcat

```bash
nc localhost 8080
```

Then type messages interactively:
```
Welcome to TCP Echo Server! Type messages and they will be echoed back.
Test 1
ECHO: Test 1
Test 2
ECHO: Test 2
```

## Testing with telnet

```bash
telnet localhost 8080
```

## Multiple Concurrent Clients

Open multiple terminals and run the client in each:

```bash
# Terminal 1
go run cmd/tcp-client/main.go

# Terminal 2
go run cmd/tcp-client/main.go

# Terminal 3
go run cmd/tcp-client/main.go
```

Each client can send messages independently. Watch the server terminal to see all connections being handled.

## Running Tests

### Test the solution

```bash
cd exercise
go test -tags=solution -v
```

### Test your implementation

```bash
cd exercise
go test -v
```

### Run benchmarks

```bash
cd exercise
go test -tags=solution -bench=. -benchmem
```

Example output:
```
BenchmarkEchoClient-8                    1234    954321 ns/op    1024 B/op    12 allocs/op
BenchmarkConcurrentClients-8             5678    234567 ns/op    2048 B/op    24 allocs/op
BenchmarkPersistentConnection-8         12345     98765 ns/op     512 B/op     6 allocs/op
```

## Example Session Logs

### Server Side

```
2024/01/15 10:30:45 TCP echo server listening on [::]:8080
2024/01/15 10:30:45 Press Ctrl+C to stop the server
2024/01/15 10:30:50 Client connected: 127.0.0.1:54321
2024/01/15 10:30:52 [127.0.0.1:54321] Received: "Hello, World!"
2024/01/15 10:30:55 [127.0.0.1:54321] Received: "This is a test"
2024/01/15 10:30:58 Client connected: 127.0.0.1:54322
2024/01/15 10:31:00 [127.0.0.1:54322] Received: "Another client"
2024/01/15 10:31:05 Client disconnected: 127.0.0.1:54321
2024/01/15 10:31:10 Client disconnected: 127.0.0.1:54322
^C
2024/01/15 10:31:15 Shutting down server gracefully...
2024/01/15 10:31:15 Waiting for active connections to close...
2024/01/15 10:31:15 Server stopped
```

### Client Side

```
2024/01/15 10:30:50 Connected to 127.0.0.1:8080
Type messages to send to the server (Ctrl+D or Ctrl+C to quit)

> Hello, World!
< Welcome to TCP Echo Server! Type messages and they will be echoed back.
< ECHO: Hello, World!
> This is a test
< ECHO: This is a test
> ^D

Goodbye!
```

## Troubleshooting

### "bind: address already in use"

Another process is using port 8080. Either:
- Stop that process
- Use a different port: `go run cmd/tcp-server/main.go :9000`

### "connection refused"

Server is not running. Start it first with `go run cmd/tcp-server/main.go`

### "connection timeout"

- Check firewall settings
- Verify server is listening on the correct interface
- Use `127.0.0.1` instead of `localhost` or vice versa

## Network Monitoring

### View active connections (Linux/Mac)

```bash
# Show all TCP connections on port 8080
lsof -i TCP:8080

# Or using netstat
netstat -an | grep 8080
```

### View active connections (Windows)

```cmd
netstat -an | findstr 8080
```

### Monitor traffic with tcpdump (requires root)

```bash
sudo tcpdump -i lo -A port 8080
```

This shows you the actual TCP packets being sent and received.

## Performance Testing

### Using Apache Bench (ab)

While ab is designed for HTTP, you can use it to measure connection establishment:

```bash
ab -n 1000 -c 10 http://localhost:8080/
```

### Using custom Go script

Create `bench.go`:

```go
package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/example/go-10x-minis/minis/33-tcp-echo-server-client/exercise"
)

func main() {
	const (
		numRequests = 1000
		concurrency = 10
	)

	start := time.Now()

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			_, err := exercise.EchoClient("localhost:8080", "Test")
			if err != nil {
				log.Printf("Request %d failed: %v", id, err)
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("Completed %d requests in %v\n", numRequests, elapsed)
	fmt.Printf("Requests per second: %.2f\n", float64(numRequests)/elapsed.Seconds())
}
```

Run it:
```bash
go run bench.go
```

## Next Steps

1. Complete the exercises in `exercise/EXERCISES.md`
2. Implement your own protocol commands
3. Add TLS encryption
4. Build a chat server
5. Create a load balancer

Happy coding!
