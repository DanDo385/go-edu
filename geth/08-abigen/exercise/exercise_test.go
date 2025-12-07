package exercise

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"testing"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type mockBackend struct {
	abi       abi.ABI
	responses map[string]interface{}
	callErr   error
	lastCalls []ethereum.CallMsg
}

func newMockBackend(t *testing.T) *mockBackend {
	t.Helper()
	parsed, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		t.Fatalf("parse abi: %v", err)
	}
	return &mockBackend{
		abi: parsed,
		responses: map[string]interface{}{
			"name":        "Mock Token",
			"symbol":      "MOCK",
			"decimals":    uint8(18),
			"totalSupply": big.NewInt(42),
			"balanceOf":   big.NewInt(7),
		},
	}
}

func (m *mockBackend) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	m.lastCalls = append(m.lastCalls, call)
	if m.callErr != nil {
		return nil, m.callErr
	}
	if len(call.Data) < 4 {
		return nil, errors.New("data too short")
	}
	method, err := m.abi.MethodById(call.Data[:4])
	if err != nil {
		return nil, err
	}
	value := m.responses[method.Name]
	if method.Name == "balanceOf" && len(call.Data) >= 36 {
		addr := common.BytesToAddress(call.Data[16:36])
		if addr == (common.Address{}) {
			value = big.NewInt(0)
		}
	}
	return method.Outputs.Pack(value)
}

func (m *mockBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return []byte{0x1}, nil
}

func TestRunSuccess(t *testing.T) {
	mock := newMockBackend(t)
	holder := common.HexToAddress("0x1000000000000000000000000000000000000001")
	cfg := Config{
		Contract: common.HexToAddress("0x2000000000000000000000000000000000000002"),
		Holder:   &holder,
	}
	res, err := Run(context.Background(), mock, cfg)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if res.Name != "Mock Token" || res.Symbol != "MOCK" {
		t.Fatalf("unexpected metadata: %+v", res)
	}
	if res.Decimals != 18 {
		t.Fatalf("unexpected decimals %d", res.Decimals)
	}
	if res.TotalSupply.Cmp(big.NewInt(42)) != 0 {
		t.Fatalf("unexpected total supply %s", res.TotalSupply)
	}
	if res.Balance == nil || res.Balance.Cmp(big.NewInt(7)) != 0 {
		t.Fatalf("unexpected balance %v", res.Balance)
	}
	if len(mock.lastCalls) != 5 {
		t.Fatalf("expected 5 calls, got %d", len(mock.lastCalls))
	}
}

func TestRunErrors(t *testing.T) {
	mock := newMockBackend(t)
	mock.callErr = errors.New("boom")
	if _, err := Run(context.Background(), mock, Config{Contract: common.HexToAddress("0x1")}); err == nil {
		t.Fatalf("expected error")
	}
	if _, err := Run(context.Background(), nil, Config{Contract: common.HexToAddress("0x1")}); err == nil {
		t.Fatalf("expected nil backend error")
	}
	if _, err := Run(context.Background(), mock, Config{}); err == nil {
		t.Fatalf("expected empty contract error")
	}
	_, err := Run(context.Background(), mock, Config{Contract: common.HexToAddress("0x1"), ABI: "not json"})
	if err == nil {
		t.Fatalf("expected ABI parse error")
	}
}
