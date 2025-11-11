package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/example/go-10x-minis/minis/10-grpc-telemetry-service/exercise"
	pb "github.com/example/go-10x-minis/minis/10-grpc-telemetry-service/proto"
)

func main() {
	agg := exercise.NewAggregator(5 * time.Minute)

	// Simulate pushing metrics
	points := []*pb.Point{
		{Metric: "cpu.usage", Value: 45.2, Timestamp: time.Now().Unix()},
		{Metric: "cpu.usage", Value: 52.1, Timestamp: time.Now().Unix()},
		{Metric: "memory.used", Value: 1024.5, Timestamp: time.Now().Unix()},
		{Metric: "memory.used", Value: 1100.3, Timestamp: time.Now().Unix()},
		{Metric: "cpu.usage", Value: 48.7, Timestamp: time.Now().Unix()},
	}

	ctx := context.Background()
	for _, p := range points {
		if err := agg.PushPoint(ctx, p); err != nil {
			log.Fatalf("Push failed: %v", err)
		}
	}

	// Get summary
	report := agg.Summary(ctx)

	fmt.Println("=== Telemetry Summary ===\n")
	for metric, summary := range report.Metrics {
		fmt.Printf("Metric: %s\n", metric)
		fmt.Printf("  Count: %d\n", summary.Count)
		fmt.Printf("  Sum:   %.2f\n", summary.Sum)
		fmt.Printf("  Avg:   %.2f\n", summary.Avg)
		fmt.Printf("  Min:   %.2f\n", summary.Min)
		fmt.Printf("  Max:   %.2f\n", summary.Max)
		fmt.Println()
	}
}
