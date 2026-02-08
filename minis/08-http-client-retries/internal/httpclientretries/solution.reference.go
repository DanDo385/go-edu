//go:build reference

package httpclientretries

/*
Reference Solution - HTTP Client with Retry and Exponential Backoff
==================================================================

This file implements a custom http.RoundTripper that wraps another RoundTripper
and adds retry logic with exponential backoff and jitter. Any HTTP client using
this transport automatically retries on transient failures.

This connects to Go's HTTP stack:
- RoundTripper: the interface that executes a single HTTP transaction
- http.Client.Transport: defaults to DefaultTransport, can be replaced
- Request context: we respect req.Context().Done() so retries can be canceled

The exercise teaches:
- Middleware/wrapper pattern: wrap existing transport, add behavior
- Exponential backoff: delay doubles each attempt (1s, 2s, 4s, ...)
- Jitter: randomize delay to avoid thundering herd when many clients retry
- Retryable vs non-retryable: network errors retry; 4xx client errors don't

Teaching notes:
- Body ownership: on retry, we must close the response body before retrying,
  otherwise we leak connections. The client only receives the final response.
- Context: select on req.Context().Done() during sleep so user cancellation
  is honored. Timer.Stop() prevents timer leak when we exit early.
*/

import (
	"math/rand"
	"net/http"
	"time"
)

// retryRoundTripper wraps a RoundTripper and adds retries with exponential backoff.
type retryRoundTripper struct {
	next       http.RoundTripper
	maxRetries int
	baseDelay  time.Duration
	jitter     float64
}

/*
NewRetryClient - HTTP Client with Retry Transport

Returns an http.Client whose Transport retries failed requests.
maxRetries: total attempts = 1 + maxRetries (e.g. 3 retries = 4 total attempts)
baseDelay: initial delay before first retry
jitter: 0..1 multiplier; 0.2 means ±20% randomization to stagger retries
*/
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

/*
isRetryable - Should We Retry?

Network errors: always retry (timeout, connection refused, etc.)
5xx / 429 / 502 / 503 / 504: retry (transient server/rate-limit issues)
4xx (except 429): don't retry (client error, retrying won't help)
2xx, 3xx: not an error, don't retry
*/
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

/*
RoundTrip - Execute Request with Retries

Attempts up to maxRetries+1 times. On retryable failure: close response body
(to release connection), sleep with exponential backoff + jitter, retry.
Respects request context: if canceled during sleep, returns immediately.
*/
func (rt *retryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	for attempt := 0; attempt <= rt.maxRetries; attempt++ {
		resp, err := rt.next.RoundTrip(req)
		if !isRetryable(err, resp) || attempt == rt.maxRetries {
			return resp, err
		}

		// Must close body before retry - otherwise connection leaks
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}

		// Exponential backoff: 1<<0=1, 1<<1=2, 1<<2=4, ...
		backoff := rt.baseDelay * time.Duration(1<<attempt)
		delay := backoff
		if rt.jitter > 0 && backoff > 0 {
			// Jitter: add random offset in [-maxJitter, +maxJitter]
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
