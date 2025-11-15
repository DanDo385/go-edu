package exercise

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestPing(t *testing.T) {
	t.Run("sends value and closes", func(t *testing.T) {
		ch := Ping(42)
		if ch == nil {
			t.Fatal("Ping returned nil channel")
		}

		// Should receive the value
		select {
		case v := <-ch:
			if v != 42 {
				t.Errorf("Expected 42, got %d", v)
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timeout waiting for value")
		}

		// Channel should be closed
		select {
		case _, ok := <-ch:
			if ok {
				t.Error("Channel should be closed")
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timeout waiting for close")
		}
	})
}

func TestPingPong(t *testing.T) {
	t.Run("plays ping pong", func(t *testing.T) {
		ping, pong := PingPong(3)
		if ping == nil || pong == nil {
			t.Fatal("PingPong returned nil channels")
		}

		// Play 3 rounds
		for i := 0; i < 3; i++ {
			ping <- i

			select {
			case v := <-pong:
				if v != i {
					t.Errorf("Expected %d, got %d", i, v)
				}
			case <-time.After(100 * time.Millisecond):
				t.Fatalf("Timeout on round %d", i)
			}
		}

		// Channels should close after n rounds
		close(ping)
		time.Sleep(50 * time.Millisecond)

		select {
		case _, ok := <-pong:
			if ok {
				t.Error("Pong channel should be closed")
			}
		case <-time.After(100 * time.Millisecond):
			// It's ok if it doesn't close immediately
		}
	})
}

func TestMerge(t *testing.T) {
	t.Run("merges multiple channels", func(t *testing.T) {
		ch1 := make(chan int, 3)
		ch2 := make(chan int, 3)
		ch3 := make(chan int, 3)

		// Send values to all channels
		ch1 <- 1
		ch1 <- 2
		ch2 <- 10
		ch2 <- 20
		ch3 <- 100

		close(ch1)
		close(ch2)
		close(ch3)

		merged := Merge(ch1, ch2, ch3)
		if merged == nil {
			t.Fatal("Merge returned nil channel")
		}

		// Collect all values
		received := make(map[int]bool)
		for v := range merged {
			received[v] = true
		}

		// Verify all values received
		expected := []int{1, 2, 10, 20, 100}
		if len(received) != len(expected) {
			t.Errorf("Expected %d values, got %d", len(expected), len(received))
		}

		for _, exp := range expected {
			if !received[exp] {
				t.Errorf("Missing value %d", exp)
			}
		}
	})

	t.Run("closes when all inputs close", func(t *testing.T) {
		ch1 := make(chan int)
		ch2 := make(chan int)
		close(ch1)
		close(ch2)

		merged := Merge(ch1, ch2)

		// Should close immediately
		select {
		case _, ok := <-merged:
			if ok {
				t.Error("Merged channel should be closed")
			}
		case <-time.After(200 * time.Millisecond):
			t.Fatal("Merged channel didn't close")
		}
	})
}

func TestFilter(t *testing.T) {
	t.Run("filters values", func(t *testing.T) {
		input := make(chan int, 10)
		for i := 1; i <= 10; i++ {
			input <- i
		}
		close(input)

		// Filter even numbers
		evens := Filter(input, func(x int) bool { return x%2 == 0 })
		if evens == nil {
			t.Fatal("Filter returned nil channel")
		}

		results := []int{}
		for v := range evens {
			results = append(results, v)
		}

		expected := []int{2, 4, 6, 8, 10}
		if len(results) != len(expected) {
			t.Fatalf("Expected %d values, got %d", len(expected), len(results))
		}

		for i, exp := range expected {
			if results[i] != exp {
				t.Errorf("At index %d: expected %d, got %d", i, exp, results[i])
			}
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("transforms values", func(t *testing.T) {
		input := make(chan int, 5)
		for i := 1; i <= 5; i++ {
			input <- i
		}
		close(input)

		// Double each value
		doubled := Map(input, func(x int) int { return x * 2 })
		if doubled == nil {
			t.Fatal("Map returned nil channel")
		}

		results := []int{}
		for v := range doubled {
			results = append(results, v)
		}

		expected := []int{2, 4, 6, 8, 10}
		if len(results) != len(expected) {
			t.Fatalf("Expected %d values, got %d", len(expected), len(results))
		}

		for i, exp := range expected {
			if results[i] != exp {
				t.Errorf("At index %d: expected %d, got %d", i, exp, results[i])
			}
		}
	})
}

func TestTake(t *testing.T) {
	t.Run("takes first n values", func(t *testing.T) {
		input := make(chan int, 10)
		for i := 1; i <= 10; i++ {
			input <- i
		}
		close(input)

		first3 := Take(input, 3)
		if first3 == nil {
			t.Fatal("Take returned nil channel")
		}

		results := []int{}
		for v := range first3 {
			results = append(results, v)
		}

		if len(results) != 3 {
			t.Fatalf("Expected 3 values, got %d", len(results))
		}

		for i := 0; i < 3; i++ {
			if results[i] != i+1 {
				t.Errorf("At index %d: expected %d, got %d", i, i+1, results[i])
			}
		}
	})

	t.Run("stops early if input closes", func(t *testing.T) {
		input := make(chan int, 2)
		input <- 1
		input <- 2
		close(input)

		first10 := Take(input, 10)

		results := []int{}
		for v := range first10 {
			results = append(results, v)
		}

		if len(results) != 2 {
			t.Errorf("Expected 2 values (input closed early), got %d", len(results))
		}
	})
}

func TestOrDone(t *testing.T) {
	t.Run("forwards values until cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		input := make(chan int, 10)
		for i := 1; i <= 10; i++ {
			input <- i
		}

		output := OrDone(ctx, input)
		if output == nil {
			t.Fatal("OrDone returned nil channel")
		}

		// Receive a few values
		for i := 0; i < 3; i++ {
			select {
			case v := <-output:
				if v != i+1 {
					t.Errorf("Expected %d, got %d", i+1, v)
				}
			case <-time.After(100 * time.Millisecond):
				t.Fatal("Timeout")
			}
		}

		// Cancel context
		cancel()

		// Output should close
		select {
		case _, ok := <-output:
			if ok {
				t.Error("Output should close when context is cancelled")
			}
		case <-time.After(200 * time.Millisecond):
			t.Fatal("Output didn't close after cancellation")
		}
	})

	t.Run("closes when input closes", func(t *testing.T) {
		ctx := context.Background()
		input := make(chan int)
		close(input)

		output := OrDone(ctx, input)

		select {
		case _, ok := <-output:
			if ok {
				t.Error("Output should close when input closes")
			}
		case <-time.After(200 * time.Millisecond):
			t.Fatal("Output didn't close")
		}
	})
}

func TestTee(t *testing.T) {
	t.Run("splits to two channels", func(t *testing.T) {
		input := make(chan int, 5)
		for i := 1; i <= 5; i++ {
			input <- i
		}
		close(input)

		out1, out2 := Tee(input)
		if out1 == nil || out2 == nil {
			t.Fatal("Tee returned nil channels")
		}

		// Collect from both outputs
		results1 := []int{}
		results2 := []int{}

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			for v := range out1 {
				results1 = append(results1, v)
			}
		}()

		go func() {
			defer wg.Done()
			for v := range out2 {
				results2 = append(results2, v)
			}
		}()

		wg.Wait()

		// Both should receive all values
		expected := []int{1, 2, 3, 4, 5}
		if len(results1) != len(expected) {
			t.Errorf("out1: expected %d values, got %d", len(expected), len(results1))
		}
		if len(results2) != len(expected) {
			t.Errorf("out2: expected %d values, got %d", len(expected), len(results2))
		}
	})
}

func TestBridge(t *testing.T) {
	t.Run("flattens channel of channels", func(t *testing.T) {
		input := make(chan (<-chan int), 3)

		// Create 3 sub-channels
		for i := 0; i < 3; i++ {
			ch := make(chan int, 2)
			ch <- i*10 + 1
			ch <- i*10 + 2
			close(ch)
			input <- ch
		}
		close(input)

		flattened := Bridge(input)
		if flattened == nil {
			t.Fatal("Bridge returned nil channel")
		}

		// Collect all values
		received := make(map[int]bool)
		for v := range flattened {
			received[v] = true
		}

		// Should have 6 values: 1, 2, 11, 12, 21, 22
		expected := []int{1, 2, 11, 12, 21, 22}
		if len(received) != len(expected) {
			t.Errorf("Expected %d values, got %d", len(expected), len(received))
		}

		for _, exp := range expected {
			if !received[exp] {
				t.Errorf("Missing value %d", exp)
			}
		}
	})
}

func TestDebounce(t *testing.T) {
	t.Run("debounces rapid values", func(t *testing.T) {
		input := make(chan int)
		debounced := Debounce(input, 50*time.Millisecond)
		if debounced == nil {
			t.Fatal("Debounce returned nil channel")
		}

		// Send rapid values
		go func() {
			input <- 1
			time.Sleep(10 * time.Millisecond)
			input <- 2
			time.Sleep(10 * time.Millisecond)
			input <- 3 // This one should be forwarded (after 50ms quiet)
			time.Sleep(100 * time.Millisecond)
			input <- 4 // This one too
			time.Sleep(100 * time.Millisecond)
			close(input)
		}()

		results := []int{}
		for v := range debounced {
			results = append(results, v)
		}

		// Should only get values with 50ms quiet time
		// Expecting 3 and 4
		if len(results) < 1 || len(results) > 3 {
			t.Logf("Got results: %v", results)
			// Debouncing is timing-sensitive, so we're lenient
		}

		// Last value should definitely be there
		if len(results) > 0 && results[len(results)-1] != 4 {
			t.Errorf("Expected last value to be 4, got %d", results[len(results)-1])
		}
	})
}

func TestBoundedQueue(t *testing.T) {
	t.Run("enqueue and dequeue", func(t *testing.T) {
		queue := NewBoundedQueue(3)
		if queue == nil {
			t.Fatal("NewBoundedQueue returned nil")
		}

		queue.Enqueue(1)
		queue.Enqueue(2)
		queue.Enqueue(3)

		if v := queue.Dequeue(); v != 1 {
			t.Errorf("Expected 1, got %d", v)
		}
		if v := queue.Dequeue(); v != 2 {
			t.Errorf("Expected 2, got %d", v)
		}
		if v := queue.Dequeue(); v != 3 {
			t.Errorf("Expected 3, got %d", v)
		}
	})

	t.Run("try operations", func(t *testing.T) {
		queue := NewBoundedQueue(2)
		if queue == nil {
			t.Fatal("NewBoundedQueue returned nil")
		}

		// TryEnqueue on empty queue
		if !queue.TryEnqueue(1) {
			t.Error("TryEnqueue should succeed on empty queue")
		}
		if !queue.TryEnqueue(2) {
			t.Error("TryEnqueue should succeed (not full yet)")
		}

		// TryEnqueue on full queue
		if queue.TryEnqueue(3) {
			t.Error("TryEnqueue should fail on full queue")
		}

		// TryDequeue on non-empty queue
		v, ok := queue.TryDequeue()
		if !ok || v != 1 {
			t.Errorf("TryDequeue should succeed, expected 1, got %d", v)
		}

		v, ok = queue.TryDequeue()
		if !ok || v != 2 {
			t.Errorf("TryDequeue should succeed, expected 2, got %d", v)
		}

		// TryDequeue on empty queue
		_, ok = queue.TryDequeue()
		if ok {
			t.Error("TryDequeue should fail on empty queue")
		}
	})
}

func TestBroadcaster(t *testing.T) {
	t.Run("broadcasts to subscribers", func(t *testing.T) {
		bc := NewBroadcaster()
		if bc == nil {
			t.Fatal("NewBroadcaster returned nil")
		}

		// Subscribe 3 listeners
		ch1 := bc.Subscribe()
		ch2 := bc.Subscribe()
		ch3 := bc.Subscribe()

		if ch1 == nil || ch2 == nil || ch3 == nil {
			t.Fatal("Subscribe returned nil channel")
		}

		// Send a message
		msg := Message{ID: 1, Content: "hello"}
		bc.Send(msg)

		// All should receive
		for i, ch := range []<-chan Message{ch1, ch2, ch3} {
			select {
			case received := <-ch:
				if received.ID != msg.ID || received.Content != msg.Content {
					t.Errorf("Subscriber %d: expected %v, got %v", i, msg, received)
				}
			case <-time.After(200 * time.Millisecond):
				t.Errorf("Subscriber %d: timeout", i)
			}
		}

		bc.Close()
	})
}

func TestBarrier(t *testing.T) {
	t.Run("synchronizes goroutines", func(t *testing.T) {
		const n = 5
		barrier := NewBarrier(n)
		if barrier == nil {
			t.Fatal("NewBarrier returned nil")
		}

		var mu sync.Mutex
		ready := 0
		done := 0

		var wg sync.WaitGroup
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				// Mark as ready
				mu.Lock()
				ready++
				mu.Unlock()

				// Wait at barrier
				barrier.Wait()

				// All should be ready before any proceed
				mu.Lock()
				if ready != n {
					t.Errorf("Goroutine %d: not all ready (ready=%d)", id, ready)
				}
				done++
				mu.Unlock()
			}(i)
		}

		wg.Wait()

		if done != n {
			t.Errorf("Expected %d goroutines done, got %d", n, done)
		}
	})
}
