package main

import (
	"context"
	"flag"   // Command-line argument parsing
	"fmt"    // Formatted I/O
	"log"    // Logging
	"net/http"
	"time"   // Time operations

	"github.com/example/go-10x-minis/minis/08-http-client-retries/exercise"
)

/*
===================================================================
HTTP Client with Exponential Backoff Retries
===================================================================

This program demonstrates a robust HTTP client with automatic retry
logic using exponential backoff for transient failures.

USAGE:
    go run ./cmd/http-client-retries/main.go
    go run ./cmd/http-client-retries/main.go -url="https://httpbin.org/delay/1"
    go run ./cmd/http-client-retries/main.go -retries=5 -timeout=10s
    go run ./cmd/http-client-retries/main.go -url="https://httpbin.org/status/500" -retries=3
    go run ./cmd/http-client-retries/main.go -delay=200ms -timeout=5s

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging (select "Debug: Run main.go")
- Step Into (F11) to enter exercise package functions
- Watch Variables panel to see retry attempts and backoff delays
- Use Debug Console to inspect HTTP responses
- See /RUN_DEBUG.md for comprehensive debugging guide

RETRY BEHAVIOR:
- Exponential backoff: delay doubles after each retry
- Maximum retries: configurable via -retries flag
- Timeout: per-request timeout via -timeout flag
- Base delay: initial retry delay via -delay flag
*/

type Response struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect default values before parsing
	// DEBUG: Before flag.Parse(), all flags have their default values
	// DEBUG: In Variables panel, you'll see default URL, retries, timeout, delay

	url := flag.String("url", "https://httpbin.org/get", "URL to fetch")
	maxRetries := flag.Int("retries", 3, "Maximum number of retry attempts")
	timeout := flag.Duration("timeout", 5*time.Second, "HTTP request timeout")
	baseDelay := flag.Duration("delay", 100*time.Millisecond, "Base delay for exponential backoff")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")

	flag.Parse()

	// BREAKPOINT: Set breakpoint here to inspect parsed values
	// DEBUG: After flag.Parse(), variables now hold actual command-line values
	// DEBUG: Watch 'url', 'maxRetries', 'timeout', 'baseDelay' - they may have changed
	// DEBUG: Hover over *maxRetries to see the dereferenced int value

	if *verbose {
		log.SetFlags(log.Ltime | log.Lmicroseconds)
	}

	// ============================================================================
	// STEP 2: Create HTTP Client with Retry Logic
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before client creation
	// DEBUG: Watch how the Client struct is configured
	// DEBUG: Note the timeout, maxRetries, and baseDelay settings

	client := &exercise.Client{
		HTTP:       &http.Client{Timeout: *timeout},
		MaxRetries: *maxRetries,
		BaseDelay:  *baseDelay,
	}

	// BREAKPOINT: Set breakpoint here after client creation
	// DEBUG: Inspect 'client' variable in Variables panel
	// DEBUG: Expand client to see HTTP client, MaxRetries, BaseDelay
	// DEBUG: Expand client.HTTP to see the underlying http.Client configuration

	fmt.Println("=== HTTP Client with Retries Demo ===\n")
	fmt.Printf("Configuration:\n")
	fmt.Printf("  URL:         %s\n", *url)
	fmt.Printf("  Max Retries: %d\n", *maxRetries)
	fmt.Printf("  Timeout:     %v\n", *timeout)
	fmt.Printf("  Base Delay:  %v\n", *baseDelay)
	fmt.Println()

	// ============================================================================
	// STEP 3: Create Request Context
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before context creation
	// DEBUG: Context allows us to control cancellation and timeouts
	// DEBUG: Step Into (F11) to see context creation (though it's a built-in)

	ctx := context.Background()

	// BREAKPOINT: Set breakpoint here after context creation
	// DEBUG: Inspect 'ctx' - it's the root context with no deadline
	// DEBUG: Later, the HTTP client will add timeout from client.HTTP.Timeout

	// For demonstration, we can add a deadline to the context
	// Uncomment to test context cancellation:
	// ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel()

	// ============================================================================
	// STEP 4: Make HTTP Request with Automatic Retries
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before HTTP request
	// DEBUG: Step Into (F11) on exercise.GetJSON() to see retry logic
	// DEBUG: Watch how the function handles failures and retries

	fmt.Println("Fetching from URL...")
	startTime := time.Now()

	// BREAKPOINT: Set breakpoint here at the GetJSON call
	// DEBUG: Step Into (F11) to enter the retry logic
	// DEBUG: Watch the retry counter increment on failures
	// DEBUG: Observe exponential backoff delays between retries

	var result map[string]interface{}
	result, err := exercise.GetJSON[map[string]interface{}](ctx, client, *url)

	// BREAKPOINT: Set breakpoint here after GetJSON returns
	// DEBUG: Inspect 'err' - nil on success, error on failure
	// DEBUG: Inspect 'result' - should contain parsed JSON response
	// DEBUG: Check elapsed time to see total duration including retries

	elapsed := time.Since(startTime)

	// ============================================================================
	// STEP 5: Handle Response or Error
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before error check
	// DEBUG: Inspect 'err' to determine if request succeeded

	if err != nil {
		// BREAKPOINT: Set breakpoint here on error path
		// DEBUG: Inspect 'err' to see the error message
		// DEBUG: Error could be: timeout, network error, HTTP error, max retries exceeded
		// DEBUG: Use Call Stack to see where the error originated

		log.Printf("Request failed after %v: %v", elapsed, err)

		// BREAKPOINT: Set breakpoint here after logging error
		// DEBUG: Check if all retries were exhausted
		// DEBUG: Look at the error message for retry count info
	} else {
		// BREAKPOINT: Set breakpoint here on success path
		// DEBUG: Request succeeded (possibly after retries)
		// DEBUG: Inspect 'result' to see the parsed JSON response
		// DEBUG: Expand result to see all fields

		fmt.Printf("Success after %v:\n", elapsed)
		fmt.Printf("%+v\n", result)

		// BREAKPOINT: Set breakpoint here after printing result
		// DEBUG: Verify the response data is what you expected
	}

	// ============================================================================
	// STEP 6: Additional Test Cases
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before additional tests
	// DEBUG: The following demonstrates different scenarios

	fmt.Println("\n--- Testing Different Scenarios ---\n")

	// Test 1: Delayed response (should succeed but take time)
	// BREAKPOINT: Set breakpoint here before delayed request
	// DEBUG: Step Into (F11) to see how client handles slow responses

	testURLs := []struct {
		url         string
		description string
	}{
		{"https://httpbin.org/delay/1", "Slow response (1s delay)"},
		{"https://httpbin.org/status/200", "Success response"},
		{"https://httpbin.org/status/500", "Server error (will retry)"},
		{"https://httpbin.org/status/429", "Rate limited (will retry)"},
	}

	for i, test := range testURLs {
		// BREAKPOINT: Set breakpoint here inside test loop
		// DEBUG: Watch 'i' and 'test' change with each iteration
		// DEBUG: Step Into (F11) to see how each scenario is handled

		fmt.Printf("%d. %s\n", i+1, test.description)
		fmt.Printf("   URL: %s\n", test.url)

		// BREAKPOINT: Set breakpoint here before each test request
		// DEBUG: Note the test URL and expected behavior
		// DEBUG: Step Into (F11) to see retry logic for failures

		startTime := time.Now()
		var testResult map[string]interface{}
		testResult, err := exercise.GetJSON[map[string]interface{}](ctx, client, test.url)
		elapsed := time.Since(startTime)

		// BREAKPOINT: Set breakpoint here after each test request
		// DEBUG: Check if request succeeded or failed
		// DEBUG: For failures, observe retry behavior
		// DEBUG: For successes, check response time

		if err != nil {
			fmt.Printf("   Result: FAILED after %v - %v\n", elapsed, err)

			// BREAKPOINT: Set breakpoint here on test failure
			// DEBUG: Inspect the error type and message
			// DEBUG: Verify retries were attempted
		} else {
			fmt.Printf("   Result: SUCCESS after %v\n", elapsed)
			if *verbose {
				fmt.Printf("   Data: %+v\n", testResult)
			}

			// BREAKPOINT: Set breakpoint here on test success
			// DEBUG: Verify response data is valid
		}
		fmt.Println()

		// BREAKPOINT: Set breakpoint here between tests
		// DEBUG: Pause between requests to avoid rate limiting
		time.Sleep(200 * time.Millisecond)
	}

	// ============================================================================
	// DEBUGGING TIPS
	// ============================================================================
	//
	// KEY DEBUGGING TECHNIQUES:
	//
	// 1. Set breakpoints at "// BREAKPOINT:" comments
	// 2. Press F5 to start debugging
	// 3. Use F10 (Step Over) to execute line by line
	// 4. Use F11 (Step Into) to enter exercise package functions
	// 5. Watch Variables panel to see retry attempts and delays
	// 6. Add watch expressions: err, result, elapsed
	// 7. Use Call Stack panel to see function hierarchy
	//
	// HTTP CLIENT DEBUGGING:
	// - Step Into exercise.GetJSON() to see retry loop
	// - Watch retry counter increment on failures
	// - Observe exponential backoff: delay doubles each retry
	// - Check context for cancellation/timeout
	// - Inspect HTTP response status codes
	//
	// RETRY LOGIC DEBUGGING:
	// - First attempt uses no delay
	// - Second attempt: baseDelay (e.g., 100ms)
	// - Third attempt: 2 * baseDelay (e.g., 200ms)
	// - Fourth attempt: 4 * baseDelay (e.g., 400ms)
	// - Continue until maxRetries or success
	//
	// TESTING FAILURES:
	// - Use httpbin.org/status/500 to test error retries
	// - Use httpbin.org/delay/N to test timeout behavior
	// - Use httpbin.org/status/429 to test rate limiting
	// - Set low timeout to force timeout errors
	//
	// CONTEXT DEBUGGING:
	// - Context carries deadline and cancellation
	// - Use ctx.Deadline() in Debug Console
	// - Use ctx.Err() to check if context is cancelled
	// - Context timeout is separate from HTTP client timeout
	//
	// See /RUN_DEBUG.md for comprehensive debugging guide
}
