package main

import (
	"fmt"
	"runtime"
)

func main() {
	printHeader("Build Tags & Conditional Compilation Demo")

	// Section 1: Platform-specific code
	printSection("1. Platform-Specific Code")
	fmt.Println("Platform:", GetPlatformName())
	fmt.Println("Current User:", GetUsername())
	fmt.Println("\nDetailed System Information:")
	fmt.Print(GetSystemInfo())
	fmt.Println(GetSpecialFeature())

	// Section 2: Architecture-specific code
	printSection("2. Architecture-Specific Optimizations")
	fmt.Print(GetArchitectureInfo())
	fmt.Println("\nPerformance Test: Vector Sum")
	testData := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}
	sum := VectorSum(testData)
	fmt.Printf("Sum of %v = %.2f\n", testData, sum)
	fmt.Println(GetOptimizationHint())

	// Section 3: Debug features
	printSection("3. Debug Features")
	fmt.Println("Debug Mode Enabled:", IsDebugEnabled())
	fmt.Println(GetDebugInfo())
	if IsDebugEnabled() {
		EnableVerboseLogging()
		LogDebug("main", "This is a debug message")
		LogDebug("main", "Debug logging includes file and line numbers")
		DumpMemoryStats()
	} else {
		fmt.Println("\n(Build with -tags=debug to enable debug features)")
	}

	// Section 4: Premium features
	printSection("4. Feature Flags (Premium Tier)")
	fmt.Println("Premium Enabled:", IsPremiumEnabled())
	fmt.Println(GetPremiumFeatures())
	fmt.Println("\nFeature Limits:")
	fmt.Printf("  Max Concurrent Users: %d\n", GetMaxConcurrentUsers())
	fmt.Printf("  API Rate Limit: %d requests/hour\n", GetAPIRateLimit())
	fmt.Println()
	EnableAdvancedAnalytics()
	fmt.Println()
	EnablePrioritySupport()
	if !IsPremiumEnabled() {
		fmt.Println("\n(Build with -tags=premium to enable premium features)")
	}

	// Section 5: Advanced constraints (only on certain platforms)
	printSection("5. Advanced Build Constraints")
	if (runtime.GOOS == "linux" || runtime.GOOS == "darwin") && runtime.GOARCH == "amd64" {
		// This function only exists when complex constraints are met
		// ShowAdvancedConstraints()
		fmt.Println("Complex constraint example:")
		fmt.Println("  advanced_constraints.go is compiled on Unix-like AMD64 systems without CGO")
	} else {
		fmt.Println("Advanced constraints not met on this platform")
		fmt.Printf("  Current: OS=%s, Arch=%s\n", runtime.GOOS, runtime.GOARCH)
	}

	// Section 6: Build instructions
	printSection("6. How to Build with Different Tags")
	fmt.Println("Try building this program with different build tags:")
	fmt.Println()
	fmt.Println("Default build (production, free tier):")
	fmt.Println("  go run cmd/demo/main.go")
	fmt.Println()
	fmt.Println("Debug build:")
	fmt.Println("  go run -tags=debug cmd/demo/main.go")
	fmt.Println()
	fmt.Println("Premium features:")
	fmt.Println("  go run -tags=premium cmd/demo/main.go")
	fmt.Println()
	fmt.Println("Debug + Premium:")
	fmt.Println("  go run -tags=\"debug,premium\" cmd/demo/main.go")
	fmt.Println()
	fmt.Println("Cross-compile for different platforms:")
	fmt.Println("  GOOS=windows go build cmd/demo/main.go")
	fmt.Println("  GOOS=darwin go build cmd/demo/main.go")
	fmt.Println("  GOOS=linux GOARCH=arm64 go build cmd/demo/main.go")
	fmt.Println()
	fmt.Println("See which files get compiled:")
	fmt.Println("  go list -f '{{.GoFiles}}'")
	fmt.Println("  go list -tags=debug -f '{{.GoFiles}}'")

	printFooter()
}

func printHeader(title string) {
	line := "═══════════════════════════════════════════════════════════════════════"
	fmt.Println(line)
	fmt.Printf("  %s\n", title)
	fmt.Println(line)
	fmt.Println()
}

func printSection(title string) {
	fmt.Println()
	fmt.Println("───────────────────────────────────────────────────────────────────────")
	fmt.Println(title)
	fmt.Println("───────────────────────────────────────────────────────────────────────")
}

func printFooter() {
	fmt.Println()
	line := "═══════════════════════════════════════════════════════════════════════"
	fmt.Println(line)
	fmt.Println("  Build tags allow you to conditionally compile code based on:")
	fmt.Println("    • Operating System (linux, darwin, windows)")
	fmt.Println("    • Architecture (amd64, arm64, 386)")
	fmt.Println("    • Custom features (debug, premium, etc.)")
	fmt.Println("    • Go version, CGO status, and more")
	fmt.Println()
	fmt.Println("  Benefits:")
	fmt.Println("    ✓ Smaller binaries (unused code not compiled)")
	fmt.Println("    ✓ No runtime overhead")
	fmt.Println("    ✓ Platform-specific optimizations")
	fmt.Println("    ✓ Feature flag management")
	fmt.Println(line)
}
