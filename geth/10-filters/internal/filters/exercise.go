//go:build !solution && !reference

package filters

import (
	"context"
	"time"
)

const defaultMaxHeads = 5
const defaultPollInterval = time.Second

// Run - TODO: implement this function
func Run(ctx context.Context, client HeadClient, cfg Config) (*Result, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil, nil
}

// subscribeHeads - TODO: implement this function
func subscribeHeads(ctx context.Context, client HeadClient, cfg Config) (*Result, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil, nil
}

// pollHeads - TODO: implement this function
func pollHeads(ctx context.Context, client HeadClient, cfg Config) (*Result, error) {
	// TODO: Implement this function
	// Refer to solution.reference.go for the complete implementation with detailed explanations
	return nil, nil
}
