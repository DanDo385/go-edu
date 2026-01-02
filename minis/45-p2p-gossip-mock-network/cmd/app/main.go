package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/45-p2p-gossip-mock-network/internal/p2pgossipmocknetwork/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
