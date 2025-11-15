// Package main demonstrates atomic operations from sync/atomic.
//
// MACRO-COMMENT: What This Program Teaches
// =========================================
// This program demonstrates:
// 1. Atomic Add operations (counters, increments)
// 2. Atomic Load/Store operations (reading/writing shared state)
// 3. Atomic Swap operations (atomic exchange)
// 4. Atomic CompareAndSwap (lock-free algorithms)
// 5. atomic.Value for arbitrary types
// 6. Performance comparison: Atomic vs Mutex
// 7. Common patterns (reference counting, flags, state machines)
//
// ATOMIC VS MUTEX:
// Atomic operations are CPU-level instructions that are lock-free and faster
// than mutexes for simple operations. However, they only work on specific types
// (int32, int64, uint32, uint64, uintptr, unsafe.Pointer) and provide limited
// composability compared to mutexes.

package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// ============================================================================
// SECTION 1: Atomic Add
// ============================================================================

// demonstrateAtomicAdd shows atomic increment operations.
//
// MACRO-COMMENT: Atomic Add
// AddInt64 performs an atomic read-modify-write operation in a single
// CPU instruction. It's equivalent to:
//   *addr += delta
//   return *addr
// But atomic (no race condition).
//
// WHY IT'S NEEDED:
// Regular increment (counter++) is three separate operations:
//   1. LOAD: temp = counter
//   2. ADD:  temp = temp + 1
//   3. STORE: counter = temp
// Two goroutines can interleave these operations, losing updates.
//
// PERFORMANCE:
// ~5-10 ns per operation (vs ~50-100 ns for mutex)
func demonstrateAtomicAdd() {
	fmt.Println("=== 1. Atomic Add ===")

	var counter int64
	const numGoroutines = 1000
	const incrementsPerGoroutine = 1000

	var wg sync.WaitGroup

	// MACRO-COMMENT: Launch goroutines that increment counter
	// Without atomic operations, this would be a race condition.
	// With atomics, all increments are correctly counted.
	fmt.Printf("Starting with counter = %d\n", counter)
	fmt.Printf("Launching %d goroutines, each incrementing %d times...\n",
		numGoroutines, incrementsPerGoroutine)

	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				// MICRO-COMMENT: Atomic increment
				// Returns new value (after increment)
				atomic.AddInt64(&counter, 1)
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	expected := int64(numGoroutines * incrementsPerGoroutine)
	fmt.Printf("Final counter: %d (expected: %d)\n", counter, expected)
	fmt.Printf("Time taken: %v\n", elapsed)

	if counter == expected {
		fmt.Println("✓ Correct! All increments were counted.")
	} else {
		fmt.Println("✗ RACE CONDITION! Some increments were lost.")
	}

	// MACRO-COMMENT: Decrement example
	// AddInt64 works with negative values too
	fmt.Println("\nDecrementing by 100...")
	newVal := atomic.AddInt64(&counter, -100)
	fmt.Printf("After decrement: %d (returned value: %d)\n",
		atomic.LoadInt64(&counter), newVal)

	fmt.Println()
}

// ============================================================================
// SECTION 2: Atomic Load and Store
// ============================================================================

// demonstrateLoadStore shows atomic read and write operations.
//
// MACRO-COMMENT: Load and Store
// LoadInt64 and StoreInt64 provide atomic read and write operations.
//
// WHY NEEDED:
// On 32-bit architectures, reading/writing int64 is NOT atomic without
// these operations. You might read half-updated values (torn reads).
//
// MEMORY ORDERING:
// Store has "release" semantics: all writes before Store are visible after Load
// Load has "acquire" semantics: all writes before Store are visible after Load
//
// PATTERN: Configuration Updates
// Writer uses Store to publish, readers use Load to observe.
func demonstrateLoadStore() {
	fmt.Println("=== 2. Atomic Load and Store ===")

	var config int64
	var stopFlag int64

	// MACRO-COMMENT: Writer goroutine
	// Simulates periodic configuration updates
	go func() {
		for version := int64(1); version <= 5; version++ {
			// MICRO-COMMENT: Check stop flag
			if atomic.LoadInt64(&stopFlag) == 1 {
				return
			}

			// MICRO-COMMENT: Atomically store new config version
			atomic.StoreInt64(&config, version)
			fmt.Printf("  Writer: Published config version %d\n", version)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// MACRO-COMMENT: Reader goroutine
	// Periodically reads current configuration
	go func() {
		lastSeen := int64(0)
		for {
			// MICRO-COMMENT: Check stop flag
			if atomic.LoadInt64(&stopFlag) == 1 {
				return
			}

			// MICRO-COMMENT: Atomically load current config
			current := atomic.LoadInt64(&config)
			if current != lastSeen {
				fmt.Printf("  Reader: Observed config version %d\n", current)
				lastSeen = current
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()

	// MICRO-COMMENT: Let reader and writer run
	time.Sleep(600 * time.Millisecond)

	// MICRO-COMMENT: Stop both goroutines
	fmt.Println("  Stopping...")
	atomic.StoreInt64(&stopFlag, 1)
	time.Sleep(100 * time.Millisecond)

	fmt.Printf("Final config version: %d\n", atomic.LoadInt64(&config))
	fmt.Println()
}

// ============================================================================
// SECTION 3: Atomic Swap
// ============================================================================

// demonstrateSwap shows atomic exchange operations.
//
// MACRO-COMMENT: Atomic Swap
// SwapInt64 atomically reads the old value and writes a new value.
// It's equivalent to:
//   old := *addr
//   *addr = new
//   return old
// But atomic.
//
// USE CASES:
// - Spinlocks (repeatedly swap 1, check if old was 0)
// - Resetting counters (swap 0, get old count)
// - State transitions
func demonstrateSwap() {
	fmt.Println("=== 3. Atomic Swap ===")

	var requestCount int64

	// MACRO-COMMENT: Simulate requests accumulating
	fmt.Println("Simulating incoming requests...")
	for i := 0; i < 100; i++ {
		atomic.AddInt64(&requestCount, 1)
	}

	// MACRO-COMMENT: Reset counter and get old value atomically
	// This pattern is useful for periodic reporting
	old := atomic.SwapInt64(&requestCount, 0)
	fmt.Printf("Processed %d requests (counter reset to 0)\n", old)
	fmt.Printf("Current counter: %d\n", atomic.LoadInt64(&requestCount))

	// MACRO-COMMENT: Spinlock example
	// This is a simple (but inefficient) lock implementation
	fmt.Println("\nSpinlock example:")
	var lock int64

	// MICRO-COMMENT: Try to acquire lock
	fmt.Println("  Acquiring lock...")
	for atomic.SwapInt64(&lock, 1) != 0 {
		// MICRO-COMMENT: Lock was held (old value = 1)
		// In real code, you'd use runtime.Gosched() here
		fmt.Println("    Lock held by another goroutine, spinning...")
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Println("  Lock acquired!")

	// MICRO-COMMENT: Critical section
	fmt.Println("  In critical section...")
	time.Sleep(50 * time.Millisecond)

	// MICRO-COMMENT: Release lock
	atomic.StoreInt64(&lock, 0)
	fmt.Println("  Lock released")

	fmt.Println()
}

// ============================================================================
// SECTION 4: Atomic CompareAndSwap (CAS)
// ============================================================================

// demonstrateCompareAndSwap shows CAS operations for lock-free algorithms.
//
// MACRO-COMMENT: CompareAndSwap (CAS)
// CompareAndSwapInt64 is the building block for lock-free algorithms.
//
// HOW IT WORKS:
//   if *addr == old {
//       *addr = new
//       return true
//   } else {
//       return false
//   }
// The entire operation is atomic (no race possible).
//
// PATTERN: CAS Loop
// For complex operations, use a CAS loop:
//   for {
//       old := atomic.LoadInt64(&value)
//       new := compute(old)
//       if atomic.CompareAndSwapInt64(&value, old, new) {
//           break  // Success
//       }
//       // Retry (value was modified by another goroutine)
//   }
func demonstrateCompareAndSwap() {
	fmt.Println("=== 4. Atomic CompareAndSwap ===")

	const (
		StateIdle    = 0
		StateRunning = 1
		StateStopped = 2
	)

	var state int64 = StateIdle

	// MACRO-COMMENT: Worker 1 tries to start
	// Only succeeds if state is Idle
	fmt.Printf("Initial state: %d (Idle)\n", atomic.LoadInt64(&state))

	success := atomic.CompareAndSwapInt64(&state, StateIdle, StateRunning)
	if success {
		fmt.Println("Worker 1: Successfully transitioned Idle → Running")
	} else {
		fmt.Println("Worker 1: Failed (state was not Idle)")
	}

	// MACRO-COMMENT: Worker 2 tries to start
	// Will fail because state is already Running
	success = atomic.CompareAndSwapInt64(&state, StateIdle, StateRunning)
	if success {
		fmt.Println("Worker 2: Successfully transitioned Idle → Running")
	} else {
		fmt.Printf("Worker 2: Failed (state is %d, not Idle)\n", atomic.LoadInt64(&state))
	}

	// MACRO-COMMENT: CAS Loop Example
	// Doubling a counter atomically (can't use AddInt64 for this)
	fmt.Println("\nCAS Loop Example: Doubling a counter")
	var counter int64 = 5
	fmt.Printf("Initial counter: %d\n", counter)

	for {
		old := atomic.LoadInt64(&counter)
		new := old * 2
		if atomic.CompareAndSwapInt64(&counter, old, new) {
			fmt.Printf("Successfully doubled: %d → %d\n", old, new)
			break
		}
		// In real code, another goroutine modified counter, retry
		fmt.Println("CAS failed, retrying...")
	}

	fmt.Println()
}

// ============================================================================
// SECTION 5: atomic.Value
// ============================================================================

// demonstrateAtomicValue shows using atomic.Value for arbitrary types.
//
// MACRO-COMMENT: atomic.Value
// atomic.Value allows atomic operations on any type (not just integers).
//
// METHODS:
//   Store(val any)              - Store a value
//   Load() any                   - Load the value
//   Swap(new any) (old any)      - Swap and return old
//   CompareAndSwap(old, new any) - CAS for any type
//
// RULES:
// 1. Must store consistent types (can't store int then string)
// 2. Type assertion required on Load
// 3. Stores pointer to value (not a copy)
func demonstrateAtomicValue() {
	fmt.Println("=== 5. atomic.Value ===")

	type Config struct {
		Timeout    int
		MaxRetries int
	}

	var config atomic.Value

	// MICRO-COMMENT: Store initial config
	config.Store(&Config{Timeout: 30, MaxRetries: 3})
	fmt.Println("Stored initial config: Timeout=30, MaxRetries=3")

	// MICRO-COMMENT: Load config
	cfg := config.Load().(*Config)
	fmt.Printf("Loaded config: Timeout=%d, MaxRetries=%d\n", cfg.Timeout, cfg.MaxRetries)

	// MACRO-COMMENT: Update config atomically
	// All readers will see either old or new config, never partial update
	fmt.Println("\nUpdating config...")
	config.Store(&Config{Timeout: 60, MaxRetries: 5})

	cfg = config.Load().(*Config)
	fmt.Printf("New config: Timeout=%d, MaxRetries=%d\n", cfg.Timeout, cfg.MaxRetries)

	// MACRO-COMMENT: Swap example
	oldCfg := config.Swap(&Config{Timeout: 90, MaxRetries: 10}).(*Config)
	fmt.Printf("Swapped: old Timeout=%d, new Timeout=%d\n",
		oldCfg.Timeout, config.Load().(*Config).Timeout)

	fmt.Println()
}

// ============================================================================
// SECTION 6: Performance Comparison
// ============================================================================

// benchmarkAtomic measures atomic operation performance.
func benchmarkAtomic(iterations int, numGoroutines int) time.Duration {
	var counter int64
	var wg sync.WaitGroup

	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations/numGoroutines; j++ {
				atomic.AddInt64(&counter, 1)
			}
		}()
	}

	wg.Wait()
	return time.Since(start)
}

// benchmarkMutex measures mutex operation performance.
func benchmarkMutex(iterations int, numGoroutines int) time.Duration {
	var counter int64
	var mu sync.Mutex
	var wg sync.WaitGroup

	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations/numGoroutines; j++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	return time.Since(start)
}

// demonstratePerformance compares atomic vs mutex performance.
//
// MACRO-COMMENT: Performance Comparison
// Atomic operations are significantly faster than mutexes for simple
// operations like incrementing a counter.
//
// TYPICAL RESULTS (8-core CPU):
//   Atomic:  50-100 ms  for 10M operations
//   Mutex:   200-500 ms for 10M operations
//   Speedup: 3-5x
//
// WHY ATOMIC IS FASTER:
// 1. No kernel involvement (mutex may need syscalls)
// 2. No lock contention (atomics use CPU's cache coherency protocol)
// 3. No context switching (goroutines don't block on atomics)
func demonstratePerformance() {
	fmt.Println("=== 6. Performance Comparison: Atomic vs Mutex ===")

	const iterations = 10_000_000
	numGoroutines := runtime.NumCPU()

	fmt.Printf("Running %d iterations with %d goroutines...\n\n",
		iterations, numGoroutines)

	// MICRO-COMMENT: Benchmark atomic operations
	fmt.Println("Benchmarking atomic operations...")
	atomicTime := benchmarkAtomic(iterations, numGoroutines)
	fmt.Printf("  Atomic:  %v\n", atomicTime)

	// MICRO-COMMENT: Benchmark mutex operations
	fmt.Println("Benchmarking mutex operations...")
	mutexTime := benchmarkMutex(iterations, numGoroutines)
	fmt.Printf("  Mutex:   %v\n", mutexTime)

	// MICRO-COMMENT: Calculate speedup
	speedup := float64(mutexTime) / float64(atomicTime)
	fmt.Printf("\nAtomic is %.1fx faster than mutex\n", speedup)

	fmt.Println()
}

// ============================================================================
// SECTION 7: Lock-Free Stack
// ============================================================================

// Node represents a node in the lock-free stack.
type Node struct {
	value int
	next  unsafe.Pointer // *Node
}

// LockFreeStack is a simple lock-free stack using CAS.
type LockFreeStack struct {
	head unsafe.Pointer // *Node
}

// Push adds a value to the stack (lock-free).
//
// MACRO-COMMENT: Lock-Free Push
// This implementation uses CAS to push a new node onto the stack.
//
// ALGORITHM:
// 1. Create new node
// 2. Read current head
// 3. Set new node's next to current head
// 4. Try to CAS head from old to new node
// 5. If CAS fails, retry (another goroutine modified head)
//
// This is a classic lock-free algorithm pattern.
func (s *LockFreeStack) Push(value int) {
	node := &Node{value: value}

	for {
		// MICRO-COMMENT: Read current head
		oldHead := atomic.LoadPointer(&s.head)

		// MICRO-COMMENT: Set new node's next to current head
		node.next = oldHead

		// MICRO-COMMENT: Try to update head to new node
		// If another goroutine modified head, CAS fails and we retry
		if atomic.CompareAndSwapPointer(&s.head, oldHead, unsafe.Pointer(node)) {
			return // Success
		}
		// CAS failed, retry (rare in low contention)
	}
}

// Pop removes and returns a value from the stack (lock-free).
//
// MACRO-COMMENT: Lock-Free Pop
// Similar to Push, but in reverse.
//
// ALGORITHM:
// 1. Read current head
// 2. If head is nil, stack is empty
// 3. Read next node
// 4. Try to CAS head from old to next
// 5. If CAS succeeds, return old head's value
// 6. If CAS fails, retry
func (s *LockFreeStack) Pop() (int, bool) {
	for {
		// MICRO-COMMENT: Read current head
		headPtr := atomic.LoadPointer(&s.head)
		if headPtr == nil {
			return 0, false // Stack is empty
		}

		// MICRO-COMMENT: Convert to *Node and read next
		head := (*Node)(headPtr)
		next := atomic.LoadPointer(&head.next)

		// MICRO-COMMENT: Try to update head to next
		if atomic.CompareAndSwapPointer(&s.head, headPtr, next) {
			return head.value, true // Success
		}
		// CAS failed, retry
	}
}

// demonstrateLockFreeStack shows a lock-free stack implementation.
//
// MACRO-COMMENT: Lock-Free Data Structures
// Lock-free algorithms use atomic operations (especially CAS) to avoid locks.
//
// BENEFITS:
// - No lock contention
// - No deadlocks
// - Better scalability under high concurrency
//
// CHALLENGES:
// - Complex to implement correctly
// - ABA problem (value changes A→B→A, CAS succeeds but state is different)
// - Memory reclamation (can't free nodes that other threads might be reading)
func demonstrateLockFreeStack() {
	fmt.Println("=== 7. Lock-Free Stack ===")

	stack := &LockFreeStack{}

	// MICRO-COMMENT: Push values
	fmt.Println("Pushing values: 1, 2, 3")
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	// MICRO-COMMENT: Pop values
	fmt.Print("Popping values: ")
	for {
		value, ok := stack.Pop()
		if !ok {
			break
		}
		fmt.Printf("%d ", value)
	}
	fmt.Println("(LIFO order)")

	// MACRO-COMMENT: Concurrent stress test
	fmt.Println("\nConcurrent stress test: 8 goroutines, 1000 ops each")
	const numOps = 1000
	var wg sync.WaitGroup

	start := time.Now()

	// MICRO-COMMENT: Half push, half pop
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				stack.Push(id*1000 + j)
			}
		}(i)
	}

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				stack.Pop()
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	fmt.Printf("Completed in %v\n", elapsed)
	fmt.Println("All operations completed without locks!")

	fmt.Println()
}

// ============================================================================
// SECTION 8: Common Patterns
// ============================================================================

// demonstrateReferenceCounting shows atomic reference counting.
//
// MACRO-COMMENT: Reference Counting Pattern
// Use atomic operations to track how many references exist to a resource.
// When count reaches zero, clean up the resource.
//
// This is how sync.WaitGroup works internally!
func demonstrateReferenceCounting() {
	fmt.Println("=== 8. Reference Counting Pattern ===")

	type Resource struct {
		refCount int64
		name     string
	}

	cleanup := func(r *Resource) {
		fmt.Printf("  Cleaning up resource '%s' (refCount = 0)\n", r.name)
	}

	resource := &Resource{refCount: 0, name: "database-connection"}

	// MICRO-COMMENT: Acquire references
	fmt.Println("Acquiring 3 references...")
	for i := 0; i < 3; i++ {
		count := atomic.AddInt64(&resource.refCount, 1)
		fmt.Printf("  Reference %d acquired (count = %d)\n", i+1, count)
	}

	// MICRO-COMMENT: Release references
	fmt.Println("\nReleasing references...")
	for i := 0; i < 3; i++ {
		count := atomic.AddInt64(&resource.refCount, -1)
		fmt.Printf("  Reference %d released (count = %d)\n", i+1, count)

		if count == 0 {
			cleanup(resource)
		}
	}

	fmt.Println()
}

// ============================================================================
// MAIN FUNCTION
// ============================================================================

func main() {
	fmt.Println("=== Atomic Operations Demonstration ===\n")

	// SECTION 1: Basic atomic operations
	demonstrateAtomicAdd()
	demonstrateLoadStore()
	demonstrateSwap()
	demonstrateCompareAndSwap()

	// SECTION 2: Advanced features
	demonstrateAtomicValue()
	demonstratePerformance()

	// SECTION 3: Real-world patterns
	demonstrateLockFreeStack()
	demonstrateReferenceCounting()

	fmt.Println("=== All Demonstrations Complete ===")
}
