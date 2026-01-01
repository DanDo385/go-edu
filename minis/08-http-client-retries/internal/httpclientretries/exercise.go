//go:build !solution && !reference

package httpclientretries

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type Client struct {
	HTTP       *http.Client
	MaxRetries int
	BaseDelay  time.Duration
}

// GetJSON implements the exercise.
//
// TODO: Implement this function
func GetJSON[T any](ctx context.Context, c *Client, url string) (T, error) {
	// TODO: Implement
	return *new(T), nil
}
