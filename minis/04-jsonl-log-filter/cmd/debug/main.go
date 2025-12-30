// Simple debug program to test your functions directly
// without running the full test suite.
//
// This file uses fixed default values - no CLI arguments needed!
// Perfect for quick debugging sessions.
//
// To use:
// 1. Set breakpoints in ../../internal/exercise/exercise.go
// 2. Open this file (cmd/debug/main.go)
// 3. Use "Debug Main Program (Current Package)" configuration
// 4. Press F5 - that's it! The debugger will stop at your breakpoints
//
// Usage:
//   go run ./cmd/debug
//   # Or just press F5 in VS Code

package main

import (
	"fmt"
	"strings"

	"github.com/example/go-10x-minis/minis/04-jsonl-log-filter/internal/exercise"
)

func main() {
	// Fixed default values - modify these directly if you want to test different inputs
	level := exercise.Warn
	jsonlData := `{"level":"info","message":"Application started"}
{"level":"error","message":"Failed to connect"}
{"level":"debug","message":"Processing request"}
{"level":"warn","message":"Warning message"}`

	fmt.Println("=== Debugging FilterLogs ===")
	fmt.Printf("Filter level: %v\n", level)
	fmt.Printf("JSONL Data:\n%s\n\n", jsonlData)

	reader := strings.NewReader(jsonlData)

	// Set breakpoint in exercise.go at FilterLogs function
	entries, err := exercise.FilterLogs(reader, level)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Filtered entries: %+v\n", entries)
}

