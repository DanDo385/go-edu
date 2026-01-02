package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/23-bounded-channel-semaphore/internal/boundedchannelsemaphore/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
