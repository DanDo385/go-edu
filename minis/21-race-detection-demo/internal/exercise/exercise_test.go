package exercise

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// ============================================================================
// Test 1: Safe Counter
// ============================================================================

func TestSafeCounter(t *testing.T) {
	counter := NewSafeCounter()

	const numGoroutines = 10
	const incrementsPerGoroutine = 1000
	const expected = numGoroutines * incrementsPerGoroutine

	var wg sync.WaitGroup
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				counter.Increment()
			}
		}()
	}

	wg.Wait()

	if got := counter.Value(); got != expected {
		t.Errorf("SafeCounter = %d, want %d (race detected!)", got, expected)
	}
}

func TestSafeCounterConcurrentRead(t *testing.T) {
	counter := NewSafeCounter()

	// Writer goroutines
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				counter.Increment()
			}
		}()
	}

	// Reader goroutines (should not race)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				_ = counter.Value()
			}
		}()
	}

	wg.Wait()

	if got := counter.Value(); got != 5000 {
		t.Errorf("SafeCounter = %d, want 5000", got)
	}
}

// ============================================================================
// Test 2: Safe Map
// ============================================================================

func TestSafeMap(t *testing.T) {
	m := NewSafeMap()

	const numGoroutines = 10
	const opsPerGoroutine = 100

	var wg sync.WaitGroup

	// Concurrent writers
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				m.Set(key, id*1000+j)
			}
		}(i)
	}

	// Concurrent readers
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				// Read may or may not find the key (depends on timing)
				m.Get(key)
			}
		}(i)
	}

	wg.Wait()

	expectedLen := numGoroutines * opsPerGoroutine
	if got := m.Len(); got != expectedLen {
		t.Errorf("SafeMap.Len() = %d, want %d", got, expectedLen)
	}
}

func TestSafeMapGetSet(t *testing.T) {
	m := NewSafeMap()

	m.Set("foo", 42)
	m.Set("bar", 99)

	if val, ok := m.Get("foo"); !ok || val != 42 {
		t.Errorf("Get(foo) = (%d, %v), want (42, true)", val, ok)
	}

	if val, ok := m.Get("bar"); !ok || val != 99 {
		t.Errorf("Get(bar) = (%d, %v), want (99, true)", val, ok)
	}

	if _, ok := m.Get("missing"); ok {
		t.Errorf("Get(missing) should return false")
	}

	if got := m.Len(); got != 2 {
		t.Errorf("Len() = %d, want 2", got)
	}
}

// ============================================================================
// Test 3: Lazy Initialization
// ============================================================================

func TestLazyInit(t *testing.T) {
	l := NewLazyInit()

	callCount := 0
	var mu sync.Mutex

	init := func() interface{} {
		mu.Lock()
		callCount++
		mu.Unlock()
		time.Sleep(10 * time.Millisecond) // Simulate slow initialization
		return "initialized"
	}

	const numGoroutines = 10
	var wg sync.WaitGroup
	results := make([]interface{}, numGoroutines)

	// Call GetOrInit concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			results[id] = l.GetOrInit(init)
		}(i)
	}

	wg.Wait()

	// Check that init was called exactly once
	if callCount != 1 {
		t.Errorf("init function called %d times, want 1 (not thread-safe!)", callCount)
	}

	// Check that all goroutines got the same value
	for i, result := range results {
		if result != "initialized" {
			t.Errorf("result[%d] = %v, want 'initialized'", i, result)
		}
	}
}

// ============================================================================
// Test 4: Safe Slice
// ============================================================================

func TestSafeSlice(t *testing.T) {
	s := NewSafeSlice()

	const numGoroutines = 10
	const appendsPerGoroutine = 100

	var wg sync.WaitGroup

	// Concurrent appends
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < appendsPerGoroutine; j++ {
				s.Append(id*1000 + j)
			}
		}(i)
	}

	wg.Wait()

	expectedLen := numGoroutines * appendsPerGoroutine
	if got := s.Len(); got != expectedLen {
		t.Errorf("SafeSlice.Len() = %d, want %d (race detected!)", got, expectedLen)
	}
}

func TestSafeSliceGet(t *testing.T) {
	s := NewSafeSlice()

	s.Append(10)
	s.Append(20)
	s.Append(30)

	if val, ok := s.Get(0); !ok || val != 10 {
		t.Errorf("Get(0) = (%d, %v), want (10, true)", val, ok)
	}

	if val, ok := s.Get(1); !ok || val != 20 {
		t.Errorf("Get(1) = (%d, %v), want (20, true)", val, ok)
	}

	if val, ok := s.Get(2); !ok || val != 30 {
		t.Errorf("Get(2) = (%d, %v), want (30, true)", val, ok)
	}

	if _, ok := s.Get(3); ok {
		t.Errorf("Get(3) should return false for out of bounds")
	}

	if got := s.Len(); got != 3 {
		t.Errorf("Len() = %d, want 3", got)
	}
}

// ============================================================================
// Test 5: Process IDs (Loop Variable Capture)
// ============================================================================

func TestProcessIDs(t *testing.T) {
	ids := []int{1, 2, 3, 4, 5}

	// Simple square function
	square := func(n int) int {
		return n * n
	}

	results := ProcessIDs(ids, square)

	// Check length
	if len(results) != len(ids) {
		t.Fatalf("len(results) = %d, want %d", len(results), len(ids))
	}

	// Check that all results are correct
	expected := []int{1, 4, 9, 16, 25}
	for i, want := range expected {
		if results[i] != want {
			t.Errorf("results[%d] = %d, want %d (loop variable capture bug?)", i, results[i], want)
		}
	}
}

// ============================================================================
// Test 6: URL Cache
// ============================================================================

func TestURLCache(t *testing.T) {
	fetchCount := make(map[string]int)
	var mu sync.Mutex

	fetcher := func(url string) (string, error) {
		mu.Lock()
		fetchCount[url]++
		mu.Unlock()
		time.Sleep(10 * time.Millisecond) // Simulate network delay
		return "content of " + url, nil
	}

	cache := NewURLCache(fetcher)

	urls := []string{
		"http://example.com",
		"http://example.com", // Duplicate
		"http://google.com",
		"http://example.com", // Duplicate
		"http://github.com",
	}

	var wg sync.WaitGroup
	results := make([]string, len(urls))

	// Fetch concurrently
	for i, url := range urls {
		wg.Add(1)
		go func(idx int, u string) {
			defer wg.Done()
			content, err := cache.Fetch(u)
			if err != nil {
				t.Errorf("Fetch(%s) error: %v", u, err)
				return
			}
			results[idx] = content
		}(i, url)
	}

	wg.Wait()

	// Check that example.com was fetched only once
	mu.Lock()
	exampleCount := fetchCount["http://example.com"]
	mu.Unlock()

	if exampleCount != 1 {
		t.Errorf("example.com fetched %d times, want 1 (caching not working)", exampleCount)
	}

	// Check results
	expectedResults := []string{
		"content of http://example.com",
		"content of http://example.com",
		"content of http://google.com",
		"content of http://example.com",
		"content of http://github.com",
	}

	for i, want := range expectedResults {
		if results[i] != want {
			t.Errorf("results[%d] = %s, want %s", i, results[i], want)
		}
	}
}

// ============================================================================
// Test 7: Concurrent Metrics
// ============================================================================

func TestMetrics(t *testing.T) {
	m := NewMetrics()

	const numGoroutines = 10
	const requestsPerGoroutine = 1000
	const errorRate = 10 // 1 error every 10 requests

	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < requestsPerGoroutine; j++ {
				m.IncrementRequests()
				if j%errorRate == 0 {
					m.IncrementErrors()
				}
			}
		}()
	}

	wg.Wait()

	requests, errors := m.GetStats()

	expectedRequests := int64(numGoroutines * requestsPerGoroutine)
	expectedErrors := int64(numGoroutines * (requestsPerGoroutine / errorRate))

	if requests != expectedRequests {
		t.Errorf("requests = %d, want %d (race detected!)", requests, expectedRequests)
	}

	if errors != expectedErrors {
		t.Errorf("errors = %d, want %d (race detected!)", errors, expectedErrors)
	}
}

// ============================================================================
// Test 8: Bank Account
// ============================================================================

func TestBankAccount(t *testing.T) {
	account := NewBankAccount(1000)

	const numDepositors = 5
	const numWithdrawers = 5
	const opsPerGoroutine = 100
	const depositAmount = 10
	const withdrawAmount = 10

	var wg sync.WaitGroup

	// Depositors
	for i := 0; i < numDepositors; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				account.Deposit(depositAmount)
			}
		}()
	}

	// Withdrawers
	for i := 0; i < numWithdrawers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				account.Withdraw(withdrawAmount)
			}
		}()
	}

	wg.Wait()

	// Initial: 1000
	// Deposits: 5 * 100 * 10 = 5000
	// Withdrawals: 5 * 100 * 10 = 5000
	// Final: 1000 + 5000 - 5000 = 1000
	expected := int64(1000)

	if got := account.Balance(); got != expected {
		t.Errorf("Balance() = %d, want %d (race detected!)", got, expected)
	}
}

func TestBankAccountInsufficientFunds(t *testing.T) {
	account := NewBankAccount(100)

	// Try to withdraw more than balance
	if account.Withdraw(150) {
		t.Error("Withdraw(150) should return false (insufficient funds)")
	}

	// Balance should be unchanged
	if got := account.Balance(); got != 100 {
		t.Errorf("Balance() = %d, want 100", got)
	}

	// Withdraw exact amount should succeed
	if !account.Withdraw(100) {
		t.Error("Withdraw(100) should return true")
	}

	if got := account.Balance(); got != 0 {
		t.Errorf("Balance() = %d, want 0", got)
	}
}

// ============================================================================
// Test 9: Pipeline Pattern
// ============================================================================

func TestPipeline(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}

	// Expected: 1^2 + 2^2 + 3^2 + 4^2 + 5^2 = 1 + 4 + 9 + 16 + 25 = 55
	expected := 55

	got := Pipeline(numbers)

	if got != expected {
		t.Errorf("Pipeline(%v) = %d, want %d", numbers, got, expected)
	}
}

func TestPipelineEmpty(t *testing.T) {
	numbers := []int{}
	expected := 0

	got := Pipeline(numbers)

	if got != expected {
		t.Errorf("Pipeline(%v) = %d, want %d", numbers, got, expected)
	}
}

func TestPipelineLarge(t *testing.T) {
	numbers := make([]int, 100)
	expected := 0
	for i := 0; i < 100; i++ {
		numbers[i] = i + 1
		expected += (i + 1) * (i + 1)
	}

	got := Pipeline(numbers)

	if got != expected {
		t.Errorf("Pipeline(1..100) = %d, want %d", got, expected)
	}
}

// ============================================================================
// Test 10: Worker Pool
// ============================================================================

func TestWorkerPool(t *testing.T) {
	jobs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	square := func(n int) int {
		return n * n
	}

	results := WorkerPool(3, jobs, square)

	// Check length
	if len(results) != len(jobs) {
		t.Fatalf("len(results) = %d, want %d", len(results), len(jobs))
	}

	// Check that all results are present (order doesn't matter)
	expected := make(map[int]bool)
	for _, n := range jobs {
		expected[n*n] = true
	}

	for _, result := range results {
		if !expected[result] {
			t.Errorf("unexpected result: %d", result)
		}
		delete(expected, result)
	}

	if len(expected) > 0 {
		t.Errorf("missing results: %v", expected)
	}
}

func TestWorkerPoolLarge(t *testing.T) {
	jobs := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		jobs[i] = i
	}

	double := func(n int) int {
		return n * 2
	}

	results := WorkerPool(10, jobs, double)

	if len(results) != 1000 {
		t.Fatalf("len(results) = %d, want 1000", len(results))
	}

	// Check that all results are present
	expected := make(map[int]bool)
	for i := 0; i < 1000; i++ {
		expected[i*2] = true
	}

	for _, result := range results {
		if !expected[result] {
			t.Errorf("unexpected result: %d", result)
		}
		delete(expected, result)
	}

	if len(expected) > 0 {
		t.Errorf("missing %d results", len(expected))
	}
}

// ============================================================================
// Benchmark Tests (to measure performance)
// ============================================================================

func BenchmarkSafeCounterMutex(b *testing.B) {
	counter := NewSafeCounter()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Increment()
		}
	})
}

func BenchmarkSafeMapWrite(b *testing.B) {
	m := NewSafeMap()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Set(fmt.Sprintf("key-%d", i), i)
			i++
		}
	})
}

func BenchmarkSafeMapRead(b *testing.B) {
	m := NewSafeMap()

	// Prepopulate
	for i := 0; i < 1000; i++ {
		m.Set(fmt.Sprintf("key-%d", i), i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			m.Get(fmt.Sprintf("key-%d", i%1000))
			i++
		}
	})
}
