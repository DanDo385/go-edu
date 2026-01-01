//go:build !solution && !reference

package workerpoolwordcount

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"strings"
	"sync"
	"unicode"
)

// WordCount implements the exercise.
//
// TODO: Implement this function
func WordCount(ctx context.Context, urls []string, workers int) (map[string]int, error) {
	// TODO: Implement
	return nil, nil
}

// WordCountWithErrGroup implements the exercise.
//
// TODO: Implement this function
func WordCountWithErrGroup(ctx context.Context, urls []string, workers int) (map[string]int, error) {
	// TODO: Implement
	return nil, nil
}
