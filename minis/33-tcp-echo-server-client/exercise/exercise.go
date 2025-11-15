//go:build !solution
// +build !solution

package exercise

import (
	"bufio"
	"fmt"
	"net"
)

// StartEchoServer starts a TCP echo server on the given address.
// It returns when the server stops or encounters a fatal error.
//
// The server should:
//  1. Listen on the given address (e.g., ":8080")
//  2. Accept multiple concurrent client connections
//  3. For each client, read lines of text and echo them back with "ECHO: " prefix
//  4. Handle errors gracefully
//  5. Close connections properly
//
// Example session:
//
//	Client sends: "Hello\n"
//	Server responds: "ECHO: Hello\n"
//	Client sends: "World\n"
//	Server responds: "ECHO: World\n"
//
// Parameters:
//   - addr: Address to listen on (e.g., ":8080", "127.0.0.1:9000")
//
// Returns:
//   - error: Non-nil if the server fails to start
func StartEchoServer(addr string) error {
	// TODO: Implement the TCP echo server
	//
	// Steps:
	// 1. Create a TCP listener using net.Listen("tcp", addr)
	// 2. Loop forever to accept connections
	// 3. For each connection, spawn a goroutine to handle it
	// 4. In the handler:
	//    a. Use bufio.Scanner to read lines
	//    b. Echo each line back with "ECHO: " prefix
	//    c. Use bufio.Writer and Flush() for efficient writing
	//    d. Handle errors and close the connection when done
	//
	// Hints:
	// - defer listener.Close()
	// - defer conn.Close() in the handler
	// - scanner := bufio.NewScanner(conn)
	// - writer := bufio.NewWriter(conn)
	// - writer.Flush() after each write

	return fmt.Errorf("not implemented")
}

// EchoClient connects to a TCP echo server and sends a single message.
// It returns the server's response.
//
// This function:
//  1. Connects to the server at the given address
//  2. Sends the message (with newline)
//  3. Reads the response line
//  4. Closes the connection
//  5. Returns the response (without the trailing newline)
//
// Parameters:
//   - addr: Server address (e.g., "localhost:8080")
//   - message: Message to send (newline will be added automatically)
//
// Returns:
//   - string: Server's response (without trailing newline)
//   - error: Non-nil if connection or I/O fails
func EchoClient(addr, message string) (string, error) {
	// TODO: Implement the TCP echo client
	//
	// Steps:
	// 1. Connect to the server using net.Dial("tcp", addr)
	// 2. Create bufio.Writer and write the message with newline
	// 3. Flush the writer to ensure data is sent
	// 4. Create bufio.Reader and read the response line
	// 5. Remove the trailing newline from the response
	// 6. Close the connection
	// 7. Return the response
	//
	// Hints:
	// - conn, err := net.Dial("tcp", addr)
	// - defer conn.Close()
	// - writer := bufio.NewWriter(conn)
	// - fmt.Fprintf(writer, "%s\n", message)
	// - writer.Flush()
	// - reader := bufio.NewReader(conn)
	// - response, err := reader.ReadString('\n')
	// - strings.TrimRight(response, "\n")

	return "", fmt.Errorf("not implemented")
}

// SendMessage sends a message to an already-established connection and reads the response.
// This is a helper function for interactive clients.
//
// Parameters:
//   - conn: Established TCP connection
//   - message: Message to send
//
// Returns:
//   - string: Server's response (without trailing newline)
//   - error: Non-nil if I/O fails
func SendMessage(conn net.Conn, message string) (string, error) {
	// TODO: Implement message send/receive
	//
	// Steps:
	// 1. Create a bufio.Writer and write the message with newline
	// 2. Flush the writer
	// 3. Create a bufio.Reader and read the response line
	// 4. Remove trailing newline from response
	// 5. Return the response
	//
	// Hints:
	// - Similar to EchoClient but uses existing connection
	// - Don't close the connection (caller manages it)

	return "", fmt.Errorf("not implemented")
}

// ReadResponse reads a single line response from the connection.
// This is a helper function for reading server responses.
//
// Parameters:
//   - reader: bufio.Reader for the connection
//
// Returns:
//   - string: Response line (without trailing newline)
//   - error: Non-nil if read fails
func ReadResponse(reader *bufio.Reader) (string, error) {
	// TODO: Implement response reading
	//
	// Steps:
	// 1. Read until newline using reader.ReadString('\n')
	// 2. Handle EOF and other errors
	// 3. Remove trailing newline
	// 4. Return the response
	//
	// Hints:
	// - line, err := reader.ReadString('\n')
	// - Check for io.EOF
	// - strings.TrimRight(line, "\n")

	return "", fmt.Errorf("not implemented")
}
