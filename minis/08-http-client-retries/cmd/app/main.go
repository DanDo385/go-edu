package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/08-http-client-retries/internal/httpclientretries/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
