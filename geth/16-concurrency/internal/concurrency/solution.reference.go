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
Reference Solution - Concurrent Endpoint Probing
=================================================

This file demonstrates a worker-pool pattern for probing multiple RPC endpoints
concurrently. Bounded workers pull jobs from a channel; results collected with
mutex protection. Context for cancellation; per-probe timeout.

This connects to the Ethereum ecosystem by showing:
- Typical pattern: wallet/node lists probe endpoints to find healthy ones
- sync.WaitGroup: wait for all workers to drain jobs before returning
- context.WithTimeout: per-probe deadline; cancel() to release timer
- sync.Mutex: protect shared Result maps from concurrent writes

The exercise builds understanding of:
- Worker pool: N goroutines, single jobs channel, close to signal done
- jobs <- endpoint: send work; workers range over channel
- mu.Lock/Unlock around res.Successes/Failures: map not concurrent-safe

Teaching notes (per .cursorrules):
- close(jobs): signals workers to exit after range; must close after all sent.
- cancel() after Probe: releases WithTimeout resources; defer would be after
  err check but we call immediately since we're done with probeCtx.
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
