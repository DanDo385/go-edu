//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"time"
)

// Prober is a small interface for health checks (could be HTTP, RPC, etc).
type Prober interface {
	Probe(ctx context.Context, endpoint string) error
}

// Config controls the concurrent probe run.
type Config struct {
	Endpoints []string
	Workers   int
	Timeout   time.Duration
}

// Result aggregates successes and failures.
type Result struct {
	Successes map[string]time.Duration
	Failures  map[string]error
}

// Run is the student entry point for module 16-concurrency.
//
// TODOs for students:
// 1) Validate inputs: ctx/prober non-nil; default workers>0; default timeout.
// 2) Create a child context with timeout to bound the entire run.
// 3) Start a worker pool (goroutines) reading endpoints from a channel.
// 4) For each endpoint, call Probe with a per-request timeout; record latency or error.
// 5) Aggregate results into the Result struct and return it.
func Run(ctx context.Context, p Prober, cfg Config) (*Result, error) {
	panic("TODO: implement Run for 16-concurrency")
}
