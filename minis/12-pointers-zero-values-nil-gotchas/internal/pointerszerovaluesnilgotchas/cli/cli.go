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

var available = []string{}

func Examples() []Example {
	return []Example{
	{Name: "example-1", Args: []string{"-list"}},
}
}

func RunCLI(args []string) int {
	fs := flag.NewFlagSet("app", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	list := fs.Bool("list", false, "list available demo targets")
	fn := fs.String("fn", "", "demo target name")
	in := fs.String("in", "", "string input (for demos)")
	n := fs.Int("n", 5, "int input (for demos)")
	f := fs.Float64("f", 12.34, "float64 input (for demos)")
	b := fs.Bool("b", true, "bool input (for demos)")
	file := fs.String("file", "", "file input (for demos)")
	stdin := fs.Bool("stdin", false, "stdin input (for demos)")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "parse error:", err)
		return 2
	}

	_ = in
	_ = n
	_ = f
	_ = b
	_ = file
	_ = stdin

	if *list || *fn == "" {
			if len(available) > 0 {
				fmt.Println("Available exported functions:")
				for _, n := range available {
					fmt.Println(" -", n)
				}
			} else {
				fmt.Println("No exported functions detected for this module.")
			}
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
