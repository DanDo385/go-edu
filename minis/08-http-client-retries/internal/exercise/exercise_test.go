package exercise

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetJSON_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"message":"success"}`)
	}))
	defer server.Close()

	client := &Client{
		HTTP:       &http.Client{},
		MaxRetries: 3,
		BaseDelay:  10 * time.Millisecond,
	}

	type Response struct {
		Message string `json:"message"`
	}

	result, err := GetJSON[Response](context.Background(), client, server.URL)
	if err != nil {
		t.Fatalf("Expected success, got error: %v", err)
	}
	if result.Message != "success" {
		t.Errorf("Expected message='success', got %q", result.Message)
	}
}

func TestGetJSON_RetrySuccess(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"message":"success"}`)
	}))
	defer server.Close()

	client := &Client{
		HTTP:       &http.Client{},
		MaxRetries: 3,
		BaseDelay:  10 * time.Millisecond,
	}

	type Response struct {
		Message string `json:"message"`
	}

	result, err := GetJSON[Response](context.Background(), client, server.URL)
	if err != nil {
		t.Fatalf("Expected success after retries, got error: %v", err)
	}
	if result.Message != "success" {
		t.Errorf("Expected message='success', got %q", result.Message)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestGetJSON_AllRetriesFail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := &Client{
		HTTP:       &http.Client{},
		MaxRetries: 2,
		BaseDelay:  10 * time.Millisecond,
	}

	type Response struct {
		Message string `json:"message"`
	}

	_, err := GetJSON[Response](context.Background(), client, server.URL)
	if err == nil {
		t.Error("Expected error after all retries fail")
	}
}

func TestGetJSON_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		fmt.Fprintln(w, `{"message":"success"}`)
	}))
	defer server.Close()

	client := &Client{
		HTTP:       &http.Client{},
		MaxRetries: 3,
		BaseDelay:  10 * time.Millisecond,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	type Response struct {
		Message string `json:"message"`
	}

	_, err := GetJSON[Response](ctx, client, server.URL)
	if err == nil {
		t.Error("Expected context timeout error")
	}
}
