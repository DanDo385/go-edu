package exercise

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestLoadAndExecute tests basic plugin loading and execution.
func TestLoadAndExecute(t *testing.T) {
	// Skip if not on Linux/macOS
	if !isPluginSupported() {
		t.Skip("Plugin system not supported on this platform")
	}

	// Build test plugin
	pluginPath := buildTestPlugin(t, "greeter")
	defer os.Remove(pluginPath)

	// Test loading and executing
	result, err := LoadAndExecute(pluginPath, "TestUser")
	if err != nil {
		t.Fatalf("LoadAndExecute failed: %v", err)
	}

	// Verify result is a string
	greeting, ok := result.(string)
	if !ok {
		t.Fatalf("Expected string result, got %T", result)
	}

	// Verify greeting contains the name
	if !strings.Contains(greeting, "TestUser") {
		t.Errorf("Expected greeting to contain 'TestUser', got: %s", greeting)
	}

	t.Logf("Successfully loaded and executed plugin: %s", greeting)
}

// TestLoadAndExecute_InvalidPath tests error handling for invalid paths.
func TestLoadAndExecute_InvalidPath(t *testing.T) {
	if !isPluginSupported() {
		t.Skip("Plugin system not supported on this platform")
	}

	_, err := LoadAndExecute("/nonexistent/plugin.so", "test")
	if err == nil {
		t.Fatal("Expected error for nonexistent plugin, got nil")
	}

	t.Logf("Correctly returned error: %v", err)
}

// TestDiscoverPlugins tests plugin discovery in a directory.
func TestDiscoverPlugins(t *testing.T) {
	if !isPluginSupported() {
		t.Skip("Plugin system not supported on this platform")
	}

	// Create temp directory
	tmpDir := t.TempDir()

	// Build multiple test plugins
	plugin1 := filepath.Join(tmpDir, "greeter.so")
	plugin2 := filepath.Join(tmpDir, "math.so")

	buildTestPluginTo(t, "greeter", plugin1)
	buildTestPluginTo(t, "math", plugin2)

	// Also create a non-.so file (should be ignored)
	txtFile := filepath.Join(tmpDir, "readme.txt")
	if err := os.WriteFile(txtFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	// Discover plugins
	plugins, err := DiscoverPlugins(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverPlugins failed: %v", err)
	}

	// Verify we found exactly 2 plugins
	if len(plugins) != 2 {
		t.Fatalf("Expected 2 plugins, found %d: %v", len(plugins), plugins)
	}

	// Verify all are .so files
	for _, p := range plugins {
		if !strings.HasSuffix(p, ".so") {
			t.Errorf("Non-.so file in results: %s", p)
		}
	}

	t.Logf("Successfully discovered %d plugins: %v", len(plugins), plugins)
}

// TestDiscoverPlugins_EmptyDir tests discovery in an empty directory.
func TestDiscoverPlugins_EmptyDir(t *testing.T) {
	if !isPluginSupported() {
		t.Skip("Plugin system not supported on this platform")
	}

	tmpDir := t.TempDir()

	plugins, err := DiscoverPlugins(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverPlugins failed: %v", err)
	}

	if len(plugins) != 0 {
		t.Errorf("Expected 0 plugins in empty directory, found %d", len(plugins))
	}
}

// TestReloadPlugin tests plugin reloading.
func TestReloadPlugin(t *testing.T) {
	if !isPluginSupported() {
		t.Skip("Plugin system not supported on this platform")
	}

	// Build test plugin
	pluginPath := buildTestPlugin(t, "greeter")
	defer os.Remove(pluginPath)

	// Load plugin
	plugin, err := ReloadPlugin(pluginPath)
	if err != nil {
		t.Fatalf("ReloadPlugin failed: %v", err)
	}

	// Verify plugin interface
	if plugin.Name() == "" {
		t.Error("Plugin name is empty")
	}
	if plugin.Version() == "" {
		t.Error("Plugin version is empty")
	}

	// Try to use the plugin
	result, err := plugin.Process("TestUser")
	if err != nil {
		t.Fatalf("Plugin.Process failed: %v", err)
	}

	if result == nil {
		t.Error("Plugin.Process returned nil result")
	}

	t.Logf("Successfully reloaded plugin: %s v%s", plugin.Name(), plugin.Version())
}

// Helper: Build a test plugin and return its path
func buildTestPlugin(t *testing.T, name string) string {
	t.Helper()

	tmpDir := t.TempDir()
	pluginPath := filepath.Join(tmpDir, name+".so")
	buildTestPluginTo(t, name, pluginPath)
	return pluginPath
}

// Helper: Build a test plugin to a specific path
func buildTestPluginTo(t *testing.T, name, outputPath string) {
	t.Helper()

	// Find the source file
	sourcePath := filepath.Join("..", "plugins", name, name+".go")
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		t.Skipf("Plugin source not found: %s", sourcePath)
	}

	// Build the plugin
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", outputPath, sourcePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build plugin %s: %v\nOutput: %s", name, err, output)
	}
}

// Helper: Check if plugin system is supported on this platform
func isPluginSupported() bool {
	// Plugins are only supported on Linux and macOS
	// Not supported on Windows
	cmd := exec.Command("go", "env", "GOOS")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	os := strings.TrimSpace(string(output))
	return os == "linux" || os == "darwin"
}

// BenchmarkPluginLoad benchmarks plugin loading performance.
func BenchmarkPluginLoad(b *testing.B) {
	if !isPluginSupported() {
		b.Skip("Plugin system not supported on this platform")
	}

	// Build test plugin once
	pluginPath := buildTestPlugin(&testing.T{}, "greeter")
	defer os.Remove(pluginPath)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ReloadPlugin(pluginPath)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPluginExecute benchmarks plugin execution performance.
func BenchmarkPluginExecute(b *testing.B) {
	if !isPluginSupported() {
		b.Skip("Plugin system not supported on this platform")
	}

	// Build and load plugin once
	pluginPath := buildTestPlugin(&testing.T{}, "greeter")
	defer os.Remove(pluginPath)

	plugin, err := ReloadPlugin(pluginPath)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := plugin.Process("BenchUser")
		if err != nil {
			b.Fatal(err)
		}
	}
}
