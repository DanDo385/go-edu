package main

import (
	"fmt"
	"strings"

	"github.com/example/go-10x-minis/minis/04-jsonl-log-filter/internal/jsonllogfilter"
)

/*
Debug Harness for JSONL Log Filter Module

This file runs the log filter with sample data.
Set breakpoints in exercise.go and press F5 in VS Code.
*/
func main() {
	fmt.Println("=== JSONL Log Filter Debug Harness ===")
	fmt.Println()

	// Sample JSONL data
	jsonlData := `{"ts":"2024-01-01T10:00:00Z","level":"debug","msg":"Starting application"}
{"ts":"2024-01-01T10:00:01Z","level":"info","msg":"Server listening on :8080"}
{"ts":"2024-01-01T10:00:02Z","level":"warn","msg":"Connection pool near capacity"}
{"ts":"2024-01-01T10:00:03Z","level":"error","msg":"Database connection failed"}
{"ts":"2024-01-01T10:00:04Z","level":"info","msg":"Retrying connection"}
{"ts":"2024-01-01T10:00:05Z","level":"error","msg":"Max retries exceeded"}`

	fmt.Println("--- Sample JSONL Data ---")
	fmt.Println(jsonlData)
	fmt.Println()

	// Filter at different levels
	levels := []struct {
		name  string
		level jsonllogfilter.Level
	}{
		{"debug", jsonllogfilter.LevelDebug},
		{"warn", jsonllogfilter.LevelWarn},
		{"error", jsonllogfilter.LevelError},
	}

	for _, l := range levels {
		fmt.Printf("--- Filtering at level: %s ---\n", l.name)
		reader := strings.NewReader(jsonlData)
		entries, err := jsonllogfilter.FilterLogs(reader, l.level)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Printf("Found %d entries:\n", len(entries))
		for _, entry := range entries {
			fmt.Printf("  [%s] %s\n", entry.TS.Format("15:04:05"), entry.Msg)
		}
		fmt.Println()
	}

	fmt.Println("=== What You Learned ===")
	fmt.Println("1. JSONL = one JSON object per line")
	fmt.Println("2. Custom UnmarshalJSON for enum types")
	fmt.Println("3. sort.Slice for custom ordering")
	fmt.Println("4. time.Time parses RFC3339 automatically")
	fmt.Println()
	fmt.Println("Next: Proceed to minis/05-cli-todo-files")
}
