package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/37-http-middleware-chain/internal/httpmiddlewarechain/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
