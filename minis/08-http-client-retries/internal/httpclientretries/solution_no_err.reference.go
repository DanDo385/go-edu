//go:build reference
// +build reference

package httpclientretries

/*
Core Retry Algorithm - Simplified Version
==========================================

This file focuses on the core exponential backoff algorithm without error handling.
Use this to understand the fundamental retry logic.

ALGORITHM STEPS:

1. Exponential Backoff Calculation:
   - delay = BaseDelay * (2^attempt)
   - attempt 0: BaseDelay * 1
   - attempt 1: BaseDelay * 2
   - attempt 2: BaseDelay * 4
   - attempt 3: BaseDelay * 8

2. Jitter Calculation:
   - jitter = delay * (rand [-0.2, 0.2])
   - finalDelay = delay + jitter
   - Spreads retries to prevent thundering herd

3. Retry Loop:
   - Try request
   - If success, return
   - Calculate backoff delay
   - Wait delay duration
   - Retry up to MaxRetries times

DEBUGGING GUIDE:
- Set breakpoint in retry loop to see attempt progression
- Watch delay calculation to understand exponential growth
- Watch jitter to see randomization effect
*/

// TODO: Implement SimpleExponentialBackoff
// Core algorithm for calculating exponential backoff delay
// Steps:
// 1. Calculate base delay: baseDelay * (2^attempt)
// 2. Add jitter: delay * (random value between -0.2 and 0.2)
// 3. Return final delay
// Signature: func SimpleExponentialBackoff(baseDelay time.Duration, attempt int) time.Duration

// TODO: Implement SimpleRetryLoop
// Core algorithm for retry loop without error handling
// Steps:
// 1. For each attempt from 0 to maxRetries:
//    - Try operation
//    - If success, return
//    - Calculate backoff delay
//    - Sleep for delay duration
// 2. Return after maxRetries exhausted
// Signature: func SimpleRetryLoop(maxRetries int, baseDelay time.Duration, operation func() bool)

// TODO: Implement SimpleJitterCalculation
// Core algorithm for adding jitter to delay
// Steps:
// 1. Generate random float between 0 and 1
// 2. Scale to range [-0.2, 0.2]: rand * 0.4 - 0.2
// 3. Multiply by delay to get jitter amount
// 4. Return jitter
// Signature: func SimpleJitterCalculation(delay time.Duration) time.Duration
