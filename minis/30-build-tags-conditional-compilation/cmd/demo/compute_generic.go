//go:build !amd64 && !arm64

package main

import "runtime"

// GetArchitectureInfo returns generic architecture information
func GetArchitectureInfo() string {
	return "Architecture: " + runtime.GOARCH + "\n" +
		"Optimizations: Using generic Go code (no arch-specific optimizations)\n" +
		"Note: This is a fallback for less common architectures\n"
}

// VectorSum provides a generic implementation
func VectorSum(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum
}

// GetOptimizationHint returns generic optimization advice
func GetOptimizationHint() string {
	return "Generic: Focus on algorithmic optimizations and efficient memory usage"
}
