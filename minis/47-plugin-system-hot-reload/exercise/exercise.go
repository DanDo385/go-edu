//go:build !solution
// +build !solution

package exercise

import (
	"github.com/example/go-10x-minis/minis/47-plugin-system-hot-reload/shared"
)

// LoadAndExecute loads a plugin from the given path, initializes it,
// executes it with the provided input, and returns the result.
//
// Steps:
// 1. Open the plugin file using plugin.Open()
// 2. Look up the "Plugin" symbol
// 3. Type assert the symbol to shared.Plugin
// 4. Call Init() on the plugin
// 5. Call Process() with the input
// 6. Return the result
//
// Parameters:
//   - pluginPath: Path to the .so file
//   - input: Data to pass to the plugin's Process() method
//
// Returns:
//   - output: Result from plugin.Process()
//   - error: Non-nil if any step fails
//
// Example:
//   result, err := LoadAndExecute("plugins/greeter.so", "Alice")
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Println(result) // "Hello, Alice! Great to see you!"
func LoadAndExecute(pluginPath string, input interface{}) (interface{}, error) {
	// TODO: Implement plugin loading and execution
	// Hint 1: Use plugin.Open(pluginPath) to load the .so file
	// Hint 2: Use p.Lookup("Plugin") to find the exported symbol
	// Hint 3: Type assert: plug, ok := sym.(shared.Plugin)
	// Hint 4: Call plug.Init() before plug.Process()
	return nil, nil
}

// DiscoverPlugins scans the given directory for .so files
// and returns a list of their absolute paths.
//
// This is useful for automatically finding all plugins in a directory.
//
// Parameters:
//   - dir: Directory to scan for .so files
//
// Returns:
//   - []string: List of absolute paths to .so files
//   - error: Non-nil if directory cannot be read
//
// Example:
//   plugins, err := DiscoverPlugins("./plugins")
//   // Returns: ["./plugins/greeter.so", "./plugins/math.so"]
func DiscoverPlugins(dir string) ([]string, error) {
	// TODO: Implement plugin discovery
	// Hint 1: Use filepath.Glob() or filepath.Walk()
	// Hint 2: Filter for files ending in ".so"
	// Hint 3: Return absolute paths (use filepath.Abs if needed)
	return nil, nil
}

// ReloadPlugin simulates reloading a plugin by loading it fresh.
//
// In a real implementation, you would:
// 1. Track loaded plugins in a map
// 2. Clean up the old version
// 3. Load the new version
// 4. Swap the old for the new
//
// For this exercise, simply load the plugin and return it.
//
// Parameters:
//   - pluginPath: Path to the .so file
//
// Returns:
//   - shared.Plugin: The loaded and initialized plugin
//   - error: Non-nil if loading or initialization fails
//
// Example:
//   plugin, err := ReloadPlugin("plugins/greeter.so")
//   fmt.Printf("Loaded %s v%s\n", plugin.Name(), plugin.Version())
func ReloadPlugin(pluginPath string) (shared.Plugin, error) {
	// TODO: Implement plugin reload
	// Hint 1: Open the plugin with plugin.Open()
	// Hint 2: Lookup and assert to shared.Plugin
	// Hint 3: Call Init() before returning
	return nil, nil
}
