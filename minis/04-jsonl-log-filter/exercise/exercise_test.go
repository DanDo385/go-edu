package exercise

import (
	"strings"
	"testing"
	"time"
)

func TestFilterLogs(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		minLevel     Level
		wantCount    int
		wantFirstMsg string
		wantErr      bool
	}{
		{
			name: "filter warn and above",
			input: `{"ts":"2024-01-01T12:00:00Z","level":"info","msg":"Server started"}
{"ts":"2024-01-01T12:00:05Z","level":"error","msg":"Database failed"}
{"ts":"2024-01-01T12:00:02Z","level":"debug","msg":"Processing request"}
{"ts":"2024-01-01T12:00:10Z","level":"warn","msg":"High memory"}`,
			minLevel:     Warn,
			wantCount:    2,
			wantFirstMsg: "Database failed", // Sorted by timestamp
			wantErr:      false,
		},
		{
			name: "all levels",
			input: `{"ts":"2024-01-01T12:00:00Z","level":"info","msg":"A"}
{"ts":"2024-01-01T12:00:01Z","level":"debug","msg":"B"}`,
			minLevel:     Debug,
			wantCount:    2,
			wantFirstMsg: "A",
			wantErr:      false,
		},
		{
			name:      "empty input",
			input:     "",
			minLevel:  Info,
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "malformed json (skip line)",
			input: `{"ts":"2024-01-01T12:00:00Z","level":"error","msg":"A"}
{this is not json
{"ts":"2024-01-01T12:00:01Z","level":"error","msg":"B"}`,
			minLevel:     Error,
			wantCount:    2,
			wantFirstMsg: "A",
			wantErr:      true, // Should report skipped lines
		},
		{
			name: "all filtered out",
			input: `{"ts":"2024-01-01T12:00:00Z","level":"debug","msg":"A"}
{"ts":"2024-01-01T12:00:01Z","level":"info","msg":"B"}`,
			minLevel:  Error,
			wantCount: 0,
			wantErr:   false,
		},
		{
			name: "sorting by timestamp",
			input: `{"ts":"2024-01-01T12:00:05Z","level":"error","msg":"Third"}
{"ts":"2024-01-01T12:00:01Z","level":"error","msg":"First"}
{"ts":"2024-01-01T12:00:03Z","level":"error","msg":"Second"}`,
			minLevel:     Error,
			wantCount:    3,
			wantFirstMsg: "First", // Should be sorted
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			got, err := FilterLogs(r, tt.minLevel)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterLogs() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check count
			if len(got) != tt.wantCount {
				t.Errorf("FilterLogs() returned %d entries, want %d", len(got), tt.wantCount)
				return
			}

			// Check first message (if any)
			if tt.wantCount > 0 && got[0].Msg != tt.wantFirstMsg {
				t.Errorf("FilterLogs() first message = %q, want %q", got[0].Msg, tt.wantFirstMsg)
			}
		})
	}
}

func TestLevel_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    Level
		wantErr bool
	}{
		{
			name:    "debug",
			json:    `"debug"`,
			want:    Debug,
			wantErr: false,
		},
		{
			name:    "info",
			json:    `"info"`,
			want:    Info,
			wantErr: false,
		},
		{
			name:    "warn",
			json:    `"warn"`,
			want:    Warn,
			wantErr: false,
		},
		{
			name:    "error",
			json:    `"error"`,
			want:    Error,
			wantErr: false,
		},
		{
			name:    "case insensitive",
			json:    `"ERROR"`,
			want:    Error,
			wantErr: false,
		},
		{
			name:    "invalid level",
			json:    `"critical"`,
			want:    0,
			wantErr: true,
		},
		{
			name:    "not a string",
			json:    `123`,
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Level
			err := got.UnmarshalJSON([]byte(tt.json))

			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("UnmarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterLogs_TimestampParsing(t *testing.T) {
	input := `{"ts":"2024-01-01T12:00:00Z","level":"info","msg":"UTC"}
{"ts":"2024-01-01T12:00:00-05:00","level":"info","msg":"EST"}`

	r := strings.NewReader(input)
	entries, err := FilterLogs(r, Debug)

	if err != nil {
		t.Fatalf("FilterLogs() error = %v, want nil", err)
	}

	if len(entries) != 2 {
		t.Fatalf("FilterLogs() returned %d entries, want 2", len(entries))
	}

	// Check that timestamps are parsed correctly
	utcTime, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00Z")
	estTime, _ := time.Parse(time.RFC3339, "2024-01-01T12:00:00-05:00")

	if !entries[0].TS.Equal(utcTime) {
		t.Errorf("First entry timestamp = %v, want %v", entries[0].TS, utcTime)
	}

	if !entries[1].TS.Equal(estTime) {
		t.Errorf("Second entry timestamp = %v, want %v", entries[1].TS, estTime)
	}
}
