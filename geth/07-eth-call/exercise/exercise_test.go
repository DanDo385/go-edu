package exercise

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

type mockCallClient struct {
	responses map[string][]byte
	err       error
	calls     []ethereum.CallMsg
}

func (m *mockCallClient) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	m.calls = append(m.calls, msg)
	if m.err != nil {
		return nil, m.err
	}
	if msg.Data == nil || len(msg.Data) < 4 {
		return nil, errors.New("invalid selector")
	}
	key := common.Bytes2Hex(msg.Data[:4])
	resp, ok := m.responses[key]
	if !ok {
		return nil, errors.New("missing response")
	}
	return resp, nil
}

func TestRunSuccess(t *testing.T) {
	contract := common.HexToAddress("0x1111111111111111111111111111111111111111")
	mock := &mockCallClient{
		responses: map[string][]byte{
			common.Bytes2Hex(selectorName):        encodeString("Dai Stablecoin"),
			common.Bytes2Hex(selectorSymbol):      encodeString("DAI"),
			common.Bytes2Hex(selectorDecimals):    encodeUint256(big.NewInt(18)),
			common.Bytes2Hex(selectorTotalSupply): encodeUint256(big.NewInt(1000)),
		},
	}

	res, err := Run(context.Background(), mock, Config{Contract: contract})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if res.Name != "Dai Stablecoin" || res.Symbol != "DAI" {
		t.Fatalf("unexpected metadata: %+v", res)
	}
	if res.Decimals != 18 {
		t.Fatalf("unexpected decimals %d", res.Decimals)
	}
	if res.TotalSupply.Cmp(big.NewInt(1000)) != 0 {
		t.Fatalf("unexpected total supply %s", res.TotalSupply)
	}
	if len(mock.calls) != 4 {
		t.Fatalf("expected 4 calls, got %d", len(mock.calls))
	}
	for _, call := range mock.calls {
		if call.To == nil || *call.To != contract {
			t.Fatalf("call not directed to contract")
		}
	}
}

func TestRunErrors(t *testing.T) {
	contract := common.HexToAddress("0x1")
	mock := &mockCallClient{err: errors.New("rpc")}
	if _, err := Run(context.Background(), mock, Config{Contract: contract}); err == nil {
		t.Fatalf("expected error")
	}
	if _, err := Run(context.Background(), nil, Config{Contract: contract}); err == nil {
		t.Fatalf("expected nil client error")
	}
	if _, err := Run(context.Background(), mock, Config{}); err == nil {
		t.Fatalf("expected empty contract error")
	}
}

func encodeString(s string) []byte {
	strBytes := []byte(s)
	length := len(strBytes)
	paddedLen := ((length + 31) / 32) * 32
	out := make([]byte, 64+paddedLen)
	putUint(out[0:32], big.NewInt(32))             // offset
	putUint(out[32:64], big.NewInt(int64(length))) // length
	copy(out[64:64+length], strBytes)
	return out
}

func encodeUint256(v *big.Int) []byte {
	out := make([]byte, 32)
	putUint(out, v)
	return out
}

func putUint(dst []byte, v *big.Int) {
	for i := range dst {
		dst[i] = 0
	}
	b := v.Bytes()
	copy(dst[32-len(b):], b)
}
