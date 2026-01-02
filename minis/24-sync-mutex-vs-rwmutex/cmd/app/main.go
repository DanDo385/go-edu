package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/24-sync-mutex-vs-rwmutex/internal/syncmutexvsrwmutex/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
