//go:build reference

/*
Problem: Build a gRPC telemetry aggregator with streaming and time windows

Requirements:
1. Accept streaming telemetry points (metric name, value, timestamp)
2. Aggregate statistics per metric (count, sum, avg, min, max)
3. Support rolling time window (exclude old data)
4. Thread-safe concurrent access
5. gRPC service implementation

Algorithm: Rolling Time Window Aggregation
- Store points with timestamps
- Filter points within time window on Summary
- Calculate statistics on filtered points
- Use RWMutex for concurrent access

Time Window Algorithm:
- cutoff = now - window duration
- For each point: if timestamp > cutoff, include in stats
- This is "lazy" cleanup (filter on read, not proactive deletion)

Statistics Calculation:
- Iterate filtered points once
- Accumulate sum, track min/max
- Calculate average = sum / count

Concurrency Strategy:
- Write lock for PushPoint (modifies map/slices)
- Read lock for Summary (reads data, doesn't modify)
- RWMutex allows concurrent readers
*/

package grpctelemetryservice

import (
	"context"
	"math"
	"sync"
	"time"

	pb "github.com/example/go-10x-minis/minis/10-grpc-telemetry-service/proto"
)

// Aggregator interface for telemetry collection.
// BREAKPOINT: Set breakpoint in implementations to trace aggregation
type Aggregator interface {
	PushPoint(ctx context.Context, p *pb.Point) error
	Summary(ctx context.Context) *pb.Report
}

// aggregator implements thread-safe telemetry aggregation.
// BREAKPOINT: Set breakpoint in methods to trace operations
type aggregator struct {
	mu     sync.RWMutex
	window time.Duration
	points map[string][]pointWithTime
}

// pointWithTime stores a metric value with timestamp.
// BREAKPOINT: Set breakpoint when creating these to see data flow
type pointWithTime struct {
	value     float64
	timestamp time.Time
}

// NewAggregator creates a new aggregator with rolling time window.
//
// BREAKPOINT: Set breakpoint here to trace aggregator creation
// DEBUG: Watch 'window' to see time window configuration
func NewAggregator(window time.Duration) Aggregator {
	// DEBUG: Creating aggregator with time window
	return &aggregator{
		window: window,
		points: make(map[string][]pointWithTime),
	}
}

// PushPoint adds a telemetry point.
//
// Algorithm:
// 1. Acquire write lock (modifying shared data)
// 2. Convert protobuf timestamp to time.Time
// 3. Append point to metric's slice
//
// BREAKPOINT: Set breakpoint at function entry to trace point ingestion
// DEBUG: Watch 'p.Metric' to see metric name
// DEBUG: Watch 'p.Value' to see metric value
// DEBUG: Watch 'p.Timestamp' to see point timestamp
func (a *aggregator) PushPoint(ctx context.Context, p *pb.Point) error {
	// BREAKPOINT: Set breakpoint here before acquiring lock
	a.mu.Lock()
	defer a.mu.Unlock()

	// DEBUG: Converting Unix timestamp to time.Time
	// BREAKPOINT: Set breakpoint here to see timestamp conversion
	ts := time.Unix(p.Timestamp, 0)
	// DEBUG: Watch 'ts' to see converted time

	// BREAKPOINT: Set breakpoint here to see point storage
	// DEBUG: Watch 'a.points[p.Metric]' before append to see existing points
	a.points[p.Metric] = append(a.points[p.Metric], pointWithTime{
		value:     p.Value,
		timestamp: ts,
	})
	// DEBUG: Watch 'a.points[p.Metric]' after append to see new point added

	return nil
}

// Summary calculates aggregated statistics for all metrics.
//
// Algorithm:
// 1. Acquire read lock (concurrent reads allowed)
// 2. Calculate cutoff time (now - window)
// 3. For each metric:
//    - Filter points within time window
//    - Calculate sum, min, max, count
//    - Calculate average = sum / count
// 4. Return report with all metrics
//
// BREAKPOINT: Set breakpoint at function entry to trace summary generation
// DEBUG: Watch 'a.window' to see time window
// DEBUG: Watch 'cutoff' to see filtering threshold
func (a *aggregator) Summary(ctx context.Context) *pb.Report {
	// BREAKPOINT: Set breakpoint here before acquiring read lock
	a.mu.RLock()
	defer a.mu.RUnlock()

	// Calculate cutoff time for rolling window
	// BREAKPOINT: Set breakpoint here to see cutoff calculation
	// DEBUG: Watch 'time.Now()' to see current time
	// DEBUG: Watch 'a.window' to see window duration
	cutoff := time.Now().Add(-a.window)
	// DEBUG: Watch 'cutoff' to see threshold time

	// Create report
	report := &pb.Report{Metrics: make(map[string]*pb.MetricSummary)}

	// BREAKPOINT: Set breakpoint here to trace metric iteration
	// DEBUG: Watch 'metric' to see current metric name
	// DEBUG: Watch 'len(pts)' to see total points for metric
	for metric, pts := range a.points {
		var validPoints []float64

		// Filter points within time window
		// BREAKPOINT: Set breakpoint here to trace filtering
		// DEBUG: Watch 'pt.timestamp' vs 'cutoff' for each point
		for _, pt := range pts {
			// DEBUG: Checking if point is within time window
			if pt.timestamp.After(cutoff) {
				// DEBUG: Point is valid, including in statistics
				validPoints = append(validPoints, pt.value)
			}
			// DEBUG: Old point excluded
		}

		// Skip metrics with no valid points
		// BREAKPOINT: Set breakpoint here to check if metric has data
		// DEBUG: Watch 'len(validPoints)' to see count
		if len(validPoints) == 0 {
			// DEBUG: No valid points for this metric, skipping
			continue
		}

		// Calculate statistics
		// BREAKPOINT: Set breakpoint here to trace stats calculation
		sum := 0.0
		min := math.MaxFloat64
		max := -math.MaxFloat64

		// DEBUG: Iterating through valid points to calculate stats
		// BREAKPOINT: Set breakpoint here to see accumulation
		for _, v := range validPoints {
			// DEBUG: Watch 'v' to see current value
			sum += v
			// DEBUG: Watch 'sum' to see running total

			if v < min {
				// DEBUG: New minimum found
				min = v
			}
			if v > max {
				// DEBUG: New maximum found
				max = v
			}
		}
		// DEBUG: Watch 'sum', 'min', 'max' to see final values

		// Calculate average
		// BREAKPOINT: Set breakpoint here to see average calculation
		// DEBUG: Watch 'sum' and 'len(validPoints)' to see division
		avg := sum / float64(len(validPoints))
		// DEBUG: Watch 'avg' to see result

		// Create metric summary
		// BREAKPOINT: Set breakpoint here to see summary creation
		// DEBUG: Watch summary fields to see final statistics
		report.Metrics[metric] = &pb.MetricSummary{
			Count: int32(len(validPoints)),
			Sum:   sum,
			Avg:   avg,
			Min:   min,
			Max:   max,
		}
		// DEBUG: Summary added to report
	}

	// DEBUG: Watch 'report' to see complete summary
	// BREAKPOINT: Set breakpoint here to inspect final report
	return report
}
