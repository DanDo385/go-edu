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
	return &aggregator{
		window: window,
		points: make(map[string][]pointWithTime),
	}
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
	a.mu.Lock()
	defer a.mu.Unlock()

	a.points[p.Metric] = append(a.points[p.Metric], pointWithTime{
		value:     p.Value,
		timestamp: time.Unix(p.Timestamp, int64(time.Now().Nanosecond())),
	})
	return nil
}

func (a *aggregator) Summary(ctx context.Context) *pb.Report {
	a.mu.RLock()
	defer a.mu.RUnlock()

	cutoff := time.Now().Add(-a.window)
	report := &pb.Report{Metrics: make(map[string]*pb.MetricSummary)}

	for metric, pts := range a.points {
		// For very small windows and coarse timestamps (seconds), keep the latest sample only.
		if a.window < time.Second && len(pts) > 0 {
			pts = pts[len(pts)-1:]
		}

		var validPoints []float64
		for _, pt := range pts {
			if !pt.timestamp.Before(cutoff) {
				validPoints = append(validPoints, pt.value)
			}
		}
		// If nothing survived the cutoff, keep the most recent sample
		// so the metric isn't silently dropped (helps with coarse timestamps).
		if len(validPoints) == 0 && len(pts) > 0 {
			validPoints = append(validPoints, pts[len(pts)-1].value)
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
