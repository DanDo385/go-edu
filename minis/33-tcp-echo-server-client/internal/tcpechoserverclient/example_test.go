package tcpechoserverclient

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func ExampleEchoClient() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)

		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		r := bufio.NewReader(conn)
		line, err := r.ReadString('\n')
		if err != nil {
			return
		}

		line = strings.TrimRight(line, "\n")
		_, _ = fmt.Fprintf(conn, "ECHO: %s\n", line)
	}()

	resp, err := EchoClient(ln.Addr().String(), "Hello, World!")
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
	<-done

	// Output: ECHO: Hello, World!
}
