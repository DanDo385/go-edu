//go:build !solution && !reference


package tcpechoserverclient

import (
	"bufio"
	"fmt"
	"net"
)

// StartEchoServer starts a TCP echo server on the given address.
func StartEchoServer(addr string) error {
	// TODO: Implement the TCP echo server.

	// Step 1: Create a TCP listener on the given address.
	// - `listener, err := net.Listen("tcp", addr)`
	// - Handle the error if the server can't start.
	// - `defer listener.Close()` to ensure the listener is closed when the function exits.

	// Step 2: Create an infinite loop to accept new connections.
	// - `for { ... }`
	// - `conn, err := listener.Accept()`
	// - `listener.Accept()` will block until a new client connects.
	// - If `Accept` returns an error, log it and `continue` to the next iteration to keep the server running.

	// Step 3: Handle each connection concurrently.
	// - For each successful connection, launch a new goroutine to handle it: `go handleClient(conn)`.
	// - This allows the server to handle many clients at the same time.

	// Step 4: Implement the `handleClient` function (you can create this as a helper).
	// - `func handleClient(conn net.Conn) { ... }`
	// - Inside `handleClient`:
	//   - `defer conn.Close()` to ensure the client connection is closed.
	//   - Create a `bufio.Scanner` to read lines from the connection: `scanner := bufio.NewScanner(conn)`.
	//   - Create a `bufio.Writer` for writing responses: `writer := bufio.NewWriter(conn)`.
	//   - Loop through the scanner: `for scanner.Scan()`.
	//   - Get the received line: `line := scanner.Text()`.
	//   - Write the response back to the client: `fmt.Fprintf(writer, "ECHO: %s\n", line)`.
	//   - **Crucially**, flush the writer's buffer to send the data: `writer.Flush()`.
	//   - After the loop, check `scanner.Err()` for any read errors.

	return fmt.Errorf("not implemented")
}

// EchoClient connects to a TCP echo server and sends a single message.
func EchoClient(addr, message string) (string, error) {
	// TODO: Implement the TCP echo client.

	// Step 1: Connect to the server.
	// - `conn, err := net.Dial("tcp", addr)`
	// - Handle connection errors.
	// - `defer conn.Close()` to ensure the connection is closed.

	// Step 2: Send the message.
	// - Create a `bufio.Writer` for the connection.
	// - Write the message, making sure to add a newline `\n` at the end, as the server reads line-by-line.
	//   - `fmt.Fprintf(writer, "%s\n", message)` is a good way to do this.
	// - Flush the writer to ensure the data is sent over the network: `writer.Flush()`.

	// Step 3: Read the response.
	// - Create a `bufio.Reader` for the connection.
	// - Read a single line from the server: `response, err := reader.ReadString('\n')`.
	// - Handle any read errors.

	// Step 4: Clean up and return the response.
	// - The response will include the trailing newline. Remove it using `strings.TrimRight(response, "\n")`.
	// - Return the cleaned response.

	return "", fmt.Errorf("not implemented")
}

// SendMessage sends a message to an already-established connection and reads the response.
func SendMessage(conn net.Conn, message string) (string, error) {
	// TODO: Implement message send/receive on an existing connection.

	// This function is very similar to `EchoClient`, but it assumes the connection is already established and does not close it.

	// Step 1: Write the message.
	// - Create a `bufio.Writer`, write the message with a newline, and flush.

	// Step 2: Read the response.
	// - Create a `bufio.Reader` and call the `ReadResponse` helper function you will also implement.
	return "", fmt.Errorf("not implemented")
}

// ReadResponse reads a single line response from the connection.
func ReadResponse(reader *bufio.Reader) (string, error) {
	// TODO: Implement a single line read.

	// Step 1: Read from the reader until a newline.
	// - `line, err := reader.ReadString('\n')`

	// Step 2: Handle errors.
	// - Pay special attention to `io.EOF`, which means the server closed the connection.
	// - For other errors, just return them.

	// Step 3: Clean up and return the line.
	// - `return strings.TrimRight(line, "\n"), nil`
	return "", fmt.Errorf("not implemented")
}
