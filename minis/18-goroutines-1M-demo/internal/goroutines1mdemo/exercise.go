//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package goroutines1mdemo

import "context"
// TODO: implement ParallelSum.
func ParallelSum(n int, numWorkers int) int64 { panic("TODO: implement") }
// TODO: implement FanOut.
func FanOut(input <-chan int, numWorkers int) []<-chan int { panic("TODO: implement") }
// TODO: implement FanIn.
func FanIn(inputs ...<-chan int) <-chan int { panic("TODO: implement") }
// TODO: implement NewWorkerPool.
func NewWorkerPool(numWorkers int) *WorkerPool { panic("TODO: implement") }
// TODO: implement Submit.
func (p *WorkerPool) Submit(job func()) { panic("TODO: implement") }
// TODO: implement Stop.
func (p *WorkerPool) Stop() { panic("TODO: implement") }
// TODO: implement NewRateLimiter.
func NewRateLimiter(maxOps int) *RateLimiter { panic("TODO: implement") }
// TODO: implement Wait.
func (r *RateLimiter) Wait() { panic("TODO: implement") }
// TODO: implement Increment.
func (c *ConcurrentCounter) Increment() { panic("TODO: implement") }
// TODO: implement Decrement.
func (c *ConcurrentCounter) Decrement() { panic("TODO: implement") }
// TODO: implement Value.
func (c *ConcurrentCounter) Value() int64 { panic("TODO: implement") }
// TODO: implement NewGracefulWorker.
func NewGracefulWorker(ctx context.Context) *GracefulWorker { panic("TODO: implement") }
// TODO: implement Start.
func (w *GracefulWorker) Start() { panic("TODO: implement") }
// TODO: implement WorkDone.
func (w *GracefulWorker) WorkDone() int64 { panic("TODO: implement") }
// TODO: implement Pipeline.
func Pipeline(input <-chan int, stages ...func(<-chan int) <-chan int) <-chan int {
	panic("TODO: implement")
}
// TODO: implement BoundedParallel.
func BoundedParallel(maxConcurrent int, fns ...func()) { panic("TODO: implement") }
