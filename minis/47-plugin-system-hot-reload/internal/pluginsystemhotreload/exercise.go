//go:build !solution && !reference

package pluginsystemhotreload

import (
	"fmt"
	"github.com/example/go-10x-minis/minis/47-plugin-system-hot-reload/shared"
	"path/filepath"
	"plugin"
)

// LoadAndExecute implements the exercise.
//
// TODO: Implement this function
func LoadAndExecute(pluginPath string, input interface{}) (interface{}, error) {
	// TODO: Implement
	return nil, nil
}

// DiscoverPlugins implements the exercise.
//
// TODO: Implement this function
func DiscoverPlugins(dir string) ([]string, error) {
	// TODO: Implement
	return nil, nil
}

// ReloadPlugin implements the exercise.
//
// TODO: Implement this function
func ReloadPlugin(pluginPath string) (shared.Plugin, error) {
	// TODO: Implement
	return shared.Plugin{}, nil
}
