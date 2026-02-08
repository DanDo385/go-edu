//go:build reference

package tcpechoserverclient

/*
Reference Solution
==================

This file is the canonical reference for this exercise. It keeps failure paths
explicit when an operation can fail, so callers can decide how to handle
errors at API boundaries.

Read this alongside exercise.go and the tests to understand the intended data
flow, ownership boundaries, and invariants that keep behavior deterministic.

Teaching notes:
- Memory/ownership: make copies when returning mutable data that should not
  alias internal state; share references only when aliasing is intentional.
- Invariants: establish assumptions close to construction, and rely on them in
  smaller helper functions to keep logic easy to audit.
- Error surfaces: prefer explicit returns over hidden panics so learners can
  reason about control flow in production-style code.
*/

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

// StartEchoServer starts a TCP echo server on the given address.
func StartEchoServer(addr string) error {
	// Create TCP listener
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	defer listener.Close()

	log.Printf("Echo server listening on %s", listener.Addr())

	// Accept loop: handle incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			// Log error but continue accepting
			log.Printf("Accept error: %v", err)
			continue
		}

		// Handle client in separate goroutine
		go handleClient(conn)
	}
}

// handleClient processes a single client connection
func handleClient(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	log.Printf("Client connected: %s", clientAddr)
	defer log.Printf("Client disconnected: %s", clientAddr)

	// Buffered reader for efficient line reading
	scanner := bufio.NewScanner(conn)

	// Buffered writer for efficient writing
	writer := bufio.NewWriter(conn)

	// Process each line
	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("[%s] Received: %q", clientAddr, line)

		// Echo back with prefix
		response := fmt.Sprintf("ECHO: %s\n", line)
		if _, err := writer.WriteString(response); err != nil {
			log.Printf("[%s] Write error: %v", clientAddr, err)
			return
		}

		// Flush immediately so client receives response
		if err := writer.Flush(); err != nil {
			log.Printf("[%s] Flush error: %v", clientAddr, err)
			return
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		log.Printf("[%s] Read error: %v", clientAddr, err)
	}
}

// EchoClient connects to a TCP echo server and sends a single message.
func EchoClient(addr, message string) (string, error) {
	// Connect to server
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return "", fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	// Send message
	writer := bufio.NewWriter(conn)
	if _, err := fmt.Fprintf(writer, "%s\n", message); err != nil {
		return "", fmt.Errorf("write error: %w", err)
	}
	if err := writer.Flush(); err != nil {
		return "", fmt.Errorf("flush error: %w", err)
	}

	// Read response
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read error: %w", err)
	}

	// Remove trailing newline
	return strings.TrimRight(response, "\n"), nil
}

// SendMessage sends a message to an already-established connection and reads the response.
func SendMessage(conn net.Conn, message string) (string, error) {
	// Send message
	writer := bufio.NewWriter(conn)
	if _, err := fmt.Fprintf(writer, "%s\n", message); err != nil {
		return "", fmt.Errorf("write error: %w", err)
	}
	if err := writer.Flush(); err != nil {
		return "", fmt.Errorf("flush error: %w", err)
	}

	// Read response
	reader := bufio.NewReader(conn)
	return ReadResponse(reader)
}

// ReadResponse reads a single line response from the connection.
func ReadResponse(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return "", fmt.Errorf("connection closed by server")
		}
		return "", fmt.Errorf("read error: %w", err)
	}

	// Remove trailing newline
	return strings.TrimRight(line, "\n"), nil
}
