// Simple debug program to test your functions directly
// without running the full test suite.
//
// This file uses fixed default values - no CLI arguments needed!
// Perfect for quick debugging sessions.
//
// To use:
// 1. Set breakpoints in ../../internal/exercise/httpclientretries.go
// 2. Open this file (cmd/debug/main.go)
// 3. Use "Debug Main Program (Current Package)" configuration
// 4. Press F5 - that's it! The debugger will stop at your breakpoints
//
// Usage:
//   go run ./cmd/debug
//   # Or just press F5 in VS Code

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/example/go-10x-minis/minis/08-http-client-retries/internal/httpclientretries"
)

func main() {
	// Fixed default values - modify these directly if you want to test different inputs
	maxRetries := 3
	baseDelay := 100 * time.Millisecond
	url := "https://httpbin.org/get"

	fmt.Println("=== Debugging HTTP Client with Retries ===")
	fmt.Printf("MaxRetries: %d, BaseDelay: %v, URL: %s\n\n", maxRetries, baseDelay, url)

	// Set breakpoint in httpclientretries.go at Client creation and GetJSON function
	ctx := context.Background()
	client := &httpclientretries.Client{
		HTTP:       nil, // Will be set by your implementation
		MaxRetries: maxRetries,
		BaseDelay:  baseDelay,
	}

	type Response struct {
		Message string `json:"message"`
	}

	var result Response
	err := httpclientretries.GetJSON[Response](ctx, client, url)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Result: %+v\n", result)
}

