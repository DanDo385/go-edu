package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/03-keys-addresses/internal/keysaddresses/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
