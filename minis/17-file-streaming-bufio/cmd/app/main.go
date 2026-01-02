package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/17-file-streaming-bufio/internal/filestreamingbufio/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
