//go:build !solution
// +build !solution

package exercise

import (
	"context"
	"math"
	"sync"
	"time"

	pb "github.com/example/go-10x-minis/minis/10-grpc-telemetry-service/proto"
)

// Aggregator collects and aggregates telemetry points.
type Aggregator interface {
	// PushPoint adds a telemetry point
	PushPoint(ctx context.Context, p *pb.Point) error

	// Summary returns aggregated statistics for all metrics
	Summary(ctx context.Context) *pb.Report
}

// NewAggregator creates a thread-safe aggregator with a rolling time window.
// Points older than 'window' are excluded from statistics.
func NewAggregator(window time.Duration) Aggregator {
	// TODO: Implement this function.
	// - Return a new instance of the `aggregator` struct.
	// - The `window` field should be set to the `window` parameter.
	// - The `points` field should be initialized as a new map: `make(map[string][]pointWithTime)`.
	return nil
}

type aggregator struct {
	mu     sync.RWMutex
	window time.Duration
	points map[string][]pointWithTime
}

type pointWithTime struct {
	value     float64
	timestamp time.Time
}

func (a *aggregator) PushPoint(ctx context.Context, p *pb.Point) error {
	// TODO: Implement this function.

	// This method adds a new data point to our aggregator.

	// Step 1: Lock the mutex.
	// - `a.mu.Lock()`
	// - Since multiple goroutines (from multiple gRPC requests) could be calling this method at the same time, we need to use a mutex to protect our `points` map from concurrent writes.
	// - A `defer a.mu.Unlock()` statement should be used to ensure the mutex is always unlocked, even if an error occurs.

	// Step 2: Append the new point.
	// - Convert the protobuf `Point`'s timestamp to a `time.Time` object: `time.Unix(p.Timestamp, 0)`.
	// - Create a `pointWithTime` struct with the value and timestamp.
	// - Append this new struct to the slice for the given metric: `a.points[p.Metric] = append(a.points[p.Metric], ...)`
	// - If the metric doesn't exist in the map, Go's `append` will gracefully handle creating a new slice.

	return nil
}

func (a *aggregator) Summary(ctx context.Context) *pb.Report {
	// TODO: Implement this function.

	// This method calculates the statistics for all metrics within the time window.

	// Step 1: Lock the mutex for reading.
	// - `a.mu.RLock()`
	// - We use a read lock (`RLock`) because we are not modifying the data, only reading it. This allows multiple readers to access the data concurrently, which is more performant than a full `Lock`.
	// - `defer a.mu.RUnlock()` ensures the lock is released.

	// Step 2: Determine the cutoff time.
	// - `cutoff := time.Now().Add(-a.window)`
	// - Any points with a timestamp before this time will be excluded.

	// Step 3: Iterate through the metrics.
	// - `for metric, pts := range a.points { ... }`

	// Step 4: For each metric, filter the points that are within the time window.
	// - Create a new slice `validPoints []float64`.
	// - Loop through the `pts` and if a point's timestamp is after the `cutoff`, add its value to `validPoints`.

	// Step 5: Calculate the summary statistics.
	// - If `len(validPoints)` is 0, `continue` to the next metric.
	// - Calculate the sum, min, max, and average of the `validPoints`.
	//   - Initialize `min` to `math.MaxFloat64` and `max` to `-math.MaxFloat64`.
	// - Create a `pb.MetricSummary` protobuf message with the results.

	// Step 6: Add the summary to the report.
	// - `report.Metrics[metric] = &pb.MetricSummary{...}`

	// Step 7: Return the report.

	// --- Memory & Language Comparison ---
	//
	// - Go: The use of `sync.RWMutex` is a key feature for controlling access to shared memory in a concurrent environment. It's more efficient than a simple `sync.Mutex` when there are many reads and few writes. The garbage collector will handle the memory of the slices that are no longer referenced.
	// - C: You would have to implement your own locking mechanism (e.g., with pthreads) and be very careful about memory management, especially when resizing arrays of points. This is very powerful but extremely error-prone.
	// - Rust: The `Mutex` or `RwLock` types from the standard library would be used. The borrow checker would provide compile-time guarantees that you are accessing the data safely, preventing data races.
	// - Java: You would use `synchronized` blocks or `ReentrantReadWriteLock` to control access. The JVM's garbage collector would handle memory management.
	// - Python: While Python has a Global Interpreter Lock (GIL) that prevents multiple threads from executing Python bytecode at the same time, you would still need locks to protect shared data structures from race conditions during I/O operations or when using libraries that release the GIL.

	return nil
}
