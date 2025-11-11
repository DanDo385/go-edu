package exercise

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWordCount_Basic(t *testing.T) {
	// Create test servers
	servers := []*httptest.Server{
		httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "hello world")
		})),
		httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "hello go")
		})),
	}
	defer closeServers(servers)

	urls := []string{servers[0].URL, servers[1].URL}
	ctx := context.Background()

	counts, err := WordCount(ctx, urls, 2)
	if err != nil {
		t.Fatalf("WordCount failed: %v", err)
	}

	// Check expected words
	if counts["hello"] != 2 {
		t.Errorf("Expected 'hello' count=2, got %d", counts["hello"])
	}
	if counts["world"] != 1 {
		t.Errorf("Expected 'world' count=1, got %d", counts["world"])
	}
	if counts["go"] != 1 {
		t.Errorf("Expected 'go' count=1, got %d", counts["go"])
	}
}

func TestWordCount_EmptyURLs(t *testing.T) {
	ctx := context.Background()
	counts, err := WordCount(ctx, []string{}, 2)
	if err != nil {
		t.Fatalf("WordCount failed: %v", err)
	}
	if len(counts) != 0 {
		t.Errorf("Expected empty map, got %d words", len(counts))
	}
}

func TestWordCount_ContextCancellation(t *testing.T) {
	// Create a slow server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		fmt.Fprintln(w, "slow response")
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := WordCount(ctx, []string{server.URL}, 1)
	if err == nil {
		t.Error("Expected error due to context timeout")
	}
}

func TestWordCount_ServerError(t *testing.T) {
	// Server returns 500 error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	ctx := context.Background()
	_, err := WordCount(ctx, []string{server.URL}, 1)
	if err == nil {
		t.Error("Expected error for 500 status")
	}
}

func TestWordCount_Punctuation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world! How are you?")
	}))
	defer server.Close()

	ctx := context.Background()
	counts, err := WordCount(ctx, []string{server.URL}, 1)
	if err != nil {
		t.Fatalf("WordCount failed: %v", err)
	}

	// Punctuation should be removed
	if counts["hello"] != 1 {
		t.Errorf("Expected 'hello' (no comma) count=1, got %d", counts["hello"])
	}
	if counts["world"] != 1 {
		t.Errorf("Expected 'world' (no exclamation) count=1, got %d", counts["world"])
	}
}

func TestWordCount_CaseInsensitive(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Go GO go")
	}))
	defer server.Close()

	ctx := context.Background()
	counts, err := WordCount(ctx, []string{server.URL}, 1)
	if err != nil {
		t.Fatalf("WordCount failed: %v", err)
	}

	// All should be counted as "go"
	if counts["go"] != 3 {
		t.Errorf("Expected 'go' count=3 (case-insensitive), got %d", counts["go"])
	}
}

func TestWordCount_MultipleWorkers(t *testing.T) {
	// Create 10 servers
	servers := make([]*httptest.Server, 10)
	for i := range servers {
		i := i // Capture loop variable
		servers[i] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "server %d response", i)
		}))
	}
	defer closeServers(servers)

	urls := make([]string, len(servers))
	for i, srv := range servers {
		urls[i] = srv.URL
	}

	ctx := context.Background()
	counts, err := WordCount(ctx, urls, 3) // 3 workers for 10 URLs
	if err != nil {
		t.Fatalf("WordCount failed: %v", err)
	}

	// Should have "server" and "response" from all
	if counts["server"] != 10 {
		t.Errorf("Expected 'server' count=10, got %d", counts["server"])
	}
	if counts["response"] != 10 {
		t.Errorf("Expected 'response' count=10, got %d", counts["response"])
	}
}

func closeServers(servers []*httptest.Server) {
	for _, srv := range servers {
		srv.Close()
	}
}
