//go:build !solution && !reference

package workerpoolwordcount

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"unicode"
	"golang.org/x/sync/errgroup"
)

func WordCount(ctx context.Context, urls []string, workers int) (map[string]int, error) {
	// TODO: Implement this function
	// TODO: Handle errors appropriately
	// TODO: Add necessary validations
	panic("not implemented")
}

func WordCountWithErrGroup(ctx context.Context, urls []string, workers int) (map[string]int, error) {
	// TODO: Implement this function
	panic("not implemented")
}

func fetchAndCount(ctx context.Context, url string) (map[string]int, error) {
	// TODO: Implement this function
	panic("not implemented")
}

func tokenizeAndCount(text string) map[string]int {
	// TODO: Implement this function
	panic("not implemented")
}
