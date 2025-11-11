//go:build solution
// +build solution

/*
Problem: Build a gRPC telemetry aggregator with streaming and time windows

Requirements:
1. Accept streaming telemetry points (metric name, value, timestamp)
2. Aggregate statistics per metric (count, sum, avg, min, max)
3. Support rolling time window (exclude old data)
4. Thread-safe concurrent access
5. gRPC service implementation

Why Go is well-suited:
- gRPC: First-class support with protoc-gen-go
- Concurrency: Goroutines + channels for streaming
- Performance: Fast aggregation with mutexes

Compared to other languages:
- Python: grpcio similar, but slower aggregation
- Node.js: grpc-js works, but single-threaded
- Rust: tonic is excellent, more complex
*/

package exercise

import (
	"context"
	"math"
	"sync"
	"time"

	pb "github.com/example/go-10x-minis/minis/10-grpc-telemetry-service/proto"
)

type Aggregator interface {
	PushPoint(ctx context.Context, p *pb.Point) error
	Summary(ctx context.Context) *pb.Report
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

func NewAggregator(window time.Duration) Aggregator {
	return &aggregator{
		window: window,
		points: make(map[string][]pointWithTime),
	}
}

func (a *aggregator) PushPoint(ctx context.Context, p *pb.Point) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	ts := time.Unix(p.Timestamp, 0)
	a.points[p.Metric] = append(a.points[p.Metric], pointWithTime{
		value:     p.Value,
		timestamp: ts,
	})

	return nil
}

func (a *aggregator) Summary(ctx context.Context) *pb.Report {
	a.mu.RLock()
	defer a.mu.RUnlock()

	cutoff := time.Now().Add(-a.window)
	report := &pb.Report{Metrics: make(map[string]*pb.MetricSummary)}

	for metric, pts := range a.points {
		var validPoints []float64
		for _, pt := range pts {
			if pt.timestamp.After(cutoff) {
				validPoints = append(validPoints, pt.value)
			}
		}

		if len(validPoints) == 0 {
			continue
		}

		sum := 0.0
		min := math.MaxFloat64
		max := -math.MaxFloat64

		for _, v := range validPoints {
			sum += v
			if v < min {
				min = v
			}
			if v > max {
				max = v
			}
		}

		report.Metrics[metric] = &pb.MetricSummary{
			Count: int32(len(validPoints)),
			Sum:   sum,
			Avg:   sum / float64(len(validPoints)),
			Min:   min,
			Max:   max,
		}
	}

	return report
}
