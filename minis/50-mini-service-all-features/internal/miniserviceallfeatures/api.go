package miniserviceallfeatures

import (
	"context"
	"time"
)

// Cache is the learner-facing cache type.
type Cache = SolutionCache

func NewCache() *Cache { return NewSolutionCache() }

// CircuitBreaker is the learner-facing circuit breaker type.
type CircuitBreaker = SolutionCircuitBreaker

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return NewSolutionCircuitBreaker(threshold, timeout)
}

// WorkerPool is the learner-facing worker pool type.
type WorkerPool = SolutionWorkerPool

func NewWorkerPool(numWorkers int) *WorkerPool { return NewSolutionWorkerPool(numWorkers) }

// RetryWithBackoff retries fn with exponential backoff.
func RetryWithBackoff(ctx context.Context, maxRetries int, initialDelay time.Duration, fn func() error) error {
	return SolutionRetryWithBackoff(ctx, maxRetries, initialDelay, fn)
}

// UserRateLimiter limits per-user actions.
type UserRateLimiter = SolutionUserRateLimiter

func NewUserRateLimiter(requestsPerSecond float64, burst int) *UserRateLimiter {
	return NewSolutionUserRateLimiter(requestsPerSecond, burst)
}

// AppError is a structured application error.
type AppError = SolutionAppError

func NewNotFoundError(message string) AppError   { return SolutionNewNotFoundError(message) }
func NewBadRequestError(message string) AppError { return SolutionNewBadRequestError(message) }
func NewInternalError(message string) AppError   { return SolutionNewInternalError(message) }

// EventBus is the learner-facing event bus type.
type EventBus = SolutionEventBus

func NewEventBus() *EventBus { return NewSolutionEventBus() }
