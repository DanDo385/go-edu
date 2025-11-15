package exercise

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
	"time"
)

// Test Exercise 1: Basic Buffer Pool
func TestBufferPool(t *testing.T) {
	pool := NewBufferPool()

	// Test Get returns a buffer
	buf1 := pool.Get()
	if buf1 == nil {
		t.Fatal("Get() returned nil")
	}

	// Test buffer can be used
	buf1.WriteString("test data")
	if buf1.String() != "test data" {
		t.Errorf("Expected 'test data', got '%s'", buf1.String())
	}

	// Test Put and reuse
	pool.Put(buf1)

	buf2 := pool.Get()
	if buf2 == nil {
		t.Fatal("Get() returned nil after Put")
	}

	// Buffer should be reset (empty)
	if buf2.Len() != 0 {
		t.Errorf("Expected empty buffer, got length %d", buf2.Len())
	}

	// Should ideally be the same buffer (but not guaranteed)
	// We can't test pointer equality reliably due to GC
}

func TestBufferPoolConcurrent(t *testing.T) {
	pool := NewBufferPool()

	var wg sync.WaitGroup
	const goroutines = 100
	const iterations = 1000

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				buf := pool.Get()
				buf.WriteString(fmt.Sprintf("goroutine %d iteration %d", id, j))
				pool.Put(buf)
			}
		}(i)
	}

	wg.Wait()
}

// Test Exercise 2: Slice Pool
func TestSlicePool(t *testing.T) {
	const capacity = 1024
	pool := NewSlicePool(capacity)

	// Test Get returns a slice
	slice1 := pool.Get()
	if slice1 == nil {
		t.Fatal("Get() returned nil")
	}

	// Test capacity
	if cap(*slice1) < capacity {
		t.Errorf("Expected capacity >= %d, got %d", capacity, cap(*slice1))
	}

	// Test slice can be used
	*slice1 = append(*slice1, []byte("test")...)
	if len(*slice1) != 4 {
		t.Errorf("Expected length 4, got %d", len(*slice1))
	}

	// Test Put and reuse
	pool.Put(slice1)

	slice2 := pool.Get()
	if slice2 == nil {
		t.Fatal("Get() returned nil after Put")
	}

	// Slice should be reset (length 0)
	if len(*slice2) != 0 {
		t.Errorf("Expected empty slice, got length %d", len(*slice2))
	}

	// Capacity should be preserved
	if cap(*slice2) < capacity {
		t.Errorf("Expected capacity >= %d, got %d", capacity, cap(*slice2))
	}
}

// Test Exercise 3: Generic Typed Pool
func TestGenericPool(t *testing.T) {
	type TestStruct struct {
		Data  []int
		Value string
	}

	resetFunc := func(s *TestStruct) {
		s.Data = s.Data[:0]
		s.Value = ""
	}

	pool := NewPool(
		func() *TestStruct {
			return &TestStruct{
				Data: make([]int, 0, 10),
			}
		},
		resetFunc,
	)

	// Test Get returns an object
	obj1 := pool.Get()
	if obj1 == nil {
		t.Fatal("Get() returned nil")
	}

	// Test object can be used
	obj1.Data = append(obj1.Data, 1, 2, 3)
	obj1.Value = "test"

	if len(obj1.Data) != 3 {
		t.Errorf("Expected Data length 3, got %d", len(obj1.Data))
	}

	// Test Put and reset
	pool.Put(obj1)

	obj2 := pool.Get()
	if obj2 == nil {
		t.Fatal("Get() returned nil after Put")
	}

	// Object should be reset
	if len(obj2.Data) != 0 {
		t.Errorf("Expected empty Data, got length %d", len(obj2.Data))
	}
	if obj2.Value != "" {
		t.Errorf("Expected empty Value, got '%s'", obj2.Value)
	}
}

// Test Exercise 4: Pool with Metrics
func TestMetricsPool(t *testing.T) {
	pool := NewMetricsPool(func() interface{} {
		return new(bytes.Buffer)
	})

	// First Get should create new object
	obj1 := pool.Get()
	if obj1 == nil {
		t.Fatal("Get() returned nil")
	}

	stats1 := pool.Stats()
	if stats1.Gets != 1 {
		t.Errorf("Expected 1 Get, got %d", stats1.Gets)
	}
	if stats1.News != 1 {
		t.Errorf("Expected 1 New, got %d", stats1.News)
	}

	// Put and Get again should reuse
	pool.Put(obj1)

	stats2 := pool.Stats()
	if stats2.Puts != 1 {
		t.Errorf("Expected 1 Put, got %d", stats2.Puts)
	}

	obj2 := pool.Get()
	if obj2 == nil {
		t.Fatal("Get() returned nil after Put")
	}

	stats3 := pool.Stats()
	if stats3.Gets != 2 {
		t.Errorf("Expected 2 Gets, got %d", stats3.Gets)
	}
	if stats3.News != 1 {
		t.Errorf("Expected 1 New (reused), got %d", stats3.News)
	}

	// Hit rate should be 50% (1 hit, 1 miss)
	if stats3.HitRate < 40 || stats3.HitRate > 60 {
		t.Errorf("Expected hit rate around 50%%, got %.1f%%", stats3.HitRate)
	}
}

func TestMetricsPoolHitRate(t *testing.T) {
	pool := NewMetricsPool(func() interface{} {
		return new(bytes.Buffer)
	})

	// Get 100 objects, put them back, get them again
	objects := make([]interface{}, 100)
	for i := 0; i < 100; i++ {
		objects[i] = pool.Get()
	}
	for i := 0; i < 100; i++ {
		pool.Put(objects[i])
	}
	for i := 0; i < 100; i++ {
		objects[i] = pool.Get()
	}

	stats := pool.Stats()
	if stats.Gets != 200 {
		t.Errorf("Expected 200 Gets, got %d", stats.Gets)
	}

	// Hit rate should be ~50% (100 initial creates, 100 reuses)
	if stats.HitRate < 40 || stats.HitRate > 60 {
		t.Errorf("Expected hit rate around 50%%, got %.1f%%", stats.HitRate)
	}
}

// Test Exercise 5: Size-Classed Buffer Pool
func TestSizeClassedPool(t *testing.T) {
	pool := NewSizeClassedPool()

	tests := []struct {
		requestSize  int
		minCapacity  int
		description  string
	}{
		{500, 1024, "small request should get 1KB buffer"},
		{2000, 4096, "medium request should get 4KB buffer"},
		{10000, 16384, "large request should get 16KB buffer"},
		{50000, 65536, "huge request should get 64KB buffer"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			buf := pool.Get(tt.requestSize)
			if buf == nil {
				t.Fatal("Get() returned nil")
			}

			if cap(*buf) < tt.minCapacity {
				t.Errorf("Expected capacity >= %d, got %d", tt.minCapacity, cap(*buf))
			}

			// Test Put
			pool.Put(buf)
		})
	}
}

func TestSizeClassedPoolReuse(t *testing.T) {
	pool := NewSizeClassedPool()

	// Get and put buffers of different sizes
	buf1 := pool.Get(2000)  // 4KB pool
	pool.Put(buf1)

	buf2 := pool.Get(3000)  // Should reuse from 4KB pool
	if cap(*buf2) < 4096 {
		t.Errorf("Expected reused buffer with capacity >= 4096, got %d", cap(*buf2))
	}
}

// Test Exercise 6: Bounded Pool
func TestBoundedPool(t *testing.T) {
	const maxSize = 5
	pool := NewBoundedPool(maxSize, func() interface{} {
		return new(bytes.Buffer)
	})

	// Get maxSize objects
	objects := make([]interface{}, maxSize)
	for i := 0; i < maxSize; i++ {
		objects[i] = pool.Get()
		if objects[i] == nil {
			t.Fatalf("Get() returned nil at iteration %d", i)
		}
	}

	// InUse should be maxSize
	if pool.InUse() != maxSize {
		t.Errorf("Expected InUse=%d, got %d", maxSize, pool.InUse())
	}

	// Test that Get blocks when limit reached (test with timeout)
	done := make(chan struct{})
	go func() {
		_ = pool.Get() // This should block
		close(done)
	}()

	// Give it a moment to block
	select {
	case <-done:
		t.Error("Get() should have blocked when limit reached")
	case <-time.After(100 * time.Millisecond):
		// Good, it blocked
	}

	// Put one back
	pool.Put(objects[0])

	// Now the blocked Get should succeed
	select {
	case <-done:
		// Good, it unblocked
	case <-time.After(1 * time.Second):
		t.Error("Get() should have unblocked after Put()")
	}
}

// Test Exercise 8: Worker Pool
func TestWorkerPool(t *testing.T) {
	pool := NewWorkerPool()

	// Test basic processing
	result := pool.Process("test data")
	if result == "" {
		t.Error("Process() returned empty string")
	}

	// Test concurrent processing
	var wg sync.WaitGroup
	const goroutines = 10
	const iterations = 100

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				data := fmt.Sprintf("worker %d iteration %d", id, j)
				result := pool.Process(data)
				if result == "" {
					t.Errorf("Process() returned empty string for %s", data)
				}
			}
		}(i)
	}

	wg.Wait()
}

// Benchmarks for Exercise 7

func BenchmarkWithoutPool(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf := new(bytes.Buffer)
		buf.WriteString("test data for benchmarking")
		buf.WriteString(" with multiple writes")
		buf.WriteString(" to simulate real usage")
		_ = buf.String()
	}
}

func BenchmarkWithPool(b *testing.B) {
	pool := NewBufferPool()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		buf := pool.Get()
		buf.WriteString("test data for benchmarking")
		buf.WriteString(" with multiple writes")
		buf.WriteString(" to simulate real usage")
		_ = buf.String()
		pool.Put(buf)
	}
}

func BenchmarkParallelWithoutPool(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := new(bytes.Buffer)
			buf.WriteString("test data for benchmarking")
			buf.WriteString(" with multiple writes")
			buf.WriteString(" to simulate real usage")
			_ = buf.String()
		}
	})
}

func BenchmarkParallelWithPool(b *testing.B) {
	pool := NewBufferPool()
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := pool.Get()
			buf.WriteString("test data for benchmarking")
			buf.WriteString(" with multiple writes")
			buf.WriteString(" to simulate real usage")
			_ = buf.String()
			pool.Put(buf)
		}
	})
}

func BenchmarkSliceWithoutPool(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		slice := make([]byte, 0, 1024)
		slice = append(slice, []byte("test data")...)
		slice = append(slice, []byte(" more data")...)
		_ = slice
	}
}

func BenchmarkSliceWithPool(b *testing.B) {
	pool := NewSlicePool(1024)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		slice := pool.Get()
		*slice = append(*slice, []byte("test data")...)
		*slice = append(*slice, []byte(" more data")...)
		_ = *slice
		pool.Put(slice)
	}
}

func BenchmarkSizeClassedPool(b *testing.B) {
	pool := NewSizeClassedPool()
	b.ReportAllocs()
	b.ResetTimer()

	sizes := []int{500, 2000, 10000, 50000}
	for i := 0; i < b.N; i++ {
		size := sizes[i%len(sizes)]
		buf := pool.Get(size)
		*buf = append(*buf, make([]byte, size)...)
		pool.Put(buf)
	}
}
