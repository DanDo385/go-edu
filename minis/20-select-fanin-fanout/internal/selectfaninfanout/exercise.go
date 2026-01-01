//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package selectfaninfanout

import "time"
// TODO: implement SelectFirst.
func SelectFirst(ch1, ch2 <-chan string, timeout time.Duration) (string, bool) {
	panic("TODO: implement")
}
// TODO: implement NonBlockingSend.
func NonBlockingSend(ch chan<- int, value int) bool { panic("TODO: implement") }
// TODO: implement FanIn.
func FanIn(channels ...<-chan int) <-chan int { panic("TODO: implement") }
// TODO: implement FanOut.
func FanOut(input <-chan int, numWorkers int, process func(int) int) <-chan int {
	panic("TODO: implement")
}
// TODO: implement OrChannel.
func OrChannel(channels ...<-chan struct{}) <-chan struct{} { panic("TODO: implement") }
// TODO: implement TryReceiveAll.
func TryReceiveAll(channels []<-chan int) map[int]int { panic("TODO: implement") }

type RateLimiter struct {
	tokens chan struct{}
	rate   int
}
// TODO: implement NewRateLimiter.
func NewRateLimiter(rate int) *RateLimiter { panic("TODO: implement") }
// TODO: implement refill.
func (rl *RateLimiter) refill() { panic("TODO: implement") }
// TODO: implement Wait.
func (rl *RateLimiter) Wait() { panic("TODO: implement") }
// TODO: implement TryWait.
func (rl *RateLimiter) TryWait() bool { panic("TODO: implement") }
// TODO: implement Pipeline.
func Pipeline(n int, numWorkers int) <-chan int { panic("TODO: implement") }
// TODO: implement generate.
func generate(n int) <-chan int { panic("TODO: implement") }
// TODO: implement filter.
func filter(in <-chan int, predicate func(int) bool) <-chan int { panic("TODO: implement") }
// TODO: implement SelectWithPriority.
func SelectWithPriority(high, low <-chan int) (int, bool) { panic("TODO: implement") }
// TODO: implement Timeout.
func Timeout(ch <-chan int, timeout time.Duration) (int, bool) { panic("TODO: implement") }
