package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/21-race-detection-demo/internal/racedetectiondemo/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
