package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/26-sync-once-singleton/internal/synconcesingleton/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
