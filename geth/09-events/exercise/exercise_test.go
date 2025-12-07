package exercise

import (
	"context"
	"errors"
	"math/big"
	"testing"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type mockLogClient struct {
	logs []types.Log
	err  error
	last ethereum.FilterQuery
}

func (m *mockLogClient) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	m.last = q
	if m.err != nil {
		return nil, m.err
	}
	return m.logs, nil
}

func TestRunDecodesLogs(t *testing.T) {
	token := common.HexToAddress("0x000000000000000000000000000000000000dead")
	from := common.HexToAddress("0x0000000000000000000000000000000000000001")
	to := common.HexToAddress("0x0000000000000000000000000000000000000002")

	mock := &mockLogClient{
		logs: []types.Log{
			{
				Address:     token,
				Topics:      []common.Hash{transferSigHash, addressTopic(from), addressTopic(to)},
				Data:        common.LeftPadBytes(big.NewInt(123).Bytes(), 32),
				BlockNumber: 100,
				TxHash:      common.HexToHash("0xabc"),
				Index:       1,
			},
		},
	}

	res, err := Run(context.Background(), mock, Config{
		Token:      token,
		FromHolder: &from,
		ToHolder:   &to,
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if len(res.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(res.Events))
	}
	ev := res.Events[0]
	if ev.From != from || ev.To != to {
		t.Fatalf("unexpected addresses: %+v", ev)
	}
	if ev.Value.Cmp(big.NewInt(123)) != 0 {
		t.Fatalf("unexpected value %s", ev.Value)
	}
	if len(mock.last.Topics) < 1 || mock.last.Topics[0][0] != transferSigHash {
		t.Fatalf("missing transfer topic filter")
	}
}

func TestRunErrors(t *testing.T) {
	token := common.HexToAddress("0xdead")
	mock := &mockLogClient{
		err: errors.New("boom"),
	}
	if _, err := Run(context.Background(), mock, Config{Token: token}); err == nil {
		t.Fatalf("expected error from FilterLogs")
	}
	if _, err := Run(context.Background(), nil, Config{Token: token}); err == nil {
		t.Fatalf("expected nil client error")
	}
	if _, err := Run(context.Background(), mock, Config{}); err == nil {
		t.Fatalf("expected missing token error")
	}
}
