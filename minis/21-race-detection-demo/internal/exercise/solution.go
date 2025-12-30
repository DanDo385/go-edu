package exercise

// This file contains solutions to all the race detection exercises.
// Students should implement these in exercise.go before looking at this file.

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// ============================================================================
// Solution 1: Safe Counter
// ============================================================================

// Solution: Use atomic.Int64 for lock-free counter
type SafeCounterSolution struct {
	value atomic.Int64
}

func NewSafeCounterSolution() *SafeCounterSolution {
	return &SafeCounterSolution{}
}

func (c *SafeCounterSolution) Increment() {
	c.value.Add(1)
}

func (c *SafeCounterSolution) Value() int64 {
	return c.value.Load()
}

// Alternative solution: Use mutex
type SafeCounterMutexSolution struct {
	value int64
	mu    sync.Mutex
}

func NewSafeCounterMutexSolution() *SafeCounterMutexSolution {
	return &SafeCounterMutexSolution{}
}

func (c *SafeCounterMutexSolution) Increment() {
	c.mu.Lock()
	c.value++
	c.mu.Unlock()
}

func (c *SafeCounterMutexSolution) Value() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// ============================================================================
// Solution 2: Safe Map
// ============================================================================

// Solution: Use sync.RWMutex to protect map access
type SafeMapSolution struct {
	data map[string]int
	mu   sync.RWMutex
}

func NewSafeMapSolution() *SafeMapSolution {
	return &SafeMapSolution{
		data: make(map[string]int),
	}
}

func (m *SafeMapSolution) Set(key string, value int) {
	m.mu.Lock()
	m.data[key] = value
	m.mu.Unlock()
}

func (m *SafeMapSolution) Get(key string) (int, bool) {
	m.mu.RLock()
	value, ok := m.data[key]
	m.mu.RUnlock()
	return value, ok
}

func (m *SafeMapSolution) Len() int {
	m.mu.RLock()
	length := len(m.data)
	m.mu.RUnlock()
	return length
}

// Alternative solution: Use sync.Map (built-in concurrent map)
type SafeMapSyncMapSolution struct {
	data sync.Map
}

func NewSafeMapSyncMapSolution() *SafeMapSyncMapSolution {
	return &SafeMapSyncMapSolution{}
}

func (m *SafeMapSyncMapSolution) Set(key string, value int) {
	m.data.Store(key, value)
}

func (m *SafeMapSyncMapSolution) Get(key string) (int, bool) {
	val, ok := m.data.Load(key)
	if !ok {
		return 0, false
	}
	return val.(int), true
}

func (m *SafeMapSyncMapSolution) Len() int {
	count := 0
	m.data.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// ============================================================================
// Solution 3: Lazy Initialization
// ============================================================================

// Solution: Use sync.Once for thread-safe lazy initialization
type LazyInitSolution struct {
	once  sync.Once
	value interface{}
}

func NewLazyInitSolution() *LazyInitSolution {
	return &LazyInitSolution{}
}

func (l *LazyInitSolution) GetOrInit(init func() interface{}) interface{} {
	l.once.Do(func() {
		l.value = init()
	})
	return l.value
}

// ============================================================================
// Solution 4: Safe Slice
// ============================================================================

// Solution: Use mutex to protect slice operations
type SafeSliceSolution struct {
	data []int
	mu   sync.RWMutex
}

func NewSafeSliceSolution() *SafeSliceSolution {
	return &SafeSliceSolution{
		data: make([]int, 0),
	}
}

func (s *SafeSliceSolution) Append(value int) {
	s.mu.Lock()
	s.data = append(s.data, value)
	s.mu.Unlock()
}

func (s *SafeSliceSolution) Get(index int) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if index < 0 || index >= len(s.data) {
		return 0, false
	}
	return s.data[index], true
}

func (s *SafeSliceSolution) Len() int {
	s.mu.RLock()
	length := len(s.data)
	s.mu.RUnlock()
	return length
}

// ============================================================================
// Solution 5: Process IDs (Loop Variable Capture)
// ============================================================================

// Solution: Pass loop variables as arguments to goroutine
func ProcessIDsSolution(ids []int, process func(int) int) []int {
	var wg sync.WaitGroup
	results := make([]int, len(ids))

	for i, id := range ids {
		wg.Add(1)
		// Pass i and id as arguments to avoid capturing loop variables
		go func(index int, value int) {
			defer wg.Done()
			results[index] = process(value)
		}(i, id) // Pass as arguments
	}

	wg.Wait()
	return results
}

// Alternative solution: Shadow the loop variables (Go 1.22+ does this automatically)
func ProcessIDsSolutionShadow(ids []int, process func(int) int) []int {
	var wg sync.WaitGroup
	results := make([]int, len(ids))

	for i, id := range ids {
		// Shadow the loop variables
		i := i
		id := id

		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i] = process(id)
		}()
	}

	wg.Wait()
	return results
}

// ============================================================================
// Solution 6: Concurrent URL Cache
// ============================================================================

// Solution: Use RWMutex for cache access
func (c *URLCache) FetchSolution(url string) (string, error) {
	// First, try to read from cache (read lock)
	c.mu.RLock()
	if content, ok := c.cache[url]; ok {
		c.mu.RUnlock()
		return content, nil
	}
	c.mu.RUnlock()

	// Not in cache, fetch it (write lock)
	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check: another goroutine might have fetched it while we waited for lock
	if content, ok := c.cache[url]; ok {
		return content, nil
	}

	// Fetch and cache
	content, err := c.fetcher(url)
	if err != nil {
		return "", err
	}

	c.cache[url] = content
	return content, nil
}

// Advanced solution: Prevent duplicate fetches using a "flight group" pattern
type URLCacheAdvanced struct {
	cache   map[string]string
	mu      sync.RWMutex
	fetcher func(url string) (string, error)
	// Track in-flight requests
	inflight map[string]*sync.WaitGroup
	inflmu   sync.Mutex
}

func NewURLCacheAdvanced(fetcher func(url string) (string, error)) *URLCacheAdvanced {
	return &URLCacheAdvanced{
		cache:    make(map[string]string),
		fetcher:  fetcher,
		inflight: make(map[string]*sync.WaitGroup),
	}
}

func (c *URLCacheAdvanced) Fetch(url string) (string, error) {
	// Check cache
	c.mu.RLock()
	if content, ok := c.cache[url]; ok {
		c.mu.RUnlock()
		return content, nil
	}
	c.mu.RUnlock()

	// Check if another goroutine is already fetching this URL
	c.inflmu.Lock()
	if wg, ok := c.inflight[url]; ok {
		// Wait for the other goroutine to finish
		c.inflmu.Unlock()
		wg.Wait()

		// Now it should be in cache
		c.mu.RLock()
		content := c.cache[url]
		c.mu.RUnlock()
		return content, nil
	}

	// We're the first, create a WaitGroup for others to wait on
	wg := &sync.WaitGroup{}
	wg.Add(1)
	c.inflight[url] = wg
	c.inflmu.Unlock()

	// Fetch (outside of any locks)
	content, err := c.fetcher(url)

	// Store in cache
	c.mu.Lock()
	if err == nil {
		c.cache[url] = content
	}
	c.mu.Unlock()

	// Remove from inflight and signal waiters
	c.inflmu.Lock()
	delete(c.inflight, url)
	c.inflmu.Unlock()
	wg.Done()

	return content, err
}

// ============================================================================
// Solution 7: Concurrent Metrics
// ============================================================================

// Solution: Use atomic counters
type MetricsSolution struct {
	requests atomic.Int64
	errors   atomic.Int64
}

func NewMetricsSolution() *MetricsSolution {
	return &MetricsSolution{}
}

func (m *MetricsSolution) IncrementRequests() {
	m.requests.Add(1)
}

func (m *MetricsSolution) IncrementErrors() {
	m.errors.Add(1)
}

func (m *MetricsSolution) GetStats() (requests int64, errors int64) {
	return m.requests.Load(), m.errors.Load()
}

// ============================================================================
// Solution 8: Bank Account
// ============================================================================

// Solution: Use mutex to protect balance
type BankAccountSolution struct {
	balance int64
	mu      sync.Mutex
}

func NewBankAccountSolution(initialBalance int64) *BankAccountSolution {
	return &BankAccountSolution{
		balance: initialBalance,
	}
}

func (b *BankAccountSolution) Deposit(amount int64) {
	b.mu.Lock()
	b.balance += amount
	b.mu.Unlock()
}

func (b *BankAccountSolution) Withdraw(amount int64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.balance >= amount {
		b.balance -= amount
		return true
	}
	return false
}

func (b *BankAccountSolution) Balance() int64 {
	b.mu.Lock()
	balance := b.balance
	b.mu.Unlock()
	return balance
}

// ============================================================================
// Solution 9: Pipeline Pattern
// ============================================================================

// Solution: Use channels to connect pipeline stages
func PipelineSolution(numbers []int) int {
	// Stage 1: Generate numbers
	gen := make(chan int)
	go func() {
		defer close(gen)
		for _, n := range numbers {
			gen <- n
		}
	}()

	// Stage 2: Square numbers
	sq := make(chan int)
	go func() {
		defer close(sq)
		for n := range gen {
			sq <- n * n
		}
	}()

	// Stage 3: Sum all squared numbers
	sum := 0
	for n := range sq {
		sum += n
	}

	return sum
}

// Alternative: More functional/composable pipeline
func PipelineSolutionComposable(numbers []int) int {
	// Stage 1: Generate
	generate := func(nums []int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for _, n := range nums {
				out <- n
			}
		}()
		return out
	}

	// Stage 2: Square
	square := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				out <- n * n
			}
		}()
		return out
	}

	// Stage 3: Sum
	sumAll := func(in <-chan int) int {
		sum := 0
		for n := range in {
			sum += n
		}
		return sum
	}

	// Compose pipeline
	gen := generate(numbers)
	sq := square(gen)
	return sumAll(sq)
}

// ============================================================================
// Solution 10: Worker Pool
// ============================================================================

// Solution: Use channels and WaitGroup for worker pool
func WorkerPoolSolution(numWorkers int, jobs []int, process func(int) int) []int {
	// Create channels
	jobsCh := make(chan int, len(jobs))
	resultsCh := make(chan int, len(jobs))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Process jobs from channel until it's closed
			for job := range jobsCh {
				result := process(job)
				resultsCh <- result
			}
		}()
	}

	// Send all jobs
	for _, job := range jobs {
		jobsCh <- job
	}
	close(jobsCh) // Signal no more jobs

	// Close results channel when all workers finish
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// Collect results
	results := make([]int, 0, len(jobs))
	for result := range resultsCh {
		results = append(results, result)
	}

	return results
}

// Alternative: Pre-allocated results slice (if order matters)
func WorkerPoolSolutionOrdered(numWorkers int, jobs []int, process func(int) int) []int {
	type result struct {
		index int
		value int
	}

	jobsCh := make(chan struct {
		index int
		value int
	}, len(jobs))
	resultsCh := make(chan result, len(jobs))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobsCh {
				resultsCh <- result{
					index: job.index,
					value: process(job.value),
				}
			}
		}()
	}

	// Send jobs
	for i, job := range jobs {
		jobsCh <- struct {
			index int
			value int
		}{index: i, value: job}
	}
	close(jobsCh)

	// Close results when done
	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	// Collect results in order
	results := make([]int, len(jobs))
	for res := range resultsCh {
		results[res.index] = res.value
	}

	return results
}

// ============================================================================
// Explanation: Why These Solutions Work
// ============================================================================

/*
KEY PRINCIPLES FOR RACE-FREE CODE:

1. **Mutual Exclusion (Locks)**
   - Use sync.Mutex for exclusive access to shared state
   - Use sync.RWMutex when you have many readers, few writers
   - Always defer mu.Unlock() to ensure unlock happens even on panic

2. **Atomic Operations**
   - Use atomic.Int64, atomic.Bool, etc. for simple counters/flags
   - Lock-free, very fast, but limited to simple types
   - Perfect for counters, flags, and simple state

3. **Single Ownership (Channels)**
   - One goroutine owns the data, others send requests via channels
   - No locks needed because there's no sharing
   - Most idiomatic Go approach

4. **Immutability**
   - Data that never changes can be safely shared
   - Use atomic.Value to swap immutable configs

5. **Confinement**
   - Each goroutine has its own data, no sharing
   - Use channels to transfer ownership when needed

6. **sync.Once**
   - Guarantees a function runs exactly once, even with concurrent calls
   - Perfect for lazy initialization

COMMON RACE PATTERNS TO AVOID:

1. **Unsynchronized counter++**
   → Fix: Use atomic.Add() or mutex

2. **Concurrent map access**
   → Fix: Use sync.RWMutex or sync.Map

3. **Loop variable capture in goroutines**
   → Fix: Pass as argument or shadow the variable

4. **Double-checked locking**
   → Fix: Use sync.Once

5. **Concurrent append to slice**
   → Fix: Use mutex or use channel to collect results

6. **Reading while writing struct fields**
   → Fix: Protect entire struct with mutex or use atomic fields

TESTING FOR RACES:

1. Always run tests with: go test -race
2. The race detector only finds races that actually execute
3. Achieve high code coverage to maximize race detection
4. Use stress tests with many goroutines
5. Test with different GOMAXPROCS values

PERFORMANCE CONSIDERATIONS:

1. Atomic operations > Mutex > Channels (for simple counters)
2. RWMutex > Mutex (for read-heavy workloads)
3. Channels are best for complex coordination, not simple sharing
4. Profile before optimizing (premature optimization is evil)

Remember: "Don't communicate by sharing memory; share memory by communicating."
*/

// ============================================================================
// Example: Complete Race-Free Server Metrics
// ============================================================================

// ServerMetrics demonstrates a complete race-free metrics system
type ServerMetrics struct {
	requests        atomic.Int64
	errors          atomic.Int64
	activeConns     atomic.Int64
	responseTimes   []int64
	responseTimesMu sync.Mutex
}

func NewServerMetrics() *ServerMetrics {
	return &ServerMetrics{
		responseTimes: make([]int64, 0),
	}
}

func (m *ServerMetrics) RecordRequest() {
	m.requests.Add(1)
}

func (m *ServerMetrics) RecordError() {
	m.errors.Add(1)
}

func (m *ServerMetrics) RecordResponseTime(ms int64) {
	m.responseTimesMu.Lock()
	m.responseTimes = append(m.responseTimes, ms)
	m.responseTimesMu.Unlock()
}

func (m *ServerMetrics) ConnOpened() {
	m.activeConns.Add(1)
}

func (m *ServerMetrics) ConnClosed() {
	m.activeConns.Add(-1)
}

func (m *ServerMetrics) Snapshot() string {
	m.responseTimesMu.Lock()
	avgResponseTime := int64(0)
	if len(m.responseTimes) > 0 {
		sum := int64(0)
		for _, t := range m.responseTimes {
			sum += t
		}
		avgResponseTime = sum / int64(len(m.responseTimes))
	}
	m.responseTimesMu.Unlock()

	return fmt.Sprintf(
		"Requests: %d, Errors: %d, Active: %d, Avg Response: %dms",
		m.requests.Load(),
		m.errors.Load(),
		m.activeConns.Load(),
		avgResponseTime,
	)
}
