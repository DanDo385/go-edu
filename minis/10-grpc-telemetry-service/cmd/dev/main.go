// Simple debug program to test your functions directly
// without running the full test suite.
//
// This file uses fixed default values - no CLI arguments needed!
// Perfect for quick debugging sessions.
//
// To use:
// 1. Set breakpoints in ../../internal/exercise/grpctelemetryservice.go
// 2. Open this file (cmd/debug/main.go)
// 3. Use "Debug Main Program (Current Package)" configuration
// 4. Press F5 - that's it! The debugger will stop at your breakpoints
//
// Usage:
//   go run ./cmd/debug
//   # Or just press F5 in VS Code

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/example/go-10x-minis/minis/10-grpc-telemetry-service/internal/grpctelemetryservice"
	pb "github.com/example/go-10x-minis/minis/10-grpc-telemetry-service/proto"
)

func main() {
	// Fixed default values - modify these directly if you want to test different inputs
	window := 5 * time.Minute
	values := []float64{10.5, 20.3, 15.7}

	fmt.Println("=== Debugging GRPC Telemetry Aggregator ===")
	fmt.Printf("Window: %v\n", window)

	ctx := context.Background()

	// Set breakpoint in grpctelemetryservice.go at NewAggregator function
	agg := grpctelemetryservice.NewAggregator(window)

	// Push test points
	now := time.Now()
	points := make([]*pb.Point, len(values))
	for i, val := range values {
		points[i] = &pb.Point{
			Metric:    "test.metric",
			Value:     val,
			Timestamp: now.Add(time.Duration(-i) * time.Minute).Unix(),
		}
	}

	fmt.Printf("Pushing %d points...\n", len(points))
	for _, p := range points {
		err := agg.PushPoint(ctx, p)
		if err != nil {
			fmt.Printf("Error pushing point: %v\n", err)
			return
		}
		fmt.Printf("  Point: value=%.2f, timestamp=%d\n", p.Value, p.Timestamp)
	}

	// Get summary
	report := agg.Summary(ctx)
	fmt.Printf("\nSummary: %+v\n", report)
}

