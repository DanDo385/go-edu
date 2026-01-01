//go:build !solution && !reference

package tcpechoserverclient

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

// StartEchoServer implements the exercise.
//
// TODO: Implement this function
func StartEchoServer(addr string) error {
	// TODO: Implement
	return nil
}

// EchoClient implements the exercise.
//
// TODO: Implement this function
func EchoClient(addr string, message string) (string, error) {
	// TODO: Implement
	return "", nil
}

// SendMessage implements the exercise.
//
// TODO: Implement this function
func SendMessage(conn net.Conn, message string) (string, error) {
	// TODO: Implement
	return "", nil
}

// ReadResponse implements the exercise.
//
// TODO: Implement this function
func ReadResponse(reader *bufio.Reader) (string, error) {
	// TODO: Implement
	return "", nil
}
