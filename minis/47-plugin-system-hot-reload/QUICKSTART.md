# Quick Start Guide

## Prerequisites

- Go 1.18 or later
- Linux or macOS (Go plugins not supported on Windows)
- Terminal access

## 1. Build the Plugins

```bash
cd /home/user/go-edu/minis/47-plugin-system-hot-reload

# Option 1: Use the build script
./build-plugins.sh

# Option 2: Use make
make build-plugins

# Option 3: Build manually
go build -buildmode=plugin -o plugins/greeter.so plugins/greeter/greeter.go
go build -buildmode=plugin -o plugins/math.so plugins/math/math.go
go build -buildmode=plugin -o plugins/transformer.so plugins/transformer/transformer.go
```

## 2. Run the Demo

```bash
# Run the interactive plugin manager
make run

# Or directly with go run
go run cmd/plugin-demo/main.go
```

You should see output like:

```
=== Go Plugin System with Hot Reload Demo ===

Plugin directory: ./plugins
Watching for .so file changes...

Loaded plugin greeter v1.0.0 from greeter.so
Loaded plugin math v1.0.0 from math.so
Loaded plugin transformer v1.0.0 from transformer.so

Loaded Plugins (3):
----------------------------------------------------------------------
1. greeter (v1.0.0)
   Path: plugins/greeter.so
   Loaded: 2025-01-15T10:30:00Z

2. math (v1.0.0)
   Path: plugins/math.so
   Loaded: 2025-01-15T10:30:00Z

3. transformer (v1.0.0)
   Path: plugins/transformer.so
   Loaded: 2025-01-15T10:30:00Z
```

## 3. Try the Interactive Commands

```
> list
# Shows all loaded plugins

> run greeter Alice
# Result: Hello, Alice! Great to see you!

> run math {"op":"add","a":10,"b":5}
# Result: 15

> run transformer "hello world"
# Result: HELLO WORLD
```

## 4. Test Hot Reload

### Terminal 1: Run the demo
```bash
make run
```

### Terminal 2: Modify and rebuild a plugin

```bash
# Edit the greeter plugin to change the greeting style
vim plugins/greeter/greeter.go

# Find the Init() function and change:
# p.greetingStyle = "friendly"
# to:
# p.greetingStyle = "enthusiastic"

# Update the version:
# return "1.0.0"
# to:
# return "1.1.0"

# Rebuild
make build-plugins
```

### Terminal 1: Watch the magic!

You should see:
```
Detected change in greeter.so, reloading...
Replaced plugin greeter: 1.0.0 -> 1.1.0
Successfully hot-reloaded greeter.so
```

Now try running the greeter again:
```
> run greeter Alice
# Result: HI ALICE!!! SO HAPPY TO MEET YOU!!!
```

## 5. Run the Exercises

```bash
# Run exercise tests (will fail until implemented)
go test ./exercise/

# Run solution tests
go test -tags solution ./exercise/

# Implement the exercises in exercise/exercise.go
# Then test your implementation
go test -v ./exercise/
```

## 6. Create Your Own Plugin

Create a new plugin directory:
```bash
mkdir -p plugins/mygreeter
```

Create `plugins/mygreeter/mygreeter.go`:
```go
package main

import (
    "fmt"
    "github.com/example/go-10x-minis/minis/47-plugin-system-hot-reload/shared"
)

var Plugin MyGreeterPlugin

type MyGreeterPlugin struct{}

func (p *MyGreeterPlugin) Name() string {
    return "mygreeter"
}

func (p *MyGreeterPlugin) Version() string {
    return "1.0.0"
}

func (p *MyGreeterPlugin) Init() error {
    return nil
}

func (p *MyGreeterPlugin) Process(input interface{}) (interface{}, error) {
    name, ok := input.(string)
    if !ok {
        return nil, fmt.Errorf("expected string")
    }
    return fmt.Sprintf("Greetings, %s!", name), nil
}

func (p *MyGreeterPlugin) Cleanup() error {
    return nil
}
```

Build and test:
```bash
go build -buildmode=plugin -o plugins/mygreeter.so plugins/mygreeter/mygreeter.go

# Run the demo - it will automatically load your new plugin!
make run
```

## Platform-Specific Notes

### Linux
Everything should work out of the box.

### macOS
If you encounter code signing issues:
```bash
# Disable code signing validation (development only)
export CGO_ENABLED=1
export GOFLAGS="-ldflags=-w"
```

### Windows
Go plugins are not supported on Windows. Alternatives:
1. Use WSL2 (Windows Subsystem for Linux)
2. Use a VM with Linux
3. Use RPC-based plugins (e.g., HashiCorp go-plugin)

## Troubleshooting

### Error: "plugin was built with a different version of package"
- Ensure host and plugins are built with the same Go version
- Check: `go version`
- Rebuild everything: `make clean && make build-plugins`

### Error: "plugin: symbol Plugin not found"
- Ensure your plugin exports `var Plugin`
- Variable must start with uppercase `P`
- Must be in `package main`

### Error: "not a Go plugin"
- File is not a valid .so file
- Try rebuilding: `go build -buildmode=plugin ...`

### Hot reload not working
- Check file watcher is running (should see "Watching for .so file changes...")
- Ensure you're rebuilding in the same directory
- Try manually reloading: `> reload greeter.so`

## Next Steps

1. Read the full [README.md](README.md) for deep dive into plugin concepts
2. Complete the exercises in `exercise/exercise.go`
3. Create your own plugins
4. Experiment with hot reload
5. Try building a plugin chain (one plugin's output feeds into another)

Happy plugin hacking! 🔌
