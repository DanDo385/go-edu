package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/19-channels-basics/internal/channelsbasics/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
