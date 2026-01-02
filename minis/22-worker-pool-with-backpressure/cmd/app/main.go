package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/22-worker-pool-with-backpressure/internal/workerpoolwithbackpressure/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
