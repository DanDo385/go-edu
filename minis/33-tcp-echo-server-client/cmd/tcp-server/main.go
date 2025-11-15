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
	// Configuration: Use command-line argument or default to :8080
	addr := ":8080"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	// Create TCP listener
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Printf("TCP echo server listening on %s", listener.Addr())
	log.Println("Press Ctrl+C to stop the server")

	// Track active connections for graceful shutdown
	var wg sync.WaitGroup
	done := make(chan struct{})

	// Handle graceful shutdown on SIGINT/SIGTERM
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh

		log.Println("\nShutting down server gracefully...")
		close(done)       // Signal all handlers to stop
		listener.Close()  // Stop accepting new connections
	}()

	// Accept loop: Handle incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-done:
				// Shutdown initiated, wait for active connections to finish
				log.Println("Waiting for active connections to close...")
				wg.Wait()
				log.Println("Server stopped")
				return
			default:
				// Temporary error, log and continue
				log.Printf("Accept error: %v", err)
				continue
			}
		}

		// Handle this client in a separate goroutine
		wg.Add(1)
		go handleClient(conn, &wg, done)
	}
}

// handleClient processes a single client connection
func handleClient(conn net.Conn, wg *sync.WaitGroup, done <-chan struct{}) {
	defer wg.Done()
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("Client connected: %s", clientAddr)
	defer log.Printf("Client disconnected: %s", clientAddr)

	// Create buffered reader for efficient line reading
	scanner := bufio.NewScanner(conn)

	// Create buffered writer for efficient writing
	writer := bufio.NewWriter(conn)

	// Send welcome message
	welcomeMsg := "Welcome to TCP Echo Server! Type messages and they will be echoed back.\n"
	if _, err := writer.WriteString(welcomeMsg); err != nil {
		log.Printf("[%s] Error sending welcome: %v", clientAddr, err)
		return
	}
	if err := writer.Flush(); err != nil {
		log.Printf("[%s] Error flushing welcome: %v", clientAddr, err)
		return
	}

	// Read loop: Process each line from the client
	for scanner.Scan() {
		// Check if server is shutting down
		select {
		case <-done:
			log.Printf("[%s] Server shutting down, closing connection", clientAddr)
			writer.WriteString("Server is shutting down. Goodbye!\n")
			writer.Flush()
			return
		default:
		}

		line := scanner.Text()
		log.Printf("[%s] Received: %q", clientAddr, line)

		// Echo back with prefix
		response := fmt.Sprintf("ECHO: %s\n", line)
		if _, err := writer.WriteString(response); err != nil {
			log.Printf("[%s] Write error: %v", clientAddr, err)
			return
		}

		// Flush immediately so client receives the response
		if err := writer.Flush(); err != nil {
			log.Printf("[%s] Flush error: %v", clientAddr, err)
			return
		}
	}

	// Check for scanner errors (vs normal EOF)
	if err := scanner.Err(); err != nil {
		log.Printf("[%s] Read error: %v", clientAddr, err)
	}
}
