//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package tcpechoserverclient

import (
	"bufio"

	"net"
)
// TODO: implement StartEchoServer.
func StartEchoServer(addr string) error { panic("TODO: implement") }
// TODO: implement handleClient.
func handleClient(conn net.Conn) { panic("TODO: implement") }
// TODO: implement EchoClient.
func EchoClient(addr, message string) (string, error) { panic("TODO: implement") }
// TODO: implement SendMessage.
func SendMessage(conn net.Conn, message string) (string, error) { panic("TODO: implement") }
// TODO: implement ReadResponse.
func ReadResponse(reader *bufio.Reader) (string, error) { panic("TODO: implement") }
