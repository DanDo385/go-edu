package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Example struct {
	Name string
	Args []string
}


func Examples() []Example {
	return []Example{
	{Name: "example-1", Args: []string{"-rpc", "https://eth.llamarpc.com"}},
}
}

func RunCLI(args []string) int {
	fs := flag.NewFlagSet("app", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	list := fs.Bool("list", false, "list available demo targets")
	fn := fs.String("fn", "", "demo target name")
	rpc := fs.String("rpc", getenvFirst("RPC_URL", "ETH_RPC_URL", "INFURA_RPC_URL", "ALCHEMY_RPC_URL", "https://eth.llamarpc.com"), "Ethereum JSON-RPC URL")
	timeout := fs.String("timeout", "30s", "timeout (informational)")
	jsonOut := fs.Bool("json", false, "print result as JSON (informational)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err)
		return 2
	}

	_ = rpc
	_ = timeout
	_ = jsonOut

	if *list || *fn == "" {
			fmt.Println("This module uses flags like -rpc/-timeout; see README.md for copy/paste examples.")
		if *fn == "" {
			return 0
		}
	}

	fmt.Printf("Selected: %s\n", *fn)
	fmt.Println("\nNote: this CLI is a learning harness.")
	fmt.Println("Implement the TODOs in internal/.../exercise.go and run: go test ./...")
	fmt.Println("\nArgs:")
	fmt.Printf("  %s\n", strings.Join(args, " "))
	return 0
}

func getenvFirst(keys ...string) string {
	for _, k := range keys {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}
