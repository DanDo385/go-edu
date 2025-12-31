package main

import (
	"context"
	"fmt"    // Formatted I/O
	"log"    // Logging
	"math/rand"
	"os"     // Operating system interface
	"strconv" // String to int conversion
	"time"   // Time operations

	"github.com/example/go-10x-minis/minis/10-grpc-telemetry-service/internal/grpctelemetryservice"
	pb "github.com/example/go-10x-minis/minis/10-grpc-telemetry-service/proto"
)

/*
===================================================================
gRPC Telemetry Service with Aggregation
===================================================================

This program demonstrates a telemetry aggregation service that
collects metrics points and computes statistics (count, sum, avg, min, max).

USAGE:
    go run ./cmd/app/main.go
    go run ./cmd/app/main.go 10m
    go run ./cmd/app/main.go 5m 10
    go run ./cmd/app/main.go 5m 100 5

Arguments:
    [window]  - Time window for aggregation, e.g., "5m", "10m" (default: "5m")
    [points]  - Number of data points to generate (default: 5)
    [metrics] - Number of different metrics to generate (default: 2)
    [seed]    - Random seed (default: current time)

DEBUGGING:
- Set breakpoints at "// BREAKPOINT:" comments
- Use F5 to start debugging (select "Debug: Run main.go")
- Step Into (F11) to enter exercise package functions
- Watch Variables panel to see metrics aggregation
- Inspect aggregator state to see running totals
- See /RUN_DEBUG.md for comprehensive debugging guide

AGGREGATION BEHAVIOR:
- Collects metrics points with metric name, value, timestamp
- Groups points by metric name
- Computes statistics: count, sum, avg, min, max
- Time window: points older than window are discarded
*/

func main() {
	// ============================================================================
	// STEP 1: Parse Command-Line Arguments
	// ============================================================================
	// BREAKPOINT: Set breakpoint here to inspect command-line arguments
	// DEBUG: os.Args[0] is the program name
	// DEBUG: os.Args[1] is the first argument (if provided)
	// DEBUG: In Variables panel, expand 'os.Args' to see all arguments

	// Default values
	window := 5 * time.Minute
	numPoints := 5
	numMetrics := 2
	seed := time.Now().UnixNano()

	// Override defaults if arguments provided
	if len(os.Args) > 1 {
		if w, err := time.ParseDuration(os.Args[1]); err == nil {
			window = w
		}
	}
	if len(os.Args) > 2 {
		if p, err := strconv.Atoi(os.Args[2]); err == nil {
			numPoints = p
		}
	}
	if len(os.Args) > 3 {
		if m, err := strconv.Atoi(os.Args[3]); err == nil {
			numMetrics = m
		}
	}
	if len(os.Args) > 4 {
		if s, err := strconv.ParseInt(os.Args[4], 10, 64); err == nil {
			seed = s
		}
	}

	// BREAKPOINT: Set breakpoint here to inspect parsed values
	// DEBUG: After parsing, variables now hold actual command-line values
	// DEBUG: Watch 'window', 'numPoints', 'numMetrics' - they may have changed

	// Seed random number generator for reproducible results
	rand.Seed(seed)

	log.Printf("Configuration: window=%v, points=%d, metrics=%d, seed=%d",
		window, numPoints, numMetrics, seed)

	// ============================================================================
	// STEP 2: Create Aggregator Instance
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before aggregator creation
	// DEBUG: Step Into (F11) on grpctelemetryservice.NewAggregator() to see initialization
	// DEBUG: Watch how the aggregator struct is allocated

	agg := grpctelemetryservice.NewAggregator(window)

	// BREAKPOINT: Set breakpoint here after aggregator creation
	// DEBUG: Inspect 'agg' variable in Variables panel
	// DEBUG: Expand agg to see internal fields: window, metrics map, mutex
	// DEBUG: Initially, metrics map should be empty

	fmt.Println("=== Telemetry Aggregation Demo ===\n")
	fmt.Printf("Time window: %v\n", window)
	fmt.Printf("Generating %d points across %d metrics\n\n", numPoints, numMetrics)

	// ============================================================================
	// STEP 3: Generate Sample Metrics Data
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before data generation
	// DEBUG: We'll create sample metrics to simulate telemetry data
	// DEBUG: Watch how points are created with different metric names

	metricNames := generateMetricNames(numMetrics)

	// BREAKPOINT: Set breakpoint here after generating metric names
	// DEBUG: Expand 'metricNames' to see all generated names
	// DEBUG: These will be used to group data points

	points := make([]*pb.Point, 0, numPoints)
	now := time.Now()

	// BREAKPOINT: Set breakpoint here before point generation loop
	// DEBUG: We'll generate random data points for our metrics
	// DEBUG: Points will have realistic values and timestamps

	for i := 0; i < numPoints; i++ {
		// BREAKPOINT: Set breakpoint here inside generation loop
		// DEBUG: Watch 'i' increment with each iteration
		// DEBUG: Observe how random metric names and values are chosen

		metricName := metricNames[rand.Intn(len(metricNames))]
		value := generateRandomValue(metricName)
		timestamp := now.Add(-time.Duration(rand.Intn(60)) * time.Second).Unix()

		point := &pb.Point{
			Metric:    metricName,
			Value:     value,
			Timestamp: timestamp,
		}

		// BREAKPOINT: Set breakpoint here after creating point
		// DEBUG: Inspect 'point' to see metric name, value, timestamp
		// DEBUG: Expand point to see all fields

		points = append(points, point)

		// Verbose logging removed for os.Args version
	}

	// BREAKPOINT: Set breakpoint here after generating all points
	// DEBUG: Expand 'points' to see all generated data points
	// DEBUG: Check len(points) == numPoints

	// ============================================================================
	// STEP 4: Push Points to Aggregator
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before pushing points
	// DEBUG: We'll send all points to the aggregator
	// DEBUG: Aggregator will group by metric and update statistics

	ctx := context.Background()

	fmt.Println("Pushing data points to aggregator...")

	for i, p := range points {
		// BREAKPOINT: Set breakpoint here inside push loop
		// DEBUG: Watch 'i' and 'p' change with each iteration
		// DEBUG: Step Into (F11) on agg.PushPoint() to see aggregation logic

		if err := agg.PushPoint(ctx, p); err != nil {
			// BREAKPOINT: Set breakpoint here on push error
			// DEBUG: Inspect 'err' to see what went wrong
			// DEBUG: Check if point data is valid

			log.Fatalf("Push failed for point %d: %v", i, err)
		}

		// BREAKPOINT: Set breakpoint here after successful push
		// DEBUG: Point has been added to aggregator
		// DEBUG: Aggregator statistics have been updated
		// DEBUG: Step Into (F11) to see how aggregator updates running totals

		// Verbose logging removed for os.Args version
	}

	// BREAKPOINT: Set breakpoint here after pushing all points
	// DEBUG: All points have been processed
	// DEBUG: Aggregator contains grouped statistics

	fmt.Printf("Pushed %d points\n\n", len(points))

	// ============================================================================
	// STEP 5: Get Aggregated Summary
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before getting summary
	// DEBUG: Step Into (F11) on agg.Summary() to see how stats are computed
	// DEBUG: Watch how the report is built from aggregated data

	report := agg.Summary(ctx)

	// BREAKPOINT: Set breakpoint here after getting summary
	// DEBUG: Inspect 'report' variable in Variables panel
	// DEBUG: Expand report to see all metric summaries
	// DEBUG: Each metric should have count, sum, avg, min, max

	// ============================================================================
	// STEP 6: Display Results
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before displaying results
	// DEBUG: We'll iterate through each metric and print statistics
	// DEBUG: Watch how summary data is formatted for output

	fmt.Println("=== Telemetry Summary ===\n")

	if len(report.Metrics) == 0 {
		// BREAKPOINT: Set breakpoint here if no metrics found
		// DEBUG: This could mean all points were outside the time window
		// DEBUG: Or no points were successfully pushed

		fmt.Println("No metrics found in the time window")
	}

	for metric, summary := range report.Metrics {
		// BREAKPOINT: Set breakpoint here inside display loop
		// DEBUG: Watch 'metric' and 'summary' change with each iteration
		// DEBUG: Inspect summary fields: Count, Sum, Avg, Min, Max

		fmt.Printf("Metric: %s\n", metric)
		fmt.Printf("  Count: %d\n", summary.Count)
		fmt.Printf("  Sum:   %.2f\n", summary.Sum)
		fmt.Printf("  Avg:   %.2f\n", summary.Avg)
		fmt.Printf("  Min:   %.2f\n", summary.Min)
		fmt.Printf("  Max:   %.2f\n", summary.Max)

		// BREAKPOINT: Set breakpoint here after printing metric stats
		// DEBUG: Verify the statistics look correct
		// DEBUG: Check that Count > 0, Min <= Avg <= Max, Sum = Count * Avg

		fmt.Println()
	}

	// ============================================================================
	// STEP 7: Additional Analysis
	// ============================================================================
	// BREAKPOINT: Set breakpoint here before additional analysis
	// DEBUG: We'll compute some additional statistics

	fmt.Println("=== Additional Analysis ===\n")

	totalPoints := 0
	totalMetrics := len(report.Metrics)

	for _, summary := range report.Metrics {
		totalPoints += int(summary.Count)
	}

	// BREAKPOINT: Set breakpoint here after computing totals
	// DEBUG: Inspect 'totalPoints' and 'totalMetrics'
	// DEBUG: totalPoints should equal len(points) (if all within window)

	fmt.Printf("Total metrics: %d\n", totalMetrics)
	fmt.Printf("Total points:  %d\n", totalPoints)
	fmt.Printf("Points per metric (avg): %.2f\n\n", float64(totalPoints)/float64(totalMetrics))

	// ============================================================================
	// STEP 8: Demonstrate Time Window Behavior
	// ============================================================================
	// Time window demonstration
	fmt.Println("=== Time Window Demonstration ===\n")

	// Add a very old point (outside window)
	oldPoint := &pb.Point{
		Metric:    "old.metric",
		Value:     999.0,
		Timestamp: time.Now().Add(-10 * window).Unix(),
	}

		// BREAKPOINT: Set breakpoint here before pushing old point
		// DEBUG: This point is outside the time window
		// DEBUG: Step Into (F11) to see how aggregator handles it

		if err := agg.PushPoint(ctx, oldPoint); err != nil {
			log.Printf("Old point rejected (expected): %v", err)
		}

		// Get updated summary
		report = agg.Summary(ctx)

		// BREAKPOINT: Set breakpoint here after pushing old point
		// DEBUG: Check if "old.metric" appears in report
		// DEBUG: It should NOT appear (outside time window)

		if _, ok := report.Metrics["old.metric"]; ok {
			fmt.Println("Old metric found (unexpected!)")
		} else {
			fmt.Println("Old metric correctly excluded (outside window)")
		}
	}

	// ============================================================================
	// DEBUGGING TIPS
	// ============================================================================
	//
	// KEY DEBUGGING TECHNIQUES:
	//
	// 1. Set breakpoints at "// BREAKPOINT:" comments
	// 2. Press F5 to start debugging
	// 3. Use F10 (Step Over) to execute line by line
	// 4. Use F11 (Step Into) to enter exercise package functions
	// 5. Watch Variables panel to see aggregator state
	// 6. Add watch expressions: len(points), report.Metrics
	// 7. Use Call Stack panel to see function hierarchy
	//
	// AGGREGATOR DEBUGGING:
	// - Step Into grpctelemetryservice.NewAggregator() to see initialization
	// - Step Into agg.PushPoint() to see how points are aggregated
	// - Watch the metrics map grow as points are added
	// - Inspect running totals: count, sum, min, max
	// - Use Debug Console: agg.Summary(ctx)
	//
	// STATISTICS VERIFICATION:
	// - Count: number of points for each metric
	// - Sum: total of all values
	// - Avg: sum / count
	// - Min: smallest value seen
	// - Max: largest value seen
	// - Verify: Min <= Avg <= Max
	// - Verify: Sum = Count * Avg (approximately)
	//
	// TIME WINDOW DEBUGGING:
	// - Points older than window are excluded
	// - Check point timestamps vs current time
	// - Use time.Now().Unix() in Debug Console
	// - Verify old points are rejected or filtered
	//
	// PROTOBUF DEBUGGING:
	// - pb.Point is a generated protobuf message
	// - Inspect Point fields: Metric, Value, Timestamp
	// - Use proto.MarshalText() to see text representation
	//
	// DATA GENERATION:
	// - Random seed controls reproducibility
	// - Same seed = same data points
	// - Change seed to get different data
	// - Use fixed seed for consistent debugging
	//
	// See /RUN_DEBUG.md for comprehensive debugging guide
}

// ============================================================================
// Helper Functions for Data Generation
// ============================================================================

// generateMetricNames creates realistic metric names
func generateMetricNames(count int) []string {
	// BREAKPOINT: Set breakpoint here to see metric name generation
	// DEBUG: This creates realistic telemetry metric names

	prefixes := []string{"cpu", "memory", "disk", "network", "response_time"}
	suffixes := []string{"usage", "used", "available", "rate", "latency"}

	names := make([]string, 0, count)
	for i := 0; i < count; i++ {
		prefix := prefixes[i%len(prefixes)]
		suffix := suffixes[i%len(suffixes)]
		names = append(names, fmt.Sprintf("%s.%s", prefix, suffix))
	}

	return names
}

// generateRandomValue creates realistic values based on metric type
func generateRandomValue(metricName string) float64 {
	// BREAKPOINT: Set breakpoint here to see value generation
	// DEBUG: Different metrics have different value ranges

	switch {
	case contains(metricName, "cpu"):
		return rand.Float64() * 100 // 0-100%
	case contains(metricName, "memory"):
		return rand.Float64() * 16384 // 0-16GB
	case contains(metricName, "disk"):
		return rand.Float64() * 1024 // 0-1TB
	case contains(metricName, "network"):
		return rand.Float64() * 1000 // 0-1000 Mbps
	case contains(metricName, "response_time"):
		return rand.Float64() * 500 // 0-500ms
	default:
		return rand.Float64() * 100
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
