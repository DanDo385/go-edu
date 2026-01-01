package smartcontracts

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// mockContractCaller simulates contract call responses for testing.
type mockContractCaller struct {
	responses map[string][]byte
	err       error
	calls     []ethereum.CallMsg
}

func (m *mockContractCaller) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
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
		return nil, errors.New("missing response for selector: " + key)
	}
	return resp, nil
}

func TestRunSuccess(t *testing.T) {
	contract := common.HexToAddress("0x1111111111111111111111111111111111111111")

	// Pre-compute selectors
	selectorName := crypto.Keccak256([]byte("name()"))[:4]
	selectorSymbol := crypto.Keccak256([]byte("symbol()"))[:4]
	selectorDecimals := crypto.Keccak256([]byte("decimals()"))[:4]
	selectorTotalSupply := crypto.Keccak256([]byte("totalSupply()"))[:4]

	mock := &mockContractCaller{
		responses: map[string][]byte{
			common.Bytes2Hex(selectorName):        encodeString("Test Token"),
			common.Bytes2Hex(selectorSymbol):      encodeString("TEST"),
			common.Bytes2Hex(selectorDecimals):    encodeUint256(big.NewInt(18)),
			common.Bytes2Hex(selectorTotalSupply): encodeUint256(big.NewInt(1000000)),
		},
	}

	res, err := Run(context.Background(), mock, Config{Contract: contract})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if res.Name != "Test Token" {
		t.Errorf("expected name 'Test Token', got '%s'", res.Name)
	}
	if res.Symbol != "TEST" {
		t.Errorf("expected symbol 'TEST', got '%s'", res.Symbol)
	}
	if res.Decimals != 18 {
		t.Errorf("expected decimals 18, got %d", res.Decimals)
	}
	if res.TotalSupply.Cmp(big.NewInt(1000000)) != 0 {
		t.Errorf("expected totalSupply 1000000, got %s", res.TotalSupply)
	}
	if len(mock.calls) != 4 {
		t.Errorf("expected 4 calls, got %d", len(mock.calls))
	}
}

func TestRunNilClient(t *testing.T) {
	contract := common.HexToAddress("0x1111111111111111111111111111111111111111")
	_, err := Run(context.Background(), nil, Config{Contract: contract})
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestRunEmptyContract(t *testing.T) {
	mock := &mockContractCaller{}
	_, err := Run(context.Background(), mock, Config{})
	if err == nil {
		t.Fatal("expected error for empty contract address")
	}
}

func TestRunCallError(t *testing.T) {
	contract := common.HexToAddress("0x1111111111111111111111111111111111111111")
	mock := &mockContractCaller{err: errors.New("rpc error")}
	_, err := Run(context.Background(), mock, Config{Contract: contract})
	if err == nil {
		t.Fatal("expected error when call fails")
	}
}

func TestDecodeString(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    string
		wantErr bool
	}{
		{
			name:  "valid string",
			input: encodeString("Hello"),
			want:  "Hello",
		},
		{
			name:    "too short",
			input:   make([]byte, 32),
			wantErr: true,
		},
		{
			name:  "empty string",
			input: encodeString(""),
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("decodeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeUint8(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    uint8
		wantErr bool
	}{
		{
			name:  "value 18",
			input: encodeUint256(big.NewInt(18)),
			want:  18,
		},
		{
			name:  "value 6",
			input: encodeUint256(big.NewInt(6)),
			want:  6,
		},
		{
			name:    "too short",
			input:   make([]byte, 16),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeUint8(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeUint8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("decodeUint8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeUint256(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    *big.Int
		wantErr bool
	}{
		{
			name:  "small value",
			input: encodeUint256(big.NewInt(1000)),
			want:  big.NewInt(1000),
		},
		{
			name:  "large value",
			input: encodeUint256(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
			want:  new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil),
		},
		{
			name:    "too short",
			input:   make([]byte, 16),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeUint256(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeUint256() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Cmp(tt.want) != 0 {
				t.Errorf("decodeUint256() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper functions for encoding test data

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
