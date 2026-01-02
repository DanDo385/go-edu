package main

import (
	"os"

	"github.com/example/go-10x-minis/minis/10-grpc-telemetry-service/internal/grpctelemetryservice/cli"
)

func main() {
	os.Exit(cli.RunCLI(os.Args[1:]))
}
