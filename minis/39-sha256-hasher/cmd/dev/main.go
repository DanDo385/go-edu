package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/example/go-10x-minis/minis/39-sha256-hasher/internal/sha256hasher/cli"
)

func main() {
	fmt.Println("=== minis/39-sha256-hasher: cmd/dev ===")
	fmt.Println()

	for _, ex := range cli.Examples() {
		fmt.Printf("$ go run ./cmd/app %s\n", strings.Join(ex.Args, " "))
		code := cli.RunCLI(ex.Args)
		if code != 0 {
			fmt.Printf("-> exited with code %d\n", code)
		}
		fmt.Println()
	}

	// If someone passes args manually to cmd/dev, forward them too.
	if len(os.Args) > 1 {
		fmt.Printf("$ (forwarded) go run ./cmd/app %s\n", strings.Join(os.Args[1:], " "))
		os.Exit(cli.RunCLI(os.Args[1:]))
	}
}
