package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/example/go-10x-minis/minis/08-http-client-retries/exercise"
)

type Response struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}

func main() {
	client := &exercise.Client{
		HTTP:       &http.Client{Timeout: 5 * time.Second},
		MaxRetries: 3,
		BaseDelay:  100 * time.Millisecond,
	}

	ctx := context.Background()

	// Example: Fetch from httpbin (public testing API)
	fmt.Println("Fetching from httpbin.org...")

	var result map[string]interface{}
	result, err := exercise.GetJSON[map[string]interface{}](ctx, client, "https://httpbin.org/get")
	if err != nil {
		log.Printf("Request failed: %v", err)
	} else {
		fmt.Printf("Success: %+v\n", result)
	}
}
