//go:build !solution && !reference

package pluginsystemhotreload



import (
	"fmt"
	"path/filepath"
	"plugin"

	"github.com/example/go-10x-minis/minis/47-plugin-system-hot-reload/shared"
)


func LoadAndExecute(pluginPath string, input interface{}) (interface{}, error) {
	// TODO: Implement this function
	panic("unimplemented")
}


func DiscoverPlugins(dir string) ([]string, error) {
	// TODO: Implement this function
	panic("unimplemented")
}


func ReloadPlugin(pluginPath string) (shared.Plugin, error) {
	// TODO: Implement this function
	panic("unimplemented")
}


