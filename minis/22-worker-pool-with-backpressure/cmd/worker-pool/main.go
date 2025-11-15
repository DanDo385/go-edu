package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== Worker Pool with Backpressure Demonstrations ===\n")

	demo1_blockingBackpressure()
	demo2_dropStrategy()
	demo3_rejectStrategy()
	demo4_timeoutStrategy()
	demo5_rateLimiting()
	demo6_adaptivePooling()

	fmt.Println("\n=== All demonstrations complete ===")
}

// Demo 1: Blocking backpressure (default channel behavior)
func demo1_blockingBackpressure() {
	fmt.Println("--- Demo 1: Blocking Backpressure ---")
	fmt.Println("Producer sends fast, consumer processes slow")
	fmt.Println("Bounded channel causes producer to block\n")

	jobs := make(chan int, 3) // Small buffer (3 items)
	done := make(chan bool)

	// Fast producer
	go func() {
		for i := 1; i <= 10; i++ {
			fmt.Printf("  [Producer] Sending job %d...\n", i)
			start := time.Now()
			jobs <- i // Blocks when buffer is full
			elapsed := time.Since(start)
			if elapsed > 10*time.Millisecond {
				fmt.Printf("  [Producer] Job %d BLOCKED for %v\n", i, elapsed.Round(time.Millisecond))
			} else {
				fmt.Printf("  [Producer] Job %d sent immediately\n", i)
			}
		}
		close(jobs)
		done <- true
	}()

	// Slow consumer
	go func() {
		for job := range jobs {
			fmt.Printf("  [Consumer] Processing job %d...\n", job)
			time.Sleep(200 * time.Millisecond) // Simulate slow work
			fmt.Printf("  [Consumer] Completed job %d\n", job)
		}
	}()

	<-done
	time.Sleep(100 * time.Millisecond) // Let consumer finish

	fmt.Println("✓ Result: Producer automatically slowed to match consumer speed")
	fmt.Println("✓ No data lost, but producer experienced blocking\n")
}

// Demo 2: Drop strategy (non-blocking, lossy)
func demo2_dropStrategy() {
	fmt.Println("--- Demo 2: Drop Strategy ---")
	fmt.Println("When queue is full, drop new items (good for metrics/logs)\n")

	jobs := make(chan int, 5)
	var wg sync.WaitGroup

	// Statistics
	var sent, dropped int
	var mu sync.Mutex

	// Workers (slow consumers)
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobs {
				fmt.Printf("  [Worker %d] Processing job %d\n", workerID, job)
				time.Sleep(100 * time.Millisecond)
			}
		}(i)
	}

	// Fast producer with drop logic
	go func() {
		for i := 1; i <= 20; i++ {
			select {
			case jobs <- i:
				mu.Lock()
				sent++
				mu.Unlock()
				fmt.Printf("  [Producer] Sent job %d\n", i)
			default:
				mu.Lock()
				dropped++
				mu.Unlock()
				fmt.Printf("  [Producer] ⚠️  DROPPED job %d (queue full)\n", i)
			}
			time.Sleep(20 * time.Millisecond) // Fast production
		}
		close(jobs)
	}()

	wg.Wait()

	fmt.Printf("\n✓ Result: Sent %d jobs, Dropped %d jobs\n", sent, dropped)
	fmt.Println("✓ Producer never blocked, but some data was lost\n")
}

// Demo 3: Reject strategy (return error to caller)
func demo3_rejectStrategy() {
	fmt.Println("--- Demo 3: Reject Strategy ---")
	fmt.Println("When queue is full, reject request with error (good for HTTP APIs)\n")

	type Result struct {
		JobID int
		Err   error
	}

	jobs := make(chan int, 3)
	results := make(chan Result, 20)

	// Worker
	go func() {
		for job := range jobs {
			fmt.Printf("  [Worker] Processing job %d\n", job)
			time.Sleep(150 * time.Millisecond)
		}
	}()

	// Submit function with reject logic
	submit := func(jobID int) error {
		select {
		case jobs <- jobID:
			return nil
		default:
			return fmt.Errorf("queue full")
		}
	}

	// Producer attempts
	for i := 1; i <= 15; i++ {
		err := submit(i)
		if err != nil {
			fmt.Printf("  [Producer] ❌ Job %d REJECTED: %v\n", i, err)
			results <- Result{JobID: i, Err: err}
		} else {
			fmt.Printf("  [Producer] ✓ Job %d accepted\n", i)
			results <- Result{JobID: i, Err: nil}
		}
		time.Sleep(50 * time.Millisecond)
	}

	close(jobs)
	close(results)

	// Count results
	var accepted, rejected int
	for r := range results {
		if r.Err != nil {
			rejected++
		} else {
			accepted++
		}
	}

	fmt.Printf("\n✓ Result: Accepted %d jobs, Rejected %d jobs\n", accepted, rejected)
	fmt.Println("✓ Caller can retry rejected jobs or handle gracefully\n")
}

// Demo 4: Timeout strategy (wait limited time)
func demo4_timeoutStrategy() {
	fmt.Println("--- Demo 4: Timeout Strategy ---")
	fmt.Println("Try to send for up to 100ms, then give up\n")

	jobs := make(chan int, 2)

	// Worker (very slow)
	go func() {
		for job := range jobs {
			fmt.Printf("  [Worker] Processing job %d\n", job)
			time.Sleep(300 * time.Millisecond)
		}
	}()

	// Submit with timeout
	submitWithTimeout := func(jobID int, timeout time.Duration) error {
		select {
		case jobs <- jobID:
			return nil
		case <-time.After(timeout):
			return fmt.Errorf("timeout after %v", timeout)
		}
	}

	// Producer attempts
	var succeeded, timedOut int
	for i := 1; i <= 10; i++ {
		err := submitWithTimeout(i, 100*time.Millisecond)
		if err != nil {
			fmt.Printf("  [Producer] ⏱️  Job %d TIMEOUT\n", i)
			timedOut++
		} else {
			fmt.Printf("  [Producer] ✓ Job %d sent\n", i)
			succeeded++
		}
		time.Sleep(50 * time.Millisecond)
	}

	close(jobs)
	time.Sleep(200 * time.Millisecond)

	fmt.Printf("\n✓ Result: Succeeded %d, Timed out %d\n", succeeded, timedOut)
	fmt.Println("✓ Balance between blocking and dropping\n")
}

// Demo 5: Rate limiting with token bucket
func demo5_rateLimiting() {
	fmt.Println("--- Demo 5: Rate Limiting (Token Bucket) ---")
	fmt.Println("Allow max 5 requests per second\n")

	type RateLimiter struct {
		tokens   chan struct{}
		rate     time.Duration
		capacity int
	}

	newRateLimiter := func(rps int) *RateLimiter {
		rl := &RateLimiter{
			tokens:   make(chan struct{}, rps),
			rate:     time.Second / time.Duration(rps),
			capacity: rps,
		}

		// Fill bucket initially
		for i := 0; i < rps; i++ {
			rl.tokens <- struct{}{}
		}

		// Refill tokens
		go func() {
			ticker := time.NewTicker(rl.rate)
			defer ticker.Stop()
			for range ticker.C {
				select {
				case rl.tokens <- struct{}{}:
				default:
				}
			}
		}()

		return rl
	}

	limiter := newRateLimiter(5) // 5 requests per second

	// Attempt to send 15 requests quickly
	start := time.Now()
	for i := 1; i <= 15; i++ {
		requestStart := time.Now()
		<-limiter.tokens // Wait for token
		waited := time.Since(requestStart)

		fmt.Printf("  [Client] Request %d sent (waited %v)\n", i, waited.Round(time.Millisecond))
	}
	elapsed := time.Since(start)

	fmt.Printf("\n✓ Result: 15 requests in %v\n", elapsed.Round(time.Millisecond))
	fmt.Printf("✓ Rate: %.1f requests/second\n", 15/elapsed.Seconds())
	fmt.Println("✓ Rate limiter prevented overwhelming the system\n")
}

// Demo 6: Adaptive worker pool (scales based on load)
func demo6_adaptivePooling() {
	fmt.Println("--- Demo 6: Adaptive Worker Pool ---")
	fmt.Println("Add workers when queue fills up, remove when idle\n")

	type AdaptivePool struct {
		jobs       chan int
		workers    int
		maxWorkers int
		mu         sync.Mutex
		ctx        context.Context
		cancel     context.CancelFunc
	}

	newAdaptivePool := func(initial, max, queueSize int) *AdaptivePool {
		ctx, cancel := context.WithCancel(context.Background())
		p := &AdaptivePool{
			jobs:       make(chan int, queueSize),
			workers:    initial,
			maxWorkers: max,
			ctx:        ctx,
			cancel:     cancel,
		}

		// Define worker function as closure
		var startWorker func(int)
		startWorker = func(id int) {
			fmt.Printf("  [Pool] Worker %d started\n", id)
			for {
				select {
				case <-p.ctx.Done():
					fmt.Printf("  [Pool] Worker %d stopped\n", id)
					return
				case job, ok := <-p.jobs:
					if !ok {
						return
					}
					// Simulate work
					time.Sleep(time.Duration(50+rand.Intn(50)) * time.Millisecond)
					_ = job
				}
			}
		}

		// Define monitor function as closure
		monitor := func() {
			ticker := time.NewTicker(200 * time.Millisecond)
			defer ticker.Stop()

			for {
				select {
				case <-p.ctx.Done():
					return
				case <-ticker.C:
					queueLen := len(p.jobs)
					queueCap := cap(p.jobs)
					utilization := float64(queueLen) / float64(queueCap)

					p.mu.Lock()
					currentWorkers := p.workers
					p.mu.Unlock()

					if utilization > 0.7 && currentWorkers < p.maxWorkers {
						p.mu.Lock()
						p.workers++
						newWorkerID := p.workers
						p.mu.Unlock()
						fmt.Printf("  [Pool] ⬆️  Queue %d/%d (%.0f%%) - Adding worker %d\n",
							queueLen, queueCap, utilization*100, newWorkerID)
						go startWorker(newWorkerID)
					}
				}
			}
		}

		// Start initial workers
		for i := 0; i < initial; i++ {
			go startWorker(i)
		}

		// Monitor and adjust
		go monitor()

		return p
	}

	pool := newAdaptivePool(2, 6, 10)

	// Send burst of jobs
	go func() {
		for i := 1; i <= 30; i++ {
			pool.jobs <- i
			if i%10 == 0 {
				time.Sleep(500 * time.Millisecond) // Burst pattern
			} else {
				time.Sleep(20 * time.Millisecond)
			}
		}
		close(pool.jobs)
	}()

	time.Sleep(3 * time.Second)
	pool.cancel()

	pool.mu.Lock()
	finalWorkers := pool.workers
	pool.mu.Unlock()

	fmt.Printf("\n✓ Result: Started with 2 workers, scaled to %d workers\n", finalWorkers)
	fmt.Println("✓ Adaptive scaling handles variable load\n")
}

func init() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)
}
