package workerpoolwordcount

import (
	"context"           // Context for cancellation
	"fmt"               // String formatting
	"net/http"          // HTTP server
	"net/http/httptest" // Test HTTP server
	"testing"           // Testing framework
	"time"              // Time durations
)

// TestWordCount_Basic tests basic functionality with two URLs
func TestWordCount_Basic(t *testing.T) {
	// Create test servers (mock HTTP servers for testing)
	servers := []*httptest.Server{
		httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "hello world") // Server 1 returns "hello world"
		})),
		httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "hello go") // Server 2 returns "hello go"
		})),
	}
	defer closeServers(servers) // Cleanup: close servers when test finishes

	urls := []string{servers[0].URL, servers[1].URL} // Get URLs from test servers
	ctx := context.Background()                      // Create background context

	counts, err := WordCount(ctx, urls, 2) // Call WordCount with 2 workers
	if err != nil {                        // Check for errors
		t.Fatalf("WordCount failed: %v", err) // Fail test if error occurred
	}

	// Verify word counts
	if counts["hello"] != 2 { // "hello" appears in both responses
		t.Errorf("Expected 'hello' count=2, got %d", counts["hello"])
	}
	if counts["world"] != 1 { // "world" appears once
		t.Errorf("Expected 'world' count=1, got %d", counts["world"])
	}
	if counts["go"] != 1 { // "go" appears once
		t.Errorf("Expected 'go' count=1, got %d", counts["go"])
	}
}

// TestWordCount_EmptyURLs tests behavior with no URLs
func TestWordCount_EmptyURLs(t *testing.T) {
	ctx := context.Background()                  // Create context
	counts, err := WordCount(ctx, []string{}, 2) // Call with empty URL list
	if err != nil {                              // Check error
		t.Fatalf("WordCount failed: %v", err) // Fail if error
	}
	if len(counts) != 0 { // Should return empty map
		t.Errorf("Expected empty map, got %d words", len(counts))
	}
}

// TestWordCount_ContextCancellation tests timeout behavior
func TestWordCount_ContextCancellation(t *testing.T) {
	// Create slow server (takes 2 seconds to respond)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Simulate slow response
		fmt.Fprintln(w, "slow response")
	}))
	defer server.Close() // Cleanup

	// Create context with 100ms timeout (should timeout before server responds)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel() // Cleanup

	_, err := WordCount(ctx, []string{server.URL}, 1) // Call with timeout context
	if err == nil {                                   // Should return error
		t.Error("Expected error due to context timeout") // Fail if no error
	}
}

// TestWordCount_ServerError tests error handling for HTTP errors
func TestWordCount_ServerError(t *testing.T) {
	// Create server that returns 500 error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError) // Return 500 status
	}))
	defer server.Close() // Cleanup

	ctx := context.Background()                       // Create context
	_, err := WordCount(ctx, []string{server.URL}, 1) // Call WordCount
	if err == nil {                                   // Should return error
		t.Error("Expected error for 500 status") // Fail if no error
	}
}

// TestWordCount_Punctuation tests that punctuation is removed
func TestWordCount_Punctuation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world! How are you?") // Text with punctuation
	}))
	defer server.Close() // Cleanup

	ctx := context.Background()                            // Create context
	counts, err := WordCount(ctx, []string{server.URL}, 1) // Call WordCount
	if err != nil {                                        // Check error
		t.Fatalf("WordCount failed: %v", err) // Fail if error
	}

	// Verify punctuation removed (should find "hello" not "Hello,")
	if counts["hello"] != 1 { // Should be lowercase, no comma
		t.Errorf("Expected 'hello' (no comma) count=1, got %d", counts["hello"])
	}
	if counts["world"] != 1 { // Should be lowercase, no exclamation
		t.Errorf("Expected 'world' (no exclamation) count=1, got %d", counts["world"])
	}
}

// TestWordCount_CaseInsensitive tests case-insensitive counting
func TestWordCount_CaseInsensitive(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Go GO go") // Mixed case
	}))
	defer server.Close() // Cleanup

	ctx := context.Background()                            // Create context
	counts, err := WordCount(ctx, []string{server.URL}, 1) // Call WordCount
	if err != nil {                                        // Check error
		t.Fatalf("WordCount failed: %v", err) // Fail if error
	}

	// All should be counted as "go" (case-insensitive)
	if counts["go"] != 3 { // Should count all three
		t.Errorf("Expected 'go' count=3 (case-insensitive), got %d", counts["go"])
	}
}

// TestWordCount_MultipleWorkers tests with multiple workers and URLs
func TestWordCount_MultipleWorkers(t *testing.T) {
	// Create 10 test servers
	servers := make([]*httptest.Server, 10) // Slice of 10 servers
	for i := range servers {                // Loop over indices
		i := i // Capture loop variable (avoid closure bug)
		servers[i] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "server %d response", i) // Each server returns unique text
		}))
	}
	defer closeServers(servers) // Cleanup all servers

	// Create URLs slice
	urls := make([]string, len(servers)) // Slice of URLs
	for i, srv := range servers {        // Range over servers
		urls[i] = srv.URL // Store server URL
	}

	ctx := context.Background()            // Create context
	counts, err := WordCount(ctx, urls, 3) // Call with 3 workers, 10 URLs
	if err != nil {                        // Check error
		t.Fatalf("WordCount failed: %v", err) // Fail if error
	}

	// Verify all URLs processed (each contains "server" and "response")
	if counts["server"] != 10 { // Should appear 10 times
		t.Errorf("Expected 'server' count=10, got %d", counts["server"])
	}
	if counts["response"] != 10 { // Should appear 10 times
		t.Errorf("Expected 'response' count=10, got %d", counts["response"])
	}
}

// closeServers closes all test servers (helper function)
func closeServers(servers []*httptest.Server) {
	for _, srv := range servers { // Range over servers
		srv.Close() // Close each server
	}
}
