package miniserviceallfeatures

// CircuitState represents the state of a circuit breaker.
type CircuitState int

const (
	StateClosed CircuitState = iota
	StateOpen
	StateHalfOpen
)

// Job is a unit of work for the worker pool.
type Job func() error

// EventHandler handles an event published on the event bus.
type EventHandler func(event interface{})
