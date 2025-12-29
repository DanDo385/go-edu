# How the Worker Pool Wordcount Solution Works

This document explains in plain language how the worker pool wordcount solution is implemented. The solution uses Go's concurrency primitives to efficiently process multiple URLs concurrently while maintaining bounded resource usage.

## The Big Picture

The solution implements a **worker pool pattern** to fetch URLs and count words across their content. Instead of creating one goroutine per URL (which could overwhelm the system), we create a fixed number of worker goroutines that process URLs from a shared queue. This provides:

- **Bounded concurrency**: Only N workers run simultaneously
- **Resource efficiency**: Memory usage stays constant regardless of URL count
- **Error handling**: First error cancels all work automatically
- **Graceful shutdown**: Context cancellation stops all workers cleanly

## Architecture Overview

The solution consists of three main components:

1. **WordCount**: The main orchestrator that sets up workers and coordinates work
2. **fetchAndCount**: Fetches a single URL and counts words in its content
3. **tokenizeAndCount**: Splits text into words and counts their frequencies

## Component 1: WordCount - The Main Orchestrator

The `WordCount` function is the heart of the solution. It uses Go's `errgroup` package to simplify worker management and error handling.

### Step 1: Initialize errgroup and Channels

```go
g, ctx := errgroup.WithContext(ctx)
jobs := make(chan string, workers)
results := make(chan map[string]int, workers)
```

**What happens here:**
- `errgroup.WithContext` creates a cancellable context and an error group. The error group automatically manages worker lifecycle and propagates the first error.
- Two buffered channels are created:
  - `jobs`: Carries URLs to be processed (buffered to allow N URLs to be queued)
  - `results`: Carries word count maps from completed work (buffered for efficiency)

**Why buffered channels?** They allow senders to continue working without blocking immediately, improving throughput.

### Step 2: Launch Worker Goroutines

```go
for i := 0; i < workers; i++ {
    g.Go(func() error {
        // Worker loop
    })
}
```

**What happens here:**
- A fixed number of worker goroutines are launched (typically 3-10 workers)
- Each worker runs concurrently and processes URLs from the `jobs` channel
- `g.Go()` automatically handles WaitGroup operations internally - no manual `Add()`/`Done()` needed

**The Worker Loop:**

Each worker runs an infinite loop that:
1. **Checks for cancellation**: If the context is cancelled (timeout or error), the worker exits immediately
2. **Receives a job**: Reads a URL from the `jobs` channel
3. **Detects channel closure**: When `jobs` closes (no more URLs), the worker exits normally
4. **Processes the job**: Calls `fetchAndCount` to fetch the URL and count words
5. **Handles errors**: If fetching fails, returns an error (errgroup cancels all workers automatically)
6. **Sends results**: Sends the word count map to the `results` channel (with cancellation check)

**Key detail - Channel closure detection:**
```go
case url, ok := <-jobs:
    if !ok {
        return nil  // Channel closed, exit normally
    }
```
The two-value receive (`url, ok`) tells us if the channel is still open. When `ok` is false, the channel is closed and there are no more jobs.

### Step 3: Send Jobs in Background

```go
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
```

**What happens here:**
- A separate goroutine sends all URLs to the `jobs` channel
- It checks for cancellation before each send (allows early exit if context cancelled)
- When done, it closes the `jobs` channel (signals workers that no more jobs are coming)

**Why a separate goroutine?** The main goroutine needs to continue and start collecting results. If we sent jobs synchronously, we'd block until all URLs were sent.

### Step 4: Close Results Channel When Workers Finish

```go
go func() {
    g.Wait()
    close(results)
}()
```

**What happens here:**
- A goroutine waits for all workers to finish using `g.Wait()`
- Once all workers are done, it closes the `results` channel
- This signals the aggregation loop that no more results are coming

**Why a separate goroutine?** We need to start collecting results immediately while workers are still processing. This goroutine ensures the results channel closes only after all work completes.

### Step 5: Aggregate Results

```go
finalCounts := make(map[string]int)
for counts := range results {
    for word, count := range counts {
        finalCounts[word] += count
    }
}
```

**What happens here:**
- The main goroutine reads from the `results` channel
- For each word count map received, it merges the counts into `finalCounts`
- The loop exits when `results` channel closes (after all workers finish)

**Why range over channel?** The `range` loop automatically exits when the channel is closed and empty, making this pattern clean and idiomatic.

### Step 6: Check for Errors

```go
if err := g.Wait(); err != nil {
    return nil, err
}
return finalCounts, nil
```

**What happens here:**
- `g.Wait()` is called again to check if any worker returned an error
- If an error occurred, it returns the error (errgroup stores the first error)
- Otherwise, returns the aggregated word counts

**Why call Wait() twice?** The first call (in the goroutine) waits for workers to finish so we can close the results channel. The second call checks if any errors occurred. This is safe because `Wait()` is idempotent - calling it multiple times returns the same result.

## Component 2: fetchAndCount - Fetch and Process a URL

This function handles the HTTP request and delegates word counting to `tokenizeAndCount`.

### Step 1: Create HTTP Request with Context

```go
req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
```

**What happens here:**
- Creates an HTTP GET request with the cancellable context
- If the context is cancelled, the HTTP client will abort the request

### Step 2: Execute Request

```go
resp, err := http.DefaultClient.Do(req)
defer resp.Body.Close()
```

**What happens here:**
- Sends the HTTP request and waits for response
- Uses `defer` to ensure the response body is always closed (prevents resource leaks)
- If the context is cancelled during the request, `Do()` returns an error

### Step 3: Check Status Code

```go
if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
}
```

**What happens here:**
- Verifies the response was successful (HTTP 200)
- Returns an error for non-200 status codes (4xx, 5xx, etc.)

### Step 4: Read Response Body

```go
body, err := io.ReadAll(resp.Body)
```

**What happens here:**
- Reads the entire response body into memory
- This is necessary because we need the full text to count words

### Step 5: Tokenize and Count

```go
return tokenizeAndCount(string(body)), nil
```

**What happens here:**
- Converts the byte slice to a string
- Calls `tokenizeAndCount` to process the text
- Returns the word count map

## Component 3: tokenizeAndCount - Process Text

This function splits text into words, normalizes them, and counts frequencies.

### Step 1: Split into Words

```go
for _, word := range strings.Fields(text) {
```

**What happens here:**
- `strings.Fields()` splits text by whitespace (spaces, tabs, newlines)
- Returns a slice of words

### Step 2: Normalize Each Word

```go
word = strings.ToLower(word)
word = strings.Map(func(r rune) rune {
    if unicode.IsLetter(r) {
        return r
    }
    return -1
}, word)
```

**What happens here:**
- Converts to lowercase ("Hello" → "hello")
- Removes non-letter characters using `strings.Map`
  - For each character (rune), checks if it's a letter
  - Keeps letters, removes everything else (punctuation, numbers, etc.)
  - Returns -1 to delete a character

**Example:** "Hello, world!" → "hello" and "world"

### Step 3: Skip Empty Words

```go
if word == "" {
    continue
}
```

**What happens here:**
- After removing non-letters, a word might become empty (e.g., "123" → "")
- Skips empty words to avoid counting them

### Step 4: Count Words

```go
counts[word]++
```

**What happens here:**
- Increments the count for this word
- If the word doesn't exist in the map, Go creates it with value 0, then increments to 1
- This is a common Go idiom for counting

## Error Handling Flow

The solution uses `errgroup` for automatic error propagation:

1. **Worker encounters error**: Returns error from `g.Go()` function
2. **errgroup cancels context**: Automatically cancels the context when first error occurs
3. **All workers stop**: Workers checking `ctx.Done()` exit immediately
4. **Error stored**: errgroup stores the first error
5. **Main function checks**: `g.Wait()` returns the stored error
6. **Error returned**: Main function returns the error to caller

**Key benefit:** No manual error channel or cancellation logic needed. errgroup handles it all automatically.

## Concurrency Safety

The solution is safe for concurrent access:

- **Channels are thread-safe**: Multiple goroutines can safely send/receive
- **Maps are not shared**: Each worker creates its own word count map, then sends it through a channel
- **Aggregation is sequential**: The main goroutine aggregates results sequentially (no race conditions)
- **Context is safe**: Context cancellation is thread-safe and can be checked concurrently

## Performance Characteristics

- **Memory**: O(workers) - Only N workers exist regardless of URL count
- **Time**: O(total_words) - Must process all words, but does so concurrently
- **Network**: Bounded by worker count - Only N simultaneous HTTP requests
- **CPU**: Utilizes all available CPU cores efficiently

## Example Execution Flow

For 5 URLs and 3 workers:

1. **T=0ms**: 3 workers start, begin waiting for jobs
2. **T=0ms**: Job sender goroutine sends first 3 URLs to `jobs` channel
3. **T=0ms**: 3 workers receive URLs and start fetching concurrently
4. **T=100ms**: Worker 1 finishes, sends results, receives next URL
5. **T=150ms**: Worker 2 finishes, sends results, receives next URL
6. **T=200ms**: Worker 3 finishes, sends results, receives next URL
7. **T=250ms**: Worker 1 finishes last URL, exits (channel closed)
8. **T=300ms**: All workers done, results channel closes
9. **T=300ms**: Aggregation completes, function returns

**Key insight:** Workers process URLs concurrently, so total time is roughly `max(individual_url_times)` rather than `sum(individual_url_times)`.

## Summary

The solution implements a robust, efficient worker pool pattern using:

- **errgroup**: Simplifies worker management and error handling
- **Channels**: Type-safe communication between goroutines
- **Context**: Enables cancellation and timeout support
- **Bounded concurrency**: Prevents resource exhaustion
- **Graceful shutdown**: Clean exit on errors or cancellation

This pattern is widely used in production Go systems for processing work queues, web scraping, batch processing, and any scenario requiring concurrent processing with bounded resources.
