//go:build !solution && !reference

package racedetectiondemo

import (
	"sync"
)

// ============================================================================
// Exercise 1: Safe Counter
// ============================================================================

func NewSafeCounter() *SafeCounter {
	return &SafeCounter{}
}

func (c *SafeCounter) Increment() {
	c.value.Add(1)
}

func (c *SafeCounter) Value() int64 {
	return c.value.Load()
}

// ============================================================================
// Exercise 2: Safe Map
// ============================================================================

func NewSafeMap() *SafeMap {
	return &SafeMap{data: make(map[string]int)}
}

func (m *SafeMap) Set(key string, value int) {
	m.mu.Lock()
	m.data[key] = value
	m.mu.Unlock()
}

func (m *SafeMap) Get(key string) (int, bool) {
	m.mu.RLock()
	v, ok := m.data[key]
	m.mu.RUnlock()
	return v, ok
}

func (m *SafeMap) Len() int {
	m.mu.RLock()
	n := len(m.data)
	m.mu.RUnlock()
	return n
}

// ============================================================================
// Exercise 3: Lazy Initialization
// ============================================================================

func NewLazyInit() *LazyInit {
	return &LazyInit{}
}

func (l *LazyInit) GetOrInit(init func() any) any {
	l.once.Do(func() {
		l.value = init()
	})
	return l.value
}

// ============================================================================
// Exercise 4: Safe Slice
// ============================================================================

func NewSafeSlice() *SafeSlice {
	return &SafeSlice{data: make([]int, 0)}
}

func (s *SafeSlice) Append(v int) {
	s.mu.Lock()
	s.data = append(s.data, v)
	s.mu.Unlock()
}

func (s *SafeSlice) Get(index int) (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if index < 0 || index >= len(s.data) {
		return 0, false
	}
	return s.data[index], true
}

func (s *SafeSlice) Len() int {
	s.mu.RLock()
	n := len(s.data)
	s.mu.RUnlock()
	return n
}

// ============================================================================
// Exercise 5: Loop Variable Capture
// ============================================================================

func ProcessIDs(ids []int, process func(int) int) []int {
	var wg sync.WaitGroup
	out := make([]int, len(ids))

	for i, id := range ids {
		wg.Add(1)
		go func(i int, id int) {
			defer wg.Done()
			out[i] = process(id)
		}(i, id)
	}

	wg.Wait()
	return out
}

// ============================================================================
// Exercise 6: Concurrent URL Cache
// ============================================================================

func NewURLCache(fetcher func(url string) (string, error)) *URLCache {
	return &URLCache{
		cache:   make(map[string]string),
		fetcher: fetcher,
	}
}

func (c *URLCache) Fetch(url string) (string, error) {
	// Fast path: read lock for cache hit.
	c.mu.RLock()
	if v, ok := c.cache[url]; ok {
		c.mu.RUnlock()
		return v, nil
	}
	c.mu.RUnlock()

	// Slow path: single-writer update.
	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after waiting for write lock.
	if v, ok := c.cache[url]; ok {
		return v, nil
	}

	v, err := c.fetcher(url)
	if err != nil {
		return "", err
	}
	c.cache[url] = v
	return v, nil
}

// ============================================================================
// Exercise 7: Concurrent Metrics
// ============================================================================

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (m *Metrics) IncrementRequests() {
	m.requests.Add(1)
}

func (m *Metrics) IncrementErrors() {
	m.errors.Add(1)
}

func (m *Metrics) GetStats() (requests int64, errors int64) {
	return m.requests.Load(), m.errors.Load()
}

// ============================================================================
// Exercise 8: Bank Account
// ============================================================================

func NewBankAccount(initialBalance int64) *BankAccount {
	return &BankAccount{balance: initialBalance}
}

func (b *BankAccount) Deposit(amount int64) {
	b.mu.Lock()
	b.balance += amount
	b.mu.Unlock()
}

func (b *BankAccount) Withdraw(amount int64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.balance < amount {
		return false
	}
	b.balance -= amount
	return true
}

func (b *BankAccount) Balance() int64 {
	b.mu.Lock()
	v := b.balance
	b.mu.Unlock()
	return v
}

// ============================================================================
// Exercise 9: Pipeline Pattern
// ============================================================================

func Pipeline(numbers []int) int {
	gen := make(chan int)
	go func() {
		defer close(gen)
		for _, n := range numbers {
			gen <- n
		}
	}()

	sq := make(chan int)
	go func() {
		defer close(sq)
		for n := range gen {
			sq <- n * n
		}
	}()

	sum := 0
	for n := range sq {
		sum += n
	}
	return sum
}

// ============================================================================
// Exercise 10: Worker Pool
// ============================================================================

func WorkerPool(numWorkers int, jobs []int, process func(int) int) []int {
	if numWorkers <= 0 {
		numWorkers = 1
	}

	jobsCh := make(chan int)
	resultsCh := make(chan int)

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobsCh {
				resultsCh <- process(j)
			}
		}()
	}

	go func() {
		for _, j := range jobs {
			jobsCh <- j
		}
		close(jobsCh)
	}()

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	out := make([]int, 0, len(jobs))
	for r := range resultsCh {
		out = append(out, r)
	}
	return out
}
