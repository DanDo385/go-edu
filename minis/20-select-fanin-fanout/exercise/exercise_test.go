package exercise

import (
	"sort"
	"testing"
	"time"
)

// ============================================================================
// Tests for Exercise 1: SelectFirst
// ============================================================================

func TestSelectFirst_Ch1First(t *testing.T) {
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	ch1 <- "from ch1"

	value, ok := SelectFirst(ch1, ch2, 1*time.Second)
	if !ok {
		t.Fatal("Expected to receive value")
	}
	if value != "from ch1" {
		t.Errorf("Expected 'from ch1', got '%s'", value)
	}
}

func TestSelectFirst_Ch2First(t *testing.T) {
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	ch2 <- "from ch2"

	value, ok := SelectFirst(ch1, ch2, 1*time.Second)
	if !ok {
		t.Fatal("Expected to receive value")
	}
	if value != "from ch2" {
		t.Errorf("Expected 'from ch2', got '%s'", value)
	}
}

func TestSelectFirst_Timeout(t *testing.T) {
	ch1 := make(chan string)
	ch2 := make(chan string)

	value, ok := SelectFirst(ch1, ch2, 100*time.Millisecond)
	if ok {
		t.Errorf("Expected timeout, got value: %s", value)
	}
	if value != "" {
		t.Errorf("Expected empty string on timeout, got '%s'", value)
	}
}

func TestSelectFirst_SlowSender(t *testing.T) {
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch1 <- "delayed"
	}()

	value, ok := SelectFirst(ch1, ch2, 200*time.Millisecond)
	if !ok {
		t.Fatal("Expected to receive value")
	}
	if value != "delayed" {
		t.Errorf("Expected 'delayed', got '%s'", value)
	}
}

// ============================================================================
// Tests for Exercise 2: NonBlockingSend
// ============================================================================

func TestNonBlockingSend_UnbufferedNoReceiver(t *testing.T) {
	ch := make(chan int)

	sent := NonBlockingSend(ch, 42)
	if sent {
		t.Error("Expected send to fail (no receiver)")
	}
}

func TestNonBlockingSend_BufferedWithSpace(t *testing.T) {
	ch := make(chan int, 1)

	sent := NonBlockingSend(ch, 42)
	if !sent {
		t.Error("Expected send to succeed (buffered channel)")
	}

	received := <-ch
	if received != 42 {
		t.Errorf("Expected 42, got %d", received)
	}
}

func TestNonBlockingSend_BufferedFull(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 1 // Fill buffer

	sent := NonBlockingSend(ch, 42)
	if sent {
		t.Error("Expected send to fail (buffer full)")
	}
}

func TestNonBlockingSend_WithReceiver(t *testing.T) {
	ch := make(chan int)

	go func() {
		<-ch
	}()

	// Give receiver time to start
	time.Sleep(10 * time.Millisecond)

	sent := NonBlockingSend(ch, 42)
	if !sent {
		t.Error("Expected send to succeed (receiver ready)")
	}
}

// ============================================================================
// Tests for Exercise 3: FanIn
// ============================================================================

func TestFanIn_TwoChannels(t *testing.T) {
	ch1 := make(chan int, 3)
	ch2 := make(chan int, 3)

	ch1 <- 1
	ch1 <- 2
	ch1 <- 3
	close(ch1)

	ch2 <- 4
	ch2 <- 5
	ch2 <- 6
	close(ch2)

	merged := FanIn(ch1, ch2)

	var results []int
	for v := range merged {
		results = append(results, v)
	}

	if len(results) != 6 {
		t.Errorf("Expected 6 values, got %d", len(results))
	}

	sort.Ints(results)
	expected := []int{1, 2, 3, 4, 5, 6}
	for i, v := range expected {
		if results[i] != v {
			t.Errorf("Expected %d at index %d, got %d", v, i, results[i])
		}
	}
}

func TestFanIn_EmptyChannels(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)

	close(ch1)
	close(ch2)

	merged := FanIn(ch1, ch2)

	count := 0
	for range merged {
		count++
	}

	if count != 0 {
		t.Errorf("Expected 0 values, got %d", count)
	}
}

func TestFanIn_MultipleChannels(t *testing.T) {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch3 := make(chan int, 1)

	ch1 <- 1
	ch2 <- 2
	ch3 <- 3
	close(ch1)
	close(ch2)
	close(ch3)

	merged := FanIn(ch1, ch2, ch3)

	var results []int
	for v := range merged {
		results = append(results, v)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 values, got %d", len(results))
	}

	sort.Ints(results)
	expected := []int{1, 2, 3}
	for i, v := range expected {
		if results[i] != v {
			t.Errorf("Expected %d, got %d", v, results[i])
		}
	}
}

// ============================================================================
// Tests for Exercise 4: FanOut
// ============================================================================

func TestFanOut_Basic(t *testing.T) {
	input := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		input <- i
	}
	close(input)

	square := func(n int) int { return n * n }
	results := FanOut(input, 3, square)

	var values []int
	for v := range results {
		values = append(values, v)
	}

	if len(values) != 5 {
		t.Errorf("Expected 5 results, got %d", len(values))
	}

	sort.Ints(values)
	expected := []int{1, 4, 9, 16, 25}
	for i, v := range expected {
		if values[i] != v {
			t.Errorf("Expected %d, got %d", v, values[i])
		}
	}
}

func TestFanOut_MoreWorkersThanTasks(t *testing.T) {
	input := make(chan int, 2)
	input <- 1
	input <- 2
	close(input)

	double := func(n int) int { return n * 2 }
	results := FanOut(input, 10, double)

	var values []int
	for v := range results {
		values = append(values, v)
	}

	if len(values) != 2 {
		t.Errorf("Expected 2 results, got %d", len(values))
	}

	sort.Ints(values)
	expected := []int{2, 4}
	for i, v := range expected {
		if values[i] != v {
			t.Errorf("Expected %d, got %d", v, values[i])
		}
	}
}

func TestFanOut_EmptyInput(t *testing.T) {
	input := make(chan int)
	close(input)

	identity := func(n int) int { return n }
	results := FanOut(input, 3, identity)

	count := 0
	for range results {
		count++
	}

	if count != 0 {
		t.Errorf("Expected 0 results, got %d", count)
	}
}

// ============================================================================
// Tests for Exercise 5: OrChannel
// ============================================================================

func TestOrChannel_FirstCloses(t *testing.T) {
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	ch3 := make(chan struct{})

	done := OrChannel(ch1, ch2, ch3)

	go func() {
		time.Sleep(50 * time.Millisecond)
		close(ch1)
	}()

	select {
	case <-done:
		// Success - done closed
	case <-time.After(200 * time.Millisecond):
		t.Error("OrChannel did not close when input closed")
	}
}

func TestOrChannel_AllClosed(t *testing.T) {
	ch1 := make(chan struct{})
	ch2 := make(chan struct{})

	close(ch1)
	close(ch2)

	done := OrChannel(ch1, ch2)

	select {
	case <-done:
		// Success - done closed immediately
	case <-time.After(100 * time.Millisecond):
		t.Error("OrChannel did not close for already-closed channels")
	}
}

func TestOrChannel_NoChannels(t *testing.T) {
	done := OrChannel()

	select {
	case <-done:
		// Success - done closed immediately
	case <-time.After(100 * time.Millisecond):
		t.Error("OrChannel should close immediately with no inputs")
	}
}

func TestOrChannel_SingleChannel(t *testing.T) {
	ch := make(chan struct{})

	done := OrChannel(ch)

	go func() {
		time.Sleep(50 * time.Millisecond)
		close(ch)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(200 * time.Millisecond):
		t.Error("OrChannel did not close")
	}
}

// ============================================================================
// Tests for Exercise 6: TryReceiveAll
// ============================================================================

func TestTryReceiveAll_SomeHaveValues(t *testing.T) {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch3 := make(chan int, 1)

	ch1 <- 10
	ch3 <- 30

	result := TryReceiveAll([]<-chan int{ch1, ch2, ch3})

	if len(result) != 2 {
		t.Errorf("Expected 2 values, got %d", len(result))
	}

	if result[0] != 10 {
		t.Errorf("Expected result[0]=10, got %d", result[0])
	}

	if result[2] != 30 {
		t.Errorf("Expected result[2]=30, got %d", result[2])
	}

	if _, exists := result[1]; exists {
		t.Error("Expected result[1] not to exist")
	}
}

func TestTryReceiveAll_AllEmpty(t *testing.T) {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	result := TryReceiveAll([]<-chan int{ch1, ch2})

	if len(result) != 0 {
		t.Errorf("Expected empty map, got %d entries", len(result))
	}
}

func TestTryReceiveAll_AllHaveValues(t *testing.T) {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)

	ch1 <- 1
	ch2 <- 2

	result := TryReceiveAll([]<-chan int{ch1, ch2})

	if len(result) != 2 {
		t.Errorf("Expected 2 values, got %d", len(result))
	}

	if result[0] != 1 {
		t.Errorf("Expected result[0]=1, got %d", result[0])
	}

	if result[1] != 2 {
		t.Errorf("Expected result[1]=2, got %d", result[1])
	}
}

// ============================================================================
// Tests for Exercise 7: RateLimiter
// ============================================================================

func TestRateLimiter_Basic(t *testing.T) {
	limiter := NewRateLimiter(10) // 10 ops/sec

	// First 10 operations should be instant
	start := time.Now()
	for i := 0; i < 10; i++ {
		limiter.Wait()
	}
	elapsed := time.Since(start)

	if elapsed > 100*time.Millisecond {
		t.Errorf("First 10 operations took too long: %v", elapsed)
	}
}

func TestRateLimiter_RateLimit(t *testing.T) {
	limiter := NewRateLimiter(10) // 10 ops/sec

	// Consume initial tokens
	for i := 0; i < 10; i++ {
		limiter.Wait()
	}

	// Next 10 should take ~1 second
	start := time.Now()
	for i := 0; i < 10; i++ {
		limiter.Wait()
	}
	elapsed := time.Since(start)

	if elapsed < 900*time.Millisecond {
		t.Errorf("Rate limiting too lenient: %v", elapsed)
	}
	if elapsed > 1200*time.Millisecond {
		t.Errorf("Rate limiting too strict: %v", elapsed)
	}
}

func TestRateLimiter_TryWait(t *testing.T) {
	limiter := NewRateLimiter(5)

	// First 5 should succeed
	for i := 0; i < 5; i++ {
		if !limiter.TryWait() {
			t.Errorf("TryWait() failed on operation %d", i)
		}
	}

	// Next should fail (no tokens)
	if limiter.TryWait() {
		t.Error("TryWait() should have failed (no tokens)")
	}

	// Wait for refill
	time.Sleep(250 * time.Millisecond)

	// Should succeed now
	if !limiter.TryWait() {
		t.Error("TryWait() should have succeeded after refill")
	}
}

// ============================================================================
// Tests for Exercise 8: Pipeline
// ============================================================================

func TestPipeline_Basic(t *testing.T) {
	results := Pipeline(10, 3)

	var values []int
	for v := range results {
		values = append(values, v)
	}

	// Expected: squares of 1-10 that are even
	// 1^2=1 (odd), 2^2=4 (even), 3^2=9 (odd), 4^2=16 (even), 5^2=25 (odd),
	// 6^2=36 (even), 7^2=49 (odd), 8^2=64 (even), 9^2=81 (odd), 10^2=100 (even)
	// Result: 4, 16, 36, 64, 100

	if len(values) != 5 {
		t.Errorf("Expected 5 values, got %d", len(values))
	}

	sort.Ints(values)
	expected := []int{4, 16, 36, 64, 100}
	for i, v := range expected {
		if values[i] != v {
			t.Errorf("Expected %d, got %d", v, values[i])
		}
	}
}

func TestPipeline_SmallN(t *testing.T) {
	results := Pipeline(3, 2)

	var values []int
	for v := range results {
		values = append(values, v)
	}

	// 1^2=1 (odd), 2^2=4 (even), 3^2=9 (odd)
	// Result: 4

	if len(values) != 1 {
		t.Errorf("Expected 1 value, got %d", len(values))
	}

	if values[0] != 4 {
		t.Errorf("Expected 4, got %d", values[0])
	}
}

func TestPipeline_AllOdd(t *testing.T) {
	results := Pipeline(1, 1)

	var values []int
	for v := range results {
		values = append(values, v)
	}

	// 1^2=1 (odd)
	// Result: empty

	if len(values) != 0 {
		t.Errorf("Expected 0 values, got %d", len(values))
	}
}

// ============================================================================
// Tests for Exercise 9: SelectWithPriority
// ============================================================================

func TestSelectWithPriority_BothReady(t *testing.T) {
	high := make(chan int, 1)
	low := make(chan int, 1)

	high <- 1
	low <- 2

	value, isHigh := SelectWithPriority(high, low)

	if !isHigh {
		t.Error("Expected high priority to be selected")
	}
	if value != 1 {
		t.Errorf("Expected 1, got %d", value)
	}
}

func TestSelectWithPriority_OnlyLow(t *testing.T) {
	high := make(chan int, 1)
	low := make(chan int, 1)

	low <- 2

	value, isHigh := SelectWithPriority(high, low)

	if isHigh {
		t.Error("Expected low priority to be selected")
	}
	if value != 2 {
		t.Errorf("Expected 2, got %d", value)
	}
}

func TestSelectWithPriority_OnlyHigh(t *testing.T) {
	high := make(chan int, 1)
	low := make(chan int, 1)

	high <- 1

	value, isHigh := SelectWithPriority(high, low)

	if !isHigh {
		t.Error("Expected high priority to be selected")
	}
	if value != 1 {
		t.Errorf("Expected 1, got %d", value)
	}
}

func TestSelectWithPriority_Blocking(t *testing.T) {
	high := make(chan int, 1)
	low := make(chan int, 1)

	go func() {
		time.Sleep(50 * time.Millisecond)
		low <- 42
	}()

	done := make(chan bool)
	go func() {
		value, isHigh := SelectWithPriority(high, low)
		if isHigh {
			t.Error("Expected low priority")
		}
		if value != 42 {
			t.Errorf("Expected 42, got %d", value)
		}
		done <- true
	}()

	select {
	case <-done:
		// Success
	case <-time.After(200 * time.Millisecond):
		t.Error("SelectWithPriority did not receive value")
	}
}

// ============================================================================
// Tests for Exercise 10: Timeout
// ============================================================================

func TestTimeout_ValueBeforeTimeout(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 42

	value, ok := Timeout(ch, 1*time.Second)
	if !ok {
		t.Error("Expected to receive value")
	}
	if value != 42 {
		t.Errorf("Expected 42, got %d", value)
	}
}

func TestTimeout_TimeoutOccurs(t *testing.T) {
	ch := make(chan int)

	value, ok := Timeout(ch, 50*time.Millisecond)
	if ok {
		t.Errorf("Expected timeout, got value: %d", value)
	}
	if value != 0 {
		t.Errorf("Expected 0 on timeout, got %d", value)
	}
}

func TestTimeout_ValueAfterDelay(t *testing.T) {
	ch := make(chan int, 1)

	go func() {
		time.Sleep(50 * time.Millisecond)
		ch <- 99
	}()

	value, ok := Timeout(ch, 200*time.Millisecond)
	if !ok {
		t.Error("Expected to receive value")
	}
	if value != 99 {
		t.Errorf("Expected 99, got %d", value)
	}
}

func TestTimeout_ValueArrivesTooLate(t *testing.T) {
	ch := make(chan int, 1)

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch <- 99
	}()

	value, ok := Timeout(ch, 50*time.Millisecond)
	if ok {
		t.Errorf("Expected timeout, got value: %d", value)
	}
	if value != 0 {
		t.Errorf("Expected 0 on timeout, got %d", value)
	}
}
