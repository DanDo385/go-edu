package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	fmt.Println("=== Mutex vs RWMutex Demonstrations ===")
	fmt.Println()

	demo1RaceCondition()
	demo2MutexBasics()
	demo3RWMutexBasics()
	demo4PerformanceComparison()
	demo5LockContention()
	demo6ShardedMap()
	demo7DeadlockPrevention()
}

// Demo 1: Race Condition (UNSAFE)
func demo1RaceCondition() {
	fmt.Println("--- Demo 1: Race Condition (UNSAFE) ---")

	// This demonstrates WHY we need synchronization
	var counter int
	var wg sync.WaitGroup

	// Start 100 goroutines, each incrementing counter 100 times
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				counter++ // RACE CONDITION!
			}
		}()
	}

	wg.Wait()

	fmt.Printf("Expected: 10000\n")
	fmt.Printf("Actual:   %d\n", counter)
	fmt.Printf("Lost increments: %d\n", 10000-counter)
	fmt.Println("Note: Run with 'go run -race main.go' to detect the race")
	fmt.Println()
}

// Demo 2: Mutex Basics (SAFE)
func demo2MutexBasics() {
	fmt.Println("--- Demo 2: Mutex Basics (SAFE) ---")

	type SafeCounter struct {
		mu    sync.Mutex
		value int
	}

	counter := &SafeCounter{}
	var wg sync.WaitGroup

	// Start 100 goroutines, each incrementing counter 100 times
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				counter.mu.Lock()
				counter.value++
				counter.mu.Unlock()
			}
		}()
	}

	wg.Wait()

	fmt.Printf("Expected: 10000\n")
	fmt.Printf("Actual:   %d\n", counter.value)
	fmt.Printf("Lost increments: %d\n", 10000-counter.value)
	fmt.Println("Success: Mutex prevented race conditions!")
	fmt.Println()
}

// Demo 3: RWMutex Basics
func demo3RWMutexBasics() {
	fmt.Println("--- Demo 3: RWMutex Basics ---")

	type Cache struct {
		mu   sync.RWMutex
		data map[string]int
	}

	cache := &Cache{
		data: make(map[string]int),
	}

	// Pre-populate cache
	for i := 0; i < 10; i++ {
		cache.mu.Lock()
		cache.data[fmt.Sprintf("key%d", i)] = i * 10
		cache.mu.Unlock()
	}

	var wg sync.WaitGroup
	reads := atomic.Int64{}
	writes := atomic.Int64{}

	// 90 reader goroutines
	for i := 0; i < 90; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("key%d", j%10)

				cache.mu.RLock()
				_ = cache.data[key]
				cache.mu.RUnlock()

				reads.Add(1)
			}
		}(i)
	}

	// 10 writer goroutines
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("key%d", j%10)

				cache.mu.Lock()
				cache.data[key] = j
				cache.mu.Unlock()

				writes.Add(1)
			}
		}(i)
	}

	wg.Wait()

	fmt.Printf("Total reads:  %d\n", reads.Load())
	fmt.Printf("Total writes: %d\n", writes.Load())
	fmt.Printf("Read/Write ratio: %.1f%%/%.1f%%\n",
		float64(reads.Load())/float64(reads.Load()+writes.Load())*100,
		float64(writes.Load())/float64(reads.Load()+writes.Load())*100)
	fmt.Println("RWMutex allowed concurrent reads while serializing writes!")
	fmt.Println()
}

// Demo 4: Performance Comparison
func demo4PerformanceComparison() {
	fmt.Println("--- Demo 4: Performance Comparison ---")

	// Mutex-based cache
	type MutexCache struct {
		mu   sync.Mutex
		data map[string]int
	}

	// RWMutex-based cache
	type RWMutexCache struct {
		mu   sync.RWMutex
		data map[string]int
	}

	// Benchmark Mutex
	mutexCache := &MutexCache{data: make(map[string]int)}
	for i := 0; i < 100; i++ {
		mutexCache.data[fmt.Sprintf("key%d", i)] = i
	}

	start := time.Now()
	var wg sync.WaitGroup

	// 100 readers
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10000; j++ {
				key := fmt.Sprintf("key%d", j%100)
				mutexCache.mu.Lock()
				_ = mutexCache.data[key]
				mutexCache.mu.Unlock()
			}
		}()
	}

	wg.Wait()
	mutexDuration := time.Since(start)

	// Benchmark RWMutex
	rwMutexCache := &RWMutexCache{data: make(map[string]int)}
	for i := 0; i < 100; i++ {
		rwMutexCache.data[fmt.Sprintf("key%d", i)] = i
	}

	start = time.Now()

	// 100 readers
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10000; j++ {
				key := fmt.Sprintf("key%d", j%100)
				rwMutexCache.mu.RLock()
				_ = rwMutexCache.data[key]
				rwMutexCache.mu.RUnlock()
			}
		}()
	}

	wg.Wait()
	rwMutexDuration := time.Since(start)

	fmt.Printf("Mutex duration:   %v\n", mutexDuration)
	fmt.Printf("RWMutex duration: %v\n", rwMutexDuration)
	fmt.Printf("Speedup: %.2fx faster\n", float64(mutexDuration)/float64(rwMutexDuration))
	fmt.Println()
}

// Demo 5: Lock Contention Visualization
func demo5LockContention() {
	fmt.Println("--- Demo 5: Lock Contention Visualization ---")

	type ContentionMetrics struct {
		mu           sync.Mutex
		totalWait    time.Duration
		maxWait      time.Duration
		lockAttempts int64
	}

	metrics := &ContentionMetrics{}
	var wg sync.WaitGroup

	// Simulate high contention
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < 100; j++ {
				start := time.Now()

				metrics.mu.Lock()
				waitTime := time.Since(start)

				// Track metrics
				metrics.totalWait += waitTime
				if waitTime > metrics.maxWait {
					metrics.maxWait = waitTime
				}
				metrics.lockAttempts++

				// Simulate work
				time.Sleep(time.Microsecond * 100)

				metrics.mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	avgWait := metrics.totalWait / time.Duration(metrics.lockAttempts)
	fmt.Printf("Total lock attempts: %d\n", metrics.lockAttempts)
	fmt.Printf("Average wait time: %v\n", avgWait)
	fmt.Printf("Max wait time: %v\n", metrics.maxWait)
	fmt.Printf("Total time spent waiting: %v\n", metrics.totalWait)
	fmt.Println()
}

// Demo 6: Sharded Map for Reduced Contention
func demo6ShardedMap() {
	fmt.Println("--- Demo 6: Sharded Map (Reduced Contention) ---")

	type Shard struct {
		mu   sync.RWMutex
		data map[string]int
	}

	type ShardedMap struct {
		shards [16]*Shard
	}

	newShardedMap := func() *ShardedMap {
		sm := &ShardedMap{}
		for i := 0; i < 16; i++ {
			sm.shards[i] = &Shard{
				data: make(map[string]int),
			}
		}
		return sm
	}

	getShard := func(sm *ShardedMap, key string) *Shard {
		hash := 0
		for _, c := range key {
			hash += int(c)
		}
		return sm.shards[hash%16]
	}

	shardedMap := newShardedMap()

	// Pre-populate
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key%d", i)
		shard := getShard(shardedMap, key)
		shard.mu.Lock()
		shard.data[key] = i
		shard.mu.Unlock()
	}

	start := time.Now()
	var wg sync.WaitGroup

	// 100 goroutines doing random reads/writes
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r := rand.New(rand.NewSource(time.Now().UnixNano()))

			for j := 0; j < 10000; j++ {
				key := fmt.Sprintf("key%d", r.Intn(1000))
				shard := getShard(shardedMap, key)

				if r.Float32() < 0.9 {
					// 90% reads
					shard.mu.RLock()
					_ = shard.data[key]
					shard.mu.RUnlock()
				} else {
					// 10% writes
					shard.mu.Lock()
					shard.data[key] = r.Intn(1000)
					shard.mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	fmt.Printf("Operations: 1,000,000\n")
	fmt.Printf("Duration: %v\n", duration)
	fmt.Printf("Throughput: %.0f ops/sec\n", 1000000/duration.Seconds())
	fmt.Println("Sharding reduced contention by distributing load across 16 independent locks!")
	fmt.Println()
}

// Demo 7: Deadlock Prevention
func demo7DeadlockPrevention() {
	fmt.Println("--- Demo 7: Deadlock Prevention ---")

	type Account struct {
		mu      sync.Mutex
		balance int
		id      int
	}

	// Safe transfer with lock ordering
	transfer := func(from, to *Account, amount int) {
		// Always lock accounts in order of ID to prevent deadlock
		first, second := from, to
		if from.id > to.id {
			first, second = to, from
		}

		first.mu.Lock()
		defer first.mu.Unlock()

		second.mu.Lock()
		defer second.mu.Unlock()

		// Transfer
		from.balance -= amount
		to.balance += amount
	}

	// Create accounts
	account1 := &Account{id: 1, balance: 1000}
	account2 := &Account{id: 2, balance: 1000}
	account3 := &Account{id: 3, balance: 1000}

	var wg sync.WaitGroup

	// Goroutine 1: Transfer 1→2 and 2→3
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			transfer(account1, account2, 10)
			transfer(account2, account3, 5)
		}
	}()

	// Goroutine 2: Transfer 3→1 and 1→2
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			transfer(account3, account1, 8)
			transfer(account1, account2, 3)
		}
	}()

	// Goroutine 3: Transfer 2→1 and 3→2
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			transfer(account2, account1, 7)
			transfer(account3, account2, 4)
		}
	}()

	wg.Wait()

	totalBalance := account1.balance + account2.balance + account3.balance

	fmt.Printf("Account 1: $%d\n", account1.balance)
	fmt.Printf("Account 2: $%d\n", account2.balance)
	fmt.Printf("Account 3: $%d\n", account3.balance)
	fmt.Printf("Total: $%d (should be $3000)\n", totalBalance)

	if totalBalance == 3000 {
		fmt.Println("Success: No deadlock, money conserved!")
	} else {
		fmt.Println("Error: Money was lost or created!")
	}
	fmt.Println()
}
