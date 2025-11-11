package exercise

import (
	"context"
	"testing"
	"time"

	pb "github.com/example/go-10x-minis/minis/10-grpc-telemetry-service/proto"
)

func TestAggregator_BasicStats(t *testing.T) {
	agg := NewAggregator(1 * time.Hour)
	ctx := context.Background()

	points := []*pb.Point{
		{Metric: "cpu", Value: 10.0, Timestamp: time.Now().Unix()},
		{Metric: "cpu", Value: 20.0, Timestamp: time.Now().Unix()},
		{Metric: "cpu", Value: 30.0, Timestamp: time.Now().Unix()},
	}

	for _, p := range points {
		agg.PushPoint(ctx, p)
	}

	report := agg.Summary(ctx)
	summary := report.Metrics["cpu"]

	if summary.Count != 3 {
		t.Errorf("Expected count=3, got %d", summary.Count)
	}
	if summary.Sum != 60.0 {
		t.Errorf("Expected sum=60, got %.2f", summary.Sum)
	}
	if summary.Avg != 20.0 {
		t.Errorf("Expected avg=20, got %.2f", summary.Avg)
	}
	if summary.Min != 10.0 {
		t.Errorf("Expected min=10, got %.2f", summary.Min)
	}
	if summary.Max != 30.0 {
		t.Errorf("Expected max=30, got %.2f", summary.Max)
	}
}

func TestAggregator_TimeWindow(t *testing.T) {
	agg := NewAggregator(100 * time.Millisecond)
	ctx := context.Background()

	// Old point (will expire)
	oldPoint := &pb.Point{
		Metric:    "cpu",
		Value:     10.0,
		Timestamp: time.Now().Add(-200 * time.Millisecond).Unix(),
	}
	agg.PushPoint(ctx, oldPoint)

	// New point (within window)
	newPoint := &pb.Point{
		Metric:    "cpu",
		Value:     20.0,
		Timestamp: time.Now().Unix(),
	}
	agg.PushPoint(ctx, newPoint)

	report := agg.Summary(ctx)
	summary := report.Metrics["cpu"]

	if summary.Count != 1 {
		t.Errorf("Expected count=1 (old point expired), got %d", summary.Count)
	}
	if summary.Avg != 20.0 {
		t.Errorf("Expected avg=20 (only new point), got %.2f", summary.Avg)
	}
}

func TestAggregator_MultipleMetrics(t *testing.T) {
	agg := NewAggregator(1 * time.Hour)
	ctx := context.Background()

	agg.PushPoint(ctx, &pb.Point{Metric: "cpu", Value: 50.0, Timestamp: time.Now().Unix()})
	agg.PushPoint(ctx, &pb.Point{Metric: "memory", Value: 1024.0, Timestamp: time.Now().Unix()})

	report := agg.Summary(ctx)

	if len(report.Metrics) != 2 {
		t.Errorf("Expected 2 metrics, got %d", len(report.Metrics))
	}

	if report.Metrics["cpu"].Avg != 50.0 {
		t.Error("CPU metric incorrect")
	}
	if report.Metrics["memory"].Avg != 1024.0 {
		t.Error("Memory metric incorrect")
	}
}
