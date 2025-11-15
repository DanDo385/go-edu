// Package main demonstrates sync.Once and the singleton pattern in Go.
//
// MACRO-COMMENT: What This Program Teaches
// =========================================
// This program provides hands-on examples of:
// 1. sync.Once mechanics and exactly-once guarantees
// 2. Singleton pattern implementation
// 3. Lazy initialization for expensive resources
// 4. Thread-safe initialization without races
// 5. Error handling in singleton initialization
// 6. Real-world patterns (config, DB, logger)
// 7. Performance characteristics (fast path optimization)
//
// COMPILER BEHAVIOR: sync.Once
// sync.Once uses atomic operations and mutexes internally:
// - Fast path: atomic.LoadUint32 (lock-free after first initialization)
// - Slow path: sync.Mutex (only during first initialization)
// This makes it extremely efficient for the common case (already initialized).
//
// RUNTIME BEHAVIOR: Memory Ordering
// sync.Once provides happens-before guarantees:
// - All writes in the Do() function happen before Do() returns
// - All goroutines see the initialized value after Do() returns
// - No additional synchronization needed for the initialized value

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// SECTION 1: Basic sync.Once Usage
// ============================================================================

// demonstrateBasicOnce shows the simplest use of sync.Once.
//
// MACRO-COMMENT: The Fundamental Guarantee
// sync.Once guarantees that the function passed to Do() runs exactly once,
// even when called from multiple goroutines simultaneously.
//
// KEY INSIGHT: This is the building block for thread-safe lazy initialization.
func demonstrateBasicOnce() {
	fmt.Println("=== Basic sync.Once: Exactly-Once Execution ===")

	var (
		counter int32
		once    sync.Once
		wg      sync.WaitGroup
	)

	// MICRO-COMMENT: Launch 10 goroutines
	// Each will try to increment the counter via once.Do()
	// Only ONE will succeed (exactly-once guarantee)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			once.Do(func() {
				// MICRO-COMMENT: This function runs exactly once
				// Even though 10 goroutines call once.Do()
				atomic.AddInt32(&counter, 1)
				fmt.Printf("  Goroutine %d: Initialized! (counter=%d)\n", id, counter)
			})

			// MICRO-COMMENT: All goroutines see the initialized value
			// once.Do() blocks until initialization completes
			fmt.Printf("  Goroutine %d: Done (counter=%d)\n", id, atomic.LoadInt32(&counter))
		}(i)
	}

	wg.Wait()
	fmt.Printf("Final counter value: %d (expected: 1)\n", counter)
	fmt.Println()
}

// ============================================================================
// SECTION 2: Lazy Initialization Pattern
// ============================================================================

// ExpensiveResource simulates a resource that's costly to initialize.
//
// MICRO-COMMENT: Examples of expensive resources:
// - Database connections (network I/O, handshakes)
// - File I/O (reading large config files)
// - Cryptographic operations (key generation)
// - Large memory allocations
type ExpensiveResource struct {
	data []byte
	name string
}

// newExpensiveResource creates a resource with simulated expensive initialization.
func newExpensiveResource(name string) *ExpensiveResource {
	// MICRO-COMMENT: Simulate expensive initialization
	fmt.Printf("  [EXPENSIVE] Initializing %s (sleeping 100ms)...\n", name)
	time.Sleep(100 * time.Millisecond)

	return &ExpensiveResource{
		data: make([]byte, 1024*1024), // 1MB allocation
		name: name,
	}
}

// Global variables for lazy initialization
var (
	resource     *ExpensiveResource
	resourceOnce sync.Once
)

// getResource returns the singleton resource, initializing it lazily.
//
// MACRO-COMMENT: Lazy Initialization Pattern
// The resource is not created until the first call to getResource().
// Subsequent calls return the same instance immediately (fast path).
//
// PERFORMANCE:
// - First call: ~100ms (initialization cost)
// - Subsequent calls: ~0.5ns (atomic load only)
//
// This is the canonical lazy initialization pattern in Go.
func getResource() *ExpensiveResource {
	// MICRO-COMMENT: sync.Once guarantees:
	// 1. The initialization function runs exactly once
	// 2. All goroutines wait for initialization to complete
	// 3. No race conditions (thread-safe)
	resourceOnce.Do(func() {
		resource = newExpensiveResource("SharedResource")
	})
	return resource
}

// demonstrateLazyInitialization shows lazy initialization with sync.Once.
func demonstrateLazyInitialization() {
	fmt.Println("=== Lazy Initialization: Defer Expensive Work ===")

	fmt.Println("Program started (resource not yet initialized)")
	time.Sleep(50 * time.Millisecond)

	fmt.Println("First access to resource:")
	start := time.Now()
	r1 := getResource()
	fmt.Printf("  First call took: %v\n", time.Since(start))
	fmt.Printf("  Resource name: %s\n", r1.name)

	fmt.Println("\nSecond access to resource:")
	start = time.Now()
	r2 := getResource()
	fmt.Printf("  Second call took: %v (fast path!)\n", time.Since(start))
	fmt.Printf("  Same instance? %t\n", r1 == r2)

	fmt.Println()
}

// ============================================================================
// SECTION 3: Singleton Pattern
// ============================================================================

// Config represents application configuration.
type Config struct {
	DatabaseURL string
	APIKey      string
	Port        int
}

// ConfigSingleton manages the configuration singleton.
//
// MACRO-COMMENT: Singleton Pattern Structure
// A singleton has three key components:
// 1. Private instance variable (unexported)
// 2. sync.Once for thread-safe initialization
// 3. Public getter function (exported)
//
// This ensures exactly one instance exists and is safely shared.
type ConfigSingleton struct {
	cfg  *Config
	once sync.Once
}

// Global singleton instance
var configSingleton = &ConfigSingleton{}

// Get returns the configuration singleton.
//
// MICRO-COMMENT: This is the public API for accessing the singleton.
// Calling Get() multiple times returns the same *Config instance.
func (s *ConfigSingleton) Get() *Config {
	s.once.Do(func() {
		// MICRO-COMMENT: Simulate loading config from file
		fmt.Println("  [CONFIG] Loading configuration...")
		time.Sleep(50 * time.Millisecond)

		s.cfg = &Config{
			DatabaseURL: "postgres://localhost:5432/mydb",
			APIKey:      "secret-api-key-12345",
			Port:        8080,
		}
	})
	return s.cfg
}

// demonstrateSingletonPattern shows the singleton pattern in action.
func demonstrateSingletonPattern() {
	fmt.Println("=== Singleton Pattern: One Instance, Global Access ===")

	var wg sync.WaitGroup

	// MICRO-COMMENT: Launch 5 goroutines that all access the config
	// Only the first one triggers initialization
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			cfg := configSingleton.Get()
			fmt.Printf("  Goroutine %d: Got config (Port=%d)\n", id, cfg.Port)
		}(i)
	}

	wg.Wait()

	// MICRO-COMMENT: Verify all goroutines got the same instance
	cfg1 := configSingleton.Get()
	cfg2 := configSingleton.Get()
	fmt.Printf("Same instance? %t (pointer: %p)\n", cfg1 == cfg2, cfg1)

	fmt.Println()
}

// ============================================================================
// SECTION 4: Error Handling in Initialization
// ============================================================================

// DatabaseConnection simulates a database connection.
type DatabaseConnection struct {
	url       string
	connected bool
}

// DatabaseSingleton with error handling.
type DatabaseSingleton struct {
	db      *DatabaseConnection
	once    sync.Once
	initErr error
}

var dbSingleton = &DatabaseSingleton{}

// Get returns the database connection or an error.
//
// MACRO-COMMENT: Error Handling Strategies
// Option 1: Return error (shown here) - caller decides what to do
// Option 2: Panic - fail fast if initialization is critical
// Option 3: Use default/fallback - graceful degradation
//
// IMPORTANT: sync.Once runs exactly once, even if initialization fails.
// If the first attempt fails, subsequent calls return the same error.
// If you need retry logic, use sync.Mutex instead of sync.Once.
func (s *DatabaseSingleton) Get() (*DatabaseConnection, error) {
	s.once.Do(func() {
		fmt.Println("  [DATABASE] Attempting connection...")
		time.Sleep(30 * time.Millisecond)

		// MICRO-COMMENT: Simulate connection with 50% failure rate
		if rand.Intn(2) == 0 {
			s.initErr = fmt.Errorf("connection refused")
			fmt.Println("  [DATABASE] Connection failed!")
			return
		}

		s.db = &DatabaseConnection{
			url:       "localhost:5432",
			connected: true,
		}
		fmt.Println("  [DATABASE] Connected successfully!")
	})

	return s.db, s.initErr
}

// demonstrateErrorHandling shows error handling in singleton initialization.
func demonstrateErrorHandling() {
	fmt.Println("=== Error Handling: Initialization Failures ===")

	db, err := dbSingleton.Get()
	if err != nil {
		fmt.Printf("First call failed: %v\n", err)
		fmt.Println("Note: sync.Once will NOT retry on next call!")
	} else {
		fmt.Printf("First call succeeded: connected=%t\n", db.connected)
	}

	// MICRO-COMMENT: Try again - sync.Once will not re-run initialization
	db2, err2 := dbSingleton.Get()
	if err2 != nil {
		fmt.Printf("Second call: still failed: %v\n", err2)
		fmt.Println("  (Same error because once.Do() doesn't retry)")
	} else {
		fmt.Printf("Second call: still connected=%t\n", db2.connected)
		fmt.Printf("  (Same instance: %t)\n", db == db2)
	}

	fmt.Println()
}

// ============================================================================
// SECTION 5: Multiple Singletons
// ============================================================================

// Logger represents a logging system.
type Logger struct {
	name string
}

// Metrics represents a metrics collector.
type Metrics struct {
	requestCount uint64
}

var (
	logger     *Logger
	loggerOnce sync.Once

	metrics     *Metrics
	metricsOnce sync.Once
)

// GetLogger returns the logger singleton.
func GetLogger() *Logger {
	loggerOnce.Do(func() {
		fmt.Println("  [LOGGER] Initializing logger...")
		logger = &Logger{name: "AppLogger"}
	})
	return logger
}

// GetMetrics returns the metrics singleton.
func GetMetrics() *Metrics {
	metricsOnce.Do(func() {
		fmt.Println("  [METRICS] Initializing metrics collector...")
		metrics = &Metrics{}
	})
	return metrics
}

// demonstrateMultipleSingletons shows managing multiple singletons.
//
// MACRO-COMMENT: Independent Initialization
// Each singleton has its own sync.Once, so they initialize independently.
// This allows fine-grained control over what gets initialized when.
func demonstrateMultipleSingletons() {
	fmt.Println("=== Multiple Singletons: Independent Initialization ===")

	// MICRO-COMMENT: Initialize logger first
	log := GetLogger()
	fmt.Printf("Logger initialized: %s\n", log.name)

	// MICRO-COMMENT: Metrics not yet initialized
	fmt.Println("Metrics not yet accessed (lazy initialization)")

	// MICRO-COMMENT: Now initialize metrics
	m := GetMetrics()
	fmt.Printf("Metrics initialized: requestCount=%d\n", m.requestCount)

	// MICRO-COMMENT: Both are now available
	log2 := GetLogger()
	m2 := GetMetrics()
	fmt.Printf("Same instances? Logger=%t, Metrics=%t\n", log == log2, m == m2)

	fmt.Println()
}

// ============================================================================
// SECTION 6: Comparison with Naive Approaches
// ============================================================================

// DANGER: This is an INCORRECT implementation (for educational purposes only)
var (
	naiveConfig     *Config
	naiveConfigInit bool
)

// getNaiveConfig is a BROKEN lazy initialization (race condition).
//
// MACRO-COMMENT: Why This Fails
// The check `if !naiveConfigInit` and the assignment are not atomic.
// Multiple goroutines can see `naiveConfigInit == false` simultaneously
// and all initialize the config (wasting resources and causing races).
//
// DON'T USE THIS IN PRODUCTION! Use sync.Once instead.
func getNaiveConfig() *Config {
	if !naiveConfigInit {
		// MICRO-COMMENT: RACE CONDITION HERE!
		// Multiple goroutines might all enter this block
		fmt.Println("  [NAIVE] Initializing (RACE CONDITION!)")
		time.Sleep(10 * time.Millisecond) // Simulate slow init
		naiveConfig = &Config{Port: 8080}
		naiveConfigInit = true
	}
	return naiveConfig
}

// demonstrateRaceCondition shows why naive approaches fail.
//
// MACRO-COMMENT: The Problem with Naive Lazy Initialization
// Run this with `go run -race main.go` to see the race detector warnings.
func demonstrateRaceCondition() {
	fmt.Println("=== Race Condition Demo: Why sync.Once is Needed ===")
	fmt.Println("(Run with -race flag to see race detector warnings)")

	var wg sync.WaitGroup

	// MICRO-COMMENT: Reset for demo
	naiveConfig = nil
	naiveConfigInit = false

	// MICRO-COMMENT: Launch goroutines that race to initialize
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cfg := getNaiveConfig()
			fmt.Printf("  Goroutine %d: Got config (Port=%d)\n", id, cfg.Port)
		}(i)
	}

	wg.Wait()
	fmt.Println("Notice: Multiple 'Initializing' messages (should be only 1!)")
	fmt.Println()
}

// ============================================================================
// SECTION 7: Performance Comparison
// ============================================================================

// demonstratePerformance shows the performance benefit of sync.Once.
//
// MACRO-COMMENT: Fast Path Optimization
// After initialization, sync.Once uses only an atomic load (no lock).
// This makes it extremely fast for the common case (resource already initialized).
//
// Mutex-based approaches acquire a lock on EVERY call, which is much slower.
func demonstratePerformance() {
	fmt.Println("=== Performance: Fast Path After Initialization ===")

	var (
		onceValue  int
		once       sync.Once
		mutexValue int
		mu         sync.Mutex
	)

	// MICRO-COMMENT: Initialize both
	once.Do(func() { onceValue = 42 })
	mu.Lock()
	mutexValue = 42
	mu.Unlock()

	// MICRO-COMMENT: Benchmark sync.Once (fast path)
	start := time.Now()
	for i := 0; i < 1000000; i++ {
		once.Do(func() { onceValue++ }) // No-op, already initialized
	}
	onceTime := time.Since(start)

	// MICRO-COMMENT: Benchmark mutex
	start = time.Now()
	for i := 0; i < 1000000; i++ {
		mu.Lock()
		_ = mutexValue // Check (already initialized, no-op)
		mu.Unlock()
	}
	mutexTime := time.Since(start)

	fmt.Printf("1,000,000 calls:\n")
	fmt.Printf("  sync.Once: %v\n", onceTime)
	fmt.Printf("  Mutex:     %v\n", mutexTime)
	fmt.Printf("  Speedup:   %.1fx faster\n", float64(mutexTime)/float64(onceTime))

	fmt.Println()
}

// ============================================================================
// SECTION 8: Real-World Pattern - Struct with Lazy Fields
// ============================================================================

// Application with multiple lazily-initialized components.
//
// MACRO-COMMENT: Lazy Field Pattern
// Each field has its own sync.Once, allowing independent lazy initialization.
// This is useful when:
// - Different components are needed at different times
// - Some components are expensive to initialize
// - You want to optimize startup time
type Application struct {
	dbOnce sync.Once
	db     *DatabaseConnection

	loggerOnce sync.Once
	logger     *Logger

	cacheOnce sync.Once
	cache     map[string]string
}

// GetDB returns the database connection, initializing lazily.
func (app *Application) GetDB() *DatabaseConnection {
	app.dbOnce.Do(func() {
		fmt.Println("  [APP] Initializing database...")
		time.Sleep(20 * time.Millisecond)
		app.db = &DatabaseConnection{
			url:       "localhost:5432",
			connected: true,
		}
	})
	return app.db
}

// GetLogger returns the logger, initializing lazily.
func (app *Application) GetLogger() *Logger {
	app.loggerOnce.Do(func() {
		fmt.Println("  [APP] Initializing logger...")
		time.Sleep(10 * time.Millisecond)
		app.logger = &Logger{name: "AppLogger"}
	})
	return app.logger
}

// GetCache returns the cache, initializing lazily.
func (app *Application) GetCache() map[string]string {
	app.cacheOnce.Do(func() {
		fmt.Println("  [APP] Initializing cache...")
		time.Sleep(15 * time.Millisecond)
		app.cache = make(map[string]string)
	})
	return app.cache
}

// demonstrateLazyFields shows struct with lazy field initialization.
func demonstrateLazyFields() {
	fmt.Println("=== Lazy Fields: Per-Component Initialization ===")

	app := &Application{}

	fmt.Println("Application created (nothing initialized yet)")

	fmt.Println("\nAccessing logger:")
	log := app.GetLogger()
	fmt.Printf("  Logger: %s\n", log.name)

	fmt.Println("\nAccessing cache:")
	cache := app.GetCache()
	cache["key"] = "value"
	fmt.Printf("  Cache: %d items\n", len(cache))

	fmt.Println("\nNote: Database not yet initialized (lazy!)")
	fmt.Println("Only initialize what you use!")

	fmt.Println()
}

// ============================================================================
// MAIN FUNCTION: Orchestrates All Demonstrations
// ============================================================================

// main executes all demonstration functions in order.
//
// MACRO-COMMENT: Learning Progression
// The demos are ordered to build understanding:
// 1. Basic sync.Once (foundation)
// 2. Lazy initialization (common pattern)
// 3. Singleton pattern (design pattern)
// 4. Error handling (production readiness)
// 5. Multiple singletons (real-world scenarios)
// 6. Race conditions (why sync.Once is needed)
// 7. Performance (fast path optimization)
// 8. Lazy fields (advanced pattern)
//
// RUN WITH RACE DETECTOR:
//   go run -race cmd/once-demo/main.go
//
// This will catch the race condition in demonstrateRaceCondition().
func main() {
	fmt.Println("sync.Once and Singleton Pattern Demonstrations")
	fmt.Println("==============================================\n")

	demonstrateBasicOnce()
	demonstrateLazyInitialization()
	demonstrateSingletonPattern()
	demonstrateErrorHandling()
	demonstrateMultipleSingletons()
	demonstrateRaceCondition()
	demonstratePerformance()
	demonstrateLazyFields()

	// MACRO-COMMENT: Key Insights
	// After running this program, you should understand:
	// 1. sync.Once guarantees exactly-once execution (thread-safe)
	// 2. Lazy initialization defers work until needed (fast startup)
	// 3. Singleton pattern ensures one instance (resource sharing)
	// 4. Naive approaches have race conditions (use sync.Once!)
	// 5. sync.Once fast path is extremely efficient (lock-free)
	// 6. Per-field sync.Once enables granular lazy initialization

	fmt.Println("==============================================")
	fmt.Println("Demonstrations complete!")
	fmt.Println("\nKey Takeaways:")
	fmt.Println("1. sync.Once guarantees exactly-once execution")
	fmt.Println("2. Use it for thread-safe lazy initialization")
	fmt.Println("3. Perfect for singleton pattern implementation")
	fmt.Println("4. Fast path is lock-free (very efficient)")
	fmt.Println("5. Handle initialization errors explicitly")
}
