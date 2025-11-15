//go:build solution
// +build solution

package exercise

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// ParallelSum calculates the sum using multiple workers.
func ParallelSum(n int, numWorkers int) int64 {
	var total atomic.Int64
	var wg sync.WaitGroup

	// Calculate range for each worker
	rangeSize := n / numWorkers
	remainder := n % numWorkers

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		// Calculate start and end for this worker
		start := i*rangeSize + 1
		end := (i + 1) * rangeSize

		// Give remainder to last worker
		if i == numWorkers-1 {
			end += remainder
		}

		go func(s, e int) {
			defer wg.Done()

			// Sum this worker's range
			var sum int64
			for j := s; j <= e; j++ {
				sum += int64(j)
			}

			// Add to total atomically
			total.Add(sum)
		}(start, end)
	}

	wg.Wait()
	return total.Load()
}

// FanOut distributes values to multiple channels.
func FanOut(input <-chan int, numWorkers int) []<-chan int {
	outputs := make([]chan int, numWorkers)
	readOnlyOutputs := make([]<-chan int, numWorkers)

	for i := 0; i < numWorkers; i++ {
		outputs[i] = make(chan int)
		readOnlyOutputs[i] = outputs[i]
	}

	// Distributor goroutine
	go func() {
		defer func() {
			for _, ch := range outputs {
				close(ch)
			}
		}()

		i := 0
		for v := range input {
			// Round-robin distribution
			outputs[i%numWorkers] <- v
			i++
		}
	}()

	return readOnlyOutputs
}

// FanIn merges multiple channels into one.
func FanIn(inputs ...<-chan int) <-chan int {
	output := make(chan int)
	var wg sync.WaitGroup

	// Launch a goroutine for each input
	for _, in := range inputs {
		wg.Add(1)
		go func(ch <-chan int) {
			defer wg.Done()
			for v := range ch {
				output <- v
			}
		}(in)
	}

	// Close output when all inputs are done
	go func() {
		wg.Wait()
		close(output)
	}()

	return output
}

// NewWorkerPool creates a worker pool.
func NewWorkerPool(numWorkers int) *WorkerPool {
	pool := &WorkerPool{
		jobs: make(chan func(), 100), // Buffered channel
	}

	// Launch workers
	for i := 0; i < numWorkers; i++ {
		pool.wg.Add(1)
		go func() {
			defer pool.wg.Done()
			for job := range pool.jobs {
				job()
			}
		}()
	}

	return pool
}

// Submit adds a job to the pool.
func (p *WorkerPool) Submit(job func()) {
	if p.stopped.Load() {
		return // Pool is stopped, don't accept new jobs
	}
	p.jobs <- job
}

// Stop shuts down the pool.
func (p *WorkerPool) Stop() {
	p.stopped.Store(true)
	close(p.jobs) // Signal workers to stop
	p.wg.Wait()   // Wait for all workers to finish
}

// NewRateLimiter creates a rate limiter.
func NewRateLimiter(maxOps int) *RateLimiter {
	limiter := &RateLimiter{
		ticker: time.NewTicker(time.Second / time.Duration(maxOps)),
		tokens: make(chan struct{}, maxOps),
	}

	// Pre-fill tokens
	for i := 0; i < maxOps; i++ {
		limiter.tokens <- struct{}{}
	}

	// Refill tokens at rate
	go func() {
		for range limiter.ticker.C {
			select {
			case limiter.tokens <- struct{}{}:
			default:
				// Token bucket full, skip
			}
		}
	}()

	return limiter
}

// Wait blocks until a token is available.
func (r *RateLimiter) Wait() {
	<-r.tokens
}

// ConcurrentCounter implementation.

// Increment atomically increments.
func (c *ConcurrentCounter) Increment() {
	c.value.Add(1)
}

// Decrement atomically decrements.
func (c *ConcurrentCounter) Decrement() {
	c.value.Add(-1)
}

// Value returns current value.
func (c *ConcurrentCounter) Value() int64 {
	return c.value.Load()
}

// NewGracefulWorker creates a new graceful worker.
func NewGracefulWorker(ctx context.Context) *GracefulWorker {
	return &GracefulWorker{
		ctx:      ctx,
		workDone: atomic.Int64{},
	}
}

// Start begins execution.
func (w *GracefulWorker) Start() {
	go func() {
		for {
			select {
			case <-w.ctx.Done():
				return // Exit when cancelled
			default:
				// Do work
				w.workDone.Add(1)
				time.Sleep(1 * time.Microsecond)
			}
		}
	}()
}

// WorkDone returns the total work completed.
func (w *GracefulWorker) WorkDone() int64 {
	return w.workDone.Load()
}

// Pipeline chains stages together.
func Pipeline(input <-chan int, stages ...func(<-chan int) <-chan int) <-chan int {
	output := input
	for _, stage := range stages {
		output = stage(output)
	}
	return output
}

// BoundedParallel executes with limited concurrency.
func BoundedParallel(maxConcurrent int, fns ...func()) {
	sem := make(chan struct{}, maxConcurrent) // Semaphore
	var wg sync.WaitGroup

	for _, fn := range fns {
		wg.Add(1)

		go func(f func()) {
			defer wg.Done()

			// Acquire semaphore
			sem <- struct{}{}
			defer func() { <-sem }() // Release semaphore

			// Execute function
			f()
		}(fn)
	}

	wg.Wait()
}
