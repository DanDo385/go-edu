package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/28-pprof-cpu-mem-benchmarks/internal/pprofcpumembenchmarks/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
