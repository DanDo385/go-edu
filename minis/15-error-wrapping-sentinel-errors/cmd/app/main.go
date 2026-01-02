package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/15-error-wrapping-sentinel-errors/internal/errorwrappingsentinelerrors/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
