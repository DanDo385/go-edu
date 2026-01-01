//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package grpctelemetryservice

import (
	"context"

	pb "github.com/example/go-10x-minis/minis/10-grpc-telemetry-service/proto"
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
// TODO: implement NewAggregator.
func NewAggregator(window time.Duration) Aggregator { panic("TODO: implement") }
// TODO: implement PushPoint.
func (a *aggregator) PushPoint(ctx context.Context, p *pb.Point) error { panic("TODO: implement") }
// TODO: implement Summary.
func (a *aggregator) Summary(ctx context.Context) *pb.Report { panic("TODO: implement") }
