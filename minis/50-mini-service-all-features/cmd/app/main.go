package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/50-mini-service-all-features/internal/miniserviceallfeatures/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
