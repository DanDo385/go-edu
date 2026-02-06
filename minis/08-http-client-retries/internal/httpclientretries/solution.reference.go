//go:build reference

package httpclientretries

import (
	"math/rand"
	"net/http"
	"time"
)

/*
Reference Solution
==================

This file is the canonical reference for this exercise. It keeps error paths
explicit when an operation can fail, so callers decide how to handle failure.

Read this alongside exercise.go and the tests to understand the intended data
flow, boundary checks, and invariants.
*/

// retryRoundTripper is a struct that implements the http.RoundTripper interface
// to provide retry logic.
type retryRoundTripper struct {
	next       http.RoundTripper
	maxRetries int
	baseDelay  time.Duration
	jitter     float64
}

// NewRetryClient creates a new http.Client that automatically retries failed
// requests using exponential backoff with jitter.
func NewRetryClient(maxRetries int, baseDelay time.Duration, jitter float64) *http.Client {
	if maxRetries < 0 {
		maxRetries = 0
	}
	if jitter < 0 {
		jitter = 0
	}
	return &http.Client{
		Transport: &retryRoundTripper{
			next:       http.DefaultTransport,
			maxRetries: maxRetries,
			baseDelay:  baseDelay,
			jitter:     jitter,
		},
	}
}

// isRetryable checks if an error or response indicates that a request is retryable.
func isRetryable(err error, resp *http.Response) bool {
	if err != nil {
		return true
	}
	if resp == nil {
		return false
	}
	switch resp.StatusCode {
	case http.StatusTooManyRequests, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

// RoundTrip executes a single HTTP transaction, adding retry logic.
func (rt *retryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	for attempt := 0; attempt <= rt.maxRetries; attempt++ {
		resp, err := rt.next.RoundTrip(req)
		if !isRetryable(err, resp) || attempt == rt.maxRetries {
			return resp, err
		}

		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}

		backoff := rt.baseDelay * time.Duration(1<<attempt)
		delay := backoff
		if rt.jitter > 0 && backoff > 0 {
			maxJitter := time.Duration(float64(backoff) * rt.jitter)
			j := time.Duration(rand.Int63n(int64(maxJitter)*2+1)) - maxJitter
			delay += j
			if delay < 0 {
				delay = 0
			}
		}

		timer := time.NewTimer(delay)
		select {
		case <-req.Context().Done():
			timer.Stop()
			return nil, req.Context().Err()
		case <-timer.C:
		}
	}

	return nil, nil
}
