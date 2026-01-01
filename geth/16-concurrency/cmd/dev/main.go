package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	// BREAKPOINT: deterministic inputs
	endpoints := []string{"https://eth.llamarpc.com", "https://rpc.ankr.com/eth"}
	workers := 4
	timeout := 5 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	p := jsonRPCProber{hc: &http.Client{Timeout: timeout / 2}}

	out, err := concurrency.Run(ctx, p, concurrency.Config{Endpoints: endpoints, Workers: workers, Timeout: timeout})
	if err != nil {
		fmt.Println("run error:", err)
	}

	fmt.Println("Successes:", len(out.Successes), "Failures:", len(out.Failures))
	for ep, d := range out.Successes {
		fmt.Println("OK", strings.TrimSpace(ep), d)
	}
}
