package exercise

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestParallelSum(t *testing.T) {
	tests := []struct {
		name       string
		n          int
		numWorkers int
		expected   int64
	}{
		{
			name:       "sum 1 to 100 with 4 workers",
			n:          100,
			numWorkers: 4,
			expected:   5050,
		},
		{
			name:       "sum 1 to 1000 with 10 workers",
			n:          1000,
			numWorkers: 10,
			expected:   500500,
		},
		{
			name:       "sum 1 to 10 with 1 worker",
			n:          10,
			numWorkers: 1,
			expected:   55,
		},
		{
			name:       "sum 1 to 100 with 100 workers",
			n:          100,
			numWorkers: 100,
			expected:   5050,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParallelSum(tt.n, tt.numWorkers)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestFanOutFanIn(t *testing.T) {
	t.Run("fan out distributes values", func(t *testing.T) {
		input := make(chan int, 10)
		outputs := FanOut(input, 3)

		if len(outputs) != 3 {
			t.Fatalf("Expected 3 output channels, got %d", len(outputs))
		}

		// Send values
		for i := 0; i < 9; i++ {
			input <- i
		}
		close(input)

		// Collect values from all outputs
		received := make(map[int]bool)
		var mu sync.Mutex

		var wg sync.WaitGroup
		for _, out := range outputs {
			wg.Add(1)
			go func(ch <-chan int) {
				defer wg.Done()
				for v := range ch {
					mu.Lock()
					received[v] = true
					mu.Unlock()
				}
			}(out)
		}
		wg.Wait()

		// Verify all values received
		if len(received) != 9 {
			t.Errorf("Expected 9 unique values, got %d", len(received))
		}
	})

	t.Run("fan in merges channels", func(t *testing.T) {
		ch1 := make(chan int, 5)
		ch2 := make(chan int, 5)
		ch3 := make(chan int, 5)

		// Send values to inputs
		for i := 0; i < 5; i++ {
			ch1 <- i
			ch2 <- i + 10
			ch3 <- i + 20
		}
		close(ch1)
		close(ch2)
		close(ch3)

		output := FanIn(ch1, ch2, ch3)

		// Collect all values
		received := make(map[int]bool)
		for v := range output {
			received[v] = true
		}

		// Verify count
		if len(received) != 15 {
			t.Errorf("Expected 15 unique values, got %d", len(received))
		}
	})
}

func TestWorkerPool(t *testing.T) {
	t.Run("processes jobs concurrently", func(t *testing.T) {
		pool := NewWorkerPool(5)
		if pool == nil {
			t.Fatal("NewWorkerPool returned nil")
		}

		var counter ConcurrentCounter
		var wg sync.WaitGroup

		// Submit 100 jobs
		for i := 0; i < 100; i++ {
			wg.Add(1)
			pool.Submit(func() {
				defer wg.Done()
				counter.Increment()
				time.Sleep(1 * time.Millisecond)
			})
		}

		wg.Wait()
		pool.Stop()

		if counter.Value() != 100 {
			t.Errorf("Expected 100 jobs processed, got %d", counter.Value())
		}
	})

	t.Run("stop waits for jobs to finish", func(t *testing.T) {
		pool := NewWorkerPool(2)
		if pool == nil {
			t.Fatal("NewWorkerPool returned nil")
		}

		var counter ConcurrentCounter

		// Submit slow jobs
		for i := 0; i < 10; i++ {
			pool.Submit(func() {
				time.Sleep(10 * time.Millisecond)
				counter.Increment()
			})
		}

		pool.Stop()

		// All jobs should have completed
		if counter.Value() != 10 {
			t.Errorf("Expected 10 jobs completed, got %d", counter.Value())
		}
	})
}

func TestRateLimiter(t *testing.T) {
	t.Run("limits rate correctly", func(t *testing.T) {
		limiter := NewRateLimiter(10) // 10 ops/sec
		if limiter == nil {
			t.Fatal("NewRateLimiter returned nil")
		}

		start := time.Now()
		count := 0

		// Try to do 25 operations
		// First 10 happen immediately (token bucket pre-filled)
		// Next 15 require waiting: 15 ops at 10 ops/sec = 1.5 seconds
		for i := 0; i < 25; i++ {
			limiter.Wait()
			count++
		}

		elapsed := time.Since(start)

		// Should take at least 1.4 seconds (accounting for timing variations)
		// First 10 ops are immediate, next 15 ops take 1.5s at 10 ops/sec
		if elapsed < 1400*time.Millisecond {
			t.Errorf("Expected at least 1.4s for rate-limited execution, got %v", elapsed)
		}

		if count != 25 {
			t.Errorf("Expected 25 operations, got %d", count)
		}
	})
}

func TestConcurrentCounter(t *testing.T) {
	t.Run("increments and decrements correctly", func(t *testing.T) {
		counter := &ConcurrentCounter{}

		counter.Increment()
		counter.Increment()
		counter.Increment()

		if counter.Value() != 3 {
			t.Errorf("Expected 3, got %d", counter.Value())
		}

		counter.Decrement()

		if counter.Value() != 2 {
			t.Errorf("Expected 2, got %d", counter.Value())
		}
	})

	t.Run("is thread-safe", func(t *testing.T) {
		counter := &ConcurrentCounter{}
		var wg sync.WaitGroup

		// Launch 100 goroutines, each incrementing 100 times
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					counter.Increment()
				}
			}()
		}

		wg.Wait()

		expected := int64(100 * 100)
		if counter.Value() != expected {
			t.Errorf("Expected %d, got %d (race condition detected)", expected, counter.Value())
		}
	})
}

func TestGracefulWorker(t *testing.T) {
	t.Run("stops when context is cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		worker := NewGracefulWorker(ctx)

		worker.Start()

		// Let it run for a bit
		time.Sleep(50 * time.Millisecond)

		// Cancel and give it time to stop
		cancel()
		time.Sleep(10 * time.Millisecond)

		workDone := worker.WorkDone()
		if workDone == 0 {
			t.Error("Expected worker to do some work")
		}

		// Verify it stopped (work shouldn't increase significantly)
		time.Sleep(50 * time.Millisecond)
		workAfter := worker.WorkDone()

		if workAfter > workDone+100 {
			t.Errorf("Worker didn't stop promptly (work increased from %d to %d)", workDone, workAfter)
		}
	})
}

func TestPipeline(t *testing.T) {
	t.Run("processes values through stages", func(t *testing.T) {
		input := make(chan int, 10)

		// Stage 1: double the value
		double := func(in <-chan int) <-chan int {
			out := make(chan int)
			go func() {
				defer close(out)
				for v := range in {
					out <- v * 2
				}
			}()
			return out
		}

		// Stage 2: add 10
		addTen := func(in <-chan int) <-chan int {
			out := make(chan int)
			go func() {
				defer close(out)
				for v := range in {
					out <- v + 10
				}
			}()
			return out
		}

		output := Pipeline(input, double, addTen)

		// Send values
		for i := 1; i <= 5; i++ {
			input <- i
		}
		close(input)

		// Collect results
		results := []int{}
		for v := range output {
			results = append(results, v)
		}

		// Expected: (1*2)+10=12, (2*2)+10=14, (3*2)+10=16, (4*2)+10=18, (5*2)+10=20
		expected := []int{12, 14, 16, 18, 20}

		if len(results) != len(expected) {
			t.Fatalf("Expected %d results, got %d", len(expected), len(results))
		}

		for i, exp := range expected {
			if results[i] != exp {
				t.Errorf("At index %d: expected %d, got %d", i, exp, results[i])
			}
		}
	})
}

func TestBoundedParallel(t *testing.T) {
	t.Run("limits concurrency", func(t *testing.T) {
		const maxConcurrent = 5
		var currentConcurrent int
		var maxObserved int
		var mu sync.Mutex

		var fns []func()
		for i := 0; i < 50; i++ {
			fns = append(fns, func() {
				mu.Lock()
				currentConcurrent++
				if currentConcurrent > maxObserved {
					maxObserved = currentConcurrent
				}
				mu.Unlock()

				time.Sleep(10 * time.Millisecond)

				mu.Lock()
				currentConcurrent--
				mu.Unlock()
			})
		}

		BoundedParallel(maxConcurrent, fns...)

		if maxObserved > maxConcurrent {
			t.Errorf("Expected max concurrency %d, observed %d", maxConcurrent, maxObserved)
		}

		if maxObserved < maxConcurrent {
			t.Errorf("Expected to reach max concurrency %d, only reached %d", maxConcurrent, maxObserved)
		}
	})

	t.Run("executes all functions", func(t *testing.T) {
		var counter ConcurrentCounter

		var fns []func()
		for i := 0; i < 100; i++ {
			fns = append(fns, func() {
				counter.Increment()
			})
		}

		BoundedParallel(10, fns...)

		if counter.Value() != 100 {
			t.Errorf("Expected all 100 functions to execute, got %d", counter.Value())
		}
	})
}
