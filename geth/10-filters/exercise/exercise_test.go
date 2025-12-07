package exercise

import (
	"context"
	"math/big"
	"testing"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type mockHeadClient struct {
	subHeaders  []*types.Header
	pollHeaders []*types.Header
	pollIndex   int
	subErr      error
	pollErr     error
}

func (m *mockHeadClient) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	if m.subErr != nil {
		return nil, m.subErr
	}
	sub := newMockSubscription()
	go func() {
		for _, h := range m.subHeaders {
			select {
			case <-ctx.Done():
				return
			case ch <- h:
			}
		}
	}()
	return sub, nil
}

func (m *mockHeadClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	if m.pollErr != nil {
		return nil, m.pollErr
	}
	if m.pollIndex >= len(m.pollHeaders) {
		if len(m.pollHeaders) == 0 {
			return nil, nil
		}
		return m.pollHeaders[len(m.pollHeaders)-1], nil
	}
	h := m.pollHeaders[m.pollIndex]
	m.pollIndex++
	return h, nil
}

type mockSubscription struct {
	errCh chan error
}

func newMockSubscription() *mockSubscription {
	return &mockSubscription{errCh: make(chan error)}
}

func (m *mockSubscription) Unsubscribe() {
	close(m.errCh)
}

func (m *mockSubscription) Err() <-chan error {
	return m.errCh
}

func TestRunSubscription(t *testing.T) {
	genesis := common.HexToHash("0x01")
	h1 := makeHeader(1, genesis)
	h2 := makeHeader(2, h1.Hash())
	// Force reorg: parent hash doesn't match previous head hash
	h3 := makeHeader(3, common.HexToHash("0xdead"))

	mock := &mockHeadClient{
		subHeaders: []*types.Header{h1, h2, h3},
	}

	res, err := Run(context.Background(), mock, Config{MaxHeads: 3})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if res.Mode != "subscription" {
		t.Fatalf("expected subscription mode, got %s", res.Mode)
	}
	if len(res.Heads) != 3 {
		t.Fatalf("expected 3 heads, got %d", len(res.Heads))
	}
	if res.Heads[2].Reorg != true {
		t.Fatalf("expected reorg on third head")
	}
}

func TestRunPolling(t *testing.T) {
	genesis := common.HexToHash("0x42")
	h1 := makeHeader(10, genesis)
	h2 := makeHeader(11, h1.Hash())

	mock := &mockHeadClient{
		pollHeaders: []*types.Header{h1, h1, h2},
	}

	cfg := Config{
		MaxHeads:     2,
		PollMode:     true,
		PollInterval: 1 * time.Millisecond,
	}

	res, err := Run(context.Background(), mock, cfg)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if res.Mode != "polling" {
		t.Fatalf("expected polling mode, got %s", res.Mode)
	}
	if len(res.Heads) != 2 {
		t.Fatalf("expected 2 heads, got %d", len(res.Heads))
	}
	if res.Heads[0].Number != 10 || res.Heads[1].Number != 11 {
		t.Fatalf("unexpected head numbers: %+v", res.Heads)
	}
}

func makeHeader(number uint64, parent common.Hash) *types.Header {
	return &types.Header{
		Number:     new(big.Int).SetUint64(number),
		ParentHash: parent,
		Time:       number,
	}
}
