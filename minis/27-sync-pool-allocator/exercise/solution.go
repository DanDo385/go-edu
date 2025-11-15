//go:build solution
// +build solution

package exercise

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
)

// Exercise 1: Basic Buffer Pool
type BufferPool struct {
	pool sync.Pool
}

func NewBufferPool() *BufferPool {
	return &BufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

func (bp *BufferPool) Get() *bytes.Buffer {
	return bp.pool.Get().(*bytes.Buffer)
}

func (bp *BufferPool) Put(buf *bytes.Buffer) {
	buf.Reset()
	bp.pool.Put(buf)
}

// Exercise 2: Slice Pool
type SlicePool struct {
	pool     sync.Pool
	capacity int
}

func NewSlicePool(capacity int) *SlicePool {
	return &SlicePool{
		capacity: capacity,
		pool: sync.Pool{
			New: func() interface{} {
				slice := make([]byte, 0, capacity)
				return &slice
			},
		},
	}
}

func (sp *SlicePool) Get() *[]byte {
	return sp.pool.Get().(*[]byte)
}

func (sp *SlicePool) Put(slice *[]byte) {
	*slice = (*slice)[:0] // Reset length, keep capacity
	sp.pool.Put(slice)
}

// Exercise 3: Generic Typed Pool
type Pool[T any] struct {
	pool  sync.Pool
	reset func(*T)
}

func NewPool[T any](newFunc func() *T, resetFunc func(*T)) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() interface{} {
				return newFunc()
			},
		},
		reset: resetFunc,
	}
}

func (p *Pool[T]) Get() *T {
	return p.pool.Get().(*T)
}

func (p *Pool[T]) Put(obj *T) {
	if p.reset != nil {
		p.reset(obj)
	}
	p.pool.Put(obj)
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
	HitRate float64
}

func NewMetricsPool(newFunc func() interface{}) *MetricsPool {
	mp := &MetricsPool{}
	mp.pool.New = func() interface{} {
		mp.news.Add(1)
		return newFunc()
	}
	return mp
}

func (mp *MetricsPool) Get() interface{} {
	mp.gets.Add(1)
	return mp.pool.Get()
}

func (mp *MetricsPool) Put(obj interface{}) {
	mp.puts.Add(1)
	mp.pool.Put(obj)
}

func (mp *MetricsPool) Stats() PoolStats {
	gets := mp.gets.Load()
	puts := mp.puts.Load()
	news := mp.news.Load()

	var hitRate float64
	if gets > 0 {
		hitRate = float64(gets-news) / float64(gets) * 100
	}

	return PoolStats{
		Gets:    gets,
		Puts:    puts,
		News:    news,
		HitRate: hitRate,
	}
}

// Exercise 5: Size-Classed Buffer Pool
type SizeClassedPool struct {
	pools [4]sync.Pool
}

func NewSizeClassedPool() *SizeClassedPool {
	scp := &SizeClassedPool{}

	// 1KB pool
	scp.pools[0].New = func() interface{} {
		buf := make([]byte, 0, 1024)
		return &buf
	}

	// 4KB pool
	scp.pools[1].New = func() interface{} {
		buf := make([]byte, 0, 4096)
		return &buf
	}

	// 16KB pool
	scp.pools[2].New = func() interface{} {
		buf := make([]byte, 0, 16384)
		return &buf
	}

	// 64KB pool
	scp.pools[3].New = func() interface{} {
		buf := make([]byte, 0, 65536)
		return &buf
	}

	return scp
}

func (scp *SizeClassedPool) Get(size int) *[]byte {
	var poolIdx int
	switch {
	case size <= 1024:
		poolIdx = 0
	case size <= 4096:
		poolIdx = 1
	case size <= 16384:
		poolIdx = 2
	default:
		poolIdx = 3
	}

	return scp.pools[poolIdx].Get().(*[]byte)
}

func (scp *SizeClassedPool) Put(buf *[]byte) {
	*buf = (*buf)[:0] // Reset length

	capacity := cap(*buf)
	var poolIdx int
	switch {
	case capacity <= 1024:
		poolIdx = 0
	case capacity <= 4096:
		poolIdx = 1
	case capacity <= 16384:
		poolIdx = 2
	default:
		poolIdx = 3
	}

	scp.pools[poolIdx].Put(buf)
}

// Exercise 6: Bounded Pool with Semaphore
type BoundedPool struct {
	pool      sync.Pool
	semaphore chan struct{}
	maxSize   int
}

func NewBoundedPool(maxSize int, newFunc func() interface{}) *BoundedPool {
	return &BoundedPool{
		pool: sync.Pool{
			New: newFunc,
		},
		semaphore: make(chan struct{}, maxSize),
		maxSize:   maxSize,
	}
}

func (bp *BoundedPool) Get() interface{} {
	bp.semaphore <- struct{}{} // Acquire semaphore (blocks if full)
	return bp.pool.Get()
}

func (bp *BoundedPool) Put(obj interface{}) {
	bp.pool.Put(obj)
	<-bp.semaphore // Release semaphore
}

func (bp *BoundedPool) InUse() int {
	return len(bp.semaphore)
}

// Exercise 8: Worker Pool Pattern
type Worker struct {
	buf  *bytes.Buffer
	temp []byte
}

type WorkerPool struct {
	pool sync.Pool
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &Worker{
					buf:  new(bytes.Buffer),
					temp: make([]byte, 0, 1024),
				}
			},
		},
	}
}

func (wp *WorkerPool) Process(data string) string {
	worker := wp.pool.Get().(*Worker)
	defer wp.pool.Put(worker)

	// Reset worker state
	worker.Reset()

	// Process data using worker's resources
	worker.buf.WriteString("Processed: ")
	worker.buf.WriteString(data)

	// Use temp buffer for additional processing
	worker.temp = append(worker.temp, []byte(" [transformed]")...)

	result := worker.buf.String() + string(worker.temp)
	return result
}

func (w *Worker) Reset() {
	if w.buf != nil {
		w.buf.Reset()
	}
	if w.temp != nil {
		w.temp = w.temp[:0]
	}
}

// Bonus: Enhanced MetricsPool with detailed stats
type EnhancedMetricsPool struct {
	pool         sync.Pool
	gets         atomic.Int64
	puts         atomic.Int64
	news         atomic.Int64
	reuses       atomic.Int64
}

type EnhancedStats struct {
	Gets        int64
	Puts        int64
	News        int64
	Reuses      int64
	HitRate     float64
	MissRate    float64
	Efficiency  float64
}

func NewEnhancedMetricsPool(newFunc func() interface{}) *EnhancedMetricsPool {
	emp := &EnhancedMetricsPool{}
	emp.pool.New = func() interface{} {
		emp.news.Add(1)
		return newFunc()
	}
	return emp
}

func (emp *EnhancedMetricsPool) Get() interface{} {
	emp.gets.Add(1)
	obj := emp.pool.Get()

	// Try to determine if this was a reuse
	// (In real implementation, we'd check if New was called)
	currentNews := emp.news.Load()
	currentGets := emp.gets.Load()
	if currentGets > currentNews {
		emp.reuses.Add(1)
	}

	return obj
}

func (emp *EnhancedMetricsPool) Put(obj interface{}) {
	emp.puts.Add(1)
	emp.pool.Put(obj)
}

func (emp *EnhancedMetricsPool) Stats() EnhancedStats {
	gets := emp.gets.Load()
	puts := emp.puts.Load()
	news := emp.news.Load()
	reuses := emp.reuses.Load()

	var hitRate, missRate, efficiency float64
	if gets > 0 {
		hitRate = float64(gets-news) / float64(gets) * 100
		missRate = float64(news) / float64(gets) * 100
	}

	if puts > 0 {
		efficiency = float64(reuses) / float64(puts) * 100
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

func (es EnhancedStats) String() string {
	return fmt.Sprintf(
		"Gets: %d, Puts: %d, News: %d, Reuses: %d, Hit Rate: %.1f%%, Miss Rate: %.1f%%, Efficiency: %.1f%%",
		es.Gets, es.Puts, es.News, es.Reuses, es.HitRate, es.MissRate, es.Efficiency,
	)
}

// Bonus: Copy-on-Write Pool Pattern
type CopyOnWritePool[T any] struct {
	pool sync.Pool
	copy func(*T) *T
}

func NewCopyOnWritePool[T any](newFunc func() *T, copyFunc func(*T) *T) *CopyOnWritePool[T] {
	return &CopyOnWritePool[T]{
		pool: sync.Pool{
			New: func() interface{} {
				return newFunc()
			},
		},
		copy: copyFunc,
	}
}

func (cp *CopyOnWritePool[T]) Get() *T {
	obj := cp.pool.Get().(*T)
	// Return a copy so original stays in pool
	if cp.copy != nil {
		return cp.copy(obj)
	}
	return obj
}

func (cp *CopyOnWritePool[T]) Put(obj *T) {
	cp.pool.Put(obj)
}
