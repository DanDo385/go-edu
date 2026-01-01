//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package channelsbasics

import (
	"context"

	"time"
)
// TODO: implement Ping.
func Ping(value int) <-chan int { panic("TODO: implement") }
// TODO: implement PingPong.
func PingPong(n int) (chan<- int, <-chan int) { panic("TODO: implement") }
// TODO: implement Merge.
func Merge(channels ...<-chan int) <-chan int { panic("TODO: implement") }
// TODO: implement Filter.
func Filter(input <-chan int, predicate func(int) bool) <-chan int { panic("TODO: implement") }
// TODO: implement Map.
func Map(input <-chan int, transform func(int) int) <-chan int { panic("TODO: implement") }
// TODO: implement Take.
func Take(input <-chan int, n int) <-chan int { panic("TODO: implement") }
// TODO: implement OrDone.
func OrDone(ctx context.Context, input <-chan int) <-chan int { panic("TODO: implement") }
// TODO: implement Tee.
func Tee(input <-chan int) (<-chan int, <-chan int) { panic("TODO: implement") }
// TODO: implement Bridge.
func Bridge(input <-chan (<-chan int)) <-chan int { panic("TODO: implement") }
// TODO: implement Debounce.
func Debounce(input <-chan int, duration time.Duration) <-chan int { panic("TODO: implement") }
// TODO: implement NewBoundedQueue.
func NewBoundedQueue(capacity int) *BoundedQueue { panic("TODO: implement") }
// TODO: implement Enqueue.
func (q *BoundedQueue) Enqueue(value int) { panic("TODO: implement") }
// TODO: implement Dequeue.
func (q *BoundedQueue) Dequeue() int { panic("TODO: implement") }
// TODO: implement TryEnqueue.
func (q *BoundedQueue) TryEnqueue(value int) bool { panic("TODO: implement") }
// TODO: implement TryDequeue.
func (q *BoundedQueue) TryDequeue() (int, bool) { panic("TODO: implement") }
// TODO: implement NewBroadcaster.
func NewBroadcaster() *Broadcaster { panic("TODO: implement") }
// TODO: implement Subscribe.
func (b *Broadcaster) Subscribe() <-chan Message { panic("TODO: implement") }
// TODO: implement Unsubscribe.
func (b *Broadcaster) Unsubscribe(ch <-chan Message) { panic("TODO: implement") }
// TODO: implement Send.
func (b *Broadcaster) Send(msg Message) { panic("TODO: implement") }
// TODO: implement Close.
func (b *Broadcaster) Close() { panic("TODO: implement") }
// TODO: implement NewBarrier.
func NewBarrier(n int) *Barrier { panic("TODO: implement") }
// TODO: implement Wait.
func (b *Barrier) Wait() { panic("TODO: implement") }
