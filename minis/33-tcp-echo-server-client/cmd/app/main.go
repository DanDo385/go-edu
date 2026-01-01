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
	// Configuration: Use command-line argument or default to localhost:8080
	addr := "localhost:8080"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	// Connect to the TCP server
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	log.Printf("Connected to %s", conn.RemoteAddr())
	fmt.Println("Type messages to send to the server (Ctrl+D or Ctrl+C to quit)")
	fmt.Println()

	// Channel to communicate errors from the reader goroutine
	errCh := make(chan error, 1)

	// Start a goroutine to read server responses concurrently
	// This allows us to receive messages while typing
	go func() {
		reader := bufio.NewReader(conn)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					errCh <- fmt.Errorf("read error: %w", err)
				}
				return
			}
			// Print server response
			fmt.Print("< " + line)
		}
	}()

	// Read user input from stdin
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(conn)

	for {
		fmt.Print("> ")

		// Check for errors from the reader goroutine
		select {
		case err := <-errCh:
			log.Fatal(err)
		default:
		}

		// Read line from stdin
		if !scanner.Scan() {
			// EOF or error
			break
		}

		line := scanner.Text()

		// Send to server
		if _, err := fmt.Fprintf(writer, "%s\n", line); err != nil {
			log.Fatalf("Write error: %v", err)
		}

		// Flush to ensure data is sent immediately
		if err := writer.Flush(); err != nil {
			log.Fatalf("Flush error: %v", err)
		}
	}

	// Check for scanner errors (vs normal EOF/Ctrl+D)
	if err := scanner.Err(); err != nil {
		log.Fatalf("Input error: %v", err)
	}

	fmt.Println("\nGoodbye!")
}
