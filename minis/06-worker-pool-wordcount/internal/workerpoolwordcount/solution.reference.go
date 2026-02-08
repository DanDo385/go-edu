//go:build reference

package workerpoolwordcount

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"unicode"
)

/*
Reference Solution
==================

This file is the canonical reference for this exercise. It keeps failure paths
explicit when an operation can fail, so callers can decide how to handle
errors at API boundaries.

Read this alongside exercise.go and the tests to understand the intended data
flow, ownership boundaries, and invariants that keep behavior deterministic.

Teaching notes:
- Memory/ownership: make copies when returning mutable data that should not
  alias internal state; share references only when aliasing is intentional.
- Invariants: establish assumptions close to construction, and rely on them in
  smaller helper functions to keep logic easy to audit.
- Error surfaces: prefer explicit returns over hidden panics so learners can
  reason about control flow in production-style code.
*/

// WordCount is the main function you will implement.
func WordCount(ctx context.Context, urls []string, workers int) (map[string]int, error) {
	if workers <= 0 {
		workers = 1
	}

	type result struct {
		counts map[string]int
		err    error
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan string)
	results := make(chan result)

	var wg sync.WaitGroup
	workerFn := func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case url, ok := <-jobs:
				if !ok {
					return
				}
				counts, err := fetchAndCount(ctx, url)
				select {
				case results <- result{counts: counts, err: err}:
				case <-ctx.Done():
					return
				}
			}
		}
	}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go workerFn()
	}

	go func() {
		defer close(jobs)
		for _, url := range urls {
			select {
			case jobs <- url:
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	merged := make(map[string]int)
	for res := range results {
		if res.err != nil {
			cancel()
			return nil, res.err
		}
		for word, count := range res.counts {
			merged[word] += count
		}
	}

	if err := ctx.Err(); err != nil && err != context.Canceled {
		return nil, err
	}
	return merged, nil
}

// fetchAndCount fetches the content of a URL and counts the words.
func fetchAndCount(ctx context.Context, url string) (map[string]int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return tokenizeAndCount(string(body)), nil
}

// tokenizeAndCount takes a string, tokenizes it into words, and counts them.
func tokenizeAndCount(text string) map[string]int {
	counts := make(map[string]int)
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		token := normalizeToken(scanner.Text())
		if token == "" {
			continue
		}
		counts[token]++
	}
	return counts
}

// normalizeToken implements the reference behavior for this exercise.
//
// Algorithm steps:
// 1. Validate prerequisites and invariants before mutating state.
// 2. Execute the core operation while keeping ownership/aliasing explicit.
// 3. Return explicit values/errors so callers control failure behavior.
func normalizeToken(word string) string {
	var b strings.Builder
	for _, r := range word {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}
