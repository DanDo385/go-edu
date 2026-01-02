package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/33-tcp-echo-server-client/internal/tcpechoserverclient/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
