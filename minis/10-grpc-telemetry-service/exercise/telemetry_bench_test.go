package exercise

import (
	"context"
	"testing"
	"time"

	pb "github.com/example/go-10x-minis/minis/10-grpc-telemetry-service/proto"
)

func BenchmarkAggregator_Push(b *testing.B) {
	agg := NewAggregator(1 * time.Hour)
	ctx := context.Background()
	point := &pb.Point{
		Metric:    "cpu",
		Value:     50.0,
		Timestamp: time.Now().Unix(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		agg.PushPoint(ctx, point)
	}
}

func BenchmarkAggregator_Summary(b *testing.B) {
	agg := NewAggregator(1 * time.Hour)
	ctx := context.Background()

	// Pre-fill with data
	for i := 0; i < 1000; i++ {
		agg.PushPoint(ctx, &pb.Point{
			Metric:    "cpu",
			Value:     float64(i),
			Timestamp: time.Now().Unix(),
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		agg.Summary(ctx)
	}
}

func BenchmarkAggregator_PushAndSummary(b *testing.B) {
	agg := NewAggregator(1 * time.Hour)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%10 < 9 {
			agg.PushPoint(ctx, &pb.Point{
				Metric:    "cpu",
				Value:     float64(i),
				Timestamp: time.Now().Unix(),
			})
		} else {
			agg.Summary(ctx)
		}
	}
}
