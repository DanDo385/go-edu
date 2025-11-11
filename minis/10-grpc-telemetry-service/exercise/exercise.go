package exercise

import (
	"context"
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
	// TODO: implement
	return nil
}
