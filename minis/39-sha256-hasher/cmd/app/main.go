package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/39-sha256-hasher/internal/sha256hasher/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
