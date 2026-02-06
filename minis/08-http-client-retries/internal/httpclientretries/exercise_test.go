package httpclientretries

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestRetryClient_Success(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		fmt.Fprintln(w, "success")
	}))
	defer server.Close()

	client := NewRetryClient(3, 10*time.Millisecond, 0.1)

	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if atomic.LoadInt32(&attempts) != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}

func TestRetryClient_RetryAndSucceed(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentAttempt := atomic.AddInt32(&attempts, 1)
		if currentAttempt < 3 {
			w.WriteHeader(http.StatusServiceUnavailable) // 503 is retryable
			return
		}
		fmt.Fprintln(w, "success")
	}))
	defer server.Close()

	client := NewRetryClient(3, 10*time.Millisecond, 0.1)
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Expected success after retries, got error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestRetryClient_NonRetryableError(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusBadRequest) // 400 is not retryable
	}))
	defer server.Close()

	client := NewRetryClient(3, 10*time.Millisecond, 0.1)
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Expected success (with bad status) but got error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", resp.StatusCode)
	}
	if atomic.LoadInt32(&attempts) != 1 {
		t.Errorf("Expected only 1 attempt for non-retryable error, got %d", attempts)
	}
}

func TestRetryClient_MaxRetriesExceeded(t *testing.T) {
	var attempts int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewRetryClient(2, 10*time.Millisecond, 0.1) // Max 2 retries
	resp, err := client.Get(server.URL)
	if err != nil {
		t.Fatalf("Expected success (with bad status) but got error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("Expected status 503, got %d", resp.StatusCode)
	}
	// The number of attempts will be maxRetries + 1 (the initial attempt)
	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("Expected 3 attempts (1 initial + 2 retries), got %d", attempts)
	}
}

func TestRetryClient_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This handler will cause the first request to fail, triggering a retry with a delay.
		// The context will be cancelled during that delay.
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewRetryClient(3, 50*time.Millisecond, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
	_, err := client.Do(req)

	if err == nil {
		t.Fatal("Expected an error due to context cancellation, but got nil")
	}
	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Errorf("Expected context deadline error, got: %v", err)
	}
}

func TestIsRetryable(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		resp     *http.Response
		expected bool
	}{
		{"Network Error", io.EOF, nil, true},
		{"Status 503", nil, &http.Response{StatusCode: 503}, true},
		{"Status 429", nil, &http.Response{StatusCode: 429}, true},
		{"Status 404", nil, &http.Response{StatusCode: 404}, false},
		{"Status 200", nil, &http.Response{StatusCode: 200}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRetryable(tt.err, tt.resp); got != tt.expected {
				t.Errorf("isRetryable() = %v, want %v", got, tt.expected)
			}
		})
	}
}