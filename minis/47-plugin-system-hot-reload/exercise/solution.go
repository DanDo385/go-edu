//go:build solution
// +build solution

/*
Problem: Dynamic plugin loading and execution

We need to:
1. Load Go plugins (.so files) dynamically at runtime
2. Look up exported symbols from the loaded plugins
3. Type assert symbols to the expected interface
4. Initialize and execute plugins
5. Discover all plugins in a directory
6. Support reloading plugins

Constraints:
- Plugins must be built with the same Go version as the host
- Plugins must export a symbol named "Plugin"
- The symbol must implement the shared.Plugin interface
- Only works on Linux and macOS (not Windows)

Time/Space Complexity:
- plugin.Open(): O(n) where n = size of .so file (disk I/O + linking)
- Lookup(): O(1) - hash table lookup in symbol table
- Type assertion: O(1)
- Space: O(m) where m = number of loaded plugins in memory

Why Go plugins are powerful:
- No need to recompile the host application
- Extend functionality at runtime
- Isolate plugin code from host code
- Enable hot reloading without process restart

Challenges:
- Platform-specific (Linux/macOS only)
- Version compatibility (Go version must match exactly)
- Type identity (shared types must be from same package)
- Memory management (plugins can't be unloaded from memory)
*/

package exercise

import (
	"fmt"
	"path/filepath"
	"plugin"

	"github.com/example/go-10x-minis/minis/47-plugin-system-hot-reload/shared"
)

// LoadAndExecute loads a plugin, initializes it, and executes it.
//
// Go Concepts Demonstrated:
// - plugin.Open(): Dynamic library loading
// - Symbol lookup: Finding exported variables/functions
// - Type assertion: Runtime type checking
// - Interface satisfaction: Verifying plugin implements interface
//
// Three-Input Iteration Table:
//
// Input 1: Valid plugin, valid input
//   1. Open("greeter.so") → plugin loaded
//   2. Lookup("Plugin") → symbol found
//   3. Assert to shared.Plugin → success
//   4. Init() → plugin initialized
//   5. Process("Alice") → "Hello, Alice!"
//   Result: greeting string, nil error
//
// Input 2: Invalid path
//   1. Open("/invalid.so") → error
//   Result: nil, "plugin.Open: no such file"
//
// Input 3: Plugin missing symbol
//   1. Open("bad.so") → plugin loaded
//   2. Lookup("Plugin") → error
//   Result: nil, "plugin: symbol Plugin not found"
func LoadAndExecute(pluginPath string, input interface{}) (interface{}, error) {
	// Step 1: Open the plugin file
	// This loads the .so file into the process's memory space
	// The dynamic linker resolves dependencies and initializes globals
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin %s: %w", pluginPath, err)
	}

	// Step 2: Look up the exported symbol
	// Plugins must export a variable named "Plugin" (case-sensitive)
	// The symbol table is searched for this name
	sym, err := p.Lookup("Plugin")
	if err != nil {
		return nil, fmt.Errorf("failed to find Plugin symbol in %s: %w", pluginPath, err)
	}

	// Step 3: Type assert to the Plugin interface
	// This verifies that the symbol implements shared.Plugin
	// If the type doesn't match, ok will be false
	plug, ok := sym.(shared.Plugin)
	if !ok {
		return nil, fmt.Errorf("symbol 'Plugin' is not of type shared.Plugin (got %T)", sym)
	}

	// Step 4: Initialize the plugin
	// This gives the plugin a chance to set up resources, validate config, etc.
	if err := plug.Init(); err != nil {
		return nil, &shared.PluginError{
			PluginName: plug.Name(),
			Operation:  "init",
			Err:        err,
		}
	}

	// Step 5: Execute the plugin with the input
	result, err := plug.Process(input)
	if err != nil {
		return nil, &shared.PluginError{
			PluginName: plug.Name(),
			Operation:  "process",
			Err:        err,
		}
	}

	return result, nil
}

// DiscoverPlugins finds all .so files in a directory.
//
// Go Concepts Demonstrated:
// - filepath.Glob(): Pattern-based file matching
// - filepath.Walk(): Recursive directory traversal (alternative approach)
// - String manipulation: Filtering by extension
//
// Three-Input Iteration Table:
//
// Input 1: Directory with 2 .so files and 1 .txt file
//   1. Glob("dir/*.so") → ["dir/a.so", "dir/b.so"]
//   2. Filter .so files → keep both
//   3. Return list
//   Result: ["dir/a.so", "dir/b.so"], nil
//
// Input 2: Empty directory
//   1. Glob("empty/*.so") → []
//   Result: [], nil
//
// Input 3: Invalid directory
//   1. Glob("/invalid/*.so") → error
//   Result: nil, "no such file or directory"
func DiscoverPlugins(dir string) ([]string, error) {
	// Use filepath.Glob to find all .so files
	// Pattern: "dir/*.so" matches all files ending in .so in the directory
	pattern := filepath.Join(dir, "*.so")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to glob %s: %w", pattern, err)
	}

	// matches is already a slice of absolute/relative paths
	// No additional filtering needed since glob handles the .so extension
	return matches, nil
}

// ReloadPlugin loads a fresh instance of a plugin.
//
// Go Concepts Demonstrated:
// - Plugin lifecycle: Open → Lookup → Init
// - Error handling at each step
// - Returning interface types
//
// Note on actual hot reload:
// Go plugins cannot be truly "unloaded" from memory. The plugin.Open()
// function will return the same plugin instance if called multiple times
// with the same path. For true hot reload, you need to:
// 1. Build new plugin with different filename (e.g., plugin-v2.so)
// 2. Load the new file
// 3. Swap references from old to new
//
// Three-Input Iteration Table:
//
// Input 1: Valid plugin file
//   1. Open("greeter.so") → success
//   2. Lookup("Plugin") → found
//   3. Assert type → shared.Plugin
//   4. Init() → initialized
//   Result: plugin instance, nil
//
// Input 2: File exists but isn't a plugin
//   1. Open("notaplugin.txt") → error
//   Result: nil, "not a Go plugin"
//
// Input 3: Plugin with Init() error
//   1. Open("badplugin.so") → success
//   2. Lookup("Plugin") → found
//   3. Assert type → success
//   4. Init() → error
//   Result: nil, PluginError{operation: "init"}
func ReloadPlugin(pluginPath string) (shared.Plugin, error) {
	// Open the plugin
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin: %w", err)
	}

	// Look up the Plugin symbol
	sym, err := p.Lookup("Plugin")
	if err != nil {
		return nil, fmt.Errorf("failed to lookup Plugin symbol: %w", err)
	}

	// Type assert to shared.Plugin interface
	plug, ok := sym.(shared.Plugin)
	if !ok {
		return nil, fmt.Errorf("Plugin symbol has wrong type: %T", sym)
	}

	// Initialize the plugin
	if err := plug.Init(); err != nil {
		return nil, &shared.PluginError{
			PluginName: plug.Name(),
			Operation:  "init",
			Err:        err,
		}
	}

	return plug, nil
}

/*
Alternatives & Trade-offs:

1. Static linking (compile plugins into binary):
   var plugins = map[string]Plugin{
       "greeter": GreeterPlugin{},
       "math": MathPlugin{},
   }
   Pros: Simpler, no platform issues, type-safe at compile time
   Cons: Must recompile to add plugins, larger binary, no hot reload

2. RPC-based plugins (HashiCorp go-plugin):
   Plugin runs as separate process, communicates via gRPC
   Pros: Works on Windows, better isolation, true hot reload
   Cons: IPC overhead, more complex, requires protobuf definitions

3. Embedded scripting (Lua, JavaScript via goja):
   Execute scripts in an embedded interpreter
   Pros: Works on all platforms, sandboxed, dynamic
   Cons: Different language, slower, limited Go interop

4. WebAssembly plugins (WASM):
   Compile plugins to WASM, load via runtime (wazero, wasmer)
   Pros: Platform-independent, sandboxed, language-agnostic
   Cons: Limited Go support, performance overhead, new tech

5. Shared libraries with C API (cgo):
   Build plugins with C-compatible API, load with dlopen
   Pros: Language-agnostic (can use C/C++/Rust plugins)
   Cons: Lose Go type safety, complex FFI, manual memory management

Go vs X:

Go vs Python (importlib):
  import importlib
  plugin = importlib.import_module("plugins.greeter")
  result = plugin.greet("Alice")
  Pros: Simpler (no compilation), no type assertions
  Cons: Runtime overhead, no compile-time checks, GIL limits parallelism
  Go: Compiled plugins are faster, type-safe, true parallelism

Go vs Java (ServiceLoader):
  ServiceLoader<Plugin> loader = ServiceLoader.load(Plugin.class);
  for (Plugin plugin : loader) {
      plugin.execute();
  }
  Pros: Clean API, classpath-based loading
  Cons: Verbose, reflection overhead, classpath complexity
  Go: Simpler build process, better performance

Go vs Rust (libloading):
  let lib = libloading::Library::new("plugin.so")?;
  let func: Symbol<fn() -> i32> = lib.get(b"my_func")?;
  Pros: Zero-cost abstractions, memory safety guarantees
  Cons: Complex type system (need extern "C" for ABI), steeper learning curve
  Go: Easier to use, no unsafe code needed, simpler ABI

Go vs C (dlopen):
  void *handle = dlopen("plugin.so", RTLD_LAZY);
  void *symbol = dlsym(handle, "Plugin");
  Pros: Maximum control, works everywhere, no runtime
  Cons: Manual memory management, no type safety, error-prone
  Go: Type-safe, memory-safe, cleaner API
*/
