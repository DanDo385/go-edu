package exercise

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// =============================================================================
// Context Tests
// =============================================================================

func TestContextHelpers(t *testing.T) {
	ctx := context.Background()

	t.Run("RequestID", func(t *testing.T) {
		// Test WithRequestID and GetRequestID
		requestID := "test-request-id-123"
		ctx = WithRequestID(ctx, requestID)

		got, ok := GetRequestID(ctx)
		if !ok {
			t.Fatal("GetRequestID returned false, expected true")
		}

		if got != requestID {
			t.Errorf("GetRequestID() = %q, want %q", got, requestID)
		}
	})

	t.Run("User", func(t *testing.T) {
		// Test WithUser and GetUser
		user := &User{ID: 1, Name: "Alice"}
		ctx = WithUser(ctx, user)

		got, ok := GetUser(ctx)
		if !ok {
			t.Fatal("GetUser returned false, expected true")
		}

		if got.ID != user.ID || got.Name != user.Name {
			t.Errorf("GetUser() = %+v, want %+v", got, user)
		}
	})

	t.Run("Missing values", func(t *testing.T) {
		emptyCtx := context.Background()

		if _, ok := GetRequestID(emptyCtx); ok {
			t.Error("GetRequestID should return false for empty context")
		}

		if _, ok := GetUser(emptyCtx); ok {
			t.Error("GetUser should return false for empty context")
		}
	})
}

// =============================================================================
// ResponseWriter Tests
// =============================================================================

func TestResponseWriter(t *testing.T) {
	t.Run("captures status code", func(t *testing.T) {
		w := httptest.NewRecorder()
		rw := NewResponseWriter(w)

		rw.WriteHeader(http.StatusNotFound)

		if rw.StatusCode() != http.StatusNotFound {
			t.Errorf("StatusCode() = %d, want %d", rw.StatusCode(), http.StatusNotFound)
		}
	})

	t.Run("default status code is 200", func(t *testing.T) {
		w := httptest.NewRecorder()
		rw := NewResponseWriter(w)

		rw.Write([]byte("test"))

		if rw.StatusCode() != http.StatusOK {
			t.Errorf("StatusCode() = %d, want %d", rw.StatusCode(), http.StatusOK)
		}
	})

	t.Run("counts bytes written", func(t *testing.T) {
		w := httptest.NewRecorder()
		rw := NewResponseWriter(w)

		data := []byte("Hello, World!")
		rw.Write(data)

		if rw.BytesWritten() != len(data) {
			t.Errorf("BytesWritten() = %d, want %d", rw.BytesWritten(), len(data))
		}
	})

	t.Run("accumulates bytes across multiple writes", func(t *testing.T) {
		w := httptest.NewRecorder()
		rw := NewResponseWriter(w)

		rw.Write([]byte("Hello, "))
		rw.Write([]byte("World!"))

		expected := len("Hello, World!")
		if rw.BytesWritten() != expected {
			t.Errorf("BytesWritten() = %d, want %d", rw.BytesWritten(), expected)
		}
	})
}

// =============================================================================
// Middleware Tests
// =============================================================================

func TestLoggingMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	middleware := LoggingMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
	}

	body := w.Body.String()
	if body != "test response" {
		t.Errorf("Body = %q, want %q", body, "test response")
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	middleware := RecoveryMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Should not panic
	middleware.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestRequestIDMiddleware(t *testing.T) {
	var capturedRequestID string

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := GetRequestID(r.Context())
		if !ok {
			t.Error("Request ID not found in context")
		}
		capturedRequestID = id
		w.WriteHeader(http.StatusOK)
	})

	middleware := RequestIDMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	middleware.ServeHTTP(w, req)

	// Check X-Request-ID header
	requestID := w.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Error("X-Request-ID header not set")
	}

	if capturedRequestID != requestID {
		t.Errorf("Context request ID = %q, want %q", capturedRequestID, requestID)
	}
}

func TestAuthMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUser(r.Context())
		if !ok {
			t.Error("User not found in context")
		}
		if user == nil {
			t.Error("User is nil")
		}
		w.WriteHeader(http.StatusOK)
	})

	middleware := AuthMiddleware(handler)

	t.Run("valid token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusUnauthorized)
		}
	})

	t.Run("missing token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusUnauthorized)
		}
	})
}

func TestCORSMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := CORSMiddleware("*")(handler)

	t.Run("adds CORS headers", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		if origin := w.Header().Get("Access-Control-Allow-Origin"); origin != "*" {
			t.Errorf("Access-Control-Allow-Origin = %q, want %q", origin, "*")
		}

		methods := w.Header().Get("Access-Control-Allow-Methods")
		if methods == "" {
			t.Error("Access-Control-Allow-Methods header not set")
		}

		headers := w.Header().Get("Access-Control-Allow-Headers")
		if headers == "" {
			t.Error("Access-Control-Allow-Headers header not set")
		}
	})

	t.Run("handles OPTIONS preflight", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/test", nil)
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
		}
	})
}

func TestMethodMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := MethodMiddleware("GET", "POST")(handler)

	t.Run("allows GET", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("allows POST", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/test", nil)
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("rejects PUT", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/test", nil)
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusMethodNotAllowed)
		}
	})

	t.Run("rejects DELETE", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/test", nil)
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusMethodNotAllowed)
		}
	})
}

// =============================================================================
// Chain Tests
// =============================================================================

func TestChain(t *testing.T) {
	var executionOrder []string

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		executionOrder = append(executionOrder, "handler")
		w.WriteHeader(http.StatusOK)
	})

	middlewareA := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "A-before")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "A-after")
		})
	}

	middlewareB := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "B-before")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "B-after")
		})
	}

	middlewareC := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "C-before")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "C-after")
		})
	}

	// Chain: A → B → C → handler
	chained := Chain(handler, middlewareA, middlewareB, middlewareC)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	chained.ServeHTTP(w, req)

	expected := []string{
		"A-before",
		"B-before",
		"C-before",
		"handler",
		"C-after",
		"B-after",
		"A-after",
	}

	if len(executionOrder) != len(expected) {
		t.Fatalf("Execution order length = %d, want %d", len(executionOrder), len(expected))
	}

	for i, step := range expected {
		if executionOrder[i] != step {
			t.Errorf("Step %d: got %q, want %q", i, executionOrder[i], step)
		}
	}
}

func TestChainWithShortCircuit(t *testing.T) {
	var executed []string

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		executed = append(executed, "handler")
		w.WriteHeader(http.StatusOK)
	})

	middlewareA := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executed = append(executed, "A-before")
			next.ServeHTTP(w, r)
			executed = append(executed, "A-after")
		})
	}

	// This middleware short-circuits (doesn't call next)
	middlewareB := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executed = append(executed, "B-before")
			w.WriteHeader(http.StatusUnauthorized)
			// Not calling next.ServeHTTP!
			executed = append(executed, "B-after")
		})
	}

	chained := Chain(handler, middlewareA, middlewareB)

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	chained.ServeHTTP(w, req)

	// Should execute: A-before → B-before → B-after → A-after
	// Handler should NOT execute
	expected := []string{
		"A-before",
		"B-before",
		"B-after",
		"A-after",
	}

	if len(executed) != len(expected) {
		t.Fatalf("Execution length = %d, want %d\nGot: %v", len(executed), len(expected), executed)
	}

	for i, step := range expected {
		if executed[i] != step {
			t.Errorf("Step %d: got %q, want %q", i, executed[i], step)
		}
	}

	// Verify handler was not executed
	for _, step := range executed {
		if step == "handler" {
			t.Error("Handler should not have been executed due to short-circuit")
		}
	}

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Status code = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

// =============================================================================
// Integration Tests
// =============================================================================

func TestFullMiddlewareStack(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request ID is in context
		if _, ok := GetRequestID(r.Context()); !ok {
			t.Error("Request ID not in context")
		}

		// Verify user is in context
		if _, ok := GetUser(r.Context()); !ok {
			t.Error("User not in context")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	// Build full middleware stack
	fullStack := Chain(
		handler,
		RecoveryMiddleware,
		LoggingMiddleware,
		RequestIDMiddleware,
		AuthMiddleware,
		CORSMiddleware("*"),
		MethodMiddleware("GET", "POST"),
	)

	t.Run("successful request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()

		fullStack.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
		}

		if body := w.Body.String(); body != "success" {
			t.Errorf("Body = %q, want %q", body, "success")
		}

		// Verify CORS header
		if origin := w.Header().Get("Access-Control-Allow-Origin"); origin != "*" {
			t.Errorf("CORS header = %q, want %q", origin, "*")
		}

		// Verify Request ID header
		if requestID := w.Header().Get("X-Request-ID"); requestID == "" {
			t.Error("X-Request-ID header not set")
		}
	})

	t.Run("unauthorized request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		// No Authorization header
		w := httptest.NewRecorder()

		fullStack.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusUnauthorized)
		}
	})

	t.Run("method not allowed", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/test", nil)
		req.Header.Set("Authorization", "Bearer valid-token")
		w := httptest.NewRecorder()

		fullStack.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusMethodNotAllowed)
		}
	})
}

func TestMiddlewareComposition(t *testing.T) {
	t.Run("multiple middleware layers", func(t *testing.T) {
		counter := 0

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter++
			w.WriteHeader(http.StatusOK)
		})

		// Apply same middleware multiple times (for testing)
		stacked := Chain(
			handler,
			LoggingMiddleware,
			LoggingMiddleware,
			LoggingMiddleware,
		)

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		stacked.ServeHTTP(w, req)

		if counter != 1 {
			t.Errorf("Handler executed %d times, want 1", counter)
		}

		if w.Code != http.StatusOK {
			t.Errorf("Status code = %d, want %d", w.Code, http.StatusOK)
		}
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
