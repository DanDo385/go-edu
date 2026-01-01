package main

import (
	"fmt"

	"github.com/example/go-10x-minis/minis/39-sha256-hasher/internal/sha256hasher"
)

func main() {
	fmt.Println("=== Debug Harness ===")
	fmt.Println("Fixed inputs for debugging - perfect for stepping through code.")
	fmt.Println()
	
	// TODO: Add debug harness code using the sha256hasher package
	// This file uses fixed, deterministic inputs - no CLI arguments needed!
	// Perfect for setting breakpoints and stepping through logic.
	
	_ = sha256hasher
	
	fmt.Println("Debug harness complete.")
}
