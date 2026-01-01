//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package toolbox

import "context"
// TODO: implement Run.
func Run(ctx context.Context, client ToolboxClient, cfg Config) (*Result, error) {
	panic("TODO: implement")
}
// TODO: implement handleStatus.
func handleStatus(ctx context.Context, client ToolboxClient) (*Result, error) {
	panic("TODO: implement")
}
// TODO: implement handleBlock.
func handleBlock(ctx context.Context, client ToolboxClient, args []string) (*Result, error) {
	panic("TODO: implement")
}
// TODO: implement handleTx.
func handleTx(ctx context.Context, client ToolboxClient, args []string) (*Result, error) {
	panic("TODO: implement")
}
