//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package pluginsystemhotreload

import "github.com/example/go-10x-minis/minis/47-plugin-system-hot-reload/shared"
// TODO: implement LoadAndExecute.
func LoadAndExecute(pluginPath string, input interface{}) (interface{}, error) {
	panic("TODO: implement")
}
// TODO: implement DiscoverPlugins.
func DiscoverPlugins(dir string) ([]string, error) { panic("TODO: implement") }
// TODO: implement ReloadPlugin.
func ReloadPlugin(pluginPath string) (shared.Plugin, error) { panic("TODO: implement") }
