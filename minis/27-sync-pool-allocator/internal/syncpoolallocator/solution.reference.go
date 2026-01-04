//go:build reference

package syncpoolallocator

/*
Problem: Implement memory pool allocators using sync.Pool for zero-allocation reuse

Requirements:
1. Reuse allocated objects to reduce GC pressure
2. Thread-safe pool access (sync.Pool is already concurrent-safe per-P)
3. Proper object reset before reuse (critical for correctness)
4. Metrics tracking for pool efficiency
5. Size-classed pools for different allocation sizes
6. Bounded pools to limit memory usage
7. Generic pools with custom reset logic

Data Structure:
- sync.Pool: Per-P (processor) local storage with lock-free fast path
- Atomic counters: Thread-safe metrics without mutex overhead
- Semaphores: Bounded pool capacity control
- Size classes: Multiple pools for different allocation sizes

Time/Space Complexity:
- Get: O(1) average (per-P local storage, no locks)
- Put: O(1) average (per-P local storage, no locks)
- Space: O(number of P's × pool depth) per pool
- Reset: O(1) typically (depends on reset function)

Algorithm: sync.Pool Internal Mechanics
- Each P (Go processor) has local storage (lock-free)
- Fast path: Get from local storage (no synchronization)
- Slow path: Steal from other P's or create new (with locking)
- Objects can be GC'd between GC cycles (no guarantees)

Why sync.Pool is powerful:
- Per-P local storage eliminates contention
- Lock-free fast path for common case
- Automatic cleanup (GC can free unused objects)
- Zero-cost when pool is empty
- Works seamlessly with Go's scheduler

Critical Best Practices:
1. Always reset objects before Put (security/correctness)
2. Don't assume objects survive across GC cycles
3. Reset should return object to initial state
4. For bounded pools, track in-use count
5. Metrics help understand pool efficiency
*/

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
)

// ============================================================================
// Exercise 1: Basic Buffer Pool
// ============================================================================

/*
BufferPool wraps sync.Pool for bytes.Buffer reuse.

Use Case:
- Frequent buffer allocations in hot paths (HTTP handlers, serialization)
- Reduces GC pressure by reusing allocated buffers
- Critical: Always call Reset() before Put()

Why bytes.Buffer needs reset:
- Buffer maintains internal slice that grows over time
- Previous data must be cleared for security/correctness
- Reset() reuses underlying slice capacity
*/

// BufferPool provides a pool of bytes.Buffer instances.
// BREAKPOINT: Set breakpoint here to inspect pool initialization
type BufferPool struct {
	pool sync.Pool // Thread-safe per-P pool
}

// NewBufferPool creates a new buffer pool.
// BREAKPOINT: Set breakpoint here to trace pool creation
// DEBUG: Watch return value to verify pool.New function is set
func NewBufferPool() *BufferPool {
	// BREAKPOINT: Set breakpoint here before returning
	// DEBUG: Watch 'pool.New' function - called when pool is empty
	return &BufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				// BREAKPOINT: Set breakpoint here when new buffer is created
				// DEBUG: This is called when pool has no available buffers
				// DEBUG: Watch return value - a fresh bytes.Buffer
				return new(bytes.Buffer)
			},
		},
	}
}

// Get retrieves a buffer from the pool or creates a new one.
// BREAKPOINT: Set breakpoint here to trace buffer retrieval
// DEBUG: Watch return value - may be reused or newly allocated
// DEBUG: Fast path: Gets from per-P local storage (lock-free)
func (bp *BufferPool) Get() *bytes.Buffer {
	// BREAKPOINT: Set breakpoint here before pool.Get()
	// DEBUG: sync.Pool.Get() is thread-safe and lock-free on fast path
	buf := bp.pool.Get().(*bytes.Buffer)
	// DEBUG: Watch 'buf' - may have capacity from previous use
	// DEBUG: Length should be checked (should be 0 after Reset, but verify)
	return buf
}

// Put returns a buffer to the pool after resetting it.
// BREAKPOINT: Set breakpoint here to trace buffer return
// DEBUG: Watch 'buf' before reset - may contain previous data
func (bp *BufferPool) Put(buf *bytes.Buffer) {
	// BREAKPOINT: Set breakpoint here before reset
	// DEBUG: CRITICAL: Always reset before putting back
	// DEBUG: Watch 'buf.Len()' and 'buf.Cap()' before reset
	buf.Reset()
	// BREAKPOINT: Set breakpoint here after reset
	// DEBUG: Watch 'buf.Len()' - should be 0
	// DEBUG: Watch 'buf.Cap()' - may retain capacity (good for reuse)
	bp.pool.Put(buf)
	// DEBUG: Buffer is now in pool for reuse (per-P local storage)
}

// ============================================================================
// Exercise 2: Slice Pool
// ============================================================================

/*
SlicePool wraps sync.Pool for byte slice reuse with fixed capacity.

Use Case:
- Frequent slice allocations of known capacity
- Serialization, encoding, temporary storage
- Reset preserves capacity, reduces allocations

Why slice reset matters:
- Slice length must be reset ([:0] keeps capacity)
- Capacity is preserved for reuse efficiency
- Zero-length slice with capacity is ideal for append
*/

// SlicePool provides a pool of byte slices with fixed capacity.
// BREAKPOINT: Set breakpoint here to inspect slice pool
type SlicePool struct {
	pool     sync.Pool
	capacity int // Target capacity for slices
}

// NewSlicePool creates a slice pool with the given capacity.
// BREAKPOINT: Set breakpoint here to trace slice pool creation
// DEBUG: Watch 'capacity' parameter to verify initialization
func NewSlicePool(capacity int) *SlicePool {
	// BREAKPOINT: Set breakpoint here before returning
	// DEBUG: Watch 'capacity' stored for reference
	return &SlicePool{
		capacity: capacity,
		pool: sync.Pool{
			New: func() interface{} {
				// BREAKPOINT: Set breakpoint here when new slice is created
				// DEBUG: Watch slice creation - length 0, capacity 'capacity'
				// DEBUG: Returns pointer to slice (required for sync.Pool)
				slice := make([]byte, 0, capacity)
				return &slice
			},
		},
	}
}

// Get retrieves a slice from the pool or creates a new one.
// BREAKPOINT: Set breakpoint here to trace slice retrieval
// DEBUG: Watch return value - pointer to byte slice
func (sp *SlicePool) Get() *[]byte {
	// BREAKPOINT: Set breakpoint here before pool.Get()
	slicePtr := sp.pool.Get().(*[]byte)
	// DEBUG: Watch '*slicePtr' - length should be 0
	// DEBUG: Watch 'cap(*slicePtr)' - should be >= sp.capacity
	return slicePtr
}

// Put returns a slice to the pool after resetting it.
// BREAKPOINT: Set breakpoint here to trace slice return
// DEBUG: Watch '*slice' before reset - may have data
func (sp *SlicePool) Put(slice *[]byte) {
	// BREAKPOINT: Set breakpoint here before reset
	// DEBUG: Watch 'len(*slice)' and 'cap(*slice)' before reset
	// Reset length to 0, keep capacity (critical optimization)
	*slice = (*slice)[:0] // Reset length, keep capacity
	// BREAKPOINT: Set breakpoint here after reset
	// DEBUG: Watch 'len(*slice)' - should be 0
	// DEBUG: Watch 'cap(*slice)' - capacity preserved (good!)
	sp.pool.Put(slice)
	// DEBUG: Slice is now available for reuse
}

// ============================================================================
// Exercise 3: Generic Typed Pool
// ============================================================================

/*
Pool is a generic wrapper around sync.Pool with custom reset logic.

Use Case:
- Type-safe pool for any type T
- Custom reset functions for complex objects
- Reusable pattern for any pooled type

Why generics:
- Type safety at compile time
- No interface{} type assertions
- Cleaner API for specific types
*/

// Pool is a generic thread-safe pool with reset functionality.
// BREAKPOINT: Set breakpoint here to inspect generic pool
type Pool[T any] struct {
	pool  sync.Pool
	reset func(*T) // Optional reset function
}

// NewPool creates a generic pool with custom new and reset functions.
// BREAKPOINT: Set breakpoint here to trace generic pool creation
// DEBUG: Watch 'newFunc' and 'resetFunc' parameters
func NewPool[T any](newFunc func() *T, resetFunc func(*T)) *Pool[T] {
	// BREAKPOINT: Set breakpoint here before returning
	// DEBUG: Watch 'reset' field - may be nil if no reset needed
	return &Pool[T]{
		pool: sync.Pool{
			New: func() interface{} {
				// BREAKPOINT: Set breakpoint here when new object is created
				// DEBUG: Calls user-provided newFunc for type-specific creation
				return newFunc()
			},
		},
		reset: resetFunc,
	}
}

// Get retrieves an object from the pool or creates a new one.
// BREAKPOINT: Set breakpoint here to trace object retrieval
// DEBUG: Watch return value - may be reused or newly created
func (p *Pool[T]) Get() *T {
	// BREAKPOINT: Set breakpoint here before pool.Get()
	obj := p.pool.Get().(*T)
	// DEBUG: Watch 'obj' - type-safe access without interface{} assertion
	return obj
}

// Put returns an object to the pool, resetting it first if needed.
// BREAKPOINT: Set breakpoint here to trace object return
// DEBUG: Watch 'obj' before reset - may contain previous state
func (p *Pool[T]) Put(obj *T) {
	// BREAKPOINT: Set breakpoint here before reset
	// DEBUG: Watch 'p.reset' - may be nil (no reset needed)
	if p.reset != nil {
		// BREAKPOINT: Set breakpoint here when reset is called
		// DEBUG: Calls user-provided reset function
		// DEBUG: Watch 'obj' before and after reset
		p.reset(obj)
	}
	p.pool.Put(obj)
	// DEBUG: Object is now available for reuse
}

// ============================================================================
// Exercise 4: Pool with Metrics
// ============================================================================

/*
MetricsPool tracks pool efficiency with atomic counters.

Use Case:
- Monitor pool hit rate (reuse vs new allocations)
- Understand pool effectiveness
- Debug memory allocation patterns
- Performance optimization insights

Why atomic counters:
- Thread-safe without mutex overhead
- Lock-free metrics collection
- Minimal performance impact
*/

// MetricsPool tracks pool statistics using atomic counters.
// BREAKPOINT: Set breakpoint here to inspect metrics pool
type MetricsPool struct {
	pool sync.Pool
	gets atomic.Int64 // Total Get() calls
	puts atomic.Int64 // Total Put() calls
	news atomic.Int64 // Total new allocations (New() calls)
}

// PoolStats contains pool efficiency metrics.
type PoolStats struct {
	Gets    int64   // Total Get() calls
	Puts    int64   // Total Put() calls
	News    int64   // Total new allocations
	HitRate float64 // Percentage of Gets that reused objects (0-100)
}

// NewMetricsPool creates a metrics pool with the given new function.
// BREAKPOINT: Set breakpoint here to trace metrics pool creation
// DEBUG: Watch 'newFunc' parameter
func NewMetricsPool(newFunc func() interface{}) *MetricsPool {
	mp := &MetricsPool{}
	mp.pool.New = func() interface{} {
		// BREAKPOINT: Set breakpoint here when new object is created
		// DEBUG: Increment news counter (atomic, thread-safe)
		mp.news.Add(1)
		// DEBUG: Watch 'mp.news' increment
		return newFunc()
	}
	// DEBUG: Watch return value - pool.New is now instrumented
	return mp
}

// Get retrieves an object and increments the gets counter.
// BREAKPOINT: Set breakpoint here to trace metrics collection
// DEBUG: Watch 'mp.gets' increment atomically
func (mp *MetricsPool) Get() interface{} {
	// BREAKPOINT: Set breakpoint here before incrementing counter
	// DEBUG: Atomic increment (lock-free, thread-safe)
	mp.gets.Add(1)
	// DEBUG: Watch 'mp.gets' value
	obj := mp.pool.Get()
	// DEBUG: Watch 'obj' - may trigger news increment if pool empty
	return obj
}

// Put returns an object and increments the puts counter.
// BREAKPOINT: Set breakpoint here to trace Put metrics
// DEBUG: Watch 'mp.puts' increment atomically
func (mp *MetricsPool) Put(obj interface{}) {
	// BREAKPOINT: Set breakpoint here before incrementing counter
	// DEBUG: Atomic increment (lock-free, thread-safe)
	mp.puts.Add(1)
	// DEBUG: Watch 'mp.puts' value
	mp.pool.Put(obj)
}

// Stats returns current pool statistics.
// BREAKPOINT: Set breakpoint here to trace statistics calculation
// DEBUG: Watch all counter values and hit rate calculation
func (mp *MetricsPool) Stats() PoolStats {
	// BREAKPOINT: Set breakpoint here before loading counters
	// DEBUG: Atomic loads ensure consistent snapshot
	gets := mp.gets.Load()
	puts := mp.puts.Load()
	news := mp.news.Load()

	// BREAKPOINT: Set breakpoint here before hit rate calculation
	// DEBUG: Watch 'gets', 'news' to understand hit rate formula
	// Hit rate = (Gets - News) / Gets * 100
	// If gets=100, news=20, then 80% hit rate (80 reuses)
	var hitRate float64
	if gets > 0 {
		// BREAKPOINT: Set breakpoint here during calculation
		// DEBUG: Watch calculation - reuses / total gets * 100
		hitRate = float64(gets-news) / float64(gets) * 100
	}
	// DEBUG: Watch 'hitRate' - higher is better (more reuse)

	return PoolStats{
		Gets:    gets,
		Puts:    puts,
		News:    news,
		HitRate: hitRate,
	}
}

// ============================================================================
// Exercise 5: Size-Classed Buffer Pool
// ============================================================================

/*
SizeClassedPool uses multiple pools for different allocation sizes.

Use Case:
- Allocations vary in size (1KB, 4KB, 16KB, 64KB)
- Avoids wasting memory (don't use 64KB pool for 1KB needs)
- Common pattern in memory allocators (jemalloc, tcmalloc)

Why size classes:
- Reduces memory waste
- Better cache locality
- More predictable performance
- Matches allocation size to pool size
*/

// SizeClassedPool maintains multiple pools for different size classes.
// BREAKPOINT: Set breakpoint here to inspect size-classed pool
type SizeClassedPool struct {
	pools [4]sync.Pool // Array of pools for different sizes
}

// NewSizeClassedPool creates a pool with 4 size classes (1KB, 4KB, 16KB, 64KB).
// BREAKPOINT: Set breakpoint here to trace size-classed pool creation
// DEBUG: Watch each pool initialization
func NewSizeClassedPool() *SizeClassedPool {
	scp := &SizeClassedPool{}

	// 1KB pool (index 0)
	// BREAKPOINT: Set breakpoint here for 1KB pool setup
	// DEBUG: Watch capacity 1024
	scp.pools[0].New = func() interface{} {
		buf := make([]byte, 0, 1024)
		return &buf
	}

	// 4KB pool (index 1)
	// BREAKPOINT: Set breakpoint here for 4KB pool setup
	// DEBUG: Watch capacity 4096
	scp.pools[1].New = func() interface{} {
		buf := make([]byte, 0, 4096)
		return &buf
	}

	// 16KB pool (index 2)
	// BREAKPOINT: Set breakpoint here for 16KB pool setup
	// DEBUG: Watch capacity 16384
	scp.pools[2].New = func() interface{} {
		buf := make([]byte, 0, 16384)
		return &buf
	}

	// 64KB pool (index 3)
	// BREAKPOINT: Set breakpoint here for 64KB pool setup
	// DEBUG: Watch capacity 65536
	scp.pools[3].New = func() interface{} {
		buf := make([]byte, 0, 65536)
		return &buf
	}

	return scp
}

// Get retrieves a buffer from the appropriate size class.
// BREAKPOINT: Set breakpoint here to trace size class selection
// DEBUG: Watch 'size' parameter and pool selection logic
func (scp *SizeClassedPool) Get(size int) *[]byte {
	// BREAKPOINT: Set breakpoint here before size class selection
	// DEBUG: Watch 'size' to determine which pool to use
	var poolIdx int
	switch {
	case size <= 1024:
		// BREAKPOINT: Hit when size <= 1KB
		poolIdx = 0
	case size <= 4096:
		// BREAKPOINT: Hit when size <= 4KB
		poolIdx = 1
	case size <= 16384:
		// BREAKPOINT: Hit when size <= 16KB
		poolIdx = 2
	default:
		// BREAKPOINT: Hit when size > 16KB
		poolIdx = 3
	}
	// DEBUG: Watch 'poolIdx' to see which pool is selected

	// BREAKPOINT: Set breakpoint here before getting from selected pool
	buf := scp.pools[poolIdx].Get().(*[]byte)
	// DEBUG: Watch '*buf' - capacity should match size class
	return buf
}

// Put returns a buffer to the appropriate size class based on its capacity.
// BREAKPOINT: Set breakpoint here to trace size class determination
// DEBUG: Watch buffer capacity to determine return pool
func (scp *SizeClassedPool) Put(buf *[]byte) {
	// BREAKPOINT: Set breakpoint here before reset
	// Reset length but preserve capacity
	*buf = (*buf)[:0] // Reset length, keep capacity

	// BREAKPOINT: Set breakpoint here before determining size class
	// DEBUG: Watch 'cap(*buf)' to determine which pool to return to
	capacity := cap(*buf)
	var poolIdx int
	switch {
	case capacity <= 1024:
		// BREAKPOINT: Hit when capacity <= 1KB
		poolIdx = 0
	case capacity <= 4096:
		// BREAKPOINT: Hit when capacity <= 4KB
		poolIdx = 1
	case capacity <= 16384:
		// BREAKPOINT: Hit when capacity <= 16KB
		poolIdx = 2
	default:
		// BREAKPOINT: Hit when capacity > 16KB
		poolIdx = 3
	}
	// DEBUG: Watch 'poolIdx' to see which pool receives the buffer

	// BREAKPOINT: Set breakpoint here before putting in selected pool
	scp.pools[poolIdx].Put(buf)
}

// ============================================================================
// Exercise 6: Bounded Pool with Semaphore
// ============================================================================

/*
BoundedPool limits the maximum number of objects in use.

Use Case:
- Prevent unbounded memory growth
- Control resource usage
- Backpressure when pool is exhausted

Why semaphore:
- Limits concurrent in-use count
- Blocks when limit reached (backpressure)
- Simple bounded counter implementation
*/

// BoundedPool limits the number of objects that can be in use simultaneously.
// BREAKPOINT: Set breakpoint here to inspect bounded pool
type BoundedPool struct {
	pool      sync.Pool
	semaphore chan struct{} // Semaphore for bounding
	maxSize   int           // Maximum in-use count
}

// NewBoundedPool creates a bounded pool with the given limit.
// BREAKPOINT: Set breakpoint here to trace bounded pool creation
// DEBUG: Watch 'maxSize' parameter
func NewBoundedPool(maxSize int, newFunc func() interface{}) *BoundedPool {
	// BREAKPOINT: Set breakpoint here before returning
	// DEBUG: Watch 'semaphore' initialization - buffered channel
	return &BoundedPool{
		pool: sync.Pool{
			New: newFunc,
		},
		semaphore: make(chan struct{}, maxSize), // Buffered channel as semaphore
		maxSize:   maxSize,
	}
}

// Get retrieves an object, blocking if the pool limit is reached.
// BREAKPOINT: Set breakpoint here to trace bounded Get
// DEBUG: Watch semaphore acquire operation (may block)
func (bp *BoundedPool) Get() interface{} {
	// BREAKPOINT: Set breakpoint here before semaphore acquire
	// DEBUG: Sending to semaphore acquires permit (blocks if full)
	bp.semaphore <- struct{}{} // Acquire semaphore (blocks if full)
	// BREAKPOINT: Set breakpoint here after acquiring permit
	// DEBUG: Permit acquired - now safe to get from pool
	obj := bp.pool.Get()
	// DEBUG: Watch 'obj' - obtained from pool
	return obj
}

// Put returns an object and releases the semaphore permit.
// BREAKPOINT: Set breakpoint here to trace bounded Put
// DEBUG: Watch semaphore release operation
func (bp *BoundedPool) Put(obj interface{}) {
	// BREAKPOINT: Set breakpoint here before putting in pool
	bp.pool.Put(obj)
	// BREAKPOINT: Set breakpoint here before releasing permit
	// DEBUG: Receiving from semaphore releases permit
	<-bp.semaphore // Release semaphore
	// DEBUG: Permit released - another Get() can now proceed
}

// InUse returns the current number of objects in use.
// BREAKPOINT: Set breakpoint here to trace in-use count
// DEBUG: Watch return value - number of acquired permits
func (bp *BoundedPool) InUse() int {
	// BREAKPOINT: Set breakpoint here before reading semaphore length
	// DEBUG: Length of semaphore channel = number of in-use objects
	return len(bp.semaphore)
}

// ============================================================================
// Exercise 8: Worker Pool Pattern
// ============================================================================

/*
WorkerPool demonstrates pooling complex objects (workers with multiple resources).

Use Case:
- Workers have multiple pooled resources (buffers, temp storage)
- Reduces allocation overhead for complex objects
- Common in request handlers, processors

Why worker pools:
- Multiple resources per worker (buffer + temp slice)
- Reset must clear all resources
- More efficient than pooling resources separately
*/

// Worker contains resources that can be pooled together.
// BREAKPOINT: Set breakpoint here to inspect worker structure
type Worker struct {
	buf  *bytes.Buffer // Pooled buffer
	temp []byte        // Pooled temp slice
}

// WorkerPool pools workers for processing tasks.
// BREAKPOINT: Set breakpoint here to inspect worker pool
type WorkerPool struct {
	pool sync.Pool
}

// NewWorkerPool creates a pool of workers.
// BREAKPOINT: Set breakpoint here to trace worker pool creation
// DEBUG: Watch worker initialization with multiple resources
func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		pool: sync.Pool{
			New: func() interface{} {
				// BREAKPOINT: Set breakpoint here when new worker is created
				// DEBUG: Watch worker creation with buffer and temp slice
				return &Worker{
					buf:  new(bytes.Buffer),
					temp: make([]byte, 0, 1024),
				}
			},
		},
	}
}

// Process processes data using a pooled worker.
// BREAKPOINT: Set breakpoint here to trace worker processing
// DEBUG: Watch worker retrieval, usage, and return
func (wp *WorkerPool) Process(data string) string {
	// BREAKPOINT: Set breakpoint here before getting worker
	worker := wp.pool.Get().(*Worker)
	// DEBUG: Watch 'worker' - may be reused or newly created
	defer wp.pool.Put(worker) // Ensure worker is always returned

	// BREAKPOINT: Set breakpoint here before reset
	// DEBUG: Reset all worker resources before use
	worker.Reset()

	// BREAKPOINT: Set breakpoint here during processing
	// DEBUG: Watch 'worker.buf' usage
	worker.buf.WriteString("Processed: ")
	worker.buf.WriteString(data)

	// BREAKPOINT: Set breakpoint here during temp buffer usage
	// DEBUG: Watch 'worker.temp' append operation
	worker.temp = append(worker.temp, []byte(" [transformed]")...)

	// BREAKPOINT: Set breakpoint here before returning result
	result := worker.buf.String() + string(worker.temp)
	// DEBUG: Watch 'result' - combines both worker resources
	return result
}

// Reset clears all worker resources for reuse.
// BREAKPOINT: Set breakpoint here to trace worker reset
// DEBUG: Watch all resources being cleared
func (w *Worker) Reset() {
	// BREAKPOINT: Set breakpoint here before buffer reset
	if w.buf != nil {
		// DEBUG: Reset buffer (clears length, keeps capacity)
		w.buf.Reset()
	}
	// BREAKPOINT: Set breakpoint here before temp slice reset
	if w.temp != nil {
		// DEBUG: Reset slice length (keeps capacity)
		w.temp = w.temp[:0]
	}
	// DEBUG: All resources are now reset and ready for reuse
}

// ============================================================================
// Bonus: Enhanced Metrics Pool
// ============================================================================

/*
EnhancedMetricsPool provides detailed statistics including reuses.

Use Case:
- More detailed pool analysis
- Understanding reuse patterns
- Performance optimization insights

Why enhanced metrics:
- Tracks reuses explicitly
- Calculates miss rate and efficiency
- Better understanding of pool behavior
*/

// EnhancedMetricsPool tracks detailed pool statistics.
// BREAKPOINT: Set breakpoint here to inspect enhanced metrics
type EnhancedMetricsPool struct {
	pool   sync.Pool
	gets   atomic.Int64 // Total Get() calls
	puts   atomic.Int64 // Total Put() calls
	news   atomic.Int64 // Total new allocations
	reuses atomic.Int64 // Total reuses (gets - news)
}

// EnhancedStats contains detailed pool efficiency metrics.
type EnhancedStats struct {
	Gets       int64   // Total Get() calls
	Puts       int64   // Total Put() calls
	News       int64   // Total new allocations
	Reuses     int64   // Total object reuses
	HitRate    float64 // Percentage of Gets that reused (0-100)
	MissRate   float64 // Percentage of Gets that allocated new (0-100)
	Efficiency float64 // Percentage of Puts that were reused (0-100)
}

// NewEnhancedMetricsPool creates an enhanced metrics pool.
// BREAKPOINT: Set breakpoint here to trace enhanced pool creation
func NewEnhancedMetricsPool(newFunc func() interface{}) *EnhancedMetricsPool {
	emp := &EnhancedMetricsPool{}
	emp.pool.New = func() interface{} {
		// BREAKPOINT: Set breakpoint here when new object created
		emp.news.Add(1)
		return newFunc()
	}
	return emp
}

// Get retrieves an object and tracks metrics.
// BREAKPOINT: Set breakpoint here to trace enhanced Get
// DEBUG: Watch metrics collection and reuse tracking
func (emp *EnhancedMetricsPool) Get() interface{} {
	emp.gets.Add(1)
	obj := emp.pool.Get()

	// BREAKPOINT: Set breakpoint here before reuse calculation
	// DEBUG: Try to determine if this was a reuse
	// Note: This is an approximation - in practice, track more precisely
	currentNews := emp.news.Load()
	currentGets := emp.gets.Load()
	if currentGets > currentNews {
		// BREAKPOINT: Set breakpoint here when reuse detected
		emp.reuses.Add(1)
	}

	return obj
}

// Put returns an object and tracks metrics.
// BREAKPOINT: Set breakpoint here to trace enhanced Put
func (emp *EnhancedMetricsPool) Put(obj interface{}) {
	emp.puts.Add(1)
	emp.pool.Put(obj)
}

// Stats returns detailed pool statistics.
// BREAKPOINT: Set breakpoint here to trace statistics calculation
// DEBUG: Watch all metrics calculations
func (emp *EnhancedMetricsPool) Stats() EnhancedStats {
	// BREAKPOINT: Set breakpoint here before loading counters
	gets := emp.gets.Load()
	puts := emp.puts.Load()
	news := emp.news.Load()
	reuses := emp.reuses.Load()

	// BREAKPOINT: Set breakpoint here before rate calculations
	var hitRate, missRate, efficiency float64
	if gets > 0 {
		// BREAKPOINT: Set breakpoint here during hit rate calculation
		hitRate = float64(gets-news) / float64(gets) * 100
		// DEBUG: Hit rate = reuses / total gets
		missRate = float64(news) / float64(gets) * 100
		// DEBUG: Miss rate = new allocations / total gets
	}

	if puts > 0 {
		// BREAKPOINT: Set breakpoint here during efficiency calculation
		efficiency = float64(reuses) / float64(puts) * 100
		// DEBUG: Efficiency = reuses / total puts
	}

	return EnhancedStats{
		Gets:       gets,
		Puts:       puts,
		News:       news,
		Reuses:     reuses,
		HitRate:    hitRate,
		MissRate:   missRate,
		Efficiency: efficiency,
	}
}

// String formats enhanced statistics for display.
// BREAKPOINT: Set breakpoint here to trace string formatting
func (es EnhancedStats) String() string {
	return fmt.Sprintf(
		"Gets: %d, Puts: %d, News: %d, Reuses: %d, Hit Rate: %.1f%%, Miss Rate: %.1f%%, Efficiency: %.1f%%",
		es.Gets, es.Puts, es.News, es.Reuses, es.HitRate, es.MissRate, es.Efficiency,
	)
}

// ============================================================================
// Bonus: Copy-on-Write Pool Pattern
// ============================================================================

/*
CopyOnWritePool returns copies of objects to prevent mutation of pooled objects.

Use Case:
- Objects must not be mutated in place
- Safety when multiple goroutines use same object
- Avoid corruption of pooled objects

Why copy-on-write:
- Prevents accidental mutation of pooled objects
- Safer for concurrent access
- Trade-off: Copy cost vs safety
*/

// CopyOnWritePool returns copies of pooled objects.
// BREAKPOINT: Set breakpoint here to inspect copy-on-write pool
type CopyOnWritePool[T any] struct {
	pool sync.Pool
	copy func(*T) *T // Copy function for type T
}

// NewCopyOnWritePool creates a copy-on-write pool.
// BREAKPOINT: Set breakpoint here to trace copy-on-write pool creation
// DEBUG: Watch 'newFunc' and 'copyFunc' parameters
func NewCopyOnWritePool[T any](newFunc func() *T, copyFunc func(*T) *T) *CopyOnWritePool[T] {
	return &CopyOnWritePool[T]{
		pool: sync.Pool{
			New: func() interface{} {
				// BREAKPOINT: Set breakpoint here when new object created
				return newFunc()
			},
		},
		copy: copyFunc,
	}
}

// Get retrieves an object and returns a copy.
// BREAKPOINT: Set breakpoint here to trace copy-on-write Get
// DEBUG: Watch copy operation
func (cp *CopyOnWritePool[T]) Get() *T {
	obj := cp.pool.Get().(*T)
	// BREAKPOINT: Set breakpoint here before copy
	// DEBUG: Return a copy so original stays in pool
	if cp.copy != nil {
		// BREAKPOINT: Set breakpoint here when copying
		// DEBUG: Watch copy function call - creates new instance
		return cp.copy(obj)
	}
	return obj
}

// Put returns an object to the pool.
// BREAKPOINT: Set breakpoint here to trace copy-on-write Put
func (cp *CopyOnWritePool[T]) Put(obj *T) {
	// BREAKPOINT: Set breakpoint here before putting in pool
	// DEBUG: Original object (unmutated) goes back to pool
	cp.pool.Put(obj)
}
