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

func StartEchoServer(addr string) error {
	// TODO: Implement this function
	// TODO: Handle errors appropriately
	// TODO: Add necessary validations
	panic("not implemented")
}

func handleClient(conn net.Conn) {
	// TODO: Implement this function
	panic("not implemented")
}

func EchoClient(addr, message string) (string, error) {
	// TODO: Implement this function
	panic("not implemented")
}

func SendMessage(conn net.Conn, message string) (string, error) {
	// TODO: Implement this function
	panic("not implemented")
}

func ReadResponse(reader *bufio.Reader) (string, error) {
	// TODO: Implement this function
	panic("not implemented")
}
