# Project 47: Plugin System with Hot Reload

## 1. What Is This About?

### Real-World Scenario

Imagine you're building a data processing platform that needs to support custom transformations:

**❌ Bad approach:** Hardcode all transformations
- Every new feature requires recompiling the entire application
- Need to redeploy to add new transformations
- All customers get all features (binary bloat)
- Can't customize per customer without code branches

**✅ Better approach:** Use a **plugin system**
- Load transformations dynamically at runtime
- Add new features without recompiling core application
- Hot reload: Update plugins without restarting the service
- Each customer can have custom plugins
- Core stays stable, plugins can evolve independently

This project teaches you how to build extensible systems where functionality can be added, modified, or removed **without touching the core application**.

### What You'll Learn

1. **Go Plugin System**: Dynamic loading of compiled Go code (.so files)
2. **Plugin Interfaces**: Designing stable contracts between core and plugins
3. **Hot Reload**: Detecting changes and reloading plugins at runtime
4. **File Watching**: Monitor filesystem for plugin updates
5. **Symbol Loading**: Looking up functions and variables from shared libraries
6. **Version Management**: Handling multiple plugin versions safely

### The Challenge

Build a plugin system that:
- Loads plugins dynamically from .so files
- Defines a clear interface that plugins must implement
- Watches for plugin file changes
- Hot reloads plugins without restarting the application
- Handles errors gracefully (bad plugins shouldn't crash the host)
- Supports multiple concurrent plugins

---

## 2. First Principles: Understanding Go Plugins

### What is a Plugin?

A **plugin** is code that extends an application's functionality without being compiled into the main binary.

**Analogy**: Think of plugins like power tool attachments:
- **Main application** = Power drill (the core motor)
- **Plugins** = Different drill bits, sanders, polishers (extensions)
- **Plugin interface** = The standard chuck that accepts any compatible bit
- **Hot reload** = Swapping bits without turning off the drill

### How Do Go Plugins Work?

Go's `plugin` package allows you to:
1. Build Go code as a **shared library** (.so on Linux, .dylib on macOS)
2. Load that library at **runtime** (not compile time)
3. Look up **exported symbols** (functions, variables)
4. Call those functions as if they were native

**Key insight**: Plugins are full Go code, compiled separately, loaded dynamically.

```go
// Building a plugin
// $ go build -buildmode=plugin -o greeter.so greeter.go

// Loading a plugin
p, err := plugin.Open("greeter.so")
if err != nil {
    log.Fatal(err)
}

// Looking up a symbol
sym, err := p.Lookup("Greet")
if err != nil {
    log.Fatal(err)
}

// Type assert and use
greetFunc := sym.(func(string) string)
message := greetFunc("World")
```

### Plugin System Architecture

```
┌─────────────────────────────────────────────────────┐
│                  Host Application                   │
│                                                     │
│  ┌──────────────┐        ┌──────────────┐         │
│  │Plugin Manager│◄──────►│ File Watcher │         │
│  └──────┬───────┘        └──────────────┘         │
│         │                                           │
│         │ Load/Reload                              │
│         ▼                                           │
│  ┌─────────────────────────────────┐               │
│  │     Plugin Interface (API)       │               │
│  └─────────────────────────────────┘               │
└──────────────┬────────────┬─────────────────────────┘
               │            │
       .so file│    .so file│
               ▼            ▼
      ┌───────────┐  ┌───────────┐
      │ Plugin A  │  │ Plugin B  │
      │ (v1.0.0)  │  │ (v2.1.3)  │
      └───────────┘  └───────────┘
```

### Platform Compatibility

**Important**: Go plugins have platform restrictions:

| Platform | Support | File Extension |
|----------|---------|----------------|
| Linux    | ✅ Yes  | .so            |
| macOS    | ✅ Yes  | .so (or .dylib)|
| Windows  | ❌ No   | -              |
| FreeBSD  | ⚠️  Limited | .so        |

**Why no Windows?**
- Go plugins use C shared library mechanisms
- Windows uses DLLs with different semantics
- Go team hasn't implemented Windows support (as of Go 1.21)

**Workarounds for cross-platform**:
- Use **gRPC plugins** (separate processes communicating via RPC)
- Use **WASM plugins** (compile plugins to WebAssembly)
- Use **embedded interpreters** (Lua, JavaScript via V8)

### What is Hot Reload?

**Hot reload** = Updating code while the application is running, without restart.

**Why it matters**:
- **Zero downtime**: No service interruption
- **Faster iteration**: Change plugin logic, see results immediately
- **Production updates**: Deploy new features without bringing down the service

**Challenges**:
1. **File locking**: Can't overwrite a loaded .so file (OS locks it)
2. **Symbol conflicts**: Loading same plugin twice causes issues
3. **State management**: What happens to plugin state during reload?
4. **Type compatibility**: New version must match interface

**Solution strategy**:
```
1. Watch plugin directory for changes
2. When .so file changes, load it as a NEW plugin (different name)
3. Gracefully shut down old plugin instance
4. Swap to new plugin instance
5. Clean up old plugin (if possible)
```

---

## 3. Breaking Down the Solution

### Step 1: Define the Plugin Interface

Every plugin system needs a **stable contract** that plugins must implement.

```go
// Shared between host and plugins
package shared

type Plugin interface {
    // Name returns the plugin identifier
    Name() string

    // Version returns the plugin version
    Version() string

    // Process performs the plugin's main function
    Process(input interface{}) (interface{}, error)
}
```

**Key principle**: The interface must be:
- **Stable**: Rarely changes (breaking changes break all plugins)
- **Versioned**: Include version in interface or metadata
- **Simple**: Complex interfaces are hard to implement correctly
- **Documented**: Clear contracts prevent bugs

### Step 2: Build Plugins

Plugins are regular Go code with a special build mode:

```go
// plugins/greeter/greeter.go
package main

import "fmt"

// Must export a variable that implements the interface
var Plugin GreeterPlugin

type GreeterPlugin struct{}

func (GreeterPlugin) Name() string {
    return "greeter"
}

func (GreeterPlugin) Version() string {
    return "1.0.0"
}

func (GreeterPlugin) Process(input interface{}) (interface{}, error) {
    name, ok := input.(string)
    if !ok {
        return nil, fmt.Errorf("expected string input")
    }
    return fmt.Sprintf("Hello, %s!", name), nil
}
```

**Build command**:
```bash
go build -buildmode=plugin -o plugins/greeter.so plugins/greeter/greeter.go
```

**Critical requirements**:
- Use `package main` (plugins are executables)
- Export at least one symbol (variable or function)
- Build with exact same Go version as host
- Use same dependency versions as host

### Step 3: Load Plugins at Runtime

```go
func LoadPlugin(path string) (*plugin.Plugin, error) {
    // Open the plugin file
    p, err := plugin.Open(path)
    if err != nil {
        return nil, fmt.Errorf("failed to open plugin: %w", err)
    }

    // Look up the exported symbol
    sym, err := p.Lookup("Plugin")
    if err != nil {
        return nil, fmt.Errorf("plugin missing 'Plugin' symbol: %w", err)
    }

    // Type assert to the interface
    plug, ok := sym.(shared.Plugin)
    if !ok {
        return nil, fmt.Errorf("plugin doesn't implement Plugin interface")
    }

    return plug, nil
}
```

### Step 4: Implement File Watching

To enable hot reload, we need to detect when plugin files change:

```go
func WatchDirectory(dir string, onChange func(path string)) error {
    // Use fsnotify or poll the directory
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return err
    }

    err = watcher.Add(dir)
    if err != nil {
        return err
    }

    go func() {
        for {
            select {
            case event := <-watcher.Events:
                if event.Op&fsnotify.Write == fsnotify.Write {
                    onChange(event.Name)
                }
            case err := <-watcher.Errors:
                log.Printf("watcher error: %v", err)
            }
        }
    }()

    return nil
}
```

### Step 5: Implement Hot Reload

Hot reload requires careful coordination:

```go
type PluginManager struct {
    plugins map[string]*LoadedPlugin
    mu      sync.RWMutex
}

type LoadedPlugin struct {
    plugin   shared.Plugin
    path     string
    loadTime time.Time
}

func (pm *PluginManager) Reload(path string) error {
    pm.mu.Lock()
    defer pm.mu.Unlock()

    // Load new version
    newPlugin, err := LoadPlugin(path)
    if err != nil {
        return fmt.Errorf("failed to load new version: %w", err)
    }

    // Replace old version
    name := newPlugin.Name()
    old := pm.plugins[name]
    pm.plugins[name] = &LoadedPlugin{
        plugin:   newPlugin,
        path:     path,
        loadTime: time.Now(),
    }

    // Clean up old version (if any)
    if old != nil {
        log.Printf("Replaced plugin %s: %s -> %s",
            name, old.plugin.Version(), newPlugin.Version())
    }

    return nil
}
```

---

## 4. Complete Solution Walkthrough

### Plugin Interface Design

Our plugin system uses a simple interface:

```go
package shared

// Plugin is the interface that all plugins must implement.
type Plugin interface {
    // Name returns a unique identifier for this plugin
    Name() string

    // Version returns the semantic version of this plugin
    Version() string

    // Init is called once when the plugin is loaded
    // Use this to set up resources, validate configuration, etc.
    Init() error

    // Process is the main function of the plugin
    // It takes arbitrary input and returns arbitrary output
    Process(input interface{}) (interface{}, error)

    // Cleanup is called when the plugin is unloaded
    // Use this to release resources, close connections, etc.
    Cleanup() error
}
```

**Design rationale**:
- `Name()` and `Version()`: Required for tracking and debugging
- `Init()`: Setup hook (allocate resources)
- `Process()`: The actual work (kept generic with `interface{}`)
- `Cleanup()`: Teardown hook (prevent resource leaks)

### Example Plugin: Greeter

```go
package main

import (
    "fmt"
    "strings"
)

var Plugin GreeterPlugin

type GreeterPlugin struct {
    greetingStyle string
}

func (p *GreeterPlugin) Name() string {
    return "greeter"
}

func (p *GreeterPlugin) Version() string {
    return "1.0.0"
}

func (p *GreeterPlugin) Init() error {
    p.greetingStyle = "friendly"
    return nil
}

func (p *GreeterPlugin) Process(input interface{}) (interface{}, error) {
    name, ok := input.(string)
    if !ok {
        return nil, fmt.Errorf("expected string, got %T", input)
    }

    switch p.greetingStyle {
    case "friendly":
        return fmt.Sprintf("Hello, %s! Great to see you!", name), nil
    case "formal":
        return fmt.Sprintf("Good day, %s.", name), nil
    default:
        return fmt.Sprintf("Hi, %s", name), nil
    }
}

func (p *GreeterPlugin) Cleanup() error {
    // Nothing to clean up
    return nil
}
```

### Plugin Manager Implementation

The plugin manager handles loading, tracking, and reloading:

```go
type PluginManager struct {
    plugins   map[string]*LoadedPlugin
    pluginDir string
    mu        sync.RWMutex
    watcher   *fsnotify.Watcher
}

type LoadedPlugin struct {
    plugin    shared.Plugin
    path      string
    loadTime  time.Time
    modTime   time.Time
}

func NewPluginManager(pluginDir string) (*PluginManager, error) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return nil, err
    }

    pm := &PluginManager{
        plugins:   make(map[string]*LoadedPlugin),
        pluginDir: pluginDir,
        watcher:   watcher,
    }

    // Load existing plugins
    if err := pm.LoadAll(); err != nil {
        return nil, err
    }

    // Start watching
    if err := pm.StartWatching(); err != nil {
        return nil, err
    }

    return pm, nil
}
```

### Hot Reload Logic

```go
func (pm *PluginManager) handleFileChange(path string) {
    // Only process .so files
    if !strings.HasSuffix(path, ".so") {
        return
    }

    // Wait a bit to ensure file write is complete
    time.Sleep(100 * time.Millisecond)

    // Check modification time
    info, err := os.Stat(path)
    if err != nil {
        log.Printf("Error stating file %s: %v", path, err)
        return
    }

    pm.mu.RLock()
    for _, loaded := range pm.plugins {
        if loaded.path == path {
            // File hasn't actually changed
            if !info.ModTime().After(loaded.modTime) {
                pm.mu.RUnlock()
                return
            }
            break
        }
    }
    pm.mu.RUnlock()

    // Reload the plugin
    if err := pm.ReloadPlugin(path); err != nil {
        log.Printf("Failed to reload plugin %s: %v", path, err)
        return
    }

    log.Printf("Successfully reloaded plugin: %s", path)
}
```

---

## 5. Key Concepts Explained

### Concept 1: Shared Library Mechanics

**What happens when you load a plugin?**

1. **Dynamic linker** loads the .so file into process memory
2. **Symbol table** is parsed to find exported symbols
3. **Dependencies** are resolved (other .so files, system libs)
4. **Initialization** code runs (init functions, global variables)

**Why this matters**:
- Plugins share the host's memory space (not isolated processes)
- Crashes in plugins can crash the host
- Global state is shared (be careful!)

### Concept 2: Symbol Lookup

```go
sym, err := p.Lookup("MyFunction")
```

This searches the plugin's **export table** for a symbol named "MyFunction".

**What gets exported?**
- Package-level `var` declarations
- Package-level `func` declarations
- Must be in `package main`
- Must start with uppercase letter (Go's export rules)

**Common mistake**:
```go
// ❌ Won't be exported (lowercase)
var plugin MyPlugin

// ✅ Will be exported
var Plugin MyPlugin
```

### Concept 3: Type Identity

**Critical rule**: Types must be **identical** between host and plugin.

**What "identical" means**:
- Same package path
- Same type definition
- Compiled with same Go version
- Same struct field order

**Example**:
```go
// host/shared/types.go
package shared

type Request struct {
    ID   string
    Data string
}

// plugin/plugin.go
import "host/shared"

func Process(req shared.Request) {
    // This works - same type
}
```

**What breaks**:
```go
// plugin/types.go
type Request struct {
    ID   string
    Data string
}

// This is a DIFFERENT type, even though it looks the same!
```

### Concept 4: Plugin Lifecycle

```
Load → Init → Use → Cleanup → (Reload)
 ↓      ↓      ↓       ↓         ↓
Open  Init() Process() Cleanup() Load new version
 .so   setup   work    teardown   replace old
```

**Each stage has responsibilities**:

1. **Load**: Verify plugin file, check dependencies
2. **Init**: Set up resources, validate config
3. **Use**: Handle requests, do work
4. **Cleanup**: Release resources, flush buffers
5. **Reload**: Graceful transition to new version

### Concept 5: File Watching Strategies

**Three approaches**:

1. **Polling** (simple but inefficient):
```go
for {
    time.Sleep(5 * time.Second)
    checkForChanges()
}
```

2. **inotify/FSEvents** (efficient, platform-specific):
```go
watcher, _ := fsnotify.NewWatcher()
watcher.Add(dir)
for event := range watcher.Events {
    if event.Op&fsnotify.Write == fsnotify.Write {
        reload(event.Name)
    }
}
```

3. **Checksum comparison** (reliable, moderate overhead):
```go
lastHash := computeHash(file)
for {
    time.Sleep(1 * time.Second)
    newHash := computeHash(file)
    if newHash != lastHash {
        reload(file)
        lastHash = newHash
    }
}
```

**This project uses polling** for file watching (checks every second for changes). This approach:
- Works on all platforms without external dependencies
- Is simple and reliable
- Has minimal overhead (1-second polling interval)

For production systems, consider using `github.com/fsnotify/fsnotify` for event-based file watching:
```go
watcher, _ := fsnotify.NewWatcher()
watcher.Add(dir)
for event := range watcher.Events {
    if event.Op&fsnotify.Write == fsnotify.Write {
        reload(event.Name)
    }
}
```

---

## 6. Common Patterns You Can Reuse

### Pattern 1: Plugin Registry

```go
type PluginRegistry struct {
    plugins map[string]shared.Plugin
    mu      sync.RWMutex
}

func (r *PluginRegistry) Register(p shared.Plugin) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    name := p.Name()
    if _, exists := r.plugins[name]; exists {
        return fmt.Errorf("plugin %s already registered", name)
    }

    r.plugins[name] = p
    return nil
}

func (r *PluginRegistry) Get(name string) (shared.Plugin, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    p, ok := r.plugins[name]
    return p, ok
}
```

### Pattern 2: Plugin Chain

```go
type PluginChain struct {
    plugins []shared.Plugin
}

func (c *PluginChain) Process(input interface{}) (interface{}, error) {
    result := input
    for _, plugin := range c.plugins {
        var err error
        result, err = plugin.Process(result)
        if err != nil {
            return nil, fmt.Errorf("%s failed: %w", plugin.Name(), err)
        }
    }
    return result, nil
}
```

### Pattern 3: Graceful Reload

```go
func (pm *PluginManager) GracefulReload(name string, newPath string) error {
    // Load new version
    newPlugin, err := LoadPlugin(newPath)
    if err != nil {
        return err
    }

    // Initialize new version
    if err := newPlugin.Init(); err != nil {
        return err
    }

    // Swap atomically
    pm.mu.Lock()
    old := pm.plugins[name]
    pm.plugins[name] = &LoadedPlugin{plugin: newPlugin}
    pm.mu.Unlock()

    // Clean up old version (in background)
    if old != nil {
        go func() {
            if err := old.plugin.Cleanup(); err != nil {
                log.Printf("Cleanup error: %v", err)
            }
        }()
    }

    return nil
}
```

### Pattern 4: Plugin Sandbox (Error Isolation)

```go
func SafeProcess(p shared.Plugin, input interface{}) (output interface{}, err error) {
    // Recover from panics
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("plugin %s panicked: %v", p.Name(), r)
        }
    }()

    // Set timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    done := make(chan struct{})
    go func() {
        output, err = p.Process(input)
        close(done)
    }()

    select {
    case <-done:
        return output, err
    case <-ctx.Done():
        return nil, fmt.Errorf("plugin %s timed out", p.Name())
    }
}
```

### Pattern 5: Plugin Discovery

```go
func DiscoverPlugins(dir string) ([]string, error) {
    var plugins []string

    err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if !info.IsDir() && strings.HasSuffix(path, ".so") {
            plugins = append(plugins, path)
        }

        return nil
    })

    return plugins, err
}
```

---

## 7. Real-World Applications

### Data Processing Pipelines

**Use case**: ETL (Extract, Transform, Load) systems with custom transformations

```
Raw Data → Plugin1 (Clean) → Plugin2 (Enrich) → Plugin3 (Aggregate) → Database
```

Companies using this: Apache Beam, AWS Glue, dbt

### Web Frameworks

**Use case**: Middleware chains with custom handlers

```
Request → Auth Plugin → Logging Plugin → Business Logic → Response
```

Companies using this: Kong API Gateway, Traefik

### Monitoring Systems

**Use case**: Custom metrics collectors

```
Host Agent → CPU Plugin → Memory Plugin → Disk Plugin → Custom Plugin → Time Series DB
```

Companies using this: Telegraf, Prometheus exporters

### Game Engines

**Use case**: Mods and custom game logic

```
Game Engine → Physics Plugin → Rendering Plugin → Custom Mod Plugin → Display
```

Examples: Minecraft mods, Unity plugins

### CI/CD Systems

**Use case**: Custom build steps and deployment strategies

```
Build Pipeline → Test Plugin → Deploy Plugin → Notify Plugin → Done
```

Companies using this: Jenkins, GitLab CI, GitHub Actions

---

## 8. Common Mistakes to Avoid

### Mistake 1: Version Mismatch

**❌ Wrong**:
```bash
# Built host with Go 1.21
go build -o host main.go

# Built plugin with Go 1.20
go build -buildmode=plugin -o plugin.so plugin.go
```

**Error**: `plugin was built with a different version of package`

**✅ Correct**: Use same Go version for everything.

### Mistake 2: Type Redefinition

**❌ Wrong**:
```go
// In both host and plugin, separately:
type Config struct {
    Name string
}
```

**✅ Correct**: Share types via a common package:
```go
// shared/types.go
type Config struct {
    Name string
}

// Import in both host and plugin
```

### Mistake 3: Forgetting to Export

**❌ Wrong**:
```go
var plugin MyPlugin  // lowercase - not exported!
```

**✅ Correct**:
```go
var Plugin MyPlugin  // uppercase - exported
```

### Mistake 4: Not Handling Plugin Panics

**❌ Wrong**:
```go
result := plugin.Process(input)  // If plugin panics, host crashes
```

**✅ Correct**:
```go
defer func() {
    if r := recover(); r != nil {
        log.Printf("Plugin panicked: %v", r)
    }
}()
result := plugin.Process(input)
```

### Mistake 5: File Lock Issues on Reload

**❌ Wrong**:
```go
// Try to overwrite same .so file while it's loaded
os.Remove("plugin.so")  // Fails! File is locked
os.Rename("new.so", "plugin.so")  // Also fails
```

**✅ Correct**: Use versioned filenames:
```go
// Load plugin-v1.so
// Build new version as plugin-v2.so
// Load plugin-v2.so
// Switch to v2
```

### Mistake 6: Shared Global State

**❌ Wrong**:
```go
var counter int  // Shared between host and all plugins - race conditions!

func (p Plugin) Process(input interface{}) {
    counter++  // DATA RACE
}
```

**✅ Correct**: Use plugin-local state:
```go
type Plugin struct {
    counter int  // Each plugin instance has its own
}
```

### Mistake 7: Not Validating Plugin Interface

**❌ Wrong**:
```go
sym, _ := p.Lookup("Plugin")
plugin := sym.(shared.Plugin)  // Panics if type assertion fails
```

**✅ Correct**:
```go
sym, err := p.Lookup("Plugin")
if err != nil {
    return nil, err
}

plugin, ok := sym.(shared.Plugin)
if !ok {
    return nil, fmt.Errorf("invalid plugin type")
}
```

---

## 9. Stretch Goals

### Goal 1: Add Plugin Configuration ⭐

Support loading plugin-specific config from JSON/YAML.

**Hint**:
```go
type Plugin interface {
    Configure(config map[string]interface{}) error
    // ... other methods
}
```

### Goal 2: Implement Plugin Dependencies ⭐⭐

Allow plugins to depend on other plugins.

**Hint**:
```go
type Plugin interface {
    Dependencies() []string  // Names of required plugins
    // ... other methods
}

// Load plugins in dependency order
func (pm *PluginManager) LoadWithDeps() error {
    // Topological sort
}
```

### Goal 3: Add Plugin Marketplace ⭐⭐

Download and install plugins from a remote repository.

**Hint**:
```go
func (pm *PluginManager) Install(name, version string) error {
    // Download from registry
    // Verify signature
    // Extract to plugin dir
    // Load plugin
}
```

### Goal 4: Implement Plugin Versioning ⭐⭐⭐

Support multiple versions of the same plugin running simultaneously.

**Hint**:
```go
type PluginManager struct {
    plugins map[string]map[string]shared.Plugin  // name → version → plugin
}

func (pm *PluginManager) GetVersion(name, version string) shared.Plugin {
    return pm.plugins[name][version]
}
```

### Goal 5: Add RPC-based Plugins ⭐⭐⭐

Support plugins as separate processes (works on Windows!).

**Hint**: Use HashiCorp's go-plugin library:
```go
import "github.com/hashicorp/go-plugin"

// Define gRPC service
type GreeterPlugin struct {
    plugin.Plugin
}
```

---

## 10. Platform Notes

### Linux (Primary Target)

Go plugins work best on Linux:
```bash
# Build plugin
go build -buildmode=plugin -o greeter.so greeter.go

# Run host
./plugin-demo
```

### macOS

Plugins work but with caveats:
```bash
# May need to set DYLD_LIBRARY_PATH
export DYLD_LIBRARY_PATH=./plugins:$DYLD_LIBRARY_PATH
go build -buildmode=plugin -o greeter.so greeter.go
./plugin-demo
```

**Known issue**: Code signing on macOS can interfere with plugin loading.

### Windows

Native plugins don't work. Alternatives:

1. **Use WSL2** (Windows Subsystem for Linux)
2. **Use RPC plugins** (separate processes)
3. **Use embedded scripting** (Lua, JavaScript)

---

## How to Run

```bash
# Build all plugins
make build-plugins

# Run the demo
make run P=47-plugin-system-hot-reload

# Or manually:
cd /home/user/go-edu/minis/47-plugin-system-hot-reload

# Build plugins
go build -buildmode=plugin -o plugins/greeter.so plugins/greeter/greeter.go
go build -buildmode=plugin -o plugins/math.so plugins/math/math.go

# Run demo
go run cmd/plugin-demo/main.go

# In another terminal, modify and rebuild a plugin to see hot reload:
# Edit plugins/greeter/greeter.go
go build -buildmode=plugin -o plugins/greeter.so plugins/greeter/greeter.go
```

### Testing

```bash
# Run tests
go test ./minis/47-plugin-system-hot-reload/...

# Run exercise tests
go test ./minis/47-plugin-system-hot-reload/exercise/

# Run with race detector
go test -race ./minis/47-plugin-system-hot-reload/...
```

---

## Summary

**What you learned**:
- ✅ Go plugins enable dynamic code loading (no recompilation needed)
- ✅ Plugins are .so files built with `-buildmode=plugin`
- ✅ Symbol lookup retrieves exported functions/variables
- ✅ Hot reload requires file watching and careful version management
- ✅ Platform support is limited (Linux/macOS only, no Windows)
- ✅ Type identity is critical (shared types must be identical)

**Why this matters**:
Plugin systems enable extensibility without coupling. This pattern is used in production systems from API gateways to monitoring tools to game engines. Understanding dynamic loading also helps you appreciate how operating systems load shared libraries.

**Next steps**:
- Project 48: Learn about reflection and runtime introspection
- Project 49: Build state machines with dynamic transitions
- Project 50: Combine all concepts into a mini-service

Go forth and extend! 🔌
