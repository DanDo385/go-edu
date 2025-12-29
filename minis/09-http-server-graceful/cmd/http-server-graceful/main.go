package main

import (
	"context"
	"flag"   // Command-line argument parsing
	"log"    // Logging
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"   // Time operations

	"github.com/example/go-10x-minis/minis/09-http-server-graceful/exercise"
)

/*
===================================================================
HTTP Server with Graceful Shutdown
===================================================================

This program demonstrates an HTTP server that handles shutdown
signals gracefully, allowing in-flight requests to complete.

USAGE:
    go run ./cmd/http-server-graceful/main.go
    go run ./cmd/http-server-graceful/main.go -port=9090
    go run ./cmd/http-server-graceful/main.go -port=8080 -timeout=10s
    go run ./cmd/http-server-graceful/main.go -port=3000 -timeout=30s

    # Test with curl in another terminal:
    curl http://localhost:8080/ping
    curl http://localhost:8080/items
    curl -X POST http://localhost:8080/items -d '{"name":"test"}'

    # Send shutdown signal:
    Ctrl+C (SIGINT) or kill -TERM <pid>

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging (select "Debug: Run main.go")
- Step Into (F11) to enter exercise package functions
- Watch Variables panel to see server state and shutdown process
- Use Debug Console to inspect server status
- See /RUN_DEBUG.md for comprehensive debugging guide

GRACEFUL SHUTDOWN:
- Server listens for SIGINT (Ctrl+C) and SIGTERM signals
- On signal, server stops accepting new requests
- In-flight requests are allowed to complete (up to timeout)
- Resources are cleaned up before exit
*/

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect default values before parsing
	// DEBUG: Before flag.Parse(), all flags have their default values
	// DEBUG: In Variables panel, you'll see: port=":8080", timeout=5s

	port := flag.String("port", ":8080", "Server port (e.g., :8080 or :9090)")
	shutdownTimeout := flag.Duration("timeout", 5*time.Second, "Graceful shutdown timeout")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")

	flag.Parse()

	// BREAKPOINT: Set breakpoint here to inspect parsed values
	// DEBUG: After flag.Parse(), variables now hold actual command-line values
	// DEBUG: Watch 'port', 'shutdownTimeout' - they may have changed from defaults
	// DEBUG: Hover over *port to see the dereferenced string value

	if *verbose {
		log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)
	}

	log.Printf("Configuration: port=%s, shutdown_timeout=%v", *port, *shutdownTimeout)

	// ============================================================================
	// STEP 2: Create In-Memory Store
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before store creation
	// DEBUG: Step Into (F11) on exercise.NewMemStore() to see initialization
	// DEBUG: Watch how the store struct is allocated

	store := exercise.NewMemStore()

	// BREAKPOINT: Set breakpoint here after store creation
	// DEBUG: Inspect 'store' variable in Variables panel
	// DEBUG: Expand store to see internal fields (mutex, items map)
	// DEBUG: Initially, store should be empty

	// ============================================================================
	// STEP 3: Set Up HTTP Routes
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before route registration
	// DEBUG: Watch how routes are registered with the mux
	// DEBUG: Step Into (F11) to see handler registration

	mux := http.NewServeMux()

	// BREAKPOINT: Set breakpoint here after mux creation
	// DEBUG: Inspect 'mux' - it's the HTTP request multiplexer
	// DEBUG: Routes will be added to this mux

	exercise.RegisterRoutes(mux, store)

	// BREAKPOINT: Set breakpoint here after route registration
	// DEBUG: Routes are now registered (GET /items, POST /items, etc.)
	// DEBUG: Step Into (F11) on RegisterRoutes to see all handlers

	// ============================================================================
	// STEP 4: Create HTTP Server
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before server creation
	// DEBUG: Watch how the http.Server is configured
	// DEBUG: Note the address and handler

	srv := exercise.NewServer(*port, mux)

	// BREAKPOINT: Set breakpoint here after server creation
	// DEBUG: Inspect 'srv' variable in Variables panel
	// DEBUG: Expand srv to see Addr, Handler, ReadTimeout, WriteTimeout, etc.
	// DEBUG: The server is configured but not yet running

	log.Printf("Server configuration: addr=%s", srv.Addr)

	// ============================================================================
	// STEP 5: Start Server in Background Goroutine
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before starting server
	// DEBUG: The server will run in a separate goroutine
	// DEBUG: Main goroutine will continue to signal handling

	go func() {
		// BREAKPOINT: Set breakpoint here inside goroutine (won't pause initially)
		// DEBUG: This goroutine runs concurrently with main
		// DEBUG: Server will block here until shutdown or error

		log.Printf("Server starting on %s", *port)
		log.Println("Press Ctrl+C to shutdown gracefully")

		// BREAKPOINT: Set breakpoint here before ListenAndServe
		// DEBUG: Step Into (F11) to see server startup (if debugger supports)
		// DEBUG: Server will block here, accepting connections

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// BREAKPOINT: Set breakpoint here on server error
			// DEBUG: This is reached only if server fails to start or crashes
			// DEBUG: Normal shutdown produces http.ErrServerClosed (not logged)

			log.Fatalf("Server error: %v", err)
		}

		// BREAKPOINT: Set breakpoint here after server stops
		// DEBUG: Server has stopped (shutdown completed)
		log.Println("Server stopped")
	}()

	// BREAKPOINT: Set breakpoint here after starting goroutine
	// DEBUG: Main goroutine continues while server runs in background
	// DEBUG: Server is now accepting HTTP requests

	// Small delay to ensure server starts before we continue
	time.Sleep(100 * time.Millisecond)
	log.Println("Server is ready to accept requests")
	log.Println("Try: curl http://localhost" + *port + "/ping")

	// ============================================================================
	// STEP 6: Set Up Signal Handling for Graceful Shutdown
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before signal setup
	// DEBUG: We'll create a channel to receive OS signals
	// DEBUG: This allows the program to handle Ctrl+C gracefully

	sigCh := make(chan os.Signal, 1)

	// BREAKPOINT: Set breakpoint here after channel creation
	// DEBUG: Inspect 'sigCh' - it's a buffered channel with capacity 1
	// DEBUG: The channel will receive signals from the OS

	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// BREAKPOINT: Set breakpoint here after signal.Notify
	// DEBUG: Program is now listening for SIGINT (Ctrl+C) and SIGTERM
	// DEBUG: When signal is received, it will be sent to sigCh

	// ============================================================================
	// STEP 7: Wait for Shutdown Signal
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before waiting for signal
	// DEBUG: The following line will block until a signal is received
	// DEBUG: Send Ctrl+C or kill -TERM to trigger shutdown

	<-sigCh

	// BREAKPOINT: Set breakpoint here after receiving signal
	// DEBUG: A shutdown signal was received (SIGINT or SIGTERM)
	// DEBUG: Program will now begin graceful shutdown

	log.Println("Shutdown signal received...")

	// ============================================================================
	// STEP 8: Graceful Shutdown with Timeout
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before shutdown context creation
	// DEBUG: We'll create a context with timeout for shutdown
	// DEBUG: This ensures shutdown completes within the timeout

	ctx, cancel := context.WithTimeout(context.Background(), *shutdownTimeout)
	defer cancel()

	// BREAKPOINT: Set breakpoint here after context creation
	// DEBUG: Inspect 'ctx' - it has a deadline set to now + shutdownTimeout
	// DEBUG: Use ctx.Deadline() in Debug Console to see exact deadline

	log.Printf("Graceful shutdown initiated (timeout: %v)...", *shutdownTimeout)

	// BREAKPOINT: Set breakpoint here before calling GracefulShutdown
	// DEBUG: Step Into (F11) on exercise.GracefulShutdown() to see shutdown logic
	// DEBUG: Watch how server stops accepting new connections
	// DEBUG: Observe how in-flight requests are allowed to complete

	if err := exercise.GracefulShutdown(ctx, srv); err != nil {
		// BREAKPOINT: Set breakpoint here on shutdown error
		// DEBUG: Error could be: context deadline exceeded (timeout)
		// DEBUG: Inspect 'err' to see what went wrong during shutdown

		log.Fatalf("Shutdown error: %v", err)
	}

	// BREAKPOINT: Set breakpoint here after successful shutdown
	// DEBUG: Server has shut down gracefully
	// DEBUG: All in-flight requests completed (or timeout occurred)

	log.Println("Server stopped gracefully")

	// ============================================================================
	// DEBUGGING TIPS
	// ============================================================================
	//
	// KEY DEBUGGING TECHNIQUES:
	//
	// 1. Set breakpoints at "// BREAKPOINT:" comments
	// 2. Press F5 to start debugging
	// 3. Use F10 (Step Over) to execute line by line
	// 4. Use F11 (Step Into) to enter exercise package functions
	// 5. Watch Variables panel to see server state
	// 6. Add watch expressions: srv.Addr, ctx.Err()
	// 7. Use Call Stack panel to see goroutine hierarchy
	//
	// SERVER DEBUGGING:
	// - Server runs in a goroutine (concurrent execution)
	// - Main goroutine blocks on signal channel
	// - Use Debug Console: srv.Addr to check server address
	// - Check server state in Variables panel
	//
	// TESTING THE SERVER:
	// - Open another terminal while debugging
	// - Use curl to send requests:
	//   curl http://localhost:8080/ping
	//   curl http://localhost:8080/items
	//   curl -X POST http://localhost:8080/items -d '{"name":"test"}'
	// - Set breakpoint in handler (in exercise package)
	// - Send request to trigger breakpoint
	//
	// GRACEFUL SHUTDOWN DEBUGGING:
	// - Start server and send some requests
	// - Set breakpoint after <-sigCh
	// - Send Ctrl+C or kill signal
	// - Step through shutdown logic
	// - Watch how server stops accepting new connections
	// - Verify in-flight requests complete before shutdown
	//
	// SIGNAL HANDLING:
	// - Signal channel receives OS signals
	// - SIGINT: Ctrl+C in terminal
	// - SIGTERM: kill command (kill -TERM <pid>)
	// - Program blocks on <-sigCh until signal arrives
	// - After signal, graceful shutdown begins
	//
	// CONTEXT WITH TIMEOUT:
	// - Context has deadline for shutdown
	// - If shutdown takes too long, context expires
	// - Use ctx.Deadline() to see exact deadline
	// - Use ctx.Err() to check if context is done
	//
	// GOROUTINE DEBUGGING:
	// - Server runs in background goroutine
	// - Use Goroutines panel in debugger
	// - Switch between goroutines to see state
	// - Breakpoints in goroutines may not pause immediately
	//
	// See /RUN_DEBUG.md for comprehensive debugging guide
}
