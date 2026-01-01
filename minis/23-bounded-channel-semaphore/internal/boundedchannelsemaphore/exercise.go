//go:build !solution && !reference

package boundedchannelsemaphore

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ============================================================================
// Exercise 1-3: Counting Semaphore
// ============================================================================

type semaphore struct {
	sem chan struct{}
}

func NewSemaphore(maxPermits int) SemaphoreInterface {
	if maxPermits <= 0 {
		maxPermits = 1
	}
	return &semaphore{sem: make(chan struct{}, maxPermits)}
}

func (s *semaphore) Acquire() {
	s.sem <- struct{}{}
}

func (s *semaphore) Release() {
	<-s.sem
}

func (s *semaphore) TryAcquire() bool {
	select {
	case s.sem <- struct{}{}:
		return true
	default:
		return false
	}
}

func (s *semaphore) AcquireWithContext(ctx context.Context) error {
	select {
	case s.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// ============================================================================
// Exercise 4: Rate Limiter (token bucket)
// ============================================================================

type rateLimiter struct {
	tokens chan struct{}
	rate   time.Duration
	done   chan struct{}
}

func NewRateLimiter(maxBurst int, rate time.Duration) RateLimiterInterface {
	if maxBurst <= 0 {
		maxBurst = 1
	}
	if rate <= 0 {
		rate = time.Second
	}

	rl := &rateLimiter{
		tokens: make(chan struct{}, maxBurst),
		rate:   rate,
		done:   make(chan struct{}),
	}

	// Initial burst.
	for i := 0; i < maxBurst; i++ {
		rl.tokens <- struct{}{}
	}

	go rl.refill()
	return rl
}

func (rl *rateLimiter) refill() {
	t := time.NewTicker(rl.rate)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			// Non-blocking add.
			select {
			case rl.tokens <- struct{}{}:
			default:
			}
		case <-rl.done:
			return
		}
	}
}

func (rl *rateLimiter) Wait() {
	<-rl.tokens
}

func (rl *rateLimiter) TryAcquire() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

func (rl *rateLimiter) Stop() {
	select {
	case <-rl.done:
		// already closed
	default:
		close(rl.done)
	}
}

// ============================================================================
// Exercise 5: Weighted Semaphore
// ============================================================================

type WeightedSemaphore struct {
	permits chan struct{}
}

func NewWeightedSemaphore(maxWeight int) *WeightedSemaphore {
	if maxWeight <= 0 {
		maxWeight = 1
	}
	return &WeightedSemaphore{permits: make(chan struct{}, maxWeight)}
}

func (ws *WeightedSemaphore) Acquire(weight int) {
	for i := 0; i < weight; i++ {
		ws.permits <- struct{}{}
	}
}

func (ws *WeightedSemaphore) Release(weight int) {
	for i := 0; i < weight; i++ {
		<-ws.permits
	}
}

func (ws *WeightedSemaphore) AcquireWithContext(ctx context.Context, weight int) error {
	acquired := 0
	for i := 0; i < weight; i++ {
		select {
		case ws.permits <- struct{}{}:
			acquired++
		case <-ctx.Done():
			// Clean up partial acquisition to avoid leaking permits.
			for j := 0; j < acquired; j++ {
				<-ws.permits
			}
			return ctx.Err()
		}
	}
	return nil
}

// ============================================================================
// Exercise 6: Worker Pool
// ============================================================================

func DefaultProcessor(job Job) Result {
	return Result{JobID: job.ID, Output: fmt.Sprintf("processed:%d", job.ID)}
}

type workerPool struct {
	jobs       chan Job
	results    chan Result
	numWorkers int
	processor  func(Job) Result

	startOnce sync.Once
	stopOnce  sync.Once
	wg        sync.WaitGroup
}

func NewWorkerPool(numWorkers int, processor func(Job) Result) WorkerPoolInterface {
	if numWorkers <= 0 {
		numWorkers = 1
	}
	if processor == nil {
		processor = DefaultProcessor
	}

	return &workerPool{
		jobs:       make(chan Job, numWorkers*2),
		results:    make(chan Result, numWorkers*2),
		numWorkers: numWorkers,
		processor:  processor,
	}
}

func (wp *workerPool) Start() {
	wp.startOnce.Do(func() {
		for i := 0; i < wp.numWorkers; i++ {
			wp.wg.Add(1)
			go func() {
				defer wp.wg.Done()
				for job := range wp.jobs {
					wp.results <- wp.processor(job)
				}
			}()
		}
	})
}

func (wp *workerPool) Submit(job Job) {
	wp.jobs <- job
}

func (wp *workerPool) Results() <-chan Result {
	return wp.results
}

func (wp *workerPool) Stop() {
	wp.stopOnce.Do(func() {
		close(wp.jobs)
		wp.wg.Wait()
		close(wp.results)
	})
}
