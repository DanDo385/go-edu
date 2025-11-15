package exercise

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ============================================================================
// TEST 1: Atomic Counter
// ============================================================================

func TestAtomicCounter_Basic(t *testing.T) {
	counter := NewAtomicCounter()
	if counter == nil {
		t.Fatal("NewAtomicCounter returned nil")
	}

	// Test initial value
	if val := counter.Value(); val != 0 {
		t.Errorf("Initial value = %d, want 0", val)
	}

	// Test increment
	counter.Increment()
	if val := counter.Value(); val != 1 {
		t.Errorf("After Increment: value = %d, want 1", val)
	}

	// Test add
	counter.Add(5)
	if val := counter.Value(); val != 6 {
		t.Errorf("After Add(5): value = %d, want 6", val)
	}

	// Test decrement
	counter.Decrement()
	if val := counter.Value(); val != 5 {
		t.Errorf("After Decrement: value = %d, want 5", val)
	}

	// Test reset
	old := counter.Reset()
	if old != 5 {
		t.Errorf("Reset returned %d, want 5", old)
	}
	if val := counter.Value(); val != 0 {
		t.Errorf("After Reset: value = %d, want 0", val)
	}
}

func TestAtomicCounter_Concurrent(t *testing.T) {
	counter := NewAtomicCounter()
	if counter == nil {
		t.Fatal("NewAtomicCounter returned nil")
	}

	const numGoroutines = 100
	const incrementsPerGoroutine = 1000

	var wg sync.WaitGroup

	// Concurrent increments
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				counter.Increment()
			}
		}()
	}

	wg.Wait()

	expected := int64(numGoroutines * incrementsPerGoroutine)
	if val := counter.Value(); val != expected {
		t.Errorf("Concurrent increments: value = %d, want %d", val, expected)
	}
}

// ============================================================================
// TEST 2: Atomic Flag
// ============================================================================

func TestAtomicFlag_Basic(t *testing.T) {
	flag := NewAtomicFlag()
	if flag == nil {
		t.Fatal("NewAtomicFlag returned nil")
	}

	// Test initial state
	if flag.IsSet() {
		t.Error("Initial state should be false")
	}

	// Test set
	flag.Set()
	if !flag.IsSet() {
		t.Error("After Set: flag should be true")
	}

	// Test clear
	flag.Clear()
	if flag.IsSet() {
		t.Error("After Clear: flag should be false")
	}

	// Test TestAndSet
	old := flag.TestAndSet()
	if old {
		t.Error("TestAndSet returned true, want false (flag was clear)")
	}
	if !flag.IsSet() {
		t.Error("After TestAndSet: flag should be true")
	}

	// TestAndSet again
	old = flag.TestAndSet()
	if !old {
		t.Error("TestAndSet returned false, want true (flag was set)")
	}
}

func TestAtomicFlag_Concurrent(t *testing.T) {
	flag := NewAtomicFlag()
	if flag == nil {
		t.Fatal("NewAtomicFlag returned nil")
	}

	const numGoroutines = 100
	winners := int32(0)

	var wg sync.WaitGroup

	// Multiple goroutines try to set the flag
	// Only one should see the old value as false
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if !flag.TestAndSet() {
				// We were the first to set it
				atomic.AddInt32(&winners, 1)
			}
		}()
	}

	wg.Wait()

	if winners != 1 {
		t.Errorf("TestAndSet winners = %d, want 1 (only one goroutine should win)", winners)
	}
}

// ============================================================================
// TEST 3: Rate Limiter
// ============================================================================

func TestRateLimiter_Basic(t *testing.T) {
	limiter := NewRateLimiter(5, 10)
	if limiter == nil {
		t.Fatal("NewRateLimiter returned nil")
	}

	// Should allow 5 requests immediately (bucket starts full)
	for i := 0; i < 5; i++ {
		if !limiter.Allow() {
			t.Errorf("Request %d rejected, should be allowed", i+1)
		}
	}

	// 6th request should be rejected (bucket empty)
	if limiter.Allow() {
		t.Error("Request 6 allowed, should be rejected (bucket empty)")
	}
}

func TestRateLimiter_Refill(t *testing.T) {
	limiter := NewRateLimiter(2, 5) // 2 tokens max, refill 5/sec
	if limiter == nil {
		t.Fatal("NewRateLimiter returned nil")
	}

	// Consume all tokens
	limiter.Allow()
	limiter.Allow()

	// Should be rejected now
	if limiter.Allow() {
		t.Error("Request allowed, should be rejected (no tokens)")
	}

	// Wait for refill (1 second = 5 tokens, but max is 2)
	time.Sleep(1100 * time.Millisecond)

	// Should allow 2 requests again
	if !limiter.Allow() {
		t.Error("Request rejected after refill, should be allowed")
	}
	if !limiter.Allow() {
		t.Error("Request rejected after refill, should be allowed")
	}

	// Should reject 3rd request
	if limiter.Allow() {
		t.Error("Request allowed, should be rejected (tokens exhausted)")
	}
}

// ============================================================================
// TEST 4: Atomic Max
// ============================================================================

func TestAtomicMax_Basic(t *testing.T) {
	am := NewAtomicMax()
	if am == nil {
		t.Fatal("NewAtomicMax returned nil")
	}

	// Update with first value
	am.Update(10)
	if max := am.Max(); max != 10 {
		t.Errorf("Max = %d, want 10", max)
	}

	// Update with smaller value (should not change)
	am.Update(5)
	if max := am.Max(); max != 10 {
		t.Errorf("Max = %d, want 10 (should not decrease)", max)
	}

	// Update with larger value
	am.Update(20)
	if max := am.Max(); max != 20 {
		t.Errorf("Max = %d, want 20", max)
	}
}

func TestAtomicMax_Concurrent(t *testing.T) {
	am := NewAtomicMax()
	if am == nil {
		t.Fatal("NewAtomicMax returned nil")
	}

	const numGoroutines = 100
	var wg sync.WaitGroup

	// Each goroutine updates with its ID
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			am.Update(id)
		}(int64(i))
	}

	wg.Wait()

	// Max should be numGoroutines - 1
	expected := int64(numGoroutines - 1)
	if max := am.Max(); max != expected {
		t.Errorf("Max = %d, want %d", max, expected)
	}
}

// ============================================================================
// TEST 5: SpinLock
// ============================================================================

func TestSpinLock_Basic(t *testing.T) {
	lock := NewSpinLock()
	if lock == nil {
		t.Fatal("NewSpinLock returned nil")
	}

	// Lock should succeed
	lock.Lock()

	// Unlock should succeed
	lock.Unlock()

	// Should be able to lock again
	lock.Lock()
	lock.Unlock()
}

func TestSpinLock_Concurrent(t *testing.T) {
	lock := NewSpinLock()
	if lock == nil {
		t.Fatal("NewSpinLock returned nil")
	}

	var counter int
	const numGoroutines = 10
	const incrementsPerGoroutine = 1000

	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				lock.Lock()
				counter++
				lock.Unlock()
			}
		}()
	}

	wg.Wait()

	expected := numGoroutines * incrementsPerGoroutine
	if counter != expected {
		t.Errorf("Counter = %d, want %d (spinlock failed to protect counter)", counter, expected)
	}
}

// ============================================================================
// TEST 6: Atomic State Machine
// ============================================================================

func TestAtomicState_Basic(t *testing.T) {
	state := NewAtomicState()
	if state == nil {
		t.Fatal("NewAtomicState returned nil")
	}

	// Test initial state
	if current := state.CurrentState(); current != StateIdle {
		t.Errorf("Initial state = %d, want %d (StateIdle)", current, StateIdle)
	}

	// Transition Idle → Running (should succeed)
	if !state.Transition(StateIdle, StateRunning) {
		t.Error("Transition Idle→Running failed, should succeed")
	}

	// Verify state changed
	if current := state.CurrentState(); current != StateRunning {
		t.Errorf("State = %d, want %d (StateRunning)", current, StateRunning)
	}

	// Transition Idle → Paused (should fail, not in Idle)
	if state.Transition(StateIdle, StatePaused) {
		t.Error("Transition Idle→Paused succeeded, should fail (state is Running)")
	}

	// Transition Running → Paused (should succeed)
	if !state.Transition(StateRunning, StatePaused) {
		t.Error("Transition Running→Paused failed, should succeed")
	}

	// Verify state
	if current := state.CurrentState(); current != StatePaused {
		t.Errorf("State = %d, want %d (StatePaused)", current, StatePaused)
	}
}

// ============================================================================
// TEST 7: Reference Counter
// ============================================================================

func TestReferenceCounter_Basic(t *testing.T) {
	rc := NewReferenceCounter()
	if rc == nil {
		t.Fatal("NewReferenceCounter returned nil")
	}

	// Test initial count
	if count := rc.Count(); count != 0 {
		t.Errorf("Initial count = %d, want 0", count)
	}

	// Acquire reference
	rc.Acquire()
	if count := rc.Count(); count != 1 {
		t.Errorf("After Acquire: count = %d, want 1", count)
	}

	// Acquire more references
	rc.Acquire()
	rc.Acquire()
	if count := rc.Count(); count != 3 {
		t.Errorf("After 3 Acquires: count = %d, want 3", count)
	}

	// Release (should not reach 0)
	if rc.Release() {
		t.Error("Release returned true, should be false (count not 0)")
	}
	if count := rc.Count(); count != 2 {
		t.Errorf("After Release: count = %d, want 2", count)
	}

	// Release again
	if rc.Release() {
		t.Error("Release returned true, should be false (count not 0)")
	}

	// Final release (should reach 0)
	if !rc.Release() {
		t.Error("Release returned false, should be true (count reached 0)")
	}
	if count := rc.Count(); count != 0 {
		t.Errorf("After final Release: count = %d, want 0", count)
	}
}

// ============================================================================
// TEST 8: Config Manager
// ============================================================================

func TestConfigManager_Basic(t *testing.T) {
	cm := NewConfigManager()
	if cm == nil {
		t.Fatal("NewConfigManager returned nil")
	}

	// Test initial config
	cfg := cm.Get()
	if cfg == nil {
		t.Fatal("Get returned nil")
	}
	if cfg.Timeout != 30 || cfg.MaxRetries != 3 || cfg.BatchSize != 100 {
		t.Errorf("Initial config = %+v, want Timeout=30, MaxRetries=3, BatchSize=100", cfg)
	}

	// Update config
	newCfg := &Config{Timeout: 60, MaxRetries: 5, BatchSize: 200}
	cm.Update(newCfg)

	// Verify update
	cfg = cm.Get()
	if cfg.Timeout != 60 || cfg.MaxRetries != 5 || cfg.BatchSize != 200 {
		t.Errorf("After Update: config = %+v, want Timeout=60, MaxRetries=5, BatchSize=200", cfg)
	}
}

func TestConfigManager_Concurrent(t *testing.T) {
	cm := NewConfigManager()
	if cm == nil {
		t.Fatal("NewConfigManager returned nil")
	}

	var wg sync.WaitGroup

	// Writer goroutines
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				cm.Update(&Config{Timeout: id * 10, MaxRetries: id, BatchSize: id * 100})
			}
		}(i)
	}

	// Reader goroutines
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				cfg := cm.Get()
				if cfg == nil {
					t.Error("Get returned nil during concurrent access")
					return
				}
			}
		}()
	}

	wg.Wait()
}

// ============================================================================
// TEST 9: Load Balancer
// ============================================================================

func TestLoadBalancer_Basic(t *testing.T) {
	lb := NewLoadBalancer(3)
	if lb == nil {
		t.Fatal("NewLoadBalancer returned nil")
	}

	// Should round-robin through workers
	expected := []int64{0, 1, 2, 0, 1, 2, 0, 1, 2}
	for i, want := range expected {
		got := lb.NextWorker()
		if got != want {
			t.Errorf("NextWorker()[%d] = %d, want %d", i, got, want)
		}
	}
}

func TestLoadBalancer_Concurrent(t *testing.T) {
	const numWorkers = 5
	lb := NewLoadBalancer(numWorkers)
	if lb == nil {
		t.Fatal("NewLoadBalancer returned nil")
	}

	const numRequests = 10000
	workerCounts := make([]int32, numWorkers)

	var wg sync.WaitGroup

	// Distribute requests
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker := lb.NextWorker()
			atomic.AddInt32(&workerCounts[worker], 1)
		}()
	}

	wg.Wait()

	// Each worker should get roughly equal number of requests
	expectedPerWorker := numRequests / numWorkers
	for i, count := range workerCounts {
		// Allow 10% deviation
		if count < int32(expectedPerWorker*9/10) || count > int32(expectedPerWorker*11/10) {
			t.Errorf("Worker %d got %d requests, want ~%d (±10%%)", i, count, expectedPerWorker)
		}
	}
}

// ============================================================================
// TEST 10: Atomic Bitmap
// ============================================================================

func TestAtomicBitmap_Basic(t *testing.T) {
	bm := NewAtomicBitmap()
	if bm == nil {
		t.Fatal("NewAtomicBitmap returned nil")
	}

	// Test initial state (all bits should be 0)
	for i := 0; i < 256; i++ {
		if bm.TestBit(i) {
			t.Errorf("Bit %d is set, should be clear initially", i)
		}
	}

	// Set some bits
	bm.SetBit(0)
	bm.SetBit(63)
	bm.SetBit(64)
	bm.SetBit(127)
	bm.SetBit(255)

	// Test set bits
	setBits := []int{0, 63, 64, 127, 255}
	for _, bit := range setBits {
		if !bm.TestBit(bit) {
			t.Errorf("Bit %d is clear, should be set", bit)
		}
	}

	// Test clear bits
	for i := 0; i < 256; i++ {
		shouldBeSet := false
		for _, setBit := range setBits {
			if i == setBit {
				shouldBeSet = true
				break
			}
		}
		if shouldBeSet {
			continue
		}
		if bm.TestBit(i) {
			t.Errorf("Bit %d is set, should be clear", i)
		}
	}

	// Clear a bit
	bm.ClearBit(63)
	if bm.TestBit(63) {
		t.Error("Bit 63 is set, should be clear after ClearBit")
	}
}

func TestAtomicBitmap_Concurrent(t *testing.T) {
	bm := NewAtomicBitmap()
	if bm == nil {
		t.Fatal("NewAtomicBitmap returned nil")
	}

	var wg sync.WaitGroup

	// Multiple goroutines set different bits
	for i := 0; i < 256; i++ {
		wg.Add(1)
		go func(bit int) {
			defer wg.Done()
			bm.SetBit(bit)
		}(i)
	}

	wg.Wait()

	// All bits should be set
	for i := 0; i < 256; i++ {
		if !bm.TestBit(i) {
			t.Errorf("Bit %d is clear, should be set", i)
		}
	}
}

// ============================================================================
// BENCHMARK TESTS
// ============================================================================

func BenchmarkAtomicCounter(b *testing.B) {
	counter := NewAtomicCounter()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Increment()
		}
	})
}

func BenchmarkMutexCounter(b *testing.B) {
	var counter int64
	var mu sync.Mutex
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			counter++
			mu.Unlock()
		}
	})
}

func BenchmarkAtomicFlag(b *testing.B) {
	flag := NewAtomicFlag()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			flag.TestAndSet()
		}
	})
}

func BenchmarkCompareAndSwap(b *testing.B) {
	var value int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for {
				old := atomic.LoadInt64(&value)
				if atomic.CompareAndSwapInt64(&value, old, old+1) {
					break
				}
			}
		}
	})
}

// ============================================================================
// RACE DETECTOR TESTS
// ============================================================================

// TestRaceCondition_NonAtomic should fail with -race flag.
func TestRaceCondition_NonAtomic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping race condition test in short mode")
	}

	var counter int64
	const numGoroutines = 100
	const incrementsPerGoroutine = 100

	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				IncrementCounterNonAtomic(&counter)
			}
		}()
	}

	wg.Wait()

	// This will likely fail (race condition)
	expected := int64(numGoroutines * incrementsPerGoroutine)
	if counter != expected {
		t.Logf("Race condition detected: counter = %d, expected %d", counter, expected)
	}
}

// TestNoRaceCondition_Atomic should pass even with -race flag.
func TestNoRaceCondition_Atomic(t *testing.T) {
	var counter int64
	const numGoroutines = 100
	const incrementsPerGoroutine = 100

	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				IncrementCounterAtomic(&counter)
			}
		}()
	}

	wg.Wait()

	expected := int64(numGoroutines * incrementsPerGoroutine)
	if counter != expected {
		t.Errorf("Counter = %d, want %d (atomic operations failed)", counter, expected)
	}
}
