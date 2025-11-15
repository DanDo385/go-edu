package exercise

import (
	"runtime"
	"strings"
	"testing"
)

func TestGetPathSeparator(t *testing.T) {
	sep := GetPathSeparator()

	switch runtime.GOOS {
	case "windows":
		if sep != "\\" {
			t.Errorf("Expected path separator '\\\\' on Windows, got %q", sep)
		}
	default: // Unix-like systems
		if sep != "/" {
			t.Errorf("Expected path separator '/' on Unix-like systems, got %q", sep)
		}
	}

	t.Logf("✓ Path separator for %s: %q", runtime.GOOS, sep)
}

func TestGetHomeDirectory(t *testing.T) {
	home := GetHomeDirectory()

	if home == "" {
		t.Error("Expected non-empty home directory")
	}

	// Basic validation: should be an absolute path
	switch runtime.GOOS {
	case "windows":
		// Windows paths typically start with C:\ or similar
		if !strings.Contains(home, ":") {
			t.Errorf("Expected Windows path to contain ':', got %q", home)
		}
	default:
		// Unix paths start with /
		if !strings.HasPrefix(home, "/") {
			t.Errorf("Expected Unix path to start with '/', got %q", home)
		}
	}

	t.Logf("✓ Home directory: %s", home)
}

func TestGetStorageBackend(t *testing.T) {
	backend := GetStorageBackend()

	// This test's expectations change based on build tags
	// We can't check the specific value without knowing build tags,
	// but we can verify it's not empty
	if backend == "" {
		t.Error("Expected non-empty storage backend")
	}

	validBackends := map[string]bool{
		"S3":               true,
		"Local Filesystem": true,
	}

	if !validBackends[backend] {
		t.Errorf("Expected storage backend to be 'S3' or 'Local Filesystem', got %q", backend)
	}

	t.Logf("✓ Storage backend: %s", backend)

	// Give hints about build tags
	if backend == "Local Filesystem" {
		t.Log("  (Run 'go test -tags=cloud' to test cloud storage)")
	} else {
		t.Log("  (Run 'go test' without tags to test local storage)")
	}
}

func TestGetWordSize(t *testing.T) {
	wordSize := GetWordSize()

	expectedSize := 0
	switch runtime.GOARCH {
	case "amd64", "arm64", "ppc64", "ppc64le", "mips64", "mips64le", "s390x":
		expectedSize = 64
	case "386", "arm", "mips", "mipsle":
		expectedSize = 32
	default:
		t.Logf("Unknown architecture %s, skipping word size validation", runtime.GOARCH)
		return
	}

	if wordSize != expectedSize {
		t.Errorf("Expected word size %d for %s, got %d", expectedSize, runtime.GOARCH, wordSize)
	}

	t.Logf("✓ Word size for %s: %d bits", runtime.GOARCH, wordSize)
}

func TestIsLinux64Bit(t *testing.T) {
	result := IsLinux64Bit()

	expected := false
	if runtime.GOOS == "linux" && (runtime.GOARCH == "amd64" || runtime.GOARCH == "arm64") {
		expected = true
	}

	if result != expected {
		t.Errorf("Expected IsLinux64Bit() = %v on %s/%s, got %v",
			expected, runtime.GOOS, runtime.GOARCH, result)
	}

	if expected {
		t.Logf("✓ Running on 64-bit Linux (%s)", runtime.GOARCH)
	} else {
		t.Logf("✓ Not running on 64-bit Linux (OS: %s, Arch: %s)", runtime.GOOS, runtime.GOARCH)
	}
}

func TestLoggingEnabled(t *testing.T) {
	enabled := IsLoggingEnabled()

	// We can't assert a specific value without knowing build tags
	// but we can test that it returns a boolean and that LogMessage works

	t.Logf("✓ Logging enabled: %v", enabled)

	// Test that LogMessage doesn't panic
	LogMessage("INFO", "Test message")
	LogMessage("ERROR", "Another test message")

	if enabled {
		t.Log("  (Debug logging is enabled)")
	} else {
		t.Log("  (Debug logging is disabled - run 'go test -tags=debug' to enable)")
	}
}

// TestBuildTags is a meta-test that reports which build configuration is active
func TestBuildTags(t *testing.T) {
	t.Log("=== Active Build Configuration ===")
	t.Logf("Operating System: %s", runtime.GOOS)
	t.Logf("Architecture: %s", runtime.GOARCH)
	t.Logf("Compiler: %s", runtime.Compiler)
	t.Logf("Go Version: %s", runtime.Version())

	t.Log("\n=== Build Tag Results ===")
	t.Logf("Path Separator: %q", GetPathSeparator())
	t.Logf("Home Directory: %s", GetHomeDirectory())
	t.Logf("Storage Backend: %s", GetStorageBackend())
	t.Logf("Word Size: %d bits", GetWordSize())
	t.Logf("Linux 64-bit: %v", IsLinux64Bit())
	t.Logf("Logging Enabled: %v", IsLoggingEnabled())

	t.Log("\n=== Try These Commands ===")
	t.Log("go test                          # Default build")
	t.Log("go test -tags=debug              # Enable debug logging")
	t.Log("go test -tags=cloud              # Use cloud storage")
	t.Log("go test -tags=\"debug,cloud\"      # Multiple tags")
	t.Log("GOOS=windows go test             # Test Windows behavior")
	t.Log("GOOS=darwin go test              # Test macOS behavior")
	t.Log("GOARCH=386 go test               # Test 32-bit architecture")
}

// Benchmark to verify no runtime overhead in production builds
func BenchmarkLogging(b *testing.B) {
	for i := 0; i < b.N; i++ {
		LogMessage("INFO", "Benchmark message")
	}
}
