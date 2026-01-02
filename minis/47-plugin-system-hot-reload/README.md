# minis/47-plugin-system-hot-reload

## Problem

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

## Quickstart

```bash
cd minis/47-plugin-system-hot-reload
go test ./...
```

## CLI (cmd/app)

This project includes a small CLI wrapper around the exercise package.

### Flags

- **`-list`**: list available exported functions
- **`-fn`**: function name to run
- **`-in`**: string input (for `func(string) ...`)
- **`-n`**: int input (for `func(int) ...`)
- **`-f`**: float64 input (for `func(float64) ...`)
- **`-b`**: bool input (for `func(bool) ...`)
- **`-file`** / **`-stdin`**: input sources for `func(io.Reader) ...`

### Usage

```bash
go run ./cmd/app -h
```

### Copy/paste examples

```bash
go run ./cmd/app -list
go run ./cmd/app -fn "DiscoverPlugins" -in "Hello, 世界 👋"
go run ./cmd/app -fn "ReloadPlugin" -in "Hello, 世界 👋"
go run ./cmd/dev
```

## Debug harness (cmd/dev)

`cmd/dev` runs the same CLI with preselected example arguments, so you can just hit Run/Debug and see what the project can do.

```bash
go run ./cmd/dev
```

## Files

- `internal/pluginsystemhotreload/exercise.go`: implement the TODOs here
- `internal/pluginsystemhotreload/solution.reference.go`: full reference solution (tagged)
- `cmd/app`: CLI entrypoint
- `cmd/dev`: debug/demo harness
