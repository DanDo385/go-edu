//go:build !solution
// +build !solution

package exercise

import (
	"bytes"
)

// Exercise 1: Basic Buffer Pool
type BufferPool struct {
	pool sync.Pool
}

func NewBufferPool() *BufferPool {
	// TODO: Implement this function.
	// - `sync.Pool` needs a `New` function to be set. This function is called when `Get` is called on an empty pool.
	// - It should return a pointer to a new `bytes.Buffer`.
	return &BufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

func (bp *BufferPool) Get() *bytes.Buffer {
	// TODO: Implement this function.
	// - Get an item from the pool using `bp.pool.Get()`.
	// - The item returned from the pool is of type `interface{}`, so you must perform a type assertion to convert it to `*bytes.Buffer`.
	// - `return bp.pool.Get().(*bytes.Buffer)`
	return nil
}

func (bp *BufferPool) Put(buf *bytes.Buffer) {
	// TODO: Implement this function.
	// - It is crucial to reset the buffer before putting it back in the pool.
	// - If you don't reset it, the next user of the buffer will see the old data.
	// - `buf.Reset()`
	// - Then, put the buffer back into the pool: `bp.pool.Put(buf)`.
}

// Exercise 2: Slice Pool
type SlicePool struct {
	pool     sync.Pool
	capacity int
}

func NewSlicePool(capacity int) *SlicePool {
	// TODO: Implement this function.
	// - The `New` function for the pool should create a `*[]byte`.
	// - The slice should be created with a length of 0 but the specified `capacity`.
	// - `slice := make([]byte, 0, capacity)`
	// - The pool must store a pointer to the slice (`&slice`) because a slice header itself is a value type, and modifications to it wouldn't be seen across `Get`/`Put` calls otherwise.
	return nil
}

func (sp *SlicePool) Get() *[]byte {
	// TODO: Implement this function.
	// - Get a `*[]byte` from the pool and type-assert it.
	return nil
}

func (sp *SlicePool) Put(slice *[]byte) {
	// TODO: Implement this function.
	// - To "reset" a slice while keeping its capacity, you set its length to 0.
	// - `*slice = (*slice)[:0]`
	// - Then, put the pointer back into the pool.
}

// Exercise 3: Generic Typed Pool
type Pool[T any] struct {
	pool  sync.Pool
	reset func(*T)
}

func NewPool[T any](newFunc func() *T, resetFunc func(*T)) *Pool[T] {
	// TODO: Implement this function.

	// This combines `sync.Pool` with Go generics to create a type-safe pool.

	// Step 1: Initialize the `Pool[T]` struct.
	// - Store the `resetFunc`.
	// - The `sync.Pool.New` function should call the user-provided `newFunc`.
	//   - `New: func() interface{} { return newFunc() }`
	return nil
}

func (p *Pool[T]) Get() *T {
	// TODO: Implement this function.
	// - Get the item from the pool and type-assert it to `*T`.
	return nil
}

func (p *Pool[T]) Put(obj *T) {
	// TODO: Implement this function.
	// - Before putting the object back, check if a `reset` function was provided.
	// - If `p.reset != nil`, call `p.reset(obj)`.
	// - Then, put the object back into the pool.
}

// Exercise 4: Pool with Metrics
type MetricsPool struct {
	pool sync.Pool
	gets atomic.Int64
	puts atomic.Int64
	news atomic.Int64
}

type PoolStats struct {
	Gets    int64
	Puts    int64
	News    int64
	HitRate float64 // Percentage of Gets that were cache hits
}

func NewMetricsPool(newFunc func() interface{}) *MetricsPool {
	// TODO: Implement this function.

	// Step 1: Initialize the `MetricsPool`.
	// Step 2: The `sync.Pool.New` function needs to be special. It must wrap the user-provided `newFunc`.
	// - When `New` is called, it should first increment the `news` counter atomically.
	// - Then, it should call `newFunc()` and return its result.
	// - This is how you track a "miss" - when the pool had to create a new object.
	mp := &MetricsPool{}
	mp.pool.New = func() interface{} {
		mp.news.Add(1)
		return newFunc()
	}
	return mp
}

func (mp *MetricsPool) Get() interface{} {
	// TODO: Implement this function.
	// - First, atomically increment the `gets` counter.
	// - Then, get and return an item from the pool.
	return nil
}

func (mp *MetricsPool) Put(obj interface{}) {
	// TODO: Implement this function.
	// - First, atomically increment the `puts` counter.
	// - Then, put the object back into the pool.
}

func (mp *MetricsPool) Stats() PoolStats {
	// TODO: Implement this function.

	// Step 1: Atomically load the current values of `gets`, `puts`, and `news`.
	// - It's important to load them into local variables so you're working with a consistent snapshot.

	// Step 2: Calculate the hit rate.
	// - A "hit" is a `Get` that *didn't* result in a `New`.
	// - So, `hits = gets - news`.
	// - `hitRate = (float64(hits) / float64(gets)) * 100`.
	// - Be careful to avoid division by zero if `gets` is 0.

	// Step 3: Return the `PoolStats` struct.
	return PoolStats{}
}

// Exercise 5: Size-Classed Buffer Pool
type SizeClassedPool struct {
	pools [4]sync.Pool
}

func NewSizeClassedPool() *SizeClassedPool {
	// TODO: Implement this function.
	// - This pattern uses multiple pools, each for a different "size class". This is more memory-efficient than a single pool if you work with objects of varying sizes.

	// Step 1: Initialize the `SizeClassedPool`.
	// Step 2: For each element in the `pools` array, set its `New` function.
	// - `scp.pools[0].New = func() interface{} { ... }` should create a buffer for the smallest size class (e.g., 1KB).
	// - `scp.pools[1].New = func() interface{} { ... }` should create a buffer for the next size class (e.g., 4KB).
	// - And so on for all your size classes.
	// - Remember to pool pointers to slices: `buf := make([]byte, 0, 1024); return &buf`.
	return nil
}

func (scp *SizeClassedPool) Get(size int) *[]byte {
	// TODO: Implement this function.
	// - Based on the requested `size`, decide which pool to use.
	// - Use a `switch` statement or `if/else if` chain.
	// - `case size <= 1024:` -> get from `scp.pools[0]`.
	// - `case size <= 4096:` -> get from `scp.pools[1]`.
	// - ...
	// - Get from the selected pool and type-assert the result to `*[]byte`.
	return nil
}

func (scp *SizeClassedPool) Put(buf *[]byte) {
	// TODO: Implement this function.
	// - To return a buffer to the correct pool, you need to check its `capacity`.
	// - `capacity := cap(*buf)`
	// - Use a `switch` statement on the `capacity` to determine which pool it belongs to.
	// - Reset the buffer's length (`*buf = (*buf)[:0]`) and `Put` it back into the correct pool.
}

// Exercise 6: Bounded Pool with Semaphore
type BoundedPool struct {
	pool      sync.Pool
	semaphore chan struct{}
	maxSize   int
}

func NewBoundedPool(maxSize int, newFunc func() interface{}) *BoundedPool {
	// TODO: Implement this function.
	// - This pool uses a semaphore to limit the total number of objects that can be "in circulation" (i.e., either in use or in the pool).
	// - Initialize the `BoundedPool` struct.
	// - `semaphore`: This should be a buffered channel of `struct{}` with a capacity of `maxSize`.
	// - `pool`: Initialize this with the provided `newFunc`.
	return nil
}

func (bp *BoundedPool) Get() interface{} {
	// TODO: Implement this function.
	// - To get an object, you must first acquire a permit from the semaphore.
	// - `bp.semaphore <- struct{}{}`
	// - This will block if `maxSize` objects are already in circulation.
	// - Once the permit is acquired, get an object from the underlying `sync.Pool`.
	return nil
}

func (bp *BoundedPool) Put(obj interface{}) {
	// TODO: Implement this function.
	// - First, return the object to the underlying `sync.Pool`.
	// - Then, release the semaphore permit by receiving from the channel.
	// - `<-bp.semaphore`
	// - This will make space in the semaphore, allowing a blocked `Get()` call to proceed.
}

func (bp *BoundedPool) InUse() int {
	// TODO: Implement this function.
	// - The number of objects currently "in use" (i.e., checked out from the pool) is equal to the number of permits taken from the semaphore.
	// - This can be read by checking the length of the semaphore channel.
	// - `return len(bp.semaphore)`
	return 0
}

// Exercise 7: Benchmark Pool Performance
//
// Implement benchmark functions to measure pool performance.
//
// Requirements:
//   - BenchmarkWithoutPool: Allocate new buffers without pooling
//   - BenchmarkWithPool: Reuse buffers from pool
//   - Use b.ReportAllocs() to show allocation statistics
//   - Run with: go test -bench=. -benchmem
//
// This is implemented in exercise_test.go

// Bonus Exercise 8: Worker Pool Pattern
type Worker struct {
	buf  *bytes.Buffer
	temp []byte
}

type WorkerPool struct {
	pool sync.Pool
}

func NewWorkerPool() *WorkerPool {
	// TODO: Implement this function.
	// - Create a pool that allocates `*Worker` objects.
	// - The `New` function should create a `&Worker{...}` with its internal `buf` and `temp` slice initialized.
	return nil
}

func (wp *WorkerPool) Process(data string) string {
	// TODO: Implement this function.

	// This pattern is powerful for high-throughput processing. Instead of allocating buffers for every request, you allocate a whole "worker" object that contains all the necessary buffers.

	// Step 1: Get a worker from the pool.
	// - `worker := wp.pool.Get().(*Worker)`

	// Step 2: Defer putting the worker back.
	// - `defer wp.pool.Put(worker)`
	// - This guarantees the worker is returned to the pool when the function finishes.

	// Step 3: Reset the worker's state.
	// - `worker.Reset()`
	// - This is crucial to clear any data from the previous use.

	// Step 4: Use the worker's resources to do the processing.
	// - `worker.buf.WriteString(...)`
	// - `worker.temp = append(...)`

	// Step 5: Return the result.
	return ""
}

// Helper function to reset worker state
func (w *Worker) Reset() {
	// TODO: Implement this function.
	// - Reset the buffer: `w.buf.Reset()`.
	// - Reset the slice's length: `w.temp = w.temp[:0]`.
}
