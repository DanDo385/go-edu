package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	fmt.Println("=== Race Detection Demonstration ===")
	fmt.Println()

	// Demonstration 1: Counter Increment Race
	demo1CounterRace()

	// Demonstration 2: Map Concurrent Access
	demo2MapRace()

	// Demonstration 3: Loop Variable Capture
	demo3LoopVariableRace()

	// Demonstration 4: Lazy Initialization Race
	demo4LazyInitRace()

	// Demonstration 5: Slice Append Race
	demo5SliceAppendRace()

	// Demonstration 6: Struct Field Race
	demo6StructFieldRace()

	fmt.Println()
	fmt.Println("=== All demonstrations complete! ===")
	fmt.Println("Run with: go run -race cmd/race-demo/main.go")
	fmt.Println("To see race detector warnings on the buggy versions.")
}

// ============================================================================
// Demo 1: Counter Increment Race
// ============================================================================

func demo1CounterRace() {
	fmt.Println("[Demo 1] Counter Increment Race")
	fmt.Println("Problem: counter++ is not atomic (read-modify-write)")
	fmt.Println()

	// Buggy version
	fmt.Println("  Running BUGGY version (uncomment to see race)...")
	// Uncomment the next line to see the race:
	// runBuggyCounter()

	// Fixed version with mutex
	fmt.Println("  Running FIXED version (mutex)...")
	runFixedCounterMutex()

	// Fixed version with atomic
	fmt.Println("  Running FIXED version (atomic)...")
	runFixedCounterAtomic()

	fmt.Println()
}

// Buggy: Unsynchronized access to shared counter
func runBuggyCounter() {
	var counter int
	var wg sync.WaitGroup

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				counter++ // RACE: Read-modify-write without synchronization
			}
		}()
	}

	wg.Wait()
	fmt.Printf("    Final counter: %d (expected 2000, but varies due to race)\n", counter)
}

// Fixed: Using mutex
func runFixedCounterMutex() {
	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("    Final counter: %d (correct!)\n", counter)
}

// Fixed: Using atomic operations
func runFixedCounterAtomic() {
	var counter atomic.Int64
	var wg sync.WaitGroup

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				counter.Add(1)
			}
		}()
	}

	wg.Wait()
	fmt.Printf("    Final counter: %d (correct!)\n", counter.Load())
}

// ============================================================================
// Demo 2: Map Concurrent Access
// ============================================================================

func demo2MapRace() {
	fmt.Println("[Demo 2] Map Concurrent Read/Write")
	fmt.Println("Problem: Maps are not thread-safe, concurrent access causes panic")
	fmt.Println()

	// Buggy version
	fmt.Println("  Running BUGGY version (uncomment to see panic/race)...")
	// Uncomment the next line to see the panic:
	// runBuggyMap()

	// Fixed version
	fmt.Println("  Running FIXED version (mutex)...")
	runFixedMapMutex()

	fmt.Println("  Running FIXED version (sync.Map)...")
	runFixedSyncMap()

	fmt.Println()
}

// Buggy: Concurrent map access
func runBuggyMap() {
	cache := make(map[int]int)
	var wg sync.WaitGroup

	// Writers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				cache[id] = j // RACE: Concurrent writes
			}
		}(i)
	}

	// Readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = cache[id] // RACE: Concurrent reads during writes
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("    Map size: %d (if it didn't panic)\n", len(cache))
}

// Fixed: Using mutex
func runFixedMapMutex() {
	cache := make(map[int]int)
	var mu sync.RWMutex
	var wg sync.WaitGroup

	// Writers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				mu.Lock()
				cache[id] = j
				mu.Unlock()
			}
		}(i)
	}

	// Readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				mu.RLock()
				_ = cache[id]
				mu.RUnlock()
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("    Map size: %d (correct!)\n", len(cache))
}

// Fixed: Using sync.Map
func runFixedSyncMap() {
	var cache sync.Map
	var wg sync.WaitGroup

	// Writers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				cache.Store(id, j)
			}
		}(i)
	}

	// Readers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				cache.Load(id)
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("    sync.Map operations completed (correct!)\n")
}

// ============================================================================
// Demo 3: Loop Variable Capture
// ============================================================================

func demo3LoopVariableRace() {
	fmt.Println("[Demo 3] Loop Variable Capture")
	fmt.Println("Problem: Loop variable is shared across all goroutines")
	fmt.Println()

	// Buggy version
	fmt.Println("  Running BUGGY version...")
	runBuggyLoopCapture()

	// Fixed version
	fmt.Println("  Running FIXED version (pass as argument)...")
	runFixedLoopCapture()

	fmt.Println()
}

// Buggy: Captures loop variable by reference
func runBuggyLoopCapture() {
	var wg sync.WaitGroup

	fmt.Print("    Output: ")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// RACE: All goroutines share the same &i variable
			// By the time they run, i is likely 5
			fmt.Printf("%d ", i)
		}()
	}

	wg.Wait()
	fmt.Println("(often prints '5 5 5 5 5' - WRONG!)")
}

// Fixed: Pass loop variable as argument
func runFixedLoopCapture() {
	var wg sync.WaitGroup

	fmt.Print("    Output: ")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Each goroutine gets its own copy of i
			fmt.Printf("%d ", id)
		}(i) // Pass i as argument
	}

	wg.Wait()
	fmt.Println("(prints 0-4 in some order - correct!)")
}

// ============================================================================
// Demo 4: Lazy Initialization Race (Double-Checked Locking)
// ============================================================================

type Config struct {
	Host string
	Port int
}

func demo4LazyInitRace() {
	fmt.Println("[Demo 4] Lazy Initialization (Double-Checked Locking)")
	fmt.Println("Problem: Reading without lock creates race condition")
	fmt.Println()

	// Buggy version
	fmt.Println("  Running BUGGY version (uncomment to see race)...")
	// Uncomment the next line to see the race:
	// runBuggyLazyInit()

	// Fixed version with sync.Once
	fmt.Println("  Running FIXED version (sync.Once)...")
	runFixedLazyInitOnce()

	fmt.Println()
}

var (
	buggyConfig *Config
	buggyMu     sync.Mutex
)

// Buggy: Double-checked locking anti-pattern
func getBuggyConfig() *Config {
	if buggyConfig == nil { // RACE: Read without lock
		buggyMu.Lock()
		if buggyConfig == nil {
			// Simulate slow initialization
			time.Sleep(1 * time.Millisecond)
			buggyConfig = &Config{Host: "localhost", Port: 8080}
		}
		buggyMu.Unlock()
	}
	return buggyConfig
}

func runBuggyLazyInit() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cfg := getBuggyConfig()
			fmt.Printf("    Got config: %+v\n", cfg)
		}()
	}

	wg.Wait()
}

// Fixed: Using sync.Once (idiomatic Go)
var (
	fixedConfig *Config
	once        sync.Once
)

func getFixedConfig() *Config {
	once.Do(func() {
		// This runs exactly once, guaranteed
		time.Sleep(1 * time.Millisecond)
		fixedConfig = &Config{Host: "localhost", Port: 8080}
	})
	return fixedConfig
}

func runFixedLazyInitOnce() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cfg := getFixedConfig()
			if id == 0 {
				fmt.Printf("    Got config: %+v (correct!)\n", cfg)
			}
		}(i)
	}

	wg.Wait()
}

// ============================================================================
// Demo 5: Slice Append Race
// ============================================================================

func demo5SliceAppendRace() {
	fmt.Println("[Demo 5] Slice Append Race")
	fmt.Println("Problem: append modifies slice header (len, cap, ptr)")
	fmt.Println()

	// Buggy version
	fmt.Println("  Running BUGGY version (uncomment to see race)...")
	// Uncomment the next line to see the race:
	// runBuggySliceAppend()

	// Fixed version
	fmt.Println("  Running FIXED version (mutex)...")
	runFixedSliceAppendMutex()

	fmt.Println("  Running FIXED version (channel)...")
	runFixedSliceAppendChannel()

	fmt.Println()
}

// Buggy: Concurrent append without synchronization
func runBuggySliceAppend() {
	var results []int
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// RACE: append reads and writes the slice header
			results = append(results, id)
		}(i)
	}

	wg.Wait()
	fmt.Printf("    Results: %v (len=%d, might be < 10 due to race)\n", results, len(results))
}

// Fixed: Using mutex
func runFixedSliceAppendMutex() {
	var results []int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			mu.Lock()
			results = append(results, id)
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	fmt.Printf("    Results len: %d (correct!)\n", len(results))
}

// Fixed: Using channel (idiomatic)
func runFixedSliceAppendChannel() {
	ch := make(chan int, 10)
	var wg sync.WaitGroup

	// Producers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ch <- id
		}(i)
	}

	// Close channel when all producers finish
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Collector (single goroutine owns the slice)
	var results []int
	for val := range ch {
		results = append(results, val)
	}

	fmt.Printf("    Results len: %d (correct!)\n", len(results))
}

// ============================================================================
// Demo 6: Struct Field Race
// ============================================================================

type Stats struct {
	Requests int
	Errors   int
}

func demo6StructFieldRace() {
	fmt.Println("[Demo 6] Struct Field Race")
	fmt.Println("Problem: Even different fields can race (cache line sharing)")
	fmt.Println()

	// Buggy version
	fmt.Println("  Running BUGGY version (uncomment to see race)...")
	// Uncomment the next line to see the race:
	// runBuggyStructFields()

	// Fixed version with mutex
	fmt.Println("  Running FIXED version (mutex)...")
	runFixedStructFieldsMutex()

	// Fixed version with atomic fields
	fmt.Println("  Running FIXED version (atomic fields)...")
	runFixedStructFieldsAtomic()

	fmt.Println()
}

// Buggy: Concurrent writes to struct fields
func runBuggyStructFields() {
	var stats Stats
	var wg sync.WaitGroup

	// Goroutine 1: Increments Requests
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				stats.Requests++ // RACE
			}
		}()
	}

	// Goroutine 2: Increments Errors
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				stats.Errors++ // RACE
			}
		}()
	}

	wg.Wait()
	fmt.Printf("    Stats: Requests=%d, Errors=%d (both wrong due to race)\n",
		stats.Requests, stats.Errors)
}

// Fixed: Using mutex to protect all fields
func runFixedStructFieldsMutex() {
	var stats Stats
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Goroutine 1: Increments Requests
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mu.Lock()
				stats.Requests++
				mu.Unlock()
			}
		}()
	}

	// Goroutine 2: Increments Errors
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mu.Lock()
				stats.Errors++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	fmt.Printf("    Stats: Requests=%d, Errors=%d (correct!)\n",
		stats.Requests, stats.Errors)
}

// Fixed: Using atomic fields
type AtomicStats struct {
	Requests atomic.Int64
	Errors   atomic.Int64
}

func runFixedStructFieldsAtomic() {
	var stats AtomicStats
	var wg sync.WaitGroup

	// Goroutine 1: Increments Requests
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				stats.Requests.Add(1)
			}
		}()
	}

	// Goroutine 2: Increments Errors
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				stats.Errors.Add(1)
			}
		}()
	}

	wg.Wait()
	fmt.Printf("    Stats: Requests=%d, Errors=%d (correct!)\n",
		stats.Requests.Load(), stats.Errors.Load())
}

// ============================================================================
// Additional Example: Race-Free Pipeline Pattern
// ============================================================================

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Demonstrates a race-free pipeline using channels
func exampleRaceFreePipeline() {
	// Stage 1: Generate numbers
	generate := func(nums ...int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for _, n := range nums {
				out <- n
			}
		}()
		return out
	}

	// Stage 2: Square numbers
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

	// Stage 3: Filter even numbers
	filterEven := func(in <-chan int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for n := range in {
				if n%2 == 0 {
					out <- n
				}
			}
		}()
		return out
	}

	// Build and run pipeline
	nums := generate(1, 2, 3, 4, 5)
	squared := square(nums)
	filtered := filterEven(squared)

	// Consume results
	for result := range filtered {
		fmt.Println(result)
	}
}
