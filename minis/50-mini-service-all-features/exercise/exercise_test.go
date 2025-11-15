package exercise

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := NewCache()

	// Test Set and Get
	cache.Set("key1", "value1", 1*time.Second)
	val, ok := cache.Get("key1")
	if !ok {
		t.Fatal("expected key1 to exist")
	}
	if val != "value1" {
		t.Errorf("expected value1, got %v", val)
	}

	// Test expiration
	cache.Set("key2", "value2", 100*time.Millisecond)
	time.Sleep(200 * time.Millisecond)
	_, ok = cache.Get("key2")
	if ok {
		t.Error("expected key2 to be expired")
	}

	// Test Delete
	cache.Set("key3", "value3", 1*time.Hour)
	cache.Delete("key3")
	_, ok = cache.Get("key3")
	if ok {
		t.Error("expected key3 to be deleted")
	}
}

func TestCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker(3, 1*time.Second)

	// Test closed circuit (should allow calls)
	err := cb.Call(func() error { return nil })
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test opening circuit after failures
	for i := 0; i < 3; i++ {
		cb.Call(func() error { return errors.New("fail") })
	}

	if cb.State() != StateOpen {
		t.Errorf("expected circuit to be open, got %v", cb.State())
	}

	// Test open circuit (should reject calls immediately)
	err = cb.Call(func() error { return nil })
	if err == nil {
		t.Error("expected circuit open error")
	}
}

func TestRetryWithBackoff(t *testing.T) {
	attempts := 0
	err := RetryWithBackoff(
		context.Background(),
		3,
		10*time.Millisecond,
		func() error {
			attempts++
			if attempts < 3 {
				return errors.New("temporary error")
			}
			return nil
		},
	)

	if err != nil {
		t.Errorf("expected success after retries, got %v", err)
	}

	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestRetryWithBackoffMaxRetries(t *testing.T) {
	attempts := 0
	err := RetryWithBackoff(
		context.Background(),
		3,
		10*time.Millisecond,
		func() error {
			attempts++
			return errors.New("permanent error")
		},
	)

	if err == nil {
		t.Error("expected error after max retries")
	}

	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestUserRateLimiter(t *testing.T) {
	limiter := NewUserRateLimiter(10, 1) // 10 req/sec, burst 1

	// First request should succeed
	if !limiter.Allow("user1") {
		t.Error("expected first request to be allowed")
	}

	// Immediate second request should fail (no burst)
	if limiter.Allow("user1") {
		t.Error("expected second immediate request to be denied")
	}

	// Different user should have separate limit
	if !limiter.Allow("user2") {
		t.Error("expected user2 first request to be allowed")
	}
}

func TestAppError(t *testing.T) {
	err := NewNotFoundError("user not found")
	if err.Error() == "" {
		t.Error("expected error message")
	}

	err = NewBadRequestError("invalid input")
	if err.Error() == "" {
		t.Error("expected error message")
	}

	err = NewInternalError("database error")
	if err.Error() == "" {
		t.Error("expected error message")
	}
}

func TestEventBus(t *testing.T) {
	bus := NewEventBus()

	received := make(chan string, 1)

	// Subscribe to events
	bus.Subscribe("user.created", func(event interface{}) {
		received <- event.(string)
	})

	// Publish event
	bus.Publish("user.created", "user123")

	// Wait for event
	select {
	case msg := <-received:
		if msg != "user123" {
			t.Errorf("expected user123, got %s", msg)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("timeout waiting for event")
	}
}

func TestWorkerPool(t *testing.T) {
	pool := NewWorkerPool(3)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	pool.Start(ctx)

	completed := 0
	for i := 0; i < 10; i++ {
		pool.Submit(func() error {
			completed++
			return nil
		})
	}

	pool.Shutdown()

	if completed != 10 {
		t.Errorf("expected 10 completed jobs, got %d", completed)
	}
}
