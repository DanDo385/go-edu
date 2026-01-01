package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/example/go-10x-minis/geth/16-concurrency/internal/concurrency"
)

type jsonRPCProber struct {
	hc *http.Client
}

func (p jsonRPCProber) Probe(ctx context.Context, endpoint string) error {
	payload := []byte(`{"jsonrpc":"2.0","id":1,"method":"eth_chainId","params":[]}`)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.hc.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("non-2xx status: %s", resp.Status)
	}

	var out struct {
		Result string          `json:"result"`
		Error  json.RawMessage `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return err
	}
	if len(out.Error) != 0 {
		return fmt.Errorf("rpc error: %s", string(out.Error))
	}
	if out.Result == "" {
		return fmt.Errorf("missing result")
	}
	return nil
}

func main() {
	// geth/16-concurrency
	//
	// Usage:
	//   go run ./geth/16-concurrency/cmd/app --endpoints <url1,url2,...> [--workers N]
	//
	// BREAKPOINT: parse flags
	endpointsCSV := flag.String("endpoints", "https://eth.llamarpc.com,https://rpc.ankr.com/eth", "comma-separated RPC URLs")
	workers := flag.Int("workers", 4, "worker count")
	timeout := flag.Duration("timeout", 5*time.Second, "overall timeout")
	flag.Parse()

	endpoints := []string{}
	for _, s := range strings.Split(*endpointsCSV, ",") {
		if t := strings.TrimSpace(s); t != "" {
			endpoints = append(endpoints, t)
		}
	}
	if len(endpoints) == 0 {
		fmt.Fprintln(os.Stderr, "no endpoints")
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	p := jsonRPCProber{hc: &http.Client{Timeout: *timeout / 2}}

	// BREAKPOINT: run
	out, err := concurrency.Run(ctx, p, concurrency.Config{Endpoints: endpoints, Workers: *workers, Timeout: *timeout})
	if err != nil {
		fmt.Fprintln(os.Stderr, "run:", err)
	}

	fmt.Println("Successes:", len(out.Successes))
	fmt.Println("Failures:", len(out.Failures))
	for ep, d := range out.Successes {
		fmt.Println("OK", ep, d)
	}
	for ep, e := range out.Failures {
		fmt.Println("ERR", ep, e)
	}
}
