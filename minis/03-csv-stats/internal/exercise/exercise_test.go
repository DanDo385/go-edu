package exercise

import (
	"strings"
	"testing"
)

func TestSummarizeCSV(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    map[string]Stat
		wantErr bool
	}{
		{
			name: "basic transactions",
			input: `id,category,amount
1,groceries,12.50
2,groceries,7.50
3,books,10.00
4,groceries,5.00`,
			want: map[string]Stat{
				"groceries": {Count: 3, Sum: 25.00, Avg: 8.333333333333334},
				"books":     {Count: 1, Sum: 10.00, Avg: 10.00},
			},
			wantErr: false,
		},
		{
			name: "single category",
			input: `id,category,amount
1,electronics,150.00
2,electronics,200.00`,
			want: map[string]Stat{
				"electronics": {Count: 2, Sum: 350.00, Avg: 175.00},
			},
			wantErr: false,
		},
		{
			name: "header only (no data)",
			input: `id,category,amount
`,
			want:    map[string]Stat{},
			wantErr: false,
		},
		{
			name:    "empty file",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid header",
			input: `foo,bar,baz
1,groceries,12.50`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid amount (not a number)",
			input: `id,category,amount
1,groceries,12.50
2,books,invalid`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty category",
			input: `id,category,amount
1,,12.50`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "wrong number of columns",
			input: `id,category,amount
1,groceries`,
			want:    nil,
			wantErr: true,
		},
		{
			name: "negative amounts",
			input: `id,category,amount
1,refunds,-10.00
2,refunds,-5.50`,
			want: map[string]Stat{
				"refunds": {Count: 2, Sum: -15.50, Avg: -7.75},
			},
			wantErr: false,
		},
		{
			name: "decimal precision",
			input: `id,category,amount
1,test,0.01
2,test,0.02
3,test,0.03`,
			want: map[string]Stat{
				"test": {Count: 3, Sum: 0.06, Avg: 0.02},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			got, err := SummarizeCSV(r)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("SummarizeCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expected an error, we're done
			if tt.wantErr {
				return
			}

			// Compare results
			if len(got) != len(tt.want) {
				t.Errorf("SummarizeCSV() got %d categories, want %d", len(got), len(tt.want))
				return
			}

			for category, wantStat := range tt.want {
				gotStat, exists := got[category]
				if !exists {
					t.Errorf("SummarizeCSV() missing category %q", category)
					continue
				}

				if gotStat.Count != wantStat.Count {
					t.Errorf("SummarizeCSV() category %q Count = %d, want %d",
						category, gotStat.Count, wantStat.Count)
				}

				// Use epsilon comparison for floats (avoid exact equality)
				const epsilon = 0.0001
				if abs(gotStat.Sum-wantStat.Sum) > epsilon {
					t.Errorf("SummarizeCSV() category %q Sum = %.4f, want %.4f",
						category, gotStat.Sum, wantStat.Sum)
				}

				if abs(gotStat.Avg-wantStat.Avg) > epsilon {
					t.Errorf("SummarizeCSV() category %q Avg = %.4f, want %.4f",
						category, gotStat.Avg, wantStat.Avg)
				}
			}
		})
	}
}

// Helper function for floating-point comparison
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
