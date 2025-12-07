package exercise

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

type mockStorageClient struct {
	value    []byte
	err      error
	slot     common.Hash
	contract common.Address
	block    *big.Int
}

func (m *mockStorageClient) StorageAt(ctx context.Context, contract common.Address, slot common.Hash, blockNumber *big.Int) ([]byte, error) {
	m.slot = slot
	m.contract = contract
	m.block = blockNumber
	if m.err != nil {
		return nil, m.err
	}
	return m.value, nil
}

func TestRunSimpleSlot(t *testing.T) {
	client := &mockStorageClient{
		value: common.Hex2Bytes("deadbeef"),
	}
	slot := big.NewInt(0)
	contract := common.HexToAddress("0x000000000000000000000000000000000000abcd")

	res, err := Run(context.Background(), client, Config{
		Contract: contract,
		Slot:     slot,
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if client.slot != slotToHash(slot) {
		t.Fatalf("expected slot hash %s, got %s", slotToHash(slot), client.slot)
	}
	if res.ResolvedSlot != slotToHash(slot) {
		t.Fatalf("result slot mismatch")
	}
	if string(res.Value) != string(client.value) {
		t.Fatalf("unexpected value")
	}
}

func TestRunMappingSlot(t *testing.T) {
	client := &mockStorageClient{
		value: common.LeftPadBytes([]byte{1}, 32),
	}
	contract := common.HexToAddress("0x0000000000000000000000000000000000000123")
	key := common.Hex2Bytes("c0ffee")
	slot := big.NewInt(5)

	expectedSlot := mappingSlotHash(key, slotToHash(slot))

	res, err := Run(context.Background(), client, Config{
		Contract:   contract,
		Slot:       slot,
		MappingKey: key,
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if res.ResolvedSlot != expectedSlot {
		t.Fatalf("expected hashed slot %s, got %s", expectedSlot.Hex(), res.ResolvedSlot.Hex())
	}
	if client.slot != expectedSlot {
		t.Fatalf("client StorageAt called with wrong slot")
	}
}

func TestRunErrors(t *testing.T) {
	client := &mockStorageClient{err: errors.New("boom")}
	if _, err := Run(context.Background(), client, Config{
		Contract: common.HexToAddress("0x1"),
		Slot:     big.NewInt(1),
	}); err == nil {
		t.Fatalf("expected storage error")
	}
	if _, err := Run(context.Background(), client, Config{
		Slot: big.NewInt(1),
	}); err == nil {
		t.Fatalf("expected contract validation error")
	}
	if _, err := Run(context.Background(), client, Config{
		Contract: common.HexToAddress("0x1"),
	}); err == nil {
		t.Fatalf("expected slot validation error")
	}
}
