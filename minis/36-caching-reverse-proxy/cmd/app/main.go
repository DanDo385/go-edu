package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/36-caching-reverse-proxy/internal/cachingreverseproxy/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
