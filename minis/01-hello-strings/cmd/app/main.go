package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/01-hello-strings/internal/hellostrings/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
