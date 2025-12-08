package exercise

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

type mockTraceClient struct {
	hash common.Hash
	resp json.RawMessage
	err  error
}

func (m *mockTraceClient) TraceTransaction(ctx context.Context, txHash common.Hash) (json.RawMessage, error) {
	m.hash = txHash
	if m.err != nil {
		return nil, m.err
	}
	return m.resp, nil
}

func TestRunSuccess(t *testing.T) {
	payload := json.RawMessage(`{"calls":[{"type":"CALL"}]}`)
	client := &mockTraceClient{resp: payload}
	hash := common.HexToHash("0x1234")

	res, err := Run(context.Background(), client, Config{TxHash: hash})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if res.TxHash != hash {
		t.Fatalf("hash mismatch")
	}
	if !jsonEqual(res.Trace, payload) {
		t.Fatalf("trace payload mismatch")
	}

	// Mutating the returned slice should not affect the stored copy.
	res.Trace[0] = '{'
	if client.resp[0] != '{' {
		t.Fatalf("trace was not copied")
	}
	if client.hash != hash {
		t.Fatalf("client saw wrong hash")
	}
}

func TestRunErrors(t *testing.T) {
	if _, err := Run(context.Background(), nil, Config{TxHash: common.HexToHash("0x1")}); err == nil {
		t.Fatalf("expected nil client error")
	}
	if _, err := Run(context.Background(), &mockTraceClient{}, Config{}); err == nil {
		t.Fatalf("expected missing hash error")
	}
	if _, err := Run(context.Background(), &mockTraceClient{err: errors.New("boom")}, Config{TxHash: common.HexToHash("0x2")}); err == nil {
		t.Fatalf("expected upstream error")
	}
}

func jsonEqual(a, b json.RawMessage) bool {
	var x, y interface{}
	_ = json.Unmarshal(a, &x)
	_ = json.Unmarshal(b, &y)
	ax, _ := json.Marshal(x)
	ay, _ := json.Marshal(y)
	return bytes.Equal(ax, ay)
}
