//go:build reference

package concurrency

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var errNilProber = errors.New("nil prober")

/*
Reference Solution

Structure:
- Build a bounded worker pool over endpoints.
- Probe each endpoint with per-probe timeout.
- Aggregate successes and failures with latency.

Invariants:
- Writes to result maps are synchronized with a mutex.
- Every endpoint is processed at most once.
*/
func Run(ctx context.Context, p Prober, cfg Config) (*Result, error) {
	if p == nil {
		return nil, errNilProber
	}
	if len(cfg.Endpoints) == 0 {
		return nil, errors.New("at least one endpoint is required")
	}

	workers := cfg.Workers
	if workers <= 0 {
		workers = len(cfg.Endpoints)
		if workers > 4 {
			workers = 4
		}
	}
	if workers > len(cfg.Endpoints) {
		workers = len(cfg.Endpoints)
	}

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 2 * time.Second
	}

	res := &Result{
		Successes: make(map[string]time.Duration, len(cfg.Endpoints)),
		Failures:  make(map[string]error),
	}

	jobs := make(chan string)
	var wg sync.WaitGroup
	var mu sync.Mutex

	worker := func() {
		defer wg.Done()
		for endpoint := range jobs {
			start := time.Now()
			probeCtx, cancel := context.WithTimeout(ctx, timeout)
			err := p.Probe(probeCtx, endpoint)
			cancel()
			elapsed := time.Since(start)

			mu.Lock()
			if err != nil {
				res.Failures[endpoint] = err
			} else {
				res.Successes[endpoint] = elapsed
			}
			mu.Unlock()
		}
	}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker()
	}

	for _, endpoint := range cfg.Endpoints {
		select {
		case <-ctx.Done():
			close(jobs)
			wg.Wait()
			return nil, fmt.Errorf("dispatch endpoints: %w", ctx.Err())
		case jobs <- endpoint:
		}
	}
	close(jobs)
	wg.Wait()

	return res, nil
}
