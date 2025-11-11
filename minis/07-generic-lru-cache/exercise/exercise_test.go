package exercise

import (
	"testing"
	"time"
)

func TestCache_BasicOperations(t *testing.T) {
	cache := New[string, int](3, 0)

	// Set and Get
	cache.Set("a", 1)
	if val, ok := cache.Get("a"); !ok || val != 1 {
		t.Errorf("Expected a=1, got val=%d, ok=%v", val, ok)
	}

	// Get non-existent
	if _, ok := cache.Get("nonexistent"); ok {
		t.Error("Expected false for non-existent key")
	}
}

func TestCache_Eviction(t *testing.T) {
	cache := New[string, int](2, 0) // Capacity 2

	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3) // Should evict "a"

	if _, ok := cache.Get("a"); ok {
		t.Error("Expected 'a' to be evicted")
	}
	if val, ok := cache.Get("b"); !ok || val != 2 {
		t.Error("Expected 'b' to still exist")
	}
	if val, ok := cache.Get("c"); !ok || val != 3 {
		t.Error("Expected 'c' to exist")
	}
}

func TestCache_LRUOrder(t *testing.T) {
	cache := New[string, int](2, 0)

	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Get("a") // Access "a", making it most recent
	cache.Set("c", 3) // Should evict "b", not "a"

	if val, ok := cache.Get("a"); !ok || val != 1 {
		t.Error("Expected 'a' to still exist (was accessed recently)")
	}
	if _, ok := cache.Get("b"); ok {
		t.Error("Expected 'b' to be evicted (least recently used)")
	}
}

func TestCache_Update(t *testing.T) {
	cache := New[string, int](3, 0)

	cache.Set("a", 1)
	cache.Set("a", 10) // Update

	if val, ok := cache.Get("a"); !ok || val != 10 {
		t.Errorf("Expected a=10 after update, got %d", val)
	}

	if cache.Len() != 1 {
		t.Errorf("Expected Len()=1, got %d", cache.Len())
	}
}

func TestCache_TTL(t *testing.T) {
	cache := New[string, int](3, 100*time.Millisecond)

	cache.Set("a", 1)
	time.Sleep(150 * time.Millisecond)

	if _, ok := cache.Get("a"); ok {
		t.Error("Expected 'a' to be expired")
	}
}

func TestCache_CustomTTL(t *testing.T) {
	cache := New[string, int](3, 1*time.Hour) // Default: 1 hour

	cache.SetWithTTL("a", 1, 100*time.Millisecond) // Custom: 100ms
	cache.Set("b", 2) // Uses default TTL

	time.Sleep(150 * time.Millisecond)

	if _, ok := cache.Get("a"); ok {
		t.Error("Expected 'a' to be expired (custom TTL)")
	}
	if _, ok := cache.Get("b"); !ok {
		t.Error("Expected 'b' to still exist (default TTL)")
	}
}

func TestCache_ZeroTTL(t *testing.T) {
	cache := New[string, int](3, 0) // No expiration

	cache.Set("a", 1)
	time.Sleep(100 * time.Millisecond)

	if val, ok := cache.Get("a"); !ok || val != 1 {
		t.Error("Expected 'a' to not expire (zero TTL)")
	}
}

func TestCache_Len(t *testing.T) {
	cache := New[string, int](3, 0)

	if cache.Len() != 0 {
		t.Errorf("Expected initial Len()=0, got %d", cache.Len())
	}

	cache.Set("a", 1)
	cache.Set("b", 2)

	if cache.Len() != 2 {
		t.Errorf("Expected Len()=2, got %d", cache.Len())
	}

	cache.Set("c", 3)
	cache.Set("d", 4) // Evicts one

	if cache.Len() != 3 {
		t.Errorf("Expected Len()=3 (capacity), got %d", cache.Len())
	}
}
