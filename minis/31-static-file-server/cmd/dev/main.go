package main

import (
	"fmt"

	"github.com/example/go-10x-minis/minis/31-static-file-server/internal/staticfileserver"
)

func main() {
	fmt.Println("=== Debug Harness ===")
	fmt.Println("Fixed inputs for debugging - perfect for stepping through code.")
	fmt.Println()
	
	// TODO: Add debug harness code using the staticfileserver package
	// This file uses fixed, deterministic inputs - no CLI arguments needed!
	// Perfect for setting breakpoints and stepping through logic.
	
	_ = staticfileserver
	
	fmt.Println("Debug harness complete.")
}
