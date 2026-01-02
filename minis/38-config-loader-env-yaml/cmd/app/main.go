package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/38-config-loader-env-yaml/internal/configloaderenvyaml/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
