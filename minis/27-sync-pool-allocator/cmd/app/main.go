package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/27-sync-pool-allocator/internal/syncpoolallocator/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
