//go:build solution
// +build solution

package exercise

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"unicode"

	"golang.org/x/sync/errgroup"
)

/*
WordCount fetches a list of URLs concurrently and aggregates word frequencies.

This function is a *textbook worker-pool pattern*:

  - A bounded number of worker goroutines (N = workers)
  - A jobs channel distributes work
  - A results channel collects outputs
  - A shared context coordinates cancellation
  - A WaitGroup coordinates shutdown

IMPORTANT CONCURRENCY RULE:
---------------------------
At any instant:
  - At most `workers` HTTP requests are in flight
  - At most `workers` goroutines are actively tokenizing text
  - No goroutine ever mutates shared state unsafely

MEMORY MANAGEMENT OVERVIEW:
----------------------------
1. Channels: Each buffered channel allocates memory for its buffer
  - jobs channel: workers * sizeof(string) ≈ workers * 16-32 bytes
  - results channel: workers * sizeof(map[string]int) ≈ workers * 8 bytes (map header)
  - errCh: 1 * sizeof(error) ≈ 16 bytes
    Total channel memory: ~(workers * 24) + 16 bytes

2. Goroutines: Each goroutine starts with ~2KB stack, grows as needed
  - N workers: N * 2KB initial
  - Main goroutine: 2KB
  - Helper goroutines (job sender, results closer): 2KB each
    Total goroutine memory: ~(N + 3) * 2KB

3. Maps: Each worker creates a map[string]int per URL
  - Maps grow dynamically (rehash when load factor > 6.5)
  - Memory = O(vocabulary_size) per URL
  - Final aggregation merges all maps into one

4. HTTP responses: Each response body is read into memory
  - Memory = O(response_size) per concurrent request
  - Maximum concurrent memory = workers * average_response_size

VARIABLE FLOW THROUGH FUNCTIONS:
---------------------------------
urls ([]string) → jobs channel → worker goroutines → fetchAndCount() → HTTP request

	        ↓
	response body ([]byte)
	        ↓
	tokenizeAndCount() → map[string]int
	        ↓
	results channel → finalCounts (map[string]int)
*/
func WordCount(ctx context.Context, urls []string, workers int) (map[string]int, error) {
	/*
	   UNDERSTANDING THE CONTEXT PARAMETER
	   -----------------------------------

	   What is `ctx context.Context`?
	   -------------------------------
	   Context is Go's standard way to:
	   1. Carry cancellation signals across function boundaries
	   2. Propagate timeouts and deadlines
	   3. Pass request-scoped values (like request IDs)

	   Think of context as a "stop button" that can be pressed anywhere,
	   and all functions listening to it will stop immediately.

	   When you pass `ctx` to a function, you're giving it:
	   - The ability to check if work should be cancelled: <-ctx.Done()
	   - The ability to respect timeouts/deadlines
	   - A way to propagate cancellation to child operations

	   In this function:
	   - We receive `ctx` from the caller (might have a timeout)
	   - We create a child context with cancellation: ctx, cancel := context.WithCancel(ctx)
	   - We pass this child context to all workers and HTTP requests
	   - If any worker fails, we call cancel() → all workers stop

	   Memory: Context is a small struct (~48 bytes), passed by value (copied)
	           but internally contains pointers to shared cancellation state.
	*/

	/*
	   CHANNEL DESIGN AND MEMORY ALLOCATION
	   ------------------------------------

	   What is a Channel?
	   ------------------
	   A channel is a typed, thread-safe queue for communication between goroutines.
	   Think of it like a conveyor belt:
	   - Sender puts items on: ch <- value
	   - Receiver takes items off: value := <-ch
	   - If buffer is full, sender blocks
	   - If buffer is empty, receiver blocks

	   jobs channel:
	     - Type: chan string (carries URLs)
	     - Buffer size: workers
	     - Purpose: Distribute URLs to workers
	     - Memory: workers * ~16-32 bytes (string header + pointer)
	     - Producer: Job sender goroutine
	     - Consumers: Worker goroutines

	   results channel:
	     - Type: chan map[string]int (carries word frequency maps)
	     - Buffer size: workers
	     - Purpose: Collect results from workers
	     - Memory: workers * 8 bytes (map header pointer)
	     - Producers: Worker goroutines
	     - Consumer: Main goroutine (aggregator)

	   errCh channel:
	     - Type: chan error (carries error messages)
	     - Buffer size: 1 (only need first error)
	     - Purpose: Signal that an error occurred
	     - Memory: 1 * ~16 bytes (error interface)
	     - Why buffered size 1?
	       * Ensures first error sender never blocks
	       * If multiple workers error simultaneously, only first is recorded
	       * Non-blocking send pattern: select { case errCh <- err: ... default: }
	*/

	jobs := make(chan string, workers)
	results := make(chan map[string]int, workers)
	errCh := make(chan error, 1)

	/*
	   CONTEXT
	   -------

	   We derive a cancellable context from the caller’s context.

	   Why?
	   - If *any* worker fails, we want to:
	       * stop issuing new jobs
	       * abort in-flight HTTP requests
	       * make other workers exit quickly

	   Context cancellation is Go’s standard "broadcast shutdown signal".
	*/
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // Safety net: ensures cleanup on return

	/*
	   UNDERSTANDING sync.WaitGroup
	   -----------------------------

	   What is WaitGroup?
	   ------------------
	   WaitGroup is a counter that tracks how many goroutines are still running.
	   It provides three operations:
	     1. Add(delta): Increment counter by delta
	     2. Done(): Decrement counter by 1 (shorthand for Add(-1))
	     3. Wait(): Block until counter reaches 0

	   How it works internally:
	   ------------------------
	   WaitGroup uses atomic operations (lock-free) for Add/Done.
	   Wait() spins or blocks until the counter is zero.

	   Memory: WaitGroup is a struct with ~12 bytes (counter + semaphore)
	           Very lightweight, no heap allocation.

	   Why track only workers?
	   -----------------------
	   - Workers are the "heavy lifters" (do the actual work)
	   - Job sender finishes quickly (just sends URLs)
	   - Results closer waits for workers anyway
	   - Once workers finish, we know all work is done

	   Pattern:
	   --------
	   wg.Add(N)        // Before starting N goroutines
	   go func() {
	       defer wg.Done()  // When goroutine finishes
	       // work
	   }()
	   wg.Wait()        // Block until all Done() calls complete
	*/
	var wg sync.WaitGroup

	/*
	   WORKER SPAWN LOOP - UNDERSTANDING ITERATION AND COUNTERS
	   --------------------------------------------------------

	   This loop runs SEQUENTIALLY in the main goroutine.
	   Each iteration:
	     1. Increments WaitGroup counter: wg.Add(1)
	     2. Launches a goroutine (runs concurrently)
	     3. Increments loop counter: i++

	   The goroutines themselves run CONCURRENTLY (in parallel).

	   HOW THE COUNTER `i` WORKS:
	   ---------------------------
	   - `i` starts at 0, increments each iteration: 0, 1, 2, ..., workers-1
	   - Each iteration creates ONE goroutine
	   - After `workers` iterations, we have `workers` goroutines running

	   WHY PASS `i` AS PARAMETER `(i)` ON LINE 163?
	   --------------------------------------------
	   This is CRITICAL to avoid a closure variable capture bug!

	   ❌ WRONG (without parameter):
	   for i := 0; i < workers; i++ {
	       go func() {
	           fmt.Println(i)  // BUG: All goroutines print the SAME value!
	       }()
	   }
	   Problem: All goroutines capture the SAME variable `i` by reference.
	           By the time goroutines run, loop has finished, i = workers.
	           All goroutines see i = workers (not 0, 1, 2, ...).

	   ✅ CORRECT (with parameter):
	   for i := 0; i < workers; i++ {
	       go func(workerID int) {
	           fmt.Println(workerID)  // Each goroutine gets its own copy!
	       }(i)  // Pass i as argument
	   }
	   Solution: Each goroutine receives its own COPY of i's value.
	           Goroutine 0 gets workerID=0, goroutine 1 gets workerID=1, etc.

	   Memory: Each goroutine's closure captures the workerID parameter
	           (8 bytes for int). Without the parameter, all goroutines
	           would share a reference to the same variable (race condition).

	   VARIABLE MORPHING THROUGH THE LOOP:
	   -----------------------------------
	   Iteration 0: i=0 → spawn goroutine with workerID=0
	   Iteration 1: i=1 → spawn goroutine with workerID=1
	   Iteration 2: i=2 → spawn goroutine with workerID=2
	   ...
	   Iteration N-1: i=workers-1 → spawn goroutine with workerID=workers-1

	   Each goroutine has its own independent copy of workerID.
	*/
	for i := 0; i < workers; i++ {
		wg.Add(1) // Increment counter BEFORE starting goroutine

		/*
		   GOROUTINE CREATION
		   ------------------
		   `go func(...) { ... }()` creates a new goroutine.

		   What happens:
		   1. Go runtime allocates ~2KB stack for the goroutine
		   2. Goroutine is scheduled to run (may start immediately or later)
		   3. Main goroutine continues to next iteration (doesn't wait)
		   4. Multiple goroutines run concurrently (on different OS threads)

		   Memory allocation:
		   - Initial stack: ~2KB per goroutine
		   - Stack grows automatically if needed (up to 1GB default)
		   - Heap allocations (maps, strings) shared across goroutines
		*/
		go func(workerID int) {
			/*
			   DEFER STATEMENT
			   ----------------
			   defer schedules a function call to run when the enclosing
			   function returns, regardless of how it returns:
			   - Normal return
			   - Early return (error case)
			   - Panic

			   Here: defer wg.Done() ensures the WaitGroup counter is
			         decremented when this worker goroutine exits.

			   Memory: defer uses a linked list of deferred calls (stored on stack).
			           Very efficient, no heap allocation for simple defers.
			*/
			defer wg.Done()

			/*
			   WORKER MAIN LOOP - HOW ITERATION WORKS
			   ---------------------------------------

			   This is an INFINITE loop (for { ... }) that continues until:
			     1. Context is cancelled (someone called cancel())
			     2. Jobs channel is closed (no more work)

			   Each iteration:
			     1. Check if context cancelled: <-ctx.Done()
			     2. Try to receive a job: url, ok := <-jobs
			     3. Process the job: fetchAndCount()
			     4. Send results: results <- counts
			     5. Loop back to step 1

			   VARIABLE FLOW IN EACH ITERATION:
			   --------------------------------
			   Iteration 1:
			     url (string) ← jobs channel
			     ↓
			     counts (map[string]int), err (error) ← fetchAndCount(url)
			     ↓
			     if err: send to errCh, cancel(), return
			     if success: counts → results channel

			   Iteration 2:
			     (same flow, different url)

			   Iteration N:
			     url, ok := <-jobs
			     if !ok: channel closed, return (exit loop)
			*/
			for {
				/*
				   SELECT STATEMENT - MULTIPLEXING CHANNEL OPERATIONS
				   --------------------------------------------------

				   select is like a "switch" for channels.
				   It waits until ONE of the cases can proceed, then executes it.

				   Behavior:
				   - If multiple cases are ready, ONE is chosen randomly
				   - If NO case is ready, select BLOCKS until one becomes ready
				   - If a case is ready, it executes immediately

				   Here we're checking TWO things:
				     1. Is context cancelled? (high priority)
				     2. Is there a job available? (normal operation)
				*/
				select {

				/*
				   UNDERSTANDING <-ctx.Done()
				   --------------------------

				   ctx.Done() returns a channel.
				   When context is cancelled:
				     - This channel is CLOSED
				     - Reading from it (<-ctx.Done()) returns immediately
				     - The zero value (nil) is returned

				   Before cancellation:
				     - Reading from ctx.Done() BLOCKS (waits)

				   This is how goroutines detect cancellation:
				     select {
				     case <-ctx.Done():
				         // Context cancelled! Stop working.
				         return
				     case work := <-workQueue:
				         // Got work, process it
				     }

				   Memory: ctx.Done() returns a channel that's shared across
				           all goroutines using the same context. Very efficient.
				*/
				case <-ctx.Done():
					// Context cancelled → exit immediately
					// This unblocks the select, and we return from the goroutine
					return

				/*
				   RECEIVING FROM A CHANNEL
				   ------------------------

				   url, ok := <-jobs

				   This syntax:
				     - url: The value received from the channel
				     - ok: Boolean indicating if channel is still open

				   ok values:
				     - true: Successfully received a value, channel is open
				     - false: Channel is closed and empty (no more values)

				   What happens:
				     - If jobs channel has a value: receive it immediately, ok=true
				     - If jobs channel is empty but open: BLOCK until value arrives
				     - If jobs channel is closed: return immediately, ok=false

				   Memory: Receiving from channel copies the value.
				           For strings, this copies the string header (~16 bytes),
				           not the underlying bytes (strings are immutable).
				*/
				case url, ok := <-jobs:
					if !ok {
						// Channel closed → no more work
						// This is how we signal "all jobs sent, workers can exit"
						return
					}

					/*
					   VARIABLE OWNERSHIP
					   ------------------
					   At this point, this worker "owns" the url.
					   No other worker will receive this same URL.
					   Channels guarantee that each value is delivered to exactly one receiver.
					*/

					/*
					   CALLING fetchAndCount - FUNCTION CALL STACK
					   -------------------------------------------

					   When we call fetchAndCount(ctx, url):
					     1. Go creates a new stack frame for fetchAndCount
					     2. Parameters are copied: ctx (context), url (string)
					     3. fetchAndCount executes (may take seconds for HTTP request)
					     4. Returns: counts (map[string]int), err (error)
					     5. Stack frame is popped, control returns here

					   Memory: Each function call uses stack space (~few KB).
					           Maps are allocated on heap (returned by reference).
					           HTTP response body is read into heap memory.

					   Variable transformation:
					     url (string) → HTTP request → response body ([]byte)
					     → tokenizeAndCount() → counts (map[string]int)
					*/
					counts, err := fetchAndCount(ctx, url)
					if err != nil {
						/*
						   ERROR HANDLING WITH NON-BLOCKING SEND
						   -------------------------------------

						   We want to send the error to errCh, but:
						   - Multiple workers might error simultaneously
						   - We only care about the FIRST error
						   - We don't want to block if errCh is full

						   Solution: Non-blocking send using select with default

						   How it works:
						     select {
						     case errCh <- err:
						         // Successfully sent error
						         cancel()  // Cancel context to stop other workers
						     default:
						         // errCh is full (already has an error)
						         // Ignore this error, first error already recorded
						     }

						   Why non-blocking?
						   - If errCh already has an error, we don't need to wait
						   - We want to exit immediately and let other workers stop
						   - Prevents deadlock if error handler isn't reading

						   Memory: Error is formatted into a string (heap allocation).
						           Only one error is kept in errCh buffer.
						*/
						select {
						case errCh <- fmt.Errorf("fetching %s: %w", url, err):
							// Successfully sent error → cancel context
							cancel() // This closes ctx.Done() channel, unblocking all workers
						default:
							// errCh already has an error, ignore this one
						}
						return // Exit this worker goroutine
					}

					/*
					   SENDING RESULTS - INTERRUPTIBLE OPERATION
					   ----------------------------------------

					   We want to send results, but also check if context was cancelled.
					   Why? Because cancellation might have happened while we were
					   processing the URL (between receiving job and sending result).

					   Pattern: Check cancellation AND send result atomically

					   select {
					   case <-ctx.Done():
					       // Cancelled while trying to send → exit
					   case results <- counts:
					       // Successfully sent result
					   }

					   Memory: Sending map to channel copies the map header (~8 bytes).
					           The map's underlying hash table remains on heap.
					           Multiple goroutines can reference the same map safely
					           if only reading, but here each worker creates its own map.
					*/
					select {
					case <-ctx.Done():
						// Context cancelled while processing → exit
						return
					case results <- counts:
						// Successfully sent result → continue to next iteration
					}
				}
			}
		}(i) // ← THIS IS LINE 163: Passing i as parameter to avoid closure bug
	}

	/*
	   JOB PRODUCER GOROUTINE - HOW IT SENDS WORK
	   ------------------------------------------

	   This goroutine feeds URLs into the jobs channel.

	   Why a separate goroutine?
	   --------------------------
	   - If we sent jobs in main goroutine, we'd block until all jobs are sent
	   - Workers couldn't start processing until all URLs are queued
	   - By using a goroutine, workers can start processing while URLs are still being sent

	   HOW ITERATION WORKS:
	   -------------------
	   for _, url := range urls {
	       // Iterates over each URL in the slice
	       // url variable is a COPY of each URL string
	       // Each iteration: url = urls[0], then urls[1], then urls[2], ...
	   }

	   VARIABLE FLOW:
	   --------------
	   urls[0] → url → jobs channel → worker receives it
	   urls[1] → url → jobs channel → worker receives it
	   urls[2] → url → jobs channel → worker receives it
	   ...

	   Each URL is sent to exactly one worker (channels guarantee this).

	   MEMORY MANAGEMENT:
	   -----------------
	   - Each url is a string (header ~16 bytes, data may be shared)
	   - Sending to channel copies the string header
	   - Channel buffer holds up to `workers` URLs
	   - Once buffer is full, sender blocks until worker receives

	   CLOSING THE CHANNEL:
	   --------------------
	   close(jobs) signals to all workers: "no more jobs coming"
	   - Workers receive ok=false when channel is empty and closed
	   - This is how workers know to exit their loops
	   - Only the sender should close the channel (rule of thumb)
	*/
	go func() {
		/*
		   RANGE LOOP OVER SLICE
		   ---------------------
		   for _, url := range urls

		   This iterates over each element in the urls slice.
		   - _ means "ignore the index"
		   - url is the value at each position
		   - Each iteration: url = urls[0], urls[1], urls[2], ...

		   Memory: url is a copy of the string header (~16 bytes).
		           The underlying string data may be shared (Go's string interning).
		*/
		for _, url := range urls {
			/*
			   INTERRUPTIBLE SEND
			   -------------------
			   We check context cancellation before sending each URL.
			   If context is cancelled (error occurred), we stop sending jobs.

			   Why check here?
			   - If a worker errors, we cancel context
			   - This goroutine should stop sending new jobs immediately
			   - Prevents wasting time sending URLs that won't be processed
			*/
			select {
			case <-ctx.Done():
				// Context cancelled → stop sending jobs
				return
			case jobs <- url:
				// Successfully sent URL to jobs channel
				// If buffer is full, this blocks until a worker receives
			}
		}
		/*
		   CLOSING THE CHANNEL
		   -------------------
		   close(jobs) tells all workers: "no more URLs will be sent"

		   What happens:
		     1. All pending sends complete
		     2. Channel is marked as closed
		     3. Future sends will panic
		     4. Receivers get ok=false when channel is empty

		   This is the signal for workers to exit their loops.
		*/
		close(jobs) // Signal: no more jobs coming
	}()

	/*
	   RESULTS CLOSER GOROUTINE - COORDINATION PATTERN
	   -----------------------------------------------

	   This goroutine coordinates the shutdown of the results channel.

	   Why is this needed?
	   -------------------
	   - The aggregator loop uses `for counts := range results`
	   - `range` over a channel continues until channel is closed
	   - We can't close results until ALL workers are done sending
	   - Workers finish at different times (some URLs take longer)

	   HOW IT WORKS:
	   -------------
	   1. This goroutine calls wg.Wait()
	      - Blocks until WaitGroup counter reaches 0
	      - Counter reaches 0 when all workers call wg.Done()
	   2. Once all workers finish, close(results)
	      - Signals to aggregator: "no more results coming"
	      - Aggregator's range loop exits

	   TIMING DIAGRAM:
	   ---------------
	   Worker 1: [fetch URL1] [send result] [wg.Done()]
	   Worker 2: [fetch URL2] [send result] [wg.Done()]
	   Worker 3: [fetch URL3] [send result] [wg.Done()]
	                      ↓
	              [All workers done]
	                      ↓
	              [wg.Wait() returns]
	                      ↓
	              [close(results)]
	                      ↓
	              [Aggregator's range loop exits]

	   MEMORY: This goroutine uses ~2KB stack, waits efficiently using
	           OS-level synchronization primitives (no busy-waiting).

	   CRITICAL RULE:
	   --------------
	   Only ONE goroutine should close a channel.
	   If multiple goroutines try to close the same channel, you get a panic.
	   Here, only this goroutine closes results.
	*/
	go func() {
		/*
		   WAITGROUP.WAIT() - BLOCKING OPERATION
		   -------------------------------------
		   wg.Wait() blocks until the WaitGroup counter reaches 0.

		   How it works:
		     - Internally uses a semaphore (OS-level synchronization)
		     - Efficiently blocks the goroutine (no CPU spinning)
		     - Wakes up when counter reaches 0

		   Memory: Wait() uses a small amount of stack space for the wait operation.
		           The goroutine is descheduled by the Go runtime (doesn't consume CPU).
		*/
		wg.Wait() // Block until all workers finish

		/*
		   CLOSING RESULTS CHANNEL
		   -----------------------
		   After all workers finish, we close the results channel.
		   This signals to the aggregator that no more results will arrive.
		*/
		close(results) // Signal: no more results coming
	}()

	/*
	   AGGREGATION (MAIN GOROUTINE) - MERGING RESULTS
	   -----------------------------------------------

	   This runs in the MAIN goroutine (not a worker).
	   It collects and merges word counts from all workers.

	   WHY NO LOCKS NEEDED:
	   --------------------
	   - Only ONE goroutine (main) writes to finalCounts
	   - Workers send results through channel (synchronization built-in)
	   - Channels provide thread-safe communication
	   - No shared mutable state → no race conditions

	   HOW ITERATION WORKS:
	   --------------------
	   for counts := range results {
	       // This loop continues until results channel is closed
	       // Each iteration receives ONE map[string]int from a worker
	       // counts variable is a map reference (maps are reference types)

	       for word, count := range counts {
	           // Iterate over each word-count pair in the map
	           // word: string (the word)
	           // count: int (frequency of that word)

	           finalCounts[word] += count
	           // Add this word's count to the final aggregate
	       }
	   }

	   VARIABLE MORPHING:
	   ------------------
	   Iteration 1:
	     counts = {"hello": 2, "world": 1}  (from worker 1)
	     finalCounts = {"hello": 2, "world": 1}

	   Iteration 2:
	     counts = {"hello": 1, "go": 3}     (from worker 2)
	     finalCounts = {"hello": 3, "world": 1, "go": 3}  (merged)

	   Iteration 3:
	     counts = {"world": 2, "go": 1}     (from worker 3)
	     finalCounts = {"hello": 3, "world": 3, "go": 4}  (merged)

	   ... continues until results channel is closed ...

	   MEMORY MANAGEMENT:
	   ------------------
	   - finalCounts: One map that grows as words are added
	   - Maps use hash tables internally (O(1) average insert/lookup)
	   - Map grows dynamically (rehashes when load factor > 6.5)
	   - Memory = O(unique_words_across_all_URLs)

	   - counts: Each iteration, counts is a reference to a worker's map
	   - Worker's map remains on heap until garbage collected
	   - We're copying counts INTO finalCounts, not sharing the map

	   MAP ACCUMULATION PATTERN:
	   -------------------------
	   finalCounts[word] += count

	   How this works:
	     - If word doesn't exist in finalCounts: Go creates entry with zero value (0)
	     - Then adds count to it: 0 + count = count
	     - If word exists: Adds count to existing value

	   Example:
	     finalCounts["hello"] = 2  (from previous iteration)
	     finalCounts["hello"] += 1  (add 1 more)
	     finalCounts["hello"] = 3   (result)
	*/
	finalCounts := make(map[string]int)

	/*
	   RANGE OVER CHANNEL
	   -------------------
	   for counts := range results

	   This is a special Go syntax for iterating over channel values.

	   How it works:
	     1. Receives values from results channel
	     2. Assigns each value to counts variable
	     3. Executes loop body
	     4. Repeats until channel is closed

	   When does it exit?
	     - When results channel is closed AND empty
	     - The closer goroutine closes results after all workers finish

	   Memory: Each iteration receives a map reference (~8 bytes).
	           The map's data remains on heap (maps are reference types).
	*/
	for counts := range results {
		/*
		   RANGE OVER MAP
		   --------------
		   for word, count := range counts

		   Iterates over key-value pairs in the map.
		   - word: the key (string)
		   - count: the value (int)

		   Order: Maps in Go have RANDOM iteration order (by design).
		*/
		for word, count := range counts {
			/*
			   MAP ACCUMULATION
			   ----------------
			   finalCounts[word] += count

			   This adds count to the existing value for word.
			   If word doesn't exist, Go uses zero value (0 for int).
			*/
			finalCounts[word] += count
		}
	}

	/*
	   ERROR CHECK - NON-BLOCKING RECEIVE
	   ----------------------------------

	   At this point:
	     - All workers are done (wg.Wait() completed)
	     - All results have been merged (range loop exited)
	     - If an error occurred, it MUST be in errCh

	   WHY NON-BLOCKING RECEIVE?
	   -------------------------
	   We use select with default to check for error without blocking.

	   If we used blocking receive:
	     err := <-errCh  // Would block forever if no error occurred!

	   With non-blocking:
	     select {
	     case err := <-errCh:
	         // Error exists → return it
	     default:
	         // No error → continue to success return
	     }

	   MEMORY: Error is an interface type (~16 bytes).
	           If present, it contains a pointer to the error string.
	*/
	select {
	case err := <-errCh:
		// Error occurred → return it (cancellation already happened)
		return nil, err
	default:
		// No error → proceed to return success
	}

	/*
	   SUCCESS RETURN
	   --------------
	   Return the aggregated word counts.

	   Memory: finalCounts map is returned by value, but maps are reference types,
	           so only the map header (~8 bytes) is copied, not the entire hash table.
	           The caller gets a reference to the same underlying data structure.
	*/
	return finalCounts, nil
}

/*
fetchAndCount performs:
 1. Context-aware HTTP GET
 2. Full body read
 3. Tokenization

FUNCTION CALL STACK AND VARIABLE FLOW:
--------------------------------------
Input: ctx (context.Context), url (string)

	↓

http.NewRequestWithContext(ctx, ...)

	↓

req (*http.Request) - contains URL, method, context

	↓

http.DefaultClient.Do(req)

	↓

resp (*http.Response) - contains status, headers, body

	↓

io.ReadAll(resp.Body)

	↓

body ([]byte) - entire response body in memory

	↓

string(body) - convert bytes to string

	↓

tokenizeAndCount(string)

	↓

counts (map[string]int) - word frequencies

	↓

Return: counts, error

MEMORY ALLOCATIONS:
------------------
1. req: ~200 bytes (HTTP request struct)
2. resp: ~500 bytes (HTTP response struct)
3. body: O(response_size) - entire body read into memory
4. string(body): O(response_size) - copy of body as string
5. counts map: O(vocabulary_size) - hash table for word counts

Total per URL: ~700 bytes + 2 * response_size + vocabulary_size

CONTEXT PROPAGATION:
-------------------
Passing ctx into NewRequestWithContext ensures that:
  - If cancel() is called, the HTTP request is aborted immediately
  - The underlying TCP connection is closed
  - No time wasted waiting for slow/failed requests
  - Resources are freed promptly

Key idea:
---------
Context cancellation propagates through the HTTP client stack.
When ctx.Done() is closed, the HTTP transport layer detects it and
aborts the connection, even if the request is in-flight.
*/
func fetchAndCount(ctx context.Context, url string) (map[string]int, error) {
	/*
	   CREATING HTTP REQUEST WITH CONTEXT
	   -----------------------------------

	   http.NewRequestWithContext(ctx, method, url, body)

	   Parameters:
	     - ctx: Context for cancellation/timeout
	     - http.MethodGet: HTTP method ("GET")
	     - url: The URL to fetch
	     - nil: Request body (GET requests don't have bodies)

	   Returns:
	     - req: *http.Request (pointer to request struct)
	     - err: error if URL is invalid

	   Memory: Request struct is ~200 bytes, allocated on heap.
	           Contains URL string, headers map, context reference.
	*/
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	/*
	   EXECUTING HTTP REQUEST
	   ----------------------

	   http.DefaultClient.Do(req)

	   What happens:
	     1. HTTP client opens TCP connection to server
	     2. Sends HTTP request
	     3. Waits for response (may take seconds)
	     4. Reads response headers
	     5. Returns response with body stream

	   If ctx is cancelled during this:
	     - TCP connection is closed
	     - Do() returns immediately with error
	     - No time wasted waiting

	   Memory: Response struct is ~500 bytes.
	           Body is a stream (io.ReadCloser), not yet read into memory.
	*/
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	/*
	   DEFER STATEMENT FOR CLEANUP
	   -----------------------------

	   defer resp.Body.Close()

	   Ensures the response body is closed when function returns.
	   This is CRITICAL:
	     - Releases TCP connection back to pool
	     - Frees file descriptors
	     - Prevents resource leaks

	   Memory: Closing releases the connection, freeing ~1KB of resources.
	*/
	defer resp.Body.Close()

	/*
	   STATUS CODE CHECK
	   -----------------

	   Only accept HTTP 200 OK responses.
	   Other status codes (404, 500, etc.) are treated as errors.
	*/
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	/*
	   READING RESPONSE BODY
	   ---------------------

	   io.ReadAll(resp.Body)

	   Reads the ENTIRE response body into memory.

	   What happens:
	     1. Allocates buffer (starts small, grows as needed)
	     2. Reads all bytes from resp.Body stream
	     3. Returns []byte containing entire body

	   Memory: Allocates O(response_size) bytes.
	           For large responses (MBs), this can be significant.
	           Alternative: Use bufio.Scanner for streaming (see stretch goals).

	   Variable transformation:
	     resp.Body (io.ReadCloser stream) → body ([]byte)
	*/
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	/*
	   CONVERTING BYTES TO STRING
	   --------------------------

	   string(body)

	   Converts []byte to string.

	   Memory: In Go, strings are immutable.
	           This conversion may:
	             - Share underlying bytes if body is ASCII (no copy)
	             - Copy bytes if body contains non-ASCII (UTF-8 conversion)
	           Typically: O(response_size) memory allocation.

	   Then tokenizeAndCount processes the string and returns a map.
	*/
	return tokenizeAndCount(string(body)), nil
}

/*
tokenizeAndCount is CPU-bound, deterministic, and side-effect-free.

TOKENIZATION - WHAT IT MEANS:
-----------------------------
Tokenization is the process of breaking text into individual words (tokens).
Example: "Hello, world!" → ["hello", "world"]

This function:
 1. Splits text into words (by whitespace)
 2. Normalizes words (lowercase, remove punctuation)
 3. Counts frequency of each word

WHY SAFE FOR CONCURRENCY:
-------------------------
- Pure function: No side effects (doesn't modify global state)
- Deterministic: Same input → same output
- No shared mutable state: Each call creates its own map
- Thread-safe: Can run in multiple goroutines simultaneously

MEMORY ALLOCATIONS:
------------------
1. counts map: O(vocabulary_size) - hash table
2. strings.Fields: O(text_length) - creates slice of strings
3. strings.ToLower: May allocate new string if changes needed
4. strings.Map: Allocates new string for each word

Total: O(text_length + vocabulary_size)

VARIABLE TRANSFORMATION THROUGH ITERATION:
-----------------------------------------
Input: text = "Hello, world! Go is great."

Iteration 1:

	word = "Hello," (from strings.Fields)
	word = "hello," (after ToLower)
	word = "hello" (after Map removes comma)
	counts["hello"] = 1

Iteration 2:

	word = "world!"
	word = "world!"
	word = "world"
	counts["world"] = 1

Iteration 3:

	word = "Go"
	word = "go"
	word = "go"
	counts["go"] = 1

Iteration 4:

	word = "is"
	word = "is"
	word = "is"
	counts["is"] = 1

Iteration 5:

	word = "great."
	word = "great."
	word = "great"
	counts["great"] = 1

Result: counts = {"hello": 1, "world": 1, "go": 1, "is": 1, "great": 1}
*/
func tokenizeAndCount(text string) map[string]int {
	/*
	   MAP INITIALIZATION
	   ------------------

	   make(map[string]int)

	   Creates an empty map with string keys and int values.

	   Memory: Map starts with small hash table (~48 bytes).
	           Grows automatically as entries are added.
	           Each entry: ~8 bytes (key pointer) + ~8 bytes (value) + overhead.
	*/
	counts := make(map[string]int)

	/*
	   STRINGS.FIELDS - SPLITTING TEXT INTO WORDS
	   -------------------------------------------

	   strings.Fields(text)

	   Splits text on any Unicode whitespace character:
	     - Space (' ')
	     - Tab ('\t')
	     - Newline ('\n')
	     - Other Unicode whitespace

	   Returns: []string (slice of strings)

	   Example:
	     "hi,\nthere\tfriend" → ["hi,", "there", "friend"]

	   Memory: Allocates slice with pointers to string headers.
	           Strings may share underlying bytes (Go's string interning).

	   HOW ITERATION WORKS:
	   -------------------
	   for _, word := range strings.Fields(text)

	   This iterates over each word in the slice.
	   - _ means "ignore the index"
	   - word is the string value at each position
	   - Each iteration processes one word
	*/
	for _, word := range strings.Fields(text) {
		/*
		   NORMALIZING TO LOWERCASE
		   ------------------------

		   strings.ToLower(word)

		   Converts all characters to lowercase.

		   Memory: May allocate new string if word contains uppercase letters.
		           If word is already lowercase, may return same string (optimization).

		   Example:
		     "Hello" → "hello"
		     "WORLD" → "world"
		     "hello" → "hello" (may reuse same string)
		*/
		word = strings.ToLower(word)

		/*
		   REMOVING NON-LETTER CHARACTERS
		   ------------------------------

		   strings.Map(func(r rune) rune, word)

		   Applies a function to each rune (character) in the string.

		   Function signature: func(rune) rune
		     - Input: one character (rune)
		     - Output: replacement character (or -1 to delete)

		   How it works:
		     - Iterates over each character in word
		     - Calls function with that character
		     - If function returns the character: keeps it
		     - If function returns -1: deletes it

		   Example:
		     word = "hello,"
		     'h' → unicode.IsLetter('h') = true → return 'h' → keep
		     'e' → unicode.IsLetter('e') = true → return 'e' → keep
		     'l' → unicode.IsLetter('l') = true → return 'l' → keep
		     'l' → unicode.IsLetter('l') = true → return 'l' → keep
		     'o' → unicode.IsLetter('o') = true → return 'o' → keep
		     ',' → unicode.IsLetter(',') = false → return -1 → delete
		     Result: "hello"

		   Memory: Allocates new string with filtered characters.
		           O(word_length) allocation.

		   Why not use regex?
		   ------------------
		   - Regex is slower (compilation, matching overhead)
		   - strings.Map handles Unicode correctly (runes, not bytes)
		   - More efficient for simple character filtering
		*/
		word = strings.Map(func(r rune) rune {
			/*
			   UNICODE.ISLETTER
			   ----------------

			   Checks if a rune is a Unicode letter character.
			   Handles all languages: English, Chinese, Arabic, etc.

			   Returns:
			     - true: Character is a letter → keep it
			     - false: Character is not a letter → delete it (-1)
			*/
			if unicode.IsLetter(r) {
				return r // Keep this character
			}
			return -1 // Delete this character
		}, word)

		/*
		   SKIPPING EMPTY WORDS
		   --------------------

		   After removing punctuation, some "words" might be empty.
		   Example: "123" → "" (all digits removed)

		   We skip empty words because they're not meaningful.
		*/
		if word == "" {
			continue // Skip to next iteration
		}

		/*
		   COUNTING WORDS
		   --------------

		   counts[word]++

		   This increments the count for this word.

		   How it works:
		     - If word doesn't exist in map: Go creates entry with zero value (0), then increments → 1
		     - If word exists: Increments existing value

		   Example:
		     counts["hello"]++  // First time: counts["hello"] = 1
		     counts["hello"]++  // Second time: counts["hello"] = 2

		   Memory: Map grows as new words are added.
		           Hash table rehashes when load factor > 6.5 (automatic).
		*/
		counts[word]++
	}

	/*
	   RETURNING THE MAP
	   -----------------

	   Maps in Go are reference types.
	   Returning a map returns a reference to the underlying hash table.
	   The map remains on the heap until garbage collected.
	*/
	return counts
}

/*
============================================================================
ERRGROUP IMPLEMENTATION: A SAFER, MORE CONCISE ALTERNATIVE
============================================================================

This section demonstrates how errgroup simplifies the WordCount function
by combining WaitGroup, Context cancellation, and error handling into
a single, safer abstraction.

HOW ERRGROUP WORKS UNDER THE HOOD:
-----------------------------------

errgroup.Group is essentially a wrapper that combines three things:

1. sync.WaitGroup - Tracks goroutine completion
2. context.Context - Provides cancellation
3. Error handling - Captures first error and cancels context

Internal structure (simplified):
---------------------------------
type Group struct {
    cancel func()              // Function to cancel context
    wg     sync.WaitGroup     // Tracks active goroutines
    errOnce sync.Once         // Ensures only first error stored
    err    error              // First error encountered
    errMu  sync.Mutex         // Protects err field (thread-safe)
}

When you call g.Go(func() error { ... }):
------------------------------------------
1. g.wg.Add(1) is called automatically
2. A goroutine is launched
3. The function runs
4. If function returns error:
   a. sync.Once ensures only first error is stored
   b. errMu.Lock() protects err field
   c. err is stored
   d. cancel() is called → context cancelled
5. defer g.wg.Done() is called automatically

When you call g.Wait():
-----------------------
1. g.wg.Wait() blocks until all goroutines finish
2. g.cancel() is called (cleanup)
3. g.err is returned (first error, or nil)

WHY ERRGROUP IS SAFER:
----------------------
1. No manual WaitGroup management (can't forget Add/Done)
2. Automatic context cancellation on error
3. Thread-safe error handling (sync.Once + Mutex)
4. Always cleans up context (even on success)
5. Less code = fewer bugs

WHY ERRGROUP IS MORE EFFICIENT:
--------------------------------
1. No error channel needed (saves memory)
2. Built-in error handling (less code)
3. Automatic cleanup (no defer cancel() needed)
4. Single abstraction (easier to reason about)
*/

/*
WordCountWithErrGroup demonstrates how errgroup simplifies the worker pool pattern.

COMPARISON WITH MANUAL APPROACH:
---------------------------------

Manual Approach Problems:
- Manual WaitGroup: wg.Add(1), defer wg.Done(), wg.Wait()
- Manual error channel: errCh := make(chan error, 1)
- Manual error handling: select { case errCh <- err: cancel() default: }
- Manual context cleanup: defer cancel()
- ~100+ lines of coordination code

errgroup Approach Benefits:
- Automatic WaitGroup: g.Go() handles Add/Done internally
- No error channel: Just return error from goroutine
- Automatic error handling: errgroup captures first error
- Automatic context cleanup: g.Wait() cancels context
- ~30% less code, fewer bugs

HOW IT WORKS:
-------------

1. Create errgroup: g, ctx := errgroup.WithContext(ctx)
  - Creates cancellable context
  - Creates Group with WaitGroup + error handling

2. Launch workers: g.Go(func() error { ... })
  - Automatically calls wg.Add(1)
  - Runs function in goroutine
  - If function returns error, cancels context automatically
  - Automatically calls wg.Done() when done

3. Wait for completion: g.Wait()
  - Blocks until all goroutines finish
  - Cancels context (cleanup)
  - Returns first error (if any)

MEMORY COMPARISON:
------------------
Manual:
  - errCh: 16 bytes
  - WaitGroup: 12 bytes
  - Context management: 48 bytes
  - Total: ~76 bytes + coordination code

errgroup:
  - Group struct: ~64 bytes (includes everything)
  - Total: ~64 bytes, simpler code

Note: For worker pools, we still need channels for job distribution.

	errgroup is best for fan-out patterns, but can be adapted for pools.
*/
func WordCountWithErrGroup(ctx context.Context, urls []string, workers int) (map[string]int, error) {
	/*
		STEP 1: CREATE ERRGROUP WITH CONTEXT
		------------------------------------

		errgroup.WithContext(ctx) does two things:
		1. Creates a cancellable context: ctx, cancel := context.WithCancel(ctx)
		2. Creates a Group struct with WaitGroup + error handling

		What you get back:
		- g: *errgroup.Group (manages goroutines and errors)
		- ctx: context.Context (cancellable context for all operations)

		Memory: Group struct is ~64 bytes (WaitGroup + mutexes + function pointer)

		Key difference from manual approach:
		- No need for: ctx, cancel := context.WithCancel(ctx); defer cancel()
		- errgroup handles context creation and cleanup automatically
	*/
	g, ctx := errgroup.WithContext(ctx)

	/*
		STEP 2: CREATE CHANNELS (STILL NEEDED FOR WORKER POOL)
		-------------------------------------------------------

		Even with errgroup, we still need channels for:
		- jobs: Distribute URLs to workers
		- results: Collect word counts from workers

		Why still needed?
		- Worker pool pattern requires job distribution
		- errgroup is better for fan-out (one goroutine per task)
		- But we can adapt it for worker pools

		Memory: Same as manual approach
		- jobs: workers * 16-32 bytes
		- results: workers * 8 bytes
	*/
	jobs := make(chan string, workers)
	results := make(chan map[string]int, workers)

	/*
		STEP 3: START WORKERS USING ERRGROUP
		-------------------------------------

		KEY DIFFERENCE FROM MANUAL APPROACH:
		- No manual wg.Add(1) - errgroup does it automatically
		- No defer wg.Done() - errgroup does it automatically
		- Just return error - errgroup handles cancellation

		How g.Go() works internally:
		1. g.wg.Add(1) ← Automatic!
		2. go func() {
		3.     defer g.wg.Done() ← Automatic!
		4.     if err := f(); err != nil {
		5.         g.errOnce.Do(func() {
		6.             g.errMu.Lock()
		7.             g.err = err
		8.             g.errMu.Unlock()
		9.             g.cancel() ← Automatic cancellation!
		10.         })
		11.     }
		12. }()

		Benefits:
		- Can't forget wg.Add(1) or wg.Done()
		- Error automatically cancels context
		- Thread-safe error handling built-in
	*/
	for i := 0; i < workers; i++ {
		/*
			CAPTURE LOOP VARIABLE
			----------------------
			Still need to pass i as parameter to avoid closure bug.
			errgroup doesn't change this requirement.
		*/
		g.Go(func() error {
			/*
				WORKER LOOP - SIMPLIFIED ERROR HANDLING
				---------------------------------------

				Key simplification:
				- No manual errCh channel
				- Just return error - errgroup handles it
				- Context cancellation is automatic

				Before (manual):
				  if err != nil {
				      select {
				      case errCh <- err:
				          cancel()
				      default:
				      }
				      return
				  }

				After (errgroup):
				  if err != nil {
				      return err  // That's it! errgroup handles cancellation
				  }
			*/
			for {
				select {
				case <-ctx.Done():
					/*
						CONTEXT CANCELLATION CHECK
						-------------------------
						If context cancelled (error occurred), exit.
						errgroup automatically cancelled context when first error occurred.
					*/
					return ctx.Err()

				case url, ok := <-jobs:
					if !ok {
						// Channel closed → no more work
						return nil // Success (no error)
					}

					/*
						PROCESS JOB
						-----------
						Fetch URL and count words.
						If error occurs, just return it - errgroup handles cancellation.
					*/
					counts, err := fetchAndCount(ctx, url)
					if err != nil {
						/*
							ERROR HANDLING - SIMPLIFIED!
							---------------------------
							Just return the error. errgroup will:
							1. Store it (first error only, thread-safe)
							2. Cancel context automatically
							3. Other workers will see ctx.Done() and exit

							No need for:
							- errCh channel
							- select with default
							- manual cancel() call
						*/
						return fmt.Errorf("fetching %s: %w", url, err)
					}

					/*
						SEND RESULTS
						------------
						Send word counts to results channel.
						Check context in case cancellation happened during processing.
					*/
					select {
					case <-ctx.Done():
						return ctx.Err()
					case results <- counts:
						// Successfully sent result
					}
				}
			}
		})
	}

	/*
		STEP 4: SEND JOBS (SAME AS MANUAL APPROACH)
		-------------------------------------------

		Still need separate goroutine to send jobs.
		errgroup doesn't change this requirement for worker pools.
	*/
	go func() {
		defer close(jobs)
		for _, url := range urls {
			select {
			case <-ctx.Done():
				return
			case jobs <- url:
			}
		}
	}()

	/*
		STEP 5: CLOSE RESULTS WHEN WORKERS DONE
		---------------------------------------

		Still need to close results channel when all workers finish.
		Use g.Wait() instead of wg.Wait().
	*/
	go func() {
		/*
			WAIT FOR COMPLETION
			-------------------
			g.Wait() does three things:
			1. Blocks until all goroutines finish (wg.Wait())
			2. Cancels context (cleanup)
			3. Returns first error (if any)

			We ignore the error here because we'll check it after aggregating results.
		*/
		_ = g.Wait() // Wait for workers, ignore error (check later)
		close(results)
	}()

	/*
		STEP 6: AGGREGATE RESULTS (SAME AS MANUAL APPROACH)
		---------------------------------------------------
	*/
	finalCounts := make(map[string]int)
	for counts := range results {
		for word, count := range counts {
			finalCounts[word] += count
		}
	}

	/*
		STEP 7: CHECK FOR ERRORS (SIMPLIFIED!)
		---------------------------------------

		KEY SIMPLIFICATION:
		- No select statement needed
		- No errCh channel
		- Just call g.Wait() and check error

		Before (manual):
		  select {
		  case err := <-errCh:
		      return nil, err
		  default:
		  }

		After (errgroup):
		  if err := g.Wait(); err != nil {
		      return nil, err
		  }

		Why this works:
		- g.Wait() returns the first error encountered
		- If no error, returns nil
		- Context already cancelled (cleanup done)
	*/
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return finalCounts, nil
}
