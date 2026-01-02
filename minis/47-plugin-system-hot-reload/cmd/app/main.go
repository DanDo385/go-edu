package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/47-plugin-system-hot-reload/internal/pluginsystemhotreload/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
