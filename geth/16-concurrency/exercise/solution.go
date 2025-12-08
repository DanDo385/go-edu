//go:build solution
// +build solution

package exercise

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Run contains the reference solution for a bounded-concurrency probe runner.
//
// Concepts:
// - Context timeouts guard the whole run and each probe.
// - Worker pools (jobs/results channels + WaitGroup) avoid unbounded goroutines.
// - Safe aggregation using mutex-protected maps.
func Run(ctx context.Context, p Prober, cfg Config) (*Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if p == nil {
		return nil, errors.New("prober is nil")
	}
	if cfg.Workers <= 0 {
		cfg.Workers = 4
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 5 * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	jobs := make(chan string, cfg.Workers)
	var mu sync.Mutex
	res := &Result{
		Successes: make(map[string]time.Duration),
		Failures:  make(map[string]error),
	}

	var wg sync.WaitGroup
	for i := 0; i < cfg.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for endpoint := range jobs {
				start := time.Now()

				// Per-request timeout to avoid one slow endpoint stalling others.
				reqCtx, cancelReq := context.WithTimeout(ctx, cfg.Timeout/2)
				err := p.Probe(reqCtx, endpoint)
				cancelReq()

				latency := time.Since(start)
				mu.Lock()
				if err != nil {
					res.Failures[endpoint] = fmt.Errorf("probe failed: %w", err)
				} else {
					res.Successes[endpoint] = latency
				}
				mu.Unlock()
			}
		}()
	}

	go func() {
		defer close(jobs)
		for _, ep := range cfg.Endpoints {
			select {
			case <-ctx.Done():
				return
			case jobs <- ep:
			}
		}
	}()

	wg.Wait()

	if ctx.Err() != nil {
		return res, fmt.Errorf("run aborted: %w", ctx.Err())
	}
	return res, nil
}
