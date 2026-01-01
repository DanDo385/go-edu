package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/example/go-10x-minis/geth/06-eip1559/internal/eip1559"
)

func main() {
	// geth/06-eip1559
	//
	// Usage:
	//   go run ./geth/06-eip1559/cmd/app <RPC_URL> <private_key_hex> <to_address> <amount_wei> [--send]
	//
	// By default we do NOT broadcast. Pass --send to broadcast.
	//
	// BREAKPOINT: parse args
	if len(os.Args) < 5 {
		fmt.Fprintln(os.Stderr, "usage: <RPC_URL> <private_key_hex> <to_address> <amount_wei> [--send]")
		os.Exit(2)
	}

	rpcURL := os.Args[1]
	pkHex := strings.TrimPrefix(os.Args[2], "0x")
	to := common.HexToAddress(os.Args[3])
	amount, ok := new(big.Int).SetString(os.Args[4], 10)
	if !ok {
		fmt.Fprintln(os.Stderr, "invalid amount_wei")
		os.Exit(2)
	}

	send := len(os.Args) >= 6 && os.Args[5] == "--send"

	pk, err := crypto.HexToECDSA(pkHex)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid private key:", err)
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dial:", err)
		os.Exit(1)
	}
	defer client.Close()

	// BREAKPOINT: run
	out, err := eip1559.Run(ctx, client, eip1559.Config{
		PrivateKey: pk,
		To:         to,
		AmountWei:  amount,
		GasLimit:   21000,
		NoSend:     !send,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "run:", err)
		os.Exit(1)
	}

	fmt.Println("From:", out.FromAddress.Hex())
	fmt.Println("Nonce:", out.Nonce)
	fmt.Println("BaseFee:", out.BaseFee)
	fmt.Println("TxHash:", out.Tx.Hash().Hex())
	fmt.Println("Sent:", send)
}
