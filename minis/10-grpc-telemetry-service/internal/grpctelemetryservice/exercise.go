//go:build !solution && !reference

package grpctelemetryservice

import (
	"context"
	pb "github.com/example/go-10x-minis/minis/10-grpc-telemetry-service/proto"
	"math"
	"sync"
	"time"
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

// NewAggregator implements the exercise.
//
// TODO: Implement this function
func NewAggregator(window time.Duration) Aggregator {
	// TODO: Implement
	return Aggregator{}
}

// PushPoint implements the exercise.
//
// TODO: Implement this function
func (a *aggregator) PushPoint(ctx context.Context, p *pb.Point) error {
	// TODO: Implement
	return nil
}

// Summary implements the exercise.
//
// TODO: Implement this function
func (a *aggregator) Summary(ctx context.Context) *pb.Report {
	// TODO: Implement
	return nil
}
