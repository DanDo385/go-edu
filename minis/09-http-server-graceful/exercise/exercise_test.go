package exercise

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestKVHandler_PostAndGet(t *testing.T) {
	store := NewMemStore()
	mux := http.NewServeMux()
	RegisterRoutes(mux, store)

	// POST
	body := bytes.NewBufferString(`{"key":"name","val":"Go"}`)
	req := httptest.NewRequest("POST", "/kv", body)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", w.Code)
	}

	// GET
	req = httptest.NewRequest("GET", "/kv?k=name", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["val"] != "Go" {
		t.Errorf("Expected val=Go, got %q", resp["val"])
	}
}

func TestKVHandler_NotFound(t *testing.T) {
	store := NewMemStore()
	mux := http.NewServeMux()
	RegisterRoutes(mux, store)

	req := httptest.NewRequest("GET", "/kv?k=nonexistent", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}

func TestMiddleware_RequestCount(t *testing.T) {
	store := NewMemStore()
	mux := http.NewServeMux()
	RegisterRoutes(mux, store)

	for i := 1; i <= 3; i++ {
		req := httptest.NewRequest("GET", "/kv?k=test", nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)

		count := w.Header().Get("X-Req-Count")
		if count == "" {
			t.Error("Expected X-Req-Count header")
		}
	}
}
