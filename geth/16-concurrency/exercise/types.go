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
