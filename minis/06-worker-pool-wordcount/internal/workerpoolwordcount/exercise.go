//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package workerpoolwordcount

import "context"
// TODO: implement WordCount.
func WordCount(ctx context.Context, urls []string, workers int) (map[string]int, error) {
	panic("TODO: implement")
}
// TODO: implement WordCountWithErrGroup.
func WordCountWithErrGroup(ctx context.Context, urls []string, workers int) (map[string]int, error) {
	panic("TODO: implement")
}
// TODO: implement fetchAndCount.
func fetchAndCount(ctx context.Context, url string) (map[string]int, error) { panic("TODO: implement") }
// TODO: implement tokenizeAndCount.
func tokenizeAndCount(text string) map[string]int { panic("TODO: implement") }
