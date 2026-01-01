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

func GetJSON(ctx context.Context, c *Client, url string) (T, error) {
	// TODO: Implement this function
	// TODO: Handle errors appropriately
	// TODO: Add necessary validations
	panic("not implemented")
}

func doRequest(ctx context.Context, client *http.Client, url string) (T, error) {
	// TODO: Implement this function
	panic("not implemented")
}

func isRetryable(err error) bool {
	// TODO: Implement this function
	panic("not implemented")
}
