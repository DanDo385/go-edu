package exercise

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestEchoClient tests the basic client functionality
func TestEchoClient(t *testing.T) {
	// Start test server
	listener, addr := startTestServer(t)
	defer listener.Close()

	// Test basic echo
	response, err := EchoClient(addr, "Hello, World!")
	if err != nil {
		t.Fatalf("EchoClient failed: %v", err)
	}

	expected := "ECHO: Hello, World!"
	if response != expected {
		t.Errorf("Expected %q, got %q", expected, response)
	}
}

// TestEchoClientMultipleMessages tests sending multiple messages sequentially
func TestEchoClientMultipleMessages(t *testing.T) {
	listener, addr := startTestServer(t)
	defer listener.Close()

	messages := []string{
		"First message",
		"Second message",
		"Third message with special chars: !@#$%",
	}

	for _, msg := range messages {
		response, err := EchoClient(addr, msg)
		if err != nil {
			t.Fatalf("EchoClient failed for message %q: %v", msg, err)
		}

		expected := "ECHO: " + msg
		if response != expected {
			t.Errorf("Expected %q, got %q", expected, response)
		}
	}
}

// TestEchoClientEmptyMessage tests sending an empty message
func TestEchoClientEmptyMessage(t *testing.T) {
	listener, addr := startTestServer(t)
	defer listener.Close()

	response, err := EchoClient(addr, "")
	if err != nil {
		t.Fatalf("EchoClient failed: %v", err)
	}

	expected := "ECHO: "
	if response != expected {
		t.Errorf("Expected %q, got %q", expected, response)
	}
}

// TestEchoClientUnicodeMessage tests sending Unicode characters
func TestEchoClientUnicodeMessage(t *testing.T) {
	listener, addr := startTestServer(t)
	defer listener.Close()

	messages := []string{
		"Hello 世界",
		"Emoji test 🚀 🎉 ✨",
		"Русский текст",
		"مرحبا بالعالم",
	}

	for _, msg := range messages {
		response, err := EchoClient(addr, msg)
		if err != nil {
			t.Fatalf("EchoClient failed for message %q: %v", msg, err)
		}

		expected := "ECHO: " + msg
		if response != expected {
			t.Errorf("Expected %q, got %q", expected, response)
		}
	}
}

// TestSendMessage tests the SendMessage helper function
func TestSendMessage(t *testing.T) {
	listener, addr := startTestServer(t)
	defer listener.Close()

	// Connect to server
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Send multiple messages on same connection
	messages := []string{"First", "Second", "Third"}
	for _, msg := range messages {
		response, err := SendMessage(conn, msg)
		if err != nil {
			t.Fatalf("SendMessage failed: %v", err)
		}

		expected := "ECHO: " + msg
		if response != expected {
			t.Errorf("Expected %q, got %q", expected, response)
		}
	}
}

// TestReadResponse tests the ReadResponse helper function
func TestReadResponse(t *testing.T) {
	listener, addr := startTestServer(t)
	defer listener.Close()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	// Send a message
	fmt.Fprintf(writer, "Test message\n")
	writer.Flush()

	// Read response using ReadResponse
	response, err := ReadResponse(reader)
	if err != nil {
		t.Fatalf("ReadResponse failed: %v", err)
	}

	expected := "ECHO: Test message"
	if response != expected {
		t.Errorf("Expected %q, got %q", expected, response)
	}
}

// TestConcurrentClients tests multiple concurrent client connections
func TestConcurrentClients(t *testing.T) {
	listener, addr := startTestServer(t)
	defer listener.Close()

	const numClients = 10
	var wg sync.WaitGroup
	errors := make(chan error, numClients)

	// Start multiple concurrent clients
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			message := fmt.Sprintf("Client %d message", clientID)
			response, err := EchoClient(addr, message)
			if err != nil {
				errors <- fmt.Errorf("client %d failed: %w", clientID, err)
				return
			}

			expected := "ECHO: " + message
			if response != expected {
				errors <- fmt.Errorf("client %d: expected %q, got %q", clientID, expected, response)
			}
		}(i)
	}

	// Wait for all clients to complete
	wg.Wait()
	close(errors)

	// Check for errors
	for err := range errors {
		t.Error(err)
	}
}

// TestLongMessage tests sending a long message
func TestLongMessage(t *testing.T) {
	listener, addr := startTestServer(t)
	defer listener.Close()

	// Create a message with 1000 characters
	message := strings.Repeat("a", 1000)

	response, err := EchoClient(addr, message)
	if err != nil {
		t.Fatalf("EchoClient failed: %v", err)
	}

	expected := "ECHO: " + message
	if response != expected {
		t.Errorf("Message length mismatch: expected %d, got %d", len(expected), len(response))
	}
}

// TestMultipleMessagesOnSameConnection tests keeping a connection open for multiple messages
func TestMultipleMessagesOnSameConnection(t *testing.T) {
	listener, addr := startTestServer(t)
	defer listener.Close()

	// Connect once
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	// Send and receive multiple messages
	for i := 1; i <= 5; i++ {
		message := fmt.Sprintf("Message number %d", i)

		// Send
		fmt.Fprintf(writer, "%s\n", message)
		if err := writer.Flush(); err != nil {
			t.Fatalf("Write failed: %v", err)
		}

		// Receive
		response, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("Read failed: %v", err)
		}

		expected := fmt.Sprintf("ECHO: %s\n", message)
		if response != expected {
			t.Errorf("Expected %q, got %q", expected, response)
		}
	}
}

// TestConnectionClosure tests proper connection cleanup
func TestConnectionClosure(t *testing.T) {
	listener, addr := startTestServer(t)
	defer listener.Close()

	// Connect and send message
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	// Send one message
	fmt.Fprintf(writer, "Test\n")
	writer.Flush()

	// Read response
	reader.ReadString('\n')

	// Close connection
	conn.Close()

	// Try to send again (should fail)
	_, err = fmt.Fprintf(writer, "After close\n")
	// Note: Write might not fail immediately due to buffering
	// The flush or read will fail

	// Wait a bit for connection to fully close
	time.Sleep(100 * time.Millisecond)

	// Verify we can still connect with a new connection
	_, err = EchoClient(addr, "New connection")
	if err != nil {
		t.Errorf("Failed to establish new connection after closing previous: %v", err)
	}
}

// TestInvalidServerAddress tests connecting to invalid address
func TestInvalidServerAddress(t *testing.T) {
	// Try to connect to non-existent server
	_, err := EchoClient("localhost:9999", "Test")
	if err == nil {
		t.Error("Expected error when connecting to non-existent server")
	}
}

// BenchmarkEchoClient benchmarks the echo client performance
func BenchmarkEchoClient(b *testing.B) {
	listener, addr := startTestServer(b)
	defer listener.Close()

	message := "Benchmark test message"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := EchoClient(addr, message)
		if err != nil {
			b.Fatalf("EchoClient failed: %v", err)
		}
	}
}

// BenchmarkConcurrentClients benchmarks concurrent client connections
func BenchmarkConcurrentClients(b *testing.B) {
	listener, addr := startTestServer(b)
	defer listener.Close()

	message := "Concurrent benchmark message"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := EchoClient(addr, message)
			if err != nil {
				b.Fatalf("EchoClient failed: %v", err)
			}
		}
	})
}

// BenchmarkPersistentConnection benchmarks reusing a single connection
func BenchmarkPersistentConnection(b *testing.B) {
	listener, addr := startTestServer(b)
	defer listener.Close()

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		b.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	message := "Persistent connection message"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := SendMessage(conn, message)
		if err != nil {
			b.Fatalf("SendMessage failed: %v", err)
		}
	}
}

// Helper function to start a test echo server
// Returns the listener and the address to connect to
func startTestServer(t testing.TB) (net.Listener, string) {
	t.Helper()

	// Listen on random available port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start test server: %v", err)
	}

	// Get the actual address (with the assigned port)
	addr := listener.Addr().String()

	// Start server in background goroutine
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				// Server stopped
				return
			}
			go handleTestClient(conn)
		}
	}()

	// Wait a bit for server to be ready
	time.Sleep(10 * time.Millisecond)

	return listener, addr
}

// handleTestClient is a simple echo handler for tests
func handleTestClient(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintf(writer, "ECHO: %s\n", line)
		writer.Flush()
	}
}
