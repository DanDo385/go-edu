package main

import (
	"os"

	"github.com/example/go-10x-minis/geth/24-monitor/internal/monitor/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
