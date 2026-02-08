//go:build reference

package sync

import (
	"context"
	"errors"
	"fmt"
)

var errNilClient = errors.New("nil sync client")

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

func Run(ctx context.Context, client SyncClient, cfg Config) (*Result, error) {
	_ = cfg
	if client == nil {
		return nil, errNilClient
	}

	progress, err := client.SyncProgress(ctx)
	if err != nil {
		return nil, fmt.Errorf("sync progress: %w", err)
	}
	if progress == nil {
		return &Result{IsSyncing: false, Progress: nil}, nil
	}

	copyProgress := *progress
	return &Result{IsSyncing: true, Progress: &copyProgress}, nil
}
