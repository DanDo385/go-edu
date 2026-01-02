package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/09-http-server-graceful/internal/httpservergraceful/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
