//go:build !solution && !reference

package grpctelemetryservice

import (
	"context"
	"math"
	"sync"
	"time"
	pb "github.com/example/go-10x-minis/minis/10-grpc-telemetry-service/proto"
)

func NewAggregator(window time.Duration) Aggregator {
	// TODO: Implement this function
	// TODO: Handle errors appropriately
	// TODO: Add necessary validations
	panic("not implemented")
}

func (a *aggregator) PushPoint(ctx context.Context, p *pb.Point) error {
	// TODO: Implement this function
	panic("not implemented")
}

func (a *aggregator) Summary(ctx context.Context) *pb.Report {
	// TODO: Implement this function
	panic("not implemented")
}
