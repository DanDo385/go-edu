//go:build !solution && !reference

package selectfaninfanout

import (
	"sync"
	"time"
)

type RateLimiter struct {
	tokens chan struct{}
	rate   int
}

// SelectFirst implements the exercise.
//
// TODO: Implement this function
func SelectFirst(ch1 <-chan string, ch2 <-chan string, timeout time.Duration) (string, bool) {
	// TODO: Implement
	return "", false
}

// NonBlockingSend implements the exercise.
//
// TODO: Implement this function
func NonBlockingSend(ch chan<- int, value int) bool {
	// TODO: Implement
	return false
}

// FanIn implements the exercise.
//
// TODO: Implement this function
func FanIn(channels ...<-chan int) <-chan int {
	// TODO: Implement
	return 0
}

// FanOut implements the exercise.
//
// TODO: Implement this function
func FanOut(input <-chan int, numWorkers int, process func(int) int) <-chan int {
	// TODO: Implement
	return 0
}

// OrChannel implements the exercise.
//
// TODO: Implement this function
func OrChannel(channels ...<-chan struct{}) <-chan struct{} {
	// TODO: Implement
	return nil
}

// TryReceiveAll implements the exercise.
//
// TODO: Implement this function
func TryReceiveAll(channels []<-chan int) map[int]int {
	// TODO: Implement
	return nil
}

// NewRateLimiter implements the exercise.
//
// TODO: Implement this function
func NewRateLimiter(rate int) *RateLimiter {
	// TODO: Implement
	return nil
}

// refill implements the exercise.
//
// TODO: Implement this function
func (rl *RateLimiter) refill() {
	// TODO: Implement
}

// Wait implements the exercise.
//
// TODO: Implement this function
func (rl *RateLimiter) Wait() {
	// TODO: Implement
}

// TryWait implements the exercise.
//
// TODO: Implement this function
func (rl *RateLimiter) TryWait() bool {
	// TODO: Implement
	return false
}

// Pipeline implements the exercise.
//
// TODO: Implement this function
func Pipeline(n int, numWorkers int) <-chan int {
	// TODO: Implement
	return 0
}

// SelectWithPriority implements the exercise.
//
// TODO: Implement this function
func SelectWithPriority(high <-chan int, low <-chan int) (int, bool) {
	// TODO: Implement
	return 0, false
}

// Timeout implements the exercise.
//
// TODO: Implement this function
func Timeout(ch <-chan int, timeout time.Duration) (int, bool) {
	// TODO: Implement
	return 0, false
}
