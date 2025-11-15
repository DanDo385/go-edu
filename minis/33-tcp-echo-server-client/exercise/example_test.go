package exercise_test

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/example/go-10x-minis/minis/33-tcp-echo-server-client/exercise"
)

// Example_echoClient demonstrates basic usage of the echo client
func Example_echoClient() {
	// Note: This example requires a running server
	// Start server: go run cmd/tcp-server/main.go

	// Send a single message and receive response
	response, err := exercise.EchoClient("localhost:8080", "Hello, World!")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
	// Output: ECHO: Hello, World!
}

// Example_persistentConnection demonstrates reusing a connection for multiple messages
func Example_persistentConnection() {
	// Connect to server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Send multiple messages on the same connection
	messages := []string{"First", "Second", "Third"}
	for _, msg := range messages {
		response, err := exercise.SendMessage(conn, msg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(response)
	}

	// Output:
	// ECHO: First
	// ECHO: Second
	// ECHO: Third
}

// Example_interactiveClient demonstrates an interactive client session
func Example_interactiveClient() {
	// Connect to server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Create reader for server responses
	reader := bufio.NewReader(conn)

	// Send messages and read responses
	messages := []string{
		"Hello",
		"How are you?",
		"Goodbye",
	}

	for _, msg := range messages {
		// Send message
		writer := bufio.NewWriter(conn)
		fmt.Fprintf(writer, "%s\n", msg)
		writer.Flush()

		// Read response
		response, err := exercise.ReadResponse(reader)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(response)
	}

	// Output:
	// ECHO: Hello
	// ECHO: How are you?
	// ECHO: Goodbye
}

// Example_concurrentClients demonstrates multiple clients connecting simultaneously
func Example_concurrentClients() {
	// Number of concurrent clients
	const numClients = 5

	// Channel to collect results
	results := make(chan string, numClients)

	// Start concurrent clients
	for i := 0; i < numClients; i++ {
		go func(clientID int) {
			message := fmt.Sprintf("Client %d says hello", clientID)
			response, err := exercise.EchoClient("localhost:8080", message)
			if err != nil {
				results <- fmt.Sprintf("Error: %v", err)
				return
			}
			results <- response
		}(i)
	}

	// Collect results
	for i := 0; i < numClients; i++ {
		fmt.Println(<-results)
	}

	// Output will vary due to concurrent execution
}

// Example_errorHandling demonstrates proper error handling
func Example_errorHandling() {
	// Try to connect to non-existent server
	_, err := exercise.EchoClient("localhost:9999", "Test")
	if err != nil {
		fmt.Println("Connection failed (expected):", err != nil)
	}

	// Output: Connection failed (expected): true
}

// Example_timeout demonstrates connection timeout handling
func Example_timeout() {
	// Connect to server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Set read timeout
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	// Send message
	writer := bufio.NewWriter(conn)
	fmt.Fprintf(writer, "Test\n")
	writer.Flush()

	// Read response (will timeout if server doesn't respond in 2 seconds)
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Println("Timeout occurred")
			return
		}
		log.Fatal(err)
	}

	fmt.Println(response)
}
