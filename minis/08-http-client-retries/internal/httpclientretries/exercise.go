//go:build !solution && !reference

// TODO:
// - Read the tests in exercise_test.go to understand expected behavior.
// - Implement the exported API in this file.
// - Compare with the fully-commented reference in solution.reference.go (go test -tags=reference ./...).
package httpclientretries

import (
	"context"

	"net/http"
	"time"
)

type Client struct {
	HTTP       *http.Client
	MaxRetries int
	BaseDelay  time.Duration
}
// TODO: implement GetJSON.
func GetJSON[T any](ctx context.Context, c *Client, url string) (T, error) { panic("TODO: implement") }
// TODO: implement doRequest.
func doRequest[T any](ctx context.Context, client *http.Client, url string) (T, error) {
	panic("TODO: implement")
}
// TODO: implement isRetryable.
func isRetryable(err error) bool { panic("TODO: implement") }
