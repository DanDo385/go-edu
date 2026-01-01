//go:build !solution && !reference

package grpctelemetryservice

/*
Problem: Build a gRPC telemetry aggregator with streaming and time windows
Requirements:
1. Accept streaming telemetry points (metric name, value, timestamp)
2. Aggregate statistics per metric (count, sum, avg, min, max)
3. Support rolling time window (exclude old data)
Algorithm: Rolling Time Window Aggregation
- Store points with timestamps
- Filter points within time window on Summary
- Calculate statistics on filtered points
- Use RWMutex for concurrent access
*/

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

// NewAggregator - TODO: implement this function
func NewAggregator(window time.Duration) Aggregator {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil
}

// PushPoint - TODO: implement this function
func (a *aggregator) PushPoint(ctx context.Context, p *pb.Point) error {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil, nil
}

// Summary - TODO: implement this function
func (a *aggregator) Summary(ctx context.Context) *pb.Report {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil
}

